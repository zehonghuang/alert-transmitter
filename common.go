package main

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"net/http"
)

func BuildRequestAndDo(method string, req *bytes.Buffer, url string, headers map[string]interface{}) []byte {
	request, err := http.NewRequest(method, url, req)
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

func BuildPostRequestAndDo(req *bytes.Buffer, url string, headers map[string]interface{}) []byte {
	return BuildRequestAndDo("POST", req, url, headers)
}

func HashCode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func HashCodes(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", HashCode(buf.String()))
}

func Contains(s []string, ele string) bool {
	for _, v := range s {
		if v == ele {
			return true
		}
	}
	return false
}
