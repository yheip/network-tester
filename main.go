package main

import (
	"context"

	"github.com/yheip/network-tester/internal/application"
)

func main() {
	ctx := context.Background()
	app := application.New(ctx)
	app.Start(ctx)
}
