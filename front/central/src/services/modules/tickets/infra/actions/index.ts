'use server';

import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { TicketApiRepository } from '../repository/api-repository';
import { TicketUseCases } from '../../app/use-cases';
import {
    CreateTicketDTO,
    UpdateTicketDTO,
    ListTicketsParams,
} from '../../domain/types';

async function getUseCases() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    const repository = new TicketApiRepository(token);
    return new TicketUseCases(repository);
}

export const listTicketsAction = async (params?: ListTicketsParams) => {
    return (await getUseCases()).list(params);
};

export const getTicketAction = async (id: number, businessId?: number) => {
    return (await getUseCases()).get(id, businessId);
};

export const createTicketAction = async (data: CreateTicketDTO) => {
    const r = await (await getUseCases()).create(data);
    revalidatePath('/tickets');
    return r;
};

export const updateTicketAction = async (id: number, data: UpdateTicketDTO) => {
    const r = await (await getUseCases()).update(id, data);
    revalidatePath('/tickets');
    return r;
};

export const deleteTicketAction = async (id: number) => {
    await (await getUseCases()).remove(id);
    revalidatePath('/tickets');
};

export const changeTicketStatusAction = async (id: number, status: string, note?: string) => {
    const r = await (await getUseCases()).changeStatus(id, status, note);
    revalidatePath('/tickets');
    return r;
};

export const assignTicketAction = async (id: number, assignedToId: number | null) => {
    const r = await (await getUseCases()).assign(id, assignedToId);
    revalidatePath('/tickets');
    return r;
};

export const changeTicketAreaAction = async (id: number, area: string, note?: string) => {
    const r = await (await getUseCases()).changeArea(id, area, note);
    revalidatePath('/tickets');
    return r;
};

export const escalateTicketAction = async (id: number, note?: string) => {
    const r = await (await getUseCases()).escalate(id, note);
    revalidatePath('/tickets');
    return r;
};

export const listCommentsAction = async (id: number, businessId?: number) => {
    return (await getUseCases()).listComments(id, businessId);
};

export const addCommentAction = async (id: number, body: string, isInternal: boolean) => {
    const r = await (await getUseCases()).addComment(id, body, isInternal);
    revalidatePath('/tickets');
    return r;
};

export const listAttachmentsAction = async (id: number, businessId?: number) => {
    return (await getUseCases()).listAttachments(id, businessId);
};

export const uploadAttachmentAction = async (id: number, formData: FormData) => {
    const file = formData.get('file') as File;
    const commentIdRaw = formData.get('comment_id');
    const commentId = commentIdRaw ? Number(commentIdRaw) : undefined;
    const r = await (await getUseCases()).uploadAttachment(id, file, commentId);
    revalidatePath('/tickets');
    return r;
};

export const deleteAttachmentAction = async (attachmentId: number) => {
    await (await getUseCases()).deleteAttachment(attachmentId);
    revalidatePath('/tickets');
};

export const listTicketHistoryAction = async (id: number, businessId?: number) => {
    return (await getUseCases()).listHistory(id, businessId);
};
