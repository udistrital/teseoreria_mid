package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/tesoreria_mid/helpers"
	errorctrl "github.com/udistrital/utils_oas/errorctrl"
)

// CuentaBancariaBancoController operations for cuenta_bancaria_banco
type CuentaBancariaBancoController struct {
	beego.Controller
}

// URLMapping ...
func (c *CuentaBancariaBancoController) URLMapping() {
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
}

// GetOne ...
// @Title GetOne
// @Description get Cuenta_bancaria_banco by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.CuentaBancariaBanco
// @Failure 400 :id is empty
// @router /:id [get]
func (c *CuentaBancariaBancoController) GetOne() {
	defer errorctrl.ErrorControlController(c.Controller, "CuentaBancariaBancoController")
	// Validación entradas
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(errorctrl.Error("GetOne", "Error en los parámetros de entrada id no entero", "400"))
	}
	// Lamada a helper
	if cuentaBancaria, err := helpers.ObtenerCuentaBancariaBancoPorId(id); err == nil {
		if cuentaBancaria != nil {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": cuentaBancaria}
		} else {
			panic(errorctrl.Error("GetOne", "Cuenta no existe", "404"))
		}
	} else {
		panic(errorctrl.Error("GetOne", err, err["status"].(string)))
	}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Cuenta_bancaria_banco
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []models.CuentaBancariaBanco
// @Failure 400
// @router / [get]
func (c *CuentaBancariaBancoController) GetAll() {
	defer errorctrl.ErrorControlController(c.Controller, "CuentaBancariaBancoController")
	// Parámetros
	var limit int = -1
	var offset int = -1
	if v, err := c.GetInt("limit", -1); err == nil {
		limit = v
	} else {
		panic(helpers.Error("GetOne", "Error en los parámetros de entrada limit no entero", "400"))
	}
	if v, err := c.GetInt("offset", -1); err == nil {
		offset = v
	} else {
		panic(errorctrl.Error("GetOne", "Error en los parámetros de entrada offset no entero", "400"))
	}
	// Llamada a helper
	if cuentas, err := helpers.ObtenerCuentasBancariasBancos(limit, offset); err == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": cuentas}
	} else {
		panic(errorctrl.Error("GetAll", err, err["status"].(string)))
	}
	c.ServeJSON()
}
