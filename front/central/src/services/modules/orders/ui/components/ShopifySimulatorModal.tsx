'use client';

import { useState } from 'react';
import { Modal } from '@/shared/ui';
import { simulateShopifyAction } from '../../infra/actions/testing-actions';
import { SimulateShopifyResult } from '../../domain/types';

interface ShopifySimulatorModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSuccess: () => void;
}

const TOPICS = [
    { value: 'orders/create', label: 'orders/create (nueva orden)' },
    { value: 'orders/paid', label: 'orders/paid (pagada)' },
    { value: 'orders/updated', label: 'orders/updated (actualizada)' },
    { value: 'orders/cancelled', label: 'orders/cancelled (cancelada)' },
    { value: 'orders/fulfilled', label: 'orders/fulfilled (cumplida)' },
    { value: 'orders/partially_fulfilled', label: 'orders/partially_fulfilled (parcial)' },
];

export default function ShopifySimulatorModal({ isOpen, onClose, onSuccess }: ShopifySimulatorModalProps) {
    const [topic, setTopic] = useState('orders/create');
    const [count, setCount] = useState(1);
    const [loading, setLoading] = useState(false);
    const [result, setResult] = useState<SimulateShopifyResult | null>(null);
    const [error, setError] = useState<string | null>(null);

    const handleSimulate = async () => {
        setLoading(true);
        setResult(null);
        setError(null);

        const response = await simulateShopifyAction(topic, count);

        if (response.success && response.data) {
            setResult(response.data);
            if (response.data.sent > 0) {
                onSuccess();
            }
        } else {
            setError(response.error || 'Error desconocido');
        }

        setLoading(false);
    };

    const handleClose = () => {
        setResult(null);
        setError(null);
        setTopic('orders/create');
        setCount(1);
        onClose();
    };

    return (
        <Modal isOpen={isOpen} onClose={handleClose} title="Simular Webhook Shopify" size="md">
            <div className="space-y-4 p-4">
                {/* Topic selector */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Tipo de webhook
                    </label>
                    <select
                        value={topic}
                        onChange={(e) => setTopic(e.target.value)}
                        disabled={loading}
                        className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500 focus:border-transparent"
                    >
                        {TOPICS.map((t) => (
                            <option key={t.value} value={t.value}>
                                {t.label}
                            </option>
                        ))}
                    </select>
                </div>

                {/* Count input */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Cantidad de webhooks
                    </label>
                    <input
                        type="number"
                        min={1}
                        max={20}
                        value={count}
                        onChange={(e) => {
                            const val = parseInt(e.target.value) || 1;
                            setCount(Math.min(20, Math.max(1, val)));
                        }}
                        disabled={loading}
                        className="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500 focus:border-transparent"
                    />
                    <p className="text-xs text-gray-500 mt-1">Entre 1 y 20 webhooks</p>
                </div>

                {/* Submit button */}
                <button
                    onClick={handleSimulate}
                    disabled={loading}
                    style={{ background: '#7c3aed' }}
                    className="w-full px-4 py-2 text-sm font-semibold text-white rounded-lg hover:shadow-lg hover:scale-105 transition-all disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100"
                >
                    {loading ? 'Simulando...' : 'Simular Webhook'}
                </button>

                {/* Result */}
                {result && (
                    <div className={`rounded-lg p-3 text-sm ${result.failed > 0 ? 'bg-yellow-50 border border-yellow-200' : 'bg-green-50 border border-green-200'}`}>
                        <p className="font-medium">
                            {result.sent}/{result.total} webhooks enviados
                        </p>
                        {result.errors && result.errors.length > 0 && (
                            <ul className="mt-2 space-y-1 text-red-600">
                                {result.errors.map((err, i) => (
                                    <li key={i} className="text-xs">{err}</li>
                                ))}
                            </ul>
                        )}
                    </div>
                )}

                {/* Error */}
                {error && (
                    <div className="rounded-lg p-3 text-sm bg-red-50 border border-red-200 text-red-700">
                        {error}
                    </div>
                )}
            </div>
        </Modal>
    );
}
