package helpers

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

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
	} else if data, ok := resultado["Id"]; ok && data != nil {
		if err := MapToStruct(resultado, element); err == nil && element != nil {
			return nil
		} else {
			return Error("ResultToStruct", err, "502")
		}
	}
	return Error("ResultToStruct", resultado, "502")
}

// Add an element
func Add(element interface{}, api string, endpoint string, funcion string) (id int64, err map[string]interface{}) {
	defer ErrorControlFunction(funcion, "502")
	var resultado map[string]interface{}
	if error := request.SendJson(beego.AppConfig.String(api)+endpoint, "POST", &resultado, element); error == nil {
		if err := ResultToStruct(resultado, &element); err == nil && element != nil {
			r := reflect.ValueOf(element)
			f := reflect.Indirect(r).FieldByName("Id")
			num := f.Int()
			if num != 0 {
				return f.Int(), nil
			} else {
				return 0, Error(funcion, "Error obteniendo el campo Id del elemento", "502")
			}
		} else {
			return 0, Error(funcion, err, "502")
		}
	} else {
		return 0, Error(funcion, err, "502")
	}
}

// Get an element by id
func GetById(id int, element interface{}, api string, endpoint string, funcion string) (res interface{}, err map[string]interface{}) {
	defer ErrorControlFunction(funcion, "502")
	var resultado map[string]interface{}
	if error := request.GetJson(beego.AppConfig.String(api)+endpoint+"/"+strconv.Itoa(id), &resultado); error == nil {
		if err := ResultToStruct(resultado, &element); err == nil && element != nil {
			return element, nil
		} else {
			return nil, Error(funcion, err, "502")
		}
	} else {
		return nil, Error(funcion, err, "502")
	}
}
