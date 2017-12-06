// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	//	"fmt"
	"regexp"
	"time"

	"github.com/elastic/beats/libbeat/logp"

	"github.com/elastic/beats/libbeat/common"
)

const logSelector = "json"

type Filter struct {
	Type string `json:"type"`
	Exp  string `json:"exp"`
}

type Config struct {
	Period   time.Duration `config:"period"`
	URL      string        `config:"url"`
	Username string        `config:"username"`
	Password string        `config:"Password"`
	Overview bool          `config:"overview"`
	Nodes    bool          `config:"nodes"`
	Queues   bool          `config:"queues"`
	Filters  []Filter      `config:"filters"`
}

var DefaultConfig = Config{
	Period:   10 * time.Second,
	URL:      "http://localhost:15672",
	Username: "guest",
	Password: "guest",
	Overview: true,
	Nodes:    true,
	Queues:   true,
}

func RemoveOne(cm []common.MapStr, i int) []common.MapStr {
	cm[len(cm)-1], cm[i] = cm[i], cm[len(cm)-1]
	return cm[:len(cm)-1]
}

func (f *Filter) Do(cm []common.MapStr) []common.MapStr {
	reg, err := regexp.Compile(f.Exp)
	if err != nil {
		logp.Debug(logSelector, "regexp compile failed: %s", err.Error())
		return cm
	}

	for i, c := range cm {
		qname, err := c.GetValue("name")

		if err != nil {
			continue
		}
		//		qname = qname.(string)
		//		qname = fmt.Sprintf("%s", qname)
		if reg.Match([]byte(qname.(string))) {
			cm = RemoveOne(cm, i)
			logp.Debug("json", "removed metrics of queue: %s", qname)
		}
	}
	return cm
}
