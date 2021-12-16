package helpers

import (
	"encoding/json"
	"strconv"
	"github.com/udistrital/tesoreria_mid/models"

	//giros_crud "github.com/udistrital/giros_crud/models"
	terceros_crud "github.com/udistrital/terceros_crud/models"
	errorctrl "github.com/udistrital/utils_oas/errorctrl"
)

// ObtenerBancosConCodigos Obtiene los bancos que tienen codigos bancarios registrados en terceros
func ObtenerBancosConCodigos() (res *[]models.BancosPorCodigo, err map[string]interface{}) {
	funcion := "ObtenerBancosConCodigos"
	defer errorctrl.ErrorControlFunction(funcion, "500")
	var Bancos []terceros_crud.TerceroTipoTercero
	var BancoCod models.BancosPorCodigo
	var BancosCod []models.BancosPorCodigo
	query := make(map[string]string)
	query["TipoTerceroId__CodigoAbreviacion"] = "BANCO"
	if err := GetAll(&Bancos, "terceros_crud", "tercero_tipo_tercero", 1, query, nil, nil, nil, 0, 0); err == nil {
		for _, banco := range Bancos {
			var BancosCodTemp []terceros_crud.InfoComplementariaTercero
			var BancosNit []terceros_crud.DatosIdentificacion
			var aux models.DatosCodigos
			query1 := make(map[string]string)
			query1["InfoComplementariaId__CodigoAbreviacion"] = "COD_B"
			query1["TerceroId__Id"] = strconv.Itoa(banco.TerceroId.Id)
			if err1 := GetAll(&BancosCodTemp, "terceros_crud", "info_complementaria_tercero", 1, query1, nil, nil, nil, 0, 0); err1 == nil && BancosCodTemp[0].TerceroId != nil{
				query2 := make(map[string]string)
				query2["TerceroId__Id"] = strconv.Itoa(banco.TerceroId.Id)
				query2["TipoDocumentoId__CodigoAbreviacion"] = "NIT"
				if err2 := GetAll(&BancosNit, "terceros_crud", "datos_identificacion", 1, query2, nil, nil, nil, 0, 0); err2 == nil{
					json.Unmarshal([]byte(BancosCodTemp[0].Dato), &aux)
					BancoCod.IdTercero = BancosCodTemp[0].TerceroId.Id
					BancoCod.NIT = BancosNit[0].Numero
					BancoCod.CodigoSuper = aux.CodigoSuper
					BancoCod.CodigoAch = aux.CodigoAch
					BancoCod.NombreBanco = BancosCodTemp[0].TerceroId.NombreCompleto
					BancoCod.Activo = BancosCodTemp[0].TerceroId.Activo
					BancosCod = append(BancosCod, BancoCod)
				}
			}
		}
	}
	return &BancosCod, nil
}
