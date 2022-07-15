package models

type BancosPorCodigo struct {
	IdTercero	int
	NIT			string
	CodigoAch 	int
	CodigoSuper	int
	NombreBanco string
	Activo		bool
}

type DatosCodigos struct {
	CodigoAch	int
	CodigoSuper	int
}
