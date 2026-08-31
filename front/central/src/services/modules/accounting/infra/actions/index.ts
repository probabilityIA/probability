'use server';

import { revalidatePath } from 'next/cache';
import { cookies } from 'next/headers';
import { AccountingApiRepository } from '../repository/api-repository';
import { AccountingUseCases } from '../../app/use-cases';
import {
    ActionResult,
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
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new AccountingApiRepository(token);
    return new AccountingUseCases(repository);
}

function errorMessage(error: unknown): string {
    return error instanceof Error ? error.message : 'Ocurrió un error';
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

function revalidateServices() {
    revalidatePath('/accounting/configuracion');
    revalidatePath('/accounting/facturas');
}

export const getAccountingServicesAction = async () =>
    run(async () => (await getUseCases()).getServices());

export const createAccountingServiceAction = async (data: SaveServiceDTO) =>
    run(async () => {
        const result = await (await getUseCases()).createService(data);
        revalidateServices();
        return result;
    });

export const updateAccountingServiceAction = async (id: number, data: SaveServiceDTO) =>
    run(async () => {
        const result = await (await getUseCases()).updateService(id, data);
        revalidateServices();
        return result;
    });

export const deleteAccountingServiceAction = async (id: number) =>
    run(async () => {
        await (await getUseCases()).deleteService(id);
        revalidateServices();
        return true;
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

function revalidateInvoices(id?: number) {
    revalidatePath('/accounting/facturas');
    if (id) {
        revalidatePath(`/accounting/facturas/${id}`);
    }
}

export const getAccountingInvoicesAction = async (params?: GetInvoicesParams) =>
    run(async () => (await getUseCases()).getInvoices(params));

export const getAccountingInvoiceAction = async (id: number) =>
    run(async () => (await getUseCases()).getInvoice(id));

export const createAccountingInvoiceAction = async (data: CreateInvoiceDTO) =>
    run(async () => {
        const result = await (await getUseCases()).createInvoice(data);
        revalidateInvoices();
        return result;
    });

export const updateAccountingInvoiceAction = async (id: number, data: UpdateInvoiceDTO) =>
    run(async () => {
        const result = await (await getUseCases()).updateInvoice(id, data);
        revalidateInvoices(id);
        return result;
    });

export const deleteAccountingInvoiceAction = async (id: number) =>
    run(async () => {
        await (await getUseCases()).deleteInvoice(id);
        revalidateInvoices(id);
        return true;
    });

export const sendAccountingInvoiceAction = async (id: number, emailTo?: string) =>
    run(async () => {
        const result = await (await getUseCases()).sendInvoice(id, emailTo);
        revalidateInvoices(id);
        return result;
    });

export const payAccountingInvoiceAction = async (id: number) =>
    run(async () => {
        const result = await (await getUseCases()).payInvoice(id);
        revalidateInvoices(id);
        revalidatePath('/accounting');
        revalidatePath('/accounting/movimientos');
        return result;
    });

export const cancelAccountingInvoiceAction = async (id: number) =>
    run(async () => {
        const result = await (await getUseCases()).cancelInvoice(id);
        revalidateInvoices(id);
        revalidatePath('/accounting');
        revalidatePath('/accounting/movimientos');
        return result;
    });

export const getAccountingDianConfigAction = async () =>
    run(async () => (await getUseCases()).getDianConfig());

export const updateAccountingDianConfigAction = async (data: UpdateDianConfigDTO) =>
    run(async () => {
        const result = await (await getUseCases()).updateDianConfig(data);
        revalidatePath('/accounting/configuracion');
        return result;
    });

export const emitAccountingInvoiceDianAction = async (id: number, data: EmitDianDTO) =>
    run(async () => {
        const result = await (await getUseCases()).emitInvoiceDian(id, data);
        revalidateInvoices(id);
        return result;
    });

export const getAccountingClientProfileAction = async (businessId: number) =>
    run(async () => (await getUseCases()).getClientProfile(businessId));
