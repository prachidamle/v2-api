package server

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rancher/go-rancher/client"
	model "github.com/rancher/v2-api/model"
)

func (s *Server) StackByID(rw http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	return s.getStack(rw, r, vars["id"])
}

func (s *Server) StackList(rw http.ResponseWriter, r *http.Request) error {
	return s.getStack(rw, r, "")
}

func (s *Server) getStack(rw http.ResponseWriter, r *http.Request, id string) error {
	id = s.deobfuscate(r, "Stack", id)

	rows, err := s.namedQuery(s.getStacksSQL(r, id), map[string]interface{}{
		"account_id": s.getAccountID(r),
		"id":         id,
	})
	if err != nil {
		return err
	}
	defer rows.Close()

	response := &client.GenericCollection{
		Collection: client.Collection{
			Type:         "collection",
			ResourceType: "Stack",
		},
	}

	for rows.Next() {
		obj := &model.Stack{}
		obj.Type = "Stack"

		if err := rows.StructScan(obj); err != nil {
			return err
		}		

		// Obfuscate Ids
		obj.Id = s.obfuscate(r, "Stack", obj.Id)

		// Probably do something more with data
		response.Data = append(response.Data, obj)
	}

	return s.writeResponse(rows.Err(), r, response)
}

func (s *Server) getStacksSQL(r *http.Request, id string) string {
	q := `
		SELECT
			*  
		FROM environment
		WHERE
			account_id = :account_id
			AND removed IS NULL`

	if id != "" {
		q += " AND id = :id"
	}

	return q
}

func (s *Server) StackCreate(rw http.ResponseWriter, r *http.Request) error {
	rancherClient, err := s.getClient(r)
	if err != nil {
		return err
	}

	data := s.parseInputParameters(r)

	environ, err := rancherClient.Environment.Create(&client.Environment{
		Name: data.String("name"),
	})

	if err != nil {
		return err
	}

	return s.getStack(rw, r, environ.Id)
}
