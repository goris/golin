package config

import (
    "encoding/json"
    "fmt"
    "os"
)

type Configuration struct {
    AccountsDB DBConfig
    Schema Schema
}

type Schema struct {
    Table string
    Email string
    Password string
}

type DBConfig struct {
    Host string
    Port int
    DBName string
    Username string
    Password string
}

func ReadConfig(f string) (*Configuration, error) {

    file, err := os.Open(f)
    if err != nil {
        fmt.Println("Config - File Error: %s", err)
        return nil, fmt.Errorf("Config - File Error: %s", err)
    }

    decoder := json.NewDecoder(file)
    configuration := Configuration{}

    if err := decoder.Decode(&configuration); err != nil {
        fmt.Println("Config - Decoding Error: %s", err)
        return nil, fmt.Errorf("Config - Decoding Error: %s", err)
    }
    return &configuration, nil
}
