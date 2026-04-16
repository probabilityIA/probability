package main

import (
	"context"

	"github.com/secamc93/probability/back/central/cmd/internal/server"
)

func main() {
	_ = server.Init(context.Background())
	select {}
}

