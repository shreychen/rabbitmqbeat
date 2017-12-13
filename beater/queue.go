package beater

import (
	"github.com/elastic/beats/libbeat/common"

	"github.com/elastic/beats/libbeat/logp"
)

const logSelector = "json"

func (bt *Rabbitmqbeat) Queues() ([]common.MapStr, error) {
	queues, err := bt.api.Queues()

	logp.Debug(logSelector, "filters:%s", bt.config.Filters)

	for _, f := range bt.config.Filters {
		logp.Info("Applying filters")

		queues, _ = f.Do(queues)
	}

	logp.Debug(logSelector, "final len: %v", len(queues))

	if err != nil {
		logp.Debug(logSelector, "get queues metrics failed ,caused by: %s", err.Error())
		return nil, err
	}

	return queues, nil
}
