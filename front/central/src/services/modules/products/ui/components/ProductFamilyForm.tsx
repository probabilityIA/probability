'use client';

import { useState, useEffect } from 'react';
import { ProductFamily, CreateProductFamilyDTO, UpdateProductFamilyDTO } from '../../domain/types';
import { createProductFamilyAction, updateProductFamilyAction } from '../../infra/actions';

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
        title: '',
        description: '',
        slug: '',
        category: '',
        brand: '',
        image_url: '',
        status: 'active',
        is_active: true,
        variant_axes: '',
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (family) {
            setFormData({
                name: family.name || '',
                title: family.title || '',
                description: family.description || '',
                slug: family.slug || '',
                category: family.category || '',
                brand: family.brand || '',
                image_url: family.image_url || '',
                status: family.status || 'active',
                is_active: family.is_active ?? true,
                variant_axes: family.variant_axes ? JSON.stringify(family.variant_axes, null, 2) : '',
            });
        }
    }, [family]);

    const set = (field: string, value: any) => setFormData(prev => ({ ...prev, [field]: value }));

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
            if (isEdit) {
                const payload: UpdateProductFamilyDTO = {
                    name: formData.name || undefined,
                    title: formData.title || undefined,
                    description: formData.description || undefined,
                    slug: formData.slug || undefined,
                    category: formData.category || undefined,
                    brand: formData.brand || undefined,
                    image_url: formData.image_url || undefined,
                    status: formData.status || undefined,
                    is_active: formData.is_active,
                    variant_axes: variantAxesParsed,
                };
                res = await updateProductFamilyAction(family!.id, payload, businessId);
            } else {
                const payload: CreateProductFamilyDTO = {
                    name: formData.name,
                    title: formData.title || undefined,
                    description: formData.description || undefined,
                    slug: formData.slug || undefined,
                    category: formData.category || undefined,
                    brand: formData.brand || undefined,
                    image_url: formData.image_url || undefined,
                    status: formData.status || 'active',
                    is_active: formData.is_active,
                    variant_axes: variantAxesParsed,
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
        <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-4 py-3 text-red-700 dark:text-red-400 text-sm">
                    {error}
                </div>
            )}

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
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
                    <label className={labelClass}>Titulo</label>
                    <input
                        type="text"
                        placeholder="Titulo para el catalogo"
                        value={formData.title}
                        onChange={e => set('title', e.target.value)}
                        className={inputClass}
                    />
                </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <div>
                    <label className={labelClass}>Categoria</label>
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
            </div>

            <div>
                <label className={labelClass}>Descripcion</label>
                <textarea
                    rows={3}
                    placeholder="Descripcion de la familia de productos"
                    value={formData.description}
                    onChange={e => set('description', e.target.value)}
                    className={inputClass}
                />
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                    <label className={labelClass}>Estado</label>
                    <select value={formData.status} onChange={e => set('status', e.target.value)} className={inputClass}>
                        <option value="active">Activo</option>
                        <option value="draft">Borrador</option>
                        <option value="archived">Archivado</option>
                    </select>
                </div>
                <div>
                    <label className={labelClass}>URL de imagen</label>
                    <input
                        type="text"
                        placeholder="https://..."
                        value={formData.image_url}
                        onChange={e => set('image_url', e.target.value)}
                        className={inputClass}
                    />
                </div>
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

            <div>
                <label className={labelClass}>Ejes de variante (JSON opcional)</label>
                <textarea
                    rows={4}
                    placeholder={'[\n  {"key": "color", "label": "Color"},\n  {"key": "talla", "label": "Talla"}\n]'}
                    value={formData.variant_axes}
                    onChange={e => set('variant_axes', e.target.value)}
                    className={`${inputClass} font-mono text-xs`}
                />
                <p className="text-xs text-slate-400 mt-1">Define los ejes de variacion de esta familia. Ej: color, talla, sabor.</p>
            </div>

            <div className="flex justify-end gap-3 pt-2">
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
