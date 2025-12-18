# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

hpschd is a "Writing-Through Mesostic Generator" that transforms text into Mesostic poetry using a configurable "Spine String". It's a Go-based web service that implements John Cage-inspired chance operations and algorithmic poetry generation.

### What is a Mesostic?

A Mesostic is a form of poetry where a "spine string" runs vertically through the text. Each line contains one character from the spine string (capitalized), with specific rules:
- **50% Mesostic**: The spine character is unique between itself and the previous spine character
- **100% Mesostic**: The spine character is unique between itself, the previous, and next spine characters
- **Meso-Acrostic**: No uniqueness limitations

## Architecture

### Core Components

1. **Mesostic Engine** (`mesostic.go`): The algorithmic heart of the application
   - `mesoLine()`: Processes each line to find spine string characters, operating in modes (0=WestSide building, 1=EastSide with 50% rules)
   - `Spine()`: Converts spine string into rotatable character slice
   - `Ictus()` / `Preus()`: Rotate spine string position forward/backward through text
   - Uses a global `fragMents` hash table to store line fragments before sorting
   - Fragment keys are SHA1 hashes for consistent sizing

2. **API Layer** (`api.go`): HTTP endpoints
   - `/`: Homepage displays NASA APOD mesostics (random selection from `store/`)
   - `/app`: JSON POST endpoint for generating custom mesostics
   - `/app/{arg}`: Form POST endpoint (spine string in URL path)
   - `/ping`: Readiness check
   - `/metrics`: Prometheus metrics endpoint

3. **NASA APOD Integration** (`cron.go`, `fetch.go`): Automated content generation
   - Cron job (`fetchCron()`) runs every 666 seconds by default
   - Fetches NASA Astronomy Picture of the Day metadata via API
   - Uses title as spine string, explanation text as source material
   - Stores generated mesostics in `store/` directory
   - If current date unavailable (404) or mesostic exists, triggers randomized date fetch
   - Uses buffered channel `nasaNewMESO` to communicate new mesostic availability

4. **Data Operations** (`dataops.go`): Filesystem and utilities
   - `ichingMeso()`: Chance-based selection of existing mesostics
   - `rndDate()`: Generates random dates (2000-2020 range) for APOD queries
   - `fileTmp()`: Creates temporary files in `txrx/` for processing
   - `apodNew()`: Writes mesostics to `store/` with deduplication

5. **Observability** (`obvy.go`): Metrics and tracing
   - Prometheus histograms for function runtimes
   - `Envelope()`: Returns current execution context (file, line, function) for structured logging

### Data Flow

1. User submits text + spine string via JSON API → `JSubmit()`
2. Text written to temp file in `txrx/` via `fileTmp()`
3. `mesoMain()` launched in goroutine with result channel
4. Lines processed sequentially through `mesoLine()`:
   - Finds spine character, builds west/east fragments
   - Stores in `fragMents` map with SHA1 key
   - Rotates spine position with `Ictus()` or `Preus()`
5. Fragments sorted by LineNum and formatted with left-padding
6. Result returned through channel, temp file deleted

### Runtime Directories

Both directories are auto-created on startup by `localDirs()` in main.go:48:
- `store/`: Ephemeral cache of NASA APOD mesostics (format: `YYYY-MM-DD__Title_With_Underscores`)
- `txrx/`: Temporary scratch files for processing (cleaned after use)

## Development Commands

### Build and Run

```bash
# Build locally
go build -o hpschd

# Run with default settings (NASA APOD fetching enabled)
# On startup, app will fetch one APOD mesostic before starting web server
./hpschd

# Run with debug logging
./hpschd -debug

# Run without NASA APOD cronjob (store will remain empty)
./hpschd -nofetch

# Run tests
go test ./...

# Run specific test
go test -v -run TestFunctionName
```

### Docker

Images are published to GitHub Container Registry (ghcr.io) with semantic versioning.

```bash
# Build container for LOCAL TESTING (builds binary inside container)
docker build -f Dockerfile.local -t hpschd:latest .

# Run locally
docker run --rm --name hpschd -p 9999:9999 hpschd:latest

# Optional: Mount volume for persistent store/ directory
docker run --rm --name hpschd -p 9999:9999 -v hpschd-store:/app/store hpschd:latest

# Run published image from GitHub Container Registry
docker run --rm --name hpschd -p 9999:9999 ghcr.io/maroda/hpschd:latest
docker run --rm --name hpschd -p 9999:9999 ghcr.io/maroda/hpschd:v1.2.3  # specific version

# Production build (requires goreleaser to build binary first)
# This is automated by CI/CD on git tag push
docker build -t ghcr.io/maroda/hpschd:latest .

# Manual tag and push to GitHub Container Registry
docker tag hpschd:latest ghcr.io/maroda/hpschd:v1.2.3
docker tag hpschd:latest ghcr.io/maroda/hpschd:latest
docker push ghcr.io/maroda/hpschd:v1.2.3
docker push ghcr.io/maroda/hpschd:latest
```

### Testing the API

```bash
# JSON API (generate mesostic)
curl localhost:9999/app -d '{"text": "the quick brown\nfox jumps over\nthe lazy dog\n", "spinestring": "cra"}'

# Expected output:
#      the quiCk b
#fox jumps oveR
#        the lAzy dog

# Readiness check
curl localhost:9999/ping

# View metrics
curl localhost:9999/metrics
```

## Key Implementation Notes

### Global State
- `fragMents` map: Must be global for the Ictus rotation display to work correctly (noted in code comments)
- `fragCount`: Global counter for line fragment combinations
- `nasaNewMESO`: Buffered channel (capacity 1) for communicating new NASA mesostic filenames
  - Buffered to prevent deadlock during synchronous initial startup fetch

### NASA APOD Cronjob Timing
- Default: 666 seconds (~11 minutes)
- Rationale: Avoids API rate limits (1k/hr), prevents "EXISTENT" mesostic pileup from fast fetches
- TODO noted in code: Check for existing mesostics BEFORE generating (currently generates then checks)

### Concurrency Patterns
- `mesoMain()` runs in goroutine, returns results via channel
- NASA ETL triggers recursive randomized fetches on 404 or duplicate detection
- Homepage reads from `nasaNewMESO` channel with non-blocking select (returns "HPSCHD" signal when empty)

### Startup Sequence
- On startup (unless `-nofetch` flag is used), app synchronously fetches one NASA APOD mesostic before starting the web server
- This ensures `store/` directory is populated and prevents ENOENT errors on first homepage load
- After initial fetch completes, the cronjob scheduler starts for subsequent fetches
- If initial fetch fails (404 or network error), `NASAetl()` recursively tries a random historical date

### Dependencies
- **gorilla/mux**: HTTP routing
- **zerolog**: Structured logging
- **gocron**: Scheduler for NASA APOD fetching
- **prometheus/client_golang**: Metrics collection

## Release Process

Releases are automated via goreleaser when tags are pushed.

1. Commit, push, merge to main
2. Tag release with semantic version: `git tag vX.Y.Z && git push origin vX.Y.Z`
3. CI/CD automatically builds and pushes Docker images to GitHub Container Registry:
   - `ghcr.io/maroda/hpschd:vX.Y.Z` (specific version)
   - `ghcr.io/maroda/hpschd:latest` (updated on each release)
4. Update AWS ECS Fargate task **mesostic** with new image revision on **hpschd-mesostic** cluster
