'use client';

import { useState, useEffect } from 'react';
import { ProductFamily, ProductFamilySummary, CreateProductFamilyDTO, UpdateProductFamilyDTO } from '../../domain/types';
import { createProductFamilyAction, updateProductFamilyAction, getProductFamiliesAction } from '../../infra/actions';

interface ProductFamilyFormProps {
    family?: ProductFamily;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

const inputClass = "w-full px-4 py-2.5 border-2 border-slate-200 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-[#7c3aed] focus:border-[#7c3aed] text-slate-900 dark:text-white placeholder:text-slate-400 bg-white dark:bg-gray-700 transition-all text-sm";
const labelClass = "block text-xs font-bold text-slate-700 dark:text-slate-200 mb-1.5";

export default function ProductFamilyForm({ family, onSuccess, onCancel, businessId }: ProductFamilyFormProps) {
    const isEdit = !!family;

    const [formData, setFormData] = useState({
        name: '',
        description: '',
        slug: '',
        category: '',
        brand: '',
        image_url: '',
        status: 'active',
        is_active: true,
        variant_axes: '',
        parent_family_id: '' as string | number,
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [candidateParents, setCandidateParents] = useState<ProductFamilySummary[]>([]);
    const [creatingParent, setCreatingParent] = useState(false);
    const [newParentName, setNewParentName] = useState('');
    const [creatingParentLoading, setCreatingParentLoading] = useState(false);
    const [creatingParentError, setCreatingParentError] = useState<string | null>(null);

    useEffect(() => {
        if (family) {
            setFormData({
                name: family.name || '',
                description: family.description || '',
                slug: family.slug || '',
                category: family.category || '',
                brand: family.brand || '',
                image_url: family.image_url || '',
                status: family.status || 'active',
                is_active: family.is_active ?? true,
                variant_axes: family.variant_axes ? JSON.stringify(family.variant_axes, null, 2) : '',
                parent_family_id: family.parent_family_id || '',
            });
        }
    }, [family]);

    useEffect(() => {
        (async () => {
            const res = await getProductFamiliesAction({ page: 1, page_size: 100, business_id: businessId });
            if (res && res.data) {
                setCandidateParents(res.data.filter((f: ProductFamilySummary) => f.id !== family?.id && !f.parent_family_id));
            }
        })();
    }, [businessId, family?.id]);

    const set = (field: string, value: any) => setFormData(prev => ({ ...prev, [field]: value }));

    const selectedParentPreview = formData.parent_family_id
        ? candidateParents.find(f => f.id === Number(formData.parent_family_id))
        : undefined;

    const handleParentSelectChange = (value: string) => {
        if (value === '__new__') {
            setCreatingParent(true);
            setNewParentName('');
            setCreatingParentError(null);
            return;
        }
        set('parent_family_id', value);
    };

    const handleCreateParent = async () => {
        if (!newParentName.trim()) {
            setCreatingParentError('El nombre es requerido');
            return;
        }
        setCreatingParentLoading(true);
        setCreatingParentError(null);
        try {
            const res: any = await createProductFamilyAction({ name: newParentName.trim(), status: 'active', is_active: true }, businessId);
            if (res && res.success === false) {
                setCreatingParentError(res.message || 'Error al crear la familia madre');
                return;
            }
            const created = res.data;
            setCandidateParents(prev => [...prev, created]);
            set('parent_family_id', created.id);
            setCreatingParent(false);
        } catch (err: any) {
            setCreatingParentError(err.message || 'Error inesperado');
        } finally {
            setCreatingParentLoading(false);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!formData.name.trim()) {
            setError('El nombre es requerido');
            return;
        }
        setLoading(true);
        setError(null);

        let variantAxesParsed: any = undefined;
        if (formData.variant_axes.trim()) {
            try {
                variantAxesParsed = JSON.parse(formData.variant_axes);
            } catch {
                setError('Los ejes de variante deben ser JSON valido');
                setLoading(false);
                return;
            }
        }

        try {
            let res: any;
            const parentFamilyId = formData.parent_family_id ? Number(formData.parent_family_id) : undefined;

            if (isEdit) {
                const payload: UpdateProductFamilyDTO = {
                    name: formData.name || undefined,
                    description: formData.description || undefined,
                    slug: formData.slug || undefined,
                    category: formData.category || undefined,
                    brand: formData.brand || undefined,
                    image_url: formData.image_url || undefined,
                    status: formData.status || undefined,
                    is_active: formData.is_active,
                    variant_axes: variantAxesParsed,
                    parent_family_id: parentFamilyId,
                    clear_parent: !parentFamilyId && !!family?.parent_family_id,
                };
                res = await updateProductFamilyAction(family!.id, payload, businessId);
            } else {
                const payload: CreateProductFamilyDTO = {
                    name: formData.name,
                    description: formData.description || undefined,
                    slug: formData.slug || undefined,
                    category: formData.category || undefined,
                    brand: formData.brand || undefined,
                    image_url: formData.image_url || undefined,
                    status: formData.status || 'active',
                    is_active: formData.is_active,
                    variant_axes: variantAxesParsed,
                    parent_family_id: parentFamilyId,
                };
                res = await createProductFamilyAction(payload, businessId);
            }

            if (res && res.success === false) {
                setError(res.message || 'Error al guardar familia');
            } else {
                onSuccess();
            }
        } catch (err: any) {
            setError(err.message || 'Error inesperado');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="grid grid-cols-1 lg:grid-cols-2 gap-x-8 gap-y-5">
            {error && (
                <div className="lg:col-span-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-4 py-3 text-red-700 dark:text-red-400 text-sm">
                    {error}
                </div>
            )}

            <div className="space-y-5">
                <div className="grid grid-cols-2 gap-3">
                    <div className="col-span-2">
                        <label className={labelClass}>Nombre *</label>
                        <input
                            type="text"
                            placeholder="Ej: Tenis Runner"
                            value={formData.name}
                            onChange={e => set('name', e.target.value)}
                            className={inputClass}
                            required
                        />
                    </div>
                    <div>
                        <label className={labelClass}>Categoría</label>
                        <input
                            type="text"
                            placeholder="Ej: Calzado"
                            value={formData.category}
                            onChange={e => set('category', e.target.value)}
                            className={inputClass}
                        />
                    </div>
                    <div>
                        <label className={labelClass}>Marca</label>
                        <input
                            type="text"
                            placeholder="Ej: Nike"
                            value={formData.brand}
                            onChange={e => set('brand', e.target.value)}
                            className={inputClass}
                        />
                    </div>
                    <div>
                        <label className={labelClass}>Slug</label>
                        <input
                            type="text"
                            placeholder="tenis-runner"
                            value={formData.slug}
                            onChange={e => set('slug', e.target.value)}
                            className={inputClass}
                        />
                    </div>
                    <div>
                        <label className={labelClass}>Estado</label>
                        <select value={formData.status} onChange={e => set('status', e.target.value)} className={inputClass}>
                            <option value="active">Activo</option>
                            <option value="draft">Borrador</option>
                            <option value="archived">Archivado</option>
                        </select>
                    </div>
                </div>

                <div className="p-3.5 rounded-lg border-2 border-purple-200 dark:border-purple-800 bg-purple-50/60 dark:bg-purple-900/10">
                    <div className="flex items-center gap-1.5 mb-1.5">
                        <span className="text-base leading-none">&#128081;</span>
                        <label className="text-xs font-bold text-purple-800 dark:text-purple-300">Agrupar bajo una familia madre (opcional)</label>
                    </div>
                    {creatingParent ? (
                        <div className="flex flex-col gap-2">
                            <div className="flex gap-2">
                                <input
                                    type="text"
                                    autoFocus
                                    placeholder="Nombre de la nueva familia madre (ej: Camisas Selección Argentina)"
                                    value={newParentName}
                                    onChange={e => setNewParentName(e.target.value)}
                                    className={inputClass}
                                />
                                <button
                                    type="button"
                                    onClick={handleCreateParent}
                                    disabled={creatingParentLoading}
                                    className="px-4 py-2 text-sm font-bold text-white bg-purple-600 hover:bg-purple-700 rounded-lg transition-colors disabled:opacity-60 whitespace-nowrap"
                                >
                                    {creatingParentLoading ? 'Creando...' : 'Crear'}
                                </button>
                                <button
                                    type="button"
                                    onClick={() => { setCreatingParent(false); setCreatingParentError(null); }}
                                    className="px-3 py-2 text-sm font-medium text-slate-500 hover:text-slate-700"
                                >
                                    Cancelar
                                </button>
                            </div>
                            {creatingParentError && <p className="text-xs text-red-600">{creatingParentError}</p>}
                        </div>
                    ) : (
                        <select value={formData.parent_family_id} onChange={e => handleParentSelectChange(e.target.value)} className={inputClass}>
                            <option value="">No agrupar - esta familia es independiente</option>
                            {candidateParents.map(f => (
                                <option key={f.id} value={f.id}>{f.name}</option>
                            ))}
                            <option value="__new__">+ Crear nueva familia madre...</option>
                        </select>
                    )}
                    {!creatingParent && (selectedParentPreview ? (
                        <div className="flex items-center gap-2 mt-2 p-2 rounded-md bg-white dark:bg-gray-800 border border-purple-200 dark:border-purple-800">
                            {selectedParentPreview.image_url ? (
                                <img src={selectedParentPreview.image_url} alt={selectedParentPreview.name} className="w-8 h-8 rounded-full object-cover" />
                            ) : (
                                <span className="w-8 h-8 rounded-full bg-purple-100 dark:bg-purple-900/40" />
                            )}
                            <p className="text-xs text-slate-600 dark:text-slate-300">
                                Esta familia se sincronizará como parte de <strong>{selectedParentPreview.name}</strong>: al enviarla a un
                                canal de ventas (ej. WooCommerce), se fusiona con sus demás subfamilias en <strong>un solo producto</strong> con
                                todas sus variantes.
                            </p>
                        </div>
                    ) : (
                        <p className="text-xs text-slate-400 mt-1.5">
                            Úsalo cuando el mismo producto tenga varias líneas separadas en Probability (ej: distintos colores de
                            una camiseta creados como familias distintas) pero deban verse como un solo producto en el canal de venta.
                        </p>
                    ))}
                </div>

                <div>
                    <label className={labelClass}>Descripción</label>
                    <textarea
                        rows={5}
                        placeholder="Descripción de la familia de productos"
                        value={formData.description}
                        onChange={e => set('description', e.target.value)}
                        className={inputClass}
                    />
                </div>

                <div className="flex items-center gap-3 p-3 bg-slate-50 dark:bg-gray-700/30 rounded-lg border border-slate-200 dark:border-gray-600">
                    <input
                        type="checkbox"
                        id="family-is-active"
                        checked={formData.is_active}
                        onChange={e => set('is_active', e.target.checked)}
                        className="w-4 h-4 text-[#7c3aed] border-gray-300 rounded focus:ring-[#7c3aed]"
                    />
                    <label htmlFor="family-is-active" className="text-sm font-medium text-slate-700 dark:text-slate-200 cursor-pointer">
                        Familia activa
                    </label>
                </div>
            </div>

            <div className="space-y-5">
                <div>
                    <label className={labelClass}>URL de imagen</label>
                    <input
                        type="text"
                        placeholder="https://..."
                        value={formData.image_url}
                        onChange={e => set('image_url', e.target.value)}
                        className={inputClass}
                    />
                    {formData.image_url && (
                        <div className="rounded-lg overflow-hidden border border-slate-200 dark:border-gray-600 flex items-center justify-center p-3 mt-2" style={{ minHeight: 140 }}>
                            <img src={formData.image_url} alt="Vista previa" className="max-h-32 max-w-full object-contain" onError={e => (e.currentTarget.style.display = 'none')} />
                        </div>
                    )}
                </div>

                <div>
                    <label className={labelClass}>Ejes de variante (JSON opcional)</label>
                    <textarea
                        rows={6}
                        placeholder={'[\n  {"key": "color", "label": "Color"},\n  {"key": "talla", "label": "Talla"}\n]'}
                        value={formData.variant_axes}
                        onChange={e => set('variant_axes', e.target.value)}
                        className={`${inputClass} font-mono text-xs`}
                    />
                    <p className="text-xs text-slate-400 mt-1">Define los ejes de variación de esta familia. Ej: color, talla, sabor.</p>
                </div>
            </div>

            <div className="lg:col-span-2 flex justify-end gap-3 pt-3 border-t border-gray-200 dark:border-gray-700">
                <button
                    type="button"
                    onClick={onCancel}
                    disabled={loading}
                    className="px-5 py-2.5 text-sm font-medium text-slate-700 dark:text-slate-200 bg-white dark:bg-gray-700 border-2 border-slate-200 dark:border-gray-600 rounded-lg hover:border-slate-300 transition-all"
                >
                    Cancelar
                </button>
                <button
                    type="submit"
                    disabled={loading}
                    className="px-6 py-2.5 text-sm font-bold text-white bg-gradient-to-r from-[#7c3aed] to-[#6d28d9] hover:from-[#6d28d9] hover:to-[#5b21b6] rounded-lg shadow-lg hover:shadow-xl transition-all disabled:opacity-60 disabled:cursor-not-allowed flex items-center gap-2"
                >
                    {loading && <span className="w-4 h-4 border-2 border-white/40 border-t-white rounded-full animate-spin inline-block" />}
                    {isEdit ? 'Guardar cambios' : 'Crear familia'}
                </button>
            </div>
        </form>
    );
}
