import { AcceptLegalResult, PendingLegalDocuments } from './types';

export interface ILegalRepository {
    getPendingDocuments(): Promise<PendingLegalDocuments>;
    acceptDocuments(documentIds: number[]): Promise<AcceptLegalResult>;
}
