package main

import (
	"Momotok-Server/system"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	err := system.LoadConfigInformation()
	if err != nil {
		fmt.Printf("Failed to load config information: %s\n", err)
	}
	println("Config loaded!Token secret key is:" + system.ServerInfo.Server.SecretKey)

	r := gin.Default()
	gin.SetMode(system.ServerInfo.Server.Mode)
	initRouter(r)
	r.Run(system.ServerInfo.Server.Host + ":" + system.ServerInfo.Server.Port) //listen and serve on 0.0.0.0:8080 by default
}
