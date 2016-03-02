package model

import (
	//"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	schemaTagName = "schema"
)

var (
	V2ToV1NamesMap map[string]string
	ResourceActionsToStates map[string][]string
)

type StateTransitionObject struct {
	FromState string `json:"fromState"`
	ToState string `json:"toState"`
	Field string `json:"field"`
	Type string `json:"type"`
}

func AddType(s *client.Schemas, schemaName string, obj interface{}, operations interface{}) *client.Schema {
	t := reflect.TypeOf(obj)
	opType := reflect.TypeOf(operations).Elem()
	schema := client.Schema{
		Resource: client.Resource{
			Id:    schemaName,
			Type:  "schema",
			Links: map[string]string{},
		},
		PluralName:        guessPluralName(schemaName),
		ResourceFields:    typeToFields(t),
		CollectionMethods: []string{"GET"},
		ResourceMethods:   []string{"GET"},
		ResourceActions:   typeToActions(opType, schemaName),
	}

	if s.Data == nil {
		s.Data = []client.Schema{}
	}

	s.Data = append(s.Data, schema)

	return &s.Data[len(s.Data)-1]
}

func guessPluralName(name string) string {
	if name == "" {
		return ""
	}

	if strings.HasSuffix(name, "s") ||
		strings.HasSuffix(name, "ch") ||
		strings.HasSuffix(name, "x") {
		return name + "es"
	}
	return name + "s"
}

func parseTagString(tag string) string {
	splitted := strings.Split(tag, "=")
	if len(splitted) == 2 {
		return splitted[1]
	}
	return ""
}

func parseTagBool(tag string) bool {
	splitted := strings.Split(tag, "=")
	if len(splitted) == 2 {
		if val, err := strconv.ParseBool(splitted[1]); err == nil {
			return val
		}
	}
	return false
}

func parseTagInterface(tag string) interface{} {
	splitted := strings.Split(tag, "=")
	if len(splitted) == 2 {
		return splitted[1]
	}
	return nil
}

func typeToFields(t reflect.Type) map[string]client.Field {
	result := map[string]client.Field{}

	for i := 0; i < t.NumField(); i++ {
		schemaField := client.Field{}

		typeField := t.Field(i)
		if typeField.Anonymous && typeField.Type.Kind() == reflect.Struct {
			parentFields := typeToFields(typeField.Type)
			for k, v := range result {
				parentFields[k] = v
			}
			result = parentFields
			continue
		} else if typeField.Anonymous {
			continue
		}

		create := false
		update := false
		nullable := false
		required := false
		schemaType := ""
		var defaultValue interface{}
		tagParts := strings.Split(typeField.Tag.Get(schemaTagName), ",")
		for _, tag := range tagParts[1:] {
			switch {
			case strings.HasPrefix(tag, "create"):
				create = parseTagBool(tag)
			case strings.HasPrefix(tag, "update"):
				update = parseTagBool(tag)
			case strings.HasPrefix(tag, "nullable"):
				nullable = parseTagBool(tag)
			case strings.HasPrefix(tag, "required"):
				required = parseTagBool(tag)
			case strings.HasPrefix(tag, "type"):
				schemaType = parseTagString(tag)
			case strings.HasPrefix(tag, "default"):
				defaultValue = parseTagInterface(tag)
			}
		}
		schemaField.Create = create
		schemaField.Update = update
		schemaField.Nullable = nullable
		schemaField.Required = required
		schemaField.Default = defaultValue

		if schemaType != "" {
			schemaField.Type = schemaType
		} else {
			fieldString := strings.ToLower(typeField.Type.Kind().String())

			switch {
			case strings.HasPrefix(fieldString, "int") || strings.HasPrefix(fieldString, "uint"):
				schemaField.Type = "int"
			case fieldString == "bool":
				schemaField.Type = "boolean"
			case fieldString == "float32" || fieldString == "float64":
				schemaField.Type = "float"
			case fieldString == "string":
				schemaField.Type = "string"
			case fieldString == "map":
				schemaField.Type = "map[string]"
			case fieldString == "slice":
				schemaField.Type = "array[string]"
			}
		}

		name := strings.Split(typeField.Tag.Get("json"), ",")[0]
		if name == "" && len(typeField.Name) > 1 {
			name = strings.ToLower(typeField.Name[0:1]) + typeField.Name[1:]
		} else if name == "" {
			name = typeField.Name
		}

		if schemaField.Type != "" {
			result[name] = schemaField
		}
	}

	return result
}

func GetTransitioning(state string, trans string) string {
	if trans == "error" {
		return trans
	}
	if strings.HasSuffix(state, "ing") && strings.ToLower(state) != "running" {
		return "yes"
	}
	return "no"
}

func ToLowerCamelCase(input string) string {
	return (strings.ToLower(input[:1]) + input[1:])
}

func typeToActions(val reflect.Type, v2ResourceName string) map[string]client.Action {
	resourceActions := make(map[string]client.Action)
	

	for i := 0; i < val.NumMethod(); i++ {
		var actionName, input, output string
		
		if !strings.HasPrefix(val.Method(i).Name, "Action") {
			continue
		}

		actionName = strings.TrimPrefix(val.Method(i).Name, "Action")
		actionName = strings.ToLower(actionName)
		
		if(strings.EqualFold(actionName, "list") || strings.EqualFold(actionName, "create") || strings.EqualFold(actionName, "byid")) {
			continue
		}		

		fmt.Printf("reflecting action %s ", val.Method(i).Name)
		
		if val.Method(i).Type.NumIn() > 1 {
			if val.Method(i).Type.In(1).Kind().String() != "interface" {
				input = strings.TrimPrefix(val.Method(i).Type.In(1).Elem().Name(), "*client.")
				input = ToLowerCamelCase(input)
			}
		}

		if val.Method(i).Type.NumOut() > 0 {
			if val.Method(i).Type.Out(0).Kind().String() != "interface" {
				output = strings.TrimPrefix(val.Method(i).Type.Out(0).Elem().Name(), "*")
				output = ToLowerCamelCase(output)
			}
		}

		resourceActions[actionName] = client.Action{Input: input, Output: v2ResourceName}
	}

	return resourceActions
}
