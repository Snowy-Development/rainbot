package rainbot

import (
    "testing"
)

func TestReadConfig(t *testing.T) {
	config, err := ReadConfig("testConfig.json")

	if err != nil || config.Host != "irc.freenode.org" {
		t.Fatalf("Couldn't read config file: ",  err)
	}
}
