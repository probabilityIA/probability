import {
    ShippingMargin,
    ShippingMarginsListResponse,
    GetShippingMarginsParams,
    CreateShippingMarginDTO,
    UpdateShippingMarginDTO,
} from './types';

export interface IShippingMarginRepository {
    list(params?: GetShippingMarginsParams): Promise<ShippingMarginsListResponse>;
    getById(id: number, businessId?: number): Promise<ShippingMargin>;
    create(data: CreateShippingMarginDTO, businessId?: number): Promise<ShippingMargin>;
    update(id: number, data: UpdateShippingMarginDTO, businessId?: number): Promise<ShippingMargin>;
}
