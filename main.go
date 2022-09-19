package main

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var _cache = cache.New(5*time.Minute, 10*time.Minute)

var _client *http.Client

var cfg Config

type Config struct {
	GrafanaUrl   string `yaml:"grafanaUrl"`
	AlertManager string `yaml:"alertManager"`
	Feishu       Feishu `yaml:"feishu"`
}
type Feishu struct {
	AppId             string `yaml:"appId"`
	AppSecret         string `yaml:"appSecret"`
	ChatId            string `yaml:"chatId"`
	AppAccessTokenUrl string `yaml:"appAccessTokenUrl"`
	ReceiveMessageUrl string `yaml:"receiveMessageUrl"`
	UpdateMessageUrl  string `yaml:"updateMessageUrl"`
}

func main() {

	initConfig()
	transmis := &http.Transport{
		MaxIdleConns:       8,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	_client = &http.Client{
		Transport: transmis,
	}

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := gin.Default()

	r.POST("/api/callback", silences)

	r.POST("/api/receiveAlert", receiveAlert)

	r.Run()
}

func initConfig() {
	env := os.Getenv("GO_ENV")
	viper.AddConfigPath("./resource")
	if !IsBalnk(env) {
		viper.SetConfigName("application-" + env)
	} else {
		viper.SetConfigName("application")
	}
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config file failed, %v", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("unmarshal config file failed, %v", err)
	}
	log.Printf("%+v", &cfg)
}
