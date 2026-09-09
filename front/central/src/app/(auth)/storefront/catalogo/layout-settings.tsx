'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Squares2X2Icon, XMarkIcon } from '@heroicons/react/24/outline';
import { updateCatalogLayoutAction } from '@/services/modules/storefront/infra/actions';
import { CatalogLayout } from '@/services/modules/storefront/domain/types';
import { usePermissions } from '@/shared/contexts/permissions-context';

interface CatalogLayoutSettingsProps {
    layout: CatalogLayout;
    businessId?: number;
}

const PRESETS: CatalogLayout[] = [
    { columns: 2, rows: 6 },
    { columns: 3, rows: 4 },
    { columns: 4, rows: 3 },
    { columns: 5, rows: 3 },
];

export function CatalogLayoutSettings({ layout, businessId }: CatalogLayoutSettingsProps) {
    const router = useRouter();
    const { permissions } = usePermissions();
    const [open, setOpen] = useState(false);
    const [columns, setColumns] = useState(layout.columns);
    const [rows, setRows] = useState(layout.rows);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSave = async () => {
        setSaving(true);
        setError(null);
        const result = await updateCatalogLayoutAction({ columns, rows }, businessId);
        setSaving(false);

        if (!result.success) {
            setError(result.message || 'Error al guardar');
            return;
        }
        setOpen(false);
        router.refresh();
    };

    if (permissions?.role_name === 'cliente_final') {
        return null;
    }

    return (
        <div className="relative">
            <button
                type="button"
                onClick={() => setOpen(prev => !prev)}
                className="flex items-center gap-2 px-4 py-2 h-[42px] text-sm font-medium border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
                <Squares2X2Icon className="w-4 h-4" />
                Diseno
            </button>

            {open && (
                <>
                    <div className="fixed inset-0 z-30" onClick={() => setOpen(false)} />
                    <div className="absolute right-0 mt-2 w-80 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-xl z-40 p-4">
                        <div className="flex items-center justify-between mb-3">
                            <h3 className="text-sm font-bold text-gray-900 dark:text-white">Diseno del catalogo</h3>
                            <button onClick={() => setOpen(false)} className="text-gray-400 hover:text-gray-600">
                                <XMarkIcon className="w-4 h-4" />
                            </button>
                        </div>

                        <p className="text-xs text-gray-500 dark:text-gray-400 mb-3">
                            Elige cuantos productos se ven por fila y por pagina en el catalogo de tus clientes.
                        </p>

                        <div className="grid grid-cols-2 gap-2 mb-4">
                            {PRESETS.map(preset => (
                                <button
                                    key={`${preset.columns}x${preset.rows}`}
                                    type="button"
                                    onClick={() => { setColumns(preset.columns); setRows(preset.rows); }}
                                    className={`px-3 py-2 rounded-lg text-xs font-semibold border transition-colors ${
                                        columns === preset.columns && rows === preset.rows
                                            ? 'bg-indigo-600 border-indigo-600 text-white'
                                            : 'border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700'
                                    }`}
                                >
                                    {preset.columns} x {preset.rows}
                                </button>
                            ))}
                        </div>

                        <div className="grid grid-cols-2 gap-3 mb-4">
                            <div>
                                <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Columnas</label>
                                <input
                                    type="number"
                                    min={1}
                                    max={6}
                                    value={columns}
                                    onChange={e => setColumns(Math.min(6, Math.max(1, Number(e.target.value) || 1)))}
                                    className="w-full px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                                />
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Filas por pagina</label>
                                <input
                                    type="number"
                                    min={1}
                                    max={12}
                                    value={rows}
                                    onChange={e => setRows(Math.min(12, Math.max(1, Number(e.target.value) || 1)))}
                                    className="w-full px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                                />
                            </div>
                        </div>

                        <p className="text-xs text-gray-400 dark:text-gray-500 mb-3">{columns * rows} productos por pagina</p>

                        {error && <p className="text-xs text-red-500 mb-3">{error}</p>}

                        <button
                            onClick={handleSave}
                            disabled={saving}
                            className="w-full py-2 bg-indigo-600 text-white text-sm font-semibold rounded-lg hover:bg-indigo-700 disabled:bg-gray-300 transition-colors"
                        >
                            {saving ? 'Guardando...' : 'Guardar'}
                        </button>
                    </div>
                </>
            )}
        </div>
    );
}
