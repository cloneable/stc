package main

import (
	"context"

	"github.com/cloneable/stacker/internal/command"
)

func main() {
	ctx := context.Background()
	command.Execute(ctx)
}
