package main

import (
	"context"
	"github.com/nileshshrs/infinite-storage/application"
	"fmt"
)

func main() {
	app := application.New()
	err:=app.Start(context.TODO())
	if err != nil {
		fmt.Println("Error starting application: %v", err)
	}
}
