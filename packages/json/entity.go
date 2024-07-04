package json

type response struct {
	Meta meta        `json:"meta,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
type meta struct {
	StatusCode int         `json:"status_code,omitempty"`
	Code       string      `json:"code,omitempty"`
	Message    interface{} `json:"message,omitempty"`
}
