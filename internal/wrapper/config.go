package wrapper

import (
	"encoding/json"
	"fmt"
)

// ParseConfig parses a configuration for a datasource from the given marshaled
// JSON object.
func ParseConfig[Config any](data []byte) (*Config, error) {
	config := new(Config)

	err := json.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse datasource JSON config: %w", err)
	}

	return config, nil
}
