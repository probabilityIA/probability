'use client';

import { useState, useEffect } from 'react';
import { CustomerDetail } from '../../domain/types';
import { getCustomerByIdAction } from '../../infra/actions';
import { Spinner, Alert } from '@/shared/ui';

interface CustomerDetailProps {
    customerId: number;
    businessId?: number;
}

function StatCard({ label, value }: { label: string; value: string }) {
    return (
        <div className="bg-gray-50 rounded-lg p-4 text-center">
            <p className="text-xs text-gray-500 uppercase tracking-wide mb-1">{label}</p>
            <p className="text-lg font-semibold text-gray-900">{value}</p>
        </div>
    );
}

function Field({ label, value }: { label: string; value?: string | null }) {
    return (
        <div>
            <p className="text-xs text-gray-500 uppercase tracking-wide mb-1">{label}</p>
            <p className="text-sm text-gray-900">{value || <span className="text-gray-300">—</span>}</p>
        </div>
    );
}

export default function CustomerDetailView({ customerId, businessId }: CustomerDetailProps) {
    const [customer, setCustomer] = useState<CustomerDetail | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetch = async () => {
            setLoading(true);
            setError(null);
            try {
                const data = await getCustomerByIdAction(customerId, businessId);
                setCustomer(data);
            } catch (err: any) {
                setError(err.message || 'Error al cargar el cliente');
            } finally {
                setLoading(false);
            }
        };
        fetch();
    }, [customerId, businessId]);

    if (loading) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    if (error) {
        return (
            <Alert type="error" onClose={() => setError(null)}>
                {error}
            </Alert>
        );
    }

    if (!customer) return null;

    const totalSpentFormatted = new Intl.NumberFormat('es-CO', {
        style: 'currency',
        currency: 'COP',
        minimumFractionDigits: 0,
    }).format(customer.total_spent);

    const lastOrderFormatted = customer.last_order_at
        ? new Date(customer.last_order_at).toLocaleDateString('es-CO', {
              day: '2-digit',
              month: 'short',
              year: 'numeric',
          })
        : '—';

    return (
        <div className="space-y-6">
            {/* Stats de órdenes */}
            <div className="grid grid-cols-3 gap-4">
                <StatCard label="Órdenes" value={String(customer.order_count)} />
                <StatCard label="Total gastado" value={totalSpentFormatted} />
                <StatCard label="Última orden" value={lastOrderFormatted} />
            </div>

            {/* Datos del cliente */}
            <div className="bg-white border border-gray-200 rounded-lg p-5 space-y-4">
                <h3 className="text-sm font-semibold text-gray-700 uppercase tracking-wide border-b pb-2">
                    Información de contacto
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <Field label="Nombre" value={customer.name} />
                    <Field label="Email" value={customer.email} />
                    <Field label="Teléfono" value={customer.phone} />
                    <Field label="Documento" value={customer.dni} />
                </div>
            </div>

            {/* Metadatos */}
            <div className="text-xs text-gray-400 flex gap-4">
                <span>
                    Creado: {new Date(customer.created_at).toLocaleDateString('es-CO', {
                        day: '2-digit', month: 'short', year: 'numeric',
                    })}
                </span>
                <span>
                    Actualizado: {new Date(customer.updated_at).toLocaleDateString('es-CO', {
                        day: '2-digit', month: 'short', year: 'numeric',
                    })}
                </span>
            </div>
        </div>
    );
}
