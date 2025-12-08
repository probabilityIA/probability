'use client';

import { useState } from 'react';
import { Modal } from './modal';
import { AvatarUpload } from './avatar-upload';
import { Button } from './button';
import { Spinner } from './spinner';
import { Alert } from './alert';
import { updateUserAction } from '@/services/auth/users/infra/actions';

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
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [removeAvatar, setRemoveAvatar] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  if (!user) return null;

  const handleFileSelect = (file: File | null) => {
    setAvatarFile(file);
    // Si se selecciona un archivo, no eliminar el avatar
    if (file) {
      setRemoveAvatar(false);
    }
    setError(null);
    setSuccess(false);
  };

  const handleRemoveAvatar = () => {
    setAvatarFile(null);
    setRemoveAvatar(true);
    setError(null);
    setSuccess(false);
  };

  const handleSave = async () => {
    // Validar que haya algo que hacer
    if (!avatarFile && !removeAvatar) {
      setError('Por favor selecciona una imagen o elimina la actual');
      return;
    }

    setLoading(true);
    setError(null);
    setSuccess(false);

    try {
      const userId = parseInt(user.userId, 10);
      if (isNaN(userId)) {
        setError('ID de usuario inválido');
        return;
      }

      const updateData: { avatarFile?: File; remove_avatar?: boolean } = {};
      if (avatarFile) {
        updateData.avatarFile = avatarFile;
      }
      if (removeAvatar && user.avatarUrl) {
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

  const handleClose = () => {
    setAvatarFile(null);
    setError(null);
    setSuccess(false);
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title="Cambiar Foto de Perfil">
      <div className="space-y-6">
        {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}
        {success && <Alert type="success">Foto de perfil actualizada exitosamente</Alert>}

        <div className="flex flex-col items-center gap-4">
          <AvatarUpload
            currentAvatarUrl={removeAvatar ? null : (user.avatarUrl || null)}
            onFileSelect={handleFileSelect}
            size="lg"
          />
          
          {user.avatarUrl && !removeAvatar && (
            <Button
              type="button"
              variant="danger"
              size="sm"
              onClick={handleRemoveAvatar}
            >
              Eliminar foto actual
            </Button>
          )}
        </div>

        <div className="flex justify-end gap-2 pt-4 border-t">
          <Button type="button" variant="secondary" onClick={handleClose}>
            Cancelar
          </Button>
          <Button
            type="button"
            onClick={handleSave}
            disabled={loading || (!avatarFile && !removeAvatar)}
          >
            {loading ? <Spinner size="sm" /> : 'Guardar'}
          </Button>
        </div>
      </div>
    </Modal>
  );
}
