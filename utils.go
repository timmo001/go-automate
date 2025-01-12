package main

import (
	"crypto/rand"
	"math/big"

	"github.com/charmbracelet/log"
)

func RandomID() int {
	reader := rand.Reader
	n, err := rand.Int(reader, big.NewInt(1000))
	if err != nil {
		log.Fatalf("error generating random ID: %v", err)
	}
	return int(n.Int64())
}
