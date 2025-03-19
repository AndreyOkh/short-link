package req

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"io"
	"log"
)

func Decode[T any](body io.ReadCloser) (T, error) {
	var payload T
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		log.Printf("body decode error: %v \n", err)
		return payload, err
	}
	return payload, nil
}

func DecodeForm[T any](form map[string][]string) (T, error) {
	var decoder = schema.NewDecoder()
	var payload T
	err := decoder.Decode(&payload, form)
	if err != nil {
		return payload, err
	}
	return payload, nil
}
