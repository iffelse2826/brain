package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)


// request payload struct
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

func loadConfig(path string) (*Config, error) {
	// Candidate paths to search for config file
	paths := []string{
		path,
		"config/config.yaml",
		"../config/config.yaml",
		"backend/config/config.yaml",
	}

	var data []byte
	var err error
	var found bool

	for _, p := range paths {
		if p == "" {
			continue
		}
		data, err = os.ReadFile(p)
		if err == nil {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("config file not found")
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	configPath := "../config/config.yaml"

	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	dbConfig := config.Database

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	log.Printf("Successfully loaded configuration. Database Host: %s", dbConfig.Host)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close() // ensure the connection is closed when the function returns

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping the database: ", err)
	}

	fmt.Println("Successfully connected to the database")
}
