package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/halimath/depot/internal/generate"
)

var (
	tableName   = flag.String("table", "", "Name of the database table")
	repoName    = flag.String("repo", "", "Name of the repo")
	repoPackage = flag.String("repo-package", "", "Package name for the repo")
	readOnly    = flag.Bool("ro", false, "Generate a read-only repo")
)

// Usage: depot generate --table=messages ./test.go Message

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "%s: missing args\n", os.Args[0])
		os.Exit(1)
	}

	source, err := generate.GenerateRepository(generate.Options{
		Filename:    flag.Arg(0),
		EntityName:  flag.Arg(1),
		TableName:   *tableName,
		RepoPackage: *repoPackage,
		RepoName:    *repoName,
		ReadOnly:    *readOnly,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error generating repository: %s\n", os.Args[0], err)
		os.Exit(2)
	}

	os.Stdout.Write(source)
}
