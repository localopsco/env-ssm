package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/hashicorp/go-envparse"
)

var envFilePath = ".env"
var awsSsmPath = "/"
var keepOldEnvInSSM = true

// read .env file
func readFile(filePath string) ([]byte, error) {
	envFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return envFile, nil
}

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the env file path").
				Inline(true).
				Value(&envFilePath),

			huh.NewInput().
				Title("Enter AWS SSM path").
				Inline(true).
				Validate(huh.ValidateMinLength(1)).
				Value(&awsSsmPath),

			huh.NewConfirm().
				Title("Keep old env in SSM").
				Inline(true).
				Value(&keepOldEnvInSSM),
		),
	)

	if err := form.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[error] unable to run TUI: %v\n", err)
		os.Exit(1)
	}

	// check if the env file exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		log.Fatal("[error] .env file does not exist")
	}

	envFile, err := readFile(envFilePath)
	if err != nil {
		log.Fatal("[error] unable to read .env file", err)
	}

	envs, err := envparse.Parse(bytes.NewReader(envFile))
	if err != nil {
		log.Fatal("[error] unable to parse .env file", err)
	}

	// load envs from aws ssm
	ssmParams := loadConfigFromAws(awsSsmPath)

	// if keepOldEnvInSSM is true, we need to keep the old envs in SSM
	// otherwise, we need to delete the old envs in SSM
	keysToDelete := []string{}
	for _, param := range ssmParams {
		if _, ok := envs[*param.Name]; !ok {
			keysToDelete = append(keysToDelete, *param.Name)
		}
	}

	var deleteKeysConfirm bool
	deleteKeys := len(keysToDelete) > 0 && !keepOldEnvInSSM

	// if there are keys to delete, we need to delete them, but ask for confirmation
	if deleteKeys {
		keysToDeleteStr := "\n"
		for _, key := range keysToDelete {
			keysToDeleteStr += "  - " + awsSsmPath + key + "\n"
		}

		keysToDeleteStr += "Are you sure?"

		confirm := huh.NewConfirm().
			Inline(true).
			Title("The following keys will be deleted" + keysToDeleteStr).
			Value(&deleteKeysConfirm)

		if err := confirm.Run(); err != nil {
			log.Fatal("[error] unable to run TUI", err)
		}

		if !deleteKeys {
			log.Fatal("[info] user cancelled")
		}
	}

	// add new envs to aws ssm
	addSsmParams(awsSsmPath, envs)

	// delete old envs from aws ssm
	if deleteKeys {
		deleteSsmParams(keysToDelete)
	}

	fmt.Println("[success] Done")
}
