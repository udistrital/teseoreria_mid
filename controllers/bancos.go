package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/tesoreria_mid/helpers"
	errorctrl "github.com/udistrital/utils_oas/errorctrl"
)

// BancoController operations for cuenta_bancaria_banco
type BancoController struct {
	beego.Controller
}

// URLMapping ...
func (c *BancoController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
}

// GetAll ...
// @Title GetAll
// @Description get Banco si tiene codigos registrados
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []models.BancosPorCodigo
// @Failure 400
// @router / [get]
func (c *BancoController) GetAll() {
	defer errorctrl.ErrorControlController(c.Controller, "BancoController")
	// Parámetros
	var limit int = 0
	var offset int = 0
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	}
	if v, err := c.GetInt("offset"); err == nil {
		offset = v
	}
	// Lamada a helper
	if banco, err := helpers.ObtenerBancosConCodigos(limit, offset); err == nil {
		if banco != nil {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Request successful", "Data": banco}
		} else {
			panic(errorctrl.Error("GetAll", "Banco no existe", "404"))
		}
	} else {
		panic(errorctrl.Error("GetOne", err, err["status"].(string)))
	}
	c.ServeJSON()
}
