'use client';

import { useState, useEffect } from 'react';
import { Button } from '@/shared/ui';
import { Order } from '@/services/modules/orders/domain/types';
import { getOrdersAction } from '@/services/modules/orders/infra/actions';
import { ShipmentApiRepository } from '@/services/modules/shipments/infra/repository/api-repository';
import { EnvioClickQuoteRequest, EnvioClickRate } from '@/services/modules/shipments/domain/types';
import { getWalletBalanceAction } from '@/services/modules/wallet/infra/actions';
import danes from "@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json";

interface MassGuideGenerationModalProps {
    isOpen: boolean;
    onClose: () => void;
    onComplete?: (count: number) => void;
}

interface OrderWithQuote extends Order {
    quote?: EnvioClickRate;
    quoteError?: string;
}

export default function MassGuideGenerationModal({ isOpen, onClose, onComplete }: MassGuideGenerationModalProps) {
    const [orders, setOrders] = useState<OrderWithQuote[]>([]);
    const [selectedOrderIds, setSelectedOrderIds] = useState<Set<string>>(new Set());
    const [loading, setLoading] = useState(false);
    const [quotingProgress, setQuotingProgress] = useState(0);
    const [generatingProgress, setGeneratingProgress] = useState(0);
    const [step, setStep] = useState<'select' | 'quote' | 'confirm' | 'generate' | 'complete'>('select');
    const [error, setError] = useState<string | null>(null);
    const [totalCost, setTotalCost] = useState(0);
    const [walletBalance, setWalletBalance] = useState<number | null>(null);
    const [generatedCount, setGeneratedCount] = useState(0);
    const [failedCount, setFailedCount] = useState(0);
    const [generationErrors, setGenerationErrors] = useState<string[]>([]);
    const [selectedOrderForDetails, setSelectedOrderForDetails] = useState<OrderWithQuote | null>(null);
    const [isEditingOrder, setIsEditingOrder] = useState(false);

    const repo = new ShipmentApiRepository();

    useEffect(() => {
        if (isOpen && step === 'select') {
            loadOrders();
            loadWalletBalance();
        }
    }, [isOpen, step]);

    const loadOrders = async () => {
        setLoading(true);
        try {
            const response = await getOrdersAction({ page: 1, page_size: 100 });
            if (response.success && response.data) {
                // Filter orders without tracking numbers
                const ordersWithoutGuides = response.data.filter(order => !order.tracking_number);
                setOrders(ordersWithoutGuides);
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar órdenes');
        } finally {
            setLoading(false);
        }
    };

    const loadWalletBalance = async () => {
        try {
            const response = await getWalletBalanceAction();
            if (typeof response === 'number') {
                setWalletBalance(response);
            } else if (response && 'data' in response && response.data) {
                setWalletBalance(response.data.Balance || 0);
            }
        } catch (err) {
            console.error('Error loading wallet balance:', err);
        }
    };

    const toggleOrderSelection = (orderId: string) => {
        const newSelection = new Set(selectedOrderIds);
        if (newSelection.has(orderId)) {
            newSelection.delete(orderId);
        } else {
            newSelection.add(orderId);
        }
        setSelectedOrderIds(newSelection);
    };

    const selectAll = () => {
        setSelectedOrderIds(new Set(orders.map(o => o.id)));
    };

    const deselectAll = () => {
        setSelectedOrderIds(new Set());
    };

    const handleQuoteAll = async () => {
        setStep('quote');
        setLoading(true);
        setError(null);
        setQuotingProgress(0);

        const selectedOrders = orders.filter(o => selectedOrderIds.has(o.id));
        const quotedOrders: OrderWithQuote[] = [];
        let total = 0;

        for (let i = 0; i < selectedOrders.length; i++) {
            const order = selectedOrders[i];
            try {
                const quotePayload: EnvioClickQuoteRequest = {
                    packages: [{
                        weight: order.weight || 1,
                        height: order.height || 10,
                        width: order.width || 10,
                        length: order.length || 10,
                    }],
                    description: `Orden ${order.order_number}`,
                    contentValue: order.total_amount || 10000,
                    includeGuideCost: false,
                    codPaymentMethod: 'cash',
                    origin: {
                        daneCode: '11001000', // Default Bogotá
                        address: 'Calle 1 # 1-1',
                    },
                    destination: {
                        daneCode: (() => {
                            const city = (order.shipping_city || "").toUpperCase();
                            const foundDane = Object.entries(danes).find(([_, data]: [string, any]) =>
                                (data as any).ciudad.includes(city)
                            );
                            return foundDane ? foundDane[0] : '11001001'; // Fallback to Bogotá if not found
                        })(),
                        address: order.shipping_street || 'Dirección no especificada',
                    },
                };

                const response = await repo.quoteShipment(quotePayload);
                if (response.data?.rates && response.data.rates.length > 0) {
                    const cheapestRate = response.data.rates.reduce((prev, curr) =>
                        curr.flete < prev.flete ? curr : prev
                    );
                    quotedOrders.push({ ...order, quote: cheapestRate });
                    total += cheapestRate.flete + (cheapestRate.minimumInsurance ?? 0);
                } else {
                    quotedOrders.push({ ...order, quoteError: 'No hay tarifas disponibles' });
                }
            } catch (err: any) {
                quotedOrders.push({ ...order, quoteError: err.message || 'Error al cotizar' });
            }
            setQuotingProgress(((i + 1) / selectedOrders.length) * 100);
        }

        setOrders(quotedOrders);
        setTotalCost(total);
        setLoading(false);
        setStep('confirm');
    };

    const handleGenerateAll = async () => {
        setStep('generate');
        setLoading(true);
        setError(null);
        setGeneratingProgress(0);
        setGeneratedCount(0);
        setFailedCount(0);
        setGenerationErrors([]);

        const ordersToGenerate = orders.filter(o => selectedOrderIds.has(o.id) && o.quote);
        const errors: string[] = [];

        for (let i = 0; i < ordersToGenerate.length; i++) {
            const order = ordersToGenerate[i];
            try {
                const generatePayload: EnvioClickQuoteRequest = {
                    idRate: order.quote!.idRate,
                    myShipmentReference: order.order_number,
                    external_order_id: order.order_number,
                    order_uuid: order.id,
                    requestPickup: false,
                    pickupDate: new Date().toISOString().split('T')[0],
                    insurance: false,
                    description: `Orden ${order.order_number}`,
                    contentValue: order.total_amount || 10000,
                    includeGuideCost: false,
                    codPaymentMethod: 'cash',
                    packages: [{
                        weight: order.weight || 1,
                        height: order.height || 10,
                        width: order.width || 10,
                        length: order.length || 10,
                    }],
                    origin: {
                        daneCode: '11001000',
                        address: 'Calle 1 # 1-1',
                        company: 'Mi Empresa',
                        firstName: 'Admin',
                        lastName: 'User',
                        email: 'admin@example.com',
                        phone: '3001234567',
                        suburb: 'Centro',
                        crossStreet: 'Calle 2',
                        reference: 'Edificio principal',
                    },
                    destination: {
                        daneCode: '11001000',
                        address: order.shipping_street || 'Dirección no especificada',
                        company: order.customer_name || 'Cliente',
                        firstName: order.customer_name?.split(' ')[0] || 'Cliente',
                        lastName: order.customer_name?.split(' ')[1] || 'Apellido',
                        email: order.customer_email || 'cliente@example.com',
                        phone: order.customer_phone || '3009876543',
                        suburb: 'Barrio',
                        crossStreet: 'Calle principal',
                        reference: 'Casa',
                    },
                };

                await repo.generateGuide(generatePayload);
                setGeneratedCount(prev => prev + 1);
            } catch (err: any) {
                setFailedCount(prev => prev + 1);
                errors.push(`Orden ${order.order_number}: ${err.message || 'Error desconocido'}`);
            }
            setGeneratingProgress(((i + 1) / ordersToGenerate.length) * 100);
        }

        setGenerationErrors(errors);
        setLoading(false);
        setStep('complete');
    };

    const handleClose = () => {
        setStep('select');
        setSelectedOrderIds(new Set());
        setOrders([]);
        setError(null);
        setTotalCost(0);
        setGeneratedCount(0);
        setFailedCount(0);
        setGenerationErrors([]);
        onClose();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto p-6">
                <div className="flex justify-between items-center mb-6">
                    <h2 className="text-2xl font-bold text-gray-800">Generación Masiva de Guías</h2>
                    <button onClick={handleClose} className="text-gray-500 hover:text-gray-700 text-2xl">
                        ×
                    </button>
                </div>

                {/* Step 1: Select Orders */}
                {step === 'select' && (
                    <div className="space-y-4">
                        <div className="flex justify-between items-center">
                            <p className="text-sm text-gray-600">
                                Selecciona las órdenes para generar guías ({selectedOrderIds.size} seleccionadas)
                            </p>
                            <div className="space-x-2">
                                <Button variant="outline" size="sm" onClick={selectAll}>
                                    Seleccionar todas
                                </Button>
                                <Button variant="outline" size="sm" onClick={deselectAll}>
                                    Deseleccionar todas
                                </Button>
                            </div>
                        </div>

                        {loading ? (
                            <div className="text-center py-8">Cargando órdenes...</div>
                        ) : orders.length === 0 ? (
                            <div className="text-center py-8 text-gray-500">
                                No hay órdenes sin guía de envío
                            </div>
                        ) : (
                            <div className="border rounded-lg max-h-96 overflow-y-auto">
                                {orders.map(order => (
                                    <div
                                        key={order.id}
                                        onClick={() => toggleOrderSelection(order.id)}
                                        className="flex items-center p-3 border-b hover:bg-gray-50 cursor-pointer"
                                    >
                                        <input
                                            type="checkbox"
                                            checked={selectedOrderIds.has(order.id)}
                                            onChange={() => { }}
                                            className="mr-3"
                                        />
                                        <div className="flex-1">
                                            <div className="font-semibold">{order.order_number}</div>
                                            <div className="text-sm text-gray-600">
                                                {order.customer_name} - {order.shipping_city}
                                            </div>
                                        </div>
                                        <div className="text-right">
                                            <div className="font-semibold">${order.total_amount?.toLocaleString()}</div>
                                            <div className="text-xs text-gray-500">
                                                {order.weight}kg
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}

                        <div className="flex justify-end space-x-3">
                            <Button variant="outline" onClick={handleClose}>
                                Cancelar
                            </Button>
                            <Button
                                onClick={handleQuoteAll}
                                disabled={selectedOrderIds.size === 0 || loading}
                            >
                                Cotizar Seleccionadas ({selectedOrderIds.size})
                            </Button>
                        </div>
                    </div>
                )}

                {/* Step 2: Quoting Progress */}
                {step === 'quote' && (
                    <div className="space-y-4">
                        <p className="text-center text-gray-600">Cotizando envíos...</p>
                        <div className="w-full bg-gray-200 rounded-full h-4">
                            <div
                                className="bg-orange-500 h-4 rounded-full transition-all duration-300"
                                style={{ width: `${quotingProgress}%` }}
                            />
                        </div>
                        <p className="text-center text-sm text-gray-500">
                            {Math.round(quotingProgress)}% completado
                        </p>
                    </div>
                )}

                {/* Step 3: Confirm */}
                {step === 'confirm' && (
                    <div className="space-y-4">
                        {/* Summary Header */}
                        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                            <h3 className="font-semibold text-blue-800 mb-2">Resumen de Cotización</h3>
                            <div className="flex justify-between items-center text-sm text-blue-700">
                                <div>
                                    <p>Órdenes cotizadas: {orders.filter(o => o.quote).length}</p>
                                    <p>Órdenes con error: {orders.filter(o => o.quoteError).length}</p>
                                </div>
                                <div className="text-right">
                                    <p className="text-lg font-bold">Total: ${totalCost.toLocaleString()} COP</p>
                                    {walletBalance !== null && (
                                        <p className={walletBalance >= totalCost ? 'text-green-700' : 'text-red-700'}>
                                            Saldo: ${walletBalance.toLocaleString()} COP
                                        </p>
                                    )}
                                </div>
                            </div>
                        </div>

                        {/* Detailed Order List */}
                        <div className="border rounded-lg max-h-96 overflow-y-auto">
                            <table className="w-full text-sm">
                                <thead className="bg-gray-50 sticky top-0">
                                    <tr className="border-b">
                                        <th className="text-left p-3 font-semibold">Orden</th>
                                        <th className="text-left p-3 font-semibold">Cliente</th>
                                        <th className="text-left p-3 font-semibold">Transportadora</th>
                                        <th className="text-right p-3 font-semibold">Precio</th>
                                        <th className="text-center p-3 font-semibold">Estado</th>
                                        <th className="text-center p-3 font-semibold">Acciones</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {orders.filter(o => selectedOrderIds.has(o.id)).map(order => (
                                        <tr key={order.id} className="border-b hover:bg-gray-50">
                                            <td className="p-3 font-medium">{order.order_number}</td>
                                            <td className="p-3">
                                                <div>{order.customer_name}</div>
                                                <div className="text-xs text-gray-500">{order.shipping_city}</div>
                                            </td>
                                            <td className="p-3">
                                                {order.quote ? (
                                                    <div>
                                                        <div className="font-medium">{order.quote.carrier}</div>
                                                        <div className="text-xs text-gray-500">{order.quote.product}</div>
                                                        <div className="text-xs text-gray-500">{order.quote.deliveryDays} días</div>
                                                    </div>
                                                ) : (
                                                    <span className="text-red-600 text-xs">-</span>
                                                )}
                                            </td>
                                            <td className="p-3 text-right">
                                                {order.quote ? (
                                                    <div className="font-bold text-orange-600">
                                                        ${(order.quote.flete + (order.quote.minimumInsurance ?? 0)).toLocaleString()}
                                                    </div>
                                                ) : (
                                                    <span className="text-red-600">-</span>
                                                )}
                                            </td>
                                            <td className="p-3 text-center">
                                                {order.quote ? (
                                                    <span className="inline-block px-2 py-1 bg-green-100 text-green-700 rounded text-xs">
                                                        ✓ Cotizada
                                                    </span>
                                                ) : (
                                                    <span className="inline-block px-2 py-1 bg-red-100 text-red-700 rounded text-xs">
                                                        ✗ Error
                                                    </span>
                                                )}
                                            </td>
                                            <td className="p-3 text-center">
                                                <button
                                                    onClick={() => setSelectedOrderForDetails(order)}
                                                    className="text-gray-400 hover:text-gray-600"
                                                >
                                                    ⋮
                                                </button>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>

                        {walletBalance !== null && walletBalance < totalCost && (
                            <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
                                ⚠️ Saldo insuficiente. Necesitas ${(totalCost - walletBalance).toLocaleString()} COP adicionales.
                            </div>
                        )}

                        <div className="flex justify-end space-x-3">
                            <Button variant="outline" onClick={() => setStep('select')}>
                                Volver
                            </Button>
                            <Button
                                onClick={handleGenerateAll}
                                disabled={walletBalance !== null && walletBalance < totalCost}
                            >
                                Generar Guías ({orders.filter(o => o.quote).length})
                            </Button>
                        </div>
                    </div>
                )}

                {/* Step 4: Generating Progress */}
                {step === 'generate' && (
                    <div className="space-y-4">
                        <p className="text-center text-gray-600">Generando guías...</p>
                        <div className="w-full bg-gray-200 rounded-full h-4">
                            <div
                                className="bg-green-500 h-4 rounded-full transition-all duration-300"
                                style={{ width: `${generatingProgress}%` }}
                            />
                        </div>
                        <p className="text-center text-sm text-gray-500">
                            {Math.round(generatingProgress)}% completado
                        </p>
                        <div className="text-center text-sm">
                            <p className="text-green-600">Exitosas: {generatedCount}</p>
                            <p className="text-red-600">Fallidas: {failedCount}</p>
                        </div>
                    </div>
                )}

                {/* Step 5: Complete */}
                {step === 'complete' && (
                    <div className="space-y-4">
                        <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                            <h3 className="font-semibold text-green-800 mb-2">✅ Proceso Completado</h3>
                            <div className="space-y-1 text-sm text-green-700">
                                <p>Guías generadas exitosamente: {generatedCount}</p>
                                <p>Guías fallidas: {failedCount}</p>
                            </div>
                        </div>

                        {generationErrors.length > 0 && (
                            <div className="bg-red-50 border border-red-200 rounded-lg p-4 max-h-48 overflow-y-auto">
                                <h4 className="font-semibold text-red-800 mb-2">Errores:</h4>
                                <ul className="text-sm text-red-700 space-y-1">
                                    {generationErrors.map((err, idx) => (
                                        <li key={idx}>• {err}</li>
                                    ))}
                                </ul>
                            </div>
                        )}

                        <div className="flex justify-end">
                            <Button onClick={handleClose}>
                                Cerrar
                            </Button>
                        </div>
                    </div>
                )}

                {error && (
                    <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
                        {error}
                    </div>
                )}
            </div>

            {/* Details Modal */}
            {selectedOrderForDetails && (
                <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-[60] p-4">
                    <div className="bg-white rounded-lg shadow-2xl max-w-lg w-full p-6 space-y-4">
                        <div className="flex justify-between items-center">
                            <h3 className="text-lg font-bold">Detalles de Cotización - {selectedOrderForDetails.order_number}</h3>
                            <button onClick={() => setSelectedOrderForDetails(null)} className="text-gray-400 hover:text-gray-600 text-xl">&times;</button>
                        </div>

                        <div className="grid grid-cols-2 gap-4 text-sm">
                            <div className="space-y-1">
                                <p className="font-bold text-gray-500 uppercase text-[10px]">Origen</p>
                                <p>Bogotá D.C. (Fallback)</p>
                                <p className="text-xs text-gray-500">Calle 1 # 1-1</p>
                            </div>
                            <div className="space-y-1">
                                <p className="font-bold text-gray-500 uppercase text-[10px]">Destino</p>
                                <p>{selectedOrderForDetails.shipping_city}</p>
                                <p className="text-xs text-gray-500">{selectedOrderForDetails.shipping_street}</p>
                            </div>
                            <div className="space-y-1">
                                <p className="font-bold text-gray-500 uppercase text-[10px]">Paquete</p>
                                <p>{selectedOrderForDetails.weight}kg</p>
                                <p className="text-xs text-gray-500">{selectedOrderForDetails.height}x{selectedOrderForDetails.width}x{selectedOrderForDetails.length} cm</p>
                            </div>
                            <div className="space-y-1">
                                <p className="font-bold text-gray-500 uppercase text-[10px]">Valor Declarado</p>
                                <p>${selectedOrderForDetails.total_amount?.toLocaleString()}</p>
                            </div>
                        </div>

                        {selectedOrderForDetails.quote ? (
                            <div className="bg-orange-50 p-3 rounded-lg border border-orange-100">
                                <p className="font-bold text-orange-800 text-sm mb-2">Transportadora: {selectedOrderForDetails.quote.carrier}</p>
                                <div className="grid grid-cols-2 gap-2 text-xs text-orange-700">
                                    <p>Flete: ${selectedOrderForDetails.quote.flete.toLocaleString()}</p>
                                    <p>Seguro: ${(selectedOrderForDetails.quote.minimumInsurance ?? 0).toLocaleString()}</p>
                                    <p className="font-bold">Total: ${(selectedOrderForDetails.quote.flete + (selectedOrderForDetails.quote.minimumInsurance ?? 0)).toLocaleString()}</p>
                                    <p>Entrega: {selectedOrderForDetails.quote.deliveryDays} días</p>
                                </div>
                            </div>
                        ) : (
                            <div className="bg-red-50 p-3 rounded-lg border border-red-100 text-red-700 text-sm">
                                Error: {selectedOrderForDetails.quoteError || 'No se pudo cotizar'}
                            </div>
                        )}

                        <div className="flex justify-end gap-2 pt-2">
                            <Button variant="outline" size="sm" onClick={() => setSelectedOrderForDetails(null)}>Cerrar</Button>
                            <Button size="sm" onClick={() => {
                                // Logic to re-quote could go here
                                setSelectedOrderForDetails(null);
                            }}>Volver a Cotizar</Button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
