package parse

import (
	"reflect"

	tfjson "github.com/hashicorp/terraform-json"
)

type Action string

const (
	ActionRecreate Action = "recreate"
)

func GroupByAction(resourceChangesMap map[string]*tfjson.ResourceChange) *map[interface{}][]*tfjson.ResourceChange {

	actionGroups := make(map[interface{}][]*tfjson.ResourceChange)

	for _, change := range resourceChangesMap {
		for _, action := range change.Change.Actions {
			// skip no-op and read actions
			if action != tfjson.ActionNoop && action != tfjson.ActionRead {
				actionGroups[action] = append(actionGroups[action], change)
			}
		}
	}

	return resolveRecreates(&actionGroups)
}

func resolveRecreates(actionGroups *map[interface{}][]*tfjson.ResourceChange) *map[interface{}][]*tfjson.ResourceChange {

	createList := []*tfjson.ResourceChange{}
	deleteList := []*tfjson.ResourceChange{}
	recreateList := []*tfjson.ResourceChange{}

	for _, resourceChange := range (*actionGroups)[tfjson.ActionCreate] {
		if resourceInBothCreateAndDelete(actionGroups, resourceChange) {
			recreateList = append(recreateList, resourceChange)
		} else {
			createList = append(createList, resourceChange)
		}
	}

	for _, resourceChange := range (*actionGroups)[tfjson.ActionDelete] {
		if resourceInBothCreateAndDelete(actionGroups, resourceChange) {
			recreateList = append(recreateList, resourceChange)
		} else {
			deleteList = append(deleteList, resourceChange)
		}
	}

	(*actionGroups)[tfjson.ActionDelete] = deleteList
	(*actionGroups)[tfjson.ActionCreate] = createList
	(*actionGroups)[ActionRecreate] = recreateList

	return actionGroups
}

func resourceInBothCreateAndDelete(actionGroups *map[interface{}][]*tfjson.ResourceChange, target *tfjson.ResourceChange) bool {
	actions := *actionGroups

	createdResources, hasCreatedResources := actions[tfjson.ActionCreate]
	deletedResources, hasDeletedResources := actions[tfjson.ActionDelete]

	if !hasCreatedResources || !hasDeletedResources {
		return false
	}

	contains := func(resources []*tfjson.ResourceChange, target *tfjson.ResourceChange) bool {
		for _, resource := range resources {
			if reflect.DeepEqual(resource, target) {
				return true
			}
		}
		return false
	}

	return contains(createdResources, target) && contains(deletedResources, target)
}
