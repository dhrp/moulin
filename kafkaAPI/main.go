package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome friends.\n")
}

func getTaskLists(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if id != "" {
		fmt.Fprintf(w, "No task lists with %s found (yet)!", id)
	} else {
		fmt.Fprintf(w, "No lists found!")
	}
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
	produceFromFile(filePath)
}

// readFileLines is a stub for reading the uploaded file
// and then feeding it into the Kafka Producer
func produceFromFile(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var line string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		line = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// currently the last line
	return line
}

func main() {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/v1/task_lists/", getTaskLists)
	router.GET("/v1/task_lists/:id/", getTaskLists)
	router.POST("/v1/task_list/batch/", createTaskListBatch)

	serverAddr := "127.0.0.1:8080"

	log.Println("starting server on " + serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
