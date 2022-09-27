package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func sendMessage(message string, update bool, messageId string) bool {
	var result SendingMessageResponse

	headers := map[string]interface{}{
		"Authorization": "Bearer " + refreshToken(),
		"Content-Type":  "application/json",
	}

	if !update {
		request, err := json.Marshal(map[string]string{
			"receive_id": cfg.Feishu.ChatId,
			"content":    message,
			"msg_type":   "interactive",
		})
		if err != nil {
			log.Fatal(err)
		}

		response := BuildPostRequestAndDo(bytes.NewBuffer(request), cfg.Feishu.ReceiveMessageUrl, headers)
		log.Printf("%s", response)
		if err := json.Unmarshal(response, &result); err != nil {
			log.Fatal(err)
		}

	} else {
		if IsBlank(messageId) {
			log.Fatal("MessageId should be to have value if update is true.")
		}
		request, err := json.Marshal(map[string]string{
			"content": message,
		})
		if err != nil {
			log.Fatal(err)
		}
		response := BuildRequestAndDo("PATCH", bytes.NewBuffer(request), fmt.Sprintf(cfg.Feishu.UpdateMessageUrl, messageId), headers)
		log.Printf("Update message response: %s", string(response))
	}

	return result.Code == 0
}

func refreshToken() string {

	token, found := _cache.Get("token")
	if found {
		return token.(string)
	}
	request, err := json.Marshal(map[string]string{
		"app_id":     cfg.Feishu.AppId,
		"app_secret": cfg.Feishu.AppSecret,
	})
	if err != nil {
		log.Fatal("build request body failed.")
	}
	response, err := _client.Post(cfg.Feishu.AppAccessTokenUrl, "application/json", bytes.NewBuffer(request))
	if err != nil {
		log.Fatal("Http client established connection failed.")
	}
	var result AppAccessTokenEntity
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Fatal("read response body failed.")
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", &result)
	token = result.AppAccessToken
	_cache.Set("token", result.AppAccessToken, 2*time.Hour)
	return token.(string)
}
