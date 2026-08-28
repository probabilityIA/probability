export interface LegalDocument {
    id: number;
    code: string;
    version: string;
    title: string;
    file_url: string;
    sha256: string;
    effective_from: string;
}

export interface PendingLegalDocuments {
    requires_acceptance: boolean;
    documents: LegalDocument[];
}

export interface AcceptLegalResult {
    accepted_at: string;
    document_ids: number[];
}
