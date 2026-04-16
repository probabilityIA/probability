'use client';

import { useState, useCallback } from 'react';
import { Modal } from '@/shared/ui/modal';
import { Order } from '../../domain/types';
import { changeOrderStatusAction } from '../../infra/actions';
import {
    getValidTransitions,
    getStatusByCode,
    isTerminalStatus,
    STATUS_METADATA_FIELDS,
    type StatusStep,
    type MetadataField,
} from '../../domain/order-status-transitions';

interface ChangeStatusModalProps {
    isOpen: boolean;
    onClose: () => void;
    order: Order;
    onSuccess?: () => void;
}

export function ChangeStatusModal({ isOpen, onClose, order, onSuccess }: ChangeStatusModalProps) {
    const [step, setStep] = useState<'select' | 'confirm'>('select');
    const [selectedStatus, setSelectedStatus] = useState<StatusStep | null>(null);
    const [metadata, setMetadata] = useState<Record<string, string>>({});
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const currentStatusCode = order.order_status?.code || order.status || '';
    const currentStatus = getStatusByCode(currentStatusCode);
    const validTargets = getValidTransitions(currentStatusCode);
    const metadataFields = selectedStatus ? (STATUS_METADATA_FIELDS[selectedStatus.code] || []) : [];

    const hasRequiredMetadata = metadataFields
        .filter(f => f.required)
        .every(f => metadata[f.key]?.trim());

    const canProceed = selectedStatus && (metadataFields.length === 0 || hasRequiredMetadata);

    const handleClose = useCallback(() => {
        if (loading) return;
        setStep('select');
        setSelectedStatus(null);
        setMetadata({});
        setError(null);
        setLoading(false);
        onClose();
    }, [loading, onClose]);

    const handleSelectStatus = useCallback((status: StatusStep) => {
        setSelectedStatus(status);
        setMetadata({});
        setError(null);
    }, []);

    const handleMetadataChange = useCallback((key: string, value: string) => {
        setMetadata(prev => ({ ...prev, [key]: value }));
    }, []);

    const handleConfirm = useCallback(async () => {
        if (!selectedStatus) return;
        setLoading(true);
        setError(null);

        const metaPayload: Record<string, unknown> = {};
        for (const [key, value] of Object.entries(metadata)) {
            if (value.trim()) metaPayload[key] = value.trim();
        }

        const result = await changeOrderStatusAction(order.id, {
            status: selectedStatus.code,
            metadata: Object.keys(metaPayload).length > 0 ? metaPayload : undefined,
        });

        setLoading(false);

        if (result && 'success' in result && result.success === false) {
            setError(result.message || 'Error al cambiar estado');
            return;
        }

        // Exito
        onSuccess?.();
    }, [selectedStatus, metadata, order.id, onSuccess]);

    const isDestructive = selectedStatus?.code === 'cancelled' || selectedStatus?.code === 'refunded';

    return (
        <Modal isOpen={isOpen} onClose={handleClose} title={`Cambiar Estado — #${order.order_number}`} size="lg" showCloseButton={!loading}>
            <div className="space-y-5">
                {/* Estado actual */}
                <div className="flex items-center gap-3">
                    <span className="text-sm text-gray-500 dark:text-gray-400">Estado actual:</span>
                    {currentStatus ? (
                        <span
                            className="px-3 py-1 text-sm font-semibold rounded-full"
                            style={{ backgroundColor: currentStatus.color + '20', color: currentStatus.color }}
                        >
                            {currentStatus.name}
                        </span>
                    ) : (
                        <span className="px-3 py-1 text-sm font-medium rounded-full bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300">
                            {currentStatusCode}
                        </span>
                    )}
                </div>

                {/* Error */}
                {error && (
                    <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg p-3">
                        <p className="text-sm text-red-700 dark:text-red-300">{error}</p>
                    </div>
                )}

                {/* STEP 1: Select */}
                {step === 'select' && (
                    <>
                        {validTargets.length === 0 ? (
                            <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                                <p className="text-sm">Esta orden esta en un estado terminal y no se puede cambiar.</p>
                            </div>
                        ) : (
                            <>
                                <div>
                                    <p className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">Selecciona el nuevo estado:</p>
                                    <div className="grid gap-2">
                                        {validTargets.map((target) => {
                                            const isSelected = selectedStatus?.code === target.code;
                                            const isTargetDestructive = target.code === 'cancelled' || target.code === 'refunded';
                                            return (
                                                <button
                                                    key={target.code}
                                                    onClick={() => handleSelectStatus(target)}
                                                    className={`flex items-start gap-3 p-3 rounded-lg border-2 text-left transition-all ${
                                                        isSelected
                                                            ? 'border-purple-500 bg-purple-50 dark:bg-purple-900/20'
                                                            : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
                                                    } ${isTargetDestructive ? 'hover:border-red-300 dark:hover:border-red-700' : ''}`}
                                                >
                                                    <div
                                                        className="w-4 h-4 rounded-full mt-0.5 flex-shrink-0"
                                                        style={{ backgroundColor: target.color }}
                                                    />
                                                    <div className="flex-1 min-w-0">
                                                        <span className="text-sm font-medium text-gray-900 dark:text-white">{target.name}</span>
                                                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{target.description}</p>
                                                    </div>
                                                    {isSelected && (
                                                        <svg className="w-5 h-5 text-purple-500 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                                                            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z" clipRule="evenodd" />
                                                        </svg>
                                                    )}
                                                </button>
                                            );
                                        })}
                                    </div>
                                </div>

                                {/* Metadata fields */}
                                {selectedStatus && metadataFields.length > 0 && (
                                    <div className="space-y-3 pt-2 border-t border-gray-200 dark:border-gray-700">
                                        <p className="text-sm font-medium text-gray-700 dark:text-gray-300">Informacion adicional:</p>
                                        {metadataFields.map((field) => (
                                            <div key={field.key}>
                                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">
                                                    {field.label} {field.required && <span className="text-red-500">*</span>}
                                                </label>
                                                {field.type === 'textarea' ? (
                                                    <textarea
                                                        value={metadata[field.key] || ''}
                                                        onChange={(e) => handleMetadataChange(field.key, e.target.value)}
                                                        placeholder={field.placeholder}
                                                        rows={2}
                                                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                                                    />
                                                ) : (
                                                    <input
                                                        type="text"
                                                        value={metadata[field.key] || ''}
                                                        onChange={(e) => handleMetadataChange(field.key, e.target.value)}
                                                        placeholder={field.placeholder}
                                                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                                                    />
                                                )}
                                            </div>
                                        ))}
                                    </div>
                                )}

                                {/* Boton siguiente */}
                                <div className="flex justify-end pt-2">
                                    <button
                                        onClick={() => { setError(null); setStep('confirm'); }}
                                        disabled={!canProceed}
                                        className="px-5 py-2 text-sm font-semibold text-white bg-purple-600 rounded-lg hover:bg-purple-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
                                    >
                                        Siguiente
                                    </button>
                                </div>
                            </>
                        )}
                    </>
                )}

                {/* STEP 2: Confirm */}
                {step === 'confirm' && selectedStatus && (
                    <>
                        {/* Resumen */}
                        <div className="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-4 space-y-3">
                            <div className="flex items-center gap-2">
                                {currentStatus && (
                                    <span
                                        className="px-2 py-0.5 text-xs font-medium rounded-full"
                                        style={{ backgroundColor: currentStatus.color + '20', color: currentStatus.color }}
                                    >
                                        {currentStatus.name}
                                    </span>
                                )}
                                <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                                </svg>
                                <span
                                    className="px-2 py-0.5 text-xs font-medium rounded-full"
                                    style={{ backgroundColor: selectedStatus.color + '20', color: selectedStatus.color }}
                                >
                                    {selectedStatus.name}
                                </span>
                            </div>

                            {/* Metadata resumen */}
                            {metadataFields.length > 0 && Object.keys(metadata).some(k => metadata[k]?.trim()) && (
                                <div className="text-xs text-gray-600 dark:text-gray-400 space-y-1 pt-2 border-t border-gray-200 dark:border-gray-600">
                                    {metadataFields.map(f => metadata[f.key]?.trim() ? (
                                        <p key={f.key}><span className="font-medium">{f.label}:</span> {metadata[f.key]}</p>
                                    ) : null)}
                                </div>
                            )}
                        </div>

                        {/* Warning para destructivos */}
                        {isDestructive && (
                            <div className="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-700 rounded-lg p-3">
                                <p className="text-sm text-amber-700 dark:text-amber-300">
                                    <strong>Atencion:</strong> {selectedStatus.code === 'cancelled' ? 'La cancelacion' : 'El reembolso'} es un estado terminal. Una vez aplicado, no se podra revertir.
                                </p>
                            </div>
                        )}

                        {/* Botones */}
                        <div className="flex justify-between pt-2">
                            <button
                                onClick={() => { setStep('select'); setError(null); }}
                                disabled={loading}
                                className="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 disabled:opacity-50 transition-colors"
                            >
                                Volver
                            </button>
                            <button
                                onClick={handleConfirm}
                                disabled={loading}
                                className={`px-5 py-2 text-sm font-semibold text-white rounded-lg transition-colors disabled:opacity-50 flex items-center gap-2 ${
                                    isDestructive
                                        ? 'bg-red-600 hover:bg-red-700'
                                        : 'bg-purple-600 hover:bg-purple-700'
                                }`}
                            >
                                {loading && (
                                    <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                                    </svg>
                                )}
                                {loading ? 'Procesando...' : 'Confirmar Cambio'}
                            </button>
                        </div>
                    </>
                )}
            </div>
        </Modal>
    );
}
