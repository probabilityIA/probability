'use server';

import { revalidatePath } from 'next/cache';
import { cookies } from 'next/headers';
import { AccountingApiRepository } from '../repository/api-repository';
import { AccountingUseCases } from '../../app/use-cases';
import {
    ActionResult,
    CreateConceptDTO,
    CreateEntryDTO,
    CreateTaxDTO,
    GetEntriesParams,
    UpdateConceptDTO,
    UpdateTaxDTO,
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new AccountingApiRepository(token);
    return new AccountingUseCases(repository);
}

function errorMessage(error: unknown): string {
    return error instanceof Error ? error.message : 'Ocurrio un error';
}

async function run<T>(fn: () => Promise<T>): Promise<ActionResult<T>> {
    try {
        return { success: true, data: await fn() };
    } catch (error) {
        return { success: false, error: errorMessage(error) };
    }
}

export const getAccountingConceptsAction = async () =>
    run(async () => (await getUseCases()).getConcepts());

export const createAccountingConceptAction = async (data: CreateConceptDTO) =>
    run(async () => {
        const result = await (await getUseCases()).createConcept(data);
        revalidatePath('/accounting/configuracion');
        return result;
    });

export const updateAccountingConceptAction = async (id: number, data: UpdateConceptDTO) =>
    run(async () => {
        const result = await (await getUseCases()).updateConcept(id, data);
        revalidatePath('/accounting/configuracion');
        return result;
    });

export const setAccountingConceptTaxAction = async (conceptId: number, taxId: number, isActive: boolean) =>
    run(async () => {
        await (await getUseCases()).setConceptTax(conceptId, taxId, isActive);
        revalidatePath('/accounting/configuracion');
        return true;
    });

export const getAccountingTaxesAction = async () =>
    run(async () => (await getUseCases()).getTaxes());

export const createAccountingTaxAction = async (data: CreateTaxDTO) =>
    run(async () => {
        const result = await (await getUseCases()).createTax(data);
        revalidatePath('/accounting/configuracion');
        return result;
    });

export const updateAccountingTaxAction = async (id: number, data: UpdateTaxDTO) =>
    run(async () => {
        const result = await (await getUseCases()).updateTax(id, data);
        revalidatePath('/accounting/configuracion');
        return result;
    });

export const getAccountingEntriesAction = async (params?: GetEntriesParams) =>
    run(async () => (await getUseCases()).getEntries(params));

export const createAccountingEntryAction = async (data: CreateEntryDTO) =>
    run(async () => {
        const result = await (await getUseCases()).createEntry(data);
        revalidatePath('/accounting/movimientos');
        revalidatePath('/accounting');
        return result;
    });

export const deleteAccountingEntryAction = async (id: number) =>
    run(async () => {
        await (await getUseCases()).deleteEntry(id);
        revalidatePath('/accounting/movimientos');
        revalidatePath('/accounting');
        return true;
    });

export const getAccountingReportAction = async (from: string, to: string) =>
    run(async () => (await getUseCases()).getReport(from, to));

export const syncAccountingAction = async () =>
    run(async () => {
        const result = await (await getUseCases()).syncNow();
        revalidatePath('/accounting');
        revalidatePath('/accounting/movimientos');
        return result;
    });
