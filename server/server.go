package server

import (
	"encoding/json"
	"net/http"
	"strings"
	//"io/ioutil"
	"fmt"
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/rancher/go-rancher/api"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/v2-api/model"
)

type Server struct {
	DB                 *sqlx.DB
	driver, driverName string
	Schemas		*client.Schemas
}

type SchemaConvertor interface {
	FromSchema(obj interface{}) (interface{}, error)
	ToSchema(obj interface{}) (interface{}, error)
}

func New(driver, driverName string) (*Server, error) {
	db, err := sqlx.Open(driver, driverName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	v2schemas := model.NewSchema()
		
	server := &Server{
		driver:     driver,
		driverName: driverName,
		DB:         db,
		Schemas:  v2schemas,
	}
	server.buildActionsToStatesMap()
	return server, err
}

func (s *Server) namedQuery(query string, args map[string]interface{}) (*sqlx.Rows, error) {
	rows, err := s.DB.Unsafe().NamedQuery(query, args)
	return rows, err
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
	return strings.TrimPrefix(id, getObfuscator(typeName))
}

func getObfuscator(typeName string) string {
	obfuscator := "1"
	return obfuscator + typeName[0:1]
}

func (s *Server) obfuscate(r *http.Request, typeName string, id string) string {
	if id == "" {
		return ""
	}
	return getObfuscator(typeName) + id
}

func (s *Server) getClient(r *http.Request) (*client.RancherClient, error) {
	return client.NewRancherClient(&client.ClientOpts{
		Url: "http://localhost:8080/v1/projects/1a5/schemas",
	})
}

func (s *Server) getBaseClient(r *http.Request) (*client.RancherClient, error) {
	return client.NewRancherClient(&client.ClientOpts{
		Url: "http://localhost:8080/v1",
	})
}

func (s *Server) parseInputParameters(r *http.Request, obj interface{}) error {

	decoder := json.NewDecoder(r.Body)


		if err := decoder.Decode(&obj); err == io.EOF {
			
		} else if err != nil {
			fmt.Printf("\n errror: %v",err)
			return err
		}


	return nil
}

func (s *Server) isValidAction(action string, resource *client.Resource) bool{
	schema := s.Schemas.Schema(resource.Type)
	_, ok := schema.ResourceActions[action]
	return ok
}


func (s *Server) buildActionsToStatesMap() error {
	for _, schema := range s.Schemas.Data {
		for name, _ := range schema.ResourceActions {
			err := s.parseResourceActions(schema.Resource.Id, name)
			if err != nil {
				fmt.Printf("\nerror: %v",err)
				return err
			}
		}
	}
	
	return nil
}

func (s *Server) parseResourceActions(resourceName, actionName string) error {
	rancherClient, err := s.getBaseClient(nil)
	if err != nil {
		return err
	}
	
	processDefId := "1pd!" + model.V2ToV1NamesMap[resourceName] + "." + actionName
	
	pdefinitions, err := rancherClient.ProcessDefinition.ById(processDefId)
	if err != nil {
		return err
	}
	
	if pdefinitions == nil {
		return nil
	}
	
	var states []string
	
	for _, transJson := range pdefinitions.StateTransitions {
		transObj := model.StateTransitionObject{}
		
		bytes, err := json.Marshal(transJson)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(bytes, &transObj); err != nil {
			return err
		}
		states = append(states, transObj.FromState)
	}
	
	model.ResourceActionsToStates[resourceName + "." + actionName] = states

	return nil
}


func (s *Server) addResourceActions(r *http.Request, resource *client.Resource, resourceState string) {
	apiContext := api.GetApiContext(r)
	resource.Actions = make(map[string]string)	
	
	resourceSchema, ok := s.Schemas.CheckSchema(resource.Type)
	if ok {
		for actionName, _ := range resourceSchema.ResourceActions {
			//if this action is based on states, check the resource state for eligibility
			addAction := false
			validStates, ok := model.ResourceActionsToStates[resource.Type + "." + actionName]
			if ok {
					for i := range validStates {
						if validStates[i] == resourceState {
							addAction = true
						}
					}
			} else {
				addAction = true
			}
			if addAction {
				actionLink :=  resource.Id + "/?action=" + actionName
				resource.Actions[actionName] = apiContext.UrlBuilder.ReferenceByIdLink(resource.Type, actionLink)
			}
		}
	}
}