package config

import (
    "encoding/json"
    "fmt"
    "os"
    "jwt-auth-service/models"
)

type Config struct {
    ServerPort string            `json:"server_port"`
    TargetAPI  string            `json:"target_api"`
    Clients    []models.Client   `json:"clients"`
}

var AppConfig Config

func LoadConfig() error {
    defaultConfig := Config{
        ServerPort: "8080",
        TargetAPI:  "https://jsonplaceholder.typicode.com",
        Clients: []models.Client{
            {
                ID:            "client1",
                Name:          "Test Client 1",
                AllowedPaths:  []string{"/posts", "/comments"},
                AllowedMethods: []string{"GET", "POST"},
                TokenTTL:      15,
                SecretKey:     "secret-key-1",
            },
            {
                ID:            "client2",
                Name:          "Test Client 2",
                AllowedPaths:  []string{"/users", "/albums"},
                AllowedMethods: []string{"GET"},
                TokenTTL:      30,
                SecretKey:     "secret-key-2",
            },
        },
    }

    file, err := os.ReadFile("config.json")
    if err != nil {
        fmt.Println("Config file not found, using default configuration")
        AppConfig = defaultConfig
        return nil
    }

    err = json.Unmarshal(file, &AppConfig)
    if err != nil {
        return fmt.Errorf("error parsing config file: %v", err)
    }

    return nil
}

func GetClientByID(clientID string) *models.Client {
    for _, client := range AppConfig.Clients {
        if client.ID == clientID {
            return &client
        }
    }
    return nil
}
EOF
