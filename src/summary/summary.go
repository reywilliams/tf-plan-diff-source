package summary

import (
	"fmt"
	"strconv"
	"strings"

	"tf-plan-diff/config"
	"tf-plan-diff/parse"

	tfjson "github.com/hashicorp/terraform-json"
	ga "github.com/sethvargo/go-githubactions"
)

func Run(cfg *config.Config, action *ga.Action) error {

	actionGroups, err := parse.Parse(cfg)
	if err != nil {
		return fmt.Errorf("could not parse file at %s: %v", cfg.FilePath, err)
	}

	writeSummary(cfg, actionGroups, action)

	return nil
}

func writeSummary(cfg *config.Config, actionGroups *map[interface{}][]*tfjson.ResourceChange, action *ga.Action) {

	if len(*actionGroups) == 0 {
		if cfg.AppName != "" {
			action.AddStepSummary(fmt.Sprintf("# `%s` Plan Contains No Pertinent Actions :green_circle:", cfg.AppName))
		} else {
			action.AddStepSummary("# Plan Contains No Pertinent Actions :green_circle:")
		}
		return
	}

	writeAppName(cfg, action)

	action.AddStepSummary("```diff")
	writeSummaryByAction(actionGroups, action)
	action.AddStepSummary("```")

	writeFooter(actionGroups, action)
}

func writeAppName(cfg *config.Config, action *ga.Action) {
	if cfg.AppName != "" {
		action.AddStepSummary(fmt.Sprintf("# `%s` Plan Diff :building_construction:", cfg.AppName))
	} else {
		action.AddStepSummary("# Plan Diff :building_construction:")
	}
}

func writeSummaryByAction(actionGroups *map[interface{}][]*tfjson.ResourceChange, ghAction *ga.Action) {
	actions := []interface{}{tfjson.ActionNoop, tfjson.ActionRead, tfjson.ActionCreate, tfjson.ActionDelete, tfjson.ActionUpdate, parse.ActionRecreate}

	for _, action := range actions {
		for _, change := range (*actionGroups)[action] {
			ghAction.AddStepSummary(transLateToMarkDown(action, change))
		}
	}

}

func writeFooter(actionGroups *map[interface{}][]*tfjson.ResourceChange, action *ga.Action) {
	actions := []interface{}{tfjson.ActionNoop, tfjson.ActionRead, tfjson.ActionCreate, tfjson.ActionDelete, tfjson.ActionUpdate, parse.ActionRecreate}

	verbStrings := []string{}

	for _, action := range actions {

		actionVerb := strings.ToUpper(resolveVerb(action, false))
		actionVerb = strings.Join([]string{"*", actionVerb, "*"}, "") // make verb bold
		numChanges := len((*actionGroups)[action])
		changeCount := strconv.Itoa(numChanges)

		if numChanges == 0 {
			continue
		}

		verbStrings = append(verbStrings, strings.Join([]string{actionVerb, changeCount}, " "))
	}

	verbString := strings.Join(verbStrings, ", ")

	footerString := strings.Join([]string{":warning:", "This plan will:", verbString, ":warning:"}, " ")

	action.AddStepSummary(footerString)
}

func transLateToMarkDown(action interface{}, resourceChange *tfjson.ResourceChange) string {

	diffMarker := resolveDiffMarker(action)
	address := resourceChange.Address
	actionVerb := resolveVerb(action, true)

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
		case tfjson.ActionNoop:
			return "#"
		case tfjson.ActionRead:
			return "#"
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

func resolveVerb(action interface{}, pastTense bool) string {
	switch action := action.(type) {
	case tfjson.Action:
		switch action {
		case tfjson.ActionDelete:
			if pastTense {
				return "deleted"
			} else {
				return "delete"
			}
		case tfjson.ActionCreate:
			if pastTense {
				return "created"
			} else {
				return "create"
			}
		case tfjson.ActionUpdate:
			if pastTense {
				return "updated"
			} else {
				return "update"
			}
		}
	case parse.Action:
		switch action {
		case parse.ActionRecreate:
			if pastTense {
				return "recreated"
			} else {
				return "recreate"
			}
		}
	default:
		return ""
	}

	return ""
}
