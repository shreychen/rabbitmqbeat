package beater

import (
	//	"regexp"

	"github.com/elastic/beats/libbeat/common"

	"github.com/elastic/beats/libbeat/logp"
)

const logSelector = "json"

func (bt *Rabbitmqbeat) Queues() ([]common.MapStr, error) {
	_, err := bt.api.Queues()

	logp.Debug("json", "filters:%s", bt.config.Filters)

	if err != nil {
		logp.Debug(logSelector, "get queues metrics failed ,caused by: %s", err.Error())
		return nil, err
	}

	return nil, nil
}
