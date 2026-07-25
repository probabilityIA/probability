import {
    AccountingReport,
    AccountingService,
    ClientProfile,
    Concept,
    CreateConceptDTO,
    CreateEntryDTO,
    CreateInvoiceDTO,
    CreateTaxDTO,
    DianConfig,
    EmitDianDTO,
    EntriesListResponse,
    Entry,
    GetEntriesParams,
    GetInvoicesParams,
    Invoice,
    InvoicesListResponse,
    SaveServiceDTO,
    SyncResult,
    Tax,
    UpdateConceptDTO,
    UpdateDianConfigDTO,
    UpdateInvoiceDTO,
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
    getServices(): Promise<AccountingService[]>;
    createService(data: SaveServiceDTO): Promise<AccountingService>;
    updateService(id: number, data: SaveServiceDTO): Promise<AccountingService>;
    deleteService(id: number): Promise<void>;
    getEntries(params?: GetEntriesParams): Promise<EntriesListResponse>;
    createEntry(data: CreateEntryDTO): Promise<Entry>;
    deleteEntry(id: number): Promise<void>;
    getReport(from: string, to: string): Promise<AccountingReport>;
    syncNow(): Promise<SyncResult>;
    getInvoices(params?: GetInvoicesParams): Promise<InvoicesListResponse>;
    getInvoice(id: number): Promise<Invoice>;
    createInvoice(data: CreateInvoiceDTO): Promise<Invoice>;
    updateInvoice(id: number, data: UpdateInvoiceDTO): Promise<Invoice>;
    deleteInvoice(id: number): Promise<void>;
    sendInvoice(id: number, emailTo?: string): Promise<Invoice>;
    payInvoice(id: number): Promise<Invoice>;
    cancelInvoice(id: number): Promise<Invoice>;
    getDianConfig(): Promise<DianConfig>;
    updateDianConfig(data: UpdateDianConfigDTO): Promise<DianConfig>;
    emitInvoiceDian(id: number, data: EmitDianDTO): Promise<Invoice>;
    getClientProfile(businessId: number): Promise<ClientProfile>;
}
