'use client';

import { useState, useEffect } from 'react';
import { TransferStockDTO } from '../../domain/types';
import { transferStockAction } from '../../infra/actions';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { Warehouse } from '@/services/modules/warehouses/domain/types';
import { Button, Alert, Input } from '@/shared/ui';

interface TransferStockModalProps {
    fromWarehouseId: number;
    businessId?: number;
    onSuccess: () => void;
    onClose: () => void;
}

export default function TransferStockModal({ fromWarehouseId, businessId, onSuccess, onClose }: TransferStockModalProps) {
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [formData, setFormData] = useState({
        product_id: '',
        to_warehouse_id: 0,
        quantity: 1,
        reason: '',
        notes: '',
    });

    const [loading, setLoading] = useState(false);
    const [loadingWarehouses, setLoadingWarehouses] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    useEffect(() => {
        const loadWarehouses = async () => {
            try {
                const response = await getWarehousesAction({
                    page: 1,
                    page_size: 100,
                    is_active: true,
                    business_id: businessId,
                });
                setWarehouses((response.data || []).filter(w => w.id !== fromWarehouseId));
            } catch {
                setError('Error al cargar bodegas');
            } finally {
                setLoadingWarehouses(false);
            }
        };
        loadWarehouses();
    }, [businessId, fromWarehouseId]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            const dto: TransferStockDTO = {
                product_id: formData.product_id.trim(),
                from_warehouse_id: fromWarehouseId,
                to_warehouse_id: formData.to_warehouse_id,
                quantity: formData.quantity,
                reason: formData.reason.trim() || undefined,
                notes: formData.notes.trim() || undefined,
            };
            await transferStockAction(dto, businessId);
            setSuccess('Transferencia realizada exitosamente');
            setTimeout(() => onSuccess(), 800);
        } catch (err: any) {
            setError(err.message || 'Error al transferir stock');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
            <div className="bg-white rounded-xl shadow-xl w-full max-w-md max-h-[90vh] overflow-y-auto">
                <div className="flex items-center justify-between px-6 py-4 border-b">
                    <h2 className="text-lg font-semibold text-gray-900">Transferir stock</h2>
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
                            Bodega destino <span className="text-red-500">*</span>
                        </label>
                        {loadingWarehouses ? (
                            <p className="text-sm text-gray-500">Cargando bodegas...</p>
                        ) : (
                            <select
                                value={formData.to_warehouse_id || ''}
                                onChange={(e) => setFormData(prev => ({ ...prev, to_warehouse_id: Number(e.target.value) }))}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                required
                            >
                                <option value="">— Selecciona bodega destino —</option>
                                {warehouses.map(w => (
                                    <option key={w.id} value={w.id}>{w.name} ({w.code})</option>
                                ))}
                            </select>
                        )}
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Cantidad <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="number"
                            value={formData.quantity.toString()}
                            onChange={(e) => setFormData(prev => ({ ...prev, quantity: parseInt(e.target.value) || 0 }))}
                            min="1"
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Razón</label>
                        <Input
                            type="text"
                            value={formData.reason}
                            onChange={(e) => setFormData(prev => ({ ...prev, reason: e.target.value }))}
                            placeholder="Reabastecimiento, redistribución, etc."
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
                            {loading ? 'Transfiriendo...' : 'Transferir'}
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    );
}
