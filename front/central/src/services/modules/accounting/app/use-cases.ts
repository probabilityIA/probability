import { IAccountingRepository } from '../domain/ports';
import {
    CreateConceptDTO,
    CreateEntryDTO,
    CreateInvoiceDTO,
    CreateTaxDTO,
    EmitDianDTO,
    GetEntriesParams,
    GetInvoicesParams,
    SaveServiceDTO,
    UpdateConceptDTO,
    UpdateDianConfigDTO,
    UpdateInvoiceDTO,
    UpdateTaxDTO,
} from '../domain/types';

export class AccountingUseCases {
    constructor(private repository: IAccountingRepository) {}

    async getConcepts() {
        return this.repository.getConcepts();
    }

    async createConcept(data: CreateConceptDTO) {
        return this.repository.createConcept(data);
    }

    async updateConcept(id: number, data: UpdateConceptDTO) {
        return this.repository.updateConcept(id, data);
    }

    async setConceptTax(conceptId: number, taxId: number, isActive: boolean) {
        return this.repository.setConceptTax(conceptId, taxId, isActive);
    }

    async getTaxes() {
        return this.repository.getTaxes();
    }

    async createTax(data: CreateTaxDTO) {
        return this.repository.createTax(data);
    }

    async updateTax(id: number, data: UpdateTaxDTO) {
        return this.repository.updateTax(id, data);
    }

    async getServices() {
        return this.repository.getServices();
    }

    async createService(data: SaveServiceDTO) {
        return this.repository.createService(data);
    }

    async updateService(id: number, data: SaveServiceDTO) {
        return this.repository.updateService(id, data);
    }

    async deleteService(id: number) {
        return this.repository.deleteService(id);
    }

    async getEntries(params?: GetEntriesParams) {
        return this.repository.getEntries(params);
    }

    async createEntry(data: CreateEntryDTO) {
        return this.repository.createEntry(data);
    }

    async deleteEntry(id: number) {
        return this.repository.deleteEntry(id);
    }

    async getReport(from: string, to: string) {
        return this.repository.getReport(from, to);
    }

    async syncNow() {
        return this.repository.syncNow();
    }

    async getInvoices(params?: GetInvoicesParams) {
        return this.repository.getInvoices(params);
    }

    async getInvoice(id: number) {
        return this.repository.getInvoice(id);
    }

    async createInvoice(data: CreateInvoiceDTO) {
        return this.repository.createInvoice(data);
    }

    async updateInvoice(id: number, data: UpdateInvoiceDTO) {
        return this.repository.updateInvoice(id, data);
    }

    async deleteInvoice(id: number) {
        return this.repository.deleteInvoice(id);
    }

    async sendInvoice(id: number, emailTo?: string) {
        return this.repository.sendInvoice(id, emailTo);
    }

    async payInvoice(id: number) {
        return this.repository.payInvoice(id);
    }

    async cancelInvoice(id: number) {
        return this.repository.cancelInvoice(id);
    }

    async getDianConfig() {
        return this.repository.getDianConfig();
    }

    async updateDianConfig(data: UpdateDianConfigDTO) {
        return this.repository.updateDianConfig(data);
    }

    async emitInvoiceDian(id: number, data: EmitDianDTO) {
        return this.repository.emitInvoiceDian(id, data);
    }

    async getClientProfile(businessId: number) {
        return this.repository.getClientProfile(businessId);
    }
}
