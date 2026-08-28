package errors

import "errors"

var (
	ErrUserRequired         = errors.New("usuario no identificado")
	ErrNoDocumentsSelected  = errors.New("no se indico ningun documento para aceptar")
	ErrDocumentNotAvailable = errors.New("uno de los documentos no existe o ya no esta vigente")
	ErrPendingDocuments     = errors.New("faltan documentos legales por aceptar")
)
