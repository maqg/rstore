package blobupload

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/utils"
	"os"
)

// BlobUpload base structure
type BlobUpload struct {
	ID         string
	FilePath   string
	RespWriter http.ResponseWriter
	Request    *http.Request
}

// CreateDir for blob
func (bu *BlobUpload) CreateDir() {
	dir := blobs.DirPath(bu.ID)
	utils.CreateDir(dir)
}

// Upload for blob
func (bu *BlobUpload) Upload() error {

	// create Dir for blob
	bu.CreateDir()

	err := CopyFullPayload(bu.RespWriter, bu.Request, bu.FilePath)
	if err != nil {
		fmt.Printf("copy full data for blob %s error\n", bu.ID)
		return nil
	}

	return nil
}

// StartBlobUpload for blob start upload
func StartBlobUpload(name string, digest string) (interface{}, error) {
	return nil, nil
}

// CopyFullPayload will copy all data from r.Body and write them to destWrite
func CopyFullPayload(responseWriter http.ResponseWriter, r *http.Request, filepath string) error {

	// Get a channel that tells us if the client disconnects
	var clientClosed <-chan bool
	if notifier, ok := responseWriter.(http.CloseNotifier); ok {
		clientClosed = notifier.CloseNotify()
	} else {
		fmt.Printf("the ResponseWriter does not implement CloseNotifier (type: %T)", responseWriter)
	}

	// Read in the data, if any.
	destWriter, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("create file of %s error\n", filepath)
		return err
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read all data from r.Body error\n")
		return err
	}

	len, err := destWriter.Write(data)
	if err != nil {
		fmt.Printf("Write data to dest writer error\n")
		return nil
	}

	copied := int64(len)
	if clientClosed != nil && (err != nil || (r.ContentLength > 0 && copied < r.ContentLength)) {
		// Didn't receive as much content as expected. Did the client
		// disconnect during the request? If so, avoid returning a 400
		// error to keep the logs cleaner.
		select {
		case <-clientClosed:
			// Set the response code to "499 Client Closed Request"
			// Even though the connection has already been closed,
			// this causes the logger to pick up a 499 error
			// instead of showing 0 for the HTTP status.
			responseWriter.WriteHeader(499)
			return errors.New("client disconnected")
		default:
		}
	}

	if err != nil {
		fmt.Printf("error got, copy data failed")
		return err
	}

	return nil
}
