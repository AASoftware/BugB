package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}

	filename := os.Args[1]

	// Ordner erstellen, falls er nicht existiert
	outputDir := "data_head_request"
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", outputDir, err)
		return
	}

	// Datei öffnen
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	var wg sync.WaitGroup
	urlChannel := make(chan string)

	// Start der Worker-Goroutines
	for i := 0; i < 70; i++ { // max. 70 gleichzeitige Worker/Threads getestet
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urlChannel {
				processURL(url, outputDir)
			}
		}()
	}

	// URLs aus der Datei lesen und in den Channel schicken
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), `"`) // Entfernt Hochkommas
		urlChannel <- line
	}

	close(urlChannel) // Channel schließen, um die Worker zu beenden
	wg.Wait()         // Warten, bis alle Worker fertig sind

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	fmt.Println("Fertig!")
}

// processURL führt die HEAD-Anfrage durch und schreibt die URL in die entsprechende Datei
func processURL(url, outputDir string) {
	// HTTP-Client mit Timeout
	client := &http.Client{
		Timeout: 5 * time.Second, // Timeout von 5 Sekunden
	}

	resp, err := client.Head(url) // HEAD-Request
	if err != nil {
		fmt.Printf("Fehler bei der URL %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Dateiname basierend auf Statuscode und Ordner
	filename := filepath.Join(outputDir, fmt.Sprintf("%d_validate_url.txt", resp.StatusCode))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Fehler beim Schreiben der Datei %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	// URL in die Datei schreiben
	_, err = file.WriteString(fmt.Sprintf("%s\n", url))
	if err != nil {
		fmt.Printf("Fehler beim Schreiben der URL %s in die Datei: %v\n", url, err)
	}
}

