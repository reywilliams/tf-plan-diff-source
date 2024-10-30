package config

import (
	"fmt"

	"github.com/sethvargo/go-githubactions"
)

const (
	FILE_PATH_INPUT_KEY    string = "file_path"
	APP_NAME_INPUT_KEY     string = "app_name"
	INCLUDE_READ_INPUT_KEY string = "include_read"
	INCLUDE_NOOP_INPUT_KEY string = "include_noop"
)

type Config struct {
	AppName            string
	FilePath           string
	IncludeReadActions bool
	IncludeNoOpActions bool
}

func ConfigFromAction(action *githubactions.Action) (*Config, error) {

	filePath := action.GetInput(FILE_PATH_INPUT_KEY)
	if filePath == "" {
		return nil, fmt.Errorf("%s input was empty", FILE_PATH_INPUT_KEY)
	}

	appName := action.GetInput(APP_NAME_INPUT_KEY)

	includeReadActions := false
	includeReadInput := action.GetInput(INCLUDE_READ_INPUT_KEY)
	if includeReadInput != "" {
		includeReadActions = true
	}

	includeNoOpActions := false
	includeNoOpInput := action.GetInput(INCLUDE_NOOP_INPUT_KEY)
	if includeNoOpInput != "" {
		includeNoOpActions = true
	}

	cfg := Config{FilePath: filePath,
		IncludeReadActions: includeReadActions,
		IncludeNoOpActions: includeNoOpActions,
		AppName:            appName}

	return &cfg, nil

}
