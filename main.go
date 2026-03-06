package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/NAbinav/hotreload/debounce"
	"github.com/NAbinav/hotreload/runner"
	"github.com/NAbinav/hotreload/watcher"
)

func main() {
	root := flag.String("root", ".", "directory to watch")
	build := flag.String("build", "", "build command")
	execCmd := flag.String("exec", "", "run command")
	flag.Parse()

	if *build == "" || *execCmd == "" {
		fmt.Fprintln(os.Stderr, "usage: hotreload --root DIR --build CMD --exec CMD")
		os.Exit(1)
	}

	w, err := watcher.New(*root)
	if err != nil {
		fmt.Println("watcher error:", err)
		os.Exit(1)
	}
	defer w.Close()

	rebuild := debounce.New(500*time.Millisecond, func() {
		runner.Run(*build, *execCmd)
	})

	runner.Run(*build, *execCmd)

	w.Watch(func() {
		rebuild()
	})
}
