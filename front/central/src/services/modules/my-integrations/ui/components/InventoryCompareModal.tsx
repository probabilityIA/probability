'use client';

import { useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { X } from 'lucide-react';
import type { Integration } from '@/services/integrations/core/domain/types';
import { InventoryCompareTable } from './InventoryCompareTable';

interface InventoryCompareModalProps {
    integration: Integration;
    businessId: number | null;
    onClose: () => void;
}

export function InventoryCompareModal({ integration, businessId, onClose }: InventoryCompareModalProps) {
    const [montado, setMontado] = useState(false);
    useEffect(() => setMontado(true), []);
    if (!montado) return null;

    return createPortal(
        <div className="fixed inset-0 z-[250] flex items-center justify-center p-3 sm:p-5" role="dialog" aria-modal="true">
            <button
                aria-label="Cerrar comparacion de inventario"
                onClick={onClose}
                className="absolute inset-0 cursor-default bg-gray-900/40 backdrop-blur-[2px]"
            />

            <div className="relative flex h-[88vh] w-full max-w-[84rem] flex-col overflow-hidden rounded-2xl bg-white shadow-2xl dark:bg-gray-900">
                <div className="flex items-start justify-between gap-4 px-6 pb-3 pt-5">
                    <div className="min-w-0">
                        <h2 className="text-[15px] font-bold text-gray-900 dark:text-white">
                            Comparar inventario con {integration.name}
                        </h2>
                        <p className="mt-0.5 text-[12px] text-gray-500 dark:text-gray-400">
                            Se lee el stock que el canal tiene ahora mismo, de a 100 productos, y tu decides cuales se
                            actualizan con el stock de Probability.
                        </p>
                    </div>
                    <button
                        onClick={onClose}
                        aria-label="Cerrar"
                        className="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-gray-100 text-gray-500 transition-colors hover:bg-gray-200 hover:text-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700"
                    >
                        <X size={15} />
                    </button>
                </div>

                <div className="flex min-h-0 flex-1 flex-col px-6 pb-4">
                    <InventoryCompareTable
                        businessId={businessId}
                        integrations={[integration]}
                        fixedIntegrationId={integration.id}
                    />
                </div>
            </div>
        </div>,
        document.body,
    );
}
