package main

import (
	"fmt"
	"github.com/rluisr/tvbit-bot/pkg/external"
)

func main() {
	err := external.Router.Run()
	if err != nil {
		panic(fmt.Errorf("failed to run router %w", err))
	}
}
