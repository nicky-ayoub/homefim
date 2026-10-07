package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"homefim/internal/baseline"
)

const scanInterval = 30 * time.Second

type Baseline map[string]string

func main() {
	loop := flag.Bool("loop", false, "continuously check integrity every 30 seconds")
	flag.Parse()

	var currentBaseline Baseline
	if err := json.Unmarshal(baseline.Raw, &currentBaseline); err != nil {
		log.Fatalf("[-] Failed to parse embedded baseline data: %v", err)
	}

	log.Printf("[+] FIM engine initialized. Loaded %d audit targets from embedded baseline.", len(currentBaseline))
	error_count := verifyIntegrity(currentBaseline)

	if !*loop {
		os.Exit(error_count)
	}

	log.Printf("[+] NIS compatibility mode active (Polling every %v)...", scanInterval)
	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()
	for range ticker.C {
		error_count += verifyIntegrity(currentBaseline)
	}
}

func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func verifyIntegrity(currentBaseline Baseline) int {
	error_count := 0
	for target, baselineHash := range currentBaseline {
		_, err := os.Stat(target)
		if os.IsNotExist(err) {
			log.Printf("[🚨 CRITICAL] %s has been REMOVED from the NIS share!", target)
			error_count++
			continue
		}

		newHash, err := calculateHash(target)
		if err != nil {
			log.Printf("[-] Error reading target %s over network: %v", target, err)
			error_count++
			continue
		}

		if newHash != baselineHash {
			log.Printf("[🚨 CRITICAL INTEGRITY COMPROMISE] Hash mismatch on NIS share: %s", target)
			log.Printf("    -> Expected: %s", baselineHash)
			log.Printf("    -> Observed: %s", newHash)
		}
	}
	return error_count
}
