package main

import (
	"github.com/meschbach/vault-k8s-unsealer/internal"
	"log"
	"os"
)

func main()  {
	keyFile, found := os.LookupEnv("UNSEAL_FILE")
	if !found {
		log.Fatal("Env var UNSEAL_FILE must be set to the key file")
	}
	if err := internal.ControlLoopUnseal(keyFile); err != nil {
		panic(err)
	}
}
