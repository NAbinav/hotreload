package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	root := flag.String("root", ".", "Directory to watch")
	build := flag.String("build", "", "Build command")
	exec := flag.String("exec", "", "Exec command")
	flag.Parse()

	if *build == "" || *exec == "" {
		fmt.Fprintln(os.Stderr, "error: --build and --exec are required")
		os.Exit(1)
	}

	// TODO: wire up watcher
	fmt.Printf("watching: %s\n", *root)
	fmt.Printf("build:    %s\n", *build)
	fmt.Printf("exec:     %s\n", *exec)
}
