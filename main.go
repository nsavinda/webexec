package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

var CONFIG_FILE = "/etc/webexec/config.yaml"

type Config struct {
	Key     string `yaml:"key"`
	Command string `yaml:"command"`
	Dir     string `yaml:"dir"`
	Port    string `yaml:"port"`
}

func loadConfig(filename string) (Config, error) {
	var config Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	return config, err
}

func verifySignature(secret, signature string, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

func webhookHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		signature := r.Header.Get("X-Hub-Signature-256")
		if signature == "" {
			http.Error(w, "Missing signature", http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		if !verifySignature(config.Key, signature, body) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		// Change to specified directory
		if err := os.Chdir(config.Dir); err != nil {
			log.Printf("Error changing directory: %v", err)
			http.Error(w, "Directory error", http.StatusInternalServerError)
			return
		}

		// Split command into name and args
		parts := strings.Fields(config.Command)
		if len(parts) == 0 {
			http.Error(w, "Invalid command", http.StatusBadRequest)
			return
		}
		cmdName := parts[0]
		cmdArgs := parts[1:]

		// Execute command
		cmd := exec.Command(cmdName, cmdArgs...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Command execution failed: %v, output: %s", err, string(output))
			http.Error(w, fmt.Sprintf("Command failed: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Command executed successfully: %s", string(output))
		fmt.Fprintf(w, "Command executed: %s", string(output))
	}
}

func main() {
	config, err := loadConfig(CONFIG_FILE)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	http.HandleFunc("/webhook", webhookHandler(config))
	log.Printf("Starting server on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
