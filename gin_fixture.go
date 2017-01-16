package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type GinFixture struct {
	router *gin.Engine
}

func (gf *GinFixture) addDummyHandler(path, content string) {
	gf.router.GET(path, func(c *gin.Context) {
		c.String(http.StatusOK, content)
	})
}

func (gf *GinFixture) init() {
	gin.SetMode(gin.TestMode)
	gf.router = gin.New()
	gf.router.POST("/checkText", checkTextHandler)
}

func (gf *GinFixture) testRequest(request interface{}, expectedCode int, expectedURL string) error {
	jsonRequest, _ := json.Marshal(request)
	httpRequest, _ := http.NewRequest("POST", "/checkText", bytes.NewBuffer(jsonRequest))
	httpRequest.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	gf.router.ServeHTTP(recorder, httpRequest)

	if recorder.Code != expectedCode {
		msg := fmt.Sprintf("Error: expected code %d, received: %d", expectedCode, recorder.Code)
		return errors.New(msg)
	}

	if recorder.Code != http.StatusOK {
		return nil
	}

	respBody, _ := ioutil.ReadAll(recorder.Body)
	var response Response
	err := json.Unmarshal(respBody, &response)

	if err != nil {
		return errors.New(fmt.Sprintf("Invalid response. Details: %q", err.Error()))
	}

	if response != (Response{expectedURL}) {
		msg := fmt.Sprintf("Error: expected url %q, received: %q", expectedURL, response.FoundAtSite)
		return errors.New(msg)
	}

	return nil
}
