package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
)

func (suite *MainTestSuite) TestUploadTaskBatch() {
	log.Println("*** testing TestUploadTaskBatch")

	suite.NotNil(suite.server.rouge, "rouge not initialized")

	filename := "./test/testtextfile.txt"
	targetURL := "http://testserver.com/v1/task_list/batch/"

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		log.Panic("error writing to buffer")
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		log.Panic("error opening test file")
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	suite.Nil(err, "io.Copy (for creating mock file upload) failed")

	contentType := bodyWriter.FormDataContentType()
	_ = contentType
	bodyWriter.Close()

	// resp, err := http.Post(targetURL, contentType, bodyBuf)
	req, err := http.NewRequest("POST", targetURL, bodyBuf)
	req.Header.Add("Content-Type", contentType)
	if err != nil {
		suite.FailNow(err.Error(), "failed to make request")
	}
	res := httptest.NewRecorder()
	// ps := httprouter.Params{}
	suite.server.createTaskListBatch(res, req) // , ps

}
