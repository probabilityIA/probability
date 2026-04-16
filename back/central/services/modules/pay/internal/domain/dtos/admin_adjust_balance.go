package dtos

// AdminAdjustBalanceDTO ajusta el saldo de una billetera (admin)
type AdminAdjustBalanceDTO struct {
	BusinessID uint
	Amount     float64 // Positivo = agregar, Negativo = restar
	Reference  string  // Motivo del ajuste
}
