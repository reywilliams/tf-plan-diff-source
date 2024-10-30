package summary

import (
	"fmt"
	"strings"

	"tf-plan-diff/config"
	"tf-plan-diff/parse"

	tfjson "github.com/hashicorp/terraform-json"
)

func Run(cfg *config.Config) error {

	actionGroups, err := parse.Parse(cfg.FilePath)
	if err != nil {
		return fmt.Errorf("could not parse file at %s: %v", cfg.FilePath, err)
	}

	writeSummary(cfg, actionGroups)

	return nil
}

func writeSummary(cfg *config.Config, actionGroups *map[interface{}][]*tfjson.ResourceChange) {

	writeAppName(cfg)
	writeSummaryByAction(actionGroups)
	writeFooter(actionGroups)
}

func writeAppName(cfg *config.Config) {
	if cfg.AppName != "" {
		fmt.Printf("# %s Plan Diff\n", cfg.AppName)
	}
}

func writeSummaryByAction(actionGroups *map[interface{}][]*tfjson.ResourceChange) {
	actions := []interface{}{tfjson.ActionNoop, tfjson.ActionRead, tfjson.ActionCreate, tfjson.ActionDelete, tfjson.ActionUpdate, parse.ActionRecreate}

	for _, action := range actions {
		for _, change := range (*actionGroups)[action] {
			fmt.Println(transLateToMarkDown(action, change))
		}
	}

}

func writeFooter(actionGroups *map[interface{}][]*tfjson.ResourceChange) {
	fmt.Println("footer")
}

func transLateToMarkDown(action interface{}, resourceChange *tfjson.ResourceChange) string {

	diffMarker := resolveDiffMarker(action)
	address := resourceChange.Address
	actionVerb := resolveVerb(action)

	return strings.Join([]string{diffMarker, address, "will be", actionVerb}, " ")
}

func resolveDiffMarker(action interface{}) string {
	switch action := action.(type) {
	case tfjson.Action:
		switch action {
		case tfjson.ActionDelete:
			return "-"
		case tfjson.ActionCreate:
			return "+"
		case tfjson.ActionUpdate:
			return "!"
		}
	case parse.Action:
		switch action {
		case parse.ActionRecreate:
			return "!"
		}
	default:
		return ""
	}

	return ""
}

func resolveVerb(action interface{}) string {
	switch action := action.(type) {
	case tfjson.Action:
		switch action {
		case tfjson.ActionDelete:
			return "deleted"
		case tfjson.ActionCreate:
			return "created"
		case tfjson.ActionUpdate:
			return "updated in place"
		}
	case parse.Action:
		switch action {
		case parse.ActionRecreate:
			return "recreated"
		}
	default:
		return ""
	}

	return ""
}
