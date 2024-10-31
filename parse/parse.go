package parse

import (
	"fmt"
	"os"
	"tf-plan-diff/config"

	gw "github.com/gruntwork-io/terratest/modules/terraform"
	tfjson "github.com/hashicorp/terraform-json"
)

type Action string

const (
	ActionRecreate Action = "recreate"
)

func Parse(cfg *config.Config) (*map[interface{}][]*tfjson.ResourceChange, error) {

	planStruct, err := getPlanStruct(cfg.FilePath)
	if err != nil {
		return nil, fmt.Errorf("could not get plan struct from file path %s: %v", cfg.FilePath, err)
	}

	return groupByAction(cfg, planStruct.ResourceChangesMap), nil
}

func getPlanStruct(filePath string) (*gw.PlanStruct, error) {

	planJsonString, err := planJsonToString(filePath)
	if err != nil {
		return nil, fmt.Errorf("file %s could not read file: %v", filePath, err)
	}

	planStruct, err := gw.ParsePlanJSON(planJsonString)
	if err != nil {
		return nil, fmt.Errorf("could not parse plan json: %v", err)
	}

	return planStruct, nil
}

func planJsonToString(filePath string) (string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("file %s could not read file: %v", filePath, err)
	}

	return string(fileBytes), nil
}

func groupByAction(cfg *config.Config, resourceChangesMap map[string]*tfjson.ResourceChange) *map[interface{}][]*tfjson.ResourceChange {

	actionGroups := make(map[interface{}][]*tfjson.ResourceChange)

	for _, resource := range resourceChangesMap {
		action := resource.Change.Actions
		if action.Create() {
			actionGroups[tfjson.ActionCreate] = append(actionGroups[tfjson.ActionCreate], resource)
		} else if action.Delete() {
			actionGroups[tfjson.ActionDelete] = append(actionGroups[tfjson.ActionDelete], resource)
		} else if action.Update() {
			actionGroups[tfjson.ActionUpdate] = append(actionGroups[tfjson.ActionUpdate], resource)
		} else if action.Read() && cfg.IncludeReadActions {
			actionGroups[tfjson.ActionRead] = append(actionGroups[tfjson.ActionRead], resource)
		} else if action.NoOp() && cfg.IncludeNoOpActions {
			actionGroups[tfjson.ActionNoop] = append(actionGroups[tfjson.ActionNoop], resource)
		} else if action.Replace() {
			actionGroups[ActionRecreate] = append(actionGroups[ActionRecreate], resource)
		}
	}

	return &actionGroups
}
