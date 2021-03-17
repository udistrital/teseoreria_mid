package helpers

import (
	"encoding/json"

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
		if _, errObj := CrearObservacion(solicitudAvance, solicitud, 5, "Objetivo"); errObj != nil {
			return nil, Error("CrearSolicitudAvance", errObj, "502")
		}
		if _, errJust := CrearObservacion(solicitudAvance, solicitud, 6, "Justificación"); errJust != nil {
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

func CrearObservacion(solicitudAvance *models.SolicitudAvance, solicitud *solicitudes_crud.Solicitud, tipoObservacionId int, titulo string) (obs *solicitudes_crud.Observacion, err map[string]interface{}) {
	defer ErrorControlFunction("CrearObservacion", "502")
	tipoObservacion := solicitudes_crud.TipoObservacion{Id: tipoObservacionId}
	observacion := solicitudes_crud.Observacion{
		TipoObservacionId: &tipoObservacion,
		SolicitudId:       solicitud,
		TerceroId:         solicitudAvance.TerceroId,
		Titulo:            titulo,
		Valor:             solicitudAvance.Objetivo,
		Activo:            true,
	}
	if _, err1 := AddObservacion(&observacion); err1 == nil {
		return &observacion, nil
	} else {
		return nil, Error("CrearObservacion", err1, "502")
	}
}

func AddSolicitudAvanceCrud(solicitud *avances_crud.SolicitudAvance) (id int64, err map[string]interface{}) {
	return Add(solicitud, "avances_crud", "solicitud_avance", "AddSolicitudAvanceCrud")
}

func AddSolicitud(solicitud *solicitudes_crud.Solicitud) (id int64, err map[string]interface{}) {
	return Add(solicitud, "solicitudes_crud", "solicitud", "AddSolicitud")
}

func AddSolicitante(solicitante *solicitudes_crud.Solicitante) (id int64, err map[string]interface{}) {
	return Add(solicitante, "solicitudes_crud", "solicitante", "AddSolicitante")
}

func AddObservacion(observacion *solicitudes_crud.Observacion) (id int64, err map[string]interface{}) {
	return Add(observacion, "solicitudes_crud", "observacion", "AddObservacion")
}

// OBTENCIÓN

// Obtiene una solicitud de avance buscando por Id
func ObtenerSolicitudAvancePorId(id int) (sol *models.SolicitudAvance, err map[string]interface{}) {
	defer ErrorControlFunction("ObtenerSolicitudAvancePorId", "502")
	solicitudAvance := models.SolicitudAvance{}
	// Consulta en avances_crud
	if _, err := ObtenerSolicitudAvanceCrudPorId(id, &solicitudAvance); err == nil {
		if _, err := ObtenerSolicitudPorId(solicitudAvance.SolicitudId, &solicitudAvance); err == nil {
			return &solicitudAvance, nil
		} else {
			return nil, Error("ObtenerSolicitudAvancePorId", err, "502")
		}
	} else {
		return nil, Error("ObtenerSolicitudAvancePorId", err, "502")
	}

	//

	// Guardado en solicitudes_crud
	// Guardar el solicitante
	/*if _, err2 := CrearSolicitante(solicitudAvance, solicitud); err2 != nil {
		return nil, Error("CrearSolicitudCrud", err2, "502")
	}
	// Guardar objetivo y justificacion (observaciones)
	if _, errObj := CrearObservacion(solicitudAvance, solicitud, 5, "Objetivo"); errObj != nil {
		return nil, Error("CrearSolicitudCrud", errObj, "502")
	}
	if _, errJust := CrearObservacion(solicitudAvance, solicitud, 6, "Justificación"); errJust != nil {
		return nil, Error("CrearSolicitudCrud", errJust, "502")
	}*/
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

func GetSolicitudAvanceCrudById(id int) (solicitud *avances_crud.SolicitudAvance, err map[string]interface{}) {
	_, err = GetById(id, &solicitud, "avances_crud", "solicitud_avance", "GetSolicitudAvanceCrudById")
	return solicitud, err
}

func GetSolicitudCrudById(id int) (solicitud *solicitudes_crud.Solicitud, err map[string]interface{}) {
	_, err = GetById(id, &solicitud, "solicitudes_crud", "solicitud", "GetSolicitudCrudById")
	return solicitud, err
}

/*func ObtenerSolicitudAvance(solicitudCrud solicitudes_crud.Solicitud, solicitudAvanceCrud avances_crud.SolicitudAvance, solicitante avances_crud.Solicitante, objetivo avances_crud.Observacion, justificacion avances_crud.Observacion) (solicitud models.SolicitudAvance, err error){

}*/
