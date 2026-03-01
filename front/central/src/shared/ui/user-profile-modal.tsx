'use client';

import { useState, useEffect } from 'react';
import { Modal } from './modal';
import { AvatarUpload } from './avatar-upload';
import { Button } from './button';
import { Spinner } from './spinner';
import { Alert } from './alert';
import { ConfirmModal } from './confirm-modal';
import { updateUserAction } from '@/services/auth/users/infra/actions';
import { ChangePasswordForm } from '@/services/auth/login/ui';
import { useDarkMode } from '@/shared/contexts/dark-mode-context';

interface UserProfileModalProps {
  isOpen: boolean;
  onClose: () => void;
  user: {
    userId: string;
    name: string;
    email: string;
    role: string;
    avatarUrl?: string;
  } | null;
  onUpdate?: () => void;
}

export function UserProfileModal({ isOpen, onClose, user, onUpdate }: UserProfileModalProps) {
  const { isDark, toggleDarkMode } = useDarkMode();
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [removeAvatar, setRemoveAvatar] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [showChangePassword, setShowChangePassword] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [showUpdateConfirm, setShowUpdateConfirm] = useState(false);
  const [pendingFile, setPendingFile] = useState<File | null>(null);

  // Resetear estado cuando el modal se abre o cambia el usuario
  useEffect(() => {
    if (isOpen && user) {
      console.log('UserProfileModal - Modal opened/User changed:', {
        userId: user.userId,
        avatarUrl: user.avatarUrl,
        isOpen
      });
      setAvatarFile(null);
      setRemoveAvatar(false);
      setError(null);
      setSuccess(false);
      setShowChangePassword(false);
      setShowDeleteConfirm(false);
      setShowUpdateConfirm(false);
      setPendingFile(null);
    }
  }, [isOpen, user?.userId]);

  if (!user) return null;

  const handleSaveAvatar = async (file?: File | null, shouldRemove?: boolean) => {
    const fileToSave = file !== undefined ? file : avatarFile;
    const shouldRemoveAvatar = shouldRemove !== undefined ? shouldRemove : removeAvatar;

    // Validar que haya algo que hacer
    if (!fileToSave && !shouldRemoveAvatar) {
      return;
    }

    setLoading(true);
    setError(null);
    setSuccess(false);

    try {
      const userId = parseInt(user.userId, 10);
      if (isNaN(userId)) {
        setError('ID de usuario inválido');
        setLoading(false);
        return;
      }

      const updateData: { avatarFile?: File; remove_avatar?: boolean } = {};
      if (fileToSave) {
        updateData.avatarFile = fileToSave;
      }
      if (shouldRemoveAvatar && user.avatarUrl) {
        updateData.remove_avatar = true;
      }

      const response = await updateUserAction(userId, updateData);

      if (response.success) {
        setSuccess(true);
        setTimeout(() => {
          if (onUpdate) onUpdate();
          onClose();
          setAvatarFile(null);
          setRemoveAvatar(false);
          setSuccess(false);
        }, 1500);
      } else {
        setError('Error al actualizar la foto de perfil');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al actualizar la foto de perfil');
    } finally {
      setLoading(false);
    }
  };

  const handleEditClick = () => {
    // Crear input file temporal para abrir selector
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*';
    input.style.display = 'none';
    
    const cleanup = () => {
      if (document.body.contains(input)) {
        document.body.removeChild(input);
      }
    };

    input.onchange = (e) => {
      const file = (e.target as HTMLInputElement).files?.[0] || null;
      if (file) {
        setPendingFile(file);
        setShowUpdateConfirm(true);
      }
      cleanup();
    };

    // Limpiar si el usuario cancela (cuando el input pierde el foco)
    input.oncancel = cleanup;
    
    // También limpiar después de un tiempo por si acaso
    setTimeout(cleanup, 1000);

    document.body.appendChild(input);
    input.click();
  };

  const handleFileSelect = (file: File | null) => {
    // Esta función ya no se usa directamente, pero la mantenemos por compatibilidad
    if (file) {
      setPendingFile(file);
      setShowUpdateConfirm(true);
    }
  };

  const handleConfirmUpdate = async () => {
    if (!pendingFile) return;
    
    setShowUpdateConfirm(false);
    setAvatarFile(pendingFile);
    setRemoveAvatar(false);
    setError(null);
    setSuccess(false);
    // Guardar automáticamente cuando se confirma la actualización
    await handleSaveAvatar(pendingFile, false);
    setPendingFile(null);
  };

  const handleRemoveClick = () => {
    // Mostrar modal de confirmación
    setShowDeleteConfirm(true);
  };

  const handleConfirmDelete = async () => {
    setShowDeleteConfirm(false);
    setAvatarFile(null);
    setRemoveAvatar(true);
    setError(null);
    setSuccess(false);
    // Guardar automáticamente cuando se confirma la eliminación
    await handleSaveAvatar(null, true);
  };

  const handleClose = () => {
    setAvatarFile(null);
    setRemoveAvatar(false);
    setError(null);
    setSuccess(false);
    setShowChangePassword(false);
    setShowDeleteConfirm(false);
    onClose();
  };

  const handlePasswordChangeSuccess = () => {
    setShowChangePassword(false);
    // Opcional: mostrar mensaje de éxito o recargar
  };

  return (
    <>
      <Modal isOpen={isOpen} onClose={handleClose} title={showChangePassword ? "Cambiar Contraseña" : "Cambiar Foto de Perfil"}>
        <div className="space-y-6">
          {/* Vista de cambio de contraseña */}
          {showChangePassword ? (
            <div>
              <button
                type="button"
                onClick={() => setShowChangePassword(false)}
                className="mb-4 flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200 transition-colors"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                </svg>
                <span>Volver a foto de perfil</span>
              </button>
              <ChangePasswordForm
                onSuccess={handlePasswordChangeSuccess}
                onCancel={() => setShowChangePassword(false)}
              />
            </div>
          ) : (
            /* Vista de foto de perfil */
            <>
              {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}
              {success && <Alert type="success">Foto de perfil actualizada exitosamente</Alert>}
              {loading && (
                <div className="flex justify-center">
                  <Spinner size="md" />
                </div>
              )}

              <div className="flex flex-col items-center gap-4">
                <AvatarUpload
                  key={`${user.userId}-${isOpen}`} // Forzar re-render cuando cambia el usuario o se abre el modal
                  currentAvatarUrl={removeAvatar ? null : (user.avatarUrl || null)}
                  onFileSelect={handleFileSelect}
                  onRemoveClick={handleRemoveClick}
                  onEditClick={handleEditClick}
                  disableClick={true}
                  size="lg"
                />

                {/* Toggle Tema Oscuro */}
                <button
                  type="button"
                  onClick={toggleDarkMode}
                  className="flex items-center gap-3 px-4 py-2.5 rounded-lg border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors w-full max-w-xs"
                >
                  {isDark ? (
                    <svg className="w-5 h-5 text-yellow-400" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                    </svg>
                  ) : (
                    <svg className="w-5 h-5 text-gray-600 dark:text-gray-300" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" />
                    </svg>
                  )}
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-200 flex-1 text-left">
                    Tema Oscuro
                  </span>
                  <div className={`relative w-10 h-5 rounded-full transition-colors ${isDark ? 'bg-purple-600' : 'bg-gray-300'}`}>
                    <div className={`absolute top-0.5 w-4 h-4 bg-white rounded-full shadow transition-transform ${isDark ? 'translate-x-5' : 'translate-x-0.5'}`} />
                  </div>
                </button>

                {/* Botón para cambiar contraseña debajo de la foto */}
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  onClick={() => setShowChangePassword(true)}
                  className="flex items-center gap-2"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                  </svg>
                  Cambiar contraseña
                </Button>
              </div>
            </>
          )}
        </div>
      </Modal>

      {/* Modal de confirmación para eliminar foto */}
      <ConfirmModal
        isOpen={showDeleteConfirm}
        onClose={() => setShowDeleteConfirm(false)}
        onConfirm={handleConfirmDelete}
        title="Eliminar foto de perfil"
        message="¿Estás seguro de que deseas eliminar tu foto de perfil? Esta acción no se puede deshacer."
        confirmText="Eliminar"
        cancelText="Cancelar"
        type="danger"
      />

      {/* Modal de confirmación para actualizar foto */}
      <ConfirmModal
        isOpen={showUpdateConfirm}
        onClose={() => {
          setShowUpdateConfirm(false);
          setPendingFile(null);
        }}
        onConfirm={handleConfirmUpdate}
        title="Actualizar foto de perfil"
        message="¿Deseas actualizar tu foto de perfil con la imagen seleccionada?"
        confirmText="Actualizar"
        cancelText="Cancelar"
        type="info"
      />
    </>
  );
}
