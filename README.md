# hyphen

[![Go Report](https://goreportcard.com/badge/github.com/nekomeowww/hyphen)](https://goreportcard.com/report/github.com/nekomeowww/hyphen)
[![Testing](https://github.com/nekomeowww/hyphen/actions/workflows/ci.yml/badge.svg)](https://github.com/nekomeowww/hyphen/actions/workflows/ci.yml)
[![Building](https://github.com/nekomeowww/hyphen/actions/workflows/build.yml/badge.svg)](https://github.com/nekomeowww/hyphen/actions/workflows/build.yml)

---

An elegant URL Shortener service

## Usage

### Run

```shell
hyphen -l 0.0.0.0:9010 --data ./data/bbolt_data.db
```

### Run with Docker

```shell
docker run -it --rm -p 9010:9010 -v <path/to/bbolt/db/data>:/var/lib/hyphen nekomeowww/hyphen:latest
```

## Build on your own

### Build with go

Remember to specify your version with `-ldflags " -X 'github.com/nekomeowww/hyphen/pkg/meta.Version=<version to be specified>'"`

```shell
go build -a -o "release/hyphen" -ldflags " -X 'github.com/nekomeowww/hyphen/pkg/meta.Version=<version>'" "github.com/nekomeowww/hyphen/cmd/hyphen"
```

### Build with Docker

Remember to specify your version with `--build-arg VERSION=<version to be specified>`

```shell
docker buildx build --platform <your/arch> -t <tag> --build-arg VERSION=<version> . -f Dockerfile
```

## API

### `GET` `/api/v1/url/full`

Query Params:

- `url` (required): Full URL

Response:

```json
{
  "shortUrl": "2f224d525e"
}
```

### `GET` `/api/v1/url/short`

Query Params:

- `url` (required): Short URL

Response:

```json
{
  "url": "https://example.com/8b4e9413fd13f4a83a2a31c8494347"
}
```

### `POST` `/api/v1/url`

Body Params:

```json
{
    "url": "https://example.com/8b4e9413fd13f4a83a2a31c8494347"
}
```

Response:

```json
{
    "url": "https://example.com/8b4e9413fd13f4a83a2a31c8494347",
    "shortUrl": "2f224d525e",
}
```

### `DELETE` `/api/v1/url`

Body Params:

```json
{
    "shortUrl": "2f224d525e"
}
```

Response:

```json
{
    "ok": true
}
```
