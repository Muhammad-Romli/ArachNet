package main

import (
	"fmt"

	"github.com/muhammad-romli/ArachNet/internal/config"
)

func readLoop(maxRetries int, data *config.Config, err error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		data, err = config.Read()
		if err == nil {
			break
		}
		fmt.Printf("Retrying reading")
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	maxRetries := 2
	var data *config.Config
	var err error
	readLoop(maxRetries, data, err)
	// The data is just place holder in the first readLoop
	data.SetUser("Justicar")
	readLoop(maxRetries, data, err)
}
