package main

import (
	"github.com/azzamt11/todoapp_backend/app"
	"github.com/azzamt11/todoapp_backend/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
