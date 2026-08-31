'use client';

import { useState } from 'react';
import { ProductFamily } from '../../domain/types';
import { importFamilyVariantsPreviewAction, updateProductAction } from '../../infra/actions';
import { XMarkIcon } from '@heroicons/react/24/outline';

interface ImportFamilyVariantsModalProps {
    businessId?: number;
    family: ProductFamily;
    onClose: () => void;
    onAssociated: () => void;
}

interface MatchedRow {
    product_id: string;
    sku: string;
    name: string;
    current_family_id?: number;
    current_family_name?: string;
    weight?: number;
    length?: number;
    width?: number;
    height?: number;
}

export default function ImportFamilyVariantsModal({ businessId, family, onClose, onAssociated }: ImportFamilyVariantsModalProps) {
    const [file, setFile] = useState<File | null>(null);
    const [uploading, setUploading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [matched, setMatched] = useState<MatchedRow[]>([]);
    const [notFound, setNotFound] = useState<string[]>([]);
    const [variantLabels, setVariantLabels] = useState<Record<string, string>>({});
    const [applying, setApplying] = useState(false);
    const [progress, setProgress] = useState(0);
    const [failedSkus, setFailedSkus] = useState<string[]>([]);

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const selectedFile = e.target.files?.[0];
        if (!selectedFile) return;
        const ext = selectedFile.name.split('.').pop()?.toLowerCase();
        if (ext === 'xlsx' || ext === 'xls') {
            setFile(selectedFile);
            setError(null);
        } else {
            setError('Solo se permiten archivos Excel (.xlsx, .xls)');
            setFile(null);
        }
    };

    const handleUpload = async () => {
        if (!file) {
            setError('Por favor selecciona un archivo');
            return;
        }
        setUploading(true);
        setError(null);
        try {
            const result = await importFamilyVariantsPreviewAction(family.id, file, businessId);
            if (!result.success) {
                setError(result.message || 'Error al procesar el archivo');
                return;
            }
            const rows: MatchedRow[] = result.data?.matched || [];
            setMatched(rows);
            setNotFound(result.data?.not_found || []);
            setVariantLabels(Object.fromEntries(rows.map((r) => [r.product_id, ''])));
            if (rows.length === 0) {
                setError('Ningún SKU del archivo coincide con productos de este negocio');
            }
        } catch (err: any) {
            setError(err.message || 'Error inesperado al cargar el archivo');
        } finally {
            setUploading(false);
        }
    };

    const allLabelsFilled = matched.length > 0 && matched.every((r) => variantLabels[r.product_id]?.trim());

    const handleConfirm = async () => {
        if (!allLabelsFilled) return;
        setApplying(true);
        setProgress(0);
        setFailedSkus([]);

        const failed: string[] = [];
        await Promise.all(matched.map(async (row) => {
            const label = variantLabels[row.product_id].trim();
            try {
                const res = await updateProductAction(
                    row.product_id,
                    {
                        family_id: family.id,
                        variant_label: label,
                        variant_attributes: { variante: label },
                        weight: row.weight,
                        length: row.length,
                        width: row.width,
                        height: row.height,
                    },
                    businessId
                );
                if (res.success === false) failed.push(row.sku);
            } catch {
                failed.push(row.sku);
            } finally {
                setProgress((p) => p + 1);
            }
        }));

        setApplying(false);
        if (failed.length > 0) {
            setFailedSkus(failed);
        } else {
            onAssociated();
            onClose();
        }
    };

    return (
        <div className="fixed inset-0 z-[60] flex items-center justify-center p-4">
            <div className="absolute inset-0 bg-black/50" onClick={applying ? undefined : onClose} />
            <div className="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-2xl max-h-[85vh] flex flex-col">
                <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
                    <div>
                        <h3 className="text-base font-semibold text-gray-900 dark:text-white">Cargar productos por Excel</h3>
                        <p className="text-xs text-gray-500 mt-0.5">a la familia &quot;{family.name}&quot;</p>
                    </div>
                    <button
                        onClick={applying ? undefined : onClose}
                        disabled={applying}
                        className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors disabled:opacity-40"
                    >
                        <XMarkIcon className="w-5 h-5" />
                    </button>
                </div>

                <div className="overflow-auto flex-1 p-5 flex flex-col gap-4">
                    {matched.length === 0 ? (
                        <>
                            <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
                                <p className="text-sm text-blue-700 dark:text-blue-300">
                                    Puedes usar el mismo Excel que descargas en &quot;Peso y dimensiones&quot; (columnas sku, id, nombre, peso_kg, largo_cm, ancho_cm, alto_cm).
                                    Peso y dimensiones son opcionales: solo se necesita el sku para asociar el producto a esta familia.
                                </p>
                            </div>

                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Seleccionar archivo</label>
                                <input
                                    type="file"
                                    accept=".xlsx,.xls"
                                    onChange={handleFileChange}
                                    className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-orange-50 file:text-orange-700 hover:file:bg-orange-100"
                                />
                                {file && (
                                    <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">
                                        Archivo seleccionado: <strong>{file.name}</strong>
                                    </p>
                                )}
                            </div>

                            {error && (
                                <div className="px-3 py-2 rounded-lg bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 text-xs">
                                    {error}
                                </div>
                            )}

                            <div className="flex items-center gap-2">
                                <button
                                    type="button"
                                    onClick={onClose}
                                    className="px-3 py-2 text-sm rounded-lg border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                                >
                                    Cancelar
                                </button>
                                <button
                                    type="button"
                                    disabled={!file || uploading}
                                    onClick={handleUpload}
                                    className="flex-1 px-3 py-2 text-sm rounded-lg bg-[var(--color-primary)] text-white font-medium disabled:opacity-50 hover:opacity-90 transition-opacity"
                                >
                                    {uploading ? 'Procesando...' : 'Subir archivo'}
                                </button>
                            </div>
                        </>
                    ) : (
                        <>
                            <p className="text-xs text-amber-600">
                                Indica la variante de cada producto antes de confirmar. Se asignaran a la familia &quot;{family.name}&quot;.
                            </p>

                            {notFound.length > 0 && (
                                <div className="px-3 py-2 rounded-lg bg-amber-50 dark:bg-amber-900/20 text-amber-700 dark:text-amber-400 text-xs">
                                    SKUs sin coincidencia (no se procesaran): {notFound.join(', ')}
                                </div>
                            )}

                            <div className="border border-gray-200 dark:border-gray-700 rounded-lg divide-y divide-gray-100 dark:divide-gray-700 max-h-72 overflow-auto">
                                {matched.map((row) => (
                                    <div key={row.product_id} className="px-4 py-2.5 flex items-center justify-between gap-3">
                                        <div className="min-w-0">
                                            <div className="text-sm font-medium text-gray-900 dark:text-white truncate">{row.name}</div>
                                            <div className="text-xs text-gray-500 dark:text-gray-400">
                                                SKU: {row.sku}
                                                {row.current_family_name && row.current_family_id !== family.id && (
                                                    <span className="text-amber-600"> · se reasignara desde &quot;{row.current_family_name}&quot;</span>
                                                )}
                                            </div>
                                        </div>
                                        <input
                                            type="text"
                                            placeholder="Ej: Talla M, Rojo..."
                                            value={variantLabels[row.product_id] || ''}
                                            onChange={(e) => setVariantLabels((prev) => ({ ...prev, [row.product_id]: e.target.value }))}
                                            className="w-48 px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder:text-gray-400 focus:ring-2 focus:ring-[var(--color-primary)] focus:border-transparent"
                                        />
                                    </div>
                                ))}
                            </div>

                            {applying && (
                                <div className="text-xs text-gray-500">
                                    Aplicando... {progress}/{matched.length}
                                </div>
                            )}

                            {failedSkus.length > 0 && (
                                <div className="px-3 py-2 rounded-lg bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 text-xs">
                                    No se pudieron asociar: {failedSkus.join(', ')}
                                </div>
                            )}

                            <div className="flex items-center gap-2">
                                <button
                                    type="button"
                                    onClick={onClose}
                                    disabled={applying}
                                    className="px-3 py-2 text-sm rounded-lg border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors disabled:opacity-50"
                                >
                                    Cancelar
                                </button>
                                <button
                                    type="button"
                                    disabled={!allLabelsFilled || applying}
                                    onClick={handleConfirm}
                                    className="flex-1 px-3 py-2 text-sm rounded-lg bg-[var(--color-primary)] text-white font-medium disabled:opacity-50 hover:opacity-90 transition-opacity"
                                >
                                    {applying ? 'Asociando...' : `Confirmar asociación de ${matched.length} producto${matched.length !== 1 ? 's' : ''}`}
                                </button>
                            </div>
                        </>
                    )}
                </div>
            </div>
        </div>
    );
}
