package main

import (
	"context"

	v1 "github.com/kamuridesu/kavideos/internal/api/v1"
)

func main() {
	ctx := context.Background()

	v1.SetupRoutes(ctx)
}
