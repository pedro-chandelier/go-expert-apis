package main

import "github.com/pedro-chandelier/go-expert-apis/configs"

func main() {
	config, _ := configs.LoadConfig("configs/.env")
	println(config.DBHost)
}
