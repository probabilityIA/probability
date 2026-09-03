import { describe, it, expect, vi, beforeEach } from 'vitest';
import { TicketUseCases } from './use-cases';
import { ITicketRepository } from '../domain/ports';
import {
    Ticket,
    TicketComment,
    TicketAttachment,
    TicketHistoryEntry,
    PaginatedTickets,
    CreateTicketDTO,
    UpdateTicketDTO,
} from '../domain/types';

const makeTicket = (overrides: Partial<Ticket> = {}): Ticket => ({
    id: 1,
    code: 'TCK-001',
    business_id: 26,
    created_by_id: 5,
    title: 'No carga el listado',
    description: 'El listado queda en blanco',
    type: 'bug',
    priority: 'high',
    status: 'open',
    source: 'business',
    area: 'soporte',
    escalated_to_dev: false,
    created_at: '2026-09-01T10:00:00Z',
    updated_at: '2026-09-01T10:00:00Z',
    comments_count: 0,
    attachments_count: 0,
    ...overrides,
});

const makeComment = (overrides: Partial<TicketComment> = {}): TicketComment => ({
    id: 10,
    ticket_id: 1,
    user_id: 5,
    user_name: 'Cam',
    body: 'Revisando',
    is_internal: false,
    created_at: '2026-09-01T11:00:00Z',
    ...overrides,
});

const makeAttachment = (overrides: Partial<TicketAttachment> = {}): TicketAttachment => ({
    id: 20,
    ticket_id: 1,
    uploaded_by_id: 5,
    uploaded_by_name: 'Cam',
    file_url: 'https://s3/f.png',
    file_name: 'f.png',
    mime_type: 'image/png',
    size: 1024,
    created_at: '2026-09-01T12:00:00Z',
    ...overrides,
});

const makeHistory = (overrides: Partial<TicketHistoryEntry> = {}): TicketHistoryEntry => ({
    id: 30,
    ticket_id: 1,
    from_status: 'open',
    to_status: 'in_review',
    changed_by_id: 5,
    changed_by_name: 'Cam',
    note: '',
    created_at: '2026-09-01T13:00:00Z',
    ...overrides,
});

const paginated: PaginatedTickets = {
    data: [makeTicket()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const makeFile = () => new File(['x'], 'evidencia.png', { type: 'image/png' });

function createMockRepository(): ITicketRepository {
    return {
        list: vi.fn(),
        listCategories: vi.fn(),
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
}

describe('TicketUseCases', () => {
    let repo: ITicketRepository;
    let useCases: TicketUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new TicketUseCases(repo);
    });

    describe('list', () => {
        it('retorna la pagina del repositorio con los filtros recibidos', async () => {
            vi.mocked(repo.list).mockResolvedValue(paginated);

            const result = await useCases.list({ page: 2, page_size: 20, status: 'open', only_mine: true });

            expect(result).toEqual(paginated);
            expect(repo.list).toHaveBeenCalledWith({ page: 2, page_size: 20, status: 'open', only_mine: true });
        });

        it('delega sin parametros', async () => {
            vi.mocked(repo.list).mockResolvedValue(paginated);

            await useCases.list();

            expect(repo.list).toHaveBeenCalledWith(undefined);
        });

        it('propaga el error del repositorio', async () => {
            vi.mocked(repo.list).mockRejectedValue(new Error('HTTP 500'));

            await expect(useCases.list()).rejects.toThrow('HTTP 500');
        });
    });

    describe('get', () => {
        it('retorna el ticket por id', async () => {
            const ticket = makeTicket({ id: 7 });
            vi.mocked(repo.get).mockResolvedValue(ticket);

            const result = await useCases.get(7);

            expect(result).toBe(ticket);
            expect(repo.get).toHaveBeenCalledWith(7, undefined);
        });

        it('propaga el business id del super admin', async () => {
            vi.mocked(repo.get).mockResolvedValue(makeTicket());

            await useCases.get(7, 26);

            expect(repo.get).toHaveBeenCalledWith(7, 26);
        });

        it('propaga el error cuando el ticket no existe', async () => {
            vi.mocked(repo.get).mockRejectedValue(new Error('ticket no encontrado'));

            await expect(useCases.get(999)).rejects.toThrow('ticket no encontrado');
        });
    });

    describe('create', () => {
        const dto: CreateTicketDTO = {
            business_id: 26,
            title: 'Falla el envio',
            description: 'No genera guia',
            type: 'bug',
            priority: 'critical',
            area: 'soporte',
        };

        it('crea el ticket y retorna el creado', async () => {
            const created = makeTicket({ id: 2, title: 'Falla el envio' });
            vi.mocked(repo.create).mockResolvedValue(created);

            const result = await useCases.create(dto);

            expect(result).toBe(created);
            expect(repo.create).toHaveBeenCalledWith(dto);
        });

        it('propaga el error de validacion', async () => {
            vi.mocked(repo.create).mockRejectedValue(new Error('titulo requerido'));

            await expect(useCases.create(dto)).rejects.toThrow('titulo requerido');
        });
    });

    describe('update', () => {
        const dto: UpdateTicketDTO = { title: 'Nuevo titulo', clear_due_date: true };

        it('actualiza el ticket con el id y el dto', async () => {
            const updated = makeTicket({ title: 'Nuevo titulo' });
            vi.mocked(repo.update).mockResolvedValue(updated);

            const result = await useCases.update(1, dto);

            expect(result).toBe(updated);
            expect(repo.update).toHaveBeenCalledWith(1, dto);
        });

        it('propaga el error', async () => {
            vi.mocked(repo.update).mockRejectedValue(new Error('no autorizado'));

            await expect(useCases.update(1, dto)).rejects.toThrow('no autorizado');
        });
    });

    describe('remove', () => {
        it('elimina el ticket por id', async () => {
            vi.mocked(repo.remove).mockResolvedValue(undefined);

            await expect(useCases.remove(3)).resolves.toBeUndefined();
            expect(repo.remove).toHaveBeenCalledWith(3);
        });

        it('propaga el error', async () => {
            vi.mocked(repo.remove).mockRejectedValue(new Error('HTTP 403'));

            await expect(useCases.remove(3)).rejects.toThrow('HTTP 403');
        });
    });

    describe('changeStatus', () => {
        it('delega estado y nota', async () => {
            const moved = makeTicket({ status: 'in_review' });
            vi.mocked(repo.changeStatus).mockResolvedValue(moved);

            const result = await useCases.changeStatus(1, 'in_review', 'lo tomo yo');

            expect(result).toBe(moved);
            expect(repo.changeStatus).toHaveBeenCalledWith(1, 'in_review', 'lo tomo yo');
        });

        it('delega sin nota', async () => {
            vi.mocked(repo.changeStatus).mockResolvedValue(makeTicket({ status: 'closed' }));

            await useCases.changeStatus(1, 'closed');

            expect(repo.changeStatus).toHaveBeenCalledWith(1, 'closed', undefined);
        });
    });

    describe('assign', () => {
        it('asigna a un usuario', async () => {
            const assigned = makeTicket({ assigned_to_id: 9 });
            vi.mocked(repo.assign).mockResolvedValue(assigned);

            const result = await useCases.assign(1, 9);

            expect(result).toBe(assigned);
            expect(repo.assign).toHaveBeenCalledWith(1, 9);
        });

        it('desasigna con null', async () => {
            vi.mocked(repo.assign).mockResolvedValue(makeTicket({ assigned_to_id: null }));

            await useCases.assign(1, null);

            expect(repo.assign).toHaveBeenCalledWith(1, null);
        });
    });

    describe('changeArea', () => {
        it('delega area y nota', async () => {
            const moved = makeTicket({ area: 'desarrollo' });
            vi.mocked(repo.changeArea).mockResolvedValue(moved);

            const result = await useCases.changeArea(1, 'desarrollo', 'necesita codigo');

            expect(result).toBe(moved);
            expect(repo.changeArea).toHaveBeenCalledWith(1, 'desarrollo', 'necesita codigo');
        });

        it('delega sin nota', async () => {
            vi.mocked(repo.changeArea).mockResolvedValue(makeTicket({ area: 'comercial' }));

            await useCases.changeArea(1, 'comercial');

            expect(repo.changeArea).toHaveBeenCalledWith(1, 'comercial', undefined);
        });
    });

    describe('escalate', () => {
        it('escala con nota', async () => {
            const escalated = makeTicket({ escalated_to_dev: true });
            vi.mocked(repo.escalate).mockResolvedValue(escalated);

            const result = await useCases.escalate(1, 'urgente');

            expect(result).toBe(escalated);
            expect(repo.escalate).toHaveBeenCalledWith(1, 'urgente');
        });

        it('escala sin nota', async () => {
            vi.mocked(repo.escalate).mockResolvedValue(makeTicket({ escalated_to_dev: true }));

            await useCases.escalate(1);

            expect(repo.escalate).toHaveBeenCalledWith(1, undefined);
        });
    });

    describe('changeSprint', () => {
        it('asigna un sprint', async () => {
            const moved = makeTicket({ sprint_id: 4 });
            vi.mocked(repo.changeSprint).mockResolvedValue(moved);

            const result = await useCases.changeSprint(1, 4);

            expect(result).toBe(moved);
            expect(repo.changeSprint).toHaveBeenCalledWith(1, 4);
        });

        it('saca del sprint con null', async () => {
            vi.mocked(repo.changeSprint).mockResolvedValue(makeTicket({ sprint_id: null }));

            await useCases.changeSprint(1, null);

            expect(repo.changeSprint).toHaveBeenCalledWith(1, null);
        });
    });

    describe('listComments', () => {
        it('retorna los comentarios del ticket', async () => {
            const comments = [makeComment()];
            vi.mocked(repo.listComments).mockResolvedValue(comments);

            const result = await useCases.listComments(1, 26);

            expect(result).toBe(comments);
            expect(repo.listComments).toHaveBeenCalledWith(1, 26);
        });

        it('delega sin business id', async () => {
            vi.mocked(repo.listComments).mockResolvedValue([]);

            await useCases.listComments(1);

            expect(repo.listComments).toHaveBeenCalledWith(1, undefined);
        });
    });

    describe('addComment', () => {
        it('agrega un comentario publico', async () => {
            const comment = makeComment({ body: 'Ya quedo' });
            vi.mocked(repo.addComment).mockResolvedValue(comment);

            const result = await useCases.addComment(1, 'Ya quedo', false);

            expect(result).toBe(comment);
            expect(repo.addComment).toHaveBeenCalledWith(1, 'Ya quedo', false);
        });

        it('agrega un comentario interno', async () => {
            vi.mocked(repo.addComment).mockResolvedValue(makeComment({ is_internal: true }));

            await useCases.addComment(1, 'nota interna', true);

            expect(repo.addComment).toHaveBeenCalledWith(1, 'nota interna', true);
        });
    });

    describe('listAttachments', () => {
        it('retorna los adjuntos del ticket', async () => {
            const attachments = [makeAttachment()];
            vi.mocked(repo.listAttachments).mockResolvedValue(attachments);

            const result = await useCases.listAttachments(1, 26);

            expect(result).toBe(attachments);
            expect(repo.listAttachments).toHaveBeenCalledWith(1, 26);
        });
    });

    describe('uploadAttachment', () => {
        it('sube el archivo asociado al ticket', async () => {
            const file = makeFile();
            const attachment = makeAttachment();
            vi.mocked(repo.uploadAttachment).mockResolvedValue(attachment);

            const result = await useCases.uploadAttachment(1, file);

            expect(result).toBe(attachment);
            expect(repo.uploadAttachment).toHaveBeenCalledWith(1, file, undefined);
        });

        it('sube el archivo asociado a un comentario', async () => {
            const file = makeFile();
            vi.mocked(repo.uploadAttachment).mockResolvedValue(makeAttachment({ comment_id: 10 }));

            await useCases.uploadAttachment(1, file, 10);

            expect(repo.uploadAttachment).toHaveBeenCalledWith(1, file, 10);
        });

        it('propaga el error de subida', async () => {
            vi.mocked(repo.uploadAttachment).mockRejectedValue(new Error('archivo muy grande'));

            await expect(useCases.uploadAttachment(1, makeFile())).rejects.toThrow('archivo muy grande');
        });
    });

    describe('deleteAttachment', () => {
        it('elimina el adjunto por su propio id', async () => {
            vi.mocked(repo.deleteAttachment).mockResolvedValue(undefined);

            await expect(useCases.deleteAttachment(20)).resolves.toBeUndefined();
            expect(repo.deleteAttachment).toHaveBeenCalledWith(20);
        });
    });

    describe('listHistory', () => {
        it('retorna el historial del ticket', async () => {
            const history = [makeHistory()];
            vi.mocked(repo.listHistory).mockResolvedValue(history);

            const result = await useCases.listHistory(1, 26);

            expect(result).toBe(history);
            expect(repo.listHistory).toHaveBeenCalledWith(1, 26);
        });

        it('delega sin business id', async () => {
            vi.mocked(repo.listHistory).mockResolvedValue([]);

            await useCases.listHistory(1);

            expect(repo.listHistory).toHaveBeenCalledWith(1, undefined);
        });
    });
});
