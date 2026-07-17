/**
 * Modal de confirmación reutilizable
 * Útil para confirmaciones de eliminación, cambios, etc.
 */

'use client';

import { ReactNode } from 'react';
import { Modal } from './modal';

interface ConfirmModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title?: string;
  message: ReactNode;
  confirmText?: string;
  cancelText?: string;
  type?: 'danger' | 'warning' | 'info';
}

export function ConfirmModal({
  isOpen,
  onClose,
  onConfirm,
  title = 'Confirmar acción',
  message,
  confirmText = 'Confirmar',
  cancelText = 'Cancelar',
  type = 'danger',
}: ConfirmModalProps) {
  const handleConfirm = () => {
    onConfirm();
    onClose();
  };

  const confirmButtonClass = {
    danger: 'btn btn-danger',
    warning: 'btn btn-quaternary',
    info: 'btn btn-primary',
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={title} size="sm" zIndex={60}>
      <div className="space-y-6">
        {/* Mensaje */}
        <div className="text-gray-700 dark:text-gray-200 text-sm">{message}</div>

        {/* Botones */}
        <div className="flex gap-3 justify-end">
          <button 
            className="btn btn-secondary btn-sm" 
            onClick={onClose}
          >
            {cancelText}
          </button>
          <button 
            className={`${confirmButtonClass[type]} btn-sm`}
            onClick={handleConfirm}
          >
            {confirmText}
          </button>
        </div>
      </div>
    </Modal>
  );
}

