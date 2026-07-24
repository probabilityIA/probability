import {
    AccountingReport,
    Concept,
    CreateConceptDTO,
    CreateEntryDTO,
    CreateTaxDTO,
    EntriesListResponse,
    Entry,
    GetEntriesParams,
    SyncResult,
    Tax,
    UpdateConceptDTO,
    UpdateTaxDTO,
} from './types';

export interface IAccountingRepository {
    getConcepts(): Promise<Concept[]>;
    createConcept(data: CreateConceptDTO): Promise<Concept>;
    updateConcept(id: number, data: UpdateConceptDTO): Promise<Concept>;
    setConceptTax(conceptId: number, taxId: number, isActive: boolean): Promise<void>;
    getTaxes(): Promise<Tax[]>;
    createTax(data: CreateTaxDTO): Promise<Tax>;
    updateTax(id: number, data: UpdateTaxDTO): Promise<Tax>;
    getEntries(params?: GetEntriesParams): Promise<EntriesListResponse>;
    createEntry(data: CreateEntryDTO): Promise<Entry>;
    deleteEntry(id: number): Promise<void>;
    getReport(from: string, to: string): Promise<AccountingReport>;
    syncNow(): Promise<SyncResult>;
}
