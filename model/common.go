package model

type ID string

type Common struct {
	Name                 string `json:"name,omitempty" db:"name"`
	Description          string `json:"description,omitempty" db:"description"`
	State                string `json:"state,omitempty" db:"state"`
	UUID                 string `json:"uuid,omitempty" db:"uuid"`
	Created              string `json:"created" db:"created" db:"created"`
	Removed              string `json:"removed,omitempty" db:"removed"`
	ParentType           string `json:"parentType,omitempty"`
	Data                 string `db:"data"`
	Transitioning        string `json:"transitioning,omitempty"`
	TransitioningMessage string `json:"transitioningMessage,omitempty"`
}
