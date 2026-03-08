# hotreload

A CLI tool that watches a Go project for file changes and automatically rebuilds and restarts the server.

## Usage

```bash
./hotreload --root <folder> --build "<build command>" --exec "<run command>"
```

## Features

### Core

- [x] Triggers first build on start
- [x] Watches for file changes
- [x] Rebuilds automatically
- [x] Restarts server automatically
- [x] Real-time log streaming

### Bonus

- [x] Recursive subdirectory watching
- [x] Detects new folders at runtime
- [x] Debounces rapid saves (500ms)
- [x] Kills entire process group (no leaked children)
- [x] Graceful shutdown (SIGTERM → 5s → SIGKILL)
- [x] Crash loop protection
- [x] Ignores `.git/`, `node_modules/`, editor temps, etc.
- [x] inotify watch limit warning

### Flags

| Flag      | Description                    |
| --------- | ------------------------------ |
| `--root`  | Directory to watch (recursive) |
| `--build` | Command to build the project   |
| `--exec`  | Command to run the server      |

## Demo

```bash
make demo
```

This builds hotreload and runs it against the included `testserver/`.

## Features

- Triggers first build immediately on start
- Watches all subdirectories recursively
- Detects and watches new folders created at runtime
- Debounces rapid file saves (500ms window) to avoid redundant rebuilds
- Streams server logs in real time
- Kills the entire process group on restart — no leaked child processes
- Graceful shutdown: SIGTERM → 3s wait → SIGKILL
- Crash loop protection: backs off if server exits within 1 second of starting
- Warns when OS inotify watch limit is reached
- Ignores `.git/`, `node_modules/`, `bin/`, `vendor/`, and editor temp files (`*.swp`, `*~`, `#*`)

## Project Structure

```
hotreload/
├── main.go           # entry point, flag parsing, wiring
├── watcher/          # recursive fsnotify wrapper
├── runner/           # build + process management
├── debounce/         # debounce utility
└── testserver/       # sample HTTP server for demo
```

## Build & Run

```bash
make demo   # builds and runs against testserver
make build  # just builds the binary
```

### Loom video link https://www.loom.com/share/7fc22b3f8bac42b796407ad81daa8f53
