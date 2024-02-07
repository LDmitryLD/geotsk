package main

import (
	"fmt"
	"log"
	"os"
	"projects/LDmitryLD/geotask/run"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := run.NewApp()

	err := app.Run()
	if err != nil {
		log.Println(fmt.Sprintf("error: %s", err))
		os.Exit(2)
	}
}
