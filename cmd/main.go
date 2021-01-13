package main

import (
	"context"
	"fmt"

	"github.com/anfilat/final-stats/internal/loadavg"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := loadavg.GetStream(ctx)

	for avg := range ch {
		fmt.Println(avg)
	}
}
