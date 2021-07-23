package resource

import (
	"gitlab.id.vin/vincart/golib/web/render"
	"net/http"
)

func Write(w http.ResponseWriter, res Resource) {
	render.Render(w, res.Meta.HttpStatus(), render.JSON{Data: res})
}

func WriteError(w http.ResponseWriter, err error) {
	Write(w, Error(err))
}
