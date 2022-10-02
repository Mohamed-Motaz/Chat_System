package utils

import (
	"log"
	"os"
)

//THIS METHOD SHOULD BE CALLED ONLY ONCE
func InitDirectories(dataFolder string) {

	err := os.MkdirAll(dataFolder, os.ModePerm)
	if err != nil {
		log.Fatalf("Unable to initialize directories %v\n", err)
	}

	log.Printf("Done initializing directories\n\n")
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
