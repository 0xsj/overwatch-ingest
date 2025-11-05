package main

import (
	"fmt"
	"log"
)

func main() {
	log.Println("Gateway service starting")

	fmt.Println("Gateway running on :8080")

	select {}
}