package config

import (
	"fmt"

	"github.com/sethvargo/go-githubactions"
)

const (
	FILE_PATH_INPUT_KEY    string = "file-path"
	INCLUDE_READ_INPUT_KEY string = "file-path"
)

type Config struct {
	filePath           string
	includeReadActions bool
}

func ConfigFromAction(action *githubactions.Action) (*Config, error) {

	filePath := action.GetInput(FILE_PATH_INPUT_KEY)
	if filePath == "" {
		return nil, fmt.Errorf("%s input was empty", FILE_PATH_INPUT_KEY)
	}

	includeReadActions := false
	includeReadInput := action.GetInput(INCLUDE_READ_INPUT_KEY)
	if includeReadInput != "" {
		includeReadActions = true
	}

	cfg := Config{filePath: filePath, includeReadActions: includeReadActions}

	return &cfg, nil

}
