package response

import (
	"gitlab.com/golibs-starter/golib/web/render"
	"net/http"
)

func Write[T any](w http.ResponseWriter, res Response[T]) {
	render.Render(w, res.Meta.HttpStatus(), render.JSON{Data: res})
}

func WriteError(w http.ResponseWriter, err error) {
	Write(w, Error(err))
}
