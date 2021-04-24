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

// Package models contains the model definitions.
package models

//go:generate depot generate-repo --table=messages --repo-package repo --out ../repo/gen-messagerepo.go $GOFILE Message

import "time"

type (
	// Message demonstrates a persistent struct showing several mapped fields.
	Message struct {
		ID         string     `depot:"id,id"`
		Text       string     `depot:"text"`
		OrderIndex int        `depot:"order_index"`
		Length     float32    `depot:"len"`
		Attachment []byte     `depot:"attachment"`
		Created    time.Time  `depot:"created"`
		Updated    *time.Time `depot:"updated,nullable"`
	}
)
