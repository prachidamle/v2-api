package model

import "github.com/rancher/go-rancher/client"

type StackCommon struct {
	Common
	DockerCompose  string `json:"dockerCompose,omitempty" yaml:"docker_compose,omitempty"`
	ExternalId     string `json:"externalId,omitempty" yaml:"external_id,omitempty" db:"external_id"`
	RancherCompose string `json:"rancherCompose,omitempty" yaml:"rancher_compose,omitempty"`
	StartOnCreate  bool   `json:"startOnCreate,omitempty" yaml:"start_on_create,omitempty"`
}

type StackV1 struct {
	client.Resource
	StackCommon
	Environment        map[string]interface{} `json:"environment,omitempty" yaml:"environment,omitempty"`
	Outputs            map[string]interface{} `json:"outputs,omitempty" yaml:"outputs,omitempty"`
	PreviousExternalId string                 `json:"previousExternalId,omitempty" yaml:"previous_external_id,omitempty"`
}

type StackV2 struct {
	client.Resource
	StackCommon
	Variables  map[string]interface{} `json:"variables,omitempty" yaml:"variables,omitempty"`
	ServiceIds []ID                   `json:"serviceIds"`
}

type StackList struct {
	client.Collection
	Data []StackV2 `json:"data,omitempty"`
}

func getStackSchema(schemas *client.Schemas) {
	stack := schemas.AddType("stack", StackV2{})
	stack.ResourceActions = map[string]client.Action{
		"create": client.Action{
			Input:  "",
			Output: "stack",
		},
	}
	stack.CollectionMethods = []string{"GET", "POST"}
}
