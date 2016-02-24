package model

import "github.com/go-sql-driver/mysql"

type ID string

type Common struct {
	Name                 string         `json:"name"`
	Description          string         `json:"description"`
	State                string         `json:"state"`
	UUID                 string         `json:"uuid"`
	Created              mysql.NullTime `json:"created" db:"created"`
	Removed              mysql.NullTime `json:"removed" db:"removed"`
	ParentType           string         `json:"parentType"`
	Data                 string         `json:"_"`
	Transitioning        string         `json:"transitioning"`
	TransitioningMessage string         `json:"transitioningMessage"`
}
