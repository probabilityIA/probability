import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import TicketForm from './TicketForm';

vi.mock('@/shared/ui', () => ({
    Button: ({ children, ...props }: any) => <button {...props}>{children}</button>,
    Input: (props: any) => <input {...props} />,
}));

vi.mock('@/services/auth/business/ui/hooks/useBusinessesSimple', () => ({
    useBusinessesSimple: vi.fn(() => ({ businesses: [], loading: false, error: null })),
}));

import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';

const sprints = [
    { id: 3, name: 'Sprint 3', goal: '', start_date: '', end_date: '', status: 'active' as const, created_by_id: 1, created_by_name: 'x', ticket_count: 0, done_count: 0, created_at: '', updated_at: '' },
];

const users = [{ id: 11, name: 'Ana' }, { id: 12, name: 'Luis' }];

const makeFile = (name: string, type: string, size?: number) => {
    const file = new File(['contenido'], name, { type });
    if (size !== undefined) Object.defineProperty(file, 'size', { value: size });
    return file;
};

const fileInput = () => document.getElementById('ticket-create-file-input') as HTMLInputElement;

const pickFiles = (files: File[]) => {
    const input = fileInput();
    Object.defineProperty(input, 'files', { value: files, configurable: true });
    fireEvent.change(input);
};

const dropZone = () => document.querySelector('label[for="ticket-create-file-input"]') as HTMLElement;

const combos = () => screen.getAllByRole('combobox');

const fillRequired = (title = '  Falla en checkout  ', description = '  No carga el pago  ') => {
    fireEvent.change(screen.getByPlaceholderText('Resumen breve'), { target: { value: title } });
    fireEvent.change(screen.getByPlaceholderText('Detalla el problema, mejora o solicitud'), { target: { value: description } });
};

const submitForm = () => fireEvent.click(screen.getByRole('button', { name: 'Crear ticket' }));

const setup = (props: Partial<React.ComponentProps<typeof TicketForm>> = {}) => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);
    const onCancel = vi.fn();
    const utils = render(
        <TicketForm
            isSuperAdmin
            users={users}
            sprints={sprints}
            onSubmit={onSubmit}
            onCancel={onCancel}
            {...props}
        />
    );
    return { onSubmit, onCancel, ...utils };
};

describe('TicketForm', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.mocked(useBusinessesSimple).mockReturnValue({ businesses: [], loading: false, error: null } as any);
    });

    describe('validacion', () => {
        it('no envia y muestra error cuando faltan titulo y descripcion', () => {
            const { onSubmit } = setup();

            submitForm();

            expect(screen.getByText('T\u00edtulo y descripci\u00f3n son obligatorios')).toBeInTheDocument();
            expect(onSubmit).not.toHaveBeenCalled();
        });

        it('no envia cuando el titulo es solo espacios', () => {
            const { onSubmit } = setup();

            fillRequired('    ', 'descripcion valida');
            submitForm();

            expect(screen.getByText('T\u00edtulo y descripci\u00f3n son obligatorios')).toBeInTheDocument();
            expect(onSubmit).not.toHaveBeenCalled();
        });

        it('no envia cuando la descripcion es solo espacios', () => {
            const { onSubmit } = setup();

            fillRequired('titulo valido', '   ');
            submitForm();

            expect(onSubmit).not.toHaveBeenCalled();
        });
    });

    describe('envio', () => {
        it('envia el dto con los valores por defecto y los textos recortados', async () => {
            const { onSubmit } = setup();

            fillRequired();
            submitForm();

            await waitFor(() => expect(onSubmit).toHaveBeenCalledTimes(1));
            expect(onSubmit.mock.calls[0][0]).toEqual({
                dto: {
                    title: 'Falla en checkout',
                    description: 'No carga el pago',
                    type: 'support',
                    priority: 'medium',
                    severity: undefined,
                    area: 'soporte',
                    source: 'internal',
                    category: undefined,
                    due_date: undefined,
                    assigned_to_id: null,
                    sprint_id: null,
                    business_id: null,
                },
                files: [],
                comment: '',
                commentInternal: false,
            });
        });

        it('envia los valores elegidos en cada control', async () => {
            const { onSubmit } = setup();

            fillRequired('Bug grave', 'Detalle del bug');
            const [type, priority, severity, area, source, assigned, sprint, business] = combos();
            fireEvent.change(type, { target: { value: 'bug' } });
            fireEvent.change(priority, { target: { value: 'critical' } });
            fireEvent.change(severity, { target: { value: 'high' } });
            fireEvent.change(area, { target: { value: 'desarrollo' } });
            fireEvent.change(source, { target: { value: 'business' } });
            fireEvent.change(assigned, { target: { value: '12' } });
            fireEvent.change(sprint, { target: { value: '3' } });
            expect(business).toBeInTheDocument();
            fireEvent.change(screen.getByPlaceholderText('Ej: env\u00edos, facturaci\u00f3n, frontend'), { target: { value: '  envios  ' } });
            fireEvent.change(screen.getByPlaceholderText('Opcional: contexto adicional para el equipo'), { target: { value: '  contexto  ' } });
            fireEvent.click(screen.getByLabelText('Nota interna (solo super admins)'));
            submitForm();

            await waitFor(() => expect(onSubmit).toHaveBeenCalledTimes(1));
            const payload = onSubmit.mock.calls[0][0];
            expect(payload.dto).toMatchObject({
                type: 'bug',
                priority: 'critical',
                severity: 'high',
                area: 'desarrollo',
                source: 'business',
                assigned_to_id: 12,
                sprint_id: 3,
                category: 'envios',
            });
            expect(payload.comment).toBe('contexto');
            expect(payload.commentInternal).toBe(true);
        });

        it('envia la fecha objetivo cuando se completa', async () => {
            const { onSubmit } = setup();

            fillRequired('Con fecha', 'Descripcion');
            const dateInput = document.querySelector('input[type="date"]') as HTMLInputElement;
            fireEvent.change(dateInput, { target: { value: '2026-03-15' } });
            submitForm();

            await waitFor(() => expect(onSubmit).toHaveBeenCalled());
            expect(onSubmit.mock.calls[0][0].dto.due_date).toBe('2026-03-15');
        });

        it('muestra el mensaje de error cuando onSubmit rechaza', async () => {
            const onSubmit = vi.fn().mockRejectedValue(new Error('El backend no responde'));
            render(<TicketForm isSuperAdmin onSubmit={onSubmit} onCancel={vi.fn()} />);

            fillRequired('Titulo', 'Descripcion');
            submitForm();

            expect(await screen.findByText('El backend no responde')).toBeInTheDocument();
        });

        it('muestra un mensaje generico cuando el rechazo no trae mensaje', async () => {
            const onSubmit = vi.fn().mockRejectedValue({});
            render(<TicketForm isSuperAdmin onSubmit={onSubmit} onCancel={vi.fn()} />);

            fillRequired('Titulo', 'Descripcion');
            submitForm();

            expect(await screen.findByText('Error al crear ticket')).toBeInTheDocument();
        });

        it('llama onCancel al presionar Cancelar', () => {
            const { onCancel, onSubmit } = setup();

            fireEvent.click(screen.getByRole('button', { name: 'Cancelar' }));

            expect(onCancel).toHaveBeenCalledTimes(1);
            expect(onSubmit).not.toHaveBeenCalled();
        });

        it('muestra el progreso recibido y deshabilita los botones mientras envia', () => {
            setup({ submitting: true, submitLabel: 'Subiendo adjuntos 1 de 2...' });

            expect(screen.getByRole('button', { name: 'Subiendo adjuntos 1 de 2...' })).toBeDisabled();
            expect(screen.getByRole('button', { name: 'Cancelar' })).toBeDisabled();
        });

        it('usa "Creando..." cuando envia sin etiqueta de progreso', () => {
            setup({ submitting: true });

            expect(screen.getByRole('button', { name: 'Creando...' })).toBeInTheDocument();
        });
    });

    describe('permisos de super admin', () => {
        it('muestra el selector de negocio con las opciones del hook y lo envia', async () => {
            vi.mocked(useBusinessesSimple).mockReturnValue({
                businesses: [{ id: 26, name: 'Demo' }, { id: 31, name: 'Otra' }],
                loading: false,
                error: null,
            } as any);
            const { onSubmit } = setup();

            const business = combos()[7];
            expect(screen.getByRole('option', { name: 'Demo' })).toBeInTheDocument();
            fireEvent.change(business, { target: { value: '31' } });
            fillRequired('Titulo', 'Descripcion');
            submitForm();

            await waitFor(() => expect(onSubmit).toHaveBeenCalled());
            expect(onSubmit.mock.calls[0][0].dto.business_id).toBe(31);
        });

        it('oculta negocio y nota interna, y fuerza origen negocio, para usuarios no super admin', async () => {
            const { onSubmit } = setup({ isSuperAdmin: false });

            expect(screen.queryByRole('option', { name: 'Interno (sin negocio)' })).toBeNull();
            expect(screen.queryByLabelText('Nota interna (solo super admins)')).toBeNull();
            expect(combos()).toHaveLength(7);

            fillRequired('Titulo', 'Descripcion');
            submitForm();

            await waitFor(() => expect(onSubmit).toHaveBeenCalled());
            expect(onSubmit.mock.calls[0][0].dto.source).toBe('business');
            expect(onSubmit.mock.calls[0][0].dto.business_id).toBeUndefined();
            expect(onSubmit.mock.calls[0][0].commentInternal).toBe(false);
        });
    });

    describe('adjuntos', () => {
        it('agrega los archivos validos con su tamano', () => {
            setup();

            pickFiles([makeFile('captura.png', 'image/png', 2048)]);

            expect(screen.getByText('captura.png')).toBeInTheDocument();
            expect(screen.getByText('2.0 KB')).toBeInTheDocument();
        });

        it('rechaza un archivo de tipo no permitido nombrandolo', () => {
            setup();

            pickFiles([makeFile('script.exe', 'application/x-msdownload', 100)]);

            expect(screen.getByText(/No se puede adjuntar script\.exe/)).toBeInTheDocument();
            expect(screen.getByText(/Solo se admiten im\u00e1genes/)).toBeInTheDocument();
            expect(screen.queryByText('script.exe')).toBeNull();
        });

        it('rechaza un archivo de mas de 10 MB indicando su peso', () => {
            setup();

            pickFiles([makeFile('manual.pdf', 'application/pdf', 11 * 1024 * 1024)]);

            expect(screen.getByText(/No se puede adjuntar manual\.pdf/)).toBeInTheDocument();
            expect(screen.getByText(/Pesa 11\.0 MB y el l\u00edmite es 10 MB/)).toBeInTheDocument();
            expect(screen.queryByText('manual.pdf')).toBeNull();
        });

        it('acepta los validos y rechaza el invalido en la misma seleccion', () => {
            setup();

            pickFiles([
                makeFile('ok.pdf', 'application/pdf', 500),
                makeFile('malo.txt', 'text/plain', 500),
            ]);

            expect(screen.getByText('ok.pdf')).toBeInTheDocument();
            expect(screen.getByText(/No se puede adjuntar malo\.txt/)).toBeInTheDocument();
        });

        it('ignora un archivo repetido con el mismo nombre y tamano', () => {
            setup();

            pickFiles([makeFile('captura.png', 'image/png', 2048)]);
            pickFiles([makeFile('captura.png', 'image/png', 2048)]);

            expect(screen.getAllByText('captura.png')).toHaveLength(1);
        });

        it('agrega un archivo con el mismo nombre pero distinto tamano', () => {
            setup();

            pickFiles([makeFile('captura.png', 'image/png', 2048)]);
            pickFiles([makeFile('captura.png', 'image/png', 4096)]);

            expect(screen.getAllByText('captura.png')).toHaveLength(2);
        });

        it('no cambia la lista cuando la seleccion viene vacia', () => {
            setup();

            pickFiles([makeFile('captura.png', 'image/png', 2048)]);
            pickFiles([]);

            expect(screen.getAllByText('captura.png')).toHaveLength(1);
        });

        it('quita un adjunto de la lista', () => {
            setup();

            pickFiles([makeFile('a.png', 'image/png', 10), makeFile('b.png', 'image/png', 20)]);
            fireEvent.click(screen.getAllByRole('button', { name: 'Quitar' })[0]);

            expect(screen.queryByText('a.png')).toBeNull();
            expect(screen.getByText('b.png')).toBeInTheDocument();
        });

        it('envia los adjuntos aceptados junto al ticket', async () => {
            const { onSubmit } = setup();

            const png = makeFile('uno.png', 'image/png', 10);
            const pdf = makeFile('dos.pdf', 'application/pdf', 20);
            pickFiles([png, pdf]);
            fillRequired('Titulo', 'Descripcion');
            submitForm();

            await waitFor(() => expect(onSubmit).toHaveBeenCalled());
            expect(onSubmit.mock.calls[0][0].files).toEqual([png, pdf]);
        });

        it('acepta archivos soltados sobre la zona de arrastre', () => {
            setup();

            fireEvent.drop(dropZone(), { dataTransfer: { files: [makeFile('soltado.png', 'image/png', 30)] } });

            expect(screen.getByText('soltado.png')).toBeInTheDocument();
        });

        it('cambia el texto de la zona mientras se arrastra encima y al salir', () => {
            setup();

            fireEvent.dragOver(dropZone());
            expect(screen.getByText('Suelta para adjuntar')).toBeInTheDocument();

            fireEvent.dragLeave(dropZone());
            expect(screen.getByText('Arrastra archivos aqu\u00ed')).toBeInTheDocument();
        });

        it('limpia el aviso de rechazo al elegir un archivo valido despues', () => {
            setup();

            pickFiles([makeFile('malo.txt', 'text/plain', 10)]);
            expect(screen.getByText(/No se puede adjuntar malo\.txt/)).toBeInTheDocument();

            pickFiles([makeFile('bueno.png', 'image/png', 10)]);

            expect(screen.queryByText(/No se puede adjuntar malo\.txt/)).toBeNull();
        });
    });
});
