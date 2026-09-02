import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { env } from '@/shared/config/env';
import { TicketApiRepository } from './api-repository';
import { Ticket } from '../../domain/types';

const BASE = env.API_BASE_URL;

const fetchMock = vi.fn();

const ok = (body: unknown, status = 200) => ({
    ok: true,
    status,
    text: async () => (body === undefined ? '' : JSON.stringify(body)),
});

const fail = (body: unknown, status = 500) => ({
    ok: false,
    status,
    text: async () => (body === undefined ? '' : JSON.stringify(body)),
});

const ticket: Ticket = {
    id: 5,
    code: 'TCK-005',
    business_id: 26,
    created_by_id: 1,
    title: 'Titulo',
    description: 'Descripcion',
    type: 'bug',
    priority: 'high',
    status: 'open',
    source: 'business',
    escalated_to_dev: false,
    created_at: '2026-09-01T10:00:00Z',
    updated_at: '2026-09-01T10:00:00Z',
    comments_count: 0,
    attachments_count: 0,
};

const lastCall = () => {
    const call = fetchMock.mock.calls[fetchMock.mock.calls.length - 1];
    return { url: call[0] as string, init: call[1] as RequestInit };
};

const headersOf = () => lastCall().init.headers as Record<string, string>;
const bodyOf = () => JSON.parse(lastCall().init.body as string);

describe('TicketApiRepository', () => {
    let repo: TicketApiRepository;

    beforeEach(() => {
        fetchMock.mockReset();
        vi.stubGlobal('fetch', fetchMock);
        repo = new TicketApiRepository('tok-123');
    });

    afterEach(() => {
        vi.unstubAllGlobals();
    });

    describe('headers y opciones comunes', () => {
        it('manda Authorization con el token recibido', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.get(5);

            expect(headersOf()).toMatchObject({
                Accept: 'application/json',
                Authorization: 'Bearer tok-123',
            });
        });

        it('omite Authorization cuando no hay token', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await new TicketApiRepository().get(5);

            expect(headersOf()).toEqual({ Accept: 'application/json' });
        });

        it('omite Authorization cuando el token es cadena vacia', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await new TicketApiRepository('').get(5);

            expect(headersOf().Authorization).toBeUndefined();
        });

        it('desactiva la cache de Next en toda peticion', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.get(5);

            expect(lastCall().init.cache).toBe('no-store');
        });
    });

    describe('list', () => {
        const page = { data: [ticket], total: 1, page: 1, page_size: 10, total_pages: 1 };

        it('pide /tickets sin query cuando no hay parametros', async () => {
            fetchMock.mockResolvedValue(ok(page));

            const result = await repo.list();

            expect(result).toEqual(page);
            expect(lastCall().url).toBe(`${BASE}/tickets`);
            expect(lastCall().init.method).toBe('GET');
        });

        it('serializa los filtros recibidos en el query string', async () => {
            fetchMock.mockResolvedValue(ok(page));

            await repo.list({ page: 2, page_size: 50, status: 'open', escalated: true, sprint_id: 4 });

            expect(lastCall().url).toBe(
                `${BASE}/tickets?page=2&page_size=50&status=open&escalated=true&sprint_id=4`
            );
        });

        it('descarta parametros undefined, null y vacios', async () => {
            fetchMock.mockResolvedValue(ok(page));

            await repo.list({
                page: 1,
                status: '',
                priority: undefined,
                area: null as any,
                search: 'guia',
            });

            expect(lastCall().url).toBe(`${BASE}/tickets?page=1&search=guia`);
        });

        it('conserva el path sin ? cuando todos los parametros se descartan', async () => {
            fetchMock.mockResolvedValue(ok(page));

            await repo.list({ status: '', search: undefined });

            expect(lastCall().url).toBe(`${BASE}/tickets`);
        });
    });

    describe('get', () => {
        it('pide el ticket por id sin query cuando no hay business', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            const result = await repo.get(5);

            expect(result).toEqual(ticket);
            expect(lastCall().url).toBe(`${BASE}/tickets/5`);
        });

        it('agrega business_id para el super admin', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.get(5, 26);

            expect(lastCall().url).toBe(`${BASE}/tickets/5?business_id=26`);
        });
    });

    describe('create', () => {
        it('hace POST a /tickets con el dto serializado', async () => {
            fetchMock.mockResolvedValue(ok(ticket, 201));

            const result = await repo.create({ title: 'Titulo', description: 'Desc', type: 'bug' });

            expect(result).toEqual(ticket);
            expect(lastCall().url).toBe(`${BASE}/tickets`);
            expect(lastCall().init.method).toBe('POST');
            expect(headersOf()['Content-Type']).toBe('application/json');
            expect(bodyOf()).toEqual({ title: 'Titulo', description: 'Desc', type: 'bug' });
        });
    });

    describe('update', () => {
        it('hace PUT al ticket con el dto serializado', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.update(5, { title: 'Otro', clear_due_date: true });

            expect(lastCall().url).toBe(`${BASE}/tickets/5`);
            expect(lastCall().init.method).toBe('PUT');
            expect(headersOf()['Content-Type']).toBe('application/json');
            expect(bodyOf()).toEqual({ title: 'Otro', clear_due_date: true });
        });
    });

    describe('remove', () => {
        it('hace DELETE al ticket y no retorna nada', async () => {
            fetchMock.mockResolvedValue(ok(undefined, 204));

            const result = await repo.remove(5);

            expect(result).toBeUndefined();
            expect(lastCall().url).toBe(`${BASE}/tickets/5`);
            expect(lastCall().init.method).toBe('DELETE');
            expect(lastCall().init.body).toBeUndefined();
        });
    });

    describe('changeStatus', () => {
        it('hace PATCH a /status con estado y nota', async () => {
            fetchMock.mockResolvedValue(ok({ ...ticket, status: 'in_review' }));

            const result = await repo.changeStatus(5, 'in_review', 'lo reviso');

            expect(result.status).toBe('in_review');
            expect(lastCall().url).toBe(`${BASE}/tickets/5/status`);
            expect(lastCall().init.method).toBe('PATCH');
            expect(bodyOf()).toEqual({ status: 'in_review', note: 'lo reviso' });
        });

        it('manda nota vacia cuando no se envia', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.changeStatus(5, 'closed');

            expect(bodyOf()).toEqual({ status: 'closed', note: '' });
        });
    });

    describe('assign', () => {
        it('hace PATCH a /assign con el usuario', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.assign(5, 9);

            expect(lastCall().url).toBe(`${BASE}/tickets/5/assign`);
            expect(lastCall().init.method).toBe('PATCH');
            expect(bodyOf()).toEqual({ assigned_to_id: 9 });
        });

        it('manda null para desasignar', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.assign(5, null);

            expect(bodyOf()).toEqual({ assigned_to_id: null });
        });
    });

    describe('changeArea', () => {
        it('usa la ruta /area sin tilde ni acentos', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.changeArea(5, 'desarrollo');

            expect(lastCall().url).toBe(`${BASE}/tickets/5/area`);
            expect(lastCall().url).not.toContain('%C3');
            expect(lastCall().init.method).toBe('PATCH');
        });

        it('manda area y nota', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.changeArea(5, 'comercial', 'lo pasa ventas');

            expect(bodyOf()).toEqual({ area: 'comercial', note: 'lo pasa ventas' });
        });

        it('manda nota vacia cuando no se envia', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.changeArea(5, 'soporte');

            expect(bodyOf()).toEqual({ area: 'soporte', note: '' });
        });
    });

    describe('changeSprint', () => {
        it('hace PATCH a /sprint con el sprint', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.changeSprint(5, 4);

            expect(lastCall().url).toBe(`${BASE}/tickets/5/sprint`);
            expect(lastCall().init.method).toBe('PATCH');
            expect(bodyOf()).toEqual({ sprint_id: 4 });
        });

        it('manda null para sacar del sprint', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.changeSprint(5, null);

            expect(bodyOf()).toEqual({ sprint_id: null });
        });
    });

    describe('escalate', () => {
        it('hace PATCH a /escalate con la nota', async () => {
            fetchMock.mockResolvedValue(ok({ ...ticket, escalated_to_dev: true }));

            const result = await repo.escalate(5, 'necesita dev');

            expect(result.escalated_to_dev).toBe(true);
            expect(lastCall().url).toBe(`${BASE}/tickets/5/escalate`);
            expect(lastCall().init.method).toBe('PATCH');
            expect(bodyOf()).toEqual({ note: 'necesita dev' });
        });

        it('manda nota vacia cuando no se envia', async () => {
            fetchMock.mockResolvedValue(ok(ticket));

            await repo.escalate(5);

            expect(bodyOf()).toEqual({ note: '' });
        });
    });

    describe('listComments', () => {
        const comment = {
            id: 10,
            ticket_id: 5,
            user_id: 1,
            user_name: 'Cam',
            body: 'Hola',
            is_internal: false,
            created_at: '2026-09-01T11:00:00Z',
        };

        it('desempaqueta el campo data de la respuesta', async () => {
            fetchMock.mockResolvedValue(ok({ data: [comment] }));

            const result = await repo.listComments(5);

            expect(result).toEqual([comment]);
            expect(lastCall().url).toBe(`${BASE}/tickets/5/comments`);
            expect(lastCall().init.method).toBe('GET');
        });

        it('agrega business_id cuando se envia', async () => {
            fetchMock.mockResolvedValue(ok({ data: [] }));

            await repo.listComments(5, 26);

            expect(lastCall().url).toBe(`${BASE}/tickets/5/comments?business_id=26`);
        });

        it('retorna arreglo vacio si la respuesta no trae data', async () => {
            fetchMock.mockResolvedValue(ok({ total: 0 }));

            await expect(repo.listComments(5)).resolves.toEqual([]);
        });

        it('retorna arreglo vacio si data viene en null', async () => {
            fetchMock.mockResolvedValue(ok({ data: null }));

            await expect(repo.listComments(5)).resolves.toEqual([]);
        });
    });

    describe('addComment', () => {
        it('hace POST del comentario publico', async () => {
            fetchMock.mockResolvedValue(ok({ id: 11 }));

            await repo.addComment(5, 'Ya quedo', false);

            expect(lastCall().url).toBe(`${BASE}/tickets/5/comments`);
            expect(lastCall().init.method).toBe('POST');
            expect(headersOf()['Content-Type']).toBe('application/json');
            expect(bodyOf()).toEqual({ body: 'Ya quedo', is_internal: false });
        });

        it('marca is_internal en el comentario interno', async () => {
            fetchMock.mockResolvedValue(ok({ id: 12 }));

            await repo.addComment(5, 'nota', true);

            expect(bodyOf()).toEqual({ body: 'nota', is_internal: true });
        });
    });

    describe('listAttachments', () => {
        it('desempaqueta el campo data', async () => {
            const attachment = { id: 20, file_name: 'f.png' };
            fetchMock.mockResolvedValue(ok({ data: [attachment] }));

            const result = await repo.listAttachments(5);

            expect(result).toEqual([attachment]);
            expect(lastCall().url).toBe(`${BASE}/tickets/5/attachments`);
        });

        it('agrega business_id cuando se envia', async () => {
            fetchMock.mockResolvedValue(ok({ data: [] }));

            await repo.listAttachments(5, 26);

            expect(lastCall().url).toBe(`${BASE}/tickets/5/attachments?business_id=26`);
        });

        it('retorna arreglo vacio si no viene data', async () => {
            fetchMock.mockResolvedValue(ok({}));

            await expect(repo.listAttachments(5)).resolves.toEqual([]);
        });
    });

    describe('uploadAttachment', () => {
        const file = () => new File(['x'], 'evidencia.png', { type: 'image/png' });

        it('hace POST multipart con el archivo', async () => {
            fetchMock.mockResolvedValue(ok({ id: 21 }));
            const f = file();

            await repo.uploadAttachment(5, f);

            const fd = lastCall().init.body as FormData;
            expect(lastCall().url).toBe(`${BASE}/tickets/5/attachments`);
            expect(lastCall().init.method).toBe('POST');
            expect(fd).toBeInstanceOf(FormData);
            expect(fd.get('file')).toBe(f);
            expect(fd.get('comment_id')).toBeNull();
        });

        it('no fija Content-Type para que el navegador ponga el boundary', async () => {
            fetchMock.mockResolvedValue(ok({ id: 21 }));

            await repo.uploadAttachment(5, file());

            expect(headersOf()['Content-Type']).toBeUndefined();
        });

        it('adjunta comment_id cuando se envia', async () => {
            fetchMock.mockResolvedValue(ok({ id: 22 }));

            await repo.uploadAttachment(5, file(), 10);

            expect((lastCall().init.body as FormData).get('comment_id')).toBe('10');
        });
    });

    describe('deleteAttachment', () => {
        it('hace DELETE a la ruta global de adjuntos', async () => {
            fetchMock.mockResolvedValue(ok(undefined, 204));

            await expect(repo.deleteAttachment(20)).resolves.toBeUndefined();
            expect(lastCall().url).toBe(`${BASE}/tickets/attachments/20`);
            expect(lastCall().init.method).toBe('DELETE');
        });
    });

    describe('listHistory', () => {
        it('desempaqueta el campo data', async () => {
            const entry = { id: 30, from_status: 'open', to_status: 'in_review' };
            fetchMock.mockResolvedValue(ok({ data: [entry] }));

            const result = await repo.listHistory(5);

            expect(result).toEqual([entry]);
            expect(lastCall().url).toBe(`${BASE}/tickets/5/history`);
        });

        it('agrega business_id cuando se envia', async () => {
            fetchMock.mockResolvedValue(ok({ data: [] }));

            await repo.listHistory(5, 26);

            expect(lastCall().url).toBe(`${BASE}/tickets/5/history?business_id=26`);
        });

        it('retorna arreglo vacio si no viene data', async () => {
            fetchMock.mockResolvedValue(ok({}));

            await expect(repo.listHistory(5)).resolves.toEqual([]);
        });
    });

    describe('manejo de errores', () => {
        it('lanza el campo error del payload', async () => {
            fetchMock.mockResolvedValue(fail({ error: 'ticket no encontrado' }, 404));

            await expect(repo.get(999)).rejects.toThrow('ticket no encontrado');
        });

        it('cae al campo message cuando no hay error', async () => {
            fetchMock.mockResolvedValue(fail({ message: 'no autorizado' }, 403));

            await expect(repo.list()).rejects.toThrow('no autorizado');
        });

        it('cae al status HTTP cuando el payload no trae mensaje', async () => {
            fetchMock.mockResolvedValue(fail({ data: null }, 500));

            await expect(repo.list()).rejects.toThrow('HTTP 500');
        });

        it('cae al status HTTP cuando el cuerpo viene vacio', async () => {
            fetchMock.mockResolvedValue(fail(undefined, 502));

            await expect(repo.remove(5)).rejects.toThrow('HTTP 502');
        });

        it('propaga la caida de red', async () => {
            fetchMock.mockRejectedValue(new TypeError('Failed to fetch'));

            await expect(repo.list()).rejects.toThrow('Failed to fetch');
        });

        it('propaga el error de las rutas que desempaquetan data', async () => {
            fetchMock.mockResolvedValue(fail({ error: 'sin permisos' }, 403));

            await expect(repo.listComments(5)).rejects.toThrow('sin permisos');
        });
    });
});
