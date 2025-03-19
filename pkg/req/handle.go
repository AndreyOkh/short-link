package req

import (
	"errors"
	"log"
	"net/http"
	"short-link/pkg/res"
)

const (
	typeAppJson = "application/json"
	typeAppForm = "application/x-www-form-urlencoded"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	var body T
	var err error
	conType := r.Header.Get("Content-Type")
	//log.Println("Content-Type: ", conType)

	switch conType {
	case typeAppJson:
		body, err = Decode[T](r.Body)
		if err != nil {
			res.Json(*w, err.Error(), http.StatusBadRequest)
			return nil, err
		}
	case typeAppForm:
		err = r.ParseForm()
		if err != nil {
			res.Json(*w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}
		body, err = DecodeForm[T](r.PostForm)
		if err != nil {
			res.Json(*w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}
	default:
		return &body, errors.New("unsupported content-type")
	}

	log.Printf("%#v", body)
	// Валидация структуры
	if err := IsValid(body); err != nil {
		res.Json(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &body, nil
}
