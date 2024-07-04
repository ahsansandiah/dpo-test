package log

import "time"

type ResponseLog struct {
	HostName      string                 `json:"host_name"`
	Path          string                 `json:"path"`
	RequestMethod string                 `json:"request_method"`
	Params        map[string]interface{} `json:"params"`
	Body          map[string]interface{} `json:"body"`
	Headers       map[string]interface{} `json:"headers"`
	UserAgent     string                 `json:"user_agent"`
	ResponseTime  time.Duration          `json:"response_time"`
	Response      interface{}            `json:"data"`
}
