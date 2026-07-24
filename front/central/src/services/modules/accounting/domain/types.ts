export type ConceptKind = 'INCOME' | 'EXPENSE';

export interface ConceptTax {
    id: number;
    tax_id: number;
    tax_code: string;
    tax_name: string;
    rate_percent: number;
    is_active: boolean;
}

export interface Concept {
    id: number;
    code: string;
    name: string;
    description: string;
    kind: ConceptKind;
    is_real_income: boolean;
    is_automatic: boolean;
    source_type: string;
    is_active: boolean;
    taxes: ConceptTax[] | null;
}

export interface Tax {
    id: number;
    code: string;
    name: string;
    description: string;
    rate_percent: number;
    is_active: boolean;
}

export interface EntryTaxDetail {
    tax_code: string;
    tax_name: string;
    rate_percent: number;
    amount: number;
}

export interface Entry {
    id: number;
    concept_id: number;
    concept_code: string;
    concept_name: string;
    business_id: number | null;
    business_name: string;
    entry_date: string;
    amount: number;
    kind: ConceptKind;
    source_type: string;
    source_id: string;
    description: string;
    is_automatic: boolean;
    tax_total: number;
    tax_detail: EntryTaxDetail[] | null;
}

export interface CreateConceptDTO {
    code: string;
    name: string;
    description: string;
    kind: ConceptKind;
    is_real_income: boolean;
}

export interface UpdateConceptDTO {
    name: string;
    description: string;
    is_real_income: boolean;
    is_active: boolean;
}

export interface CreateTaxDTO {
    code: string;
    name: string;
    description: string;
    rate_percent: number;
}

export interface UpdateTaxDTO {
    name: string;
    description: string;
    rate_percent: number;
    is_active: boolean;
}

export interface CreateEntryDTO {
    concept_id: number;
    related_business_id?: number;
    entry_date: string;
    amount: number;
    description: string;
}

export interface GetEntriesParams {
    page?: number;
    page_size?: number;
    from?: string;
    to?: string;
    concept_id?: number;
    kind?: ConceptKind;
}

export interface EntriesListResponse {
    data: Entry[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface ReportTotals {
    real_income: number;
    cash_in: number;
    real_expense: number;
    cash_out: number;
    tax_total: number;
    net: number;
}

export interface ReportConceptRow {
    concept_id: number;
    code: string;
    name: string;
    kind: ConceptKind;
    is_real_income: boolean;
    entries_count: number;
    amount: number;
    tax_total: number;
}

export interface ReportTaxRow {
    tax_code: string;
    tax_name: string;
    amount: number;
}

export interface AccountingReport {
    from: string;
    to: string;
    totals: ReportTotals;
    by_concept: ReportConceptRow[] | null;
    by_tax: ReportTaxRow[] | null;
}

export interface SyncResult {
    created: number;
    by_source: Record<string, number> | null;
}

export type InvoiceStatus = 'DRAFT' | 'SENT' | 'PAID' | 'CANCELLED';

export type DianStatus = 'NONE' | 'VALIDATED';

export interface InvoiceTaxDetail {
    tax_code: string;
    tax_name: string;
    rate_percent: number;
    amount: number;
}

export interface InvoiceItem {
    id: number;
    description: string;
    quantity: number;
    unit_price: number;
    amount: number;
}

export interface Invoice {
    id: number;
    number: string;
    business_id: number;
    business_name: string;
    concept_id: number;
    concept_name: string;
    issue_date: string;
    due_date: string;
    status: InvoiceStatus;
    notes: string;
    subtotal: number;
    tax_total: number;
    total: number;
    tax_detail: InvoiceTaxDetail[] | null;
    email_to: string;
    sent_at: string;
    paid_at: string;
    customer_document: string;
    customer_phone: string;
    customer_address: string;
    dian_status: DianStatus;
    cufe: string;
    dian_number: string;
    dian_qr: string;
    dian_emitted_at: string;
    is_test: boolean;
    items: InvoiceItem[] | null;
    created_at: string;
}

export interface InvoiceItemInputDTO {
    description: string;
    quantity: number;
    unit_price: number;
}

export interface CreateInvoiceDTO {
    business_id: number;
    concept_id: number;
    issue_date: string;
    due_date?: string;
    notes?: string;
    email_to?: string;
    customer_document?: string;
    customer_phone?: string;
    customer_address?: string;
    is_test?: boolean;
    tax_ids?: number[];
    items: InvoiceItemInputDTO[];
}

export interface UpdateInvoiceDTO {
    concept_id: number;
    issue_date: string;
    due_date?: string;
    notes?: string;
    email_to?: string;
    customer_document?: string;
    customer_phone?: string;
    customer_address?: string;
    is_test?: boolean;
    tax_ids?: number[];
    items: InvoiceItemInputDTO[];
}

export interface GetInvoicesParams {
    page?: number;
    page_size?: number;
    status?: InvoiceStatus;
    business_id?: number;
    is_test?: 'true' | 'false';
}

export interface InvoicesListResponse {
    data: Invoice[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface DianConfig {
    configured: boolean;
    integration_id: number;
    api_url: string;
    username: string;
    numbering_range_id: number;
    municipality_id: number;
    identification_document_id: number;
    legal_organization_id: number;
}

export interface UpdateDianConfigDTO {
    api_url?: string;
    client_id: string;
    client_secret: string;
    username: string;
    password: string;
    numbering_range_id?: number;
    municipality_id?: number;
    identification_document_id?: number;
    legal_organization_id?: number;
}

export interface EmitDianDTO {
    customer_document?: string;
    customer_phone?: string;
    customer_address?: string;
}

export type ActionResult<T> =
    | { success: true; data: T }
    | { success: false; error: string };
