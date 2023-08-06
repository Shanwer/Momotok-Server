package main

import (
	"Momotok-Server/common"
	"Momotok-Server/service"
	"Momotok-Server/system"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	err := system.LoadConfigInformation()
	if err != nil {
		fmt.Printf("Failed to load config information: %s\n", err)
	}
	println(common.ServerInfo.SecretKey)

	go service.RunMessageServer()
	r := gin.Default()
	initRouter(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
