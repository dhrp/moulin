package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nerdalize/moulin/rouge"
)

// type APIServer struct {
// 	// ToDo: Learn about if this should be a reference
// 	client rouge.RougeClient
// }

// Main RougeClient instance
var rougeClient *rouge.RougeClient

func getTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queueID := ps.ByName("queue")

	msg := rougeClient.Load(queueID, 300)
	mstr := msg.ToString()
	w.Write([]byte(mstr))
}

// # HTTP api of Moulin
func main() {

	rougeClient = &rouge.RougeClient{Host: "localhost:6379"}
	rougeClient.Init()
	router := httprouter.New()

	router.GET("/v1/get_task/:queue/", getTask)

	serverAddr := "127.0.0.1:8042"

	log.Println("starting server on " + serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
