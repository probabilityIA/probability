'use client';

import { useState } from 'react';
import { Cog6ToothIcon } from '@heroicons/react/24/outline';
import { ShippingConfigModal } from './ShippingConfigModal';

interface ShippingConfigButtonProps {
    businessId?: number;
    className?: string;
}

export function ShippingConfigButton({ businessId, className }: ShippingConfigButtonProps) {
    const [open, setOpen] = useState(false);

    return (
        <>
            <button
                type="button"
                onClick={() => setOpen(true)}
                title="Configuración de envíos"
                aria-label="Configuracion de envios"
                className={className || 'p-2 rounded-lg border transition-colors hover:bg-gray-100 dark:hover:bg-gray-700'}
                style={className ? undefined : { borderColor: 'var(--color-primary)', color: 'var(--color-primary)' }}
            >
                <Cog6ToothIcon className="w-5 h-5" />
            </button>
            <ShippingConfigModal isOpen={open} onClose={() => setOpen(false)} businessId={businessId} />
        </>
    );
}
