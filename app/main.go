package main

import (
	"fmt"
	"log"
	"context"

	"github.com/hashicorp/vault/api"
)

func main() {
	// Load the configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Configure the Vault client
	client, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("http://%s:8200", config.VaultServiceName),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Set the Vault token for authentication
	client.SetToken(config.VaultToken)

	// Check connection
	healthResponse, err := client.Sys().Health()
	if err != nil {
		log.Fatal(err)
	}

	if healthResponse.Initialized {
		fmt.Println("Successfully authenticated with Vault.")
	}

	// Write a secret to Vault
	data := map[string]interface{}{
		"username": "myuser",
		"password": "mypassword",
	}
	// Userdata is the name of the secret
	_, err = client.KVv2("gosecret").Put(context.Background(), "userdata", data)
	if err != nil {
		log.Fatalf("unable to write secret: %v", err)
	}
	fmt.Println("Secret written successfully.")

	// Read a secret from Vault
	secret, err := client.KVv2("gosecret").Get(context.Background(), "userdata")
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	// Extract the secret data
	if secret != nil {
		secretData := secret.Data
		fmt.Println(secretData)
		username := secretData["username"].(string)
		password := secretData["password"].(string)
		fmt.Println("Username:", username)
		fmt.Println("Password:", password)
	}
}
