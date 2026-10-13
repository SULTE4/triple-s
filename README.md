<div align="center">

# Triple-S

**`curl -X PUT localhost:8080/photos`. Upload a file. Get it back. Like S3, on your disk.**

A simplified S3-style object storage server in Go. It handles buckets and objects over a REST API,<br>
answers in S3-style XML, and keeps its metadata in CSV files.

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)
![API](https://img.shields.io/badge/API-REST%20%2B%20XML-orange)
![Storage](https://img.shields.io/badge/metadata-CSV-lightgrey)
![Dependencies](https://img.shields.io/badge/dependencies-0-brightgreen)

</div>

---

## Why

Amazon S3 is "just" buckets of files behind HTTP verbs, but a lot of design hides in that "just": name rules, metadata, content types, and what happens when you delete a bucket that still has files in it. Triple-S rebuilds that surface with only the Go standard library, to learn REST design, file I/O and networking.

## How it works

```
 curl -X PUT localhost:8080/photos/cat.jpg --data-binary @cat.jpg
        │
        ▼
 net/http mux (Go 1.22 patterns)
   PUT    /{BucketName}               GET  /
   DELETE /{BucketName}
   PUT | GET | DELETE /{BucketName}/{ObjectName}
        │
        ▼
 <dir>/
 ├── buckets.csv           name, created, last modified
 └── photos/
     ├── objects.csv       key, size, content type, last modified
     └── cat.jpg           the object itself
```

```bash
$ curl -X PUT localhost:8080/photos
<ResponseMessage><Status>201</Status><Message>Bucket created succesfully</Message></ResponseMessage>

$ curl localhost:8080/
<ListBucketResult><Buckets><Name>photos</Name><CreationDate>2026-10-13T20:40:12+05:00</CreationDate>...</Buckets></ListBucketResult>
```

## Features

- **Buckets.** Create, list and delete. Deleting a bucket that still has objects is refused with 409.
- **Objects.** Upload any bytes with a `Content-Type`, then download (served with that type) or delete them.
- **S3-like naming rules.** 3-63 characters of lowercase letters, digits, `-` and `.`, with no `..`.
- **XML responses.** Every status and listing comes back as XML, like the real S3 API.
- **Plain-file storage.** Objects are regular files, and metadata is human-readable CSV.
- **Standard library only.**

## Quick start

```bash
git clone https://github.com/SULTE4/triple-s.git
cd triple-s
go build -o triple-s ./cmd
./triple-s --port 8080 --dir data
```

Try it:

```bash
curl -X PUT localhost:8080/photos                                    # create a bucket
curl -X PUT localhost:8080/photos/hi.txt \
     -H 'Content-Type: text/plain' --data-binary 'hello'             # upload
curl localhost:8080/photos/hi.txt                                    # hello
curl localhost:8080/                                                 # list buckets
curl -X DELETE localhost:8080/photos/hi.txt                          # 204
curl -X DELETE localhost:8080/photos                                 # 204
```

## Options

| Flag | Default | Meaning |
|---|---|---|
| `--port N` | `8080` | Port, 1024-65535 |
| `--dir S` | `data` | Storage directory name next to the binary. A plain name: no `/`, no `..`, not a source folder |
| `--help` | | Show usage |

## API

| Method | Path | Success | Errors |
|---|---|---|---|
| `PUT` | `/{BucketName}` | `201` | `400` invalid name, `404` already exists |
| `GET` | `/` | `200` bucket list | |
| `DELETE` | `/{BucketName}` | `204` | `404` not found, `409` not empty |
| `PUT` | `/{BucketName}/{ObjectName}` | `201` | `404` no such bucket |
| `GET` | `/{BucketName}/{ObjectName}` | `200` object bytes | `404` not found |
| `DELETE` | `/{BucketName}/{ObjectName}` | `204` | `404` not found |

Errors look like `<ResponseMessage><Status>409</Status><Message>Bucket is not empty</Message></ResponseMessage>`.

## Project structure

```
triple-s/
├── cmd/main.go          # routes and server
├── flags/flags.go       # --port, --dir, --help
├── handlers/
│   ├── bucket.go        # bucket endpoints, XML responses
│   └── object.go        # object endpoints, objects.csv
├── metadata/            # metadata types
└── utils/utils.go       # storage dir checks, bucket name rules, usage
```

## Roadmap

- [ ] Return `409 Conflict` instead of `404` when a bucket already exists
- [ ] Allow nested paths for `--dir`
- [ ] `GET /{BucketName}` to list a bucket's objects
- [ ] Tests for the handlers
