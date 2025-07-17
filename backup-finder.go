package main

import (
	"bufio"
	"flag"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"fmt"
)

var wg sync.WaitGroup

var filenames = []string{
	"backup", "website", "site", "db", "database", "config", "admin", "source_code",
	"htdocs", "web", "data", "dump",
}

var extensions = []string{
	".zip", ".tar", ".tar.gz", ".tgz", ".rar", ".7z",
	".bak", ".backup", ".old", ".gz", ".sql", ".sql.gz", ".sqlite", ".dump", ".bz2",
}

var defaultPaths = []string{
	"/", "/backup/", "/backups/", "/old/", "/admin/", "/upload/", "/uploads/", "/export/",
}

func getSubdomainPrefix(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	parts := strings.Split(host, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		os.Stderr.WriteString("Usage: " + os.Args[0] + " <subdomain_list_file>\n")
		os.Exit(1)
	}
	inputFile := flag.Arg(0)

	subs, err := readSubdomains(inputFile)
	if err != nil {
		os.Stderr.WriteString("Error reading file: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, sub := range subs {
		prefix := getSubdomainPrefix(sub)
		scanTarget(sub, prefix)
	}

	wg.Wait()
}

func readSubdomains(file string) ([]string, error) {
	var subs []string
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && (strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://")) {
			subs = append(subs, line)
		}
	}
	return subs, scanner.Err()
}

func checkURL(url string) {
	defer wg.Done()

	client := http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		fmt.Println(url)
	}
	if resp != nil {
		resp.Body.Close()
	}
}

func scanTarget(baseURL, prefix string) {
	seen := make(map[string]bool)

	targetFilenames := []string{}
	targetFilenames = append(targetFilenames, filenames...)
	targetFilenames = append(targetFilenames, prefix)
	for _, base := range filenames {
		targetFilenames = append(targetFilenames, prefix+"_"+base)
	}

	for _, path := range defaultPaths {
		if seen[path] {
			continue
		}
		seen[path] = true

		for _, name := range targetFilenames {
			for _, ext := range extensions {
				fullURL := baseURL + path + name + ext
				fullURL = strings.ReplaceAll(fullURL, "//", "/")
				if strings.HasPrefix(fullURL, "http:/") || strings.HasPrefix(fullURL, "https:/") {
					fullURL = strings.Replace(fullURL, ":/", "://", 1)
				}
				wg.Add(1)
				go checkURL(fullURL)
			}
		}
	}
}
