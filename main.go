package main

import (
    "github.com/azzamt11/todoapp_backend/app"
    "github.com/azzamt11/todoapp_backend/config"
    "os"
)

func main() {
    // Load configuration
    cfg := config.GetConfig()

    // Initialize the app
    a := &app.App{}
    a.Initialize(cfg)

    // Get the port from the environment variable or use default (3000)
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000" // Default port
    }

    // Start the app and listen on the specified port
    a.Run(":" + port)
}
