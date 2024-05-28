package main

import (
    "log"
    "rafikichat/internal/infrastructure/config"
    "rafikichat/internal/infrastructure/server"
)

// main is the function for running main function
func main() {
    cfg := config.LoadConfig()
    srv := server.NewServer(cfg)
    log.Fatal(srv.Start())
}