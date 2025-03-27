Simple HTTP file server written in Go

## Usage

- List all files

```bash
curl -O http://localhost:8080/filetree
```

- Upload a file

```bash
curl -X POST -F "file=@/path/to/file" http://localhost:8080/upload
```

# TODO

- [ ] Add HTTPS
