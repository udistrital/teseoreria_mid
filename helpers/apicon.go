package helpers

import (
	"encoding/json"
	"strconv"
	"strings"
	_ "github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

// Cast map to struct
func MapToStruct(m interface{}, e interface{}) (err map[string]interface{}) {
	defer ErrorControlFunction("MapToStruct", "500")
	if jsonBody, err := json.Marshal(m); err == nil {
		if err := json.Unmarshal(jsonBody, e); err == nil {
			return nil
		} else {
			return Error("MapToStruct", err, "500")
		}
	} else {
		return Error("MapToStruct", err, "500")
	}
}

// Cast result data from api to struct element
func ResultToStruct(resultado map[string]interface{}, element interface{}) (err map[string]interface{}) {
	defer ErrorControlFunction("ResultToStruct", "500")
	if data, ok := resultado["Data"]; ok && data != nil {
		if err := MapToStruct(data, element); err == nil && element != nil {
			return nil
		} else {
			return err
		}
	}
	return Error("ResultToStruct", resultado, resultado["Status"].(string))
}

// Send an json element
func Send(element interface{}, url string, tipo string, v int, funcion string) (resultado map[string]interface{}, err map[string]interface{}) {
	defer ErrorControlFunction(funcion, "500")
	switch v {
	case 1:
		if error := request.SendJson(url, tipo, &element, element); error == nil {
			return nil, nil
		} else {
			return nil, Error(funcion, error, "400")
		}
	case 2:
		var resultado map[string]interface{}
		if error := request.SendJson(url, tipo, &resultado, element); error == nil {
			if err := ResultToStruct(resultado, &element); err == nil && element != nil {
				return resultado, nil
			} else {
				return nil, err
			}
		} else {
			return nil, Error(funcion, error, "400")
		}
	default:
		return nil, Error(funcion, "No se reconoce v", "400")
	}
}

// Add an element
func Add(element interface{}, api string, endpoint string, v int) (resultado map[string]interface{}, err map[string]interface{}) {
	return Send(element, beego.AppConfig.String(api)+endpoint, "POST", v, "Add")
}

// Update an element
func Update(id int, element interface{}, api string, endpoint string, v int) (resultado map[string]interface{}, err map[string]interface{}) {
	return Send(element, beego.AppConfig.String(api)+endpoint+"/"+strconv.Itoa(id), "PUT", v, "Update")
}

// Get one or more elements
func Get(element interface{}, url string, v int, funcion string) (err map[string]interface{}) {
	defer ErrorControlFunction(funcion, "500")
	switch v {
	case 1:
		if error := request.GetJson(url, &element); error == nil {
			return nil
		} else {
			return Error(funcion, err, "502")
		}
	case 2:
		var resultado map[string]interface{}
		if error := request.GetJson(url, &resultado); error == nil {
			if err := ResultToStruct(resultado, &element); err == nil && element != nil {
				return nil
			} else {
				return Error(funcion, err, "502")
			}
		} else {
			return Error(funcion, err, "502")
		}
	default:
		return Error(funcion, "No se reconoce v", "502")
	}
}

// Get an element by id
func GetById(id int, element interface{}, api string, endpoint string, v int) (err map[string]interface{}) {
	return Get(element, beego.AppConfig.String(api)+endpoint+"/"+strconv.Itoa(id), v, "GetById")
}

// Get all elements with params
func GetAllWithParams(elements interface{}, params map[string]string, api string, endpoint string, v int, funcion string) (err map[string]interface{}) {
	parametros := ""
	for name, value := range params {
		parametros += name + "=" + value + "&"
	}
	return Get(elements, beego.AppConfig.String(api)+endpoint+"?"+parametros, v, funcion)
}

// Get all
func GetAll(elements interface{}, api string, endpoint string, v int, query map[string]string, fields []string, sortby []string, order []string, limit int, offset int) (err map[string]interface{}) {
	funcion := "GetAll"
	defer ErrorControlFunction(funcion, "500")
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
	if limit >= 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if offset >= 0 {
		params["offset"] = strconv.Itoa(offset)
	}
	return GetAllWithParams(elements, params, api, endpoint, v, funcion)
}
