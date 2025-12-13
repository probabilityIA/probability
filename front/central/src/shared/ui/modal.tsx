/**
 * Componente Modal reutilizable
 * Usa clases globales definidas en globals.css
 */

'use client';

import { ReactNode, useEffect } from 'react';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  size?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '4xl' | '5xl' | '6xl' | '7xl' | 'full';
  glass?: boolean; // Efecto glassmorphism
}

const sizeClasses = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-xl',
  '2xl': 'max-w-2xl',
  '4xl': 'max-w-4xl',
  '5xl': 'max-w-5xl',
  '6xl': 'max-w-6xl',
  '7xl': 'max-w-7xl',
  'full': 'max-w-[90vw] w-[90vw] max-h-[90vh] h-[90vh] mx-auto my-[5vh]',
};

export function Modal({ isOpen, onClose, title, children, size = 'md', glass = false }: ModalProps) {
  console.log('🔧 Modal - isOpen:', isOpen, 'title:', title);

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
      <div className="modal-backdrop" onClick={onClose} />

      {/* Modal */}
      <div className={size === 'full' ? 'fixed inset-0 z-50 flex items-center justify-center' : 'modal'}>
        {size === 'full' ? (
          <div
            className={`${glass ? 'bg-white/80 backdrop-blur-xl border border-white/20' : 'bg-white'} rounded-3xl shadow-2xl flex flex-col overflow-hidden`}
            style={{
              width: '90vw',
              height: '90vh',
              maxWidth: '90vw',
              maxHeight: '90vh',
            }}
          >
            {/* Header for full screen */}
            {title && (
              <div className="flex items-center justify-between px-8 py-6 border-b border-gray-200 bg-white">
                <h2 className="text-2xl font-bold text-gray-900">{title}</h2>
                <button
                  onClick={onClose}
                  className="text-gray-400 hover:text-gray-600 transition-colors p-2 hover:bg-gray-100 rounded-lg"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            )}

            {/* Scrollable Content */}
            <div className="flex-1 overflow-y-auto px-8 py-6">
              {children}
            </div>
          </div>
        ) : (
          <div className={`${glass ? 'modal-glass' : 'modal-content'} ${sizeClasses[size]}`}>
            {/* Header */}
            {title && (
              <div className="relative mb-4">
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

            {/* Content */}
            <div>{children}</div>
          </div>
        )}
      </div>
    </>
  );
}

