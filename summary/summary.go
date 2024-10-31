package summary

import (
	"fmt"
	"strconv"
	"strings"

	"tf-plan-diff/config"
	"tf-plan-diff/parse"

	tfjson "github.com/hashicorp/terraform-json"
)

func Run(cfg *config.Config) error {

	actionGroups, err := parse.Parse(cfg)
	if err != nil {
		return fmt.Errorf("could not parse file at %s: %v", cfg.FilePath, err)
	}

	writeSummary(cfg, actionGroups)

	return nil
}

func writeSummary(cfg *config.Config, actionGroups *map[interface{}][]*tfjson.ResourceChange) {

	if len(*actionGroups) == 0 {
		if cfg.AppName != "" {
			fmt.Printf("%s Plan Contains No Pertinent Actions :green_circle:", cfg.AppName)
		} else {
			fmt.Printf("Plan Contains No Pertinent Actions :green_circle:")
		}
		return
	}

	writeAppName(cfg)

	fmt.Println("```diff")
	writeSummaryByAction(actionGroups)
	fmt.Println("```")

	writeFooter(actionGroups)
}

func writeAppName(cfg *config.Config) {
	if cfg.AppName != "" {
		fmt.Printf("# %s Plan Diff :build_construction:\n", cfg.AppName)
	} else {
		fmt.Printf("# Plan Diff :build_construction:\n", cfg.AppName)
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
	actions := []interface{}{tfjson.ActionNoop, tfjson.ActionRead, tfjson.ActionCreate, tfjson.ActionDelete, tfjson.ActionUpdate, parse.ActionRecreate}

	verbStrings := []string{}

	for _, action := range actions {

		actionVerb := strings.ToUpper(resolveVerb(action, false))
		numChanges := len((*actionGroups)[action])
		changeCount := strconv.Itoa(numChanges)

		if numChanges == 0 {
			continue
		}

		verbStrings = append(verbStrings, strings.Join([]string{actionVerb, changeCount}, " "))
	}

	verbString := strings.Join(verbStrings, ", ")

	footerString := strings.Join([]string{"This plan will:", verbString}, " ")

	fmt.Println(footerString)
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
