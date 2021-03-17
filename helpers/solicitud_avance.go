package helpers

import (
	"encoding/json"
	"strconv"

	"github.com/udistrital/tesoreria_mid/models"

	avances_crud "github.com/udistrital/avances_crud/models"
	solicitudes_crud "github.com/udistrital/solicitudes_crud/models"
)

// CREACIÓN

func CrearSolicitudAvance(solicitudAvance *models.SolicitudAvance) (sol *models.SolicitudAvance, err map[string]interface{}) {
	defer ErrorControlFunction("CrearSolicitudAvance", "502")
	// Guardado en solicitudes_crud
	if solicitud, err1 := CrearSolicitudCrud(solicitudAvance); err1 == nil && solicitud != nil {
		// Guardar el solicitante
		if _, err2 := CrearSolicitante(solicitudAvance, solicitud); err2 != nil {
			return nil, Error("CrearSolicitudAvance", err2, "502")
		}
		// Guardar objetivo y justificacion (observaciones)
		if _, errObj := CrearObservacion(solicitudAvance, solicitud, solicitudAvance.Objetivo, 5, "Objetivo"); errObj != nil {
			return nil, Error("CrearSolicitudAvance", errObj, "502")
		}
		if _, errJust := CrearObservacion(solicitudAvance, solicitud, solicitudAvance.Justificacion, 6, "Justificación"); errJust != nil {
			return nil, Error("CrearSolicitudAvance", errJust, "502")
		}
		// Guardado en avances_crud
		solicitudAvanceCrud := avances_crud.SolicitudAvance{SolicitudId: solicitud.Id, Activo: true}
		if _, err3 := AddSolicitudAvanceCrud(&solicitudAvanceCrud); err3 == nil {
			solicitudAvance.Id = solicitudAvanceCrud.Id
			solicitudAvance.FechaCreacion = solicitudAvanceCrud.FechaCreacion
			solicitudAvance.FechaModificacion = solicitudAvanceCrud.FechaCreacion
			return solicitudAvance, nil
		} else {
			return nil, Error("CrearSolicitudAvance", err3, "502")
		}
	} else {
		return nil, Error("CrearSolicitudAvance", err1, "502")
	}
}

func CrearSolicitudCrud(solicitudAvance *models.SolicitudAvance) (sol *solicitudes_crud.Solicitud, err map[string]interface{}) {
	defer ErrorControlFunction("CrearSolicitudCrud", "502")
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
		solicitud := solicitudes_crud.Solicitud{
			EstadoTipoSolicitudId: solicitudAvance.EstadoTipoSolicitud,
			Referencia:            string(referencia),
			FechaRadicacion:       solicitudAvance.FechaRadicacion,
			SolicitudFinalizada:   solicitudAvance.SolicitudFinalizada,
			Resultado:             solicitudAvance.Resultado,
			Activo:                solicitudAvance.Activo,
		}
		if _, err1 := AddSolicitud(&solicitud); err1 == nil {
			solicitudAvance.SolicitudId = solicitud.Id
			return &solicitud, nil
		} else {
			return nil, Error("CrearSolicitudCrud", err1, "502")
		}
	}
	return nil, Error("CrearSolicitudCrud", errJSON, "400")
}

func CrearSolicitante(solicitudAvance *models.SolicitudAvance, solicitud *solicitudes_crud.Solicitud) (sol *solicitudes_crud.Solicitante, err map[string]interface{}) {
	defer ErrorControlFunction("CrearSolicitante", "502")
	solicitante := solicitudes_crud.Solicitante{
		TerceroId:   solicitudAvance.TerceroId,
		SolicitudId: solicitud,
		Activo:      true,
	}
	if _, err1 := AddSolicitante(&solicitante); err1 == nil {
		return &solicitante, nil
	} else {
		return nil, Error("CrearSolicitante", err1, "502")
	}
}

func CrearObservacion(solicitudAvance *models.SolicitudAvance, solicitud *solicitudes_crud.Solicitud, valor string, tipoObservacionId int, titulo string) (obs *solicitudes_crud.Observacion, err map[string]interface{}) {
	defer ErrorControlFunction("CrearObservacion", "502")
	tipoObservacion := solicitudes_crud.TipoObservacion{Id: tipoObservacionId}
	observacion := solicitudes_crud.Observacion{
		TipoObservacionId: &tipoObservacion,
		SolicitudId:       solicitud,
		TerceroId:         solicitudAvance.TerceroId,
		Titulo:            titulo,
		Valor:             valor,
		Activo:            true,
	}
	if _, err1 := AddObservacion(&observacion); err1 == nil {
		return &observacion, nil
	} else {
		return nil, Error("CrearObservacion", err1, "502")
	}
}

func AddSolicitudAvanceCrud(solicitud *avances_crud.SolicitudAvance) (id int64, err map[string]interface{}) {
	return Add(solicitud, "avances_crud", "solicitud_avance", 2, "AddSolicitudAvanceCrud")
}

func AddSolicitud(solicitud *solicitudes_crud.Solicitud) (id int64, err map[string]interface{}) {
	return Add(solicitud, "solicitudes_crud", "solicitud", 2, "AddSolicitud")
}

func AddSolicitante(solicitante *solicitudes_crud.Solicitante) (id int64, err map[string]interface{}) {
	return Add(solicitante, "solicitudes_crud", "solicitante", 2, "AddSolicitante")
}

func AddObservacion(observacion *solicitudes_crud.Observacion) (id int64, err map[string]interface{}) {
	return Add(observacion, "solicitudes_crud", "observacion", 2, "AddObservacion")
}

// OBTENCIÓN

// Obtiene una solicitud de avance buscando por Id
func ObtenerSolicitudAvancePorId(id int) (sol *models.SolicitudAvance, err map[string]interface{}) {
	defer ErrorControlFunction("ObtenerSolicitudAvancePorId", "502")
	solicitudAvance := models.SolicitudAvance{}
	// Consulta en avances_crud
	if _, err := ObtenerSolicitudAvanceCrudPorId(id, &solicitudAvance); err == nil {
		// Consulta en solicitudes_crud
		_, err := ObtenerSolicitudPorId(solicitudAvance.SolicitudId, &solicitudAvance)
		if err != nil {
			return nil, Error("ObtenerSolicitudAvancePorId", err, "502")
		}
		_, err = ObtenerSolicitantePorSolicitudId(solicitudAvance.SolicitudId, &solicitudAvance)
		if err != nil {
			return nil, Error("ObtenerSolicitudAvancePorId", err, "502")
		}
		_, err = ObtenerObjetivoPorSolicitudId(solicitudAvance.SolicitudId, &solicitudAvance)
		if err != nil {
			return nil, Error("ObtenerSolicitudAvancePorId", err, "502")
		}
		_, err = ObtenerJustificacionPorSolicitudId(solicitudAvance.SolicitudId, &solicitudAvance)
		if err != nil {
			return nil, Error("ObtenerSolicitudAvancePorId", err, "502")
		}
		return &solicitudAvance, nil
	} else {
		return nil, Error("ObtenerSolicitudAvancePorId", err, "502")
	}
}

func ObtenerSolicitudAvanceCrudPorId(id int, solicitudAvance *models.SolicitudAvance) (sol *avances_crud.SolicitudAvance, err map[string]interface{}) {
	defer ErrorControlFunction("ObtenerSolicitudAvanceCrudPorId", "502")
	if solicitudAvanceCrud, err := GetSolicitudAvanceCrudById(id); err == nil && solicitudAvanceCrud != nil {
		if solicitudAvance != nil {
			solicitudAvance.Id = solicitudAvanceCrud.Id
			solicitudAvance.SolicitudId = solicitudAvanceCrud.SolicitudId
			solicitudAvance.FechaCreacion = solicitudAvanceCrud.FechaCreacion
			solicitudAvance.FechaModificacion = solicitudAvanceCrud.FechaModificacion
			solicitudAvance.Activo = solicitudAvanceCrud.Activo
		}
		return solicitudAvanceCrud, nil
	} else {
		return nil, Error("ObtenerSolicitudAvanceCrudPorId", err, "502")
	}
}

func ObtenerSolicitudPorId(id int, solicitudAvance *models.SolicitudAvance) (sol *solicitudes_crud.Solicitud, err map[string]interface{}) {
	defer ErrorControlFunction("ObtenerSolicitudAvancePorId", "502")
	if solicitudAvanceCrud, err := GetSolicitudCrudById(id); err == nil && solicitudAvanceCrud != nil {
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
				return solicitudAvanceCrud, nil
			} else {
				return nil, Error("ObtenerSolicitudAvancePorId", err, "502")
			}
		}
		return solicitudAvanceCrud, nil
	} else {
		return nil, Error("ObtenerSolicitudAvancePorId", err, "502")
	}
}

func ObtenerSolicitantePorSolicitudId(id int, solicitudAvance *models.SolicitudAvance) (solicitante *solicitudes_crud.Solicitante, err map[string]interface{}) {
	defer ErrorControlFunction("ObtenerSolicitantePorSolicitudId", "502")
	query := make(map[string]string)
	query["SolicitudId"] = strconv.Itoa(id)
	if solicitantes, err := GetAllSolicitantes(query); err == nil {
		if len(solicitantes) > 0 {
			sol := solicitantes[0]
			if solicitudAvance != nil {
				solicitudAvance.TerceroId = sol.TerceroId
			}
			return &sol, nil
		} else {
			return nil, Error("ObtenerSolicitantePorSolicitudId", "No existe solicitante asociado", "502")
		}
	} else {
		return nil, Error("ObtenerSolicitantePorSolicitudId", err, "502")
	}
}

func ObtenerObjetivoPorSolicitudId(id int, solicitudAvance *models.SolicitudAvance) (observacion *solicitudes_crud.Observacion, err map[string]interface{}) {
	o, e := ObtenerObservacionPorSolicitudId(id, 5, solicitudAvance)
	if e == nil && solicitudAvance != nil && o != nil {
		solicitudAvance.Objetivo = o.Valor
	}
	return o, e
}

func ObtenerJustificacionPorSolicitudId(id int, solicitudAvance *models.SolicitudAvance) (observacion *solicitudes_crud.Observacion, err map[string]interface{}) {
	o, e := ObtenerObservacionPorSolicitudId(id, 6, solicitudAvance)
	if e == nil && solicitudAvance != nil && o != nil {
		solicitudAvance.Justificacion = o.Valor
	}
	return o, e
}

func ObtenerObservacionPorSolicitudId(id int, tipo int, solicitudAvance *models.SolicitudAvance) (observacion *solicitudes_crud.Observacion, err map[string]interface{}) {
	defer ErrorControlFunction("ObtenerObservacionPorSolicitudId", "502")
	query := make(map[string]string)
	query["SolicitudId"] = strconv.Itoa(id)
	query["TipoObservacionId"] = strconv.Itoa(tipo)
	if observaciones, err := GetAllObservaciones(query); err == nil {
		if len(observaciones) > 0 {
			return &observaciones[0], nil
		} else {
			return nil, Error("ObtenerObservacionPorSolicitudId", "No existe solicitante asociado", "502")
		}
	} else {
		return nil, Error("ObtenerObservacionPorSolicitudId", err, "502")
	}
}

func GetSolicitudAvanceCrudById(id int) (solicitud *avances_crud.SolicitudAvance, err map[string]interface{}) {
	_, err = GetById(id, &solicitud, "avances_crud", "solicitud_avance", 2, "GetSolicitudAvanceCrudById")
	return solicitud, err
}

func GetSolicitudCrudById(id int) (solicitud *solicitudes_crud.Solicitud, err map[string]interface{}) {
	_, err = GetById(id, &solicitud, "solicitudes_crud", "solicitud", 1, "GetSolicitudCrudById")
	return solicitud, err
}

func GetAllSolicitantes(query map[string]string) (solicitantes []solicitudes_crud.Solicitante, err map[string]interface{}) {
	_, err = GetAll(&solicitantes, "solicitudes_crud", "solicitante", 1, "GetAllSolicitantes", query, nil, nil, nil, 0, 0)
	return solicitantes, err
}

func GetAllObservaciones(query map[string]string) (observaciones []solicitudes_crud.Observacion, err map[string]interface{}) {
	_, err = GetAll(&observaciones, "solicitudes_crud", "observacion", 1, "GetAllObservaciones", query, nil, nil, nil, 0, 0)
	return observaciones, err
}
