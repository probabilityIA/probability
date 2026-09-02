import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, createEvent, waitFor, within } from '@testing-library/react';
import TicketsBacklog from './TicketsBacklog';
import { Ticket } from '../../domain/types';
import { Sprint } from '@/services/modules/sprints/domain/types';

vi.mock('@/shared/ui', () => ({
    Modal: ({ isOpen, title, onClose, children }: any) =>
        isOpen ? (
            <div data-testid="modal">
                <h4>{title}</h4>
                <button onClick={onClose}>modal-close</button>
                {children}
            </div>
        ) : null,
}));

vi.mock('../../infra/actions', () => ({
    listTicketsAction: vi.fn(),
    changeTicketSprintAction: vi.fn(),
}));

vi.mock('@/services/modules/sprints/infra/actions', () => ({
    createSprintAction: vi.fn(),
    updateSprintAction: vi.fn(),
    changeSprintStatusAction: vi.fn(),
}));

import { listTicketsAction, changeTicketSprintAction } from '../../infra/actions';
import {
    createSprintAction,
    updateSprintAction,
    changeSprintStatusAction,
} from '@/services/modules/sprints/infra/actions';

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
    comments_count: 0,
    attachments_count: 0,
    ...over,
});

const makeSprint = (over: Partial<Sprint> & { id: number }): Sprint => ({
    name: 'Sprint ' + over.id,
    goal: '',
    start_date: '2026-02-01T12:00:00',
    end_date: '2026-02-15T12:00:00',
    status: 'planned',
    created_by_id: 1,
    created_by_name: 'admin',
    ticket_count: 0,
    done_count: 0,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...over,
});

const activeSprint = makeSprint({ id: 1, name: 'Sprint Uno', status: 'active', ticket_count: 4, done_count: 2 });
const plannedSprint = makeSprint({ id: 2, name: 'Sprint Dos' });

const mockBuckets = (bySprint: Record<string, Ticket[]>, totals: Record<string, number> = {}) => {
    vi.mocked(listTicketsAction).mockImplementation(async (params: any) => {
        const key = String(params?.sprint_id);
        const data = bySprint[key] || [];
        return { data, total: totals[key] ?? data.length, page: 1, page_size: 100, total_pages: 1 } as any;
    });
};

const makeDataTransfer = () => {
    const store: Record<string, string> = {};
    return {
        effectAllowed: '',
        dropEffect: '',
        setData: (key: string, value: string) => { store[key] = value; },
        getData: (key: string) => store[key] ?? '',
    };
};

const rowOf = (code: string): HTMLElement => screen.getByText(code).closest('[draggable]') as HTMLElement;

const backlogZone = (): HTMLElement =>
    screen.getByRole('heading', { name: 'Backlog' }).parentElement!.parentElement as HTMLElement;

const sprintZone = (name: string): HTMLElement =>
    screen.getByText(name).closest('button')!.parentElement!.parentElement as HTMLElement;

const dragTo = (code: string, target: HTMLElement) => {
    const dataTransfer = makeDataTransfer();
    fireEvent.dragStart(rowOf(code), { dataTransfer });
    fireEvent.dragOver(target, { dataTransfer });
    fireEvent.drop(target, { dataTransfer });
};

const baseParams = { sort_by: 'created_at', sort_order: 'desc' as const };

const setup = (props: Partial<React.ComponentProps<typeof TicketsBacklog>> = {}) => {
    const onOpenTicket = vi.fn();
    const onSprintsChanged = vi.fn().mockResolvedValue(undefined);
    const utils = render(
        <TicketsBacklog
            sprints={[activeSprint, plannedSprint]}
            sprintsLoading={false}
            sprintsError={null}
            canManage
            baseParams={baseParams}
            onOpenTicket={onOpenTicket}
            onSprintsChanged={onSprintsChanged}
            {...props}
        />
    );
    return { onOpenTicket, onSprintsChanged, ...utils };
};

describe('TicketsBacklog', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        mockBuckets({});
        vi.mocked(changeTicketSprintAction).mockImplementation(async (id: number, sprintId: number | null) =>
            makeTicket({ id, sprint_id: sprintId }) as any
        );
        vi.mocked(createSprintAction).mockResolvedValue({} as any);
        vi.mocked(updateSprintAction).mockResolvedValue({} as any);
        vi.mocked(changeSprintStatusAction).mockResolvedValue({} as any);
    });

    describe('carga de buckets', () => {
        it('pide el backlog con sprint_id none y los filtros activos', async () => {
            setup();

            await waitFor(() => expect(listTicketsAction).toHaveBeenCalledWith({
                sort_by: 'created_at',
                sort_order: 'desc',
                sprint_id: 'none',
                page: 1,
                page_size: 100,
            }));
        });

        it('descarta del backlog los tickets ya terminados', async () => {
            mockBuckets({
                none: [
                    makeTicket({ id: 1, status: 'open' }),
                    makeTicket({ id: 2, status: 'resolved' }),
                    makeTicket({ id: 3, status: 'closed' }),
                    makeTicket({ id: 4, status: 'wont_fix' }),
                ],
            });
            setup();

            expect(await screen.findByText('TCK-1')).toBeInTheDocument();
            expect(screen.queryByText('TCK-2')).toBeNull();
            expect(screen.queryByText('TCK-3')).toBeNull();
            expect(screen.queryByText('TCK-4')).toBeNull();
        });

        it('conserva los tickets terminados dentro de un sprint', async () => {
            mockBuckets({ none: [], '1': [makeTicket({ id: 9, status: 'closed' })] });
            setup();

            expect(await screen.findByText('TCK-9')).toBeInTheDocument();
        });

        it('carga automaticamente el sprint activo y no el planeado', async () => {
            mockBuckets({ none: [], '1': [makeTicket({ id: 9 })] });
            setup();

            await waitFor(() => expect(listTicketsAction).toHaveBeenCalledWith(expect.objectContaining({ sprint_id: 1 })));
            expect(listTicketsAction).not.toHaveBeenCalledWith(expect.objectContaining({ sprint_id: 2 }));
        });

        it('carga el sprint planeado solo cuando se despliega', async () => {
            mockBuckets({ none: [], '2': [makeTicket({ id: 20 })] });
            setup();

            await screen.findByText('Sprint Dos');
            fireEvent.click(screen.getByText('Sprint Dos').closest('button') as HTMLElement);

            await waitFor(() => expect(listTicketsAction).toHaveBeenCalledWith(expect.objectContaining({ sprint_id: 2 })));
            expect(await screen.findByText('TCK-20')).toBeInTheDocument();
        });

        it('oculta los tickets del sprint al plegarlo de nuevo', async () => {
            mockBuckets({ none: [], '1': [makeTicket({ id: 9 })] });
            setup();

            expect(await screen.findByText('TCK-9')).toBeInTheDocument();
            fireEvent.click(screen.getByText('Sprint Uno').closest('button') as HTMLElement);

            await waitFor(() => expect(screen.queryByText('TCK-9')).toBeNull());
        });

        it('muestra el mensaje de error cuando falla el listado del bucket', async () => {
            vi.mocked(listTicketsAction).mockRejectedValue(new Error('500'));
            setup();

            expect(await screen.findAllByText('No se pudieron cargar los tickets.')).toHaveLength(2);
        });

        it('avisa cuando el backlog esta vacio', async () => {
            setup();

            expect(await screen.findByText('El backlog est\u00e1 vac\u00edo')).toBeInTheDocument();
        });

        it('avisa cuando el sprint no tiene tickets', async () => {
            setup();

            expect(await screen.findByText('Sin tickets en este sprint')).toBeInTheDocument();
        });

        it('avisa cuando hay mas tickets de los que se muestran', async () => {
            mockBuckets({ none: [makeTicket({ id: 1 })] }, { none: 240 });
            setup();

            expect(await screen.findByText('Mostrando los primeros 100 de 240. Usa los filtros para acotar la b\u00fasqueda.')).toBeInTheDocument();
        });

        it('recarga los buckets cuando cambian los filtros', async () => {
            const { rerender } = setup();

            await waitFor(() => expect(listTicketsAction).toHaveBeenCalled());
            const before = vi.mocked(listTicketsAction).mock.calls.length;

            rerender(
                <TicketsBacklog
                    sprints={[activeSprint, plannedSprint]}
                    sprintsLoading={false}
                    sprintsError={null}
                    canManage
                    baseParams={{ ...baseParams, status: 'open' }}
                    onOpenTicket={vi.fn()}
                    onSprintsChanged={vi.fn()}
                />
            );

            await waitFor(() => expect(vi.mocked(listTicketsAction).mock.calls.length).toBeGreaterThan(before));
            expect(listTicketsAction).toHaveBeenCalledWith(expect.objectContaining({ status: 'open', sprint_id: 'none' }));
        });
    });

    describe('lista de sprints', () => {
        it('ordena los sprints por estado: activo, planeado y cerrado', async () => {
            const closed = makeSprint({ id: 3, name: 'Sprint Tres', status: 'closed' });
            setup({ sprints: [closed, plannedSprint, activeSprint] });

            const names = screen.getAllByText(/Sprint (Uno|Dos|Tres)/).map((n) => n.textContent);
            expect(names).toEqual(['Sprint Uno', 'Sprint Dos', 'Sprint Tres']);
            await waitFor(() => expect(listTicketsAction).toHaveBeenCalled());
        });

        it('muestra el progreso, el estado y el objetivo del sprint', async () => {
            setup({ sprints: [makeSprint({ id: 1, name: 'Sprint Uno', status: 'active', ticket_count: 5, done_count: 3, goal: 'Cerrar la facturacion' })] });

            expect(screen.getByText('3 de 5 terminados')).toBeInTheDocument();
            expect(screen.getByText('Activo')).toBeInTheDocument();
            expect(screen.getByText('Objetivo: Cerrar la facturacion')).toBeInTheDocument();
            await waitFor(() => expect(listTicketsAction).toHaveBeenCalled());
        });

        it('muestra el aviso de sprints en carga', () => {
            setup({ sprintsLoading: true, sprints: [] });

            expect(screen.getByText('Cargando sprints...')).toBeInTheDocument();
        });

        it('muestra el error de sprints en vez del estado de carga', () => {
            setup({ sprintsLoading: true, sprints: [], sprintsError: 'El modulo no responde' });

            expect(screen.getByText('El modulo no responde')).toBeInTheDocument();
            expect(screen.queryByText('Cargando sprints...')).toBeNull();
        });

        it('ordena los sprints cerrados del mas reciente al mas antiguo', async () => {
            setup({
                sprints: [
                    makeSprint({ id: 3, name: 'Cerrado Viejo', status: 'closed', start_date: '2026-01-01T12:00:00' }),
                    makeSprint({ id: 4, name: 'Cerrado Nuevo', status: 'closed', start_date: '2026-05-01T12:00:00' }),
                ],
            });

            const names = screen.getAllByText(/Cerrado (Viejo|Nuevo)/).map((n) => n.textContent);
            expect(names).toEqual(['Cerrado Nuevo', 'Cerrado Viejo']);
            await waitFor(() => expect(listTicketsAction).toHaveBeenCalled());
        });

        it('usa la fecha de creacion para ordenar sprints sin fecha de inicio', async () => {
            setup({
                sprints: [
                    makeSprint({ id: 3, name: 'Sprint Tarde', start_date: '', created_at: '2026-06-01T12:00:00' }),
                    makeSprint({ id: 4, name: 'Sprint Temprano', start_date: '', created_at: '2026-01-01T12:00:00' }),
                ],
            });

            const names = screen.getAllByText(/Sprint (Tarde|Temprano)/).map((n) => n.textContent);
            expect(names).toEqual(['Sprint Temprano', 'Sprint Tarde']);
            await waitFor(() => expect(listTicketsAction).toHaveBeenCalled());
        });

        it('muestra las fechas tal cual cuando no son fechas validas', async () => {
            setup({ sprints: [makeSprint({ id: 3, name: 'Sprint Raro', start_date: 'sin-fecha', end_date: '' })] });

            expect(screen.getByText('sin-fecha -')).toBeInTheDocument();
            await waitFor(() => expect(listTicketsAction).toHaveBeenCalled());
        });

        it('usa cero como contador cuando el sprint no trae totales', async () => {
            setup({ sprints: [makeSprint({ id: 3, name: 'Sprint Sin Datos', ticket_count: undefined as any, done_count: undefined as any })] });

            expect(screen.getByText('0 de 0 terminados')).toBeInTheDocument();
            await waitFor(() => expect(listTicketsAction).toHaveBeenCalled());
        });

        it('cae al estilo de planeado ante un estado desconocido', async () => {
            setup({ sprints: [makeSprint({ id: 3, name: 'Sprint Raro', status: 'archivado' as any })] });

            expect(screen.getByText('Planeado')).toBeInTheDocument();
            await waitFor(() => expect(listTicketsAction).toHaveBeenCalled());
        });

        it('avisa cuando no hay sprints creados', () => {
            setup({ sprints: [] });

            expect(screen.getByText('No hay sprints creados todav\u00eda.')).toBeInTheDocument();
        });
    });

    describe('filas de ticket', () => {
        it('abre el ticket al hacer click en su fila', async () => {
            const ticket = makeTicket({ id: 1 });
            mockBuckets({ none: [ticket] });
            const { onOpenTicket } = setup();

            fireEvent.click(await screen.findByText('TCK-1'));

            expect(onOpenTicket).toHaveBeenCalledWith(ticket);
        });

        it('muestra el avatar cuando hay url y la inicial cuando no', async () => {
            mockBuckets({ none: [makeTicket({ id: 1, assigned_to_name: 'bruno diaz' }), makeTicket({ id: 2 })] });
            setup({ getAvatarUrl: (t) => (t.id === 1 ? 'https://cdn.test/b.png' : '') });

            await screen.findByText('TCK-1');
            expect(rowOf('TCK-1').querySelector('img')).toHaveAttribute('src', 'https://cdn.test/b.png');
            expect(within(rowOf('TCK-2')).getByText('-')).toBeInTheDocument();
        });

        it('muestra el estado del ticket en la fila', async () => {
            mockBuckets({ none: [makeTicket({ id: 1, status: 'blocked' })] });
            setup();

            expect(await screen.findByText('Bloqueado')).toBeInTheDocument();
        });
    });

    describe('mover tickets entre sprint y backlog', () => {
        it('mueve un ticket del backlog al sprint con el id del sprint', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })] });
            const { onSprintsChanged } = setup();

            await screen.findByText('TCK-5');
            dragTo('TCK-5', sprintZone('Sprint Uno'));

            await waitFor(() => expect(changeTicketSprintAction).toHaveBeenCalledWith(5, 1));
            await waitFor(() => expect(onSprintsChanged).toHaveBeenCalled());
        });

        it('mueve un ticket del sprint al backlog con sprint nulo', async () => {
            mockBuckets({ none: [], '1': [makeTicket({ id: 6 })] });
            setup();

            await screen.findByText('TCK-6');
            dragTo('TCK-6', backlogZone());

            await waitFor(() => expect(changeTicketSprintAction).toHaveBeenCalledWith(6, null));
        });

        it('no llama a la accion al soltar en el mismo contenedor', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })] });
            setup();

            await screen.findByText('TCK-5');
            dragTo('TCK-5', backlogZone());

            expect(changeTicketSprintAction).not.toHaveBeenCalled();
        });

        it('no llama a la accion si el id soltado no esta en ningun bucket', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })] });
            setup();

            await screen.findByText('TCK-5');
            const dataTransfer = makeDataTransfer();
            dataTransfer.setData('text/plain', '404');
            fireEvent.drop(sprintZone('Sprint Uno'), { dataTransfer });

            expect(changeTicketSprintAction).not.toHaveBeenCalled();
        });

        it('no llama a la accion cuando el drop no trae ningun id', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })] });
            setup();

            await screen.findByText('TCK-5');
            fireEvent.drop(sprintZone('Sprint Uno'), { dataTransfer: makeDataTransfer() });

            expect(changeTicketSprintAction).not.toHaveBeenCalled();
        });

        it('mueve la fila al sprint destino antes de que responda el servidor', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })], '1': [] });
            let resolveMove: (value: any) => void = () => {};
            vi.mocked(changeTicketSprintAction).mockReturnValue(new Promise((res) => { resolveMove = res; }) as any);
            setup();

            await screen.findByText('TCK-5');
            expect(within(backlogZone()).getByText('TCK-5')).toBeInTheDocument();

            dragTo('TCK-5', sprintZone('Sprint Uno'));

            await waitFor(() => expect(within(sprintZone('Sprint Uno')).getByText('TCK-5')).toBeInTheDocument());
            expect(within(backlogZone()).queryByText('TCK-5')).toBeNull();

            resolveMove(makeTicket({ id: 5, sprint_id: 1 }));
            await waitFor(() => expect(within(sprintZone('Sprint Uno')).getByText('TCK-5')).toBeInTheDocument());
        });

        it('devuelve la fila a su sitio y avisa cuando la accion falla', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })], '1': [] });
            vi.mocked(changeTicketSprintAction).mockRejectedValue(new Error('403'));
            const { onSprintsChanged } = setup();

            await screen.findByText('TCK-5');
            dragTo('TCK-5', sprintZone('Sprint Uno'));

            expect(await screen.findByText('No se pudo mover el ticket de sprint. Intenta de nuevo.')).toBeInTheDocument();
            expect(within(backlogZone()).getByText('TCK-5')).toBeInTheDocument();
            expect(within(sprintZone('Sprint Uno')).queryByText('TCK-5')).toBeNull();
            expect(onSprintsChanged).not.toHaveBeenCalled();
        });

        it('permite cerrar el aviso de error de movimiento', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })], '1': [] });
            vi.mocked(changeTicketSprintAction).mockRejectedValue(new Error('403'));
            setup();

            await screen.findByText('TCK-5');
            dragTo('TCK-5', sprintZone('Sprint Uno'));
            const aviso = await screen.findByText('No se pudo mover el ticket de sprint. Intenta de nuevo.');
            fireEvent.click(within(aviso.parentElement as HTMLElement).getByRole('button', { name: 'Cerrar' }));

            await waitFor(() => expect(aviso).not.toBeInTheDocument());
        });

        it('mueve el ticket a un sprint plegado aunque no tenga sus tickets cargados', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })] });
            setup();

            await screen.findByText('TCK-5');
            dragTo('TCK-5', sprintZone('Sprint Dos'));

            await waitFor(() => expect(changeTicketSprintAction).toHaveBeenCalledWith(5, 2));
            await waitFor(() => expect(screen.queryByText('TCK-5')).toBeNull());
        });

        it('mantiene el resaltado si el arrastre pasa por un hijo del destino', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })] });
            setup();

            await screen.findByText('TCK-5');
            const zone = backlogZone();
            fireEvent.dragOver(zone, { dataTransfer: makeDataTransfer() });
            const event = createEvent.dragLeave(zone);
            Object.defineProperty(event, 'relatedTarget', { value: rowOf('TCK-5') });
            fireEvent(zone, event);

            expect(backlogZone().className).toContain('border-purple-400');
        });

        it('resalta el destino mientras se arrastra encima', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })] });
            setup();

            await screen.findByText('TCK-5');
            fireEvent.dragOver(sprintZone('Sprint Uno'), { dataTransfer: makeDataTransfer() });
            expect(sprintZone('Sprint Uno').className).toContain('border-purple-400');

            fireEvent.dragLeave(sprintZone('Sprint Uno'));
            expect(sprintZone('Sprint Uno').className).not.toContain('border-purple-400');
        });

        it('marca la fila arrastrada y la libera al terminar', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })] });
            setup();

            await screen.findByText('TCK-5');
            fireEvent.dragStart(rowOf('TCK-5'), { dataTransfer: makeDataTransfer() });
            expect(rowOf('TCK-5').className).toContain('opacity-40');

            fireEvent.dragEnd(rowOf('TCK-5'));
            expect(rowOf('TCK-5').className).not.toContain('opacity-40');
        });
    });

    describe('sin permiso de gestion', () => {
        it('no permite arrastrar ni soltar tickets', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })] });
            setup({ canManage: false });

            await screen.findByText('TCK-5');
            expect(rowOf('TCK-5')).toHaveAttribute('draggable', 'false');

            dragTo('TCK-5', sprintZone('Sprint Uno'));

            expect(changeTicketSprintAction).not.toHaveBeenCalled();
        });

        it('oculta los botones de gestion de sprints', async () => {
            setup({ canManage: false });

            expect(screen.queryByRole('button', { name: '+ Nuevo sprint' })).toBeNull();
            expect(screen.queryByRole('button', { name: 'Editar' })).toBeNull();
            expect(screen.queryByRole('button', { name: 'Cerrar' })).toBeNull();
            await waitFor(() => expect(listTicketsAction).toHaveBeenCalled());
        });
    });

    describe('gestion de sprints', () => {
        it('valida que el nombre y las fechas sean obligatorios', async () => {
            setup();

            fireEvent.click(screen.getByRole('button', { name: '+ Nuevo sprint' }));
            fireEvent.click(screen.getByRole('button', { name: 'Guardar' }));

            expect(await screen.findByText('Nombre, fecha de inicio y fecha de fin son obligatorios.')).toBeInTheDocument();
            expect(createSprintAction).not.toHaveBeenCalled();
        });

        it('crea un sprint con los datos del formulario', async () => {
            const { onSprintsChanged } = setup();

            fireEvent.click(screen.getByRole('button', { name: '+ Nuevo sprint' }));
            fireEvent.change(screen.getByPlaceholderText('Sprint 1'), { target: { value: '  Sprint Nuevo  ' } });
            fireEvent.change(screen.getByLabelText('Fecha de inicio'), { target: { value: '2026-03-01' } });
            fireEvent.change(screen.getByLabelText('Fecha de fin'), { target: { value: '2026-03-15' } });
            fireEvent.change(screen.getByLabelText('Estado'), { target: { value: 'active' } });
            fireEvent.click(screen.getByRole('button', { name: 'Guardar' }));

            await waitFor(() => expect(createSprintAction).toHaveBeenCalledWith({
                name: 'Sprint Nuevo',
                goal: undefined,
                start_date: '2026-03-01',
                end_date: '2026-03-15',
                status: 'active',
            }));
            await waitFor(() => expect(onSprintsChanged).toHaveBeenCalled());
            await waitFor(() => expect(screen.queryByTestId('modal')).toBeNull());
        });

        it('envia el objetivo cuando se escribe', async () => {
            setup();

            fireEvent.click(screen.getByRole('button', { name: '+ Nuevo sprint' }));
            fireEvent.change(screen.getByPlaceholderText('Sprint 1'), { target: { value: 'Sprint Goal' } });
            fireEvent.change(screen.getByLabelText('Objetivo'), { target: { value: 'terminar pagos' } });
            fireEvent.change(screen.getByLabelText('Fecha de inicio'), { target: { value: '2026-03-01' } });
            fireEvent.change(screen.getByLabelText('Fecha de fin'), { target: { value: '2026-03-15' } });
            fireEvent.click(screen.getByRole('button', { name: 'Guardar' }));

            await waitFor(() => expect(createSprintAction).toHaveBeenCalledWith(expect.objectContaining({ goal: 'terminar pagos' })));
        });

        it('precarga el formulario al editar y actualiza el sprint', async () => {
            const { onSprintsChanged } = setup();

            fireEvent.click(screen.getAllByRole('button', { name: 'Editar' })[0]);

            expect((screen.getByPlaceholderText('Sprint 1') as HTMLInputElement).value).toBe('Sprint Uno');
            expect((screen.getByLabelText('Fecha de inicio') as HTMLInputElement).value).toBe('2026-02-01');
            expect((screen.getByLabelText('Fecha de fin') as HTMLInputElement).value).toBe('2026-02-15');

            fireEvent.change(screen.getByPlaceholderText('Sprint 1'), { target: { value: 'Sprint Editado' } });
            fireEvent.click(screen.getByRole('button', { name: 'Guardar' }));

            await waitFor(() => expect(updateSprintAction).toHaveBeenCalledWith(1, expect.objectContaining({ name: 'Sprint Editado', status: 'active' })));
            expect(createSprintAction).not.toHaveBeenCalled();
            await waitFor(() => expect(onSprintsChanged).toHaveBeenCalled());
        });

        it('avisa cuando el guardado falla y deja el modal abierto', async () => {
            vi.mocked(createSprintAction).mockRejectedValue(new Error('409'));
            setup();

            fireEvent.click(screen.getByRole('button', { name: '+ Nuevo sprint' }));
            fireEvent.change(screen.getByPlaceholderText('Sprint 1'), { target: { value: 'Sprint X' } });
            fireEvent.change(screen.getByLabelText('Fecha de inicio'), { target: { value: '2026-03-01' } });
            fireEvent.change(screen.getByLabelText('Fecha de fin'), { target: { value: '2026-03-15' } });
            fireEvent.click(screen.getByRole('button', { name: 'Guardar' }));

            expect(await screen.findByText('No se pudo guardar el sprint. Revisa los datos e intenta de nuevo.')).toBeInTheDocument();
            expect(screen.getByTestId('modal')).toBeInTheDocument();
        });

        it('cierra el modal con el boton cancelar', async () => {
            setup();

            fireEvent.click(screen.getByRole('button', { name: '+ Nuevo sprint' }));
            expect(screen.getByTestId('modal')).toBeInTheDocument();

            fireEvent.click(within(screen.getByTestId('modal')).getByRole('button', { name: 'Cancelar' }));

            expect(screen.queryByTestId('modal')).toBeNull();
        });

        it('cierra el modal desde el propio contenedor', async () => {
            setup();

            fireEvent.click(screen.getByRole('button', { name: '+ Nuevo sprint' }));
            fireEvent.click(screen.getByRole('button', { name: 'modal-close' }));

            expect(screen.queryByTestId('modal')).toBeNull();
        });

        it('activa un sprint planeado y recarga los buckets', async () => {
            mockBuckets({ none: [makeTicket({ id: 5 })] });
            const { onSprintsChanged } = setup();

            await screen.findByText('TCK-5');
            fireEvent.click(screen.getAllByRole('button', { name: 'Activar' })[0]);

            await waitFor(() => expect(changeSprintStatusAction).toHaveBeenCalledWith(2, 'active'));
            await waitFor(() => expect(onSprintsChanged).toHaveBeenCalled());
        });

        it('cierra un sprint activo', async () => {
            const { onSprintsChanged } = setup();

            fireEvent.click(screen.getAllByRole('button', { name: 'Cerrar' })[0]);

            await waitFor(() => expect(changeSprintStatusAction).toHaveBeenCalledWith(1, 'closed'));
            await waitFor(() => expect(onSprintsChanged).toHaveBeenCalled());
        });

        it('avisa cuando falla el cambio de estado del sprint', async () => {
            vi.mocked(changeSprintStatusAction).mockRejectedValue(new Error('500'));
            setup();

            fireEvent.click(screen.getAllByRole('button', { name: 'Cerrar' })[0]);

            expect(await screen.findByText('No se pudo cambiar el estado del sprint. Intenta de nuevo.')).toBeInTheDocument();
        });

        it('no ofrece activar un sprint que ya esta activo ni cerrar uno cerrado', async () => {
            setup({ sprints: [activeSprint, makeSprint({ id: 3, name: 'Sprint Tres', status: 'closed' })] });

            const activar = screen.getAllByRole('button', { name: 'Activar' });
            const cerrar = screen.getAllByRole('button', { name: 'Cerrar' });
            expect(activar).toHaveLength(1);
            expect(cerrar).toHaveLength(1);
            await waitFor(() => expect(listTicketsAction).toHaveBeenCalled());
        });
    });
});
