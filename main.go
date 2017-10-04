package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nerdalize/moulin/rouge"
)

// Main RougeClient instance
var rougeClient *rouge.RougeClient

func loadTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queueID := ps.ByName("queue")

	msg := rougeClient.Load(queueID, 300)
	mstr := msg.ToString()
	w.Write([]byte(mstr))
}

func pushTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if r.Header.Get("Content-type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	queueID := ps.ByName("queue")

	jsonPayload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Bad request: could not read body.", http.StatusBadRequest)
		return
	}

	task, err := createTaskFromJSON(jsonPayload)
	if err != nil {
		http.Error(w, "Bad request: json needs to have at least 'body' field.", http.StatusBadRequest)
		return
	}

	result := fmt.Sprintf("Task with body: %v, set ID %v \n", task.Body, task.ID)
	log.Println(result)
	fmt.Fprintf(w, result)

	length := rougeClient.AddTask(queueID, *task)
	fmt.Fprintf(w, "New length is %v", length)
}

// # HTTP api of Moulin
func main() {

	rougeClient = &rouge.RougeClient{Host: "localhost:6379"}
	rougeClient.Init()

	router := httprouter.New()
	router.GET("/v1/queue/:queue/load/", loadTask)
	router.POST("/v1/queue/:queue/push/", pushTask)

	serverAddr := "127.0.0.1:8042"
	log.Println("starting server on " + serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
