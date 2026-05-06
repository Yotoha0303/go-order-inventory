package main

import (
	"fmt"
	"go-order-inventory/config"
)

func main() {
	config.LoadEnv()

	fmt.Println("hello go-order-invertory")
}
