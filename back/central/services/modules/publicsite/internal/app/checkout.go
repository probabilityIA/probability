package app

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

const agreedReferencePrefix = "SFA"

func (uc *UseCase) CreateCheckoutSession(ctx context.Context, slug string, dto *dtos.CreateCheckoutDTO, userID uint) (*dtos.CheckoutSessionDTO, error) {
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

	active, err := uc.repo.IsIntegrationActive(ctx, business.ID, tiendaWebIntegrationTypeID)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, domainerrors.ErrPublicSiteNotActive
	}

	integrationID, err := uc.repo.GetTiendaIntegrationID(ctx, business.ID)
	if err != nil {
		return nil, err
	}

	customPrices := map[string]float64{}
	if userID > 0 {
		session, serr := uc.repo.GetClientSession(ctx, business.ID, userID)
		if serr == nil && session != nil && session.CustomerID > 0 {
			if prices, perr := uc.repo.GetCustomerPrices(ctx, business.ID, session.CustomerID); perr == nil {
				customPrices = prices
			}
		}
	}

	items := make([]entities.CheckoutItem, 0, len(dto.Items))
	var total float64
	for _, it := range dto.Items {
		product, perr := uc.repo.GetProductByID(ctx, business.ID, it.ProductID)
		if perr != nil {
			return nil, perr
		}
		unitPrice := product.Price
		if custom, ok := customPrices[product.ID]; ok {
			unitPrice = custom
		}
		lineTotal := unitPrice * float64(it.Quantity)
		total += lineTotal
		items = append(items, entities.CheckoutItem{
			ProductID: product.ID,
			SKU:       product.SKU,
			Name:      product.Name,
			Quantity:  it.Quantity,
			UnitPrice: unitPrice,
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

	if dto.PaymentMethod == "online" {
		return nil, domainerrors.ErrOnlinePayNotReady
	}

	reference := agreedReferencePrefix + strings.ReplaceAll(uuid.New().String(), "-", "")[:20]

	checkout := &entities.PublicCheckout{
		BusinessID:      business.ID,
		IntegrationID:   integrationID,
		Slug:            slug,
		Reference:       reference,
		Status:          entities.CheckoutStatusAgreed,
		Amount:          total,
		Currency:        "COP",
		Items:           items,
		CustomerName:    dto.CustomerName,
		CustomerEmail:   dto.CustomerEmail,
		CustomerPhone:   dto.CustomerPhone,
		CustomerDni:     dto.CustomerDni,
		ShippingAddress: address,
		BoldOrderID:     reference,
	}
	if err := uc.repo.CreateCheckout(ctx, checkout); err != nil {
		return nil, err
	}

	if err := uc.orders.PublishAgreedStorefrontOrder(ctx, reference); err != nil {
		uc.logger.Error(ctx).Err(err).Str("reference", reference).Msg("error publicando orden acordada")
		return nil, err
	}

	return &dtos.CheckoutSessionDTO{
		PaymentMethod: "agree",
		Reference:     reference,
		Amount:        total,
		Currency:      "COP",
	}, nil
}

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
