'use client';

import { useState, useRef, useEffect } from 'react';
import { QrCodeIcon, ArrowPathIcon } from '@heroicons/react/24/outline';
import { Alert, Spinner } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import { useScan } from '@/services/modules/inventory/ui/hooks/useCapture';
import { ScanResult } from '@/services/modules/inventory/domain/capture-types';

const codeTypeLabels: Record<string, { label: string; color: string }> = {
    lpn: { label: 'License Plate', color: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200' },
    location: { label: 'Ubicación', color: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200' },
    serial: { label: 'Número de serie', color: 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900 dark:text-cyan-200' },
    lot: { label: 'Lote', color: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200' },
    barcode: { label: 'Código de barras (UoM)', color: 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200' },
    product: { label: 'Producto', color: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' },
    unknown: { label: 'Desconocido', color: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300' },
};

export default function MobileScanPage() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    const { result, loading, error, scan, reset } = useScan(businessId);
    const [code, setCode] = useState('');
    const [history, setHistory] = useState<ScanResult[]>([]);
    const inputRef = useRef<HTMLInputElement>(null);

    const requiresBusiness = isSuperAdmin && selectedBusinessId === null;

    useEffect(() => { inputRef.current?.focus(); }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!code.trim()) return;
        const r = await scan({ code: code.trim(), device_id: 'web-ui', action: 'scan' });
        if (r) setHistory((h) => [r, ...h].slice(0, 10));
        setCode('');
        inputRef.current?.focus();
    };

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="max-w-3xl mx-auto space-y-4">                {requiresBusiness ? <Alert type="info">Selecciona un negocio.</Alert> : (
                    <>
                        <div className="bg-gradient-to-br from-purple-600 to-purple-800 rounded-xl shadow-lg p-6">
                            <form onSubmit={handleSubmit} className="space-y-3">
                                <label className="block text-xs font-medium text-purple-100 uppercase tracking-wide">Código escaneado</label>
                                <div className="flex gap-2">
                                    <input
                                        ref={inputRef}
                                        autoFocus
                                        value={code}
                                        onChange={(e) => setCode(e.target.value)}
                                        placeholder="Escanea o escribe un código..."
                                        className="flex-1 px-4 py-3 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-lg font-mono focus:ring-4 focus:ring-purple-300 focus:outline-none"
                                    />
                                    <button type="submit" disabled={loading || !code.trim()} className="px-6 py-3 bg-white text-purple-700 font-bold rounded-lg hover:bg-purple-50 disabled:opacity-50 shadow-lg">
                                        {loading ? <Spinner size="sm" /> : 'Escanear'}
                                    </button>
                                </div>
                            </form>
                        </div>

                        {error && <Alert type="error">{error}</Alert>}

                        {result && (
                            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm overflow-hidden">
                                <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
                                    <div>
                                        <h2 className="text-base font-semibold text-gray-900 dark:text-white">Último escaneo</h2>
                                        <p className="text-sm text-gray-500 dark:text-gray-400">
                                            {result.resolved ? '✅ Resuelto' : '⚠️ Sin resolver'}
                                        </p>
                                    </div>
                                    <button onClick={() => { reset(); setHistory([]); }} className="flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
                                        <ArrowPathIcon className="w-4 h-4" /> Limpiar
                                    </button>
                                </div>
                                <div className="p-6 space-y-4">
                                    {result.resolution && (
                                        <>
                                            <div className="flex items-center gap-3">
                                                <span className={`px-3 py-1 rounded-full text-xs font-semibold ${codeTypeLabels[result.resolution.code_type]?.color || codeTypeLabels.unknown.color}`}>
                                                    {codeTypeLabels[result.resolution.code_type]?.label || result.resolution.code_type}
                                                </span>
                                                <span className="font-mono text-lg font-bold text-gray-900 dark:text-white">{result.resolution.code}</span>
                                            </div>
                                            <dl className="grid grid-cols-2 md:grid-cols-3 gap-3">
                                                {result.resolution.matched_id !== undefined && (
                                                    <div>
                                                        <dt className="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">Match ID</dt>
                                                        <dd className="text-sm font-mono text-gray-900 dark:text-white">#{result.resolution.matched_id}</dd>
                                                    </div>
                                                )}
                                                {result.resolution.product_id && (
                                                    <div>
                                                        <dt className="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">Producto</dt>
                                                        <dd className="text-sm font-mono text-gray-900 dark:text-white">{result.resolution.product_id}</dd>
                                                    </div>
                                                )}
                                                {result.resolution.location_id && (
                                                    <div>
                                                        <dt className="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">Ubicación</dt>
                                                        <dd className="text-sm text-gray-900 dark:text-white">#{result.resolution.location_id}</dd>
                                                    </div>
                                                )}
                                                {result.resolution.lot_id && (
                                                    <div>
                                                        <dt className="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">Lote</dt>
                                                        <dd className="text-sm text-gray-900 dark:text-white">#{result.resolution.lot_id}</dd>
                                                    </div>
                                                )}
                                                {result.resolution.serial_id && (
                                                    <div>
                                                        <dt className="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">Serie</dt>
                                                        <dd className="text-sm text-gray-900 dark:text-white">#{result.resolution.serial_id}</dd>
                                                    </div>
                                                )}
                                                {result.resolution.lpn_id && (
                                                    <div>
                                                        <dt className="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">LPN</dt>
                                                        <dd className="text-sm text-gray-900 dark:text-white">#{result.resolution.lpn_id}</dd>
                                                    </div>
                                                )}
                                            </dl>
                                        </>
                                    )}
                                    {result.event && (
                                        <p className="text-xs text-gray-400 border-t border-gray-100 dark:border-gray-700 pt-3">
                                            Evento registrado #{result.event.id} · {new Date(result.event.scanned_at).toLocaleString()}
                                        </p>
                                    )}
                                </div>
                            </div>
                        )}

                        {history.length > 1 && (
                            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm">
                                <h3 className="px-6 py-3 border-b border-gray-200 dark:border-gray-700 text-sm font-semibold text-gray-900 dark:text-white">Historial (últimos 10)</h3>
                                <ul className="divide-y divide-gray-100 dark:divide-gray-700/50">
                                    {history.slice(1).map((h, i) => (
                                        <li key={i} className="px-6 py-3 flex items-center justify-between">
                                            <div className="flex items-center gap-3">
                                                <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${codeTypeLabels[h.resolution?.code_type || 'unknown']?.color || codeTypeLabels.unknown.color}`}>
                                                    {codeTypeLabels[h.resolution?.code_type || 'unknown']?.label || 'Sin resolver'}
                                                </span>
                                                <span className="font-mono text-sm text-gray-700 dark:text-gray-200">{h.resolution?.code}</span>
                                            </div>
                                            <span className="text-xs text-gray-400">{h.event ? new Date(h.event.scanned_at).toLocaleTimeString() : ''}</span>
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        )}
                    </>
                )}
            </div>
        </div>
    );
}
