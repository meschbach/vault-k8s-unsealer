package main

import (
	"log"
	"os"
)

func main()  {
	keyFile, found := os.LookupEnv("UNSEAL_FILE")
	if !found {
		log.Fatal("Env var UNSEAL_FILE must be set to the key file")
	}
	if err := controlLoopUnseal(keyFile); err != nil {
		panic(err)
	}
}
