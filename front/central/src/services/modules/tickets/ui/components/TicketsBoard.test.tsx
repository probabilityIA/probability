import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, createEvent, within } from '@testing-library/react';
import TicketsBoard, { BOARD_COLUMNS } from './TicketsBoard';
import { Ticket, TicketStatus } from '../../domain/types';

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

const makeDataTransfer = () => {
    const store: Record<string, string> = {};
    return {
        effectAllowed: '',
        dropEffect: '',
        setData: (key: string, value: string) => { store[key] = value; },
        getData: (key: string) => store[key] ?? '',
    };
};

const columnOf = (title: string): HTMLElement =>
    screen.getByRole('heading', { name: title }).parentElement!.parentElement as HTMLElement;

const cardOf = (code: string): HTMLElement =>
    screen.getByText(code).closest('[draggable]') as HTMLElement;

const dragCardTo = (code: string, toColumn: string) => {
    const dataTransfer = makeDataTransfer();
    fireEvent.dragStart(cardOf(code), { dataTransfer });
    const target = columnOf(toColumn);
    fireEvent.dragOver(target, { dataTransfer });
    fireEvent.drop(target, { dataTransfer });
    return dataTransfer;
};

const dragLeaveTo = (target: HTMLElement, relatedTarget: EventTarget | null) => {
    const event = createEvent.dragLeave(target);
    Object.defineProperty(event, 'relatedTarget', { value: relatedTarget });
    fireEvent(target, event);
};

const baseProps = {
    loading: false,
    canDrag: true,
    updatingId: null as number | null,
    onOpen: vi.fn(),
    onMove: vi.fn(),
};

describe('TicketsBoard', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe('render de columnas', () => {
        it('agrupa cada ticket en la columna que corresponde a su estado', () => {
            const tickets = [
                makeTicket({ id: 1, status: 'open' }),
                makeTicket({ id: 2, status: 'in_review' }),
                makeTicket({ id: 3, status: 'in_development' }),
                makeTicket({ id: 4, status: 'testing' }),
                makeTicket({ id: 5, status: 'blocked' }),
                makeTicket({ id: 6, status: 'closed' }),
            ];
            render(<TicketsBoard {...baseProps} tickets={tickets} />);

            expect(within(columnOf('To Do')).getByText('TCK-1')).toBeInTheDocument();
            expect(within(columnOf('To Do')).getByText('TCK-2')).toBeInTheDocument();
            expect(within(columnOf('In Progress')).getByText('TCK-3')).toBeInTheDocument();
            expect(within(columnOf('In Testing')).getByText('TCK-4')).toBeInTheDocument();
            expect(within(columnOf('Blocked')).getByText('TCK-5')).toBeInTheDocument();
            expect(within(columnOf('Done')).getByText('TCK-6')).toBeInTheDocument();
        });

        it('muestra el contador de tickets por columna', () => {
            const tickets = [
                makeTicket({ id: 1, status: 'open' }),
                makeTicket({ id: 2, status: 'in_review' }),
                makeTicket({ id: 3, status: 'in_development' }),
            ];
            render(<TicketsBoard {...baseProps} tickets={tickets} />);

            expect(within(columnOf('To Do')).getByText('2')).toBeInTheDocument();
            expect(within(columnOf('In Progress')).getByText('1')).toBeInTheDocument();
            expect(within(columnOf('Blocked')).getByText('0')).toBeInTheDocument();
        });

        it('muestra "Sin tickets" solo en las columnas vacias', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 1, status: 'open' })]} />);

            expect(within(columnOf('To Do')).queryByText('Sin tickets')).toBeNull();
            expect(within(columnOf('Blocked')).getByText('Sin tickets')).toBeInTheDocument();
            expect(screen.getAllByText('Sin tickets')).toHaveLength(BOARD_COLUMNS.length - 1);
        });

        it('muestra el badge de estado solo en columnas que agrupan varios estados', () => {
            const tickets = [
                makeTicket({ id: 1, status: 'in_review' }),
                makeTicket({ id: 2, status: 'in_development' }),
            ];
            render(<TicketsBoard {...baseProps} tickets={tickets} />);

            expect(within(columnOf('To Do')).getByText('En revisi\u00f3n')).toBeInTheDocument();
            expect(within(columnOf('In Progress')).queryByText('En desarrollo')).toBeNull();
        });

        it('muestra el overlay de carga cuando loading es true', () => {
            render(<TicketsBoard {...baseProps} loading tickets={[]} />);

            expect(screen.getByText('Actualizando...')).toBeInTheDocument();
        });

        it('no muestra el overlay de carga cuando loading es false', () => {
            render(<TicketsBoard {...baseProps} tickets={[]} />);

            expect(screen.queryByText('Actualizando...')).toBeNull();
        });
    });

    describe('datos de la tarjeta', () => {
        it('muestra el avatar cuando getAvatarUrl devuelve una url', () => {
            render(
                <TicketsBoard
                    {...baseProps}
                    tickets={[makeTicket({ id: 1, assigned_to_name: 'Ana Lopez' })]}
                    getAvatarUrl={() => 'https://cdn.test/ana.png'}
                />
            );

            const card = cardOf('TCK-1');
            expect(card.querySelector('img')).toHaveAttribute('src', 'https://cdn.test/ana.png');
            expect(within(card).getByText('Ana Lopez')).toBeInTheDocument();
        });

        it('muestra la inicial del asignado cuando no hay avatar', () => {
            render(
                <TicketsBoard
                    {...baseProps}
                    tickets={[makeTicket({ id: 1, assigned_to_name: 'ana lopez' })]}
                    getAvatarUrl={() => ''}
                />
            );

            const card = cardOf('TCK-1');
            expect(card.querySelector('img')).toBeNull();
            expect(within(card).getByText('A')).toBeInTheDocument();
        });

        it('muestra "Sin asignar" y un guion cuando el ticket no tiene responsable', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 1 })]} />);

            const card = cardOf('TCK-1');
            expect(within(card).getByText('Sin asignar')).toBeInTheDocument();
            expect(within(card).getByText('-')).toBeInTheDocument();
        });

        it('muestra el negocio, el id del negocio o "Interno" segun el ticket', () => {
            render(
                <TicketsBoard
                    {...baseProps}
                    tickets={[
                        makeTicket({ id: 1, business_name: 'Tienda Uno' }),
                        makeTicket({ id: 2, business_id: 42 }),
                        makeTicket({ id: 3 }),
                    ]}
                />
            );

            expect(within(cardOf('TCK-1')).getByText('Tienda Uno')).toBeInTheDocument();
            expect(within(cardOf('TCK-2')).getByText('#42')).toBeInTheDocument();
            expect(within(cardOf('TCK-3')).getByText('Interno')).toBeInTheDocument();
        });

        it('resalta la tarjeta que se esta actualizando', () => {
            render(
                <TicketsBoard
                    {...baseProps}
                    updatingId={2}
                    tickets={[makeTicket({ id: 1 }), makeTicket({ id: 2 })]}
                />
            );

            expect(cardOf('TCK-2').className).toContain('animate-pulse');
            expect(cardOf('TCK-1').className).not.toContain('animate-pulse');
        });

        it('llama onOpen con el ticket al hacer click en la tarjeta', () => {
            const ticket = makeTicket({ id: 7 });
            render(<TicketsBoard {...baseProps} tickets={[ticket]} />);

            fireEvent.click(cardOf('TCK-7'));

            expect(baseProps.onOpen).toHaveBeenCalledWith(ticket);
        });
    });

    describe('drag and drop', () => {
        it('mueve la tarjeta a otra columna con el estado canonico de destino', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 5, status: 'open' })]} />);

            dragCardTo('TCK-5', 'In Progress');

            expect(baseProps.onMove).toHaveBeenCalledWith(5, 'in_development');
        });

        it('usa el estado de drop de cada columna, no el estado visible del ticket', () => {
            const cases: Array<[TicketStatus, string, TicketStatus]> = [
                ['blocked', 'To Do', 'open'],
                ['open', 'In Progress', 'in_development'],
                ['open', 'In Testing', 'testing'],
                ['open', 'Blocked', 'blocked'],
                ['open', 'Done', 'resolved'],
            ];
            cases.forEach(([from, title, expected]) => {
                vi.clearAllMocks();
                const { unmount } = render(
                    <TicketsBoard {...baseProps} tickets={[makeTicket({ id: 9, status: from })]} />
                );

                dragCardTo('TCK-9', title);

                expect(baseProps.onMove).toHaveBeenCalledWith(9, expected);
                unmount();
            });
        });

        it('no llama onMove al soltar en la misma columna', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 3, status: 'in_review' })]} />);

            dragCardTo('TCK-3', 'To Do');

            expect(baseProps.onMove).not.toHaveBeenCalled();
        });

        it('no llama onMove si el id soltado no corresponde a ningun ticket', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 1, status: 'open' })]} />);

            const dataTransfer = makeDataTransfer();
            dataTransfer.setData('text/plain', '999');
            fireEvent.drop(columnOf('Done'), { dataTransfer });

            expect(baseProps.onMove).not.toHaveBeenCalled();
        });

        it('no llama onMove cuando no hay id ni en el dataTransfer ni en el arrastre en curso', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 1, status: 'open' })]} />);

            fireEvent.drop(columnOf('Done'), { dataTransfer: makeDataTransfer() });

            expect(baseProps.onMove).not.toHaveBeenCalled();
        });

        it('usa el id del arrastre en curso cuando el dataTransfer llega vacio', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 4, status: 'open' })]} />);

            fireEvent.dragStart(cardOf('TCK-4'), { dataTransfer: makeDataTransfer() });
            fireEvent.drop(columnOf('Blocked'), { dataTransfer: makeDataTransfer() });

            expect(baseProps.onMove).toHaveBeenCalledWith(4, 'blocked');
        });

        it('guarda el id del ticket en el dataTransfer al iniciar el arrastre', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 8, status: 'open' })]} />);

            const dataTransfer = makeDataTransfer();
            fireEvent.dragStart(cardOf('TCK-8'), { dataTransfer });

            expect(dataTransfer.getData('text/plain')).toBe('8');
            expect(dataTransfer.effectAllowed).toBe('move');
        });

        it('resalta la columna destino mientras se arrastra encima y lo quita al salir', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 1, status: 'open' })]} />);

            const target = columnOf('Done');
            expect(target.className).not.toContain('border-purple-400');

            fireEvent.dragOver(target, { dataTransfer: makeDataTransfer() });
            expect(columnOf('Done').className).toContain('border-purple-400');
            expect(columnOf('Blocked').className).not.toContain('border-purple-400');

            dragLeaveTo(columnOf('Done'), document.body);
            expect(columnOf('Done').className).not.toContain('border-purple-400');
        });

        it('mantiene el resaltado si el dragLeave ocurre hacia un hijo de la columna', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 1, status: 'open' })]} />);

            const target = columnOf('To Do');
            fireEvent.dragOver(target, { dataTransfer: makeDataTransfer() });
            dragLeaveTo(target, cardOf('TCK-1'));

            expect(columnOf('To Do').className).toContain('border-purple-400');
        });

        it('quita el resaltado y limpia la opacidad al terminar el arrastre', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 1, status: 'open' })]} />);

            fireEvent.dragStart(cardOf('TCK-1'), { dataTransfer: makeDataTransfer() });
            expect(cardOf('TCK-1').className).toContain('opacity-40');
            fireEvent.dragOver(columnOf('Done'), { dataTransfer: makeDataTransfer() });

            fireEvent.dragEnd(cardOf('TCK-1'));

            expect(cardOf('TCK-1').className).not.toContain('opacity-40');
            expect(columnOf('Done').className).not.toContain('border-purple-400');
        });

        it('quita el resaltado de la columna despues de soltar', () => {
            render(<TicketsBoard {...baseProps} tickets={[makeTicket({ id: 1, status: 'open' })]} />);

            dragCardTo('TCK-1', 'Done');

            expect(columnOf('Done').className).not.toContain('border-purple-400');
        });
    });

    describe('sin permiso para arrastrar', () => {
        it('marca las tarjetas como no arrastrables', () => {
            render(<TicketsBoard {...baseProps} canDrag={false} tickets={[makeTicket({ id: 1 })]} />);

            expect(cardOf('TCK-1')).toHaveAttribute('draggable', 'false');
            expect(cardOf('TCK-1').className).toContain('cursor-pointer');
        });

        it('no llama onMove al soltar en otra columna', () => {
            render(<TicketsBoard {...baseProps} canDrag={false} tickets={[makeTicket({ id: 1, status: 'open' })]} />);

            dragCardTo('TCK-1', 'Done');

            expect(baseProps.onMove).not.toHaveBeenCalled();
        });

        it('no resalta la columna al arrastrar encima', () => {
            render(<TicketsBoard {...baseProps} canDrag={false} tickets={[makeTicket({ id: 1, status: 'open' })]} />);

            fireEvent.dragOver(columnOf('Done'), { dataTransfer: makeDataTransfer() });

            expect(columnOf('Done').className).not.toContain('border-purple-400');
        });

        it('no marca la tarjeta como arrastrandose al iniciar el drag', () => {
            render(<TicketsBoard {...baseProps} canDrag={false} tickets={[makeTicket({ id: 1, status: 'open' })]} />);

            fireEvent.dragStart(cardOf('TCK-1'), { dataTransfer: makeDataTransfer() });

            expect(cardOf('TCK-1').className).not.toContain('opacity-40');
        });
    });
});
