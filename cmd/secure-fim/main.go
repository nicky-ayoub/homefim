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
	verifyIntegrity(currentBaseline)

	if !*loop {
		return
	}

	log.Printf("[+] NIS compatibility mode active (Polling every %v)...", scanInterval)
	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()
	for range ticker.C {
		verifyIntegrity(currentBaseline)
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

func verifyIntegrity(currentBaseline Baseline) {
	for target, baselineHash := range currentBaseline {
		_, err := os.Stat(target)
		if os.IsNotExist(err) {
			log.Printf("[🚨 CRITICAL] %s has been REMOVED from the NIS share!", target)
			continue
		}

		newHash, err := calculateHash(target)
		if err != nil {
			log.Printf("[-] Error reading target %s over network: %v", target, err)
			continue
		}

		if newHash != baselineHash {
			log.Printf("[🚨 CRITICAL INTEGRITY COMPROMISE] Hash mismatch on NIS share: %s", target)
			log.Printf("    -> Expected: %s", baselineHash)
			log.Printf("    -> Observed: %s", newHash)
		}
	}
}
