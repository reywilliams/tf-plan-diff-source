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
	action.Debugf("usings file path: %s", filePath)

	appName := action.GetInput(APP_NAME_INPUT_KEY)
	action.Debugf("usings app name: %s", appName)

	var includeReadActions bool
	includeReadInput := action.GetInput(INCLUDE_READ_INPUT_KEY)
	if includeReadInput != "" {
		includeReadActions = true
		action.Debugf("read actions enabled and will be included in plan diff")
	} else {
		includeReadActions = false
		action.Debugf("read actions disabled and will be not included in plan diff")
	}

	var includeNoOpActions bool
	includeNoOpInput := action.GetInput(INCLUDE_NOOP_INPUT_KEY)
	if includeNoOpInput != "" {
		includeReadActions = true
		action.Debugf("no-op actions enabled and will be included in plan diff")
	} else {
		includeNoOpActions = false
		action.Debugf("no-op actions disabled and will be not included in plan diff")
	}

	cfg := Config{FilePath: filePath,
		IncludeReadActions: includeReadActions,
		IncludeNoOpActions: includeNoOpActions,
		AppName:            appName}

	return &cfg, nil

}
