# Backup Finder Scanner

A fast Go-based tool to scan a list of subdomains for common website backup files.

---

## Features

- Accepts a list of subdomains with protocols (`http://` or `https://`)
- Dynamically adds subdomain name as a backup filename (e.g., `testabc.zip` for `testabc.example.com`)
- Checks common backup filenames and extensions (zip, tar.gz, bak, old, sql, etc.)
- Supports HTTP methods: `HEAD`, `GET`, or `both` (HEAD with fallback to GET)
- Outputs results to console and optionally to a file
- Uses concurrency for faster scanning

---

## Usage

### Build and Run

```bash
go run backup-finder.go probed.txt
