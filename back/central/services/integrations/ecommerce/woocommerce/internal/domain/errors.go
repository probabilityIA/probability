package domain

import "errors"

var (
	ErrIntegrationNotFound     = errors.New("integracion de WooCommerce no encontrada")
	ErrInvalidCredentials      = errors.New("credenciales invalidas: verifica el Consumer Key y el Consumer Secret de tu tienda")
	ErrMissingConsumerKey      = errors.New("falta el Consumer Key en las credenciales")
	ErrMissingConsumerSecret   = errors.New("falta el Consumer Secret en las credenciales")
	ErrMissingStoreURL         = errors.New("falta la URL de la tienda en la configuracion")
	ErrWebhookInvalidSignature = errors.New("la firma del webhook de WooCommerce no es valida")
	ErrWebhookMissingSecret    = errors.New("falta el secret del webhook de WooCommerce")
	ErrNoOrdersFound           = errors.New("no se encontraron ordenes en WooCommerce")
)
