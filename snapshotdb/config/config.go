package config

import (
	"flag"
	"fmt"
)

type Config struct {
	APIPort    int
	DBPath     string
	ParserPort int
	FreqSecs   int
}

func Parse() (*Config, error) {
	apiPort := flag.Int("api", 0, "Port number for REST API")
	dbPath := flag.String("db", "", "Path to SQLite database file")
	parserPort := flag.Int("parser", 0, "Port number of Parser service on localhost")
	freqSecs := flag.Int("freq", 0, "Fetch interval in seconds")

	flag.Parse()

	if *apiPort == 0 {
		return nil, fmt.Errorf("--api is required")
	}
	if *dbPath == "" {
		return nil, fmt.Errorf("--db is required")
	}
	if *parserPort == 0 {
		return nil, fmt.Errorf("--parser is required")
	}
	if *freqSecs == 0 {
		return nil, fmt.Errorf("--freq is required")
	}

	if *apiPort < 1 || *apiPort > 65535 {
		return nil, fmt.Errorf("--api must be between 1 and 65535")
	}
	if *parserPort < 1 || *parserPort > 65535 {
		return nil, fmt.Errorf("--parser must be between 1 and 65535")
	}
	if *freqSecs < 1 {
		return nil, fmt.Errorf("--freq must be at least 1 second")
	}

	return &Config{
		APIPort:    *apiPort,
		DBPath:     *dbPath,
		ParserPort: *parserPort,
		FreqSecs:   *freqSecs,
	}, nil
}

func (c *Config) ParserURL() string {
	return fmt.Sprintf("http://localhost:%d", c.ParserPort)
}

func (c *Config) ErrorLogPath() string {
	return c.DBPath + ".errors.jsonl"
}
