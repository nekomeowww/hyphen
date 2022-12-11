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

### Run with docker-compose

```shell
docker-compose up -d
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

Query one short URL by original URL.

Query Params:

- `url` `String`  (required): original URL

Response:

```json
{
  "data": {
    "shortUrl": "2f224d525e"
  }
}
```

### `GET` `/api/v1/url/short`

Query one original URL by short URL.

Query Params:

- `url` `String` (required): Short URL
- `redirect` `Boolean`: Whether to redirect to the original URL automatically (by returning `301 Permanently Moved`)

Response:

```json
{
  "data": {
    "url": "https://example.com/8b4e9413fd13f4a83a2a31c8494347"
  }
}
```

### `POST` `/api/v1/url`

Create a new URL. The short URL will be created based on the first 10 letter of sha512 hash of
the original URL.

Body Params:

```json
{
    "url": "https://example.com/8b4e9413fd13f4a83a2a31c8494347"
}
```

Response:

```json
{
  "data": {
    "url": "https://example.com/8b4e9413fd13f4a83a2a31c8494347",
    "shortUrl": "2f224d525e",
  }
}
```

### `DELETE` `/api/v1/url`

Revoke the short URL that created in the past.

**NOTICE: since `hyphen` is creating short URL base on the sha512 hash of the original URL,
revoke will disconnect all the original URL to the short URL that revoked, however, you could
create a new one to reconnect them together, it will return the same short URL as before.**

Body Params:

```json
{
    "shortUrl": "2f224d525e"
}
```

Response:

```json
{
  "data": {
    "ok": true
  }
}
```
