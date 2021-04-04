package models

//go:generate depot generate-repo --table=messages --repo-package repo --out ../repo/gen-messagerepo.go $GOFILE Message

import "time"

type (
	// Message demonstrates a persistent struct showing several mapped fields.
	Message struct {
		ID         string    `depot:"id,id"`
		Text       string    `depot:"text"`
		OrderIndex int       `depot:"order_index"`
		Length     float32   `depot:"len"`
		Attachment []byte    `depot:"attachment"`
		Created    time.Time `depot:"created"`
	}
)
