package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

func simpleHTTPHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is a test endpoint"))
}

func createTaskListBatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	// here we take a next action with the uploaded file
	// produceFromFile(filePath)
}
