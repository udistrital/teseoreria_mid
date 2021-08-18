package helpers

import (
	"encoding/json"
	"strconv"
	"fmt"
	mapstructure "github.com/mitchellh/mapstructure"
	"github.com/udistrital/tesoreria_mid/models"

	avances_crud "github.com/udistrital/avances_crud/models"
	solicitudes_crud "github.com/udistrital/solicitudes_crud/models"
)

// Crea una solicitud de avance
func CrearSolicitudAvance(solicitudAvance *models.SolicitudAvance) (sol *models.SolicitudAvance, err map[string]interface{}) {
	funcion := "CrearSolicitudAvance"
	defer ErrorControlFunction(funcion, "502")
	// Guardado en solicitudes_crud
	solicitudAvance.Id = 0
	solicitud := solicitudes_crud.Solicitud{}
	if res1, err1 := EnviarSolicitudCrud(&solicitud, solicitudAvance); err1 == nil && res1 != nil {
		// Guardar el solicitante
		if res2, err2 := EnviarSolicitante(&solicitudes_crud.Solicitante{},
			solicitudAvance, &solicitud); err2 != nil && res2 != nil {
			return nil, Error(funcion, err2, "502")
		}
		// Guardar objetivo y justificacion (observaciones)
		if resObj, errObj := EnviarObservacion(&solicitudes_crud.Observacion{},
			solicitudAvance, &solicitud, solicitudAvance.Objetivo,
			5, "Objetivo"); errObj != nil && resObj != nil {
			return nil, Error(funcion, errObj, "502")
		}
		if resJust, errJust := EnviarObservacion(&solicitudes_crud.Observacion{},
			solicitudAvance, &solicitud, solicitudAvance.Justificacion,
			6, "Justificación"); errJust != nil && resJust != nil {
			return nil, Error(funcion, errJust, "502")
		}
		// Guardado en avances_crud
		solicitudAvanceCrud := avances_crud.SolicitudAvance{
			SolicitudId: solicitud.Id, Activo: true}
		if res3, err3 := EnviarSolicitudAvanceCrud(&solicitudAvanceCrud); err3 == nil && res3 != nil {
			solicitudAvance.Id = solicitudAvanceCrud.Id
			solicitudAvance.FechaCreacion = solicitudAvanceCrud.FechaCreacion
			solicitudAvance.FechaModificacion = solicitudAvanceCrud.FechaCreacion
			if res4, err4 := EnviarSolicitudYEspecificacionTipoAvanceCrud(solicitudAvance, res3); err4 == nil && res4 != nil {
				return solicitudAvance, nil
			}
			return solicitudAvance, nil
		} else {
			return nil, Error(funcion, err3, "502")
		}
	} else {
		return nil, Error(funcion, err1, "502")
	}
}

// Actualiza una solicitud de avance
func ActualizarSolicitudAvance(solicitudAvance *models.SolicitudAvance) (err map[string]interface{}) {
	funcion := "ActualizarSolicitudAvance"
	defer ErrorControlFunction(funcion, "502")
	// Consulta en avances_crud
	if solAvanceCrud, err := ObtenerSolicitudAvanceCrudPorId(solicitudAvance.Id, nil); err == nil {
		solicitudAvance.Id = solAvanceCrud.Id
		solicitudAvance.SolicitudId = solAvanceCrud.SolicitudId
		solicitudAvance.FechaCreacion = solAvanceCrud.FechaCreacion
		// Actualizacion en solicitudes_crud
		solicitud := solicitudes_crud.Solicitud{
			Id: solicitudAvance.SolicitudId, FechaCreacion: solicitudAvance.FechaCreacion}
		if res, err := EnviarSolicitudCrud(&solicitud, solicitudAvance); err != nil  && res != nil{
			return Error(funcion, err, "502")
		}
		// Guardar el solicitante
		if solicitante, err := ObtenerSolicitantePorSolicitudId(solicitudAvance.SolicitudId, nil); err == nil {
			var res map[string]interface{}
			if res, err = EnviarSolicitante(&solicitudes_crud.Solicitante{
				Id: solicitante.Id, FechaCreacion: solicitante.FechaCreacion},
				solicitudAvance, &solicitud); err != nil && res != nil {
				return Error(funcion, err, "502")
			}
		} else {
			return Error(funcion, err, "502")
		}
		// Guardar objetivo
		if obj, err := ObtenerObjetivoPorSolicitudId(solicitudAvance.SolicitudId, nil); err == nil {
			if resObj, errObj := EnviarObservacion(&solicitudes_crud.Observacion{
				Id: obj.Id, FechaCreacion: obj.FechaCreacion},
				solicitudAvance, &solicitud, solicitudAvance.Objetivo, 5, "Objetivo"); errObj != nil  && resObj != nil{
				return Error(funcion, errObj, "502")
			}
		} else {
			return Error(funcion, err, "502")
		}
		// Guardar justificación
		if just, err := ObtenerJustificacionPorSolicitudId(solicitudAvance.SolicitudId, nil); err == nil {
			if resJust, errJust := EnviarObservacion(&solicitudes_crud.Observacion{
				Id: just.Id, FechaCreacion: just.FechaCreacion},
				solicitudAvance, &solicitud, solicitudAvance.Justificacion, 6, "Justificación"); errJust != nil && resJust != nil {
				return Error(funcion, errJust, "502")
			}
		} else {
			return Error(funcion, err, "502")
		}
		return nil
	} else {
		return Error(funcion, err, "502")
	}
}

// Obtiene una solicitud de avance buscando por Id
func ObtenerSolicitudAvancePorId(id int) (sol *models.SolicitudAvance, err map[string]interface{}) {
	funcion := "ObtenerSolicitudAvancePorId"
	defer ErrorControlFunction(funcion, "502")
	solicitudAvance := models.SolicitudAvance{}
	// Consulta en avances_crud
	if _, err := ObtenerSolicitudAvanceCrudPorId(id, &solicitudAvance); err == nil {
		// Consulta en solicitudes_crud
		_, err := ObtenerSolicitudPorId(solicitudAvance.SolicitudId, &solicitudAvance)
		if err != nil {
			return nil, Error(funcion, err, "502")
		}
		_, err = ObtenerSolicitantePorSolicitudId(solicitudAvance.SolicitudId, &solicitudAvance)
		if err != nil {
			return nil, Error(funcion, err, "502")
		}
		_, err = ObtenerObjetivoPorSolicitudId(solicitudAvance.SolicitudId, &solicitudAvance)
		if err != nil {
			return nil, Error(funcion, err, "502")
		}
		_, err = ObtenerJustificacionPorSolicitudId(solicitudAvance.SolicitudId, &solicitudAvance)
		if err != nil {
			return nil, Error(funcion, err, "502")
		}
		return &solicitudAvance, nil
	} else {
		return nil, Error(funcion, err, "502")
	}
}

// Obtiene múltiples solicitudes de avance
func ObtenerSolicitudesAvance(limit int, offset int) (solicitudes []models.SolicitudAvance, err map[string]interface{}) {
	funcion := "ObtenerSolicitudesAvance"
	defer ErrorControlFunction(funcion, "502")
	var solsAvance []avances_crud.SolicitudAvance
	if err := GetAll(&solsAvance, "avances_crud", "solicitud_avance", 2, nil, nil, nil, nil, limit, offset); err == nil {
		solicitudes := make([]models.SolicitudAvance, len(solsAvance))
		query := make(map[string]string)
		query["TipoSolicitud"] = "5"
		var estadosTipo []solicitudes_crud.EstadoTipoSolicitud
		if err := GetAll(&estadosTipo, "solicitudes_crud", "estado_tipo_solicitud", 2, query, nil, nil, nil, 0, -1); err != nil {
			return nil, Error(funcion, err, "502")
		}
		for i, solAvance := range solsAvance {
			soli := models.SolicitudAvance{}
			SetSolicitudAvancePorSolicitudAvanceCrud(&solAvance, &soli)
			if _, err := ObtenerSolicitudPorId(solAvance.SolicitudId, &soli); err != nil {
				return nil, Error(funcion, err, "502")
			}
			for _, estadoTipo := range estadosTipo {
				if soli.EstadoTipoSolicitud.Id == estadoTipo.Id {
					soli.EstadoTipoSolicitud = &estadoTipo
				}
			}
			solicitudes[i] = soli
		}
		return solicitudes, nil
	} else {
		return nil, Error(funcion, err, "502")
	}
}

// Envío

func EnviarSolicitudCrud(solicitud *solicitudes_crud.Solicitud, solicitudAvance *models.SolicitudAvance) (res map[string]interface{}, err map[string]interface{}) {
	funcion := "EnviarSolicitudCrud"
	defer ErrorControlFunction(funcion, "502")
	referencia, errJSON := json.Marshal(map[string]interface{}{
		"VigenciaId":            solicitudAvance.VigenciaId,
		"AreaFuncionalId":       solicitudAvance.AreaFuncionalId,
		"CargoOrdenadorGastoId": solicitudAvance.CargoOrdenadorGastoId,
		"DependenciaId":         solicitudAvance.DependenciaId,
		"FacultadId":            solicitudAvance.FacultadId,
		"ProyectoCurricularId":  solicitudAvance.ProyectoCurricularId,
		"ConvenioId":            solicitudAvance.ConvenioId,
		"ProyectoInvId":         solicitudAvance.ProyectoInvId,
		"FechaEvento":           solicitudAvance.FechaEvento,
	})
	if errJSON == nil {
		solicitud.EstadoTipoSolicitudId = solicitudAvance.EstadoTipoSolicitud
		solicitud.Referencia = string(referencia)
		solicitud.FechaRadicacion = solicitudAvance.FechaRadicacion
		solicitud.SolicitudFinalizada = solicitudAvance.SolicitudFinalizada
		solicitud.Resultado = solicitudAvance.Resultado
		solicitud.Activo = solicitudAvance.Activo
		var err1 map[string]interface{}
		if solicitud.Id == 0 {
			res, err1 = Add(solicitud, "solicitudes_crud", "solicitud", 2)
		} else {
			res, err1 = Update(solicitud.Id, solicitud, "solicitudes_crud", "solicitud", 2)
		}
		if err1 == nil {
			solicitudAvance.SolicitudId = solicitud.Id
			return res, nil
		} else {
			return nil, Error(funcion, err1, "502")
		}
	}
	return nil, Error(funcion, errJSON, "400")
}

func EnviarSolicitante(solicitante *solicitudes_crud.Solicitante, solicitudAvance *models.SolicitudAvance, solicitud *solicitudes_crud.Solicitud) (res map[string]interface{}, err map[string]interface{}) {
	funcion := "CrearSolicitante"
	defer ErrorControlFunction(funcion, "502")
	solicitante.TerceroId = solicitudAvance.TerceroId
	solicitante.SolicitudId = solicitud
	solicitante.Activo = true
	var err1 map[string]interface{}
	if solicitante.Id == 0 {
		res, err1 = Add(solicitante, "solicitudes_crud", "solicitante", 2)
	} else {
		res, err1 = Update(solicitante.Id, solicitante, "solicitudes_crud", "solicitante", 2)
	}
	if err1 == nil {
		return res, nil
	} else {
		return nil, Error(funcion, err1, "502")
	}
}

func EnviarObservacion(observacion *solicitudes_crud.Observacion, solicitudAvance *models.SolicitudAvance, solicitud *solicitudes_crud.Solicitud, valor string, tipoObservacionId int, titulo string) (res map[string]interface{}, err map[string]interface{}) {
	funcion := "EnviarObservacion"
	defer ErrorControlFunction(funcion, "502")
	tipoObservacion := solicitudes_crud.TipoObservacion{Id: tipoObservacionId}
	observacion.TipoObservacionId = &tipoObservacion
	observacion.SolicitudId = solicitud
	observacion.TerceroId = solicitudAvance.TerceroId
	observacion.Titulo = titulo
	observacion.Valor = valor
	observacion.Activo = true
	var err1 map[string]interface{}
	if observacion.Id == 0 {
		res, err1 = Add(observacion, "solicitudes_crud", "observacion", 2)
	} else {
		res, err1 = Update(observacion.Id, observacion, "solicitudes_crud", "observacion", 2)
	}
	if err1 == nil {
		return res, nil
	} else {
		return nil, Error(funcion, err1, "502")
	}
}

func EnviarSolicitudAvanceCrud(solicitud *avances_crud.SolicitudAvance) (res map[string]interface{}, err map[string]interface{}) {
	funcion := "EnviarSolicitudAvanceCrud"
	defer ErrorControlFunction(funcion, "502")
	var err1 map[string]interface{}
	if solicitud.Id == 0 {
		res, err1 = Add(solicitud, "avances_crud", "solicitud_avance", 2)
	} else {
		res, err1 = Update(solicitud.Id, solicitud, "avances_crud", "solicitud_avance", 2)
	}
	if err1 == nil {
		return res, nil
	} else {
		return nil, Error(funcion, err1, "502")
	}
}

func EnviarSolicitudYEspecificacionTipoAvanceCrud(solicitudAvance *models.SolicitudAvance, res3 map[string]interface{}) (res map[string]interface{}, err map[string]interface{}) {
	funcion := "EnviarSolicitudYEspecificacionTipoAvanceCrud"
	defer ErrorControlFunction(funcion, "502")
	//fmt.Println("Solicitud ", solicitudAvance)
	//fmt.Println("res3 ", res3["Data"])
	//fmt.Println("tipoAvance ", solicitudAvance.AvanceTotal[0].Id)
	//fmt.Println("tipoAvanceLen ", len(solicitudAvance.AvanceTotal))
	//resAux := map[string]interface{}
	//resAux := res3["Data"]
	solAvance := avances_crud.SolicitudAvance{}
	errDecod := mapstructure.Decode(res3["Data"], &solAvance)
	if errDecod != nil{
		panic(errDecod)
	}
	solTipoAvance := avances_crud.SolicitudTipoAvance{}
	//solEspecificacionesAvance := avances_crud.EspecificacionTipoAvance{}
	solTipoAvance.SolicitudAvanceId = &solAvance
	
	//solAvance.SolicitudAvanceId = solicitudAvance.Id
	for i := 0; i < len(solicitudAvance.AvanceTotal); i++ {

	//------------  SOLICITUD TIPO AVANCE
		solTipoAvance.TipoAvanceId = solicitudAvance.AvanceTotal[i]
		//REVISAR DESCRIPCION
		solTipoAvance.Descripcion = ""
		jsonEspecificaciones, _ := json.Marshal(solicitudAvance.TipoAvance[i]["especificaciones"])
		var especificaciones []map[string]interface{}
		errJson := json.Unmarshal([]byte(jsonEspecificaciones), &especificaciones)
		if errJson != nil{
			panic(errJson)
		}
		var valorTotal float64
		for j := 0; j < len(especificaciones); j++ {
			valorTotal += (especificaciones[j]["valor"]).(float64)
		}
		solTipoAvance.Valor = valorTotal
		solTipoAvance.FechaCreacion = solAvance.FechaCreacion
		solTipoAvance.FechaModificacion = solAvance.FechaModificacion
		solTipoAvance.Activo = true

	// ----------- POST
		var err1 map[string]interface{}
		if solTipoAvance.Id == 0 {
			res, err1 = Add(solTipoAvance, "avances_crud", "solicitud_tipo_avance", 2)
			if err1 == nil{
				if res5, err5 := EnviarSolicitudRequisitoTipoAvanceCrud(solicitudAvance, res, solTipoAvance.TipoAvanceId, i); err5 == nil && res5 != nil {
					return nil, nil
				}else{
					return nil, err5
				}
			} else{
				return nil, err1
			}
		} else {
			// Falta update
			fmt.Println("FALLO por id")
			//res, err1 = Update(solicitud.Id, solicitud, "solicitudes_crud", "solicitud", 2)
		}

		/*if err1 == nil {
			//solicitudAvance.SolicitudId = solicitud.Id
			return res, err1
		} else {
			return nil, Error(funcion, err1, "502")
		}*/
		fmt.Println(err1)

	}
	return nil, nil
}

func EnviarSolicitudRequisitoTipoAvanceCrud(solicitudAvance *models.SolicitudAvance, res map[string]interface{}, solTipoAvanceId *avances_crud.TipoAvance, i int)(res5 map[string]interface{}, err map[string]interface{}) {
	funcion := "EnviarSolicitudRequisitoTipoAvanceCrud"
	defer ErrorControlFunction(funcion, "502")
	solRequisitoTipoAvance := avances_crud.SolicitudRequisitoTipoAvance{}
	solRequisito := avances_crud.SolicitudTipoAvance{}
	errDecod := mapstructure.Decode(res["Data"], &solRequisito)
	if errDecod != nil {
		panic(errDecod)
	}
	solRequisitoTipoAvance.SolicitudTipoAvanceId = &solRequisito
	jsonRequisitos, _ := json.Marshal(solicitudAvance.TipoAvance[i]["requisitos"])
	var requisitos []map[string]interface{}
	errJson := json.Unmarshal([]byte(jsonRequisitos), &requisitos)
	if errJson != nil {
		panic(errJson)
	}
	for j := 0; j < len(requisitos); j++ {
		query := make(map[string]string)
		query["TipoAvanceId"] = strconv.Itoa(solTipoAvanceId.Id)
		query["RequisitoAvanceId"] = strconv.FormatFloat((requisitos[j]["Id"]).(float64), 'f', -1, 64)

		var requisitoTipoAvance []avances_crud.RequisitoTipoAvance
		if err := GetAll(&requisitoTipoAvance, "avances_crud", "requisito_tipo_avance", 2,
			query, nil, nil, nil, -1, -1); err == nil {
			if len(requisitoTipoAvance) > 0 {
				sol := requisitoTipoAvance[0]
				solRequisitoTipoAvance.RequisitoTipoAvanceId = &sol
				//return &sol, nil
			} else {
			return nil, Error(funcion, "No existe solicitante asociado", "502")
			}
		} else {
			return nil, Error(funcion, err, "502")
		}
		solRequisitoTipoAvance.Observaciones = ""
		solRequisitoTipoAvance.Documento = fmt.Sprintf("%g", requisitos[j]["idDocumento"].(float64))
		solRequisitoTipoAvance.Activo = true
	

		var err1 map[string]interface{}
		if solRequisitoTipoAvance.Id == 0 {
			res, err1 = Add(solRequisitoTipoAvance, "avances_crud", "solicitud_requisito_tipo_avance", 2)
			fmt.Println("res ", res)
		} else {
			// Falta update
			fmt.Println("FALLO por id")
			//res, err1 = Update(solicitud.Id, solicitud, "solicitudes_crud", "solicitud", 2)
		}
		fmt.Println(err1)
	}
	return nil, nil
}

// Obtención

func ObtenerSolicitudAvanceCrudPorId(id int, solicitudAvance *models.SolicitudAvance) (sol *avances_crud.SolicitudAvance, err map[string]interface{}) {
	funcion := "ObtenerSolicitudAvanceCrudPorId"
	defer ErrorControlFunction(funcion, "502")
	var solicitudAvanceCrud *avances_crud.SolicitudAvance
	if err := GetById(id, &solicitudAvanceCrud, "avances_crud", "solicitud_avance",
		2); err == nil && solicitudAvanceCrud != nil {
		SetSolicitudAvancePorSolicitudAvanceCrud(solicitudAvanceCrud, solicitudAvance)
		return solicitudAvanceCrud, nil
	} else {
		return nil, Error(funcion, err, "502")
	}
}

func SetSolicitudAvancePorSolicitudAvanceCrud(solicitudAvanceCrud *avances_crud.SolicitudAvance, solicitudAvance *models.SolicitudAvance) {
	funcion := "SetSolicitudAvancePorSolicitudAvanceCrud"
	defer ErrorControlFunction(funcion, "502")
	if solicitudAvance != nil {
		solicitudAvance.Id = solicitudAvanceCrud.Id
		solicitudAvance.SolicitudId = solicitudAvanceCrud.SolicitudId
		solicitudAvance.FechaCreacion = solicitudAvanceCrud.FechaCreacion
		solicitudAvance.FechaModificacion = solicitudAvanceCrud.FechaModificacion
		solicitudAvance.Activo = solicitudAvanceCrud.Activo
	}
}

func ObtenerSolicitudPorId(id int, solicitudAvance *models.SolicitudAvance) (sol *solicitudes_crud.Solicitud, err map[string]interface{}) {
	funcion := "ObtenerSolicitudPorId"
	defer ErrorControlFunction(funcion, "502")
	var solicitudAvanceCrud *solicitudes_crud.Solicitud
	if err := GetById(id, &solicitudAvanceCrud, "solicitudes_crud", "solicitud",
		1); err == nil && solicitudAvanceCrud != nil {
		SetSolicitudAvancePorSolicitudCrud(solicitudAvanceCrud, solicitudAvance)
		return solicitudAvanceCrud, nil
	} else {
		return nil, Error(funcion, err, "502")
	}
}

func SetSolicitudAvancePorSolicitudCrud(solicitudAvanceCrud *solicitudes_crud.Solicitud, solicitudAvance *models.SolicitudAvance) {
	funcion := "SetSolicitudAvancePorSolicitudCrud"
	defer ErrorControlFunction(funcion, "502")
	if solicitudAvance != nil {
		solicitudAvance.EstadoTipoSolicitud = solicitudAvanceCrud.EstadoTipoSolicitudId
		solicitudAvance.FechaRadicacion = solicitudAvanceCrud.FechaRadicacion
		solicitudAvance.SolicitudFinalizada = solicitudAvanceCrud.SolicitudFinalizada
		solicitudAvance.Resultado = solicitudAvanceCrud.Resultado
		var referencia map[string]interface{}
		if err := json.Unmarshal([]byte(solicitudAvanceCrud.Referencia), &referencia); err == nil {
			solicitudAvance.VigenciaId = int(referencia["VigenciaId"].(float64))
			solicitudAvance.AreaFuncionalId = int(referencia["AreaFuncionalId"].(float64))
			solicitudAvance.CargoOrdenadorGastoId = int(referencia["CargoOrdenadorGastoId"].(float64))
			solicitudAvance.DependenciaId = int(referencia["DependenciaId"].(float64))
			solicitudAvance.FacultadId = int(referencia["FacultadId"].(float64))
			solicitudAvance.ProyectoCurricularId = int(referencia["ProyectoCurricularId"].(float64))
			solicitudAvance.ConvenioId = int(referencia["ConvenioId"].(float64))
			solicitudAvance.ProyectoInvId = int(referencia["ProyectoInvId"].(float64))
			solicitudAvance.FechaEvento = referencia["FechaEvento"].(string)
		}
	}
}

func ObtenerSolicitantePorSolicitudId(id int, solicitudAvance *models.SolicitudAvance) (solicitante *solicitudes_crud.Solicitante, err map[string]interface{}) {
	funcion := "ObtenerSolicitantePorSolicitudId"
	defer ErrorControlFunction(funcion, "502")
	query := make(map[string]string)
	query["SolicitudId"] = strconv.Itoa(id)
	var solicitantes []solicitudes_crud.Solicitante
	if err := GetAll(&solicitantes, "solicitudes_crud", "solicitante", 1,
		query, nil, nil, nil, -1, -1); err == nil {
		if len(solicitantes) > 0 {
			sol := solicitantes[0]
			SetSolicitudAvancePorSolicitante(&sol, solicitudAvance)
			return &sol, nil
		} else {
			return nil, Error(funcion, "No existe solicitante asociado", "502")
		}
	} else {
		return nil, Error(funcion, err, "502")
	}
}

func SetSolicitudAvancePorSolicitante(sol *solicitudes_crud.Solicitante, solicitudAvance *models.SolicitudAvance) {
	funcion := "SetSolicitudAvancePorSolicitante"
	defer ErrorControlFunction(funcion, "502")
	if solicitudAvance != nil {
		solicitudAvance.TerceroId = sol.TerceroId
	}
}

func ObtenerObjetivoPorSolicitudId(id int, solicitudAvance *models.SolicitudAvance) (observacion *solicitudes_crud.Observacion, err map[string]interface{}) {
	funcion := "ObtenerObjetivoPorSolicitudId"
	defer ErrorControlFunction(funcion, "502")
	o, e := ObtenerObservacionPorSolicitudId(id, 5, solicitudAvance)
	if e == nil && solicitudAvance != nil && o != nil {
		solicitudAvance.Objetivo = o.Valor
	}
	return o, e
}

func ObtenerJustificacionPorSolicitudId(id int, solicitudAvance *models.SolicitudAvance) (observacion *solicitudes_crud.Observacion, err map[string]interface{}) {
	funcion := "ObtenerJustificacionPorSolicitudId"
	defer ErrorControlFunction(funcion, "502")
	o, e := ObtenerObservacionPorSolicitudId(id, 6, solicitudAvance)
	if e == nil && solicitudAvance != nil && o != nil {
		solicitudAvance.Justificacion = o.Valor
	}
	return o, e
}

func ObtenerObservacionPorSolicitudId(id int, tipo int, solicitudAvance *models.SolicitudAvance) (observacion *solicitudes_crud.Observacion, err map[string]interface{}) {
	funcion := "ObtenerObservacionPorSolicitudId"
	defer ErrorControlFunction(funcion, "502")
	query := make(map[string]string)
	query["SolicitudId"] = strconv.Itoa(id)
	query["TipoObservacionId"] = strconv.Itoa(tipo)
	var observaciones []solicitudes_crud.Observacion
	if err := GetAll(&observaciones, "solicitudes_crud", "observacion", 1,
		query, nil, nil, nil, -1, -1); err == nil {
		if len(observaciones) > 0 {
			return &observaciones[0], nil
		} else {
			return nil, Error(funcion, "No existe solicitante asociado", "502")
		}
	} else {
		return nil, Error(funcion, err, "502")
	}
}
