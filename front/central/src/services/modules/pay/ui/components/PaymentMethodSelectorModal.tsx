'use client';

import { useEffect, useState } from 'react';
import { Modal } from '@/shared/ui';
import { getPaymentGatewayTypesAction } from '../../infra/actions';
import { PaymentGatewayType } from '../../domain/types';
import { PaymentMethodGrid } from './PaymentMethodGrid';

const formatCurrency = (amount: string | number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP' }).format(Number(amount));

interface PaymentMethodSelectorModalProps {
    isOpen: boolean;
    onClose: () => void;
    amount: string;
    onSelectNequi: () => void;
    onSelectOther: (gatewayName: string) => void;
}

export function PaymentMethodSelectorModal({
    isOpen,
    onClose,
    amount,
    onSelectNequi,
    onSelectOther,
}: PaymentMethodSelectorModalProps) {
    const [gateways, setGateways] = useState<PaymentGatewayType[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (!isOpen) return;

        setLoading(true);
        getPaymentGatewayTypesAction()
            .then((data) => setGateways(data))
            .catch(() => setGateways([]))
            .finally(() => setLoading(false));
    }, [isOpen]);

    const handleSelect = (gateway: PaymentGatewayType) => {
        if (gateway.code === 'nequi') {
            onSelectNequi();
        } else {
            onSelectOther(gateway.name);
        }
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title="¿Cómo quieres pagar?"
            size="xl"
        >
            <div className="p-4 space-y-4">
                {amount && (
                    <div className="text-center">
                        <p className="text-sm text-gray-500">Monto a recargar</p>
                        <p className="text-2xl font-bold text-gray-900">{formatCurrency(amount)}</p>
                    </div>
                )}

                <p className="text-sm text-gray-600 text-center">
                    Selecciona tu método de pago preferido
                </p>

                <PaymentMethodGrid
                    gateways={gateways}
                    loading={loading}
                    onSelect={handleSelect}
                />
            </div>
        </Modal>
    );
}
