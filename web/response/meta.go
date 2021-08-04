package response

import "strconv"

type Meta struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (m Meta) HttpStatus() int {
	if m.Code <= 200 {
		return 200
	}
	if m.Code/100 < 10 {
		return m.Code
	}
	if s := strconv.Itoa(m.Code); len(s) >= 3 {
		extractCode, err := strconv.Atoi(s[0:3])
		if err != nil {
			extractCode = 200
		}
		return extractCode
	}
	return 200
}
