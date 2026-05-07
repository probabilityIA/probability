'use client';

import { useState, useEffect } from 'react';
import { Product, CreateProductDTO, UpdateProductDTO, ProductFamily } from '../../domain/types';
import { createProductAction, updateProductAction, getProductFamiliesAction, getNextSKUAction, getNextSKUBatchAction } from '../../infra/actions';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { getActionError } from '@/shared/utils/action-result';
import { useDynamicBusinessColors } from '../hooks/useDynamicBusinessColors';
import { extractSKUPrefix } from '../helpers/sku-helper';

interface ProductFormProps {
    product?: Product;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

type FormMode = 'single' | 'batch';

interface VariantRow {
    localId: string;
    sku: string;
    name: string;
    attributes: Record<string, string>;
    status: 'pending' | 'creating' | 'done' | 'error';
    error?: string;
}

const ic = 'w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:border-transparent';
const lc = 'block text-xs font-semibold text-gray-700 dark:text-gray-200 mb-1';

export default function ProductForm({ product, onSuccess, onCancel, businessId }: ProductFormProps) {
    const { permissions } = usePermissions();
    const defaultBusinessId = businessId || permissions?.business_id || 0;
    const isEdit = !!product;
    const { colors: businessColors } = useDynamicBusinessColors(defaultBusinessId);
    const primaryColor = businessColors?.primary_color || '#a855f7';
    const secondaryColor = businessColors?.secondary_color || '#9f5cf7';
    const tertiaryColor = businessColors?.tertiary_color || '#c4b5fd';
    const quaternaryColor = businessColors?.quaternary_color || '#f3e8ff';

    const [mode, setMode] = useState<FormMode>('single');
    const [families, setFamilies] = useState<ProductFamily[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);

    const [formData, setFormData] = useState<CreateProductDTO>({
        business_id: product?.business_id || defaultBusinessId,
        sku: product?.sku || '',
        name: product?.name || '',
        description: product?.description || '',
        price: product?.price || 0,
        compare_at_price: product?.compare_at_price || undefined,
        cost_price: product?.cost_price || undefined,
        currency: product?.currency || 'COP',
        stock: product?.stock || 0,
        manage_stock: product?.manage_stock ?? true,
        track_inventory: product?.track_inventory ?? true,
        is_active: product?.is_active ?? true,
        status: product?.status || 'active',
        family_id: product?.family?.id || undefined,
        variant_label: product?.variant_label || '',
        variant_attributes: product?.variant_attributes || undefined,
        weight: product?.weight || undefined,
        height: product?.height || undefined,
        width: product?.width || undefined,
        length: product?.length || undefined,
    });

    const [skuPrefix, setSkuPrefix] = useState('PROD');
    const [skuCounter, setSkuCounter] = useState(1);
    const [selectedFamilyId, setSelectedFamilyId] = useState<number | null>(null);
    const [sharedPrice, setSharedPrice] = useState(0);
    const [sharedCurrency, setSharedCurrency] = useState('COP');
    const [sharedStatus, setSharedStatus] = useState('active');
    const [sharedTrackInventory, setSharedTrackInventory] = useState(true);
    const [variants, setVariants] = useState<VariantRow[]>([]);
    const [batchResults, setBatchResults] = useState<{ done: number; failed: number } | null>(null);
    const [skuSuggestions, setSkuSuggestions] = useState<string[]>([]);
    const [suggestedNextSKU, setSuggestedNextSKU] = useState<string | null>(null);
    const [showSKUSuggestions, setShowSKUSuggestions] = useState(false);
    const [skuSearchLoading, setSkuSearchLoading] = useState(false);

    useEffect(() => {
        const load = async () => {
            const res = await getProductFamiliesAction({ page: 1, page_size: 100, business_id: businessId });
            if (res.success && res.data) setFamilies(res.data as ProductFamily[]);
        };
        load();
    }, [businessId]);

    useEffect(() => {
        if (mode === 'single' && !isEdit) {
            const autoSku = `${skuPrefix}-${String(skuCounter).padStart(3, '0')}`;
            setFormData(f => ({ ...f, sku: autoSku }));
        }
    }, [mode, skuPrefix, skuCounter, isEdit]);

    useEffect(() => {
        if (!isEdit && formData.sku && formData.sku.length > 0) {
            setSkuSearchLoading(true);
            const timer = setTimeout(async () => {
                const prefix = extractSKUPrefix(formData.sku);
                const res = await getNextSKUAction(prefix, defaultBusinessId);
                if (res.success && res.data && res.data !== formData.sku) {
                    setSuggestedNextSKU(res.data);
                    setFormData(f => ({ ...f, sku: res.data }));
                } else if (res.success && res.data === formData.sku) {
                    setSuggestedNextSKU(null);
                } else {
                    setSuggestedNextSKU(null);
                }
                setSkuSearchLoading(false);
            }, 300);
            return () => clearTimeout(timer);
        } else {
            setSkuSuggestions([]);
            setSuggestedNextSKU(null);
        }
    }, [formData.sku, defaultBusinessId, isEdit]);

    const singleFamily = families.find(f => f.id === formData.family_id);
    const batchFamily = families.find(f => f.id === selectedFamilyId);
    const activeFamily = mode === 'single' ? singleFamily : batchFamily;
    const familyAxes: { key: string; label: string }[] = activeFamily?.variant_axes ?? [];

    const addVariant = async () => {
        setSkuSearchLoading(true);
        try {
            const newCount = variants.length + 1;
            const res = await getNextSKUBatchAction(skuPrefix, newCount, defaultBusinessId);

            if (res.success && res.data && res.data.length === newCount) {
                const newLocalId = Math.random().toString(36).slice(2);
                setVariants(prev => {
                    const updated = prev.map((v, i) => ({ ...v, sku: res.data[i] }));
                    updated.push({
                        localId: newLocalId,
                        sku: res.data[newCount - 1],
                        name: '',
                        attributes: {},
                        status: 'pending',
                    });
                    return updated;
                });
            } else {
                const newLocalId = Math.random().toString(36).slice(2);
                setVariants(prev => [...prev, {
                    localId: newLocalId,
                    sku: `${skuPrefix}-${String(prev.length + 1).padStart(3, '0')}`,
                    name: '',
                    attributes: {},
                    status: 'pending',
                }]);
            }
        } finally {
            setSkuSearchLoading(false);
        }
    };

    const removeVariant = (id: string) => setVariants(prev => prev.filter(v => v.localId !== id));

    const updateVariantField = (id: string, field: string, value: string) => {
        setVariants(prev => prev.map(v => {
            if (v.localId !== id) return v;
            if (field === 'sku' || field === 'name') return { ...v, [field]: value };
            return { ...v, attributes: { ...v.attributes, [field]: value } };
        }));
    };

    const applyPrefix = async (prefix: string) => {
        const p = prefix.toUpperCase().trim();
        setSkuPrefix(p);

        if (variants.length > 0) {
            setSkuSearchLoading(true);
            try {
                const res = await getNextSKUBatchAction(p, variants.length, defaultBusinessId);
                if (res.success && res.data && res.data.length === variants.length) {
                    setVariants(prev => prev.map((v, i) => ({ ...v, sku: res.data[i] })));
                } else {
                    setVariants(prev => prev.map((v, i) => ({ ...v, sku: `${p}-${String(i + 1).padStart(3, '0')}` })));
                }
            } finally {
                setSkuSearchLoading(false);
            }
        }
    };

    const handleSingleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        try {
            let res;
            if (isEdit) {
                const u: UpdateProductDTO = {
                    sku: formData.sku,
                    name: formData.name,
                    description: formData.description,
                    price: formData.price,
                    compare_at_price: formData.compare_at_price,
                    cost_price: formData.cost_price,
                    currency: formData.currency,
                    manage_stock: formData.manage_stock,
                    track_inventory: formData.track_inventory,
                    is_active: formData.is_active,
                    status: formData.status,
                    family_id: formData.family_id || undefined,
                    variant_label: formData.variant_label || undefined,
                    variant_attributes: formData.variant_attributes || undefined,
                    weight: formData.weight,
                    height: formData.height,
                    width: formData.width,
                    length: formData.length,
                };
                res = await updateProductAction(product!.id, u, businessId);
            } else {
                res = await createProductAction(formData, businessId);
                if (res.success) {
                    setSkuCounter(prev => prev + 1);
                }
            }
            if (res.success) onSuccess();
            else {
                if (res.message && res.message.toLowerCase().includes('sku')) {
                    setError(`El SKU "${formData.sku}" ya está en uso. Por favor, usa otro SKU.`);
                } else {
                    setError(res.message || 'Error al guardar');
                }
            }
        } catch (err: any) {
            const errorMsg = getActionError(err, 'Error al guardar');
            if (errorMsg.toLowerCase().includes('sku')) {
                setError(`El SKU "${formData.sku}" ya está en uso. Por favor, usa otro SKU.`);
            } else {
                setError(errorMsg);
            }
        } finally {
            setLoading(false);
        }
    };

    const handleBatchSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!variants.length) { setError('Agrega al menos una variante'); return; }
        if (variants.some(v => !v.sku.trim() || !v.name.trim())) { setError('Todas las variantes necesitan SKU y nombre'); return; }
        setLoading(true);
        setError(null);
        let done = 0;
        let failed = 0;
        setVariants(prev => prev.map(v => ({ ...v, status: 'creating' as const })));

        for (const v of variants) {
            const variantLabel = v.attributes.variant?.trim() || undefined;
            const hasAttrs = Object.values(v.attributes).some(val => val.trim());
            const variantAttributes = hasAttrs ? { ...v.attributes } : undefined;
            try {
                const res = await createProductAction({
                    business_id: defaultBusinessId,
                    sku: v.sku,
                    name: v.name,
                    price: sharedPrice,
                    currency: sharedCurrency,
                    status: sharedStatus,
                    is_active: sharedStatus === 'active',
                    track_inventory: sharedTrackInventory,
                    manage_stock: sharedTrackInventory,
                    family_id: selectedFamilyId || undefined,
                    variant_label: variantLabel,
                    variant_attributes: variantAttributes,
                }, businessId);

                if (res.success) {
                    done++;
                    setVariants(prev => prev.map(r => r.localId === v.localId ? { ...r, status: 'done' as const } : r));
                } else {
                    failed++;
                    setVariants(prev => prev.map(r => r.localId === v.localId ? { ...r, status: 'error' as const, error: res.message } : r));
                }
            } catch (err: any) {
                failed++;
                setVariants(prev => prev.map(r => r.localId === v.localId ? { ...r, status: 'error' as const, error: err.message } : r));
            }
        }

        setLoading(false);
        setBatchResults({ done, failed });
        if (failed === 0) setTimeout(() => onSuccess(), 1200);
    };

    const cancelBtn = (
        <button type="button" onClick={onCancel} className="px-5 py-2.5 text-sm font-medium text-gray-700 dark:text-gray-200 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
            Cancelar
        </button>
    );

    const submitBtn = (label: string, disabled = false) => (
        <button type="submit" disabled={loading || disabled} style={{ background: `linear-gradient(to right, ${primaryColor}, ${secondaryColor})` }} className="px-6 py-2.5 text-sm font-bold text-white rounded-lg shadow transition-all disabled:opacity-60 flex items-center gap-2">
            {loading && <span className="w-4 h-4 border-2 border-white/40 border-t-white rounded-full animate-spin" />}
            {label}
        </button>
    );

    return (
        <div className="space-y-5">
            {!isEdit && (
                <div className="flex gap-1 bg-gray-100 dark:bg-gray-900 p-1 rounded-xl">
                    {(['single', 'batch'] as FormMode[]).map(m => (
                        <button
                            key={m}
                            type="button"
                            onClick={() => { setMode(m); setError(null); }}
                            style={mode === m ? { color: primaryColor } : {}}
                            className={`flex-1 py-2 text-sm font-semibold rounded-lg transition-all ${mode === m ? 'bg-white dark:bg-gray-700 shadow-sm dark:opacity-90' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'}`}
                        >
                            {m === 'single' ? 'SKU individual' : 'Variantes en lote'}
                        </button>
                    ))}
                </div>
            )}

            {error && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-4 py-3 text-red-700 dark:text-red-400 text-sm flex justify-between items-start">
                    <span>{error}</span>
                    <button onClick={() => setError(null)} className="ml-3 opacity-60 hover:opacity-100 text-lg leading-none">&times;</button>
                </div>
            )}

            {mode === 'single' && (
                <form onSubmit={handleSingleSubmit} className="space-y-5">
                    <div className="grid grid-cols-2 gap-3">
                        <div>
                            <label className={lc}>SKU <span className="text-red-500">*</span></label>
                            <div className="relative">
                                <input
                                    className={ic}
                                    type="text"
                                    value={formData.sku}
                                    onChange={e => setFormData(f => ({ ...f, sku: e.target.value }))}
                                    disabled={isEdit}
                                    required
                                />
                                {skuSearchLoading && (
                                    <div className="absolute right-3 top-3 w-4 h-4 border-2 border-gray-300 border-t-current rounded-full animate-spin" style={{ color: primaryColor }} />
                                )}
                            </div>
                            {suggestedNextSKU && suggestedNextSKU !== formData.sku && !isEdit && (
                                <p style={{ color: primaryColor }} className="text-xs mt-1">Sugerencia: {suggestedNextSKU}</p>
                            )}
                        </div>
                        <div>
                            <label className={lc}>Nombre <span className="text-red-500">*</span></label>
                            <input className={ic} type="text" value={formData.name} onChange={e => setFormData(f => ({ ...f, name: e.target.value }))} required />
                        </div>
                        <div>
                            <label className={lc}>Precio <span className="text-red-500">*</span></label>
                            <input className={ic} type="number" min="0" step="0.01" value={formData.price} onChange={e => setFormData(f => ({ ...f, price: parseFloat(e.target.value) || 0 }))} required />
                        </div>
                        <div>
                            <label className={lc}>Moneda</label>
                            <select className={ic} value={formData.currency} onChange={e => setFormData(f => ({ ...f, currency: e.target.value }))}>
                                <option value="COP">COP</option>
                                <option value="USD">USD</option>
                                <option value="MXN">MXN</option>
                                <option value="EUR">EUR</option>
                            </select>
                        </div>
                        <div>
                            <label className={lc}>Estado</label>
                            <select className={ic} value={formData.status} onChange={e => setFormData(f => ({ ...f, status: e.target.value }))}>
                                <option value="active">Activo</option>
                                <option value="draft">Borrador</option>
                                <option value="archived">Archivado</option>
                            </select>
                        </div>
                        <div className="flex items-end pb-1">
                            <label className="flex items-center gap-2 cursor-pointer">
                                <input type="checkbox" checked={formData.track_inventory ?? true} onChange={e => setFormData(f => ({ ...f, track_inventory: e.target.checked, manage_stock: e.target.checked }))} style={{ accentColor: primaryColor }} className="w-4 h-4 rounded border-gray-300" />
                                <span className="text-sm text-gray-700 dark:text-gray-200">Gestionar inventario</span>
                            </label>
                        </div>
                    </div>

                    <div style={{ borderColor: tertiaryColor + '40', backgroundColor: primaryColor + '08' }} className="rounded-xl border p-4">
                        <h4 style={{ color: primaryColor }} className="text-xs font-bold uppercase tracking-wider mb-3">Familia de variantes (opcional)</h4>
                        <div className="grid grid-cols-2 gap-3">
                            <div>
                                <label className={lc}>Familia</label>
                                <select
                                    className={ic}
                                    value={formData.family_id ?? ''}
                                    onChange={e => {
                                        const id = e.target.value ? parseInt(e.target.value) : undefined;
                                        setFormData(f => ({ ...f, family_id: id, variant_attributes: undefined, variant_label: '' }));
                                    }}
                                >
                                    <option value="">Sin familia</option>
                                    {families.map(f => <option key={f.id} value={f.id}>{f.name}</option>)}
                                </select>
                            </div>
                            <div>
                                <label className={lc}>Etiqueta variante</label>
                                <input className={ic} type="text" placeholder="Ej: Vainilla - 1kg" value={formData.variant_label || ''} onChange={e => setFormData(f => ({ ...f, variant_label: e.target.value }))} />
                            </div>
                        </div>
                        {familyAxes.length > 0 && (
                            <div className="mt-3 grid grid-cols-2 sm:grid-cols-3 gap-2">
                                {familyAxes.map(ax => (
                                    <div key={ax.key}>
                                        <label className={lc}>{ax.label}</label>
                                        <input
                                            className={ic}
                                            type="text"
                                            placeholder={ax.label}
                                            value={formData.variant_attributes?.[ax.key] || ''}
                                            onChange={e => setFormData(f => ({
                                                ...f,
                                                variant_attributes: { ...f.variant_attributes, [ax.key]: e.target.value },
                                            }))}
                                        />
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    <details className="group">
                        <summary className="text-xs font-semibold text-gray-500 dark:text-gray-400 cursor-pointer select-none flex items-center gap-1 hover:text-gray-700 dark:hover:text-gray-200">
                            <svg className="w-3 h-3 transition-transform group-open:rotate-90" fill="currentColor" viewBox="0 0 20 20"><path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" /></svg>
                            Peso y dimensiones
                        </summary>
                        <div className="mt-3 grid grid-cols-2 sm:grid-cols-4 gap-3">
                            <div><label className={lc}>Peso (kg)</label><input className={ic} type="number" min="0" step="0.01" value={formData.weight ?? ''} onChange={e => setFormData(f => ({ ...f, weight: e.target.value ? parseFloat(e.target.value) : undefined }))} /></div>
                            <div><label className={lc}>Largo (cm)</label><input className={ic} type="number" min="0" step="0.1" value={formData.length ?? ''} onChange={e => setFormData(f => ({ ...f, length: e.target.value ? parseFloat(e.target.value) : undefined }))} /></div>
                            <div><label className={lc}>Ancho (cm)</label><input className={ic} type="number" min="0" step="0.1" value={formData.width ?? ''} onChange={e => setFormData(f => ({ ...f, width: e.target.value ? parseFloat(e.target.value) : undefined }))} /></div>
                            <div><label className={lc}>Alto (cm)</label><input className={ic} type="number" min="0" step="0.1" value={formData.height ?? ''} onChange={e => setFormData(f => ({ ...f, height: e.target.value ? parseFloat(e.target.value) : undefined }))} /></div>
                        </div>
                    </details>

                    <div>
                        <label className={lc}>Descripcion</label>
                        <textarea rows={3} className={ic} placeholder="Descripcion del producto..." value={formData.description || ''} onChange={e => setFormData(f => ({ ...f, description: e.target.value }))} />
                    </div>

                    <div className="flex justify-end gap-3 pt-3 border-t border-gray-200 dark:border-gray-700">
                        {cancelBtn}
                        {submitBtn(isEdit ? 'Guardar cambios' : 'Crear producto')}
                    </div>
                </form>
            )}

            {mode === 'batch' && !batchResults && (
                <form onSubmit={handleBatchSubmit} className="space-y-5">
                    <div className="grid grid-cols-2 gap-3">
                        <div>
                            <label className={lc}>Familia (opcional)</label>
                            <select
                                className={ic}
                                value={selectedFamilyId ?? ''}
                                onChange={e => {
                                    const id = e.target.value ? parseInt(e.target.value) : null;
                                    setSelectedFamilyId(id);
                                    setVariants([]);
                                }}
                            >
                                <option value="">Sin familia</option>
                                {families.map(f => <option key={f.id} value={f.id}>{f.name}</option>)}
                            </select>
                            {batchFamily && (
                                <p style={{ color: primaryColor }} className="text-xs mt-1">Las variantes se asociaran a esta familia</p>
                            )}
                        </div>
                        <div>
                            <label className={lc}>Prefijo SKU</label>
                            <input className={ic} type="text" value={skuPrefix} onChange={e => applyPrefix(e.target.value)} />
                            <p className="text-xs text-gray-400 mt-1">Genera: {skuPrefix || 'PROD'}-001, -002...</p>
                        </div>
                    </div>

                    <div className="grid grid-cols-3 gap-3">
                        <div>
                            <label className={lc}>Precio</label>
                            <input className={ic} type="number" min="0" step="0.01" value={sharedPrice} onChange={e => setSharedPrice(parseFloat(e.target.value) || 0)} />
                        </div>
                        <div>
                            <label className={lc}>Moneda</label>
                            <select className={ic} value={sharedCurrency} onChange={e => setSharedCurrency(e.target.value)}>
                                <option value="COP">COP</option>
                                <option value="USD">USD</option>
                                <option value="MXN">MXN</option>
                                <option value="EUR">EUR</option>
                            </select>
                        </div>
                        <div>
                            <label className={lc}>Estado</label>
                            <select className={ic} value={sharedStatus} onChange={e => setSharedStatus(e.target.value)}>
                                <option value="active">Activo</option>
                                <option value="draft">Borrador</option>
                                <option value="archived">Archivado</option>
                            </select>
                        </div>
                    </div>

                    <div>
                        <div className="flex items-center justify-between mb-2">
                            <label className="text-xs font-bold uppercase tracking-wider text-gray-700 dark:text-gray-200">
                                Variantes a crear {variants.length > 0 && <span style={{ color: primaryColor }}>({variants.length})</span>}
                            </label>
                            <button
                                type="button"
                                onClick={addVariant}
                                style={{ color: primaryColor }}
                                className="text-xs font-semibold flex items-center gap-1 hover:opacity-80 transition-opacity"
                            >
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" /></svg>
                                Agregar variante
                            </button>
                        </div>

                        {variants.length === 0 ? (
                            <div
                                onClick={addVariant}
                                style={{ borderColor: tertiaryColor + '40', backgroundColor: primaryColor + '08' }}
                                className="border-2 border-dashed rounded-xl p-8 text-center cursor-pointer transition-all hover:opacity-80"
                            >
                                <svg style={{ color: tertiaryColor + '80' }} className="w-8 h-8 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 4v16m8-8H4" /></svg>
                                <p className="text-sm text-gray-400">Haz click para agregar la primera variante</p>
                            </div>
                        ) : (
                            <div className="space-y-2">
                                        <div className="flex gap-2 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                                    <span style={{ width: 130 }}>SKU</span>
                                    <span className="flex-1">Nombre del producto</span>
                                    <span style={{ width: 140 }}>Variante</span>
                                    <span style={{ width: 20 }}></span>
                                </div>
                                {variants.map(v => (
                                    <div
                                        key={v.localId}
                                        className={`flex gap-2 items-center p-2 rounded-lg border transition-colors ${
                                            v.status === 'done' ? 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800' :
                                            v.status === 'error' ? 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800' :
                                            v.status === 'creating' ? 'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800' :
                                            'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700'
                                        }`}
                                    >
                                        <input
                                            style={{ width: 130 }}
                                            className={`${ic} text-xs font-mono flex-shrink-0`}
                                            value={v.sku}
                                            onChange={e => updateVariantField(v.localId, 'sku', e.target.value)}
                                            placeholder="SKU"
                                            disabled={v.status !== 'pending'}
                                        />
                                        <input
                                            className={`${ic} text-xs flex-1 min-w-0`}
                                            value={v.name}
                                            onChange={e => updateVariantField(v.localId, 'name', e.target.value)}
                                            placeholder="Nombre del producto"
                                            disabled={v.status !== 'pending'}
                                        />
                                        <input
                                            style={{ width: 140 }}
                                            className={`${ic} text-xs flex-shrink-0`}
                                            value={v.attributes.variant || ''}
                                            onChange={e => updateVariantField(v.localId, 'variant', e.target.value)}
                                            placeholder="Ej: Rojo XL, 500ml..."
                                            disabled={v.status !== 'pending'}
                                        />
                                        <div style={{ width: 20 }} className="flex-shrink-0 flex items-center justify-center">
                                            {v.status === 'done' && <span className="text-green-500 text-base">&#10003;</span>}
                                            {v.status === 'creating' && <span className="w-4 h-4 border-2 border-blue-300 border-t-blue-500 rounded-full animate-spin block" />}
                                            {v.status === 'error' && <span title={v.error} className="text-red-500 font-bold cursor-help">!</span>}
                                            {v.status === 'pending' && (
                                                <button type="button" onClick={() => removeVariant(v.localId)} className="text-gray-300 hover:text-red-500 transition-colors">
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
                                                </button>
                                            )}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    <div className="flex justify-end gap-3 pt-3 border-t border-gray-200 dark:border-gray-700">
                        {cancelBtn}
                        {submitBtn(
                            variants.length > 0 ? `Crear ${variants.length} variante${variants.length !== 1 ? 's' : ''}` : 'Crear variantes',
                            variants.length === 0
                        )}
                    </div>
                </form>
            )}

            {mode === 'batch' && batchResults && (
                <div className="text-center py-10 space-y-4">
                    <div className="text-5xl">{batchResults.failed === 0 ? '✅' : batchResults.done > 0 ? '⚠️' : '❌'}</div>
                    <p className="text-lg font-bold text-gray-900 dark:text-white">
                        {batchResults.done} creada{batchResults.done !== 1 ? 's' : ''}{batchResults.failed > 0 ? `, ${batchResults.failed} con error` : ''}
                    </p>
                    {batchResults.failed === 0 ? (
                        <p className="text-sm text-gray-400">Cerrando...</p>
                    ) : (
                        <div className="flex justify-center gap-3">
                            <button onClick={onSuccess} className="px-5 py-2.5 text-sm font-medium text-gray-700 dark:text-gray-200 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700">Ver listado</button>
                            <button
                                onClick={() => { setBatchResults(null); setVariants(prev => prev.filter(v => v.status === 'error').map(v => ({ ...v, status: 'pending' as const }))); }}
                                style={{ backgroundColor: primaryColor }}
                                className="px-5 py-2.5 text-sm font-bold text-white rounded-lg hover:opacity-90 transition-opacity"
                            >
                                Reintentar fallidas
                            </button>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}
