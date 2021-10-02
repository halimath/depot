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

// Package generate implements code generation for repository types.
package generate

import (
	"github.com/halimath/depot/internal/utils"
)

// Options defines the generation options.
type Options struct {
	// Name of the go file containing the entity definition.
	Filename string

	// Name of the type (as declared in the given file) to detect mappings for,
	EntityName string

	// Optional name of the table to map the entity to. Defaults to a SQL converted name of the entity.
	TableName string

	// Optional name of the go package containing the generated repo. Defaults to the entity's package name.
	RepoPackage string

	// Optional name of the repo type. Defaults to the entity name with a `Repo` suffix.
	RepoName string

	// Flag indicating if the repo should only provide finder methods and do not support modifications.
	ReadOnly bool
}

// GenerateRepository generates a repository implementation based
// on the given options. It returns the generated source code or
// an error.
func GenerateRepository(options Options) ([]byte, error) {
	if len(options.RepoName) == 0 {
		// If no repo name has been given use the default, which is
		// <EntityName>Repo
		options.RepoName = options.EntityName + "Repo"
	}

	if len(options.TableName) == 0 {
		// If no table has been specified we use a converted version
		// of the entity name.
		options.TableName = utils.SQLName(options.EntityName)
	}

	// Detect mapping by looking at the source code
	mapping, err := detectMapping(options.Filename, nil, options.EntityName)
	if err != nil {
		return nil, err
	}

	if len(options.RepoPackage) == 0 {
		// If no repo package has been specified we assume that it is
		// part of the entity package.
		options.RepoPackage = mapping.Package
	} else if options.RepoPackage != mapping.Package {
		// If repo and entity do not reside in the same package, we need
		// to prefix the entity's type name with its package to allow
		// correct imports.
		options.EntityName = mapping.Package + "." + options.EntityName
	}

	// Everything has been prepared. Generate the repo's source code.
	return generateRepo(mapping, &options)
}
