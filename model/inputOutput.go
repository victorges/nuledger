package model

import "encoding/json"

type OperationInput struct {
	Account     json.RawMessage `json:"account"`
	Transaction json.RawMessage `json:"transaction"`
}

type StateOutput struct {
	Account    json.RawMessage `json:"account"`
	Violations json.RawMessage `json:"violations"`
}
