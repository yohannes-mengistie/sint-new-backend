package config

import (
	"encoding/base64"
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr         string
	DatabaseURL      string
	CloudinaryURL    string
	CloudinaryFolder string
	FirebaseJSON     []byte
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}

	config := Config{
		HTTPAddr:         os.Getenv("HTTP_ADDR"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		CloudinaryURL:    os.Getenv("CLOUDINARY_URL"),
		CloudinaryFolder: os.Getenv("CLOUDINARY_FOLDER"),
	}
	if config.HTTPAddr == "" {
		config.HTTPAddr = ":" + os.Getenv("PORT")
		if config.HTTPAddr == ":" {
			config.HTTPAddr = ":8082"
		}
	}
	if config.CloudinaryFolder == "" {
		config.CloudinaryFolder = "sint-new/photos"
	}
	switch {
	case config.DatabaseURL == "":
		return Config{}, errors.New("DATABASE_URL is required")
	case config.CloudinaryURL == "":
		return Config{}, errors.New("CLOUDINARY_URL is required")
	}
	if encoded := os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON_B64"); encoded != "" {
		credentials, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return Config{}, errors.New("FIREBASE_SERVICE_ACCOUNT_JSON_B64 is invalid")
		}
		config.FirebaseJSON = credentials
	}
	return config, nil
}
