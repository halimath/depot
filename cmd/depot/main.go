// Copyright 2021 Alexander Metzner.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains the CLI for the generate command.
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
