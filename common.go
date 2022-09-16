package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func buildRequestAndDo(req *bytes.Buffer, url string, headers map[string]interface{}) []byte {
	request, err := http.NewRequest("POST", url, req)
	if err != nil {
	}
	if headers != nil {
		for k, v := range headers {
			request.Header.Set(k, v.(string))
		}
	}
	response, err := _client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return body
}
