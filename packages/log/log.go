package log

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"

	traceHelper "github.com/ahsansandiah/dpo-test/helpers/trace"
	"github.com/ahsansandiah/dpo-test/packages/config"
	log "github.com/sirupsen/logrus"
)

type Log interface {
	ErrorLog(ctx context.Context, err error)
	CustomLog(r *http.Request, level string, data interface{})
	HttpLog(ctx context.Context, r *http.Request, payload []byte, response []byte)
}

type Options struct {
}

func NewLog() Log {
	opt := new(Options)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	return opt
}

func (o *Options) ErrorLog(ctx context.Context, err error) {
	file, funcx := traceHelper.ErrorTrace(3)
	checkID := ctx.Value(config.ContextKey("id"))
	id := ""
	if checkID != nil {
		id = checkID.(string)
	}

	log.WithFields(log.Fields{
		"id":   id,
		"file": file,
		"func": funcx,
	}).Error(err.Error())
}

func (o *Options) CustomLog(r *http.Request, level string, data interface{}) {
	ctx := r.Context()

	file, funcx := traceHelper.ErrorTrace(4)
	checkID := ctx.Value(config.ContextKey("id"))
	id := ""
	if checkID != nil {
		id = checkID.(string)
	}

	getParameters, _ := url.ParseQuery(r.URL.RawQuery)
	param := make(map[string]interface{}, len(getParameters))
	for k, v := range getParameters {
		param[k] = v[0]
	}
	if len(param) == 0 {
		param = nil
	}
	checkBody := ctx.Value(config.ContextKey("body"))
	body := make(map[string]interface{})
	if checkBody != nil {
		getBody := checkBody.([]byte)
		_ = json.Unmarshal(getBody, &body)
		if len(body) == 0 {
			body = nil
		}
	}

	headers := make(map[string]interface{})
	exceptionHeaders := map[string]bool{
		"Cache-Control":   true,
		"Postman-Token":   true,
		"Content-Length":  true,
		"Host":            true,
		"User-Agent":      true,
		"Accept-Encoding": true,
		"Accept":          true,
		"Connection":      true,
	}
	for k, v := range r.Header {
		if !exceptionHeaders[k] {
			headers[k] = v[0]
		}
	}
	if len(headers) == 0 {
		headers = nil
	}

	checkStartTime := ctx.Value(config.ContextKey("startTime"))
	responseTime := time.Duration(0)
	if checkStartTime != nil {
		getStartTime := ctx.Value(config.ContextKey("startTime")).(time.Time)
		responseTime = time.Duration(time.Since(getStartTime).Milliseconds())
	}

	if level == "ERROR" {
		data = data.(error).Error()
	}

	res := ResponseLog{
		HostName:      r.Host,
		Path:          r.URL.Path,
		RequestMethod: r.Method,
		Params:        param,
		Body:          body,
		Headers:       headers,
		UserAgent:     r.UserAgent(),
		ResponseTime:  responseTime,
		Response:      data,
	}

	logResponse := log.WithFields(log.Fields{
		"attributes": res,
		"id":         id,
		"file":       file,
		"func":       funcx,
	})
	if level == "ERROR" {
		logResponse.Error(data)
	} else {
		logResponse.Info("Success")
	}
}

func (o *Options) HttpLog(ctx context.Context, r *http.Request, payload []byte, response []byte) {

	file, funcx := traceHelper.ErrorTrace(4)
	checkID := ctx.Value(config.ContextKey("id"))
	id := ""
	if checkID != nil {
		id = checkID.(string)
	}

	getParameters, _ := url.ParseQuery(r.URL.RawQuery)
	param := make(map[string]interface{}, len(getParameters))
	for k, v := range getParameters {
		param[k] = v[0]
	}
	if len(param) == 0 {
		param = nil
	}

	body := make(map[string]interface{})
	if payload != nil {
		_ = json.Unmarshal(payload, &body)
		if len(body) == 0 {
			body = nil
		}
	}

	resCall := make(map[string]interface{})
	if response != nil {
		_ = json.Unmarshal(response, &resCall)
		if len(resCall) == 0 {
			resCall = nil
		}
	}

	headers := make(map[string]interface{})
	exceptionHeaders := map[string]bool{
		"Cache-Control":   true,
		"Postman-Token":   true,
		"Content-Length":  true,
		"Host":            true,
		"User-Agent":      true,
		"Accept-Encoding": true,
		"Accept":          true,
		"Connection":      true,
	}
	for k, v := range r.Header {
		if !exceptionHeaders[k] {
			headers[k] = v[0]
		}
	}
	if len(headers) == 0 {
		headers = nil
	}

	checkStartTime := ctx.Value(config.ContextKey("startTime"))
	responseTime := time.Duration(0)
	if checkStartTime != nil {
		getStartTime := ctx.Value(config.ContextKey("startTime")).(time.Time)
		responseTime = time.Duration(time.Since(getStartTime).Milliseconds())
	}

	res := ResponseLog{
		HostName:      r.Host,
		Path:          r.URL.Path,
		RequestMethod: r.Method,
		Params:        param,
		Body:          body,
		Headers:       headers,
		UserAgent:     r.UserAgent(),
		ResponseTime:  responseTime,
		Response:      resCall,
	}

	logResponse := log.WithFields(log.Fields{
		"attributes": res,
		"id":         id,
		"file":       file,
		"func":       funcx,
	})

	logResponse.Info("HTTP Call Trace")
}
