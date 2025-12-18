# Mesostic poetry API

[![Release](https://github.com/maroda/hpschd/actions/workflows/release.yml/badge.svg)](https://github.com/maroda/hpschd/actions/workflows/release.yml)

**The Writing-Through Mesostic Generator**

A text file for input will be transmogrified into a piece of Mesostic poetry using a configured "Spine String".

## Usage

The default webserver shows a mesostic built from fetches made to the [NASA Astronomy Picture of the Day (APOD)](https://apod.nasa.gov/apod/) API.
This uses the following API endpoint, which can be used for any blob of text.

### JSON API

```zsh
curl www.hpschd.xyz:9999/app -d '{"text": "the quick brown\nfox jumps over\nthe lazy dog\n", "spinestring": "cra"}'
```

For example:

```zsh
>>> curl localhost:9999/app -d '{"text": "the quick brown\nfox jumps over\nthe lazy dog\n", "spinestring": "cra"}'
      the quiCk b
fox jumps oveR
        the lAzy dog
```

## Operations

To run this and display a Mesostic on the homepage, you will need an APOD API Key.
Visit [NASA's API pages](https://api.nasa.gov) to sign up and get a free key.

### Run Docker Locally

First set your APOD API key in the environment. If this is not set, it will default to NASA's test key: `DEMO_KEY`

```zsh
export NASA_API_KEY=<KEY>
```

Fetch the `latest` version from GitHub Container Registry and run as a local container:
```zsh
docker run --rm --name hpschd -p 9999:9999 ghcr.io/maroda/hpschd:latest
```

Now browse to <http://localhost:9999> and see an APOD mesostic!

### Docker Compose

Use `docker compose up` with the following `compose.yaml` entry:
```yaml
services:
  hpschd:
    image: ghcr.io/maroda/hpschd:latest
    container_name: hpschd
    ports:
      - "9999:9999"
    environment:
      - NASA_API_KEY=<KEY>
  restart: unless-stopped
```

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

## Mesostics

An acrostic shows a String of letters down one side of the text.
A **mesostic** shows a String of letters down the _middle_ of the text.

This vertical line of text is capitalized and centered, so we call that a _Spine String_.
Locating which letter in a line of entry-text to center on the Spine String happens using one of three algorithms:

1. **50% Mesostic**: The Spine String Letter is unique between itself and the previous one. So if this letter is K, there cannot be a K between itself and the previous letter (which also would not be a K).
2. **100% Mesostic**: The Spine String Letter is unique between itself, the previous one, and the next one. In the example, the letter K cannot exist before _or after_ itself between its partner letters.
3. **A "meso-acrostic"**: Neither of these limitations.

John Cage would run large amounts of text through a Mesostic algorithm to create poetry.
The entry-text forms the lines of poetry and the Spine String (our term) forms the vertical letters down the middle.

> The **50% Mesostic** is what _hpschd_ uses to produce the most output from small blocks of text, like the APOD descriptions.
> 
> The algorithm is fuzzy and can lead to characters causing weird shifts, very long lines, or empty space.
>
> This is intentional.

## Chance Operations
There is a two-phase operation:

1. Under certain conditions the engine will obtain new Mesostics by creating a randomized date string and requesting the APOD from that date.
2. These are stored locally in the special runtime directory **`/store/`**.
3. When the homepage is requested a random selection from the runtime directory is chosen to display.

This way the visitor is never waiting on the fetch itself, and will always get something that has been previously fetched.
This means repeats will happen, but the more time the app runs to make new fetches, the more are saved in the cache.

## Other Implementations

Mesostic creation algorithms in the wild!

- Nicki Hoffman (python) ::: http://vyh.pythonanywhere.com/psmeso/
- UPenn team (javascript) ::: http://mesostics.sas.upenn.edu/

## Acknowledgements

This algorithm was written by a human (me), but the "way to write a Mesostic" was mostly work by **John Cage**.
Claude Code has been used to assist with testing and refactoring.

