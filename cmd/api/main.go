package main

import (
	"log"

	"github.com/haramurti/valeth-urls-shortener/cmd/bootstrap"
)

func main() {
	app := bootstrap.NewApp()

	log.Fatal(app.Fiber.Listen(":8098"))
}

//i need to make a portofolio first for apple developer academy
