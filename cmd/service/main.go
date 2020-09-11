package main

import (
	"log"

	"github.com/Gorynychdo/tdligo.git/internal/client"
)

func main() {
	tc := client.NewTelegramClient()
	if err := tc.Start(); err != nil {
		log.Fatal(err)
	}
}
