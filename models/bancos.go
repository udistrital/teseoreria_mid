package models

type BancosPorCodigo struct {
	NIT			string
	CodigoAch 	int
	CodigoSuper	int
	NombreBanco string
}

type DatosCodigos struct {
	CodigoAch	int
	CodigoSuper	int
}
