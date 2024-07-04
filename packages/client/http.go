package httpClient

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ahsansandiah/dpo-test/packages/config"
	"github.com/ahsansandiah/dpo-test/packages/log"
)

type Http interface {
	Connect()
	CallURL(ctx context.Context, method, url string, header map[string]string, rawData []byte) ([]byte, error)
}

type Options struct {
	timeout int
	http    *http.Client
	log     log.Log
}

func NewHttp(cfg *config.Config, log log.Log) Http {
	opt := new(Options)
	opt.timeout = cfg.ServerHTTPReadTimeout
	opt.log = log
	return opt
}

func (o *Options) Connect() {
	httpClient := &http.Client{
		Timeout: time.Duration(o.timeout) * time.Second,
	}

	o.http = httpClient
}

func (o *Options) CallURL(ctx context.Context, method, url string, header map[string]string, rawData []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(rawData))
	if err != nil {
		err := errors.New("[CallURL-1] Failed To Prepare Request Client HTTP")
		o.log.ErrorLog(ctx, err)
		return nil, err
	}

	if len(header) > 0 {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}

	res, err := o.http.Do(req)
	if err != nil {
		err := errors.New("[CallURL-2] Failed To Request Client HTTP")
		o.log.ErrorLog(ctx, err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err := errors.New("[CallURL-3] Failed To Read Result Client HTTP")
		o.log.ErrorLog(ctx, err)
		return nil, err
	}

	o.log.HttpLog(ctx, req, rawData, body)

	if res.StatusCode != 200 {
		err := errors.New("[CallURL-4] Error Status Code Not 200")
		o.log.ErrorLog(ctx, err)
		return body, ErrCodeNot200
	}

	return body, nil
}
