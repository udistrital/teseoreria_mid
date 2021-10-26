package helpers

import (
	"encoding/json"
	"strconv"

	"github.com/udistrital/tesoreria_mid/models"

	giros_crud "github.com/udistrital/giros_crud/models"
	terceros_crud "github.com/udistrital/terceros_crud/models"
	errorctrl "github.com/udistrital/utils_oas/errorctrl"
)

// Obtiene una cuenta bancaria por Id
func ObtenerCuentaBancariaBancoPorId(id int) (res *models.CuentaBancariaBanco, err map[string]interface{}) {
	funcion := "ObtenerCuentaBancariaBancoPorId"
	defer errorctrl.ErrorControlFunction(funcion, "500")
	cuentaBancaria := models.CuentaBancariaBanco{}
	// Consulta en giros_crud de cuentas bancarias
	if cuenta, err := ObtenerCuentaBancariaPorId(id, &cuentaBancaria); err == nil {
		// Consulta en terceros_crud
		_, err := ObtenerInfoComplementariaPorId(cuenta.SucursalId, &cuentaBancaria)
		if err != nil {
			return nil, err
		}
		return &cuentaBancaria, nil
	} else {
		return nil, err
	}
}

// Obtiene múltiples cuentas bancarias
func ObtenerCuentasBancariasBancos(limit int, offset int) (cuentasBancarias []models.CuentaBancariaBanco, err map[string]interface{}) {
	funcion := "ObtenerCuentasBancariasBancos"
	defer errorctrl.ErrorControlFunction(funcion, "500")
	var cuentas []giros_crud.CuentaBancaria
	if err := GetAll(&cuentas, "giros_crud", "cuenta_bancaria", 2, nil, nil, nil, nil, limit, offset); err == nil {
		cuentasBancarias := make([]models.CuentaBancariaBanco, len(cuentas))
		for i, cuenta := range cuentas {
			cuentaBancaria := models.CuentaBancariaBanco{}
			SetCuentaBancariaBancoPorCuentaBancaria(&cuenta, &cuentaBancaria)
			if _, err := ObtenerInfoComplementariaPorId(cuenta.SucursalId, &cuentaBancaria); err != nil {
				return nil, err
			}
			cuentasBancarias[i] = cuentaBancaria
		}
		return cuentasBancarias, nil
	} else {
		return nil, err
	}
}

// Obtención

func ObtenerCuentaBancariaPorId(id int, cuentaBancariaBanco *models.CuentaBancariaBanco) (c *giros_crud.CuentaBancaria, err map[string]interface{}) {
	funcion := "ObtenerCuentaBancariaPorId"
	defer errorctrl.ErrorControlFunction(funcion, "500")
	var cuentaBancaria *giros_crud.CuentaBancaria
	if err := GetById(id, &cuentaBancaria, "giros_crud", "cuenta_bancaria",
		2); err == nil && cuentaBancaria != nil {
		SetCuentaBancariaBancoPorCuentaBancaria(cuentaBancaria, cuentaBancariaBanco)
		return cuentaBancaria, nil
	} else {
		return nil, err
	}
}

func SetCuentaBancariaBancoPorCuentaBancaria(cuentaBancaria *giros_crud.CuentaBancaria, cuentaBancariaBanco *models.CuentaBancariaBanco) {
	funcion := "SetCuentaBancariaBancoPorCuentaBancaria"
	defer errorctrl.ErrorControlFunction(funcion, "500")
	if cuentaBancariaBanco != nil {
		cuentaBancariaBanco.Id = cuentaBancaria.Id
		cuentaBancariaBanco.NumeroCuenta = cuentaBancaria.NumeroCuenta
	}
}

func ObtenerInfoComplementariaPorId(id int, cuentaBancariaBanco *models.CuentaBancariaBanco) (res *terceros_crud.InfoComplementariaTercero, err map[string]interface{}) {
	funcion := "ObtenerInfoComplementariaPorId"
	defer errorctrl.ErrorControlFunction(funcion, "500")
	query := make(map[string]string)
	query["Id"] = strconv.Itoa(id)
	var infoComplementaria []terceros_crud.InfoComplementariaTercero
	if err := GetAll(&infoComplementaria, "terceros_crud", "info_complementaria_tercero", 1,
		query, nil, nil, nil, -1, -1); err == nil {
		if len(infoComplementaria) > 0 {
			SetCuentaBancariaPorInfoComplementaria(&infoComplementaria[0], cuentaBancariaBanco)
			return &infoComplementaria[0], nil
		} else {
			return nil, errorctrl.Error(funcion, "No existe sucursal asociada", "502")
		}
	} else {
		return nil, err
	}
}

func SetCuentaBancariaPorInfoComplementaria(infoComplementaria *terceros_crud.InfoComplementariaTercero, cuentaBancariaBanco *models.CuentaBancariaBanco) {
	funcion := "SetCuentaBancariaPorInfoComplementaria"
	defer errorctrl.ErrorControlFunction(funcion, "500")
	if cuentaBancariaBanco != nil {
		cuentaBancariaBanco.NombreBanco = infoComplementaria.TerceroId.NombreCompleto
		var dato map[string]interface{}
		if err := json.Unmarshal([]byte(infoComplementaria.Dato), &dato); err == nil {
			cuentaBancariaBanco.NombreSucursal = dato["nombreSucursal"].(string)
		}
	}
}
