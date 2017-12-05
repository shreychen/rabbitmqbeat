// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

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
