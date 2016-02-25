package server

import (
	"github.com/gorilla/mux"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/v2-api/model"
	"net/http"
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
			ResourceType: "stack",
		},
	}

	for rows.Next() {
		obj := &model.Stack{}
		obj.Type = "stack"

		if err := rows.StructScan(obj); err != nil {
			return err
		}

		// Obfuscate Ids
		obj.Id = s.obfuscate(r, "Stack", obj.Id)

		if err = s.parseData(obj.Data, obj); err != nil {
			return err
		}

		response.Data = append(response.Data, obj)
	}

	return s.writeResponse(rows.Err(), r, response)
}

func (s *Server) getStacksSQL(r *http.Request, id string) string {

	commonSql := s.getSql(r, &model.Common{})

	q := `SELECT ` + commonSql + `
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
