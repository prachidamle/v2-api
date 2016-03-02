package model

import (
	"github.com/rancher/go-rancher/client"
)	

type StackCommon struct {
	Common
	DockerCompose  string `json:"dockerCompose,omitempty" yaml:"docker_compose,omitempty" schema:",create=true,update=false"`
	ExternalId     string `json:"externalId,omitempty" yaml:"external_id,omitempty" db:"external_id" schema:",type=string,create=true,update=true,nullable=true"`
	RancherCompose string `json:"rancherCompose,omitempty" yaml:"rancher_compose,omitempty" schema:",type=string,create=true,update=false"`
	StartOnCreate  bool   `json:"startOnCreate,omitempty" yaml:"start_on_create,omitempty"`
}

type StackV1 struct {
	client.Resource
	StackCommon
	Environment        map[string]interface{} `json:"environment,omitempty" yaml:"environment,omitempty" schema:",type=map[string],create=true,update=false"`
	Outputs            map[string]interface{} `json:"outputs,omitempty" yaml:"outputs,omitempty"`
	PreviousExternalId string                 `json:"previousExternalId,omitempty" yaml:"previous_external_id,omitempty"`
}

type StackV2 struct {
	client.Resource
	StackCommon
	Variables  map[string]interface{} `json:"variables,omitempty" yaml:"variables,omitempty" schema:",type=map[string],create=true,update=false"`
	ServiceIds []ID                   `json:"serviceIds" schema:",type=array[reference[service]],create=false,update=false"`
}

type StackList struct {
	client.Collection
	Data []StackV2 `json:"data,omitempty"`
}

func getStackSchema(schemas *client.Schemas) {
	stack := AddType(schemas,"stack", StackV2{}, (*client.EnvironmentOperations)(nil))
	stack.CollectionMethods = []string{"GET", "POST"}
	stack.ResourceMethods = []string{"GET", "POST"}
	V2ToV1NamesMap["stack"] = "environment"
}
