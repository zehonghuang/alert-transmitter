package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
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

func assembleMessageCard(alertmessage AlertMessage) {
	var alertMap = make(map[string]FeishuMessageTemplateEntity)
	var alertPiecesMap = make(map[string]map[string]string)

	for _, alert := range alertmessage.Alerts {
		var message FeishuMessageTemplateEntity
		if v, ok := alertMap[alert.hashcode()]; !ok {
			message = FeishuMessageTemplateEntity{
				Datetime:     time.Now().Format("2006-01-02 15:04:05"),
				AlertName:    alert.Labels["alertname"],
				Cluster:      alert.Labels["cluster"],
				Instance:     alert.Labels["instance"],
				Namespace:    alert.Labels["namespace"],
				AlertCounter: 1,
			}
			alertMap[message.hashcode()] = message
		} else {
			message = v
			message.AlertCounter++
		}

		if IsBlank(alert.Annotations.Graphs) {
			message.GraphsURL = cfg.GrafanaUrl
		} else {
			message.GraphsURL = alert.Annotations.Graphs
		}

		if message.AlertCounter == 1 {
			for key, value := range alert.Labels {
				if key == "pieces" {
					continue
				}
				message.LabelsString = message.LabelsString + ",\"" + key + "\":\"" + value + "\""
			}
			message.LabelsString = message.LabelsString[1:]
		}
		if !IsBlank(alert.Labels["pieces"]) {
			var piecesMap map[string]string
			if map_, ok := alertPiecesMap[message.hashcode()]; ok {
				piecesMap = map_
			} else {
				piecesMap = make(map[string]string)
			}

			kv := strings.Split(alert.Labels["pieces"], ",")
			for _, r := range kv {
				arr := strings.Split(r, "=")
				if v, ok := piecesMap[arr[0]]; ok {
					v = v + "|" + arr[1]
					piecesMap[arr[0]] = v
				} else {
					piecesMap[arr[0]] = arr[1]
				}
			}
			alertPiecesMap[message.hashcode()] = piecesMap
		}

		switch alert.Annotations.Level {
		case "fatal":
			message.Color = "red"
		case "warning":
			message.Color = "orange"
		default:
			message.Color = "blue"
		}

		if !IsBlank(alert.Labels["pieces"]) && (message.AlertCounter <= 3) {
			if IsBlank(message.PiecesLabels) && (message.AlertCounter == 1) {
				message.PiecesLabels = alert.Labels["pieces"]
			} else {
				message.PiecesLabels = message.PiecesLabels + "\\n" + alert.Labels["pieces"]
			}
		}

		if !IsBlank(alert.Annotations.Summary) {
			message.Content = "**Summary:**\\n" + alert.Annotations.Summary
		}
		if !IsBlank(alert.Annotations.Proposal) {
			message.Content = message.Content + "\\n**Proposal:**\\n" + alert.Annotations.Proposal
		}
		if !IsBlank(alert.Annotations.Note) {
			message.Content = message.Content + "\\n**Note:**\\n" + alert.Annotations.Note
		}

		alertMap[message.hashcode()] = message
	}

	for _, v := range alertMap {
		piecesNameArr := make([]string, 1)
		if pieces, ok := alertPiecesMap[v.hashcode()]; ok {
			for l, lv := range pieces {
				piecesNameArr = append(piecesNameArr, l)
				v.LabelsString = v.LabelsString + ",\"" + l + "\":\"" + lv + "\""
			}
		}

		template_, _ := template.ParseFiles("template/feishu-message.txt")
		var tpl bytes.Buffer
		eerr := template_.Execute(&tpl, v)
		if eerr != nil {
			log.Fatal(eerr)
		}
		log.Printf("%s", tpl.String())
		sendMessage(tpl.String(), false, "")
	}

}
