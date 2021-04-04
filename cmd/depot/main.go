package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/halimath/depot/internal/generate"
)

var (
	tableName   = flag.String("table", "", "Name of the database table")
	repoName    = flag.String("repo", "", "Name of the repo")
	repoPackage = flag.String("repo-package", "", "Package name for the repo")
	readOnly    = flag.Bool("ro", false, "Generate a read-only repo")
	out         = flag.String("out", "", "Filename to write output to (defaults to STDOUT)")
)

// Usage: depot generate-repo --table=messages ./test.go Message

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s: missing command\n", os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] != "generate-repo" {
		fmt.Fprintf(os.Stderr, "%s: unknown command: %s\n", os.Args[0], os.Args[1])
		os.Exit(1)
	}

	flag.CommandLine.Parse(os.Args[2:])

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

	if len(*out) > 0 {
		err := ioutil.WriteFile(*out, source, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: failed to write generated file: %s", os.Args[0], err)
		}
	} else {
		os.Stdout.Write(source)
	}
}
