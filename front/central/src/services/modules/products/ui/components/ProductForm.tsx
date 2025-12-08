'use client';

import { useState } from 'react';
import { Product, CreateProductDTO, UpdateProductDTO } from '../../domain/types';
import { createProductAction, updateProductAction } from '../../infra/actions';
import { Button, Alert, Input, Select } from '@/shared/ui';

interface ProductFormProps {
    product?: Product;
    onSuccess: () => void;
    onCancel: () => void;
}

export default function ProductForm({ product, onSuccess, onCancel }: ProductFormProps) {
    const [formData, setFormData] = useState<CreateProductDTO>({
        business_id: product?.business_id || 0,
        sku: product?.sku || '',
        name: product?.name || '',
        description: product?.description || '',
        price: product?.price || 0,
        currency: product?.currency || 'USD',
        stock: product?.stock || 0,
        manage_stock: product?.manage_stock ?? true,
        is_active: product?.is_active ?? true,
        status: product?.status || 'active',
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            let response;
            if (product) {
                // Update
                const updateData: UpdateProductDTO = {
                    sku: formData.sku,
                    name: formData.name,
                    description: formData.description,
                    price: formData.price,
                    currency: formData.currency,
                    stock: formData.stock,
                    manage_stock: formData.manage_stock,
                    is_active: formData.is_active,
                    status: formData.status,
                };
                response = await updateProductAction(product.id, updateData);
            } else {
                // Create
                // Note: business_id should ideally come from context or selection if creating new
                response = await createProductAction(formData);
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
            setError(err.message || 'Error al guardar el producto');
        } finally {
            setLoading(false);
        }
    };

    const handleChange = (field: keyof CreateProductDTO, value: any) => {
        setFormData({ ...formData, [field]: value });
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

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* SKU */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        SKU <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.sku}
                        onChange={(e) => handleChange('sku', e.target.value)}
                        required
                    />
                </div>

                {/* Name */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Nombre <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.name}
                        onChange={(e) => handleChange('name', e.target.value)}
                        required
                    />
                </div>

                {/* Price */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
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

                {/* Currency */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Moneda
                    </label>
                    <Select
                        value={formData.currency}
                        onChange={(e) => handleChange('currency', e.target.value)}
                        options={[
                            { value: 'USD', label: 'USD' },
                            { value: 'COP', label: 'COP' },
                            { value: 'MXN', label: 'MXN' },
                        ]}
                    />
                </div>

                {/* Stock */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Stock
                    </label>
                    <Input
                        type="number"
                        value={formData.stock}
                        onChange={(e) => handleChange('stock', parseInt(e.target.value) || 0)}
                        min="0"
                    />
                </div>

                {/* Status */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
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
            </div>

            {/* Description */}
            <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                    Descripción
                </label>
                <textarea
                    value={formData.description}
                    onChange={(e) => handleChange('description', e.target.value)}
                    rows={3}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
            </div>

            {/* Actions */}
            <div className="flex justify-end gap-3 pt-4 border-t">
                <Button
                    type="button"
                    variant="outline"
                    onClick={onCancel}
                    disabled={loading}
                >
                    Cancelar
                </Button>
                <Button
                    type="submit"
                    variant="primary"
                    disabled={loading}
                >
                    {loading ? 'Guardando...' : product ? 'Actualizar' : 'Crear'}
                </Button>
            </div>
        </form>
    );
}
