package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func simpleHTTPHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is a test endpoint"))
}

// ToDo: add function to accept a file in the easier format of text/plain

// ToDo: re-add path parameter
// func (s *server) createTaskListBatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
func (s *server) createTaskListBatch(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Panic(err)
		return
	}
	defer file.Close()

	_, fileName := filepath.Split(handler.Filename)

	filePath := "./uploads/" + fileName
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	// send file to the kafka producer
	_, err = s.kfk.ProduceFromFile(filePath)
	if err != nil {
		log.Panic(err) // ToDo: return error to the user
	}

	// and clean up the file
	err = os.Remove(filePath)
	if err != nil {
		log.Panic(err) // ToDo: return error to the user
	}

}
