package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/shreychen/rabbitmqbeat/config"
)

type Rabbitmqbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
	api    Api
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	api := NewRMQClient(config.URL, config.Username, config.Password)

	bt := &Rabbitmqbeat{
		done:   make(chan struct{}),
		config: config,
		api:    api,
	}
	return bt, nil
}

func (bt *Rabbitmqbeat) Run(b *beat.Beat) error {
	logp.Info("rabbitmqbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	//	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		//		event := beat.Event{
		//			Timestamp: time.Now(),
		//			Fields: common.MapStr{
		//				"type":    b.Info.Name,
		//				"counter": counter,
		//			},
		//		}
		//		bt.client.Publish(event)
		//		logp.Info("Event sent")
		//		counter++

		// send overview metrics
		if bt.config.Overview {
			if overview, err := bt.api.Overview(); err == nil {
				event := beat.Event{
					Timestamp: time.Now(),
					Fields: common.MapStr{
						"type":     b.Info.Name,
						"overview": overview,
					},
				}
				bt.client.Publish(event)
				logp.Info("Overview event sent")
			}
		}

		// send nodes metrics
		if bt.config.Nodes {
			if nodes, err := bt.api.Nodes(); err == nil {
				event := beat.Event{
					Timestamp: time.Now(),
					Fields: common.MapStr{
						"type":  b.Info.Name,
						"nodes": nodes,
					},
				}
				bt.client.Publish(event)
				logp.Info("Nodes event sent")
			}
		}

		// send queues metrics
		if bt.config.Queues {
			if queues, err := bt.api.Queues(); err == nil {
				event := beat.Event{
					Timestamp: time.Now(),
					Fields: common.MapStr{
						"type":   b.Info.Name,
						"queues": queues,
					},
				}
				bt.client.Publish(event)
				logp.Info("Queues event sent")
			}
		}
	}
}

func (bt *Rabbitmqbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
