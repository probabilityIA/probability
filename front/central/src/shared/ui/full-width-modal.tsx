/**
 * Modal de pantalla completa para visualización de información detallada
 * Diseñado para mostrar grandes cantidades de información sin scroll
 */

'use client';

import { ReactNode, useEffect } from 'react';

interface FullWidthModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string | ReactNode;
  children: ReactNode;
  width?: '85vw' | '90vw' | '95vw' | '98vw';
  height?: '85vh' | '90vh' | '95vh';
}

export function FullWidthModal({ 
  isOpen, 
  onClose, 
  title, 
  children, 
  width = '95vw',
  height = '90vh'
}: FullWidthModalProps) {
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
      <div 
        className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40" 
        onClick={onClose} 
      />

      {/* Modal */}
      <div className="fixed inset-0 z-50 flex items-center justify-center p-2 sm:p-4">
        <div
          className="bg-white rounded-xl shadow-2xl flex flex-col overflow-hidden"
          style={{
            width: width,
            height: height,
            maxWidth: width,
            maxHeight: height,
            minWidth: 0,
          }}
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          {title && (
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 bg-gray-50 flex-shrink-0">
              {typeof title === 'string' ? (
                <h2 className="text-xl sm:text-2xl font-bold text-gray-900">{title}</h2>
              ) : (
                <div className="text-xl sm:text-2xl font-bold text-gray-900">{title}</div>
              )}
              <button
                onClick={onClose}
                className="text-gray-400 hover:text-gray-600 transition-colors p-2 hover:bg-gray-200 rounded-lg"
                aria-label="Cerrar"
              >
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          )}

          {/* Content - Scrollable solo si es necesario */}
          <div className="flex-1 overflow-y-auto overflow-x-hidden px-6 py-4">
            {children}
          </div>
        </div>
      </div>
    </>
  );
}
