package main

import (
	"github.com/muhammad-romli/ArachNet/internal/config"
)

func retryLoop[A any](func(A)) (data, err) {
	maxRetries := 2

	for attempt := 0; attempt < maxRetries; attempt++ {
		data, err := config.Read()
		if err == nil {
			break
		}

	}
}

func main() {
	var data config.Config

}
