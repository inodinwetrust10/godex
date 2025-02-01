package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

const (
	configDir     = ".config/gophile"
	tokenFileName = "token.json"
)

func GetConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to find home directory: %v", err)
	}
	return filepath.Join(home, configDir)
}

func getTokenFilePath() string {
	return filepath.Join(GetConfigDir(), tokenFileName)
}

func getClient(config *oauth2.Config) *http.Client {
	tokenPath := getTokenFilePath()

	tok, err := tokenFromFile(tokenPath)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenPath, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("Go to the following link in your browser:")
	fmt.Println(authURL)
	fmt.Println("\nAfter approving the permissions, you'll get an authorization code.")
	fmt.Println("Paste that code here:")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalf("Unable to create config directory: %v", err)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
