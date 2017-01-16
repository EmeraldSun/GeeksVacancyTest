package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//Look for 'oo' in 'foo' and 'bar'. Expected result - localhost/foo
func TestFindTextSuccess(t *testing.T) {
	var gf GinFixture
	gf.init()
	gf.addDummyHandler("/foo", "foo")
	gf.addDummyHandler("/bar", "bar")

	server := httptest.NewServer(gf.router)
	defer server.Close()

	sites := []string{server.URL + "/foo", server.URL + "/bar"}

	err := gf.testRequest(Request{sites, "oo"}, http.StatusOK, sites[0])
	if err != nil {
		t.Error(err)
	}
}

//Look for 'aa' in 'foo' and 'bar'. Expected result - No Content
func TestFindTextFail(t *testing.T) {
	var gf GinFixture
	gf.init()
	gf.addDummyHandler("/foo", "foo")
	gf.addDummyHandler("/bar", "bar")

	server := httptest.NewServer(gf.router)
	defer server.Close()

	sites := []string{server.URL + "/foo", server.URL + "/bar"}

	err := gf.testRequest(Request{sites, "aa"}, http.StatusNoContent, "")

	if err != nil {
		t.Error(err)
	}
}

//Look for ' ' in empty link. Expected result - No Content
func TestFindTextInEmptyLink(t *testing.T) {
	var gf GinFixture
	gf.init()

	sites := []string{""}

	err := gf.testRequest(Request{sites, " "}, http.StatusNoContent, "")

	if err != nil {
		t.Error(err)
	}
}

//Pass random values instead of JSON. Expected result - BadRequest
func TestIncorrectJson(t *testing.T) {
	var gf GinFixture
	gf.init()

	randomMap := map[string]float32{"apple": 5.1, "lettuce": 7.11, "pear": 8.77}

	err := gf.testRequest(randomMap, http.StatusBadRequest, "")

	if err != nil {
		t.Error(err)
	}
}
