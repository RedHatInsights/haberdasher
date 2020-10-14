package processors

import (
	"net/http"
	"github.com/RedHatInsights/haberdasher/proxy"
)

type timerProcessor struct{}

func init() {
	var processor timerProcessor
	proxy.RegisterReverseProcessor("timer", processor)
}

func (p timerProcessor) BeforeRequest(req *http.Request) (error) {
	return nil
}

func (p timerProcessor) AfterRequest(req *http.Request) (error) {
	return nil
}
