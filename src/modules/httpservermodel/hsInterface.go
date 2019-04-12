package httpservermodel

import(
	hst "common/httpservertools"
)
type httpServerInterface interface{
	AddRegistrar(reg hst.Registrar) error
	AddHandlerFunc(path string, reg hst.HandlersFunc) error
	DoStartServer( ) error
}