package render

import (
	"github.com/json-iterator/go"
	"net/http"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

// JSON contains the given interface object.
type JSON struct {
	Data interface{}
}

// Render (JSON) writes data with custom ContentType.
func (r JSON) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		return err
	}
	return
}

// WriteContentType (JSON) writes JSON ContentType.
func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// WriteJSON marshals the given interface object and writes it with custom ContentType.
func WriteJSON[T any](w http.ResponseWriter, obj T) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := jsoniter.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}
