package handlers

import (
	"octlink/rstore/utils/configuration"
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"octlink/rstore/modules/config"
	"octlink/rstore/modules/image"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// GetManifest for api call
func getManifest(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := mux.Vars(r)["digest"]

	if name == "" || digest == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Errorf("got manifest error for bad paras, no name or digest specified\n")
		return
	}

	manifest := manifest.GetManifest(name, digest)
	if manifest == nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Warnf("manifest of %s not exist\n", digest)
		return
	}

	data := utils.JSON2String(manifest)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, data)
}

// DeleteManifest for api call
func deleteManifest(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	digest := mux.Vars(r)["digest"]

	if name == "" || digest == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Errorf("delete manifest of %s error, no name or digest specified\n", digest)
		return
	}

	manifest := manifest.GetManifest(name, digest)
	if manifest == nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Warnf("manifest of %s not exist\n", digest)
		return
	}

	err := manifest.Delete()
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		logger.Warnf("delete manifest of %s error %s\n", digest, err)
		return
	}

	logger.Debugf("delete of manifest %s OK\n", digest)

	w.WriteHeader(http.StatusOK)
}

func removeTempFile(tempName string) {
	if tempName != "" {
		os.Remove(configuration.RootDirectory() + manifest.TempDir + "/" + tempName)
	}
}

func postManifest(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := mux.Vars(r)["digest"]
	tempName := r.FormValue("tempName")	

	// bad args for manifest post
	if name == "" || digest == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Errorf("post manifest error for bad paras, name or digest not specified\n")
		return
	}

	logger.Debugf("to post manifest %s:%s:%s\n", name, digest, tempName)

	m := manifest.GetManifest(name, digest)
	if m != nil {
		err := image.UpdateImageCallback(m.Name, m.DiskSize, m.VirtualSize,
			m.BlobSum, config.ImageStatusReady)
		if err != nil {
			logger.Warnf("update image info %s error, and manifest created OK\n", m.Name)
		}

		// remove tempfile for this manifest
		removeTempFile(tempName)

		w.WriteHeader(http.StatusOK)

		logger.Debugf("manifest of %s already exist, updated it\n", digest)

		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Errorf("readall data from http body error %s\n", err)
		// remove tempfile for this manifest
		removeTempFile(tempName)
		return
	}

	m = new(manifest.Manifest)
	if err = json.Unmarshal(data, m); err != nil {
		logger.Errorf("convert data to json error %s\n", err)
		
		// remove tempfile for this manifest
		removeTempFile(tempName)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = m.Write(); err != nil {
		logger.Errorf("error happend for manifest write %s\n", err)
		// remove tempfile for this manifest
		removeTempFile(tempName)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = image.UpdateImageCallback(m.Name, m.DiskSize, m.VirtualSize,
		m.BlobSum, config.ImageStatusReady)
	if err != nil {
		logger.Infof("update image info %s error, and manifest created OK\n", m.Name)
	}

	if tempName != "" {
		imageFilePath := configuration.RootDirectory() + manifest.TempDir + "/" + tempName
		destFilePath := configuration.RootDirectory() + manifest.ManifestDir + "/" + m.BlobSum
		if utils.IsFileExist(imageFilePath) {
			utils.CreateDir(destFilePath)
			os.Rename(imageFilePath, destFilePath)
		}
	}

	w.WriteHeader(http.StatusOK)

	logger.Debugf("Post manifest of %s OK!\n", digest)
}

func manifestManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":    http.HandlerFunc(getManifest),
		"POST":   http.HandlerFunc(postManifest),
		"DELETE": http.HandlerFunc(deleteManifest),
	}

	return mhandler
}
