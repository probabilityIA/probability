import { useState } from 'react';
import { Order } from '@/services/modules/orders/domain/types';

interface ShipmentGuideModalProps {
    isOpen: boolean;
    onClose: () => void;
    order: Order;
    recommendedCarrier?: string;
    onGuideGenerated?: (trackingNumber: string) => void;
}

interface EnvioClickRate {
    idRate: number;
    carrier: string;
    product: string;
    flete: number;
    deliveryDays: number;
    quotationType: string;
}

export default function ShipmentGuideModal({ isOpen, onClose, order, recommendedCarrier, onGuideGenerated }: ShipmentGuideModalProps) {
    const [step, setStep] = useState<1 | 2 | 3>(1);
    const [origin, setOrigin] = useState({ daneCode: '', address: '' }); // Simplified origin for now
    const [loading, setLoading] = useState(false);
    const [rates, setRates] = useState<EnvioClickRate[]>([]);
    const [selectedRate, setSelectedRate] = useState<EnvioClickRate | null>(null);
    const [error, setError] = useState<string | null>(null);

    // Hardcoded DANE code map for demo purposes or manual entry
    // ideally this would be a select input with autocomplete
    const handleOriginChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setOrigin({ ...origin, [e.target.name]: e.target.value });
    };

    const handleQuote = async () => {
        setLoading(true);
        setError(null);
        try {
            // Use environment variable or default to port 3050 as found in backend config
            const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3050/api/v1';
            const res = await fetch(`${apiUrl}/shipments/quote`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    origin: { daneCode: origin.daneCode || '11001000', address: origin.address }, // Default Bogota for testing
                    destination: { daneCode: '05001000', address: order.shipping_street || 'Calle 123' }, // Default Medellin/Order address
                    packages: [{ weight: 1, height: 10, width: 10, length: 10, contentValue: order.total_amount }], // Mock package
                    description: 'E-commerce Order'
                })
            });
            const data = await res.json();

            if (res.ok && data.data && data.data.rates) {
                setRates(data.data.rates);
                setStep(2);
            } else {
                setError(data.error || 'Error al cotizar');
            }
        } catch (err: any) {
            setError(err.message || 'Error de conexión');
        } finally {
            setLoading(false);
        }
    };

    const handleGenerate = async () => {
        if (!selectedRate) return;
        setLoading(true);
        setError(null);
        try {
            const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3050/api/v1';
            const res = await fetch(`${apiUrl}/shipments/generate`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    idRate: selectedRate.idRate,
                    orderId: order.id
                    // other needed fields
                })
            });
            const data = await res.json();
            if (res.ok) {
                setStep(3);
                if (onGuideGenerated) onGuideGenerated(data.data?.trackingNumber);
            } else {
                setError(data.error || 'Error al generar guía');
            }
        } catch (err: any) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
            <div className="bg-white rounded-xl shadow-xl w-full max-w-2xl overflow-hidden max-h-[90vh] flex flex-col">
                {/* Header */}
                <div className="p-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
                    <h3 className="text-lg font-bold text-gray-800">Generar Guía de Envío</h3>
                    <button onClick={onClose} className="text-gray-400 hover:text-gray-600">✕</button>
                </div>

                {/* Body */}
                <div className="p-6 overflow-y-auto flex-1">
                    {error && (
                        <div className="bg-red-50 text-red-600 p-3 rounded-lg mb-4 text-sm">
                            ⚠️ {error}
                        </div>
                    )}

                    {step === 1 && (
                        <div className="space-y-4">
                            <h4 className="text-blue-900 font-semibold mb-2">Paso 1: Confirmar Origen</h4>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Dirección de Recogida</label>
                                <input
                                    type="text"
                                    name="address"
                                    value={origin.address}
                                    onChange={handleOriginChange}
                                    placeholder="Ej: Calle 100 # 15-20"
                                    className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 outline-none"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Código DANE Ciudad</label>
                                <input
                                    type="text"
                                    name="daneCode"
                                    value={origin.daneCode}
                                    onChange={handleOriginChange}
                                    placeholder="Ej: 11001000 (Bogotá)"
                                    className="w-full border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 outline-none"
                                />
                                <p className="text-xs text-gray-400 mt-1">* Temporalmente manual. 11001000 = Bogotá.</p>
                            </div>

                            <div className="pt-4 flex justify-end">
                                <button
                                    onClick={handleQuote}
                                    disabled={loading || !origin.address || !origin.daneCode}
                                    className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg transition-colors disabled:opacity-50 flex items-center gap-2"
                                >
                                    {loading ? 'Cotizando...' : 'Cotizar Envíos'}
                                </button>
                            </div>
                        </div>
                    )}

                    {step === 2 && (
                        <div className="space-y-4">
                            <h4 className="text-blue-900 font-semibold mb-2">Paso 2: Seleccionar Tarifa</h4>
                            <p className="text-sm text-gray-600 mb-4">Seleccione la transportadora para generar la guía.</p>

                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                                {rates.map((rate) => {
                                    const isRecommended = rate.carrier === recommendedCarrier;
                                    const isSelected = selectedRate?.idRate === rate.idRate;

                                    return (
                                        <div
                                            key={rate.idRate}
                                            onClick={() => setSelectedRate(rate)}
                                            className={`border rounded-lg p-4 cursor-pointer transition-all relative ${isSelected
                                                    ? 'border-blue-500 bg-blue-50 ring-2 ring-blue-200'
                                                    : isRecommended
                                                        ? 'border-green-300 bg-green-50/50 hover:bg-green-50'
                                                        : 'border-gray-200 hover:bg-gray-50'
                                                }`}
                                        >
                                            {isRecommended && (
                                                <div className="absolute top-0 right-0 bg-green-500 text-white text-[10px] px-2 py-0.5 rounded-bl-lg font-bold">
                                                    RECOMENDADO IA
                                                </div>
                                            )}
                                            <div className="flex justify-between items-start mb-2">
                                                <span className="font-bold text-gray-800">{rate.carrier}</span>
                                                <span className="text-lg font-bold text-blue-700">${rate.flete.toLocaleString()}</span>
                                            </div>
                                            <p className="text-sm text-gray-500">{rate.product} - {rate.deliveryDays} días</p>
                                        </div>
                                    );
                                })}
                            </div>

                            <div className="pt-6 flex justify-between items-center border-t mt-4">
                                <button onClick={() => setStep(1)} className="text-gray-500 hover:text-gray-700 text-sm underline">Volver</button>
                                <button
                                    onClick={handleGenerate}
                                    disabled={loading || !selectedRate}
                                    className="bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-6 rounded-lg transition-colors disabled:opacity-50"
                                >
                                    {loading ? 'Generando...' : 'Generar Guía Ahora'}
                                </button>
                            </div>
                        </div>
                    )}

                    {step === 3 && (
                        <div className="text-center py-8">
                            <div className="w-16 h-16 bg-green-100 text-green-600 rounded-full flex items-center justify-center mx-auto mb-4 text-3xl">
                                ✓
                            </div>
                            <h4 className="text-xl font-bold text-gray-900 mb-2">¡Guía Generada Exitosamente!</h4>
                            <p className="text-gray-600 mb-6">El número de rastreo ha sido asignado a la orden.</p>
                            <button
                                onClick={onClose}
                                className="bg-gray-900 text-white py-2 px-8 rounded-lg hover:bg-gray-800 transition-colors"
                            >
                                Cerrar
                            </button>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
