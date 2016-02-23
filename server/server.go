package server

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/rancher/go-rancher/api"

	"github.com/rancher/go-rancher/client"
)

type Server struct {
	driver, driverName string
}

func New(driver, driverName string) (*Server, error) {
	db, err := sqlx.Open(driver, driverName)
	s := &Server{
		driver:     driver,
		driverName: driverName,
	}
	db, err = s.getDb()
	if err != nil {
		db.Close()
	}
	return s, err
}

func (s *Server) getDb() (*sqlx.DB, error) {
	return sqlx.Open(s.driver, s.driverName)
}

func (s *Server) namedQuery(query string, args map[string]interface{}) (*sqlx.DB, *sqlx.Rows, error) {
	db, err := s.getDb()
	if err != nil {
		return nil, nil, err
	}

	rows, err := db.NamedQuery(query, args)
	return db, rows, err
}

func (s *Server) handleError(rw http.ResponseWriter, r *http.Request, err error) {
	apiError := client.ServerApiError{
		Type:    "error",
		Status:  500,
		Code:    "ServerError",
		Message: err.Error(),
	}
	data, err := json.Marshal(&apiError)
	if err == nil {
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(apiError.Status)
		rw.Write(data)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("Fail to marshall: %v", err)
	}
}

func (s *Server) HandlerFunc(schemas *client.Schemas, f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return api.ApiHandlerFunc(schemas, func(rw http.ResponseWriter, r *http.Request) {
		if err := f(rw, r); err != nil {
			s.handleError(rw, r, err)
		}
	})
}

func (s *Server) writeResponse(err error, r *http.Request, data interface{}) error {
	if err != nil {
		return err
	}
	api.GetApiContext(r).Write(data)
	return nil
}

func (s *Server) deobfuscate(r *http.Request, typeName string, id string) string {
	return id
}

func (s *Server) obfuscate(r *http.Request, typeName string, id string) string {
	if id == "" {
		return ""
	}
	return "1blah" + id
}

func (s *Server) getClient(r *http.Request) (*client.RancherClient, error) {
	return client.NewRancherClient(&client.ClientOpts{
		Url: "http://localhost:8080/v1/projects/1a5/schemas",
	})
}

func (s *Server) parseInputParameters(r *http.Request) InputData {
	data := InputData{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(data)
	return data
}
