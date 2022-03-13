package internal

import (
	"errors"
	vault "github.com/hashicorp/vault/api"
	"io/ioutil"
	"log"
	"time"
)

const (
	stateError = iota
	stateUnsealed
	stateSealed
)

type autoUnseal struct {
	target *vault.Sys
	keyFile string
}

func (a *autoUnseal) checkAndUnseal() (bool, error) {
	var err error
	var status *vault.SealStatusResponse

	status, err = a.target.SealStatus()
	if err != nil {
		return false, err
	}
	if status.Sealed {
		log.Print("Vault sealed.  Proceeding.")
	} else {
		log.Println("Vault not sealed")
		return true, nil
	}

	var keyBytes []byte
	keyBytes, err = ioutil.ReadFile(a.keyFile)
	if err != nil {
		log.Printf("Failed to read key file")
		return false, err
	}

	key := string(keyBytes)

	status, err = a.target.Unseal(key)
	if err != nil {
		log.Printf("Unable to unseal")
		return false, err
	}
	if status.Sealed {
		return false, errors.New("sealed after key(s) given")
	}
	return true, nil
}

func ControlLoopUnseal(keyFile string) error {
	config := vault.DefaultConfig()

	client, err := vault.NewClient(config)
	if err != nil { return err }

	sys := client.Sys()
	autoUnseal := &autoUnseal{
		target:  sys,
		keyFile: keyFile,
	}

	backoffTimer := &exponetialBackoff{
		unit: time.Second,
		base:             2,
		limit:            10 * time.Minute,
		currentIncrement: 0,
		state:            0,
	}
	for {
		log.Print("Checking seal status")
		unsealed, problem := autoUnseal.checkAndUnseal()
		var state int
		if problem != nil {
			log.Printf("Failed to unseal beacuse %s", problem.Error())
			state = stateError
		} else if unsealed {
			log.Print("Still unsealed")
			state = stateUnsealed
		} else {
			log.Print("Found sealed, unsealed successfully")
			state = stateSealed
		}
		backoffTimer.performBackoff(state)
	}
}
