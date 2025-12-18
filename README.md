# Mesostic poetry API

[![Release](https://github.com/maroda/hpschd/actions/workflows/release.yml/badge.svg)](https://github.com/maroda/hpschd/actions/workflows/release.yml)

**The Writing-Through Mesostic Generator**

A text file for input will be transmogrified into a piece of Mesostic poetry using a configured "Spine String".
This algorithm was written by a human (me), but the "way to write a Mesostic" was originated by **John Cage**.

> Claude Code has been used to assist with testing and refactoring.

## Usage

The default webserver shows a mesostic built from fetches made to the NASA Astronomy Picture of the Day (APOD) API.
This uses the same API endpoint that can be used for any blob of text.

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
### Run Docker Locally

Set your APOD API key in the environment.

Fetch the `latest` version from GitHub Container Registry and run as a local container:
```zsh
docker run --rm --name hpschd -p 9999:9999 ghcr.io/maroda/hpschd:latest
```

Or run a specific version:
```zsh
docker run --rm --name hpschd -p 9999:9999 ghcr.io/maroda/hpschd:v1.5.0
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
A mesostic shows a String of letters down the middle of the text.
Usually this vertical line of text is capitalized, here we'll call that a Spine String.

Locating which letter in a line of text to center on the Spine String comes in three forms:

1. 50% Mesostic: The Spine String Letter is unique between itself and the previous one. So if this letter is K, there cannot be a K between itself and the previous letter (which also would not be a K).
2. 100% Mesostic: The Spine String Letter is unique between itself, the previous one, and the next one. In the example, the letter K cannot exist before _or after_ itself between its partner letters.
3. A "meso-acrostic", arguably another version of a Mesostic, has neither of these limitations.

John Cage would run large amounts of text through the Mesostic algorithm to create poetry.
The original text forms the lines of poetry and the Spine String (my term) forms the vertical letters down the middle.

> The **50% Mesostic** is what _hpschd_ displays by default.
> The algorithm is fuzzy and can lead to characters causing weird shifts, very long lines, or empty space.
> This is intentional.

See also the Wikipedia page on [Mesostic](https://en.wikipedia.org/wiki/Mesostic)

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
