# hpschd

The Writing-Through Mesostic Generator

A text file for input will be transmogrified into a piece of Mesostic poetry using a configured "Spine String".

## Usage

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

### Web Form

Coming soon.

(NB: 'Multi-Part File Upload' was removed in `v1.4.3`.)

## Operations

### Tag a Release

1. Commit, push, merge, and update main as necessary, include the new version in `Dockerfile`.
2. In main: `git tag vX.Y.Z; git push --tags origin`

### Container Build for DockerHub

Once `Dockerfile` contains the new version:

```zsh
docker build -t chaquo:hpschd .
docker tag chaquo:hpschd docker.io/maroda/chaquo:hpschd
docker push docker.io/maroda/chaquo:hpschd
```

### Continuous Build for DockerHub

Currently testing continuous registry publication with hub.docker.com + github.

### Run Docker Locally

Fetch the `latest` version and run as a local container:

```
docker run --rm --name hpschd -p 9999:9999 maroda/chaquo:hpschd
```

### Running on AWS ECS Fargate

Once the container is in the registry it can be launched/updated on AWS ECS Fargate.

### Updating Fargate

The **hpschd-mesostic** cluster runs the task **mesostic** that will need to have a new revision created and loaded that downloads the new version of hpschd from DockerHub.

## Mesostics

Of course not all mesostics are "writing through" style as Cage did often, they can just as easily be written as they are.

- 50% Mesostic: The CL is unique between itself and the previous CL.
- 100% Mesostic: The CL is unique between itself, the previous CL, and the next CL.
- A "meso-acrostic", arguably another version of a Mesostic, has neither of these limitations.

## Chance Operations

### Current Implementation

There is a two-phase operation:

1. Under certain conditions the engine will obtain new Mesostics by creating a randomized date string and requesting the APOD from that date.
2. These are stored locally in the special runtime directory **`store/`**.
3. When the homepage is requested by a visitor, a random selection from the runtime directory is chosen to display.

### I Ching

There are probably dozens if not hundreds of computer programs that simulate the I Ching.

So this doesn't mean to replicate them but to provide a source of randomness for calculating values of the Mesostic that is in line with the kind of approach Cage might do. For instance:

- The property of how many words per line could be selected via chance operations.
- The SS itself could be chance derived.
- Parameters for if display fonts are ever used.


## Complexity

How could you demonstrate complexity and chaos here?


## Resources

There are other precedents:

- Nicki Hoffman (python) ::: http://vyh.pythonanywhere.com/psmeso/
- UPenn team (javascript) ::: http://mesostics.sas.upenn.edu/

