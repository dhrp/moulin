package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/suite"
)

type TaskAPITestSuite struct {
	suite.Suite
}

func (suite *TaskAPITestSuite) SetupSuite() {

}

func (suite *TaskAPITestSuite) TestUploadTaskBatch() error {
	filename := "test/testtextfile.txt"
	targetURL := "http://testserver.com/v1/task_list/batch/"

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

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
	ps := httprouter.Params{}
	createTaskListBatch(res, req, ps)

	return nil

}

func (suite *TaskAPITestSuite) TestProduceFromFile() {
	lastLine := produceFromFile("./test/testtextfile.txt")
	suite.Equal("line six", lastLine, "line from producer not what was expected")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTaskAPITestSuite(t *testing.T) {
	suite.Run(t, new(TaskAPITestSuite))
}
