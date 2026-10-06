package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("[-] Failed to get current user: %v", err)
	}
	home := usr.HomeDir

	targets := []string{
		filepath.Join(home, ".bash_profile"),
		filepath.Join(home, ".bash_login"),
		filepath.Join(home, ".profile"),
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".bash_logout"),
	}

	sshDir := filepath.Join(home, ".ssh")
	err = filepath.Walk(sshDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		baseName := filepath.Base(path)
		if baseName == "known_hosts" || baseName == "known_hosts.old" {
			log.Printf("[*] Skipping volatile file: %s", path)
			return nil
		}

		targets = append(targets, path)
		return nil
	})
	if err != nil {
		log.Printf("[-] Warning during .ssh tree walk: %v", err)
	}

	baseline := make(map[string]string)

	log.Println("[*] Scanning discovered paths to generate baseline...")
	for _, target := range targets {
		if _, err := os.Stat(target); err == nil {
			hash, err := hashFile(target)
			if err != nil {
				log.Printf("[-] Skipping %s: %v", target, err)
				continue
			}
			baseline[target] = hash
			log.Printf("[+] Baselined: %s", target)
		}
	}

	outputData, err := json.MarshalIndent(baseline, "", "  ")
	if err != nil {
		log.Fatalf("[-] Failed to marshal JSON: %v", err)
	}

	const outputPath = "internal/baseline/baseline.json"
	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		log.Fatalf("[-] Failed to write %s: %v", outputPath, err)
	}

	log.Printf("[+] Success! '%s' created with %d targets.", outputPath, len(baseline))
}

func hashFile(filePath string) (string, error) {
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
