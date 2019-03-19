// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"github.com/igor-kupczynski/gonbp"
	"time"
)

type Config struct {
	// Currencies to index
	Currencies []string `config:"currencies"`

	// Date in format 2006-01-02 when to start indexing
	StartDay string `config:"start_day"`
}

func (c *Config) startFrom() (time.Time, error) {
	return time.Parse(gonbp.DayFormat, c.StartDay)
}

// TimeSpan represents a [from, to) interval
type TimeSpan struct {
	// Begin of timespan, inclusive
	From time.Time

	// End of timespan, exclusive
	To time.Time
}

const deltaStr = "720h"

// Split [StartDay, Now] into TimeSpans, each no longer than 30 days
func (c *Config) SplitIntoTimeSpans() (ts []TimeSpan, err error) {
	delta, err := time.ParseDuration(deltaStr)
	if err != nil {
		return ts, err
	}

	from, err := c.startFrom()
	if err != nil {
		return ts, err
	}

	now := time.Now()

	for a, b := from, from.Add(delta); a.Before(now); a, b = b, b.Add(delta) {
		if b.After(now) {
			b = now
		}
		ts = append(ts, TimeSpan{From: a, To: b})
	}
	return ts, err
}

var DefaultConfig = Config{
	Currencies: []string{"EUR", "USD", "GBP", "CHF"},
	StartDay:   "2005-01-01",
}
