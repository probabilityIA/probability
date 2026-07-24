import { IAccountingRepository } from '../domain/ports';
import {
    CreateConceptDTO,
    CreateEntryDTO,
    CreateTaxDTO,
    GetEntriesParams,
    UpdateConceptDTO,
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
}
