package handlers

import (
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
		logger.Warnf("delete manifest of %s error %s\n",
			digest, err)
		return
	}

	logger.Debugf("delete of manifest %s OK\n", digest)

	w.WriteHeader(http.StatusOK)
}

func postManifest(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := mux.Vars(r)["digest"]

	if name == "" || digest == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Debugf("got manifest post request %s:%s\n", name, digest)

	m := manifest.GetManifest(name, digest)
	if m != nil {
		err := image.UpdateImageCallback(m.Name, m.DiskSize, m.VirtualSize,
			m.BlobSum, config.ImageStatusReady)
		if err != nil {
			logger.Warnf("update image info %s error, and manifest created OK\n", m.Name)
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Errorf("readall data from http body error %s\n", err)
		return
	}

	m = new(manifest.Manifest)
	if err = json.Unmarshal(data, m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Errorf("convert data to json error %s\n", err)
		return
	}

	if err = m.Write(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Errorf("error happend for manifest write %s\n", err)
		return
	}

	err = image.UpdateImageCallback(m.Name, m.DiskSize, m.VirtualSize,
		m.BlobSum, config.ImageStatusReady)
	if err != nil {
		logger.Infof("update image info %s error, and manifest created OK\n", m.Name)
	}

	w.WriteHeader(http.StatusOK)
}

func manifestManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":    http.HandlerFunc(getManifest),
		"POST":   http.HandlerFunc(postManifest),
		"DELETE": http.HandlerFunc(deleteManifest),
	}

	return mhandler
}
