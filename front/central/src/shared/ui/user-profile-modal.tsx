'use client';

import { useState, useEffect } from 'react';
import { Modal } from './modal';
import { AvatarUpload } from './avatar-upload';
import { Button } from './button';
import { Spinner } from './spinner';
import { Alert } from './alert';
import { ConfirmModal } from './confirm-modal';
import { updateUserAction } from '@/services/auth/users/infra/actions';
import { updateBusinessAction } from '@/services/auth/business/infra/actions';
import { COLOR_PALETTES, BusinessPaletteColors } from '@/services/auth/business/domain/color-palettes';
import { ChangePasswordForm } from '@/services/auth/login/ui';
import { useDarkMode } from '@/shared/contexts/dark-mode-context';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useTheme } from '@/shared/providers/theme-provider';

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
  const { permissions, isSuperAdmin } = usePermissions();
  const { setColors, getColors } = useTheme();
  const [themeSaving, setThemeSaving] = useState<string | null>(null);
  const [themeError, setThemeError] = useState<string | null>(null);
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [removeAvatar, setRemoveAvatar] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [showChangePassword, setShowChangePassword] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [showUpdateConfirm, setShowUpdateConfirm] = useState(false);
  const [pendingFile, setPendingFile] = useState<File | null>(null);

  useEffect(() => {
    if (isOpen && user) {
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

  const roleName = (permissions?.role_name || '').toLowerCase();
  const businessId = permissions?.business_id || 0;
  const canChangeBusinessTheme =
    businessId > 0 && (isSuperAdmin || roleName === 'demo' || roleName.includes('admin'));
  const currentColors = getColors();

  const handlePaletteSelect = async (palette: { name: string; colors: BusinessPaletteColors }) => {
    if (themeSaving) return;

    setThemeSaving(palette.name);
    setThemeError(null);

    try {
      const response = await updateBusinessAction(businessId, {
        primary_color: palette.colors.primary,
        secondary_color: palette.colors.secondary,
        tertiary_color: palette.colors.tertiary,
        quaternary_color: palette.colors.quaternary,
      });

      if (response.success) {
        setColors(palette.colors);
      } else {
        setThemeError('No se pudo guardar el tema del negocio');
      }
    } catch (err) {
      setThemeError(err instanceof Error ? err.message : 'No se pudo guardar el tema del negocio');
    } finally {
      setThemeSaving(null);
    }
  };

  const handleSaveAvatar = async (file?: File | null, shouldRemove?: boolean) => {
    const fileToSave = file !== undefined ? file : avatarFile;
    const shouldRemoveAvatar = shouldRemove !== undefined ? shouldRemove : removeAvatar;

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

    input.oncancel = cleanup;
    
    setTimeout(cleanup, 1000);

    document.body.appendChild(input);
    input.click();
  };

  const handleFileSelect = (file: File | null) => {
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
    await handleSaveAvatar(pendingFile, false);
    setPendingFile(null);
  };

  const handleRemoveClick = () => {
    setShowDeleteConfirm(true);
  };

  const handleConfirmDelete = async () => {
    setShowDeleteConfirm(false);
    setAvatarFile(null);
    setRemoveAvatar(true);
    setError(null);
    setSuccess(false);
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
  };

  return (
    <>
      <Modal isOpen={isOpen} onClose={handleClose} title={showChangePassword ? "Cambiar Contraseña" : "Información del Perfil"}>
        <div className="space-y-6">
          {showChangePassword ? (
            <div>
              <button
                type="button"
                onClick={() => setShowChangePassword(false)}
                className="mb-4 flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300 dark:text-gray-400 hover:text-gray-900 dark:text-white dark:hover:text-gray-200 transition-colors"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                </svg>
                <span>Volver a información del perfil</span>
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
                  key={`${user.userId}-${isOpen}`}
                  currentAvatarUrl={removeAvatar ? null : (user.avatarUrl || null)}
                  onFileSelect={handleFileSelect}
                  onRemoveClick={handleRemoveClick}
                  onEditClick={handleEditClick}
                  disableClick={true}
                  size="lg"
                />

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
                    <svg className="w-5 h-5 text-gray-600 dark:text-gray-300 dark:text-gray-300" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" />
                    </svg>
                  )}
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-200 dark:text-gray-200 flex-1 text-left">
                    Tema Oscuro
                  </span>
                  <div className={`relative w-10 h-5 rounded-full transition-colors ${isDark ? 'bg-purple-600' : 'bg-gray-300'}`}>
                    <div className={`absolute top-0.5 w-4 h-4 bg-white rounded-full shadow transition-transform ${isDark ? 'translate-x-5' : 'translate-x-0.5'}`} />
                  </div>
                </button>

                {canChangeBusinessTheme && (
                  <div className="w-full max-w-xs">
                    <div className="flex items-center gap-2 mb-2">
                      <svg className="w-4 h-4 text-gray-500 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828L11 19.172M7 17h.01" />
                      </svg>
                      <span className="text-sm font-medium text-gray-700 dark:text-gray-200">
                        Tema del negocio
                      </span>
                    </div>

                    {themeError && (
                      <div className="mb-2">
                        <Alert type="error" onClose={() => setThemeError(null)}>{themeError}</Alert>
                      </div>
                    )}

                    <div className="grid grid-cols-5 gap-2">
                      {COLOR_PALETTES.map((palette) => {
                        const isCurrent = currentColors?.primary?.toUpperCase() === palette.colors.primary.toUpperCase();
                        const isSaving = themeSaving === palette.name;
                        return (
                          <button
                            key={palette.name}
                            type="button"
                            onClick={() => handlePaletteSelect(palette)}
                            disabled={themeSaving !== null}
                            title={palette.name}
                            aria-label={`Aplicar tema ${palette.name}`}
                            className={`relative h-9 rounded-md overflow-hidden border-2 transition-all disabled:opacity-50 disabled:cursor-not-allowed ${
                              isCurrent
                                ? 'border-purple-500 ring-2 ring-purple-300 dark:ring-purple-800'
                                : 'border-gray-200 dark:border-gray-600 hover:border-gray-400 dark:hover:border-gray-400'
                            }`}
                          >
                            <span className="absolute inset-0 flex">
                              <span className="flex-1" style={{ backgroundColor: palette.colors.primary }} />
                              <span className="flex-1" style={{ backgroundColor: palette.colors.tertiary }} />
                              <span className="flex-1" style={{ backgroundColor: palette.colors.quaternary }} />
                            </span>
                            {isSaving && (
                              <span className="absolute inset-0 flex items-center justify-center bg-black/40">
                                <Spinner size="sm" />
                              </span>
                            )}
                          </button>
                        );
                      })}
                    </div>
                    <p className="mt-2 text-xs text-gray-500 dark:text-gray-400">
                      Cambia los colores de toda la plataforma para tu negocio.
                    </p>
                  </div>
                )}

                <button
                  type="button"
                  onClick={() => setShowChangePassword(true)}
                  className="w-full max-w-xs flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg bg-purple-500 hover:bg-purple-600 dark:bg-purple-600 dark:hover:bg-purple-700 text-white dark:text-white border border-purple-600 dark:border-purple-700 transition-colors font-medium text-sm"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                  </svg>
                  Cambiar contraseña
                </button>
              </div>
            </>
          )}
        </div>
      </Modal>

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
