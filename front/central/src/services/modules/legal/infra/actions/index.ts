'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { LegalApiRepository } from '../repository/api-repository';
import { AcceptLegalResult, PendingLegalDocuments } from '../../domain/types';

async function getRepository() {
    const token = await getAuthToken();
    return new LegalApiRepository(token);
}

export const getPendingLegalDocumentsAction = async (): Promise<PendingLegalDocuments> => {
    const repository = await getRepository();
    return repository.getPendingDocuments();
};

export const acceptLegalDocumentsAction = async (documentIds: number[]): Promise<AcceptLegalResult> => {
    const repository = await getRepository();
    return repository.acceptDocuments(documentIds);
};
