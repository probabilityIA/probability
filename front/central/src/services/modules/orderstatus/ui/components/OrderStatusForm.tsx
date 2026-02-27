'use client';

import { useState, useEffect } from 'react';
import { OrderStatusInfo, CreateOrderStatusDTO, UpdateOrderStatusDTO } from '../../domain/types';
import { createOrderStatusAction, updateOrderStatusAction } from '../../infra/actions';

interface OrderStatusFormProps {
    status?: OrderStatusInfo;
    onSuccess: (status: OrderStatusInfo) => void;
    onCancel: () => void;
}

const CATEGORIES = ['pending', 'processing', 'shipped', 'delivered', 'completed', 'cancelled', 'other'];

export default function OrderStatusForm({ status, onSuccess, onCancel }: OrderStatusFormProps) {
    const isEditing = !!status;

    const [form, setForm] = useState({
        code: status?.code ?? '',
        name: status?.name ?? '',
        description: status?.description ?? '',
        category: status?.category ?? '',
        color: status?.color ?? '#6B7280',
        priority: status?.priority ?? 0,
        is_active: status?.is_active ?? true,
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (status) {
            setForm({
                code: status.code,
                name: status.name,
                description: status.description ?? '',
                category: status.category ?? '',
                color: status.color ?? '#6B7280',
                priority: status.priority ?? 0,
                is_active: status.is_active ?? true,
            });
        }
    }, [status]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            const payload: CreateOrderStatusDTO | UpdateOrderStatusDTO = {
                code: form.code,
                name: form.name,
                description: form.description || undefined,
                category: form.category || undefined,
                color: form.color || undefined,
                priority: form.priority,
                is_active: form.is_active,
            };

            let result;
            if (isEditing && status) {
                result = await updateOrderStatusAction(status.id, payload as UpdateOrderStatusDTO);
            } else {
                result = await createOrderStatusAction(payload as CreateOrderStatusDTO);
            }

            if (!result.success) {
                setError(result.message || 'Error al guardar el estado');
                return;
            }

            onSuccess(result.data);
        } catch (err: any) {
            setError(err.message || 'Error inesperado');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
                    {error}
                </div>
            )}

            <div className="grid grid-cols-2 gap-4">
                {/* Código */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Código <span className="text-red-500">*</span>
                    </label>
                    <input
                        type="text"
                        value={form.code}
                        onChange={e => setForm(f => ({ ...f, code: e.target.value }))}
                        className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono"
                        placeholder="ej. pending"
                        required
                    />
                </div>

                {/* Nombre */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Nombre <span className="text-red-500">*</span>
                    </label>
                    <input
                        type="text"
                        value={form.name}
                        onChange={e => setForm(f => ({ ...f, name: e.target.value }))}
                        className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                        placeholder="ej. Pendiente"
                        required
                    />
                </div>
            </div>

            {/* Descripción */}
            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                    Descripción
                </label>
                <textarea
                    value={form.description}
                    onChange={e => setForm(f => ({ ...f, description: e.target.value }))}
                    className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    rows={2}
                    placeholder="Descripción del estado..."
                />
            </div>

            <div className="grid grid-cols-2 gap-4">
                {/* Categoría */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Categoría
                    </label>
                    <select
                        value={form.category}
                        onChange={e => setForm(f => ({ ...f, category: e.target.value }))}
                        className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                        <option value="">Sin categoría</option>
                        {CATEGORIES.map(cat => (
                            <option key={cat} value={cat}>{cat}</option>
                        ))}
                    </select>
                </div>

                {/* Prioridad */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Prioridad
                    </label>
                    <input
                        type="number"
                        value={form.priority}
                        onChange={e => setForm(f => ({ ...f, priority: parseInt(e.target.value) || 0 }))}
                        className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                        min={0}
                        placeholder="0"
                    />
                </div>
            </div>

            <div className="grid grid-cols-2 gap-4 items-end">
                {/* Color */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Color
                    </label>
                    <div className="flex items-center gap-2">
                        <input
                            type="color"
                            value={form.color}
                            onChange={e => setForm(f => ({ ...f, color: e.target.value }))}
                            className="w-10 h-10 rounded-lg border border-gray-300 cursor-pointer p-0.5"
                        />
                        <input
                            type="text"
                            value={form.color}
                            onChange={e => setForm(f => ({ ...f, color: e.target.value }))}
                            className="flex-1 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono"
                            placeholder="#6B7280"
                        />
                        <span
                            className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium"
                            style={{
                                backgroundColor: form.color + '20',
                                color: form.color,
                                border: `1px solid ${form.color}40`,
                            }}
                        >
                            Preview
                        </span>
                    </div>
                </div>

                {/* Is Active */}
                <div>
                    <label className="flex items-center gap-2 cursor-pointer">
                        <input
                            type="checkbox"
                            checked={form.is_active}
                            onChange={e => setForm(f => ({ ...f, is_active: e.target.checked }))}
                            className="w-4 h-4 rounded text-blue-600"
                        />
                        <span className="text-sm font-medium text-gray-700">Activo</span>
                    </label>
                </div>
            </div>

            {/* Actions */}
            <div className="flex justify-end gap-3 pt-2">
                <button
                    type="button"
                    onClick={onCancel}
                    disabled={loading}
                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50"
                >
                    Cancelar
                </button>
                <button
                    type="submit"
                    disabled={loading}
                    className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 disabled:opacity-50"
                >
                    {loading ? 'Guardando...' : isEditing ? 'Actualizar' : 'Crear'}
                </button>
            </div>
        </form>
    );
}
