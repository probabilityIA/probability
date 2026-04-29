import {
    ShippingMargin,
    ShippingMarginsListResponse,
    GetShippingMarginsParams,
    CreateShippingMarginDTO,
    UpdateShippingMarginDTO,
    ProfitReportParams,
    ProfitReportResponse,
    ProfitReportDetailParams,
    ProfitReportDetailResponse,
} from './types';

export interface IShippingMarginRepository {
    list(params?: GetShippingMarginsParams): Promise<ShippingMarginsListResponse>;
    getById(id: number, businessId?: number): Promise<ShippingMargin>;
    create(data: CreateShippingMarginDTO, businessId?: number): Promise<ShippingMargin>;
    update(id: number, data: UpdateShippingMarginDTO, businessId?: number): Promise<ShippingMargin>;
    profitReport(params: ProfitReportParams): Promise<ProfitReportResponse>;
    profitReportDetail(params: ProfitReportDetailParams): Promise<ProfitReportDetailResponse>;
}
