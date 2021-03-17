package helpers

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

// Get id from an element by reflect
func getId(element interface{}, funcion string) (id int64, err map[string]interface{}) {
	r := reflect.ValueOf(element)
	f := reflect.Indirect(r).FieldByName("Id")
	num := f.Int()
	if num != 0 {
		return f.Int(), nil
	} else {
		return 0, Error(funcion, "Error obteniendo el campo Id del elemento", "502")
	}
}

// Cast map[string]interface{} to struct
func MapToStruct(m map[string]interface{}, e interface{}) (err map[string]interface{}) {
	defer ErrorControlFunction("MapToStruct", "502")
	if jsonBody, err := json.Marshal(m); err == nil {
		if err := json.Unmarshal(jsonBody, e); err == nil {
			return nil
		} else {
			return Error("MapToStruct", err, "502")
		}
	} else {
		return Error("MapToStruct", err, "502")
	}
}

// Cast result data from api to struct element
func ResultToStruct(resultado map[string]interface{}, element interface{}) (err map[string]interface{}) {
	defer ErrorControlFunction("ResultToStruct", "502")
	if data, ok := resultado["Data"]; ok && data != nil {
		switch dat := data.(type) {
		case map[string]interface{}:
			if err := MapToStruct(dat, element); err == nil && element != nil {
				return nil
			} else {
				return Error("ResultToStruct", err, "502")
			}
		}
	}
	return Error("ResultToStruct", resultado, "502")
}

// Add an element
func Add(element interface{}, api string, endpoint string, v int, funcion string) (id int64, err map[string]interface{}) {
	defer ErrorControlFunction(funcion, "502")
	switch v {
	case 1:
		if error := request.SendJson(beego.AppConfig.String(api)+endpoint, "POST", &element, element); error == nil {
			return getId(element, funcion)
		} else {
			return 0, Error(funcion, err, "502")
		}
	case 2:
		var resultado map[string]interface{}
		if error := request.SendJson(beego.AppConfig.String(api)+endpoint, "POST", &resultado, element); error == nil {
			if err := ResultToStruct(resultado, &element); err == nil && element != nil {
				return getId(element, funcion)
			} else {
				return 0, Error(funcion, err, "502")
			}
		} else {
			return 0, Error(funcion, err, "502")
		}
	default:
		return 0, Error(funcion, "No se reconoce v", "502")
	}
}

// Get one or more elements
func Get(element interface{}, url string, v int, funcion string) (res interface{}, err map[string]interface{}) {
	switch v {
	case 1:
		if error := request.GetJson(url, &element); error == nil {
			return element, nil
		} else {
			return nil, Error(funcion, err, "502")
		}
	case 2:
		var resultado map[string]interface{}
		if error := request.GetJson(url, &resultado); error == nil {
			if err := ResultToStruct(resultado, &element); err == nil && element != nil {
				return element, nil
			} else {
				return nil, Error(funcion, err, "502")
			}
		} else {
			return nil, Error(funcion, err, "502")
		}
	default:
		return 0, Error(funcion, "No se reconoce v", "502")
	}
}

// Get an element by id
func GetById(id int, element interface{}, api string, endpoint string, v int, funcion string) (res interface{}, err map[string]interface{}) {
	defer ErrorControlFunction(funcion, "502")
	return Get(element, beego.AppConfig.String(api)+endpoint+"/"+strconv.Itoa(id), v, funcion)
}

// Get all elements with params
func GetAllWithParams(elements interface{}, params map[string]string, api string, endpoint string, v int, funcion string) (res interface{}, err map[string]interface{}) {
	defer ErrorControlFunction(funcion, "502")
	parametros := ""
	for name, value := range params {
		parametros += name + "=" + value + "&"
	}
	return Get(elements, beego.AppConfig.String(api)+endpoint+"?"+parametros, v, funcion)
}

// Get all
func GetAll(elements interface{}, api string, endpoint string, v int, funcion string, query map[string]string, fields []string, sortby []string, order []string, limit int, offset int) (res interface{}, err map[string]interface{}) {
	defer ErrorControlFunction(funcion, "502")
	params := make(map[string]string)
	queryString := ""
	for name, value := range query {
		queryString += name + ":" + value + ","
	}
	queryString = strings.TrimSuffix(queryString, ",")
	fieldsString := ""
	for _, field := range fields {
		fieldsString += field + ","
	}
	fieldsString = strings.TrimSuffix(fieldsString, ",")
	sortbyString := ""
	for _, sort := range sortby {
		sortbyString += sort + ","
	}
	sortbyString = strings.TrimSuffix(sortbyString, ",")
	orderString := ""
	for _, ord := range order {
		orderString += ord + ","
	}
	orderString = strings.TrimSuffix(orderString, ",")
	if queryString != "" {
		params["query"] = queryString
	}
	if fieldsString != "" {
		params["fields"] = fieldsString
	}
	if sortbyString != "" {
		params["sortby"] = sortbyString
	}
	if orderString != "" {
		params["order"] = orderString
	}
	if limit != 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if offset != 0 {
		params["offset"] = strconv.Itoa(offset)
	}
	return GetAllWithParams(elements, params, api, endpoint, v, funcion)
}
