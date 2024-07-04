package json

import (
	"encoding/json"
	"net/http"

	"github.com/ahsansandiah/dpo-test/packages/log"
)

type Json interface {
	SuccessResponse(w http.ResponseWriter, r *http.Request, statusCode int, message interface{}, data interface{})
	ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message interface{})
}

type Options struct {
	log log.Log
}

func NewJson(lg log.Log) Json {
	opt := new(Options)
	opt.log = lg

	return opt
}

// Return JSON
func (o *Options) writeJson(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	return enc.Encode(v)
}

// Return JSON Success
func (o *Options) SuccessResponse(w http.ResponseWriter, r *http.Request, statusCode int, message interface{}, data interface{}) {
	meta := meta{
		StatusCode: statusCode,
		Message:    message,
	}

	res := &response{
		Meta: meta,
		Data: data,
	}

	o.log.CustomLog(r, "SUCCESS", data)
	o.writeJson(w, statusCode, res)
}

// Return JSON Error
func (o *Options) ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message interface{}) {
	meta := meta{
		StatusCode: statusCode,
		Message:    message.(error).Error(),
	}

	res := &response{
		Meta: meta,
	}

	o.log.CustomLog(r, "ERROR", message)
	o.writeJson(w, statusCode, res)
}
