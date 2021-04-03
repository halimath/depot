package models

import "time"

type (
	Message struct {
		ID         string    `depot:"id,id"`
		Text       string    `depot:"text"`
		OrderIndex int       `depot:"order_index"`
		Length     float32   `depot:"len"`
		Attachment []byte    `depot:"attachment"`
		Created    time.Time `depot:"created"`
	}
)
