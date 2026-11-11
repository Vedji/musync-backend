package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Capabilities struct {
	Uploads   bool `json:"uploads"`
	Downloads bool `json:"downloads"`
}

type Limits struct {
	MaxUploadBytes   int64 `json:"max_upload_bytes"`
	MaxPlaylistItems int64 `json:"max_playlist_items"`
	MaxBulkItems     int64 `json:"max_bulk_items"`
}

type Info struct {
	Name                      string       `json:"name"`
	Version                   string       `json:"version"`
	APIVersion                string       `json:"api_version"`
	RegistrationEnabled       bool         `json:"registration_enabled"`
	Capabilities              Capabilities `json:"capabilities"`
	Limits                    Limits       `json:"limits"`
	SupportedUploadExtensions []string     `json:"supported_upload_extensions"`
}

func envString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	value := envString(key, strconv.FormatBool(fallback))
	result, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("%s must be a boolean: %q", key, value)
	}
	return result
}

func envLimit(key string, fallback int64) int64 {
	value := envString(key, strconv.FormatInt(fallback, 10))
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil || result < 0 {
		log.Fatalf("%s must be a non-negative integer: %q", key, value)
	}
	return result
}

func envExtensions() []string {
	value := envString(
		"SUPPORTED_UPLOAD_EXTENSIONS",
		"mp3,flac,m4a,ogg,opus,wav,aac",
	)

	result := make([]string, 0)
	for _, extension := range strings.Split(value, ",") {
		extension = strings.TrimPrefix(
			strings.ToLower(strings.TrimSpace(extension)),
			".",
		)
		if extension != "" {
			result = append(result, extension)
		}
	}
	return result
}

func LoadConfig() Info {
	return Info{
		Name:                envString("APP_NAME", "MuSync"),
		Version:             envString("APP_VERSION", "0.1.0"),
		APIVersion:          envString("API_VERSION", "v1"),
		RegistrationEnabled: envBool("REGISTRATION_ENABLED", true),
		Capabilities: Capabilities{
			Uploads:   envBool("UPLOADS_ENABLED", true),
			Downloads: envBool("DOWNLOADS_ENABLED", true),
		},
		Limits: Limits{
			MaxUploadBytes:   envLimit("MAX_UPLOAD_BYTES", 524288000),
			MaxPlaylistItems: envLimit("MAX_PLAYLIST_ITEMS", 10000),
			MaxBulkItems:     envLimit("MAX_BULK_ITEMS", 500),
		},
		SupportedUploadExtensions: envExtensions(),
	}
}
