import { IShippingMarginRepository } from '../domain/ports';
import { GetShippingMarginsParams, CreateShippingMarginDTO, UpdateShippingMarginDTO, ProfitReportParams, ProfitReportDetailParams } from '../domain/types';

export class ShippingMarginUseCases {
    constructor(private repository: IShippingMarginRepository) {}

    async list(params?: GetShippingMarginsParams) {
        return this.repository.list(params);
    }

    async getById(id: number, businessId?: number) {
        return this.repository.getById(id, businessId);
    }

    async create(data: CreateShippingMarginDTO, businessId?: number) {
        return this.repository.create(data, businessId);
    }

    async update(id: number, data: UpdateShippingMarginDTO, businessId?: number) {
        return this.repository.update(id, data, businessId);
    }

    async profitReport(params: ProfitReportParams) {
        return this.repository.profitReport(params);
    }

    async profitReportDetail(params: ProfitReportDetailParams) {
        return this.repository.profitReportDetail(params);
    }
}
