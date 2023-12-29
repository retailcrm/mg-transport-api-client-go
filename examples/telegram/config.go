package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
	"strings"
)

var AppConfig Config

type Config struct {
	Listen     string `json:"listen"`
	BaseURL    string `json:"baseUrl" validate:"required,url"`
	System     string `json:"system" validate:"required,url"`
	APIKey     string `json:"apiKey" validate:"required"`
	TGBotToken string `json:"tgBotToken" validate:"required"`
}

func LoadConfig(src string) {
	file, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()
	if err := json.NewDecoder(file).Decode(&AppConfig); err != nil {
		panic(err)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(AppConfig)
	if err != nil {
		log.Fatalln(err)
	}
	if strings.HasSuffix(AppConfig.BaseURL, "/") {
		AppConfig.BaseURL = AppConfig.BaseURL[:len(AppConfig.BaseURL)-1]
	}
	if AppConfig.Listen == "" {
		AppConfig.Listen = ":8080"
	}
	log.Println("loaded configuration from", src)
}
