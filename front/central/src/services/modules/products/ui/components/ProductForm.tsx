'use client';

import { useState, useRef } from 'react';
import { Product, CreateProductDTO, UpdateProductDTO } from '../../domain/types';
import { createProductAction, updateProductAction, uploadProductImageAction } from '../../infra/actions';
import { Button, Alert, Input, Select } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { getActionError } from '@/shared/utils/action-result';

interface ProductFormProps {
    product?: Product;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

export default function ProductForm({ product, onSuccess, onCancel, businessId }: ProductFormProps) {
    const { permissions } = usePermissions();
    const defaultBusinessId = businessId || permissions?.business_id || 0;

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
        weight: product?.weight || undefined,
        height: product?.height || undefined,
        width: product?.width || undefined,
        length: product?.length || undefined,
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [imagePreview, setImagePreview] = useState<string | null>(product?.image_url || null);
    const [uploadingImage, setUploadingImage] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            let response;
            if (product) {
                const updateData: UpdateProductDTO = {
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
                    weight: formData.weight,
                    height: formData.height,
                    width: formData.width,
                    length: formData.length,
                };
                response = await updateProductAction(product.id, updateData, businessId);
            } else {
                response = await createProductAction(formData, businessId);
            }

            if (response.success) {
                setSuccess(product ? 'Producto actualizado exitosamente' : 'Producto creado exitosamente');
                setTimeout(() => {
                    onSuccess();
                }, 1000);
            } else {
                setError(response.message || 'Error al guardar el producto');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al guardar el producto'));
        } finally {
            setLoading(false);
        }
    };

    const handleChange = (field: keyof CreateProductDTO, value: any) => {
        setFormData({ ...formData, [field]: value });
    };

    const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file || !product) return;

        const reader = new FileReader();
        reader.onloadend = () => setImagePreview(reader.result as string);
        reader.readAsDataURL(file);

        setUploadingImage(true);
        setError(null);
        try {
            const fd = new FormData();
            fd.append('image', file);
            const result = await uploadProductImageAction(product.id, fd, businessId);
            if (result.success) {
                setImagePreview(result.image_url);
                setSuccess('Imagen subida exitosamente');
                setTimeout(() => setSuccess(null), 3000);
            } else {
                setError(result.message || 'Error al subir imagen');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al subir imagen'));
        } finally {
            setUploadingImage(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}
            {success && (
                <Alert type="success" onClose={() => setSuccess(null)}>
                    {success}
                </Alert>
            )}

            {/* Informacion basica */}
            <fieldset>
                <legend className="text-sm font-semibold text-gray-900 dark:text-white mb-3">Informacion basica</legend>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            SKU <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="text"
                            value={formData.sku}
                            onChange={(e) => handleChange('sku', e.target.value)}
                            required
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Nombre <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="text"
                            value={formData.name}
                            onChange={(e) => handleChange('name', e.target.value)}
                            required
                        />
                    </div>
                </div>
            </fieldset>

            {/* Precios */}
            <fieldset>
                <legend className="text-sm font-semibold text-gray-900 dark:text-white mb-3">Precios</legend>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Precio <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="number"
                            value={formData.price}
                            onChange={(e) => handleChange('price', parseFloat(e.target.value) || 0)}
                            required
                            min="0"
                            step="0.01"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Precio comparacion
                        </label>
                        <Input
                            type="number"
                            value={formData.compare_at_price ?? ''}
                            onChange={(e) => handleChange('compare_at_price', e.target.value ? parseFloat(e.target.value) : undefined)}
                            min="0"
                            step="0.01"
                            placeholder="Precio antes del descuento"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Costo
                        </label>
                        <Input
                            type="number"
                            value={formData.cost_price ?? ''}
                            onChange={(e) => handleChange('cost_price', e.target.value ? parseFloat(e.target.value) : undefined)}
                            min="0"
                            step="0.01"
                            placeholder="Costo del producto"
                        />
                    </div>
                </div>
                <div className="mt-4 w-full max-w-[200px]">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Moneda
                    </label>
                    <Select
                        value={formData.currency}
                        onChange={(e) => handleChange('currency', e.target.value)}
                        options={[
                            { value: 'COP', label: 'COP' },
                            { value: 'USD', label: 'USD' },
                            { value: 'MXN', label: 'MXN' },
                            { value: 'EUR', label: 'EUR' },
                        ]}
                    />
                </div>
            </fieldset>

            {/* Estado */}
            <fieldset>
                <legend className="text-sm font-semibold text-gray-900 dark:text-white mb-3">Estado</legend>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Estado
                        </label>
                        <Select
                            value={formData.status}
                            onChange={(e) => handleChange('status', e.target.value)}
                            options={[
                                { value: 'active', label: 'Activo' },
                                { value: 'draft', label: 'Borrador' },
                                { value: 'archived', label: 'Archivado' },
                            ]}
                        />
                    </div>
                    <div className="flex items-end pb-1">
                        <label className="flex items-center gap-2 cursor-pointer">
                            <input
                                type="checkbox"
                                checked={formData.track_inventory ?? true}
                                onChange={(e) => handleChange('track_inventory', e.target.checked)}
                                className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                            />
                            <span className="text-sm text-gray-700 dark:text-gray-200">Gestionar inventario</span>
                        </label>
                    </div>
                </div>
            </fieldset>

            {/* Peso y dimensiones */}
            <fieldset>
                <legend className="text-sm font-semibold text-gray-900 dark:text-white mb-3">Peso y dimensiones</legend>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Peso (kg)</label>
                        <Input
                            type="number"
                            value={formData.weight ?? ''}
                            onChange={(e) => handleChange('weight', e.target.value ? parseFloat(e.target.value) : undefined)}
                            min="0"
                            step="0.01"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Largo (cm)</label>
                        <Input
                            type="number"
                            value={formData.length ?? ''}
                            onChange={(e) => handleChange('length', e.target.value ? parseFloat(e.target.value) : undefined)}
                            min="0"
                            step="0.1"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Ancho (cm)</label>
                        <Input
                            type="number"
                            value={formData.width ?? ''}
                            onChange={(e) => handleChange('width', e.target.value ? parseFloat(e.target.value) : undefined)}
                            min="0"
                            step="0.1"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">Alto (cm)</label>
                        <Input
                            type="number"
                            value={formData.height ?? ''}
                            onChange={(e) => handleChange('height', e.target.value ? parseFloat(e.target.value) : undefined)}
                            min="0"
                            step="0.1"
                        />
                    </div>
                </div>
            </fieldset>

            {/* Imagen (solo en modo edicion) */}
            {product && (
                <fieldset>
                    <legend className="text-sm font-semibold text-gray-900 dark:text-white mb-3">Imagen del producto</legend>
                    <div className="flex items-start gap-4">
                        {imagePreview ? (
                            <div className="relative w-32 h-32 rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700">
                                <img
                                    src={imagePreview}
                                    alt={product.name}
                                    className="w-full h-full object-cover"
                                />
                            </div>
                        ) : (
                            <div className="w-32 h-32 rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-600 flex items-center justify-center text-gray-400 text-sm text-center">
                                Sin imagen
                            </div>
                        )}
                        <div className="flex flex-col gap-2">
                            <input
                                ref={fileInputRef}
                                type="file"
                                accept="image/jpeg,image/png,image/gif,image/webp"
                                onChange={handleImageUpload}
                                className="hidden"
                            />
                            <Button
                                type="button"
                                variant="outline"
                                onClick={() => fileInputRef.current?.click()}
                                disabled={uploadingImage}
                            >
                                {uploadingImage ? 'Subiendo...' : 'Cambiar imagen'}
                            </Button>
                            <p className="text-xs text-gray-500 dark:text-gray-400">
                                JPG, PNG, GIF o WebP. Max 10MB.
                            </p>
                        </div>
                    </div>
                </fieldset>
            )}

            {/* Descripcion */}
            <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                    Descripcion
                </label>
                <textarea
                    value={formData.description}
                    onChange={(e) => handleChange('description', e.target.value)}
                    rows={3}
                    className="w-full px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border border-gray-300 dark:border-gray-600 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Descripcion del producto..."
                />
            </div>

            {/* Actions */}
            <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
                    Cancelar
                </Button>
                <Button type="submit" variant="primary" disabled={loading}>
                    {loading ? 'Guardando...' : product ? 'Actualizar' : 'Crear'}
                </Button>
            </div>
        </form>
    );
}
