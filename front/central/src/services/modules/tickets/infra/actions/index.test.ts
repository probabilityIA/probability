import { describe, it, expect, vi, beforeEach } from 'vitest';

const mocks = vi.hoisted(() => {
    const useCases = {
        list: vi.fn(),
        get: vi.fn(),
        create: vi.fn(),
        update: vi.fn(),
        remove: vi.fn(),
        changeStatus: vi.fn(),
        assign: vi.fn(),
        changeArea: vi.fn(),
        escalate: vi.fn(),
        changeSprint: vi.fn(),
        listComments: vi.fn(),
        addComment: vi.fn(),
        listAttachments: vi.fn(),
        uploadAttachment: vi.fn(),
        deleteAttachment: vi.fn(),
        listHistory: vi.fn(),
    };
    return {
        useCases,
        useCasesCtor: vi.fn(),
        repositoryCtor: vi.fn(),
        cookieGet: vi.fn(),
        revalidatePath: vi.fn(),
    };
});

vi.mock('next/headers', () => ({
    cookies: vi.fn(async () => ({ get: mocks.cookieGet })),
}));

vi.mock('next/cache', () => ({
    revalidatePath: mocks.revalidatePath,
}));

vi.mock('../repository/api-repository', () => ({
    TicketApiRepository: class {
        constructor(token: string | null) {
            mocks.repositoryCtor(token);
            return { marker: 'repo' };
        }
    },
}));

vi.mock('../../app/use-cases', () => ({
    TicketUseCases: class {
        constructor(repo: unknown) {
            mocks.useCasesCtor(repo);
            return mocks.useCases;
        }
    },
}));

import {
    listTicketsAction,
    getTicketAction,
    createTicketAction,
    updateTicketAction,
    deleteTicketAction,
    changeTicketStatusAction,
    assignTicketAction,
    changeTicketAreaAction,
    changeTicketSprintAction,
    escalateTicketAction,
    listCommentsAction,
    addCommentAction,
    listAttachmentsAction,
    uploadAttachmentAction,
    deleteAttachmentAction,
    listTicketHistoryAction,
} from './index';

const TICKETS_PATH = '/tickets';

describe('acciones de tickets', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        mocks.cookieGet.mockReturnValue({ value: 'tok-123' });
    });

    describe('construccion del caso de uso', () => {
        it('toma el token de la cookie session_token', async () => {
            mocks.useCases.list.mockResolvedValue({ data: [] });

            await listTicketsAction();

            expect(mocks.cookieGet).toHaveBeenCalledWith('session_token');
            expect(mocks.repositoryCtor).toHaveBeenCalledWith('tok-123');
        });

        it('pasa null cuando no existe la cookie', async () => {
            mocks.cookieGet.mockReturnValue(undefined);
            mocks.useCases.list.mockResolvedValue({ data: [] });

            await listTicketsAction();

            expect(mocks.repositoryCtor).toHaveBeenCalledWith(null);
        });

        it('pasa null cuando la cookie viene vacia', async () => {
            mocks.cookieGet.mockReturnValue({ value: '' });
            mocks.useCases.list.mockResolvedValue({ data: [] });

            await listTicketsAction();

            expect(mocks.repositoryCtor).toHaveBeenCalledWith(null);
        });

        it('inyecta el repositorio en el caso de uso', async () => {
            mocks.useCases.list.mockResolvedValue({ data: [] });

            await listTicketsAction();

            expect(mocks.useCasesCtor).toHaveBeenCalledWith({ marker: 'repo' });
        });
    });

    describe('acciones de solo lectura', () => {
        it('listTicketsAction delega los filtros y no revalida', async () => {
            const page = { data: [], total: 0, page: 1, page_size: 10, total_pages: 0 };
            mocks.useCases.list.mockResolvedValue(page);

            const result = await listTicketsAction({ page: 2, status: 'open' });

            expect(result).toBe(page);
            expect(mocks.useCases.list).toHaveBeenCalledWith({ page: 2, status: 'open' });
            expect(mocks.revalidatePath).not.toHaveBeenCalled();
        });

        it('getTicketAction delega id y business', async () => {
            const ticket = { id: 5 };
            mocks.useCases.get.mockResolvedValue(ticket);

            const result = await getTicketAction(5, 26);

            expect(result).toBe(ticket);
            expect(mocks.useCases.get).toHaveBeenCalledWith(5, 26);
            expect(mocks.revalidatePath).not.toHaveBeenCalled();
        });

        it('listCommentsAction delega id y business', async () => {
            const comments = [{ id: 1 }];
            mocks.useCases.listComments.mockResolvedValue(comments);

            const result = await listCommentsAction(5, 26);

            expect(result).toBe(comments);
            expect(mocks.useCases.listComments).toHaveBeenCalledWith(5, 26);
        });

        it('listAttachmentsAction delega id y business', async () => {
            const attachments = [{ id: 20 }];
            mocks.useCases.listAttachments.mockResolvedValue(attachments);

            const result = await listAttachmentsAction(5, 26);

            expect(result).toBe(attachments);
            expect(mocks.useCases.listAttachments).toHaveBeenCalledWith(5, 26);
        });

        it('listTicketHistoryAction delega id y business', async () => {
            const history = [{ id: 30 }];
            mocks.useCases.listHistory.mockResolvedValue(history);

            const result = await listTicketHistoryAction(5, 26);

            expect(result).toBe(history);
            expect(mocks.useCases.listHistory).toHaveBeenCalledWith(5, 26);
        });
    });

    describe('acciones de escritura', () => {
        it('createTicketAction crea, revalida y retorna el ticket', async () => {
            const created = { id: 9 };
            mocks.useCases.create.mockResolvedValue(created);
            const dto = { title: 'T', description: 'D' };

            const result = await createTicketAction(dto);

            expect(result).toBe(created);
            expect(mocks.useCases.create).toHaveBeenCalledWith(dto);
            expect(mocks.revalidatePath).toHaveBeenCalledWith(TICKETS_PATH);
        });

        it('updateTicketAction actualiza, revalida y retorna el ticket', async () => {
            const updated = { id: 5 };
            mocks.useCases.update.mockResolvedValue(updated);

            const result = await updateTicketAction(5, { title: 'X' });

            expect(result).toBe(updated);
            expect(mocks.useCases.update).toHaveBeenCalledWith(5, { title: 'X' });
            expect(mocks.revalidatePath).toHaveBeenCalledWith(TICKETS_PATH);
        });

        it('deleteTicketAction elimina, revalida y no retorna nada', async () => {
            mocks.useCases.remove.mockResolvedValue(undefined);

            const result = await deleteTicketAction(5);

            expect(result).toBeUndefined();
            expect(mocks.useCases.remove).toHaveBeenCalledWith(5);
            expect(mocks.revalidatePath).toHaveBeenCalledWith(TICKETS_PATH);
        });

        it('changeTicketStatusAction delega estado y nota', async () => {
            const moved = { id: 5, status: 'closed' };
            mocks.useCases.changeStatus.mockResolvedValue(moved);

            const result = await changeTicketStatusAction(5, 'closed', 'listo');

            expect(result).toBe(moved);
            expect(mocks.useCases.changeStatus).toHaveBeenCalledWith(5, 'closed', 'listo');
            expect(mocks.revalidatePath).toHaveBeenCalledWith(TICKETS_PATH);
        });

        it('assignTicketAction acepta null para desasignar', async () => {
            mocks.useCases.assign.mockResolvedValue({ id: 5 });

            await assignTicketAction(5, null);

            expect(mocks.useCases.assign).toHaveBeenCalledWith(5, null);
            expect(mocks.revalidatePath).toHaveBeenCalledWith(TICKETS_PATH);
        });

        it('changeTicketAreaAction delega area y nota', async () => {
            const moved = { id: 5, area: 'desarrollo' };
            mocks.useCases.changeArea.mockResolvedValue(moved);

            const result = await changeTicketAreaAction(5, 'desarrollo', 'necesita dev');

            expect(result).toBe(moved);
            expect(mocks.useCases.changeArea).toHaveBeenCalledWith(5, 'desarrollo', 'necesita dev');
            expect(mocks.revalidatePath).toHaveBeenCalledWith(TICKETS_PATH);
        });

        it('changeTicketSprintAction delega el sprint', async () => {
            mocks.useCases.changeSprint.mockResolvedValue({ id: 5 });

            await changeTicketSprintAction(5, 4);

            expect(mocks.useCases.changeSprint).toHaveBeenCalledWith(5, 4);
            expect(mocks.revalidatePath).toHaveBeenCalledWith(TICKETS_PATH);
        });

        it('escalateTicketAction delega la nota', async () => {
            const escalated = { id: 5, escalated_to_dev: true };
            mocks.useCases.escalate.mockResolvedValue(escalated);

            const result = await escalateTicketAction(5, 'urgente');

            expect(result).toBe(escalated);
            expect(mocks.useCases.escalate).toHaveBeenCalledWith(5, 'urgente');
            expect(mocks.revalidatePath).toHaveBeenCalledWith(TICKETS_PATH);
        });

        it('addCommentAction delega cuerpo y visibilidad', async () => {
            const comment = { id: 11 };
            mocks.useCases.addComment.mockResolvedValue(comment);

            const result = await addCommentAction(5, 'nota interna', true);

            expect(result).toBe(comment);
            expect(mocks.useCases.addComment).toHaveBeenCalledWith(5, 'nota interna', true);
            expect(mocks.revalidatePath).toHaveBeenCalledWith(TICKETS_PATH);
        });

        it('deleteAttachmentAction elimina el adjunto y revalida', async () => {
            mocks.useCases.deleteAttachment.mockResolvedValue(undefined);

            const result = await deleteAttachmentAction(20);

            expect(result).toBeUndefined();
            expect(mocks.useCases.deleteAttachment).toHaveBeenCalledWith(20);
            expect(mocks.revalidatePath).toHaveBeenCalledWith(TICKETS_PATH);
        });
    });

    describe('uploadAttachmentAction', () => {
        const buildFormData = (commentId?: string) => {
            const fd = new FormData();
            fd.append('file', new File(['x'], 'evidencia.png', { type: 'image/png' }));
            if (commentId !== undefined) fd.append('comment_id', commentId);
            return fd;
        };

        it('extrae el archivo del FormData y sube sin comentario', async () => {
            const attachment = { id: 21 };
            mocks.useCases.uploadAttachment.mockResolvedValue(attachment);
            const fd = buildFormData();

            const result = await uploadAttachmentAction(5, fd);

            expect(result).toBe(attachment);
            expect(mocks.useCases.uploadAttachment).toHaveBeenCalledWith(5, fd.get('file'), undefined);
            expect(mocks.revalidatePath).toHaveBeenCalledWith(TICKETS_PATH);
        });

        it('convierte comment_id a numero', async () => {
            mocks.useCases.uploadAttachment.mockResolvedValue({ id: 22 });

            await uploadAttachmentAction(5, buildFormData('10'));

            expect(mocks.useCases.uploadAttachment).toHaveBeenCalledWith(5, expect.any(File), 10);
        });

        it('ignora comment_id vacio', async () => {
            mocks.useCases.uploadAttachment.mockResolvedValue({ id: 23 });

            await uploadAttachmentAction(5, buildFormData(''));

            expect(mocks.useCases.uploadAttachment).toHaveBeenCalledWith(5, expect.any(File), undefined);
        });
    });

    describe('propagacion de errores', () => {
        it('no revalida si la creacion falla', async () => {
            mocks.useCases.create.mockRejectedValue(new Error('titulo requerido'));

            await expect(createTicketAction({ title: '', description: '' })).rejects.toThrow('titulo requerido');
            expect(mocks.revalidatePath).not.toHaveBeenCalled();
        });

        it('propaga el error de lectura', async () => {
            mocks.useCases.get.mockRejectedValue(new Error('HTTP 404'));

            await expect(getTicketAction(999)).rejects.toThrow('HTTP 404');
        });
    });
});
