package model

import "github.com/rancher/go-rancher/client"

type Stack struct {
	client.Resource
	Common

	DockerCompose  string                 `json:"dockerCompose,omitempty" yaml:"docker_compose,omitempty"`
	RancherCompose string                 `json:"rancherCompose,omitempty" yaml:"rancher_compose,omitempty"`
	Variables      map[string]interface{} `json:"variables,omitempty" yaml:"variables,omitempty"`
	StartOnCreate  bool                   `json:"startOnCreate,omitempty" yaml:"start_on_create,omitempty"`
	ExternalId     string                 `json:"externalId,omitempty" yaml:"external_id,omitempty"`
	ServiceIds     []ID                   `json:"serviceIds"`
}

type StackList struct {
	client.Collection
	Data []Stack `json:"data,omitempty"`
}
