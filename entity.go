package main

import "time"

type AlertMessage struct {
	Receiver string  `json:"receiver"`
	Status   string  `json:"status"`
	Alerts   []Alert `json:"alerts"`
	_Version string  `json:"version"`
}

type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  Annotations       `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

type Labels struct {
	Alertname string `json:"alertname"`
	Cluster   string `json:"cluster"`
	Instance  string `json:"instance"`
	Namespace string `json:"namespace"`
}

type Annotations struct {
	Level    string `json:"level"`
	Summary  string `json:"summary"`
	Proposal string `json:"proposal"`
	Graphs   string `json:"graphs"`
	Note     string `json:"note"`
}

type AppAccessTokenEntity struct {
	Code           int    `json:"code"`
	Msg            string `json:"msg"`
	AppAccessToken string `json:"app_access_token"`
	expire         int64  `json:"expire"`
}

type UpdatingCardAction struct {
	Tag   string            `json:"tag"`
	Value map[string]string `json:"value"`
}

type UpdatingCard struct {
	OpenId        string             `json:"open_id"`
	UserId        string             `json:"user_id"`
	OpenMessageId string             `json:"open_message_id"`
	TenantKey     string             `json:"tenant_key"`
	Token         string             `json:"token"`
	Action        UpdatingCardAction `json:"action"`
}

type UpdatingCardNecessity struct {
	MessageId string
	SilenceId string
	Labels    map[string]string
}

type Matcher struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	IsRegex bool   `json:"isRegex"`
	IsEqual bool   `json:"isEqual"`
}

type SilenceEntity struct {
	Matchers  []Matcher `json:"matchers"`
	StartsAt  string    `json:"startsAt"`
	EndsAt    string    `json:"endsAt"`
	CreatedBy string    `json:"createdBy"`
	Comment   string    `json:"comment"`
}

type SilenceResponse struct {
	SilenceId string `json:"silenceID"`
}

type SendingMessageResponse struct {
	Code int `json:"code"`
}

type FeishuMessageTemplateEntity struct {
	Color        string
	AlertName    string
	Datetime     string
	Cluster      string
	Instance     string
	Namespace    string
	Content      string
	GraphsURL    string
	LabelsString string
	Revive       bool
}
