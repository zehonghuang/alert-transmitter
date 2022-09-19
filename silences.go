package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"text/template"
	"time"
)

func silences(content *gin.Context) {
	cardUpdate := UpdatingCard{}
	content.BindJSON(&cardUpdate)
	log.Printf("%+v\n", &cardUpdate)

	var matchers = make([]Matcher, 0)
	for key, value := range cardUpdate.Action.Value {
		matchers = append(matchers, Matcher{
			Name:    key,
			Value:   value,
			IsEqual: true,
			IsRegex: false,
		})
	}
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

	bys := buildRequestAndDo(bytes.NewBuffer(silencesEntity), cfg.AlertManager, nil)
	var silenceResponse SilenceResponse
	if err := json.Unmarshal(bys, &silenceResponse); err != nil {
		log.Fatal(err)
	}

	updateMessage(UpdatingCardNecessity{
		MessageId: cardUpdate.OpenMessageId,
		SilenceId: silenceResponse.SilenceId,
		Labels:    cardUpdate.Action.Value,
	})
}

func updateMessage(necessity UpdatingCardNecessity) bool {
	templateEntity := FeishuMessageTemplateEntity{
		AlertName: necessity.Labels["alertname"],
		Datetime:  time.Now().Format("2006-01-02 15:04:05"),
		Cluster:   necessity.Labels["cluster"],
		Instance:  necessity.Labels["instance"],
		Namespace: necessity.Labels["namespace"],
	}

	template_, _ := template.ParseFiles("template/feishu-silenced-message.txt")

	var tpl bytes.Buffer
	eerr := template_.Execute(&tpl, templateEntity)
	if eerr != nil {
		log.Fatal(eerr)
	}
	log.Printf(tpl.String() + "\n")
	sendMessage(tpl.String(), true, necessity.MessageId)

	return false
}
