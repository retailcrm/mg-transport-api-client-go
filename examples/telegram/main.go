package main

func main() {
	LoadConfig("config.json")
	RegisterSystem()
	InitTGBotAPI()
	SetTGWebhook()
	RegisterChannel()
	Listen()
}
