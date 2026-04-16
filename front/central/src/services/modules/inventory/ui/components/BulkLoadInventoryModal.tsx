'use client';

import { useState, useRef } from 'react';
import { BulkLoadItem, BulkLoadResult, BulkLoadAccepted } from '../../domain/types';
import { bulkLoadInventoryAction } from '../../infra/actions';
import { Button, Alert } from '@/shared/ui';

interface BulkLoadInventoryModalProps {
    warehouseId: number;
    businessId?: number;
    onSuccess: () => void;
    onClose: () => void;
}

type Step = 'input' | 'preview' | 'result';

function isBulkLoadResult(r: BulkLoadResult | BulkLoadAccepted): r is BulkLoadResult {
    return 'items' in r;
}

export default function BulkLoadInventoryModal({ warehouseId, businessId, onSuccess, onClose }: BulkLoadInventoryModalProps) {
    const [step, setStep] = useState<Step>('input');
    const [items, setItems] = useState<BulkLoadItem[]>([{ sku: '', quantity: 0 }]);
    const [reason, setReason] = useState('Carga masiva de inventario');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [result, setResult] = useState<BulkLoadResult | BulkLoadAccepted | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const addRow = () => {
        setItems(prev => [...prev, { sku: '', quantity: 0 }]);
    };

    const removeRow = (index: number) => {
        setItems(prev => prev.filter((_, i) => i !== index));
    };

    const updateItem = (index: number, field: keyof BulkLoadItem, value: string | number | null) => {
        setItems(prev => prev.map((item, i) => {
            if (i !== index) return item;
            return { ...item, [field]: value };
        }));
    };

    const handleCSVImport = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        const reader = new FileReader();
        reader.onload = (event) => {
            const text = event.target?.result as string;
            const lines = text.split('\n').filter(line => line.trim());

            // Skip header if present
            const startIdx = lines[0]?.toLowerCase().includes('sku') ? 1 : 0;
            const parsed: BulkLoadItem[] = [];

            for (let i = startIdx; i < lines.length; i++) {
                const cols = lines[i].split(',').map(c => c.trim());
                if (cols.length < 2 || !cols[0]) continue;

                const item: BulkLoadItem = {
                    sku: cols[0],
                    quantity: parseInt(cols[1]) || 0,
                };
                if (cols[2] && cols[2] !== '') item.min_stock = parseInt(cols[2]) || null;
                if (cols[3] && cols[3] !== '') item.max_stock = parseInt(cols[3]) || null;
                if (cols[4] && cols[4] !== '') item.reorder_point = parseInt(cols[4]) || null;

                if (item.quantity > 0) {
                    parsed.push(item);
                }
            }

            if (parsed.length === 0) {
                setError('No se encontraron items validos en el CSV. Formato: sku,quantity,min_stock,max_stock,reorder_point');
                return;
            }

            setItems(parsed);
            setError(null);
        };
        reader.readAsText(file);
        // Reset input so same file can be re-selected
        if (fileInputRef.current) fileInputRef.current.value = '';
    };

    const validItems = items.filter(i => i.sku.trim() && i.quantity > 0);

    const handlePreview = () => {
        if (validItems.length === 0) {
            setError('Agrega al menos un item con SKU y cantidad validos');
            return;
        }
        setError(null);
        setStep('preview');
    };

    const handleSubmit = async () => {
        setLoading(true);
        setError(null);

        const res = await bulkLoadInventoryAction({
            warehouse_id: warehouseId,
            reason: reason.trim() || 'Carga masiva de inventario',
            items: validItems,
        }, businessId);

        if (!res.success) {
            setError(res.error);
        } else {
            setResult(res.data);
            setStep('result');
        }
        setLoading(false);
    };

    const handleDone = () => {
        onSuccess();
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-3xl max-h-[90vh] overflow-hidden flex flex-col">
                {/* Header */}
                <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 shrink-0">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                        {step === 'input' && 'Carga masiva de inventario'}
                        {step === 'preview' && 'Confirmar carga masiva'}
                        {step === 'result' && 'Resultado de la carga'}
                    </h2>
                    <button onClick={onClose} className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 text-xl leading-none">
                        &times;
                    </button>
                </div>

                <div className="p-6 overflow-y-auto flex-1">
                    {error && (
                        <Alert type="error" onClose={() => setError(null)}>
                            {error}
                        </Alert>
                    )}

                    {/* STEP: INPUT */}
                    {step === 'input' && (
                        <div className="space-y-4">
                            {/* CSV Import */}
                            <div className="flex items-center gap-3">
                                <input
                                    ref={fileInputRef}
                                    type="file"
                                    accept=".csv,.txt"
                                    onChange={handleCSVImport}
                                    className="hidden"
                                />
                                <Button
                                    type="button"
                                    variant="outline"
                                    onClick={() => fileInputRef.current?.click()}
                                >
                                    Importar CSV
                                </Button>
                                <span className="text-xs text-gray-500 dark:text-gray-400">
                                    Formato: sku,quantity,min_stock,max_stock,reorder_point
                                </span>
                            </div>

                            {/* Reason */}
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                                    Razon
                                </label>
                                <input
                                    type="text"
                                    value={reason}
                                    onChange={(e) => setReason(e.target.value)}
                                    className="w-full px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                    placeholder="Motivo de la carga"
                                />
                            </div>

                            {/* Table */}
                            <div className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                                <table className="w-full text-sm">
                                    <thead className="bg-gray-50 dark:bg-gray-700/50">
                                        <tr>
                                            <th className="px-3 py-2 text-left font-medium text-gray-700 dark:text-gray-200">SKU</th>
                                            <th className="px-3 py-2 text-left font-medium text-gray-700 dark:text-gray-200 w-24">Cantidad</th>
                                            <th className="px-3 py-2 text-left font-medium text-gray-700 dark:text-gray-200 w-24">Min</th>
                                            <th className="px-3 py-2 text-left font-medium text-gray-700 dark:text-gray-200 w-24">Max</th>
                                            <th className="px-3 py-2 w-10"></th>
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                                        {items.map((item, idx) => (
                                            <tr key={idx}>
                                                <td className="px-3 py-1.5">
                                                    <input
                                                        type="text"
                                                        value={item.sku}
                                                        onChange={(e) => updateItem(idx, 'sku', e.target.value)}
                                                        className="w-full px-2 py-1.5 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded text-sm"
                                                        placeholder="SKU-001"
                                                    />
                                                </td>
                                                <td className="px-3 py-1.5">
                                                    <input
                                                        type="number"
                                                        min="1"
                                                        value={item.quantity || ''}
                                                        onChange={(e) => updateItem(idx, 'quantity', parseInt(e.target.value) || 0)}
                                                        className="w-full px-2 py-1.5 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded text-sm"
                                                    />
                                                </td>
                                                <td className="px-3 py-1.5">
                                                    <input
                                                        type="number"
                                                        min="0"
                                                        value={item.min_stock ?? ''}
                                                        onChange={(e) => updateItem(idx, 'min_stock', e.target.value ? parseInt(e.target.value) : null)}
                                                        className="w-full px-2 py-1.5 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded text-sm"
                                                    />
                                                </td>
                                                <td className="px-3 py-1.5">
                                                    <input
                                                        type="number"
                                                        min="0"
                                                        value={item.max_stock ?? ''}
                                                        onChange={(e) => updateItem(idx, 'max_stock', e.target.value ? parseInt(e.target.value) : null)}
                                                        className="w-full px-2 py-1.5 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded text-sm"
                                                    />
                                                </td>
                                                <td className="px-3 py-1.5 text-center">
                                                    {items.length > 1 && (
                                                        <button
                                                            onClick={() => removeRow(idx)}
                                                            className="text-red-500 hover:text-red-700 text-lg leading-none"
                                                            title="Eliminar fila"
                                                        >
                                                            &times;
                                                        </button>
                                                    )}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>

                            <div className="flex items-center justify-between">
                                <Button type="button" variant="outline" onClick={addRow}>
                                    + Agregar fila
                                </Button>
                                <span className="text-sm text-gray-500 dark:text-gray-400">
                                    {validItems.length} item{validItems.length !== 1 ? 's' : ''} valido{validItems.length !== 1 ? 's' : ''}
                                </span>
                            </div>
                        </div>
                    )}

                    {/* STEP: PREVIEW */}
                    {step === 'preview' && (
                        <div className="space-y-4">
                            <p className="text-sm text-gray-600 dark:text-gray-300">
                                Se cargaran <strong>{validItems.length}</strong> producto{validItems.length !== 1 ? 's' : ''} a la bodega seleccionada.
                            </p>
                            <div className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden max-h-64 overflow-y-auto">
                                <table className="w-full text-sm">
                                    <thead className="bg-gray-50 dark:bg-gray-700/50 sticky top-0">
                                        <tr>
                                            <th className="px-3 py-2 text-left font-medium text-gray-700 dark:text-gray-200">#</th>
                                            <th className="px-3 py-2 text-left font-medium text-gray-700 dark:text-gray-200">SKU</th>
                                            <th className="px-3 py-2 text-right font-medium text-gray-700 dark:text-gray-200">Cantidad</th>
                                            <th className="px-3 py-2 text-right font-medium text-gray-700 dark:text-gray-200">Min</th>
                                            <th className="px-3 py-2 text-right font-medium text-gray-700 dark:text-gray-200">Max</th>
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                                        {validItems.map((item, idx) => (
                                            <tr key={idx}>
                                                <td className="px-3 py-2 text-gray-500 dark:text-gray-400">{idx + 1}</td>
                                                <td className="px-3 py-2 font-mono text-gray-900 dark:text-white">{item.sku}</td>
                                                <td className="px-3 py-2 text-right text-gray-900 dark:text-white">{item.quantity}</td>
                                                <td className="px-3 py-2 text-right text-gray-500 dark:text-gray-400">{item.min_stock ?? '-'}</td>
                                                <td className="px-3 py-2 text-right text-gray-500 dark:text-gray-400">{item.max_stock ?? '-'}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    )}

                    {/* STEP: RESULT */}
                    {step === 'result' && result && (
                        <div className="space-y-4">
                            {isBulkLoadResult(result) ? (
                                <>
                                    {/* Summary */}
                                    <div className="grid grid-cols-3 gap-4">
                                        <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-3 text-center">
                                            <p className="text-2xl font-bold text-gray-900 dark:text-white">{result.total_items}</p>
                                            <p className="text-xs text-gray-500 dark:text-gray-400">Total</p>
                                        </div>
                                        <div className="bg-green-50 dark:bg-green-900/20 rounded-lg p-3 text-center">
                                            <p className="text-2xl font-bold text-green-600">{result.success_count}</p>
                                            <p className="text-xs text-gray-500 dark:text-gray-400">Exitosos</p>
                                        </div>
                                        <div className={`rounded-lg p-3 text-center ${result.failure_count > 0 ? 'bg-red-50 dark:bg-red-900/20' : 'bg-gray-50 dark:bg-gray-700/50'}`}>
                                            <p className={`text-2xl font-bold ${result.failure_count > 0 ? 'text-red-600' : 'text-gray-400'}`}>{result.failure_count}</p>
                                            <p className="text-xs text-gray-500 dark:text-gray-400">Fallidos</p>
                                        </div>
                                    </div>

                                    {/* Detail */}
                                    <div className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden max-h-64 overflow-y-auto">
                                        <table className="w-full text-sm">
                                            <thead className="bg-gray-50 dark:bg-gray-700/50 sticky top-0">
                                                <tr>
                                                    <th className="px-3 py-2 text-left font-medium text-gray-700 dark:text-gray-200">SKU</th>
                                                    <th className="px-3 py-2 text-center font-medium text-gray-700 dark:text-gray-200">Estado</th>
                                                    <th className="px-3 py-2 text-right font-medium text-gray-700 dark:text-gray-200">Anterior</th>
                                                    <th className="px-3 py-2 text-right font-medium text-gray-700 dark:text-gray-200">Nuevo</th>
                                                    <th className="px-3 py-2 text-left font-medium text-gray-700 dark:text-gray-200">Error</th>
                                                </tr>
                                            </thead>
                                            <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                                                {result.items.map((item, idx) => (
                                                    <tr key={idx} className={!item.success ? 'bg-red-50/50 dark:bg-red-900/10' : ''}>
                                                        <td className="px-3 py-2 font-mono text-gray-900 dark:text-white">{item.sku}</td>
                                                        <td className="px-3 py-2 text-center">
                                                            {item.success ? (
                                                                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400">OK</span>
                                                            ) : (
                                                                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400">Error</span>
                                                            )}
                                                        </td>
                                                        <td className="px-3 py-2 text-right text-gray-600 dark:text-gray-300">{item.previous_qty}</td>
                                                        <td className="px-3 py-2 text-right font-medium text-gray-900 dark:text-white">{item.success ? item.new_qty : '-'}</td>
                                                        <td className="px-3 py-2 text-red-600 dark:text-red-400 text-xs">{item.error || ''}</td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    </div>
                                </>
                            ) : (
                                <Alert type="success">
                                    {result.message}. Se procesaran {result.total_items} items en segundo plano.
                                </Alert>
                            )}
                        </div>
                    )}
                </div>

                {/* Footer */}
                <div className="flex justify-end gap-3 px-6 py-4 border-t border-gray-200 dark:border-gray-700 shrink-0">
                    {step === 'input' && (
                        <>
                            <Button type="button" variant="outline" onClick={onClose}>
                                Cancelar
                            </Button>
                            <Button type="button" variant="primary" onClick={handlePreview} disabled={validItems.length === 0}>
                                Vista previa ({validItems.length})
                            </Button>
                        </>
                    )}
                    {step === 'preview' && (
                        <>
                            <Button type="button" variant="outline" onClick={() => setStep('input')} disabled={loading}>
                                Volver
                            </Button>
                            <Button type="button" variant="primary" onClick={handleSubmit} disabled={loading}>
                                {loading ? 'Procesando...' : `Cargar ${validItems.length} items`}
                            </Button>
                        </>
                    )}
                    {step === 'result' && (
                        <Button type="button" variant="primary" onClick={handleDone}>
                            Cerrar
                        </Button>
                    )}
                </div>
            </div>
        </div>
    );
}
