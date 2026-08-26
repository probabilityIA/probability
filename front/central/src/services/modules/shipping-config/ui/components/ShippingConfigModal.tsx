'use client';

import { useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { ShippingConfigForm } from './ShippingConfigForm';
import { TruckIcon, XMarkIcon } from '@heroicons/react/24/outline';

interface ShippingConfigModalProps {
    isOpen: boolean;
    onClose: () => void;
    businessId?: number;
}

export function ShippingConfigModal({ isOpen, onClose, businessId }: ShippingConfigModalProps) {
    const [mounted, setMounted] = useState(false);

    useEffect(() => setMounted(true), []);

    if (!isOpen || !mounted) return null;

    return createPortal(
        <div className="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center p-4" style={{ zIndex: 70 }}>
            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl w-full max-w-4xl max-h-[90vh] flex flex-col overflow-hidden">
                <div className="flex items-center justify-between px-6 py-4 border-b dark:border-gray-700">
                    <div className="flex items-center gap-2">
                        <TruckIcon className="w-6 h-6" style={{ color: 'var(--color-primary)' }} />
                        <h2 className="text-xl font-bold text-gray-900 dark:text-white">Configuracion de envios</h2>
                    </div>
                    <button onClick={onClose} className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700" aria-label="Cerrar">
                        <XMarkIcon className="w-5 h-5 text-gray-500 dark:text-gray-400" />
                    </button>
                </div>
                <div className="flex-1 overflow-y-auto p-6">
                    <ShippingConfigForm businessId={businessId} onClose={onClose} />
                </div>
            </div>
        </div>,
        document.body
    );
}
