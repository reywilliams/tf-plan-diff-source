package parse

import (
	"fmt"
	"os"

	gw "github.com/gruntwork-io/terratest/modules/terraform"
	tfjson "github.com/hashicorp/terraform-json"
)

func Parse(filePath string) (*map[interface{}][]*tfjson.ResourceChange, error) {

	planStruct, err := getPlanStruct(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not get plan struct from file path %s: %v", filePath, err)
	}

	return GroupByAction(planStruct.ResourceChangesMap), nil
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
