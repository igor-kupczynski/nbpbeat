package beater

import (
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/igor-kupczynski/gonbp"

	"github.com/igor-kupczynski/nbpbeat/config"
)

// Nbpbeat configuration.
type Nbpbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
	nbp    *gonbp.NbpClient
}

// New creates an instance of nbpbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	bt := &Nbpbeat{
		done:   make(chan struct{}),
		config: c,
		nbp:    gonbp.DefaultNbpClient,
	}
	return bt, nil
}

// Run starts nbpbeat.
func (bt *Nbpbeat) Run(b *beat.Beat) error {
	logp.Info("nbpbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	spans, err := bt.config.SplitIntoTimeSpans()
	if err != nil {
		return err
	}

	for _, ts := range spans {
		for _, curr := range bt.config.Currencies {
			rates, err := bt.nbp.DateRange("A", curr, ts.From, ts.To)
			if err != nil {
				return err
			}

			events := make([]beat.Event, len(rates.Rates))
			for i, rate := range rates.Rates {
				events[i] = beat.Event{
					Timestamp: rate.EffectiveDate,
					Fields: common.MapStr{
						"type":  b.Info.Name,
						"table": rate.Number,
						"mid":   rate.Mid,
						"curr":  rates.Code,
					},
				}
			}
			bt.client.PublishAll(events)
			logp.Info("Event sent")
		}
	}

	_ = <-bt.done
	return nil
}

// Stop stops nbpbeat.
func (bt *Nbpbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
