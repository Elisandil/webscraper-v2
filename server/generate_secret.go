//go:build ignore
// +build ignore

package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	secret := base64.URLEncoding.EncodeToString(bytes)
	fmt.Printf("Your new JWT secret (copy this to config.yaml):\n\n%s\n", secret)
}
