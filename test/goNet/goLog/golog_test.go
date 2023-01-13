package goLog

import (
	"log"
	"net/http"
	"os"
	"testing"
)

func TestGoLog(t *testing.T) {
	/**
	logFile, err := os.OpenFile("goLog.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	if err != nil {
		panic(err.Error())
	}
	log.SetOutput(logFile)
	*/
	log.SetOutput(os.Stdout)
	simpleHttpGet := func(url string) {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error fetching url %s : %s", url, err.Error())
		} else {
			log.Printf("Status Code for %s : %s", url, resp.Status)
			resp.Body.Close()
		}
	}
	simpleHttpGet("www.google.com")
	simpleHttpGet("http://www.google.com")
}
