package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

const checkoutReferencePrefix = "SFO"

// CreateCheckoutSession reprecia el carrito contra el catalogo real (nunca confia en el
// precio que manda el cliente), crea un PublicCheckout en estado pending, y genera la
// firma de Bold para que el frontend abra el widget de pago.
func (uc *UseCase) CreateCheckoutSession(ctx context.Context, slug string, dto *dtos.CreateCheckoutDTO) (*dtos.CheckoutSessionDTO, error) {
	if len(dto.Items) == 0 {
		return nil, domainerrors.ErrEmptyCart
	}
	for _, it := range dto.Items {
		if it.Quantity < 1 {
			return nil, domainerrors.ErrInvalidCartItem
		}
	}

	business, err := uc.repo.GetBusinessBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if business == nil {
		return nil, domainerrors.ErrBusinessNotFound
	}

	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, business.ID, tiendaWebIntegrationTypeID)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, domainerrors.ErrPublicSiteNotActive
	}

	integrationID, err := uc.repo.GetPlatformIntegrationID(ctx, business.ID)
	if err != nil {
		return nil, err
	}

	items := make([]entities.CheckoutItem, 0, len(dto.Items))
	var total float64
	for _, it := range dto.Items {
		product, perr := uc.repo.GetProductByID(ctx, business.ID, it.ProductID)
		if perr != nil {
			return nil, perr
		}
		lineTotal := product.Price * float64(it.Quantity)
		total += lineTotal
		items = append(items, entities.CheckoutItem{
			ProductID: product.ID,
			SKU:       product.SKU,
			Name:      product.Name,
			Quantity:  it.Quantity,
			UnitPrice: product.Price,
			ImageURL:  product.ImageURL,
		})
	}

	var address *entities.CheckoutAddress
	if dto.Address != nil {
		address = &entities.CheckoutAddress{
			Street:       dto.Address.Street,
			City:         dto.Address.City,
			State:        dto.Address.State,
			Country:      dto.Address.Country,
			PostalCode:   dto.Address.PostalCode,
			Instructions: dto.Address.Instructions,
		}
	}

	signature, err := uc.bold.BoldGenerateSignatureForReference(ctx, business.ID, total, "COP", checkoutReferencePrefix)
	if err != nil {
		return nil, err
	}

	checkout := &entities.PublicCheckout{
		BusinessID:      business.ID,
		IntegrationID:   integrationID,
		Slug:            slug,
		Reference:       signature.OrderID,
		Status:          entities.CheckoutStatusPending,
		Amount:          total,
		Currency:        "COP",
		Items:           items,
		CustomerName:    dto.CustomerName,
		CustomerEmail:   dto.CustomerEmail,
		CustomerPhone:   dto.CustomerPhone,
		CustomerDni:     dto.CustomerDni,
		ShippingAddress: address,
		BoldOrderID:     signature.OrderID,
	}
	if err := uc.repo.CreateCheckout(ctx, checkout); err != nil {
		return nil, err
	}

	return &dtos.CheckoutSessionDTO{
		Reference:      signature.OrderID,
		Amount:         signature.Amount,
		Currency:       signature.Currency,
		Hash:           signature.Hash,
		PublicKey:      signature.PublicKey,
		RedirectionURL: signature.RedirectionURL,
		IsSandbox:      signature.IsSandbox,
		PollingEnabled: signature.PollingEnabled,
	}, nil
}

// GetCheckoutStatus permite al frontend consultar (via polling) si el checkout ya fue
// confirmado por el webhook de Bold.
func (uc *UseCase) GetCheckoutStatus(ctx context.Context, reference string) (string, error) {
	checkout, err := uc.repo.GetCheckoutByReference(ctx, reference)
	if err != nil {
		return "", err
	}
	if checkout == nil {
		return "", domainerrors.ErrCheckoutNotFound
	}
	return checkout.Status, nil
}
