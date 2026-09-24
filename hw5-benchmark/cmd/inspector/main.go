package main

import (
	"fmt"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/adapter"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw5-benchmark/internal/application"
)

// import (
// 	"reflect"
// )

func main() {
	cfg := adapter.ParseCmdLine()
	output, err := application.StartInspect(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	err = os.WriteFile(cfg.Output, []byte(output), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)

}
