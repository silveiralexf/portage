# gomakeme

CI/CD wrapper for building stuff up -- still in progress.

## Docker Images

Latest images are available on [Docker Hub](https://hub.docker.com/u/silveiralexf/gomakeme)

## How to Run?

```sh
DOCKER_DEFAULT_PLATFORM=linux/amd64 go run main.go golang
```

## How to Build?

To build the project from source, first thing you'll need is a working Go environment with version 1.20 
or greater installed, as further described by the official docs:

- [Go: Download & Install](https://go.dev/doc/install)

### From source

```sh
go build -o gomakeme
```

### With GoMakeMe itself

Using Go locally:

```sh
go run main.go golang \ 
    --name gomakeme \  
    --tag latest \ 
    --output "binary,oci" \ 
    --from ubuntu:jammy \ 
    --runtime scratch \ 
    --build-arg FOO=bar \ 
    --build-arg BAR=foo
```

Or from a published docker image:

```sh
 docker run silveiralexf/gomakeme:latest golang \ 
    --name "gomakeme" \ 
    --tag latest \ 
    --output "binary,oci" \ 
    --from ubuntu:jammy \ 
    --runtime scratch \ 
    --build-arg FOO=bar \ 
    --build-arg BAR=foo
```
