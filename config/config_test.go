// +build !integration

package config

import (
	"github.com/elastic/beats/libbeat/common"
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	rawConfig := `
currencies: ["EUR", "USD", "GBP", "CHF"]
start_day: "2005-01-01"
`
	cfg, err := common.NewConfigFrom(rawConfig)
	if err != nil {
		t.Fatal(err)
	}
	parsed := Config{}
	if err := cfg.Unpack(&parsed); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(parsed, DefaultConfig) {
		t.Fatalf("Parsed config %s doesn't equal the default %s", parsed, DefaultConfig)
	}

}
