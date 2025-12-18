/**
 * Modal específico para formularios que necesitan más espacio
 * Versión sin limitaciones de ancho
 */

'use client';

import { ReactNode, useEffect } from 'react';

interface WideModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  width?: '90vw' | '95vw' | '98vw';
}

export function WideModal({ isOpen, onClose, title, children, width = '90vw' }: WideModalProps) {
  // Cerrar con ESC
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [isOpen, onClose]);

  // Prevenir scroll del body cuando el modal está abierto
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'unset';
    }
    return () => {
      document.body.style.overflow = 'unset';
    };
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-40" onClick={onClose} />

      {/* Modal */}
      <div className="fixed inset-0 z-50 flex items-center justify-center">
        <div
          className="bg-white rounded-2xl shadow-2xl p-6 sm:p-8 max-h-[90vh] overflow-hidden flex flex-col"
          style={{
            width: width,
            maxWidth: width,
            minWidth: 0
          }}
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          {title && (
            <div className="relative mb-4 flex-shrink-0">
              <h3 className="text-xl font-bold text-gray-900 text-center">{title}</h3>
              <button
                onClick={onClose}
                className="absolute right-0 top-0 text-gray-400 hover:text-gray-600 transition-colors"
              >
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          )}

          {/* Content - Scrollable */}
          <div className="flex-1 overflow-y-auto overflow-x-hidden pr-1 sm:pr-2 -mr-1 sm:-mr-2 w-full max-w-full">
            {children}
          </div>
        </div>
      </div>
    </>
  );
}

