'use client';

interface InventoryDisabledNoticeProps {
    businessName?: string;
}

export function InventoryDisabledNotice({ businessName }: InventoryDisabledNoticeProps) {
    return (
        <div className="rounded-lg border border-amber-300 bg-amber-50 dark:bg-amber-900/20 dark:border-amber-700 p-4 flex items-start gap-3">
            <svg className="w-6 h-6 text-amber-600 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01M5 19h14a2 2 0 001.84-2.75L13.74 4a2 2 0 00-3.48 0L3.16 16.25A2 2 0 005 19z" />
            </svg>
            <div className="flex-1">
                <p className="font-semibold text-amber-900 dark:text-amber-100">Modulo de inventario desactivado</p>
                <p className="text-sm text-amber-800 dark:text-amber-200 mt-1">
                    {businessName
                        ? `El negocio "${businessName}" no tiene activo el modulo de inventario.`
                        : 'Este negocio no tiene activo el modulo de inventario.'}
                    {' '}Las ordenes no afectan el stock y no se generan novedades de inventario.
                    Activalo desde Tus Integraciones -&gt; Modulos Internos.
                </p>
            </div>
        </div>
    );
}
