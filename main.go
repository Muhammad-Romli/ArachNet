package main

import (
	"fmt"

	"github.com/muhammad-romli/ArachNet/internal/config"
)

func readLoop(maxRetries int) *config.Config {
	var data *config.Config
	var err error
	for attempt := 0; attempt < maxRetries; attempt++ {
		data, err = config.Read()
		if err == nil {
			break
		}
		fmt.Printf("Retrying reading\n")
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return data
}

func main() {
	maxRetries := 2
	data := readLoop(maxRetries)
	if data == nil {
		fmt.Printf("There is no data inside targeted file")
		return
	}
	// The data is just placeholder in the first readLoop
	data.SetUser("Justicar")
	fmt.Println(readLoop(maxRetries))
}
