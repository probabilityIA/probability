import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { ChangePasswordForm } from './ChangePasswordForm';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    changePasswordAction: vi.fn(),
}));

vi.mock('@/shared/config', () => ({
    TokenStorage: { getSessionToken: vi.fn(() => 'token') },
}));

vi.mock('@/shared/ui/modal', () => ({
    Modal: ({ isOpen, children, onClose, title }: any) =>
        isOpen ? (
            <div data-testid="modal">
                <h3>{title}</h3>
                {children}
            </div>
        ) : null,
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any, fallback?: string) =>
        err instanceof Error ? err.message : fallback || 'Error'
    ),
}));

import { changePasswordAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const fillForm = (current: string, newPass: string, confirm: string) => {
    fireEvent.change(screen.getByPlaceholderText('Tu contraseña actual'), { target: { value: current } });
    fireEvent.change(screen.getByPlaceholderText('Mínimo 8 caracteres'), { target: { value: newPass } });
    fireEvent.change(screen.getByPlaceholderText('Confirma tu nueva contraseña'), { target: { value: confirm } });
};

const submitForm = () => {
    const form = screen.getByRole('button', { name: /cambiar contraseña/i }).closest('form')!;
    fireEvent.submit(form);
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('ChangePasswordForm', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    // ---------------------------------------------------------------
    // Validación del formulario
    // ---------------------------------------------------------------
    describe('validación del formulario', () => {
        it('debería mostrar error cuando los campos están vacíos y se envía el formulario', () => {
            render(<ChangePasswordForm />);

            submitForm();

            expect(screen.getByText('Todos los campos son requeridos')).toBeDefined();
        });

        it('debería mostrar error cuando la nueva contraseña tiene menos de 8 caracteres', () => {
            render(<ChangePasswordForm />);

            fillForm('oldPass123', 'short', 'short');
            submitForm();

            expect(screen.getByText('La nueva contraseña debe tener al menos 8 caracteres')).toBeDefined();
        });

        it('debería mostrar error cuando las contraseñas no coinciden', () => {
            render(<ChangePasswordForm />);

            fillForm('oldPass123', 'newPassword1', 'newPassword2');
            submitForm();

            expect(screen.getByText('Las contraseñas no coinciden')).toBeDefined();
        });

        it('debería mostrar error cuando la nueva contraseña es igual a la actual', () => {
            render(<ChangePasswordForm />);

            fillForm('samePassword123', 'samePassword123', 'samePassword123');
            submitForm();

            expect(screen.getByText('La nueva contraseña debe ser diferente a la actual')).toBeDefined();
        });
    });

    // ---------------------------------------------------------------
    // Envío exitoso
    // ---------------------------------------------------------------
    describe('envío exitoso', () => {
        it('debería llamar a changePasswordAction con los datos correctos', async () => {
            vi.mocked(changePasswordAction).mockResolvedValue({ success: true, message: 'OK' });

            render(<ChangePasswordForm />);

            fillForm('oldPass123', 'newPass456', 'newPass456');
            submitForm();

            await waitFor(() => {
                expect(changePasswordAction).toHaveBeenCalledWith({
                    current_password: 'oldPass123',
                    new_password: 'newPass456',
                });
            });
        });

        it('debería mostrar modal de éxito cuando la respuesta es exitosa', async () => {
            vi.mocked(changePasswordAction).mockResolvedValue({ success: true, message: 'Contraseña cambiada' });

            render(<ChangePasswordForm />);

            fillForm('oldPass123', 'newPass456', 'newPass456');
            submitForm();

            await waitFor(() => {
                expect(screen.getByTestId('modal')).toBeDefined();
            });
        });

        it('debería limpiar los campos después de un cambio exitoso', async () => {
            vi.mocked(changePasswordAction).mockResolvedValue({ success: true, message: 'OK' });

            render(<ChangePasswordForm />);

            fillForm('oldPass123', 'newPass456', 'newPass456');
            submitForm();

            await waitFor(() => {
                const currentInput = screen.getByPlaceholderText('Tu contraseña actual') as HTMLInputElement;
                expect(currentInput.value).toBe('');
            });
        });
    });

    // ---------------------------------------------------------------
    // Manejo de errores del servidor
    // ---------------------------------------------------------------
    describe('manejo de errores del servidor', () => {
        it('debería mostrar el error cuando changePasswordAction falla', async () => {
            vi.mocked(changePasswordAction).mockRejectedValue(new Error('Contraseña actual incorrecta'));

            render(<ChangePasswordForm />);

            fillForm('wrongPass', 'newPass456', 'newPass456');
            submitForm();

            await waitFor(() => {
                expect(screen.getByText('Contraseña actual incorrecta')).toBeDefined();
            });
        });
    });

    // ---------------------------------------------------------------
    // Interacción
    // ---------------------------------------------------------------
    describe('interacción', () => {
        it('debería llamar a onCancel cuando se presiona el botón cancelar', () => {
            const onCancel = vi.fn();
            render(<ChangePasswordForm onCancel={onCancel} />);

            fireEvent.click(screen.getByRole('button', { name: /cancelar/i }));

            expect(onCancel).toHaveBeenCalledOnce();
        });
    });
});
