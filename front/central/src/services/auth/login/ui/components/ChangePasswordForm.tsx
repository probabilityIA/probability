'use client';

import { useState, useTransition } from 'react';
import { changePasswordAction } from '../../infra/actions';
import { TokenStorage } from '@/shared/config';
import { Modal } from '@/shared/ui/modal';

interface ChangePasswordFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
}

export const ChangePasswordForm = ({ onSuccess, onCancel }: ChangePasswordFormProps) => {
    const [currentPassword, setCurrentPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [showCurrentPassword, setShowCurrentPassword] = useState(false);
    const [showNewPassword, setShowNewPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [isPending, startTransition] = useTransition();
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const validateForm = (): boolean => {
        if (!currentPassword || !newPassword || !confirmPassword) {
            setError('Todos los campos son requeridos');
            return false;
        }

        if (newPassword.length < 8) {
            setError('La nueva contraseña debe tener al menos 8 caracteres');
            return false;
        }

        if (newPassword !== confirmPassword) {
            setError('Las contraseñas no coinciden');
            return false;
        }

        if (currentPassword === newPassword) {
            setError('La nueva contraseña debe ser diferente a la actual');
            return false;
        }

        return true;
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setSuccess(null);

        if (!validateForm()) {
            return;
        }

        startTransition(async () => {
            try {
                const response = await changePasswordAction(
                    {
                        current_password: currentPassword,
                        new_password: newPassword,
                    }
                );

                if (response.success) {
                    setSuccess(response.message || 'Contraseña cambiada exitosamente');
                    setCurrentPassword('');
                    setNewPassword('');
                    setConfirmPassword('');
                    // No cerrar automáticamente, mostrar modal de éxito
                }
            } catch (err: any) {
                console.error(err);
                setError(err.message || 'Error al cambiar la contraseña. Por favor intenta de nuevo.');
            }
        });
    };

    return (
        <div className="w-full">
            <div className="mb-8 sm:mb-10">
                <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
                    Cambiar Contraseña
                </h2>
                <p className="text-sm text-gray-500">
                    Ingresa tu contraseña actual y la nueva contraseña que deseas usar.
                </p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-6">
                {/* Current Password Field */}
                <div className="space-y-2">
                    <label htmlFor="currentPassword" className="block text-sm font-bold text-gray-400">
                        Contraseña Actual
                    </label>
                    <div className="relative">
                        <input
                            id="currentPassword"
                            name="currentPassword"
                            type={showCurrentPassword ? "text" : "password"}
                            autoComplete="current-password"
                            required
                            className="block w-full px-4 py-3 pr-12 bg-white border border-gray-200 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500/20 focus:border-green-500 transition-all"
                            placeholder="Tu contraseña actual"
                            value={currentPassword}
                            onChange={(e) => setCurrentPassword(e.target.value)}
                            disabled={isPending}
                        />
                        <button
                            type="button"
                            onClick={() => setShowCurrentPassword(!showCurrentPassword)}
                            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 focus:outline-none transition-colors"
                            tabIndex={-1}
                        >
                            {showCurrentPassword ? (
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                                </svg>
                            ) : (
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                                </svg>
                            )}
                        </button>
                    </div>
                </div>

                {/* New Password Field */}
                <div className="space-y-2">
                    <label htmlFor="newPassword" className="block text-sm font-bold text-gray-400">
                        Nueva Contraseña
                    </label>
                    <div className="relative">
                        <input
                            id="newPassword"
                            name="newPassword"
                            type={showNewPassword ? "text" : "password"}
                            autoComplete="new-password"
                            required
                            minLength={8}
                            className="block w-full px-4 py-3 pr-12 bg-white border border-gray-200 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500/20 focus:border-green-500 transition-all"
                            placeholder="Mínimo 8 caracteres"
                            value={newPassword}
                            onChange={(e) => setNewPassword(e.target.value)}
                            disabled={isPending}
                        />
                        <button
                            type="button"
                            onClick={() => setShowNewPassword(!showNewPassword)}
                            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 focus:outline-none transition-colors"
                            tabIndex={-1}
                        >
                            {showNewPassword ? (
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                                </svg>
                            ) : (
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                                </svg>
                            )}
                        </button>
                    </div>
                    <p className="text-xs text-gray-400">
                        La contraseña debe tener al menos 8 caracteres
                    </p>
                </div>

                {/* Confirm Password Field */}
                <div className="space-y-2">
                    <label htmlFor="confirmPassword" className="block text-sm font-bold text-gray-400">
                        Confirmar Nueva Contraseña
                    </label>
                    <div className="relative">
                        <input
                            id="confirmPassword"
                            name="confirmPassword"
                            type={showConfirmPassword ? "text" : "password"}
                            autoComplete="new-password"
                            required
                            minLength={8}
                            className="block w-full px-4 py-3 pr-12 bg-white border border-gray-200 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500/20 focus:border-green-500 transition-all"
                            placeholder="Confirma tu nueva contraseña"
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                            disabled={isPending}
                        />
                        <button
                            type="button"
                            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 focus:outline-none transition-colors"
                            tabIndex={-1}
                        >
                            {showConfirmPassword ? (
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                                </svg>
                            ) : (
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                                </svg>
                            )}
                        </button>
                    </div>
                </div>

                {/* Error Message */}
                {error && (
                    <div className="p-3 rounded-lg bg-red-50 text-red-500 text-sm">
                        {error}
                    </div>
                )}


                {/* Action Buttons */}
                <div className="flex gap-3 mt-8">
                    {onCancel && (
                        <button
                            type="button"
                            onClick={onCancel}
                            disabled={isPending}
                            className="flex-1 py-3 px-4 border border-gray-300 rounded-lg text-base font-bold text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            Cancelar
                        </button>
                    )}
                    <button
                        type="submit"
                        disabled={isPending}
                        className={`${onCancel ? 'flex-1' : 'w-full'} flex justify-center py-3 px-4 border border-transparent text-base font-bold rounded-lg text-white bg-[#4ade80] hover:bg-green-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition-all shadow-sm disabled:opacity-50 disabled:cursor-not-allowed`}
                    >
                        {isPending ? 'Cambiando contraseña...' : 'Cambiar Contraseña'}
                    </button>
                </div>
            </form>

            {/* Modal de éxito */}
            <Modal
                isOpen={!!success}
                onClose={() => {
                    setSuccess(null);
                    if (onSuccess) {
                        onSuccess();
                    }
                }}
                title="Contraseña cambiada exitosamente"
                size="sm"
                zIndex={60}
            >
                <div className="space-y-6">
                    {/* Mensaje de éxito con icono */}
                    <div className="flex items-start gap-4">
                        <div className="flex-shrink-0">
                            <svg className="w-8 h-8 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                        </div>
                        <div className="flex-1">
                            <p className="text-gray-700">
                                {success}
                            </p>
                        </div>
                    </div>

                    {/* Botón Aceptar */}
                    <div className="flex justify-end">
                        <button
                            className="btn btn-primary btn-sm px-6"
                            onClick={() => {
                                setSuccess(null);
                                if (onSuccess) {
                                    onSuccess();
                                }
                            }}
                        >
                            Aceptar
                        </button>
                    </div>
                </div>
            </Modal>
        </div>
    );
};
