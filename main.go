package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Response struct {
	FoundAtSite string
}

type Request struct {
	Site       []string
	SearchText string
}

func hasTextAtURL(url, text string) bool {
	response, err := http.Get(url)

	if err != nil {
		log.Println("Error getting page. Details: ", err)
		return false
	}

	pageData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println("Error reading page. Details: ", err)
		return false
	}

	return strings.Contains(string(pageData), text)
}

func checkTextHandler(c *gin.Context) {
	var request Request
	err := c.BindJSON(&request)
	if err != nil || len(request.SearchText) == 0 || len(request.Site) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Invalid JSON"})
		return
	}

	sitesCount := len(request.Site)
	results := make(chan string, sitesCount)

	var waitGroup sync.WaitGroup
	waitGroup.Add(sitesCount)

	for _, url := range request.Site {
		go func(url string) {
			defer waitGroup.Done()
			if hasTextAtURL(url, request.SearchText) {
				results <- url
			}
		}(url)
	}

	waitGroup.Wait()

	select {
	case url := <-results:
		c.JSON(http.StatusOK, Response{url})
	default:
		c.Writer.WriteHeader(204) // No content
	}
}

func main() {
	router := gin.Default()
	router.POST("/checkText", checkTextHandler)
	router.Run(":8080")
}
