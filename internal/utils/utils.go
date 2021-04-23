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

// Package utils contains utility functions for handling database interaction
// and mapping.
package utils

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// SQLName converts the given go identifier name to a corresponding
// SQL identifier. The convention is that go identifiers follow the
// camel case convention while SQL identifiers follow the snake case
// convention.
func SQLName(ident string) string {
	if len(ident) == 0 {
		return ""
	}

	var result strings.Builder

	var current rune
	var next rune
	var l int
	var upperCount int

	for len(ident) > 0 {
		current = next
		next, l = utf8.DecodeRuneInString(ident)
		ident = ident[l:]

		if current == 0 {
			continue
		}

		isCurrentUpper := unicode.IsUpper(current)
		isNextUpper := unicode.IsUpper(next)

		if isCurrentUpper {
			upperCount++

			if upperCount > 1 && !isNextUpper {
				result.WriteByte('_')
			}
		} else {
			upperCount = 0
		}

		result.WriteRune(unicode.ToLower(current))

		if !isCurrentUpper && isNextUpper {
			result.WriteByte('_')
		}
	}

	result.WriteRune(unicode.ToLower(next))

	return result.String()
}
