import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor, within } from '@testing-library/react';
import TicketsManager from './TicketsManager';
import { PaginatedTickets, Ticket } from '../../domain/types';

vi.mock('@/shared/ui', () => ({
    TICKETS_TABS_SLOT_ID: 'tickets-tabs-slot',
    TICKETS_ACTIONS_SLOT_ID: 'tickets-actions-slot',
    TICKETS_FILTERS_SLOT_ID: 'tickets-filters-slot',
    Modal: ({ isOpen, title, onClose, children }: any) =>
        isOpen ? (
            <div data-testid="modal" data-title={String(title)}>
                <button onClick={onClose}>modal-close</button>
                {children}
            </div>
        ) : null,
    Button: ({ children, ...props }: any) => <button {...props}>{children}</button>,
    Input: (props: any) => <input {...props} />,
    DynamicFilters: ({ availableFilters, activeFilters, onAddFilter, onRemoveFilter, onSortChange, onCreate, createButtonAriaLabel }: any) => (
        <div data-testid="filters">
            <button aria-label={createButtonAriaLabel} onClick={() => onCreate()} />
            <span data-testid="available-filters">{availableFilters.map((f: any) => f.key).join(',')}</span>
            <span data-testid="active-filters">{activeFilters.map((f: any) => f.key + '=' + f.value).join('|')}</span>
            <button onClick={() => onAddFilter('search', 'pago')}>add-search</button>
            <button onClick={() => onAddFilter('status', 'open')}>add-status</button>
            <button onClick={() => onAddFilter('area', 'soporte')}>add-area</button>
            <button onClick={() => onAddFilter('priority', 'high')}>add-priority</button>
            <button onClick={() => onAddFilter('type', 'bug')}>add-type</button>
            <button onClick={() => onAddFilter('source', 'internal')}>add-source</button>
            <button onClick={() => onAddFilter('only_mine', 'true')}>add-only-mine</button>
            <button onClick={() => onAddFilter('only_mine', true)}>add-only-mine-bool</button>
            <button onClick={() => onAddFilter('status', 'zzz')}>add-status-desconocido</button>
            <button onClick={() => onAddFilter('area', 'zzz')}>add-area-desconocida</button>
            <button onClick={() => onAddFilter('priority', 'zzz')}>add-priority-desconocida</button>
            <button onClick={() => onAddFilter('type', 'zzz')}>add-type-desconocido</button>
            <button onClick={() => onAddFilter('source', 'business')}>add-source-negocio</button>
            <button onClick={() => onAddFilter('escalated', 'true')}>add-escalated</button>
            <button onClick={() => onAddFilter('assigned_to_id', '5')}>add-assigned</button>
            <button onClick={() => onAddFilter('assigned_to_id', '')}>add-assigned-empty</button>
            <button onClick={() => onRemoveFilter('status')}>remove-status</button>
            <button onClick={() => onSortChange('priority', 'asc')}>sort-priority</button>
        </div>
    ),
    TablePagination: ({ currentPage, totalPages, totalItems, pageSize, onPageChange, onPageSizeChange }: any) => (
        <div data-testid="pagination">
            <span data-testid="pagination-state">{[currentPage, totalPages, totalItems, pageSize].join('/')}</span>
            <button onClick={() => onPageChange(3)}>go-page-3</button>
            <button onClick={() => onPageSizeChange(50)}>set-size-50</button>
        </div>
    ),
}));

vi.mock('@/shared/contexts/permissions-context', () => ({
    usePermissions: vi.fn(() => ({ isSuperAdmin: true })),
}));

vi.mock('../../infra/actions', () => ({
    listTicketsAction: vi.fn(),
    createTicketAction: vi.fn(),
    getTicketAction: vi.fn(),
    changeTicketStatusAction: vi.fn(),
    changeTicketAreaAction: vi.fn(),
    assignTicketAction: vi.fn(),
    uploadAttachmentAction: vi.fn(),
    addCommentAction: vi.fn(),
    changeTicketSprintAction: vi.fn(),
    listCommentsAction: vi.fn(),
    listAttachmentsAction: vi.fn(),
    listTicketHistoryAction: vi.fn(),
    deleteAttachmentAction: vi.fn(),
    escalateTicketAction: vi.fn(),
    deleteTicketAction: vi.fn(),
}));

vi.mock('@/services/modules/sprints/infra/actions', () => ({
    listSprintsAction: vi.fn(),
    createSprintAction: vi.fn(),
    updateSprintAction: vi.fn(),
    changeSprintStatusAction: vi.fn(),
}));

vi.mock('@/services/auth/users/infra/actions', () => ({
    getUsersAction: vi.fn(),
}));

vi.mock('@/services/auth/business/ui/hooks/useBusinessesSimple', () => ({
    useBusinessesSimple: vi.fn(() => ({ businesses: [], loading: false, error: null })),
}));

import {
    listTicketsAction,
    createTicketAction,
    getTicketAction,
    changeTicketStatusAction,
    changeTicketAreaAction,
    assignTicketAction,
    uploadAttachmentAction,
    addCommentAction,
    listCommentsAction,
    listAttachmentsAction,
    listTicketHistoryAction,
} from '../../infra/actions';
import { listSprintsAction } from '@/services/modules/sprints/infra/actions';
import { getUsersAction } from '@/services/auth/users/infra/actions';
import { usePermissions } from '@/shared/contexts/permissions-context';

const makeTicket = (over: Partial<Ticket> & { id: number }): Ticket => ({
    code: 'TCK-' + over.id,
    created_by_id: 1,
    title: 'Ticket ' + over.id,
    description: 'desc',
    type: 'bug',
    priority: 'medium',
    status: 'open',
    source: 'internal',
    escalated_to_dev: false,
    created_at: '2026-01-01T10:00:00Z',
    updated_at: '2026-01-02T10:00:00Z',
    comments_count: 2,
    attachments_count: 1,
    ...over,
});

const paginated = (data: Ticket[], over: Partial<PaginatedTickets> = {}): PaginatedTickets => ({
    data,
    total: data.length,
    page: 1,
    page_size: 10,
    total_pages: 1,
    ...over,
});

const sprint = (id: number, name: string, status: 'planned' | 'active' | 'closed') => ({
    id,
    name,
    goal: '',
    start_date: '2026-02-01T12:00:00',
    end_date: '2026-02-15T12:00:00',
    status,
    created_by_id: 1,
    created_by_name: 'admin',
    ticket_count: 0,
    done_count: 0,
    created_at: '2026-01-01T12:00:00',
    updated_at: '2026-01-01T12:00:00',
});

const makeDataTransfer = () => {
    const store: Record<string, string> = {};
    return {
        effectAllowed: '',
        dropEffect: '',
        setData: (key: string, value: string) => { store[key] = value; },
        getData: (key: string) => store[key] ?? '',
    };
};

const makeFile = (name: string, type = 'image/png') => new File(['x'], name, { type });

const rowOf = (code: string): HTMLElement => screen.getByText(code).closest('tr') as HTMLElement;

const boardColumn = (title: string): HTMLElement =>
    screen.getByRole('heading', { name: title }).parentElement!.parentElement as HTMLElement;

const boardCard = (code: string): HTMLElement => screen.getByText(code).closest('[draggable]') as HTMLElement;

const lastListCall = () => {
    const calls = vi.mocked(listTicketsAction).mock.calls;
    return calls[calls.length - 1][0] as any;
};

const addSlots = () => {
    ['tickets-tabs-slot', 'tickets-actions-slot', 'tickets-filters-slot'].forEach((id) => {
        const el = document.createElement('div');
        el.id = id;
        document.body.appendChild(el);
    });
};

const renderManager = async () => {
    const utils = render(<TicketsManager />);
    await waitFor(() => expect(listTicketsAction).toHaveBeenCalled());
    return utils;
};

const openCreateModal = () => fireEvent.click(screen.getByRole('button', { name: 'Nuevo ticket' }));

const fillCreateForm = (title = 'Ticket nuevo', description = 'Descripcion larga') => {
    fireEvent.change(screen.getByPlaceholderText('Resumen breve'), { target: { value: title } });
    fireEvent.change(screen.getByPlaceholderText('Detalla el problema, mejora o solicitud'), { target: { value: description } });
};

const attachToCreateForm = (files: File[]) => {
    const input = document.getElementById('ticket-create-file-input') as HTMLInputElement;
    Object.defineProperty(input, 'files', { value: files, configurable: true });
    fireEvent.change(input);
};

const submitCreateForm = () => fireEvent.click(screen.getByRole('button', { name: 'Crear ticket' }));

describe('TicketsManager', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        window.localStorage.clear();
        window.history.replaceState(null, '', '/tickets');
        document.querySelectorAll('[id^="tickets-"]').forEach((n) => n.remove());
        vi.mocked(usePermissions).mockReturnValue({ isSuperAdmin: true } as any);
        vi.mocked(listTicketsAction).mockResolvedValue(paginated([makeTicket({ id: 1 })]) as any);
        vi.mocked(listSprintsAction).mockResolvedValue({ data: [], total: 0, page: 1, page_size: 100, total_pages: 0 } as any);
        vi.mocked(getUsersAction).mockResolvedValue({ data: [] } as any);
        vi.mocked(listCommentsAction).mockResolvedValue([]);
        vi.mocked(listAttachmentsAction).mockResolvedValue([]);
        vi.mocked(listTicketHistoryAction).mockResolvedValue([]);
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('listado en tabla', () => {
        it('pide la primera pagina con el orden por defecto', async () => {
            await renderManager();

            expect(lastListCall()).toMatchObject({
                page: 1,
                page_size: 10,
                sort_by: 'created_at',
                sort_order: 'desc',
                sprint_id: undefined,
            });
        });

        it('pinta las columnas y los datos de cada ticket', async () => {
            vi.mocked(listTicketsAction).mockResolvedValue(paginated([
                makeTicket({ id: 1, business_name: 'Tienda Uno' }),
                makeTicket({ id: 2, business_id: 7, business_name: undefined }),
                makeTicket({ id: 3, business_name: undefined }),
            ]) as any);
            await renderManager();

            expect(await screen.findByText('TCK-1')).toBeInTheDocument();
            expect(within(rowOf('TCK-1')).getByText('Ticket 1')).toBeInTheDocument();
            expect(within(rowOf('TCK-1')).getByText('2 comentarios | 1 adjuntos')).toBeInTheDocument();
            expect(within(rowOf('TCK-1')).getByText('Tienda Uno')).toBeInTheDocument();
            expect(within(rowOf('TCK-2')).getByText('#7')).toBeInTheDocument();
            expect(within(rowOf('TCK-3')).getByText('Interno')).toBeInTheDocument();
        });

        it('avisa cuando no hay tickets', async () => {
            vi.mocked(listTicketsAction).mockResolvedValue(paginated([]) as any);
            await renderManager();

            expect(await screen.findByText('No hay tickets disponibles')).toBeInTheDocument();
        });

        it('no rompe si el listado falla', async () => {
            vi.mocked(listTicketsAction).mockRejectedValue(new Error('500'));
            await renderManager();

            await waitFor(() => expect(screen.queryByText('Actualizando...')).toBeNull());
            expect(screen.queryByTestId('pagination')).toBeNull();
        });

        it('usa valores por defecto cuando el servidor no envia pagina ni total de paginas', async () => {
            vi.mocked(listTicketsAction).mockResolvedValue(
                paginated([makeTicket({ id: 1 })], { page: 0, total_pages: 0, total: 0 }) as any
            );
            await renderManager();

            expect(await screen.findByTestId('pagination-state')).toHaveTextContent('1/1/0/10');
        });

        it('cambia de pagina y de tamano de pagina', async () => {
            await renderManager();

            fireEvent.click(await screen.findByText('go-page-3'));
            await waitFor(() => expect(lastListCall()).toMatchObject({ page: 3 }));

            fireEvent.click(screen.getByText('set-size-50'));
            await waitFor(() => expect(lastListCall()).toMatchObject({ page: 1, page_size: 50 }));
        });
    });

    describe('cabecera portalizada', () => {
        it('monta pestanas, acciones y filtros en los slots del layout', async () => {
            addSlots();
            await renderManager();

            const tabs = document.getElementById('tickets-tabs-slot') as HTMLElement;
            const actions = document.getElementById('tickets-actions-slot') as HTMLElement;
            const filters = document.getElementById('tickets-filters-slot') as HTMLElement;
            await waitFor(() => expect(within(tabs).getByRole('button', { name: 'Tabla' })).toBeInTheDocument());
            expect(within(filters).getByTestId('filters')).toBeInTheDocument();
            expect(within(filters).getByRole('button', { name: 'Nuevo ticket' })).toBeInTheDocument();
            expect(screen.getAllByRole('button', { name: 'Tabla' })).toHaveLength(1);
        });

        it('dibuja la cabecera en linea cuando no existen los slots', async () => {
            await renderManager();

            expect(document.getElementById('tickets-tabs-slot')).toBeNull();
            expect(screen.getByRole('button', { name: 'Tabla' })).toBeInTheDocument();
            expect(screen.getByRole('button', { name: 'Nuevo ticket' })).toBeInTheDocument();
            expect(screen.getByTestId('filters')).toBeInTheDocument();
        });
    });

    describe('cambio de vista', () => {
        it('guarda la vista elegida en localStorage', async () => {
            await renderManager();

            fireEvent.click(screen.getByRole('button', { name: 'Board' }));

            await waitFor(() => expect(window.localStorage.getItem('tickets_view_mode')).toBe('board'));
        });

        it('restaura la vista guardada al montar', async () => {
            window.localStorage.setItem('tickets_view_mode', 'backlog');
            await renderManager();

            expect(await screen.findByRole('heading', { name: 'Backlog' })).toBeInTheDocument();
        });

        it('ignora un valor invalido guardado en localStorage', async () => {
            window.localStorage.setItem('tickets_view_mode', 'galaxia');
            await renderManager();

            expect(await screen.findByText('TCK-1')).toBeInTheDocument();
        });

        it('en board pide hasta 100 tickets en una sola pagina', async () => {
            await renderManager();

            fireEvent.click(screen.getByRole('button', { name: 'Board' }));

            await waitFor(() => expect(lastListCall()).toMatchObject({ page: 1, page_size: 100 }));
        });

        it('en backlog no vuelve a pedir el listado principal', async () => {
            await renderManager();
            const mainCalls = () =>
                vi.mocked(listTicketsAction).mock.calls.filter((c: any) => c[0]?.page_size === 10).length;
            const before = mainCalls();

            fireEvent.click(screen.getByRole('button', { name: 'Backlog' }));

            expect(await screen.findByRole('heading', { name: 'Sprints' })).toBeInTheDocument();
            expect(mainCalls()).toBe(before);
        });

        it('avisa en board cuando hay mas tickets de los que caben', async () => {
            vi.mocked(listTicketsAction).mockResolvedValue(paginated([makeTicket({ id: 1 })], { total: 350 }) as any);
            await renderManager();

            fireEvent.click(screen.getByRole('button', { name: 'Board' }));

            expect(await screen.findByText(/Mostrando los primeros 100 tickets de 350/)).toBeInTheDocument();
        });
    });

    describe('enlace directo a un ticket', () => {
        it('abre el detalle del ticket indicado en la url', async () => {
            window.history.replaceState(null, '', '/tickets?ticket=12');
            vi.mocked(getTicketAction).mockResolvedValue(makeTicket({ id: 12, title: 'Ticket enlazado' }) as any);
            await renderManager();

            await waitFor(() => expect(getTicketAction).toHaveBeenCalledWith(12));
            expect(await screen.findByText('Ticket enlazado')).toBeInTheDocument();
        });

        it('avisa cuando el ticket enlazado no existe', async () => {
            window.history.replaceState(null, '', '/tickets?ticket=12');
            vi.mocked(getTicketAction).mockResolvedValue(null as any);
            await renderManager();

            expect(await screen.findByText('No se encontr\u00f3 el ticket 12.')).toBeInTheDocument();
        });

        it('avisa cuando el backend responde 403 y limpia el parametro de la url', async () => {
            window.history.replaceState(null, '', '/tickets?ticket=12');
            vi.mocked(getTicketAction).mockRejectedValue(new Error('request failed with 403'));
            await renderManager();

            expect(await screen.findByText('No tienes permiso para ver este ticket.')).toBeInTheDocument();
            expect(window.location.search).toBe('');
        });

        it('avisa cuando el backend responde 404', async () => {
            window.history.replaceState(null, '', '/tickets?ticket=12');
            vi.mocked(getTicketAction).mockRejectedValue(new Error('status 404'));
            await renderManager();

            expect(await screen.findByText('No se encontr\u00f3 el ticket 12.')).toBeInTheDocument();
        });

        it('avisa con un mensaje generico ante cualquier otro error', async () => {
            window.history.replaceState(null, '', '/tickets?ticket=12');
            vi.mocked(getTicketAction).mockRejectedValue(new Error('timeout'));
            await renderManager();

            expect(await screen.findByText('No se pudo abrir el ticket 12.')).toBeInTheDocument();
        });

        it('usa el mensaje generico cuando el error no trae texto', async () => {
            window.history.replaceState(null, '', '/tickets?ticket=12');
            vi.mocked(getTicketAction).mockRejectedValue({});
            await renderManager();

            expect(await screen.findByText('No se pudo abrir el ticket 12.')).toBeInTheDocument();
        });

        it('ignora un parametro de ticket que no es un numero', async () => {
            window.history.replaceState(null, '', '/tickets?ticket=abc');
            await renderManager();

            expect(getTicketAction).not.toHaveBeenCalled();
        });

        it('no consulta nada cuando la url no trae ticket', async () => {
            await renderManager();

            expect(getTicketAction).not.toHaveBeenCalled();
        });
    });

    describe('apertura del detalle desde la tabla', () => {
        it('abre el detalle y escribe el ticket en la url, y lo limpia al cerrar', async () => {
            await renderManager();

            fireEvent.click(await screen.findByText('TCK-1'));

            const modal = await screen.findByTestId('modal');
            expect(within(modal).getByRole('heading', { name: 'Ticket 1' })).toBeInTheDocument();
            expect(window.location.search).toBe('?ticket=1');

            fireEvent.click(within(screen.getByTestId('modal')).getByRole('button', { name: 'Cerrar' }));

            await waitFor(() => expect(window.location.search).toBe(''));
        });

        it('cierra el detalle desde el contenedor del modal', async () => {
            await renderManager();

            fireEvent.click(await screen.findByText('TCK-1'));
            await screen.findByTestId('modal');
            fireEvent.click(screen.getByRole('button', { name: 'modal-close' }));

            await waitFor(() => expect(screen.queryByTestId('modal')).toBeNull());
            expect(window.location.search).toBe('');
        });

        it('no rompe si al refrescar falla la consulta del ticket abierto', async () => {
            vi.mocked(getTicketAction).mockRejectedValue(new Error('timeout'));
            vi.mocked(changeTicketStatusAction).mockResolvedValue(makeTicket({ id: 1, status: 'testing' }) as any);
            await renderManager();

            fireEvent.click(await screen.findByText('TCK-1'));
            const modal = await screen.findByTestId('modal');
            fireEvent.click(within(modal).getByRole('button', { name: 'Pruebas' }));

            await waitFor(() => expect(getTicketAction).toHaveBeenCalledWith(1));
            expect(screen.getByTestId('modal')).toBeInTheDocument();
        });

        it('refresca el listado y el ticket abierto cuando el detalle cambia', async () => {
            vi.mocked(getTicketAction).mockResolvedValue(makeTicket({ id: 1, title: 'Ticket 1 actualizado' }) as any);
            vi.mocked(changeTicketStatusAction).mockResolvedValue(makeTicket({ id: 1, status: 'testing' }) as any);
            await renderManager();

            fireEvent.click(await screen.findByText('TCK-1'));
            const modal = await screen.findByTestId('modal');
            fireEvent.click(within(modal).getByRole('button', { name: 'Pruebas' }));

            await waitFor(() => expect(getTicketAction).toHaveBeenCalledWith(1));
            expect(await screen.findByText('Ticket 1 actualizado')).toBeInTheDocument();
        });
    });

    describe('acciones en linea de la tabla', () => {
        it('cambia el estado desde el selector de la fila', async () => {
            vi.mocked(changeTicketStatusAction).mockResolvedValue(makeTicket({ id: 1, status: 'testing' }) as any);
            await renderManager();

            const selects = within(rowOf('TCK-1')).getAllByRole('combobox');
            fireEvent.change(selects[1], { target: { value: 'testing' } });

            await waitFor(() => expect(changeTicketStatusAction).toHaveBeenCalledWith(1, 'testing'));
            await waitFor(() => expect((within(rowOf('TCK-1')).getAllByRole('combobox')[1] as HTMLSelectElement).value).toBe('testing'));
        });

        it('cambia el area desde el selector de la fila', async () => {
            vi.mocked(changeTicketAreaAction).mockResolvedValue(makeTicket({ id: 1, area: 'desarrollo' }) as any);
            await renderManager();

            const selects = within(rowOf('TCK-1')).getAllByRole('combobox');
            fireEvent.change(selects[0], { target: { value: 'desarrollo' } });

            await waitFor(() => expect(changeTicketAreaAction).toHaveBeenCalledWith(1, 'desarrollo'));
            await waitFor(() => expect((within(rowOf('TCK-1')).getAllByRole('combobox')[0] as HTMLSelectElement).value).toBe('desarrollo'));
        });

        it('asigna y desasigna el ticket desde la fila', async () => {
            vi.mocked(getUsersAction).mockResolvedValue({
                data: [{ id: 5, name: 'Ana', email: 'a@test.com', scope_code: 'platform' }],
            } as any);
            vi.mocked(assignTicketAction).mockResolvedValue(makeTicket({ id: 1, assigned_to_id: 5, assigned_to_name: 'Ana' }) as any);
            await renderManager();

            await screen.findByRole('option', { name: 'Ana' });
            const assignSelect = within(rowOf('TCK-1')).getAllByRole('combobox')[2];
            fireEvent.change(assignSelect, { target: { value: '5' } });

            await waitFor(() => expect(assignTicketAction).toHaveBeenCalledWith(1, 5));

            fireEvent.change(within(rowOf('TCK-1')).getAllByRole('combobox')[2], { target: { value: '' } });

            await waitFor(() => expect(assignTicketAction).toHaveBeenCalledWith(1, null));
        });

        it('mantiene la fila cuando la accion falla', async () => {
            vi.mocked(changeTicketStatusAction).mockRejectedValue(new Error('500'));
            await renderManager();

            fireEvent.change(within(rowOf('TCK-1')).getAllByRole('combobox')[1], { target: { value: 'testing' } });

            await waitFor(() => expect(changeTicketStatusAction).toHaveBeenCalled());
            await waitFor(() => expect((within(rowOf('TCK-1')).getAllByRole('combobox')[1] as HTMLSelectElement).value).toBe('open'));
        });

        it('no abre el detalle al usar los selectores de la fila', async () => {
            vi.mocked(changeTicketStatusAction).mockResolvedValue(makeTicket({ id: 1, status: 'testing' }) as any);
            await renderManager();

            const statusSelect = within(rowOf('TCK-1')).getAllByRole('combobox')[1];
            fireEvent.click(statusSelect);
            fireEvent.change(statusSelect, { target: { value: 'testing' } });

            await waitFor(() => expect(changeTicketStatusAction).toHaveBeenCalled());
            expect(screen.queryByTestId('modal')).toBeNull();
            expect(window.location.search).toBe('');
        });

        it('no rompe cuando falla el cambio de area', async () => {
            vi.mocked(changeTicketAreaAction).mockRejectedValue(new Error('500'));
            await renderManager();

            fireEvent.change(within(rowOf('TCK-1')).getAllByRole('combobox')[0], { target: { value: 'comercial' } });

            await waitFor(() => expect(changeTicketAreaAction).toHaveBeenCalled());
            await waitFor(() => expect((within(rowOf('TCK-1')).getAllByRole('combobox')[0] as HTMLSelectElement).value).toBe('soporte'));
        });

        it('no rompe cuando falla la asignacion', async () => {
            vi.mocked(getUsersAction).mockResolvedValue({
                data: [{ id: 5, name: 'Ana', email: 'a@test.com', scope_code: 'platform' }],
            } as any);
            vi.mocked(assignTicketAction).mockRejectedValue(new Error('500'));
            await renderManager();

            await screen.findByRole('option', { name: 'Ana' });
            fireEvent.change(within(rowOf('TCK-1')).getAllByRole('combobox')[2], { target: { value: '5' } });

            await waitFor(() => expect(assignTicketAction).toHaveBeenCalledWith(1, 5));
            expect(await screen.findByText('TCK-1')).toBeInTheDocument();
        });

        it('conserva el usuario asignado que no esta en la lista cargada', async () => {
            vi.mocked(listTicketsAction).mockResolvedValue(
                paginated([makeTicket({ id: 1, assigned_to_id: 77, assigned_to_name: 'Externo' })]) as any
            );
            await renderManager();

            expect(await screen.findByRole('option', { name: 'Externo' })).toBeInTheDocument();
        });

        it('muestra el avatar del asignado resolviendo rutas relativas contra S3', async () => {
            vi.mocked(listTicketsAction).mockResolvedValue(paginated([
                makeTicket({ id: 1, assigned_to_avatar_url: 'avatars/ana.png' }),
                makeTicket({ id: 2, assigned_to_avatar_url: 'https://cdn.test/luis.png' }),
            ]) as any);
            await renderManager();

            await screen.findByText('TCK-1');
            expect(rowOf('TCK-1').querySelector('img')).toHaveAttribute(
                'src',
                expect.stringContaining('/avatars/ana.png')
            );
            expect(rowOf('TCK-2').querySelector('img')).toHaveAttribute('src', 'https://cdn.test/luis.png');
        });
    });

    describe('permisos', () => {
        it('bloquea los selectores y no carga usuarios si no es super admin', async () => {
            vi.mocked(usePermissions).mockReturnValue({ isSuperAdmin: false } as any);
            vi.mocked(listTicketsAction).mockResolvedValue(
                paginated([makeTicket({ id: 1, assigned_to_name: 'Ana' })]) as any
            );
            await renderManager();

            await screen.findByText('TCK-1');
            const selects = within(rowOf('TCK-1')).getAllByRole('combobox');
            expect(selects).toHaveLength(2);
            expect(selects[0]).toBeDisabled();
            expect(selects[1]).toBeDisabled();
            expect(within(rowOf('TCK-1')).getByText('Ana')).toBeInTheDocument();
            expect(getUsersAction).not.toHaveBeenCalled();
        });

        it('solo ofrece usuarios de plataforma o super usuarios en el filtro de asignado', async () => {
            vi.mocked(getUsersAction).mockResolvedValue({
                data: [
                    { id: 5, name: 'Ana', scope_code: 'platform' },
                    { id: 6, name: 'Super', is_super_user: true },
                    { id: 7, name: 'Cliente', scope_code: 'business' },
                    { id: 8, name: '', scope_code: 'platform' },
                ],
            } as any);
            await renderManager();

            await waitFor(() =>
                expect(screen.getByTestId('available-filters').textContent).toContain('assigned_to_id')
            );
            expect(screen.getAllByRole('option', { name: 'Ana' }).length).toBeGreaterThan(0);
            expect(screen.queryByRole('option', { name: 'Cliente' })).toBeNull();
        });
    });

    describe('filtros y orden', () => {
        it('traduce cada filtro al parametro del listado', async () => {
            await renderManager();

            fireEvent.click(screen.getByText('add-search'));
            await waitFor(() => expect(lastListCall()).toMatchObject({ search: 'pago' }));

            fireEvent.click(screen.getByText('add-status'));
            await waitFor(() => expect(lastListCall()).toMatchObject({ status: 'open' }));

            fireEvent.click(screen.getByText('add-only-mine'));
            await waitFor(() => expect(lastListCall()).toMatchObject({ only_mine: true }));

            fireEvent.click(screen.getByText('add-escalated'));
            await waitFor(() => expect(lastListCall()).toMatchObject({ escalated: true }));

            fireEvent.click(screen.getByText('add-assigned'));
            await waitFor(() => expect(lastListCall()).toMatchObject({ assigned_to_id: 5 }));
        });

        it('descarta el filtro de asignado cuando llega vacio', async () => {
            await renderManager();

            fireEvent.click(screen.getByText('add-assigned'));
            await waitFor(() => expect(lastListCall()).toMatchObject({ assigned_to_id: 5 }));

            fireEvent.click(screen.getByText('add-assigned-empty'));

            await waitFor(() => expect(lastListCall().assigned_to_id).toBeUndefined());
        });

        it('quita un filtro activo', async () => {
            await renderManager();

            fireEvent.click(screen.getByText('add-status'));
            await waitFor(() => expect(screen.getByTestId('active-filters').textContent).toContain('status=Abierto'));

            fireEvent.click(screen.getByText('remove-status'));

            await waitFor(() => expect(lastListCall().status).toBeUndefined());
            expect(screen.getByTestId('active-filters').textContent).not.toContain('status=');
        });

        it('describe los filtros activos con etiquetas legibles', async () => {
            vi.mocked(getUsersAction).mockResolvedValue({
                data: [{ id: 5, name: 'Ana', scope_code: 'platform' }],
            } as any);
            await renderManager();

            fireEvent.click(screen.getByText('add-search'));
            fireEvent.click(screen.getByText('add-status'));
            fireEvent.click(screen.getByText('add-area'));
            fireEvent.click(screen.getByText('add-priority'));
            fireEvent.click(screen.getByText('add-type'));
            fireEvent.click(screen.getByText('add-source'));
            fireEvent.click(screen.getByText('add-only-mine'));
            fireEvent.click(screen.getByText('add-escalated'));
            fireEvent.click(screen.getByText('add-assigned'));

            await waitFor(() => {
                const text = screen.getByTestId('active-filters').textContent || '';
                expect(text).toContain('search=pago');
                expect(text).toContain('status=Abierto');
                expect(text).toContain('area=Soporte');
                expect(text).toContain('priority=Alta');
                expect(text).toContain('type=Bug');
                expect(text).toContain('source=Interno');
                expect(text).toContain('only_mine=Si');
                expect(text).toContain('escalated=Si');
                expect(text).toContain('assigned_to_id=Ana');
            });
        });

        it('acepta un filtro booleano ya normalizado', async () => {
            await renderManager();

            fireEvent.click(screen.getByText('add-only-mine-bool'));

            await waitFor(() => expect(lastListCall()).toMatchObject({ only_mine: true }));
        });

        it('muestra el valor crudo de un filtro cuya etiqueta no conoce', async () => {
            await renderManager();

            fireEvent.click(screen.getByText('add-status-desconocido'));
            fireEvent.click(screen.getByText('add-area-desconocida'));
            fireEvent.click(screen.getByText('add-priority-desconocida'));
            fireEvent.click(screen.getByText('add-type-desconocido'));
            fireEvent.click(screen.getByText('add-source-negocio'));
            fireEvent.click(screen.getByText('add-assigned'));

            await waitFor(() => {
                const text = screen.getByTestId('active-filters').textContent || '';
                expect(text).toContain('status=zzz');
                expect(text).toContain('area=zzz');
                expect(text).toContain('priority=zzz');
                expect(text).toContain('type=zzz');
                expect(text).toContain('source=Negocio');
                expect(text).toContain('assigned_to_id=Usuario 5');
            });
        });

        it('cambia el criterio de orden', async () => {
            await renderManager();

            fireEvent.click(screen.getByText('sort-priority'));

            await waitFor(() => expect(lastListCall()).toMatchObject({ sort_by: 'priority', sort_order: 'asc', page: 1 }));
        });
    });

    describe('sprints', () => {
        it('preselecciona el sprint activo al pedir el board', async () => {
            vi.mocked(listSprintsAction).mockResolvedValue({
                data: [sprint(1, 'Sprint Uno', 'planned'), sprint(2, 'Sprint Dos', 'active')],
            } as any);
            await renderManager();

            fireEvent.click(screen.getByRole('button', { name: 'Board' }));

            expect(await screen.findByRole('option', { name: 'Sprint Dos' })).toBeInTheDocument();
            await waitFor(() => expect(lastListCall()).toMatchObject({ sprint_id: 2 }));
        });

        it('permite filtrar el board por tickets sin sprint o por todos', async () => {
            vi.mocked(listSprintsAction).mockResolvedValue({ data: [sprint(1, 'Sprint Uno', 'planned')] } as any);
            await renderManager();

            fireEvent.click(screen.getByRole('button', { name: 'Board' }));
            const selector = await screen.findByRole('combobox');

            fireEvent.change(selector, { target: { value: 'none' } });
            await waitFor(() => expect(lastListCall()).toMatchObject({ sprint_id: 'none' }));

            fireEvent.change(screen.getByRole('combobox'), { target: { value: '1' } });
            await waitFor(() => expect(lastListCall()).toMatchObject({ sprint_id: 1 }));

            fireEvent.change(screen.getByRole('combobox'), { target: { value: '' } });
            await waitFor(() => expect(lastListCall().sprint_id).toBeUndefined());
        });

        it('tolera que el listado de sprints venga sin datos', async () => {
            vi.mocked(listSprintsAction).mockResolvedValue(undefined as any);
            await renderManager();

            fireEvent.click(screen.getByRole('button', { name: 'Board' }));

            expect(await screen.findByRole('option', { name: 'Todos los sprints' })).toBeInTheDocument();
            expect(screen.queryByRole('option', { name: 'Sprint Uno' })).toBeNull();
        });

        it('propaga el error de sprints a la vista de backlog', async () => {
            vi.mocked(listSprintsAction).mockRejectedValue(new Error('modulo caido'));
            await renderManager();

            fireEvent.click(screen.getByRole('button', { name: 'Backlog' }));

            expect(await screen.findByText(/No se pudieron cargar los sprints/)).toBeInTheDocument();
        });
    });

    describe('creacion de tickets', () => {
        beforeEach(() => {
            vi.mocked(createTicketAction).mockResolvedValue(makeTicket({ id: 99, code: 'TCK-99' }) as any);
            vi.mocked(uploadAttachmentAction).mockResolvedValue({} as any);
            vi.mocked(addCommentAction).mockResolvedValue({} as any);
        });

        it('crea el ticket, sube cada adjunto y guarda el comentario en ese orden', async () => {
            await renderManager();

            openCreateModal();
            fillCreateForm();
            attachToCreateForm([makeFile('uno.png'), makeFile('dos.png')]);
            fireEvent.change(screen.getByPlaceholderText('Opcional: contexto adicional para el equipo'), { target: { value: 'contexto' } });
            submitCreateForm();

            await waitFor(() => expect(addCommentAction).toHaveBeenCalledWith(99, 'contexto', false));
            expect(createTicketAction).toHaveBeenCalledTimes(1);
            expect(uploadAttachmentAction).toHaveBeenCalledTimes(2);
            expect(vi.mocked(uploadAttachmentAction).mock.calls[0][0]).toBe(99);
            expect((vi.mocked(uploadAttachmentAction).mock.calls[0][1].get('file') as File).name).toBe('uno.png');
            expect((vi.mocked(uploadAttachmentAction).mock.calls[1][1].get('file') as File).name).toBe('dos.png');

            const createOrder = vi.mocked(createTicketAction).mock.invocationCallOrder[0];
            const firstUpload = vi.mocked(uploadAttachmentAction).mock.invocationCallOrder[0];
            const secondUpload = vi.mocked(uploadAttachmentAction).mock.invocationCallOrder[1];
            const commentOrder = vi.mocked(addCommentAction).mock.invocationCallOrder[0];
            expect(createOrder).toBeLessThan(firstUpload);
            expect(firstUpload).toBeLessThan(secondUpload);
            expect(secondUpload).toBeLessThan(commentOrder);

            await waitFor(() => expect(screen.queryByTestId('modal')).toBeNull());
        });

        it('no llama al comentario cuando el campo queda vacio', async () => {
            await renderManager();

            openCreateModal();
            fillCreateForm();
            submitCreateForm();

            await waitFor(() => expect(createTicketAction).toHaveBeenCalled());
            expect(addCommentAction).not.toHaveBeenCalled();
            expect(uploadAttachmentAction).not.toHaveBeenCalled();
        });

        it('recarga el listado despues de crear', async () => {
            await renderManager();
            const before = vi.mocked(listTicketsAction).mock.calls.length;

            openCreateModal();
            fillCreateForm();
            submitCreateForm();

            await waitFor(() => expect(vi.mocked(listTicketsAction).mock.calls.length).toBeGreaterThan(before));
        });

        it('cierra el modal y avisa cuando falla la subida de un adjunto', async () => {
            vi.mocked(uploadAttachmentAction)
                .mockResolvedValueOnce({} as any)
                .mockRejectedValueOnce(new Error('S3 caido'));
            await renderManager();

            openCreateModal();
            fillCreateForm();
            attachToCreateForm([makeFile('uno.png'), makeFile('dos.png')]);
            submitCreateForm();

            expect(await screen.findByText('Ticket TCK-99 creado, pero no se pudieron subir 1 adjunto.')).toBeInTheDocument();
            expect(screen.queryByTestId('modal')).toBeNull();
        });

        it('avisa en plural cuando fallan varios adjuntos', async () => {
            vi.mocked(uploadAttachmentAction).mockRejectedValue(new Error('S3 caido'));
            await renderManager();

            openCreateModal();
            fillCreateForm();
            attachToCreateForm([makeFile('uno.png'), makeFile('dos.png')]);
            submitCreateForm();

            expect(await screen.findByText('Ticket TCK-99 creado, pero no se pudieron subir 2 adjuntos.')).toBeInTheDocument();
        });

        it('avisa cuando falla el comentario inicial', async () => {
            vi.mocked(addCommentAction).mockRejectedValue(new Error('sin permiso'));
            await renderManager();

            openCreateModal();
            fillCreateForm();
            fireEvent.change(screen.getByPlaceholderText('Opcional: contexto adicional para el equipo'), { target: { value: 'contexto' } });
            submitCreateForm();

            expect(await screen.findByText('Ticket TCK-99 creado, pero no se pudo agregar el comentario inicial.')).toBeInTheDocument();
        });

        it('enumera los dos problemas cuando fallan adjunto y comentario', async () => {
            vi.mocked(uploadAttachmentAction).mockRejectedValue(new Error('S3 caido'));
            vi.mocked(addCommentAction).mockRejectedValue(new Error('sin permiso'));
            await renderManager();

            openCreateModal();
            fillCreateForm();
            attachToCreateForm([makeFile('uno.png')]);
            fireEvent.change(screen.getByPlaceholderText('Opcional: contexto adicional para el equipo'), { target: { value: 'contexto' } });
            submitCreateForm();

            expect(await screen.findByText(
                'Ticket TCK-99 creado, pero no se pudieron subir 1 adjunto y no se pudo agregar el comentario inicial.'
            )).toBeInTheDocument();
        });

        it('usa el id del ticket en el aviso cuando no llega el codigo', async () => {
            vi.mocked(createTicketAction).mockResolvedValue({ id: 99 } as any);
            vi.mocked(uploadAttachmentAction).mockRejectedValue(new Error('S3 caido'));
            await renderManager();

            openCreateModal();
            fillCreateForm();
            attachToCreateForm([makeFile('uno.png')]);
            submitCreateForm();

            expect(await screen.findByText('Ticket #99 creado, pero no se pudieron subir 1 adjunto.')).toBeInTheDocument();
        });

        it('deja el modal abierto y muestra el error cuando falla la creacion', async () => {
            vi.mocked(createTicketAction).mockRejectedValue(new Error('titulo duplicado'));
            await renderManager();

            openCreateModal();
            fillCreateForm();
            submitCreateForm();

            expect(await screen.findByText('titulo duplicado')).toBeInTheDocument();
            expect(screen.getByTestId('modal')).toBeInTheDocument();
        });

        it('cierra el modal al cancelar', async () => {
            await renderManager();

            openCreateModal();
            expect(screen.getByTestId('modal')).toBeInTheDocument();

            fireEvent.click(screen.getByRole('button', { name: 'Cancelar' }));

            expect(screen.queryByTestId('modal')).toBeNull();
        });

        it('cierra el modal de creacion desde el contenedor', async () => {
            await renderManager();

            openCreateModal();
            fireEvent.click(screen.getByRole('button', { name: 'modal-close' }));

            expect(screen.queryByTestId('modal')).toBeNull();
        });

        it('permite cerrar el aviso de error', async () => {
            vi.mocked(uploadAttachmentAction).mockRejectedValue(new Error('S3 caido'));
            await renderManager();

            openCreateModal();
            fillCreateForm();
            attachToCreateForm([makeFile('uno.png')]);
            submitCreateForm();

            const aviso = await screen.findByText('Ticket TCK-99 creado, pero no se pudieron subir 1 adjunto.');
            fireEvent.click(within(aviso.parentElement as HTMLElement).getByRole('button', { name: 'Cerrar' }));

            await waitFor(() => expect(aviso).not.toBeInTheDocument());
        });
    });

    describe('board: movimiento optimista', () => {
        const goToBoard = async () => {
            fireEvent.click(screen.getByRole('button', { name: 'Board' }));
            await screen.findByRole('heading', { name: 'To Do' });
        };

        const dragCard = (code: string, column: string) => {
            const dataTransfer = makeDataTransfer();
            fireEvent.dragStart(boardCard(code), { dataTransfer });
            fireEvent.drop(boardColumn(column), { dataTransfer });
        };

        it('mueve la tarjeta al instante y confirma con el servidor', async () => {
            let resolveChange: (value: any) => void = () => {};
            vi.mocked(changeTicketStatusAction).mockReturnValue(new Promise((res) => { resolveChange = res; }) as any);
            await renderManager();
            await goToBoard();

            expect(within(boardColumn('To Do')).getByText('TCK-1')).toBeInTheDocument();

            dragCard('TCK-1', 'In Progress');

            await waitFor(() => expect(within(boardColumn('In Progress')).getByText('TCK-1')).toBeInTheDocument());
            expect(within(boardColumn('To Do')).queryByText('TCK-1')).toBeNull();
            expect(changeTicketStatusAction).toHaveBeenCalledWith(1, 'in_development');

            resolveChange(makeTicket({ id: 1, status: 'in_development', title: 'Ticket 1 confirmado' }));

            await waitFor(() => expect(screen.getByText('Ticket 1 confirmado')).toBeInTheDocument());
        });

        it('devuelve la tarjeta a su columna y avisa cuando el cambio falla', async () => {
            vi.mocked(changeTicketStatusAction).mockRejectedValue(new Error('conflicto'));
            await renderManager();
            await goToBoard();

            dragCard('TCK-1', 'Done');

            expect(await screen.findByText('No se pudo cambiar el estado del ticket. Intenta de nuevo.')).toBeInTheDocument();
            expect(within(boardColumn('To Do')).getByText('TCK-1')).toBeInTheDocument();
            expect(within(boardColumn('Done')).queryByText('TCK-1')).toBeNull();
        });

        it('no llama al servidor si la tarjeta ya estaba en ese estado', async () => {
            await renderManager();
            await goToBoard();

            dragCard('TCK-1', 'To Do');

            expect(changeTicketStatusAction).not.toHaveBeenCalled();
        });

        it('abre el detalle al hacer click en una tarjeta del board', async () => {
            await renderManager();
            await goToBoard();

            fireEvent.click(boardCard('TCK-1'));

            expect(await screen.findByTestId('modal')).toBeInTheDocument();
            expect(window.location.search).toBe('?ticket=1');
        });

        it('no permite arrastrar en el board si no es super admin', async () => {
            vi.mocked(usePermissions).mockReturnValue({ isSuperAdmin: false } as any);
            await renderManager();
            await goToBoard();

            expect(boardCard('TCK-1')).toHaveAttribute('draggable', 'false');
            dragCard('TCK-1', 'Done');

            expect(changeTicketStatusAction).not.toHaveBeenCalled();
        });
    });
});
