package summary

import (
	"strings"

	"tf-plan-diff/config"
	"tf-plan-diff/parse"

	tfjson "github.com/hashicorp/terraform-json"
)

func Run(config *config.Config) error {
	return nil
}

func transLateToMarkDown(action interface{}, resourceChange *tfjson.ResourceChange) (string, error) {

	diffMarker := resolveDiffMarker(action)
	address := resourceChange.Address
	actionVerb := resolveVerb(action)

	return strings.Join([]string{diffMarker, address, "will be", actionVerb}, " "), nil
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
