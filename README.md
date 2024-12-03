# BugB
BugBounty Tools

Wichtig - Wird laufend erweitert!

# Erleichterung für Recon Aufgaben

Eine Sammlung von Befehlen und Tools zur Unterstützung bei Sicherheits- und Recon-Aufgaben.

---

## Golang Installation

Installiere Golang mit folgendem Befehl:

```bash
sudo apt install golang
```

---

## Subdomains finden mit Amass

Nutze Amass, um Subdomains zu entdecken und in einer Datei zu speichern:

```bash
amass enum -passive -d example.com > subdomains.txt
```


---

## URL-Historie extrahieren mit Waybackurls

Nutze das Tool Waybackurls, um die URL-Historie einer Domain zu sammeln:

```bash
# Installation
go install github.com/tomnomnom/waybackurls@latest

# Nutzung
cat subdomains.txt | waybackurls > urls.txt
```

---

## Technologie-Erkennung mit HTTPX

Scanne URLs, um Technologien und Header-Informationen zu identifizieren:

```bash
# Installation
go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest

# Nutzung
httpx -l urls.txt -td | tee httpx_tech.txt
```

---

## Resultate filtern

### Füge Anführungszeichen in jede Zeile ein

```bash
sed 's/.*/"&"/' input_file > output_file
```

### Entferne Zeilen mit einem bestimmten Wort und speichere das Ergebnis in einer neuen Datei

```bash
sed '/amtsblatt/d' input.txt > output.txt
```

### Entferne Zeilen mit einem bestimmten Wort und überschreibe die Originaldatei

```bash
sed -i '/amtsblatt/d' input.txt
```

### Filtere URLs mit spezifischen Stichwörtern und speichere sie in einer neuen Datei

```bash
grep -E "\(=\|?\|&|\.js|\.php|\.json|\.xml|login|admin|password|passwort|form|formular|upload|anmeldung|kontakt|user|benutzer|id|\.txt|doku|documentation|dokumentation|wiki|geheim|confidential|debug|debugger|test|beta|security|robots|policy|config|configuration|auth|session|token|api|oauth|cert|sso|staging|dev|log|backup|dump|panel|dashboard|setup|install|private|hidden|secret|key|db|passwd\)" mutiert1.txt > stark_gefiltert.txt
```

---

## Head Request Scan mit `headScan`

Scanne URLs, filtere den Output und erstelle separate Dateien basierend auf Statuscodes:

```bash
# 1. Klone das Repository
git clone https://github.com/AASoftware/BugB.git

# 2. Wechsle in den Ordner mit der Go-Datei
cd BugB/tools/golang/source/headScan

# 3. Initialisiere das Go-Modul
go mod init github.com/AASoftware/BugB/tools/golang/source/headScan

# 4. Lade alle Abhängigkeiten
go mod tidy

# 5. Baue die Go-Datei
go build headScan.go

# 6. Führe das erstellte Binary aus
./headScan <sourcefile>
```

---

## Filterung von URLs für Technologie-Scans mit `prepTsc`

Bereite URLs für einen Technologie-Scan mit HTTPX vor:

```bash
# 1. Klone das Repository
git clone https://github.com/AASoftware/BugB.git

# 2. Wechsle in den Ordner mit der Go-Datei
cd BugB/tools/golang/source/prepTsc

# 3. Initialisiere das Go-Modul
go mod init github.com/AASoftware/BugB/tools/golang/source/prepTsc

# 4. Lade alle Abhängigkeiten
go mod tidy

# 5. Baue die Go-Datei
go build prepTsc.go

# 6. Führe das erstellte Binary aus
./prepTsc <sourcefile>
```

---

## License

Dieses Projekt steht unter der MIT-Lizenz.
