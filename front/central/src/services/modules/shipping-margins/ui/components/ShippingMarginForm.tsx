'use client';

import { useState } from 'react';
import { ShippingMargin, CARRIER_OPTIONS, CreateShippingMarginDTO, UpdateShippingMarginDTO } from '../../domain/types';
import { createShippingMarginAction, updateShippingMarginAction } from '../../infra/actions';
import { Button, Alert, Input } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

interface Props {
    margin?: ShippingMargin;
    onSuccess: () => void;
    onCancel: () => void;
    businessId?: number;
}

export default function ShippingMarginForm({ margin, onSuccess, onCancel, businessId }: Props) {
    const isEdit = !!margin;
    const [formData, setFormData] = useState({
        carrier_code: margin?.carrier_code || CARRIER_OPTIONS[0].code,
        carrier_name: margin?.carrier_name || CARRIER_OPTIONS[0].name,
        margin_amount: margin?.margin_amount?.toString() || '0',
        insurance_margin: margin?.insurance_margin?.toString() || '0',
        is_active: margin?.is_active ?? true,
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleCarrierChange = (code: string) => {
        const opt = CARRIER_OPTIONS.find((c) => c.code === code);
        setFormData((p) => ({ ...p, carrier_code: code, carrier_name: opt?.name || code }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        try {
            const margin_amount = parseFloat(formData.margin_amount) || 0;
            const insurance_margin = parseFloat(formData.insurance_margin) || 0;

            if (isEdit && margin) {
                const payload: UpdateShippingMarginDTO = {
                    carrier_name: formData.carrier_name,
                    margin_amount,
                    insurance_margin,
                    is_active: formData.is_active,
                };
                await updateShippingMarginAction(margin.id, payload, businessId);
            } else {
                const payload: CreateShippingMarginDTO = {
                    carrier_code: formData.carrier_code,
                    carrier_name: formData.carrier_name,
                    margin_amount,
                    insurance_margin,
                    is_active: formData.is_active,
                };
                await createShippingMarginAction(payload, businessId);
            }
            onSuccess();
        } catch (err: any) {
            setError(getActionError(err, 'Error al guardar el margen'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Transportadora <span className="text-red-500">*</span>
                    </label>
                    {isEdit ? (
                        <Input type="text" value={formData.carrier_name} disabled />
                    ) : (
                        <select
                            value={formData.carrier_code}
                            onChange={(e) => handleCarrierChange(e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                        >
                            {CARRIER_OPTIONS.map((c) => (
                                <option key={c.code} value={c.code}>
                                    {c.name}
                                </option>
                            ))}
                        </select>
                    )}
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Estado
                    </label>
                    <select
                        value={formData.is_active ? 'true' : 'false'}
                        onChange={(e) => setFormData((p) => ({ ...p, is_active: e.target.value === 'true' }))}
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                    >
                        <option value="true">Activo</option>
                        <option value="false">Inactivo</option>
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Margen sobre flete (COP) <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="number"
                        min="0"
                        step="100"
                        value={formData.margin_amount}
                        onChange={(e) => setFormData((p) => ({ ...p, margin_amount: e.target.value }))}
                        placeholder="2000"
                        required
                    />
                    <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                        Pesos colombianos sumados al precio del proveedor por cada guia
                    </p>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                        Margen sobre seguro (COP)
                    </label>
                    <Input
                        type="number"
                        min="0"
                        step="100"
                        value={formData.insurance_margin}
                        onChange={(e) => setFormData((p) => ({ ...p, insurance_margin: e.target.value }))}
                        placeholder="0"
                    />
                    <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                        Cargo adicional sobre el seguro obligatorio del carrier (0 = sin cargo extra)
                    </p>
                </div>
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                <Button type="button" variant="outline" onClick={onCancel} disabled={loading}>
                    Cancelar
                </Button>
                <Button type="submit" variant="primary" disabled={loading}>
                    {loading ? 'Guardando...' : isEdit ? 'Actualizar' : 'Crear margen'}
                </Button>
            </div>
        </form>
    );
}
