package models

import (
	solicitudes_crud "github.com/udistrital/solicitudes_crud/models"
)

type SolicitudAvance struct {
	Id                    int                                   `orm:"column(id);pk;auto"`
	SolicitudId           int                                   `orm:"column(solicitud_id)"`
	EstadoTipoSolicitud   *solicitudes_crud.EstadoTipoSolicitud `orm:"column(estado_tipo_solicitud_id);rel(fk)"`
	FechaRadicacion       string                                `orm:"column(fecha_radicacion);type(timestamp without time zone);null"`
	SolicitudFinalizada   bool                                  `orm:"column(solicitud_finalizada);null"`
	Resultado             string                                `orm:"colum n(resultado);type(json);null"`
	VigenciaId            int
	AreaFuncionalId       int
	CargoOrdenadorGastoId int
	Objetivo              string
	Justificacion         string
	DependenciaId         int
	FacultadId            int
	ProyectoCurricularId  int
	ConvenioId            int
	ProyectoInvId         int
	TerceroId             int
	FechaEvento           string
	FechaCreacion         string `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion     string `orm:"column(fecha_modificacion);type(timestamp without time zone)"`
	Activo                bool   `orm:"column(activo)"`
}
