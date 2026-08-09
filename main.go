package main

import (
	"fmt"
	"log"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
)

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	err = config.SetUser("Henrik")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config)
}
