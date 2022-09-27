package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"text/template"
	"time"
)

var commonLabels = []string{"alertname", "cluster", "instance", "namespace"}
var excludeLabels = []string{"unsilence", "severity"}

func callback(content *gin.Context) {
	cardUpdate := UpdatingCard{}
	content.BindJSON(&cardUpdate)
	log.Printf("%+v\n", &cardUpdate)

	if !IsBlank(cardUpdate.Challenge) {
		content.JSON(http.StatusOK, map[string]string{
			"challenge": cardUpdate.Challenge,
		})
		return
	}

	if v, ok := cardUpdate.Action.Value["unsilence"]; !ok || v == "0" {
		silence(cardUpdate.OpenMessageId, cardUpdate.Action.Value)
	} else {
		unsilence(cardUpdate.OpenMessageId, cardUpdate.Action.Value)
	}

}

func silence(messageId string, value map[string]string) {
	var matchers = make([]Matcher, 0)
	for key, value := range value {
		matchers = append(matchers, Matcher{
			Name:    key,
			Value:   value,
			IsEqual: true,
			IsRegex: !Contains(commonLabels, key),
		})
	}

	silenceResponse := createSilence(matchers)
	updateMessage(UpdatingCardNecessity{
		MessageId: messageId,
		SilenceId: silenceResponse.SilenceId,
		Labels:    value,
	})
}

func unsilence(messageId string, value map[string]string) {
	beExpiredSilence(value["silenceId"])
	updateMessage(UpdatingCardNecessity{MessageId: messageId, Labels: value})
}

func createSilence(matchers []Matcher) SilenceResponse {
	silencesEntity, err := json.Marshal(SilenceEntity{
		Matchers:  matchers,
		CreatedBy: "Admin",
		Comment:   "飞书静音",
		StartsAt:  time.Now().Format(time.RFC3339),
		EndsAt:    time.Now().Add(4 * time.Hour).Format(time.RFC3339),
	})
	if err != nil {
		log.Fatal(err)
	}

	var silenceResponse SilenceResponse
	headers := map[string]interface{}{
		"Content-Type": "application/json",
	}
	bys := BuildPostRequestAndDo(bytes.NewBuffer(silencesEntity), cfg.AlertManager, headers)
	if err := json.Unmarshal(bys, &silenceResponse); err != nil {
		log.Fatal(err)
	}
	return silenceResponse
}

func beExpiredSilence(silenceId string) {
	bys := BuildRequestAndDo("DELETE", bytes.NewBuffer([]byte{}), cfg.AlertManagerSingleSilence+"/"+silenceId, nil)
	log.Printf("be Expired silence: %s", string(bys))
}

func updateMessage(necessity UpdatingCardNecessity) {
	templateEntity := FeishuMessageTemplateEntity{
		AlertName: necessity.Labels["alertname"],
		Datetime:  time.Now().Format("2006-01-02 15:04:05"),
		Cluster:   necessity.Labels["cluster"],
		Instance:  necessity.Labels["instance"],
		Namespace: necessity.Labels["namespace"],
		Recreated: IsBlank(necessity.SilenceId),
	}
	for k, v := range necessity.Labels {
		if Contains(excludeLabels, k) {
			continue
		}
		templateEntity.LabelsString = templateEntity.LabelsString + ",\"" + k + "\":\"" + v + "\""
	}
	if !IsBlank(necessity.SilenceId) {
		templateEntity.LabelsString = templateEntity.LabelsString[1:] + ",\"silenceId\":\"" + necessity.SilenceId + "\""
	}

	template_, _ := template.ParseFiles("template/feishu-silenced-message.txt")

	var tpl bytes.Buffer
	eerr := template_.Execute(&tpl, templateEntity)
	if eerr != nil {
		log.Fatal(eerr)
	}
	log.Printf("Update message: " + tpl.String() + "\n")
	sendMessage(tpl.String(), true, necessity.MessageId)

}
