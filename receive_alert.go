package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"text/template"
	"time"
)

func receiveAlert(context *gin.Context) {
	//token := refreshToken()
	request := AlertMessage{}
	err := context.BindJSON(&request)
	if err != nil {
		log.Fatal("bind json failed.")
	}
	assembleMessageCard(request)
	context.AsciiJSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

func counter() func() int {
	i := -1
	return func() int {
		i++
		return i
	}
}

func assembleMessageCard(alertmessage AlertMessage) {
	for _, alert := range alertmessage.Alerts {
		message := FeishuMessageTemplateEntity{
			Datetime:  time.Now().Format("2006-01-02 15:04:05"),
			AlertName: alert.Labels["alertname"],
			Cluster:   alert.Labels["cluster"],
			Instance:  alert.Labels["instance"],
			Namespace: alert.Labels["namespace"],
		}
		if IsBalnk(alert.Annotations.Graphs) {
			message.GraphsURL = cfg.GrafanaUrl
		} else {
			message.GraphsURL = alert.Annotations.Graphs
		}
		for key, value := range alert.Labels {
			message.LabelsString = message.LabelsString + ",\"" + key + "\":\"" + value + "\""
		}
		message.LabelsString = message.LabelsString[1:]

		switch alert.Annotations.Level {
		case "fatal":
			message.Color = "red"
		case "warning":
			message.Color = "orange"
		default:
			message.Color = "blue"
		}

		if !IsBalnk(alert.Annotations.Summary) {
			message.Content = "**Summary:**\\n" + alert.Annotations.Summary
		}
		if !IsBalnk(alert.Annotations.Proposal) {
			message.Content = message.Content + "\\n**Proposal:**\\n" + alert.Annotations.Proposal
		}
		if !IsBalnk(alert.Annotations.Note) {
			message.Content = message.Content + "\\n**Note:**\\n" + alert.Annotations.Note
		}

		template_, _ := template.ParseFiles("template/feishu-message.txt")
		var tpl bytes.Buffer
		eerr := template_.Execute(&tpl, message)
		if eerr != nil {
			log.Printf("%+v", message)
			log.Fatal(eerr)
		}
		log.Printf(tpl.String() + "\n")
		sendMessage(tpl.String(), false, "")
	}
}
