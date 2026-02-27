'use client';

import { useState } from 'react';
import { AdjustStockDTO } from '../../domain/types';
import { adjustStockAction } from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';

interface AdjustStockModalProps {
    warehouseId: number;
    businessId?: number;
    onSuccess: () => void;
    onClose: () => void;
}

export default function AdjustStockModal({ warehouseId, businessId, onSuccess, onClose }: AdjustStockModalProps) {
    const [formData, setFormData] = useState({
        product_id: '',
        quantity: 0,
        reason: '',
        notes: '',
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
            const dto: AdjustStockDTO = {
                product_id: formData.product_id.trim(),
                warehouse_id: warehouseId,
                quantity: formData.quantity,
                reason: formData.reason.trim(),
                notes: formData.notes.trim() || undefined,
            };
            await adjustStockAction(dto, businessId);
            setSuccess('Stock ajustado exitosamente');
            setTimeout(() => onSuccess(), 800);
        } catch (err: any) {
            setError(err.message || 'Error al ajustar stock');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
            <div className="bg-white rounded-xl shadow-xl w-full max-w-md max-h-[90vh] overflow-y-auto">
                <div className="flex items-center justify-between px-6 py-4 border-b">
                    <h2 className="text-lg font-semibold text-gray-900">Ajustar stock</h2>
                    <button onClick={onClose} className="text-gray-400 hover:text-gray-600 text-xl leading-none">
                        &times;
                    </button>
                </div>
                <form onSubmit={handleSubmit} className="p-6 space-y-4">
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

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            ID del producto <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="text"
                            value={formData.product_id}
                            onChange={(e) => setFormData(prev => ({ ...prev, product_id: e.target.value }))}
                            placeholder="UUID del producto"
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Cantidad <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="number"
                            value={formData.quantity.toString()}
                            onChange={(e) => setFormData(prev => ({ ...prev, quantity: parseInt(e.target.value) || 0 }))}
                            placeholder="Positivo para agregar, negativo para quitar"
                            required
                        />
                        <p className="text-xs text-gray-500 mt-1">
                            Use valores positivos para agregar stock y negativos para quitar.
                        </p>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Razón <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="text"
                            value={formData.reason}
                            onChange={(e) => setFormData(prev => ({ ...prev, reason: e.target.value }))}
                            placeholder="Conteo físico, corrección, etc."
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Notas</label>
                        <textarea
                            value={formData.notes}
                            onChange={(e) => setFormData(prev => ({ ...prev, notes: e.target.value }))}
                            placeholder="Notas adicionales..."
                            rows={3}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                        />
                    </div>

                    <div className="flex justify-end gap-3 pt-4 border-t">
                        <Button type="button" variant="outline" onClick={onClose} disabled={loading}>
                            Cancelar
                        </Button>
                        <Button type="submit" variant="primary" disabled={loading}>
                            {loading ? 'Ajustando...' : 'Ajustar stock'}
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    );
}
