package main

import "encoding/json"

type Sample struct {
	ID   int    `json:"id"`                  // IDだよ
	Name string `json:"name" unknown:"none"` // 名前だよ
}

type Sample2 struct {
	internalField int             `json:"internal_field"`
	RawMessage    json.RawMessage `json:"raw_message"`
	Embedded      Sample          `json:"embedded"`
}
