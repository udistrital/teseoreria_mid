package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/udistrital/tesoreria_mid/helpers"
	"github.com/udistrital/tesoreria_mid/models"
	//	"github.com/udistrital/utils_oas/time_bogota"
)

// SolicitudAvanceController operations for Solicitud_avance
type SolicitudAvanceController struct {
	beego.Controller
}

// URLMapping ...
func (c *SolicitudAvanceController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Solicitud_avance
// @Param	body		body 	models.Solicitud_avance	true		"body for Solicitud_avance content"
// @Success 201 {object} models.Solicitud_avance
// @Failure 403 body is empty
// @router / [post]
func (c *SolicitudAvanceController) Post() {
	defer helpers.ErrorControlController(c.Controller, "SolicitudAvanceController")
	var solicitudAvance models.SolicitudAvance
	// Decodificación de solicitud de avance
	if err1 := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitudAvance); err1 == nil {
		// Validación entradas
		if solicitudAvance.EstadoTipoSolicitud == nil {
			panic(helpers.Error("Post", "Error en los parámetros de ingreso", "400"))
		}
		// Lamada a helper
		if solicitudAvance, err3 := helpers.CrearSolicitudAvance(&solicitudAvance); err3 == nil && solicitudAvance != nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Registration successful", "Data": solicitudAvance}
		} else {
			panic(helpers.Error("Post", err3, err3["status"].(string)))
		}
	} else {
		panic(helpers.Error("Post", "Error en los parámetros de ingreso", "400"))		
	}
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Solicitud_avance by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Solicitud_avance
// @Failure 403 :id is empty
// @router /:id [get]
func (c *SolicitudAvanceController) GetOne() {
	defer helpers.ErrorControlController(c.Controller, "SolicitudAvanceController")
	// Validación entradas
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(helpers.Error("GetOne", "Error en los parámetros de entrada id no entero", "400"))
	}
	// Lamada a helper
	if solicitudAvance, err := helpers.ObtenerSolicitudAvancePorId(id); err == nil && solicitudAvance != nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": solicitudAvance}
	} else {
		panic(helpers.Error("GetOne", err, err["status"].(string)))
	}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Solicitud_avance
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Solicitud_avance
// @Failure 403
// @router / [get]
func (c *SolicitudAvanceController) GetAll() {
	defer helpers.ErrorControlController(c.Controller, "SolicitudAvanceController")
	// Parámetros
	var limit int = -1
	var offset int = -1
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	}
	if v, err := c.GetInt("offset"); err == nil {
		offset = v
	}
	// Llamada a helper
	if solicitudes, err := helpers.ObtenerSolicitudesAvance(limit, offset); err == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": solicitudes}
	} else {
		panic(helpers.Error("GetAll", err, err["status"].(string)))
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Solicitud_avance
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Solicitud_avance	true		"body for Solicitud_avance content"
// @Success 200 {object} models.Solicitud_avance
// @Failure 403 :id is not int
// @router /:id [put]
func (c *SolicitudAvanceController) Put() {
	defer helpers.ErrorControlController(c.Controller, "SolicitudAvanceController")
	funcion := "Put"
	// Validación entradas
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(helpers.Error(funcion, "Error en los parámetros de entrada id no entero", "400"))
	}
	// Decodificación de solicitud de avance
	solicitudAvance := models.SolicitudAvance{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitudAvance); err == nil {
		// Llamada a helper
		if err := helpers.ActualizarSolicitudAvance(&solicitudAvance); err == nil {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Update successful", "Data": solicitudAvance}
		} else {
			panic(helpers.Error(funcion, err, err["status"].(string)))
		}
	} else {
		panic(helpers.Error(funcion, err, "400"))
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Solicitud_avance
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SolicitudAvanceController) Delete() {
	defer helpers.ErrorControlController(c.Controller, "SolicitudAvanceController")
	panic(helpers.Error("Delete", "No implementado", "501"))
	c.ServeJSON()
}
