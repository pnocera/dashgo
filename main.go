package main

import (
	"github.com/pnocera/dashgo/cmd"
	"github.com/pnocera/dashgo/config"
)

func main() {
	port := config.New().APIPort()
	cmd.RunWebServer(port)
}
