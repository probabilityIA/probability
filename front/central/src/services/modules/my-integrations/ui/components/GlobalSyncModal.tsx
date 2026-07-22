'use client';

import { useState } from 'react';
import { RefreshCw, ArrowRightLeft, ChevronDown } from 'lucide-react';
import { Modal } from '@/shared/ui/modal';
import type { Integration } from '@/services/integrations/core/domain/types';
import { getSyncProvider } from '../providers';
import { GlobalInventoryPanel } from './GlobalInventoryPanel';
import { GlobalProductsPanel } from './GlobalProductsPanel';

interface GlobalSyncModalProps {
    isOpen: boolean;
    onClose: () => void;
    integrations: Integration[];
    businessId: number | null;
}

type PanelKey = 'inventory' | 'products';

export function GlobalSyncModal({ isOpen, onClose, integrations, businessId }: GlobalSyncModalProps) {
    const [open, setOpen] = useState<PanelKey>('inventory');

    const eligible = integrations.filter(i => i.is_active && getSyncProvider(i.integration_type_id));

    const buttonClass = (key: PanelKey, activeClasses: string, idleHover: string) =>
        `flex items-center gap-2 rounded-full border px-4 py-1.5 text-xs font-semibold transition-colors ${
            open === key
                ? activeClasses
                : `border-gray-200 text-gray-600 dark:border-gray-600 dark:text-gray-300 ${idleHover}`
        }`;

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title={(
                <span className="inline-flex items-center justify-center gap-2">
                    <RefreshCw className="h-4 w-4" />
                    Sincronizacion global
                </span>
            )}
            size="2xl"
            zIndex={60}
        >
            {eligible.length === 0 ? (
                <p className="py-8 text-center text-sm italic text-gray-400 dark:text-gray-500">
                    No hay integraciones e-commerce activas para sincronizar
                </p>
            ) : (
                <div className="flex flex-col gap-3 pt-1">
                    <div className="flex flex-wrap items-center justify-center gap-2">
                        <button
                            onClick={() => setOpen('inventory')}
                            className={buttonClass(
                                'inventory',
                                'border-cyan-500 bg-cyan-50 text-cyan-700 dark:bg-cyan-900/30 dark:text-cyan-300',
                                'hover:border-cyan-400 hover:text-cyan-700 dark:hover:text-cyan-300',
                            )}
                        >
                            <RefreshCw size={13} />
                            Sincronizar inventario
                            <ChevronDown size={13} className={`transition-transform ${open === 'inventory' ? 'rotate-180' : ''}`} />
                        </button>
                        <button
                            onClick={() => setOpen('products')}
                            className={buttonClass(
                                'products',
                                'border-indigo-500 bg-indigo-50 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-300',
                                'hover:border-indigo-400 hover:text-indigo-700 dark:hover:text-indigo-300',
                            )}
                        >
                            <ArrowRightLeft size={13} />
                            Comparar productos
                            <ChevronDown size={13} className={`transition-transform ${open === 'products' ? 'rotate-180' : ''}`} />
                        </button>
                    </div>

                    <div className="border-t border-gray-100 pt-3 dark:border-gray-700">
                        {open === 'inventory' ? (
                            <GlobalInventoryPanel integrations={eligible} businessId={businessId} />
                        ) : (
                            <GlobalProductsPanel integrations={eligible} businessId={businessId} />
                        )}
                    </div>
                </div>
            )}
        </Modal>
    );
}
