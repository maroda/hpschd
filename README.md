# Mesostic Poetry Generator

[![Release](https://github.com/maroda/hpschd/actions/workflows/release.yml/badge.svg)](https://github.com/maroda/hpschd/actions/workflows/release.yml)

**The Golang Writing-Through Mesostic Engine**

## JSON API

The Mesostic Engine can be used directly via the API:
```shell
curl localhost:9876/app -d '{
  "text": "Any long string of text that will be the content of the poem",
  "spinestring": "The capitalized text that runs through the middle"
}'
```

For example:
```zsh
>>> curl -s localhost:9876/app -d '{"text": "the quick brown; fox jumps over; the lazy dog", "spinestring": "cra"}' | jq
{
  "title": "cra",
  "date": "",
  "poem": "\n      the quiCk b\nfox jumps oveR\n        the lAzy dog"
}
```

Printed out:
```text
      the quiCk b
fox jumps oveR
        the lAzy dog
```

## Homepage

This is a simple webserver that displays a mesostic built from the [NASA Astronomy Picture of the Day (APOD)](https://apod.nasa.gov/apod/) API.
It retrieves its content from a local datastore of mesostics.
This store is filled by an automated process that queries the APOD API for a random date.

## Operations

To run this and display a Mesostic on the homepage, you will need an APOD API Key.
Visit [NASA's API pages](https://api.nasa.gov) to sign up and get a free key.

### Environment Variables
- `NASA_API_KEY` should be set to your APOD API key. If not set, the app defaults to NASA's test key ("DEMO_KEY") that has heavier rate limiting than a normal user.
- `HPSCHD_APOD_FREQUENCY` can be set to a duration in seconds (default is "88").


Place this in `./.env`:
```zsh
NASA_API_KEY=<KEY>
HPSCHD_APOD_FREQUENCY=88
```

### Run Docker Locally
Fetch the `latest` version and run as a local container:
```zsh
docker run --env-file ./.env --rm --name hpschd -p 9876:9876 ghcr.io/maroda/hpschd:latest
```

Now browse to <http://localhost:9876> and see an APOD mesostic!

### Docker Compose

Use `docker compose up` with the following `compose.yaml` entry:
```yaml
services:
  hpschd:
    image: ghcr.io/maroda/hpschd:latest
    container_name: hpschd
    ports:
      - "9876:9876"
    environment:
      - NASA_API_KEY=<KEY>
      - HPSCHD_APOD_FREQUENCY=88
  restart: unless-stopped
```

### Disable APOD Fetching
The API can be used alone by running with the `-nofetch` flag.
This disables new APOD fetches, but does not remove any existing ones,
so visits to the homepage will still show a random mesostic from the current APOD collection.

### Release Process

The project uses GitHub Actions with GoReleaser for automated releases:

1. Commit, push, merge, and update main as necessary.
2. In main: `git tag -a vX.Y.Z -m"comment" && git push --tags origin`
3. The release workflow automatically:
   - Creates a draft GitHub release
   - Builds Linux binaries (amd64 and arm64)
   - Builds and publishes multi-arch Docker images to GitHub Container Registry (ghcr.io)
   - Publishes the release with artifacts and changelog

**What gets published:**
- GitHub release with binary archives (tar.gz)
- Docker images at `ghcr.io/maroda/hpschd:latest` and `ghcr.io/maroda/hpschd:vX.Y.Z`
- Multi-architecture support (amd64 and arm64)

See the [Release workflow](.github/workflows/release.yml) and [GoReleaser config](.goreleaser.yaml) for details.

## What are Mesostics?
An _acrostic_ shows a String of letters down one side of the text.

A **mesostic** shows a String of letters down the _middle_ of the text.

> The [Mesostic](https://en.wikipedia.org/wiki/Mesostic) was first developed extensively by **John Cage**.
> His book [**M**](https://en.wikipedia.org/wiki/M_(John_Cage_book)) contains many of them,
> including a favorite "62 Mesostics re Merce Cunningham" where
> individual letters appear in typeface styles derived by chance operations
> and performed as music by a singer.

### Mesostic Algorithms
The vertical line of text is capitalized and centered, we call that a _Spine String_.
Locating which letter in a line of entry-text to center on the Spine String happens using one of three algorithms:

1. **50% Mesostic**: The Spine String Letter is unique between itself and the previous one. So if this letter is K, there cannot be a K between itself and the previous letter (which also would not be a K).
2. **100% Mesostic**: The Spine String Letter is unique between itself, the previous one, and the next one. In the example, the letter K cannot exist before _or after_ itself between its partner letters.
3. **A "meso-acrostic"**: Neither of these limitations.

John Cage would run large amounts of text through a Mesostic algorithm to create poetry.
The entry-text forms the lines of poetry and the Spine String (our term) forms the vertical letters down the middle.

### HPSCHD Mesostics
The **50% Mesostic** is used to produce good output from small blocks of text, like the APOD descriptions.

To create more interesting poetry, some text is discarded and traded for whitespace.

> _In future versions there may be configuration options to follow more strict Mesostic rules._

## Release Notes
- There may be bugs in the way the API consumes text.
For instance, embedded control characters may create unpredictable results.
- No controllable mesostic options (like rule strictness) yet.

## Other Implementations

Mesostic creation algorithms in the wild!
These have controllable options like word density and rule selection.

- Nicki Hoffman (python) ::: http://vyh.pythonanywhere.com/psmeso/
- UPenn team (javascript) ::: http://mesostics.sas.upenn.edu/

