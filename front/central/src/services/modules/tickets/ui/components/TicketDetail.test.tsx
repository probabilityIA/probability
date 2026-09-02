import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import TicketDetail from './TicketDetail';
import { Ticket, TicketAttachment, TicketComment, TicketHistoryEntry } from '../../domain/types';

vi.mock('@/shared/ui', () => ({
    Button: ({ children, ...props }: any) => <button {...props}>{children}</button>,
}));

vi.mock('../../infra/actions', () => ({
    listCommentsAction: vi.fn(),
    addCommentAction: vi.fn(),
    listAttachmentsAction: vi.fn(),
    uploadAttachmentAction: vi.fn(),
    deleteAttachmentAction: vi.fn(),
    listTicketHistoryAction: vi.fn(),
    changeTicketStatusAction: vi.fn(),
    escalateTicketAction: vi.fn(),
    deleteTicketAction: vi.fn(),
    assignTicketAction: vi.fn(),
}));

vi.mock('@/services/auth/users/infra/actions', () => ({
    getUsersAction: vi.fn(),
}));

import {
    listCommentsAction,
    addCommentAction,
    listAttachmentsAction,
    uploadAttachmentAction,
    deleteAttachmentAction,
    listTicketHistoryAction,
    changeTicketStatusAction,
    escalateTicketAction,
    deleteTicketAction,
    assignTicketAction,
} from '../../infra/actions';
import { getUsersAction } from '@/services/auth/users/infra/actions';

const makeTicket = (over: Partial<Ticket> = {}): Ticket => ({
    id: 55,
    code: 'TCK-55',
    created_by_id: 9,
    created_by_name: 'Carlos Ruiz',
    title: 'No carga el reporte',
    description: 'Al abrir el reporte sale vacio',
    type: 'bug',
    priority: 'high',
    status: 'open',
    source: 'internal',
    escalated_to_dev: false,
    created_at: '2026-01-01T10:00:00Z',
    updated_at: '2026-01-05T10:00:00Z',
    comments_count: 0,
    attachments_count: 0,
    ...over,
});

const makeComment = (over: Partial<TicketComment> & { id: number }): TicketComment => ({
    ticket_id: 55,
    user_id: 3,
    user_name: 'Ana Lopez',
    body: 'Revisando el caso',
    is_internal: false,
    created_at: '2026-01-02T10:00:00Z',
    ...over,
});

const makeAttachment = (over: Partial<TicketAttachment> & { id: number }): TicketAttachment => ({
    ticket_id: 55,
    uploaded_by_id: 3,
    uploaded_by_name: 'Ana',
    file_url: 'https://cdn.test/file-' + over.id,
    file_name: 'archivo-' + over.id + '.png',
    mime_type: 'image/png',
    size: 2048,
    created_at: '2026-01-02T10:00:00Z',
    ...over,
});

const makeHistory = (over: Partial<TicketHistoryEntry> & { id: number }): TicketHistoryEntry => ({
    ticket_id: 55,
    from_status: 'open',
    to_status: 'resolved',
    changed_by_id: 3,
    changed_by_name: 'Ana Lopez',
    note: '',
    created_at: '2026-01-03T10:00:00Z',
    ...over,
});

const makeFile = (name: string, type: string, size?: number) => {
    const file = new File(['data'], name, { type });
    if (size !== undefined) Object.defineProperty(file, 'size', { value: size });
    return file;
};

const uploadInput = () => document.getElementById('ticket-file-input') as HTMLInputElement;

const uploadFiles = (files: File[]) => {
    const input = uploadInput();
    Object.defineProperty(input, 'files', { value: files, configurable: true });
    fireEvent.change(input);
};

const dropZone = () => document.querySelector('label[for="ticket-file-input"]') as HTMLElement;

const setup = async (
    over: { ticket?: Partial<Ticket>; isSuperAdmin?: boolean } = {}
) => {
    const onClose = vi.fn();
    const onChanged = vi.fn();
    const ticket = makeTicket(over.ticket);
    const utils = render(
        <TicketDetail
            ticket={ticket}
            isSuperAdmin={over.isSuperAdmin ?? true}
            onClose={onClose}
            onChanged={onChanged}
        />
    );
    await waitFor(() => expect(listCommentsAction).toHaveBeenCalled());
    return { onClose, onChanged, ticket, ...utils };
};

describe('TicketDetail', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.mocked(listCommentsAction).mockResolvedValue([]);
        vi.mocked(listAttachmentsAction).mockResolvedValue([]);
        vi.mocked(listTicketHistoryAction).mockResolvedValue([]);
        vi.mocked(getUsersAction).mockResolvedValue({ data: [] } as any);
        vi.mocked(addCommentAction).mockResolvedValue({} as any);
        vi.mocked(uploadAttachmentAction).mockResolvedValue({} as any);
        vi.mocked(changeTicketStatusAction).mockResolvedValue({} as any);
        vi.mocked(escalateTicketAction).mockResolvedValue({} as any);
        vi.mocked(assignTicketAction).mockResolvedValue({} as any);
        vi.mocked(deleteTicketAction).mockResolvedValue(undefined as any);
        vi.mocked(deleteAttachmentAction).mockResolvedValue(undefined as any);
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('carga inicial', () => {
        it('pide comentarios, adjuntos e historial del ticket abierto', async () => {
            await setup();

            expect(listCommentsAction).toHaveBeenCalledWith(55);
            expect(listAttachmentsAction).toHaveBeenCalledWith(55);
            expect(listTicketHistoryAction).toHaveBeenCalledWith(55);
        });

        it('muestra la cabecera con codigo, titulo, autor y descripcion', async () => {
            await setup();

            expect(screen.getByText('TCK-55')).toBeInTheDocument();
            expect(screen.getByText('No carga el reporte')).toBeInTheDocument();
            expect(screen.getByText('Creado por Carlos Ruiz')).toBeInTheDocument();
            expect(screen.getByText('Al abrir el reporte sale vacio')).toBeInTheDocument();
        });

        it('usa el id del autor cuando no llega su nombre', async () => {
            await setup({ ticket: { created_by_name: undefined } });

            expect(screen.getByText('Creado por #9')).toBeInTheDocument();
        });

        it('muestra los mensajes de vacio cuando no hay comentarios ni historial', async () => {
            await setup();

            expect(await screen.findByText('Sin comentarios todav\u00eda')).toBeInTheDocument();
            expect(screen.getByText('Sin movimientos')).toBeInTheDocument();
        });

        it('pinta los comentarios con autor, iniciales y marca de nota interna', async () => {
            vi.mocked(listCommentsAction).mockResolvedValue([
                makeComment({ id: 1 }),
                makeComment({ id: 2, is_internal: true, body: 'Solo equipo', user_name: '' }),
            ]);
            await setup();

            expect(await screen.findByText('Revisando el caso')).toBeInTheDocument();
            expect(screen.getByText('Ana Lopez')).toBeInTheDocument();
            expect(screen.getByText('AL')).toBeInTheDocument();
            expect(screen.getByText('Nota interna')).toBeInTheDocument();
            expect(screen.getByText('Usuario 3')).toBeInTheDocument();
            expect(screen.getByText('?')).toBeInTheDocument();
        });

        it('pinta los adjuntos distinguiendo imagen y pdf', async () => {
            vi.mocked(listAttachmentsAction).mockResolvedValue([
                makeAttachment({ id: 1 }),
                makeAttachment({ id: 2, file_name: 'manual.pdf', mime_type: 'application/pdf', size: 1048576 }),
            ]);
            await setup();

            expect(await screen.findByText('archivo-1.png')).toBeInTheDocument();
            expect(screen.getByAltText('archivo-1.png')).toHaveAttribute('src', 'https://cdn.test/file-1');
            expect(screen.getByText('manual.pdf')).toBeInTheDocument();
            expect(screen.getByText('PDF')).toBeInTheDocument();
            expect(screen.getByText('2.0 KB')).toBeInTheDocument();
            expect(screen.getByText('1.0 MB')).toBeInTheDocument();
        });

        it('pinta el historial con la transicion de estados y la nota', async () => {
            vi.mocked(listTicketHistoryAction).mockResolvedValue([
                makeHistory({ id: 1, note: 'listo' }),
                makeHistory({ id: 2, from_status: 'open', to_status: 'open', changed_by_name: '' }),
            ]);
            await setup();

            expect(await screen.findByText(/Abierto -> Resuelto/)).toBeInTheDocument();
            expect(screen.getByText(/\(listo\)/)).toBeInTheDocument();
            expect(screen.getByText(/actualiz\u00f3 el ticket/)).toBeInTheDocument();
            expect(screen.getByText('Usuario 3')).toBeInTheDocument();
        });

        it('describe el movimiento como actualizacion cuando no hay estado de origen', async () => {
            vi.mocked(listTicketHistoryAction).mockResolvedValue([
                makeHistory({ id: 1, from_status: '', to_status: 'resolved' }),
            ]);
            await setup();

            expect(await screen.findByText(/actualiz\u00f3 el ticket/)).toBeInTheDocument();
            expect(screen.queryByText(/->/)).toBeNull();
        });

        it('muestra la marca de escalado cuando el ticket ya fue escalado', async () => {
            await setup({ ticket: { escalated_to_dev: true } });

            expect(screen.getByText('Escalado a dev')).toBeInTheDocument();
            expect(screen.queryByRole('button', { name: 'Escalar a dev' })).toBeNull();
        });

        it('muestra la categoria solo cuando el ticket la tiene', async () => {
            const { unmount } = await setup({ ticket: { category: 'facturacion' } });
            expect(screen.getByText(/Categor\u00eda: facturacion/)).toBeInTheDocument();
            unmount();

            await setup();
            expect(screen.queryByText(/Categor\u00eda:/)).toBeNull();
        });
    });

    describe('comentarios', () => {
        it('mantiene el boton deshabilitado mientras el comentario esta vacio', async () => {
            await setup();

            expect(screen.getByRole('button', { name: 'Comentar' })).toBeDisabled();
        });

        it('publica el comentario recortado y refresca el detalle', async () => {
            const { onChanged } = await setup();

            fireEvent.change(screen.getByPlaceholderText('Escribe un comentario...'), { target: { value: '  hola equipo  ' } });
            fireEvent.click(screen.getByRole('button', { name: 'Comentar' }));

            await waitFor(() => expect(addCommentAction).toHaveBeenCalledWith(55, 'hola equipo', false));
            await waitFor(() => expect(onChanged).toHaveBeenCalled());
            expect(listCommentsAction).toHaveBeenCalledTimes(2);
            expect((screen.getByPlaceholderText('Escribe un comentario...') as HTMLTextAreaElement).value).toBe('');
        });

        it('publica como nota interna cuando el super admin marca la casilla', async () => {
            await setup();

            fireEvent.change(screen.getByPlaceholderText('Escribe un comentario...'), { target: { value: 'solo equipo' } });
            fireEvent.click(screen.getByLabelText('Nota interna (solo super admins)'));
            fireEvent.click(screen.getByRole('button', { name: 'Comentar' }));

            await waitFor(() => expect(addCommentAction).toHaveBeenCalledWith(55, 'solo equipo', true));
        });

        it('nunca marca la nota como interna para un usuario que no es super admin', async () => {
            await setup({ isSuperAdmin: false });

            expect(screen.queryByLabelText('Nota interna (solo super admins)')).toBeNull();
            fireEvent.change(screen.getByPlaceholderText('Escribe un comentario...'), { target: { value: 'consulta' } });
            fireEvent.click(screen.getByRole('button', { name: 'Comentar' }));

            await waitFor(() => expect(addCommentAction).toHaveBeenCalledWith(55, 'consulta', false));
        });
    });

    describe('adjuntos', () => {
        it('sube el archivo elegido y refresca', async () => {
            const { onChanged } = await setup();

            uploadFiles([makeFile('captura.png', 'image/png', 1024)]);

            await waitFor(() => expect(uploadAttachmentAction).toHaveBeenCalledTimes(1));
            const [id, formData] = vi.mocked(uploadAttachmentAction).mock.calls[0];
            expect(id).toBe(55);
            expect((formData.get('file') as File).name).toBe('captura.png');
            await waitFor(() => expect(onChanged).toHaveBeenCalled());
        });

        it('rechaza un tipo no permitido sin llamar a la accion', async () => {
            await setup();

            uploadFiles([makeFile('script.exe', 'application/x-msdownload', 100)]);

            expect(await screen.findByText(/No se puede adjuntar script\.exe/)).toBeInTheDocument();
            expect(uploadAttachmentAction).not.toHaveBeenCalled();
        });

        it('rechaza un archivo de mas de 10 MB indicando su peso', async () => {
            await setup();

            uploadFiles([makeFile('grande.pdf', 'application/pdf', 11 * 1024 * 1024)]);

            expect(await screen.findByText(/Pesa 11\.0 MB y el l\u00edmite es 10 MB/)).toBeInTheDocument();
            expect(uploadAttachmentAction).not.toHaveBeenCalled();
        });

        it('avisa cuando la subida falla', async () => {
            vi.mocked(uploadAttachmentAction).mockRejectedValue(new Error('boom'));
            await setup();

            uploadFiles([makeFile('captura.png', 'image/png', 1024)]);

            expect(await screen.findByText('No se pudo subir captura.png. Intenta de nuevo.')).toBeInTheDocument();
        });

        it('permite cerrar el aviso de error de subida', async () => {
            await setup();

            uploadFiles([makeFile('script.exe', 'application/x-msdownload', 100)]);
            const alert = await screen.findByText(/No se puede adjuntar script\.exe/);
            fireEvent.click(screen.getByRole('button', { name: 'x' }));

            await waitFor(() => expect(alert).not.toBeInTheDocument());
        });

        it('sube el archivo soltado sobre la zona de arrastre', async () => {
            await setup();

            fireEvent.drop(dropZone(), { dataTransfer: { files: [makeFile('soltado.png', 'image/png', 512)] } });

            await waitFor(() => expect(uploadAttachmentAction).toHaveBeenCalledTimes(1));
        });

        it('no sube nada si el drop llega sin archivos', async () => {
            await setup();

            fireEvent.drop(dropZone(), { dataTransfer: { files: [] } });

            expect(uploadAttachmentAction).not.toHaveBeenCalled();
        });

        it('no sube nada cuando el input se dispara sin archivos', async () => {
            await setup();

            uploadFiles([]);

            expect(uploadAttachmentAction).not.toHaveBeenCalled();
        });

        it('muestra un icono generico para adjuntos que no son imagen ni pdf', async () => {
            vi.mocked(listAttachmentsAction).mockResolvedValue([
                makeAttachment({ id: 3, file_name: 'notas.txt', mime_type: 'text/plain' }),
            ]);
            await setup();

            expect(await screen.findByText('notas.txt')).toBeInTheDocument();
            expect(screen.queryByAltText('notas.txt')).toBeNull();
            expect(screen.queryByText('PDF')).toBeNull();
        });

        it('cambia el texto de la zona al arrastrar encima y al salir', async () => {
            await setup();

            fireEvent.dragOver(dropZone());
            expect(screen.getByText('Suelta para adjuntar')).toBeInTheDocument();

            fireEvent.dragLeave(dropZone());
            expect(screen.getByText('Arrastra un archivo aqu\u00ed')).toBeInTheDocument();
        });

        it('elimina un adjunto tras confirmar', async () => {
            vi.spyOn(window, 'confirm').mockReturnValue(true);
            vi.mocked(listAttachmentsAction).mockResolvedValue([makeAttachment({ id: 7 })]);
            await setup();

            fireEvent.click(await screen.findByRole('button', { name: 'Eliminar' }));

            await waitFor(() => expect(deleteAttachmentAction).toHaveBeenCalledWith(7));
            expect(listAttachmentsAction).toHaveBeenCalledTimes(2);
        });

        it('no elimina el adjunto si se cancela la confirmacion', async () => {
            vi.spyOn(window, 'confirm').mockReturnValue(false);
            vi.mocked(listAttachmentsAction).mockResolvedValue([makeAttachment({ id: 7 })]);
            await setup();

            fireEvent.click(await screen.findByRole('button', { name: 'Eliminar' }));

            expect(deleteAttachmentAction).not.toHaveBeenCalled();
        });

        it('no ofrece eliminar adjuntos a un usuario que no es super admin', async () => {
            vi.mocked(listAttachmentsAction).mockResolvedValue([makeAttachment({ id: 7 })]);
            await setup({ isSuperAdmin: false });

            expect(await screen.findByText('archivo-7.png')).toBeInTheDocument();
            expect(screen.queryByRole('button', { name: 'Eliminar' })).toBeNull();
        });
    });

    describe('cambio de estado', () => {
        it('cambia a otro estado enviando la nota y luego la limpia', async () => {
            const { onChanged } = await setup();

            fireEvent.change(screen.getByPlaceholderText('Nota opcional al cambiar de estado'), { target: { value: 'ya quedo' } });
            fireEvent.click(screen.getByRole('button', { name: 'Resuelto' }));

            await waitFor(() => expect(changeTicketStatusAction).toHaveBeenCalledWith(55, 'resolved', 'ya quedo'));
            await waitFor(() => expect(onChanged).toHaveBeenCalled());
            expect((screen.getByPlaceholderText('Nota opcional al cambiar de estado') as HTMLInputElement).value).toBe('');
        });

        it('no llama a la accion al pulsar el estado actual', async () => {
            await setup();

            fireEvent.click(screen.getByRole('button', { name: 'Abierto' }));

            expect(changeTicketStatusAction).not.toHaveBeenCalled();
        });
    });

    describe('acciones de super admin', () => {
        it('escala el ticket a desarrollo', async () => {
            const { onChanged } = await setup();

            fireEvent.click(screen.getByRole('button', { name: 'Escalar a dev' }));

            await waitFor(() => expect(escalateTicketAction).toHaveBeenCalledWith(55, 'Escalado a desarrollo'));
            await waitFor(() => expect(onChanged).toHaveBeenCalled());
        });

        it('elimina el ticket y cierra el detalle tras confirmar', async () => {
            vi.spyOn(window, 'confirm').mockReturnValue(true);
            const { onClose, onChanged } = await setup();

            fireEvent.click(screen.getByTitle('Eliminar ticket'));

            await waitFor(() => expect(deleteTicketAction).toHaveBeenCalledWith(55));
            expect(onClose).toHaveBeenCalledTimes(1);
            expect(onChanged).toHaveBeenCalledTimes(1);
        });

        it('no elimina el ticket si se cancela la confirmacion', async () => {
            vi.spyOn(window, 'confirm').mockReturnValue(false);
            const { onClose } = await setup();

            fireEvent.click(screen.getByTitle('Eliminar ticket'));

            expect(deleteTicketAction).not.toHaveBeenCalled();
            expect(onClose).not.toHaveBeenCalled();
        });

        it('carga solo usuarios de plataforma o super usuarios en el selector de asignado', async () => {
            vi.mocked(getUsersAction).mockResolvedValue({
                data: [
                    { id: 1, name: 'Plataforma', email: 'p@test.com', scope_code: 'platform' },
                    { id: 2, name: 'Super', email: 's@test.com', is_super_user: true },
                    { id: 3, name: 'Cliente', email: 'c@test.com', scope_code: 'business' },
                    { id: 4, name: '', email: 'sin@test.com', scope_code: 'platform' },
                ],
            } as any);
            await setup();

            expect(await screen.findByRole('option', { name: 'Plataforma (p@test.com)' })).toBeInTheDocument();
            expect(screen.getByRole('option', { name: 'Super (s@test.com)' })).toBeInTheDocument();
            expect(screen.queryByRole('option', { name: /Cliente/ })).toBeNull();
            expect(screen.queryByRole('option', { name: /sin@test.com/ })).toBeNull();
        });

        it('asigna el ticket al usuario elegido', async () => {
            vi.mocked(getUsersAction).mockResolvedValue({
                data: [{ id: 1, name: 'Plataforma', email: 'p@test.com', scope_code: 'platform' }],
            } as any);
            const { onChanged } = await setup();

            await screen.findByRole('option', { name: 'Plataforma (p@test.com)' });
            fireEvent.change(screen.getByRole('combobox'), { target: { value: '1' } });

            await waitFor(() => expect(assignTicketAction).toHaveBeenCalledWith(55, 1));
            await waitFor(() => expect(onChanged).toHaveBeenCalled());
        });

        it('desasigna el ticket cuando se elige la opcion vacia', async () => {
            await setup({ ticket: { assigned_to_id: 4 } });

            fireEvent.change(screen.getByRole('combobox'), { target: { value: '' } });

            await waitFor(() => expect(assignTicketAction).toHaveBeenCalledWith(55, null));
        });

        it('oculta asignacion, escalado y borrado a un usuario que no es super admin', async () => {
            await setup({ isSuperAdmin: false });

            expect(screen.queryByRole('combobox')).toBeNull();
            expect(screen.queryByRole('button', { name: 'Escalar a dev' })).toBeNull();
            expect(screen.queryByTitle('Eliminar ticket')).toBeNull();
            expect(getUsersAction).not.toHaveBeenCalled();
        });
    });

    describe('enlace y cierre', () => {
        it('copia el enlace del ticket y avisa que se copio', async () => {
            const writeText = vi.fn().mockResolvedValue(undefined);
            Object.defineProperty(navigator, 'clipboard', { value: { writeText }, configurable: true });
            await setup();

            fireEvent.click(screen.getByTitle('Copiar enlace a este ticket'));

            await waitFor(() => expect(writeText).toHaveBeenCalledWith(expect.stringContaining('?ticket=55')));
            expect(await screen.findByText('Copiado')).toBeInTheDocument();
        });

        it('no rompe si el portapapeles no esta disponible', async () => {
            Object.defineProperty(navigator, 'clipboard', {
                value: { writeText: vi.fn().mockRejectedValue(new Error('denied')) },
                configurable: true,
            });
            await setup();

            fireEvent.click(screen.getByTitle('Copiar enlace a este ticket'));

            await waitFor(() => expect(screen.getByText('Enlace')).toBeInTheDocument());
        });

        it('vuelve a mostrar Enlace unos segundos despues de copiar', async () => {
            vi.useFakeTimers({ shouldAdvanceTime: true });
            const writeText = vi.fn().mockResolvedValue(undefined);
            Object.defineProperty(navigator, 'clipboard', { value: { writeText }, configurable: true });
            await setup();

            fireEvent.click(screen.getByTitle('Copiar enlace a este ticket'));
            expect(await screen.findByText('Copiado')).toBeInTheDocument();

            await vi.advanceTimersByTimeAsync(2000);

            expect(await screen.findByText('Enlace')).toBeInTheDocument();
            vi.useRealTimers();
        });

        it('cierra el detalle con el boton Cerrar', async () => {
            const { onClose } = await setup();

            fireEvent.click(screen.getByRole('button', { name: 'Cerrar' }));

            expect(onClose).toHaveBeenCalledTimes(1);
        });

        it('muestra la fecha de ultima actualizacion cuando existe', async () => {
            const { unmount } = await setup();
            expect(screen.getByText(/\u00daltima actualizaci\u00f3n/)).toBeInTheDocument();
            unmount();

            await setup({ ticket: { updated_at: '' } });
            expect(screen.queryByText(/\u00daltima actualizaci\u00f3n/)).toBeNull();
        });
    });
});
