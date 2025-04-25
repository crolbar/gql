package util

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

func DeleteCache() {
	cacheDir, err := os.UserCacheDir()

	if err != nil {
		log.Fatalf("Error while getting userCacheDir: %v", err)
	}

	cacheFile := filepath.Join(cacheDir, "gql/gql_creds")

	err = os.Remove(cacheFile)
	if err != nil {
		log.Fatalf("Error removing cache file: %v", err)
	}
}

func WriteToCacheFile(cont string) {
	cacheDir, err := os.UserCacheDir()

	if err != nil {
		log.Fatalf("Error while getting userCacheDir: %v", err)
	}

	cacheDir = filepath.Join(cacheDir, "gql")
	cacheFile := filepath.Join(cacheDir, "gql_creds")

	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating cache file: %v", err)
	}

	f, err := os.Create(cacheFile)
	if err != nil {
		log.Fatalf("Error creating cache file: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(cont)

	if err != nil {
		log.Fatalf("Error writing to cache file: %v", err)
	}
}

func CacheFileExists() bool {
	cacheDir, err := os.UserCacheDir()

	if err != nil {
		log.Fatalf("Error while getting userCacheDir: %v", err)
	}

	cacheFile := filepath.Join(cacheDir, "gql/gql_creds")

	_, err = os.Stat(cacheFile)

	return !os.IsNotExist(err)
}

func ReadFromCacheFile() string {
	cacheDir, err := os.UserCacheDir()

	if err != nil {
		log.Fatalf("Error while getting userCacheDir: %v", err)
	}

	cacheFile := filepath.Join(cacheDir, "gql/gql_creds")

	f, err := os.Open(cacheFile)
	if err != nil {
		log.Fatalf("Error creating cache file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	return scanner.Text()
}
