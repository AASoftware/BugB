package main

import (
	"bufio"
	"fmt"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Struktur zur Speicherung der URL-Parameter
type ParsedURL struct {
	Base string
	X    float64
	Y    float64
	Raw  string
}

// Funktion, um eine URL zu parsen und relevante Teile zu extrahieren
func parseURL(rawURL string) (ParsedURL, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ParsedURL{}, err
	}

	query := parsed.Query()
	x, err := strconv.ParseFloat(query.Get("x"), 64)
	if err != nil {
		x = 0
	}
	y, err := strconv.ParseFloat(query.Get("y"), 64)
	if err != nil {
		y = 0
	}

	return ParsedURL{
		Base: fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, parsed.Path),
		X:    x,
		Y:    y,
		Raw:  rawURL,
	}, nil
}

// Funktion, um URLs zu vergleichen
func isSimilar(url1, url2 ParsedURL, threshold float64) bool {
	return url1.Base == url2.Base &&
		math.Abs(url1.X-url2.X) < threshold &&
		math.Abs(url1.Y-url2.Y) < threshold
}

// Levenshtein berechnet den Levenshtein-Abstand zwischen zwei Strings
func Levenshtein(a, b string) int {
	lenA, lenB := len(a), len(b)
	dist := make([][]int, lenA+1)
	for i := range dist {
		dist[i] = make([]int, lenB+1)
	}

	for i := 0; i <= lenA; i++ {
		dist[i][0] = i
	}
	for j := 0; j <= lenB; j++ {
		dist[0][j] = j
	}

	for i := 1; i <= lenA; i++ {
		for j := 1; j <= lenB; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			dist[i][j] = min(
				dist[i-1][j]+1,
				dist[i][j-1]+1,
				dist[i-1][j-1]+cost,
			)
		}
	}

	return dist[lenA][lenB]
}

// min gibt den kleineren Wert von drei Integern zurück
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func main() {
	// Überprüfen, ob ein Argument übergeben wurde
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./prepTsc <input-file>")
		fmt.Println("\nBeispiel: ./prepTsc urls.txt")
		return
	}

	// Dateiname aus dem Argument lesen
	inputFileName := os.Args[1]

	// Datei mit URLs öffnen
	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println("Fehler beim Öffnen der Datei:", err)
		return
	}
	defer file.Close()

	// URLs parsen und filtern
	var urls []ParsedURL
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rawURL := strings.TrimSpace(scanner.Text())
		parsed, err := parseURL(rawURL)
		if err != nil {
			fmt.Println("Fehler beim Parsen der URL:", err)
			continue
		}
		urls = append(urls, parsed)
	}

	// Schritt 1: Ähnliche URLs entfernen
	var filtered []ParsedURL
	threshold := 50.0 // Setze threshold für Genauigkeit!
	for _, url := range urls {
		isDuplicate := false
		for _, kept := range filtered {
			if isSimilar(url, kept, threshold) {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			filtered = append(filtered, url)
		}
	}

	// Schritt 2: URLs mit denselben ersten 30 Zeichen filtern
	var filteredByPrefix []string
	seen := make(map[string]bool)
	for _, url := range filtered {
		// Hole die ersten 30 Zeichen der URL
		if len(url.Raw) < 30 {
			continue
		}
		prefix := url.Raw[:30]

		// Wenn dieser Prefix noch nicht gesehen wurde, füge die URL hinzu
		if !seen[prefix] {
			filteredByPrefix = append(filteredByPrefix, url.Raw)
			seen[prefix] = true
		}
	}

	// Schritt 3: Levenshtein-Abstand berechnen und URLs hinzufügen, wenn der Abstand > 50 ist
	var finalFiltered []string
	for _, url := range filteredByPrefix {
		shouldAdd := true
		for _, seenURL := range finalFiltered {
			if Levenshtein(url, seenURL) == 50 { // Nur exakt gleiche URLs
				shouldAdd = false
				break
			}
		}
		if shouldAdd {
			finalFiltered = append(finalFiltered, url)
		}
	}

	// Gefilterte URLs in eine neue Datei schreiben
	outputFileName := "prepared_for_techscan.txt"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Fehler beim Erstellen der Ausgabe-Datei:", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for _, url := range finalFiltered {
		_, _ = writer.WriteString(url + "\n")
	}
	writer.Flush()

	fmt.Printf("Gefilterte URLs wurden in '%s' gespeichert. (%d URLs)\n", outputFileName, len(finalFiltered))
}
