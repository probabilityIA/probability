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
  transparent?: boolean; // NEW: Fondo transparente sin sombra
}

const sizeClasses = {
  sm: 'max-w-sm w-[95vw] sm:w-full',
  md: 'max-w-md w-[95vw] sm:w-full',
  lg: 'max-w-lg w-[95vw] sm:w-full',
  xl: 'max-w-xl w-[95vw] sm:w-full',
  '2xl': 'max-w-2xl w-[95vw] sm:w-full',
  '4xl': 'max-w-4xl w-[95vw] sm:w-full',
  '5xl': 'max-w-5xl w-[95vw] sm:w-full',
  '6xl': 'max-w-6xl w-[95vw] sm:w-full',
  '7xl': 'max-w-7xl w-[95vw] sm:w-full',
  'full': 'max-w-[95vw] w-[95vw]',
};

export function Modal({ isOpen, onClose, title, children, size = 'md', glass = false, transparent = false }: ModalProps) {
  console.log('🔧 Modal - isOpen:', isOpen, 'title:', title, 'size:', size);

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
      <div className="fixed inset-0 z-50 flex items-center justify-center">
        {size === 'full' ? (
          <div
            className={`${transparent ? 'bg-transparent shadow-none border-none' : (glass ? 'bg-white/80 backdrop-blur-xl border border-white/20' : 'bg-white shadow-2xl rounded-3xl')} flex flex-col overflow-hidden`}
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
            <div className="flex-1 overflow-y-auto px-6 sm:px-8 py-6">
              {children}
            </div>
          </div>
        ) : (
          <div
            className={`${transparent ? 'bg-transparent shadow-none border-none' : (size === 'sm' || size === 'md' ? (glass ? 'modal-glass' : 'modal-content') : 'bg-white rounded-2xl shadow-2xl p-6 sm:p-8')} max-h-[90vh] overflow-hidden flex flex-col`}
            style={
              size === 'sm' || size === 'md'
                ? {
                  maxWidth: size === 'sm' ? '28rem' : '32rem',
                  width: '95vw'
                }
                : size === '5xl' || size === '6xl' || size === '7xl'
                  ? {
                    width: size === '5xl' ? '90vw' : size === '6xl' ? '95vw' : '98vw',
                    maxWidth: size === '5xl' ? '90vw' : size === '6xl' ? '95vw' : '98vw',
                    minWidth: 0
                  }
                  : undefined
            }
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
            <div className="flex-1 overflow-y-auto overflow-x-hidden pr-1 sm:pr-2 -mr-1 sm:-mr-2 w-full max-w-full">{children}</div>
          </div>
        )}
      </div>
    </>
  );
}

