import { Modal } from '@/shared/ui';
import { useEffect, useState } from 'react';
import { getOrderRawAction } from '../../infra/actions';

interface RawOrderModalProps {
    orderId: string;
    isOpen: boolean;
    onClose: () => void;
}

export default function RawOrderModal({ orderId, isOpen, onClose }: RawOrderModalProps) {
    const [data, setData] = useState<any>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (isOpen && orderId) {
            fetchRawData();
        } else {
            setData(null);
            setError(null);
        }
    }, [isOpen, orderId]);

    const fetchRawData = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getOrderRawAction(orderId);
            if (response.success) {
                setData(response.data);
            } else {
                setError(response.message || 'Error al cargar los datos crudos');
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar los datos crudos');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title="Datos Originales de la Orden"
            size="xl"
        >
            <div className="p-4">
                {loading && <div className="text-center py-4">Cargando datos...</div>}
                {error && (
                    <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative mb-4">
                        {error}
                    </div>
                )}
                {data && (
                    <div className="bg-gray-900 rounded-lg p-4 overflow-auto max-h-[60vh]">
                        <pre className="text-green-400 font-mono text-sm whitespace-pre-wrap">
                            {JSON.stringify(data, null, 2)}
                        </pre>
                    </div>
                )}
                {!loading && !error && !data && (
                    <div className="text-center py-4 text-gray-500">
                        No hay datos disponibles
                    </div>
                )}
            </div>
        </Modal>
    );
}
