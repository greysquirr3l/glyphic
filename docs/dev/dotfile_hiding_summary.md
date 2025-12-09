# Dotfile Hiding in Go File Servers — Explained

## What is Dotfile Hiding?

Dotfile hiding is the practice of preventing access to files and directories that
start with a dot (`.`), such as `.env`, `.git`, or `.DS_Store`. These files are
often used for configuration, secrets, or system metadata and should not be
exposed over HTTP.

## Why is it Important?

- **Security:** Dotfiles may contain sensitive information (credentials, config,
  version control data).
- **Best Practice:** By default, Go's `http.FileServer` will serve all files,
  including dotfiles, unless you explicitly block them.
- **Compliance:** Hiding dotfiles helps meet security and privacy requirements.

## How to Hide Dotfiles in Go

The article describes a simple way to prevent dotfiles from being served by wrapping
Go's `http.FileSystem`:

```go
type noDotFileSystem struct {
    fs http.FileSystem
}

func (nfs noDotFileSystem) Open(name string) (http.File, error) {
    if strings.HasPrefix(filepath.Base(name), ".") {
        return nil, os.ErrNotExist
    }
    return nfs.fs.Open(name)
}

// Usage:
fs := noDotFileSystem{http.Dir("./public")}
http.Handle("/", http.FileServer(fs))
```

This wrapper checks if the requested file or directory starts with a dot and, if
so, returns a "not found" error, effectively hiding it from HTTP clients.

## Key Takeaways

- **Always hide dotfiles** in any Go web server that serves static files.
- This is a minimal, effective pattern that can be reused in any Go project.
- Test your server to ensure dotfiles are not accessible.

## Further Reading

- [Dot-file hiding in your Go file server (go-monk.beehiiv.com)](https://go-monk.beehiiv.com/p/dot-file-hiding-file-server)
- [Go http.FileServer documentation](https://pkg.go.dev/net/http#FileServer)
