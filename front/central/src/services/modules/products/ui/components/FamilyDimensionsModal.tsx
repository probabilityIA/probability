'use client';

import { useState } from 'react';
import { Product, ProductFamily } from '../../domain/types';
import { updateProductAction } from '../../infra/actions';
import { XMarkIcon } from '@heroicons/react/24/outline';

interface FamilyDimensionsModalProps {
    businessId?: number;
    family: ProductFamily;
    variants: Product[];
    onApplied: () => void;
    onClose: () => void;
}

const inputClass = "w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder:text-gray-400 focus:ring-2 focus:ring-[var(--color-primary)] focus:border-transparent";
const labelClass = "block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1";

function commonValue(variants: Product[], field: 'weight' | 'length' | 'width' | 'height'): number | undefined {
    if (variants.length === 0) return undefined;
    const first = variants[0][field];
    if (first == null) return undefined;
    const allSame = variants.every(v => v[field] === first);
    return allSame ? first : undefined;
}

export default function FamilyDimensionsModal({ businessId, family, variants, onApplied, onClose }: FamilyDimensionsModalProps) {
    const [weight, setWeight] = useState<number | undefined>(() => commonValue(variants, 'weight'));
    const [length, setLength] = useState<number | undefined>(() => commonValue(variants, 'length'));
    const [width, setWidth] = useState<number | undefined>(() => commonValue(variants, 'width'));
    const [height, setHeight] = useState<number | undefined>(() => commonValue(variants, 'height'));
    const [applying, setApplying] = useState(false);
    const [progress, setProgress] = useState(0);
    const [failedSkus, setFailedSkus] = useState<string[]>([]);

    const mixed = {
        weight: variants.length > 1 && commonValue(variants, 'weight') === undefined && variants.some(v => v.weight != null),
        length: variants.length > 1 && commonValue(variants, 'length') === undefined && variants.some(v => v.length != null),
        width: variants.length > 1 && commonValue(variants, 'width') === undefined && variants.some(v => v.width != null),
        height: variants.length > 1 && commonValue(variants, 'height') === undefined && variants.some(v => v.height != null),
    };

    const hasAnyValue = weight !== undefined || length !== undefined || width !== undefined || height !== undefined;

    const handleApply = async () => {
        if (!hasAnyValue || variants.length === 0) return;
        setApplying(true);
        setProgress(0);
        setFailedSkus([]);

        const failed: string[] = [];
        await Promise.all(variants.map(async (product) => {
            try {
                const res = await updateProductAction(
                    product.id,
                    { weight, length, width, height },
                    businessId
                );
                if (res.success === false) failed.push(product.sku);
            } catch {
                failed.push(product.sku);
            } finally {
                setProgress((p) => p + 1);
            }
        }));

        setApplying(false);
        if (failed.length > 0) {
            setFailedSkus(failed);
        } else {
            onApplied();
            onClose();
        }
    };

    return (
        <div className="fixed inset-0 z-[60] flex items-center justify-center p-4">
            <div className="absolute inset-0 bg-black/50" onClick={applying ? undefined : onClose} />
            <div className="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl w-full max-w-md flex flex-col">
                <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
                    <div>
                        <h3 className="text-base font-semibold text-gray-900 dark:text-white">Dimensiones de la familia</h3>
                        <p className="text-xs text-gray-500 mt-0.5">
                            se aplicarán a las {variants.length} variantes de &quot;{family.name}&quot;
                        </p>
                    </div>
                    <button
                        onClick={applying ? undefined : onClose}
                        disabled={applying}
                        className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors disabled:opacity-40"
                    >
                        <XMarkIcon className="w-5 h-5" />
                    </button>
                </div>

                <div className="p-5 flex flex-col gap-4">
                    <div className="grid grid-cols-2 gap-3">
                        <div>
                            <label className={labelClass}>Peso (kg)</label>
                            <input
                                type="number" min="0" step="0.01"
                                className={inputClass}
                                placeholder={mixed.weight ? 'Varios valores' : undefined}
                                value={weight ?? ''}
                                onChange={(e) => setWeight(e.target.value ? parseFloat(e.target.value) : undefined)}
                            />
                        </div>
                        <div>
                            <label className={labelClass}>Largo (cm)</label>
                            <input
                                type="number" min="0" step="0.1"
                                className={inputClass}
                                placeholder={mixed.length ? 'Varios valores' : undefined}
                                value={length ?? ''}
                                onChange={(e) => setLength(e.target.value ? parseFloat(e.target.value) : undefined)}
                            />
                        </div>
                        <div>
                            <label className={labelClass}>Ancho (cm)</label>
                            <input
                                type="number" min="0" step="0.1"
                                className={inputClass}
                                placeholder={mixed.width ? 'Varios valores' : undefined}
                                value={width ?? ''}
                                onChange={(e) => setWidth(e.target.value ? parseFloat(e.target.value) : undefined)}
                            />
                        </div>
                        <div>
                            <label className={labelClass}>Alto (cm)</label>
                            <input
                                type="number" min="0" step="0.1"
                                className={inputClass}
                                placeholder={mixed.height ? 'Varios valores' : undefined}
                                value={height ?? ''}
                                onChange={(e) => setHeight(e.target.value ? parseFloat(e.target.value) : undefined)}
                            />
                        </div>
                    </div>

                    <p className="text-xs text-amber-600">
                        Al aplicar, se sobrescriben las dimensiones actuales de cada variante de esta familia.
                    </p>

                    {applying && (
                        <div className="text-xs text-gray-500">
                            Aplicando... {progress}/{variants.length}
                        </div>
                    )}

                    {failedSkus.length > 0 && (
                        <div className="px-3 py-2 rounded-lg bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 text-xs">
                            No se pudieron actualizar: {failedSkus.join(', ')}
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
                            disabled={!hasAnyValue || applying || variants.length === 0}
                            onClick={handleApply}
                            className="flex-1 px-3 py-2 text-sm rounded-lg bg-[var(--color-primary)] text-white font-medium disabled:opacity-50 hover:opacity-90 transition-opacity"
                        >
                            {applying ? 'Aplicando...' : `Aplicar a ${variants.length} variante${variants.length !== 1 ? 's' : ''}`}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
