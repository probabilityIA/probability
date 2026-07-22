package app

import "context"

func (uc *UseCase) EcommerceChannelLimit(ctx context.Context, businessID uint) (int, error) {
	subTypeID, err := uc.repo.GetBusinessCurrentSubscriptionTypeID(ctx, businessID)
	if err != nil {
		return 0, err
	}
	if subTypeID == nil {
		return 0, nil
	}

	subType, err := uc.repo.GetSubscriptionType(ctx, *subTypeID)
	if err != nil {
		return 0, err
	}
	if subType == nil {
		return 0, nil
	}

	return subType.MaxEcommerceChannels, nil
}
