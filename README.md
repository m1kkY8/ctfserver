Simple HTTP file server written in Go

## Usage

- List all files

```bash
curl -O localhost:8080/filetree
```

- Upload a file

```bash
curl -F "file=@/path/to/file" localhost:8080/upload
```

# TODO

- [ ] Add HTTPS
