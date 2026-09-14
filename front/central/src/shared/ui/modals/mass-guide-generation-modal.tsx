'use client';

import { useState, useEffect } from 'react';
import { Button } from '@/shared/ui';
import '@/shared/ui/styles/shipment-modals.css';
import { Order } from '@/services/modules/orders/domain/types';
import { getOrdersAction, updateOrderAction } from '@/services/modules/orders/infra/actions';
import { quoteShipmentAction, generateGuideAction } from '@/services/modules/shipments/infra/actions';
import { EnvioClickQuoteRequest, EnvioClickRate } from '@/services/modules/shipments/domain/types';
import { getWalletBalanceAction } from '@/services/modules/wallet/infra/actions';
import { getWarehousesAction } from '@/services/modules/warehouses/infra/actions';
import { Warehouse } from '@/services/modules/warehouses/domain/types';
import danes from "@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json";
import { getActionError } from '@/shared/utils/action-result';
import { buildGuideDestination } from '@/shared/utils/guide-destination';

const normalizeString = (str: string) =>
    str.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toUpperCase().trim();

const findDaneCode = (city: string, state: string) => {
    const targetCity = normalizeString(city);
    const targetState = normalizeString(state);

    const entries = Object.entries(danes);

    const exactMatch = entries.find(([_, data]: [string, any]) =>
        normalizeString(data.ciudad) === targetCity &&
        normalizeString(data.departamento) === targetState
    );
    if (exactMatch) return exactMatch[0];

    const cityMatch = entries.find(([_, data]: [string, any]) =>
        normalizeString(data.ciudad) === targetCity
    );
    if (cityMatch) return cityMatch[0];

    return null;
};

interface MassGuideGenerationModalProps {
    isOpen: boolean;
    onClose: () => void;
    onComplete?: (count: number) => void;
}

interface OrderWithQuote extends Order {
    quote?: EnvioClickRate;
    quotePending?: boolean;
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
    const [selectedOrderForEdit, setSelectedOrderForEdit] = useState<OrderWithQuote | null>(null);
    const [editForm, setEditForm] = useState<Partial<OrderWithQuote>>({});
    const [editedPackageIds, setEditedPackageIds] = useState<Set<string>>(new Set());
    const [isSaving, setIsSaving] = useState(false);
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [selectedWarehouse, setSelectedWarehouse] = useState<Warehouse | null>(null);


    useEffect(() => {
        if (isOpen && step === 'select') {
            loadOrders();
            loadWalletBalance();
            loadWarehouses();
        }
    }, [isOpen, step]);

    const loadWarehouses = async () => {
        try {
            const res = await getWarehousesAction({ is_active: true, page: 1, page_size: 100 });
            if (res.data) {
                setWarehouses(res.data);
                const defaultWh = res.data.find((w: Warehouse) => w.is_default) || res.data[0];
                if (defaultWh) setSelectedWarehouse(defaultWh);
            }
        } catch (err) {
            console.error('Error loading warehouses:', err);
        }
    };

    const loadOrders = async () => {
        setLoading(true);
        try {
            const response = await getOrdersAction({ page: 1, page_size: 100 });
            if (response.success && response.data) {
                const ordersWithoutGuides = response.data.filter(order => !order.tracking_number);
                setOrders(ordersWithoutGuides);
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar \u00f3rdenes'));
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

    const handleOpenEdit = (e: React.MouseEvent, order: OrderWithQuote) => {
        e.stopPropagation();
        setSelectedOrderForEdit(order);
        setEditForm({
            shipping_city: order.shipping_city || '',
            shipping_state: order.shipping_state || '',
            shipping_street: order.shipping_street || '',
            customer_phone: order.customer_phone || '',
            weight: order.weight || 1,
            length: order.length || 10,
            height: order.height || 10,
            width: order.width || 10,
        });
    };

    const handleSaveEdit = async () => {
        if (!selectedOrderForEdit) return;
        setIsSaving(true);
        try {
            const payload = {
                shipping_city: editForm.shipping_city,
                shipping_state: editForm.shipping_state,
                shipping_street: editForm.shipping_street,
                customer_phone: editForm.customer_phone,
                weight: Number(editForm.weight),
                length: Number(editForm.length),
                height: Number(editForm.height),
                width: Number(editForm.width),
            };
            const res = await updateOrderAction(selectedOrderForEdit.id, payload);
            if (res.success) {
                setOrders(orders.map(o => o.id === selectedOrderForEdit.id ? { ...o, ...payload } as OrderWithQuote : o));
                setEditedPackageIds(prev => new Set(prev).add(selectedOrderForEdit.id));
                setSelectedOrderForEdit(null);
            } else {
                alert(res.message || 'Error guardando');
            }
        } catch (err) {
            alert('Error conectando con el servidor');
        } finally {
            setIsSaving(false);
        }
    };

    const handleQuoteAll = async () => {
        setStep('quote');
        setLoading(true);
        setError(null);
        setQuotingProgress(0);

        const selectedOrders = orders.filter(o => selectedOrderIds.has(o.id));
        const quotedOrders: OrderWithQuote[] = [];
        const total = 0;

        for (let i = 0; i < selectedOrders.length; i++) {
            const order = selectedOrders[i];
            try {
                const destDane = findDaneCode(order.shipping_city || "", order.shipping_state || "");
                const orderCodValue = (order.cod_total && order.cod_total > 0) ? order.cod_checkout_total : undefined;
                const quotePayload: EnvioClickQuoteRequest = {
                    order_uuid: order.id,
                    auto_package: !editedPackageIds.has(order.id),
                    packages: [{
                        weight: order.weight || 1,
                        height: order.height || 10,
                        width: order.width || 10,
                        length: order.length || 10,
                    }],
                    description: `Orden ${order.order_number}`,
                    contentValue: order.total_amount || 10000,
                    codValue: orderCodValue,
                    includeGuideCost: false,
                    codPaymentMethod: orderCodValue ? 'cash' : '',
                    origin: {
                        daneCode: selectedWarehouse?.city_dane_code || '',
                        address: selectedWarehouse?.street || selectedWarehouse?.address || 'Direcci\u00f3n no especificada',
                    },
                    destination: {
                        daneCode: destDane || '',
                        address: order.shipping_street || 'Direcci\u00f3n no especificada',
                    },
                };

                const response = await quoteShipmentAction(quotePayload);
                if (response.success) {
                    quotedOrders.push({ ...order, quotePending: true });
                } else {
                    quotedOrders.push({ ...order, quoteError: response.message || 'Error al solicitar cotizaci\u00f3n' });
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
                const destDane = findDaneCode(order.shipping_city || "", order.shipping_state || "");
                const destParts = buildGuideDestination(order);

                const genCodValue = (order.cod_total && order.cod_total > 0) ? order.cod_checkout_total : undefined;
                const guideTotalCost = (order.quote!.flete) + (order.quote!.minimumInsurance ?? 0) + (order.quote!.extraInsurance ?? 0);
                const generatePayload: EnvioClickQuoteRequest = {
                    idRate: order.quote!.idRate,
                    totalCost: guideTotalCost,
                    myShipmentReference: "Orden " + (order.internal_number || order.order_number),
                    external_order_id: order.order_number,
                    order_uuid: order.id,
                    auto_package: !editedPackageIds.has(order.id),
                    requestPickup: false,
                    pickupDate: new Date().toISOString().split('T')[0],
                    insurance: true,
                    description: `Orden ${order.order_number}`,
                    contentValue: order.total_amount || 10000,
                    codValue: genCodValue,
                    includeGuideCost: false,
                    codPaymentMethod: genCodValue ? 'cash' : '',
                    packages: [{
                        weight: order.weight || 1,
                        height: order.height || 10,
                        width: order.width || 10,
                        length: order.length || 10,
                    }],
                    origin: {
                        daneCode: selectedWarehouse?.city_dane_code || '',
                        address: selectedWarehouse?.street || selectedWarehouse?.address || 'Direcci\u00f3n no especificada',
                        company: selectedWarehouse?.company || selectedWarehouse?.name || 'Mi Empresa',
                        firstName: selectedWarehouse?.first_name || selectedWarehouse?.contact_name?.split(' ')[0] || 'Admin',
                        lastName: selectedWarehouse?.last_name || selectedWarehouse?.contact_name?.split(' ').slice(1).join(' ') || '',
                        email: selectedWarehouse?.email || selectedWarehouse?.contact_email || '',
                        phone: selectedWarehouse?.phone || '',
                        suburb: selectedWarehouse?.suburb || '',
                        crossStreet: selectedWarehouse?.street || selectedWarehouse?.address || '',
                        reference: '',
                    },
                    destination: {
                        daneCode: destDane || '',
                        address: destParts.address || 'Direcci\u00f3n no especificada',
                        company: order.customer_name || 'Cliente',
                        firstName: order.customer_name?.split(' ')[0] || 'Cliente',
                        lastName: order.customer_name?.split(' ').slice(1).join(' ') || 'Apellido',
                        email: order.customer_email || 'cliente@example.com',
                        phone: order.customer_phone || '3009876543',
                        suburb: destParts.suburb,
                        crossStreet: destParts.crossStreet || destParts.address,
                        reference: destParts.reference,
                    },
                };

                const genRes = await generateGuideAction(generatePayload);
                if (genRes.success && genRes.data) {
                    setGeneratedCount(prev => prev + 1);
                } else {
                    throw new Error(genRes.message || "Error generando gu\u00eda");
                }
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
        setSelectedWarehouse(null);
        onClose();
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto p-6">
                <div className="flex justify-between items-center mb-6">
                    <h2 className="text-2xl font-bold text-gray-800 dark:text-gray-100">Generaci&#243;n Masiva de Gu&#237;as</h2>
                    <button onClick={handleClose} className="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:text-gray-200 dark:text-gray-200 text-2xl">
                        &#215;
                    </button>
                </div>

                {step === 'select' && (
                    <div className="space-y-4">
                        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                            <label className="block text-sm font-semibold text-blue-800 mb-2">Bodega de origen</label>
                            {warehouses.length === 0 ? (
                                <p className="text-sm text-blue-600">Cargando bodegas...</p>
                            ) : (
                                <select
                                    value={selectedWarehouse?.id?.toString() || ''}
                                    onChange={(e) => {
                                        const wh = warehouses.find(w => w.id === parseInt(e.target.value));
                                        setSelectedWarehouse(wh || null);
                                    }}
                                    className="w-full border border-blue-300 rounded-md p-2 text-sm bg-white"
                                >
                                    {warehouses.map(wh => (
                                        <option key={wh.id} value={wh.id}>
                                            {wh.name} &#8212; {wh.city}, {wh.state} {wh.is_default ? '(Por defecto)' : ''}
                                        </option>
                                    ))}
                                </select>
                            )}
                            {selectedWarehouse && (
                                <p className="text-xs text-blue-600 mt-1">
                                    DANE: {selectedWarehouse.city_dane_code} | {selectedWarehouse.street || selectedWarehouse.address}
                                </p>
                            )}
                        </div>

                        <div className="flex justify-between items-center">
                            <p className="text-sm text-gray-600 dark:text-gray-300">
                                Selecciona las &#243;rdenes para generar gu&#237;as ({selectedOrderIds.size} seleccionadas)
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
                            <div className="text-center py-8">Cargando &#243;rdenes...</div>
                        ) : orders.length === 0 ? (
                            <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                                No hay &#243;rdenes sin gu&#237;a de env&#237;o
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
                                            <div className="flex items-center gap-2">
                                                <span className="font-semibold">{order.order_number}</span>
                                                {order.cod_total && order.cod_total > 0 && (
                                                    <span className="shipment-badge-warning text-[10px] px-1.5 py-0.5">
                                                        Contra Entrega ${(order.cod_customer_charge ?? 0).toLocaleString()}
                                                    </span>
                                                )}
                                            </div>
                                            <div className="text-sm text-gray-600 dark:text-gray-300">
                                                {order.customer_name} - {order.shipping_city || <span className="text-red-500 text-xs font-semibold">Sin ciudad (Edita para continuar)</span>}
                                            </div>
                                        </div>
                                        <div className="text-right mr-4">
                                            <div className="font-semibold">${order.total_amount?.toLocaleString()}</div>
                                            <div className="text-xs text-gray-500 dark:text-gray-400">
                                                {order.weight || 1}kg
                                            </div>
                                        </div>
                                        <button
                                            onClick={(e) => handleOpenEdit(e, order)}
                                            className="p-2 text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded-full transition-colors"
                                            title="Editar Orden"
                                        >
                                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                                <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                                            </svg>
                                        </button>
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

                {step === 'quote' && (
                    <div className="space-y-4">
                        <p className="text-center text-gray-600 dark:text-gray-300">Cotizando env&#237;os...</p>
                        <div className="shipment-progress-bar rounded-full h-4">
                            <div
                                className="shipment-progress-fill"
                                style={{ width: `${quotingProgress}%` }}
                            />
                        </div>
                        <p className="text-center text-sm text-gray-500 dark:text-gray-400">
                            {Math.round(quotingProgress)}% completado
                        </p>
                    </div>
                )}

                {step === 'confirm' && (
                    <div className="space-y-4">
                        <div className="shipment-alert shipment-alert-primary rounded-lg">
                            <h3 className="font-semibold mb-2" style={{ color: 'var(--color-primary)' }}>Resumen de Cotizaci&#243;n</h3>
                            <div className="flex justify-between items-center text-sm" style={{ color: 'var(--color-primary)' }}>
                                <div>
                                    <p>&#211;rdenes cotizadas: {orders.filter(o => o.quote).length}</p>
                                    <p>&#211;rdenes con error: {orders.filter(o => o.quoteError).length}</p>
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

                        <div className="border rounded-lg max-h-96 overflow-y-auto">
                            <table className="w-full text-sm">
                                <thead className="sticky top-0 shipment-table-header">
                                    <tr className="border-b">
                                        <th className="text-left p-3">Orden</th>
                                        <th className="text-left p-3">Cliente</th>
                                        <th className="text-left p-3">Transportadora</th>
                                        <th className="text-right p-3">Precio</th>
                                        <th className="text-center p-3">Estado</th>
                                        <th className="text-center p-3">Acciones</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {orders.filter(o => selectedOrderIds.has(o.id)).map(order => (
                                        <tr key={order.id} className="border-b hover:bg-gray-50">
                                            <td className="p-3">
                                                <div className="flex items-center gap-1.5">
                                                    <span className="font-medium">{order.order_number}</span>
                                                    {order.cod_total && order.cod_total > 0 && (
                                                        <span className="shipment-badge-warning text-[9px] px-1 py-0.5">Contra Entrega</span>
                                                    )}
                                                </div>
                                            </td>
                                            <td className="p-3">
                                                <div>{order.customer_name}</div>
                                                <div className="text-xs text-gray-500 dark:text-gray-400">{order.shipping_city}</div>
                                            </td>
                                            <td className="p-3">
                                                {order.quote ? (
                                                    <div>
                                                        <div className="font-medium">{order.quote.carrier}</div>
                                                        <div className="text-xs text-gray-500 dark:text-gray-400">{order.quote.product}</div>
                                                        <div className="text-xs text-gray-500 dark:text-gray-400">{order.quote.deliveryDays} d&#237;as</div>
                                                    </div>
                                                ) : (
                                                    <span className="text-red-600 text-xs">-</span>
                                                )}
                                            </td>
                                            <td className="p-3 text-right">
                                                {order.quote ? (
                                                    <div className="font-bold text-orange-600">
                                                        ${(order.quote.flete + (order.quote.minimumInsurance ?? 0) + (order.quote.extraInsurance ?? 0)).toLocaleString()}
                                                    </div>
                                                ) : (
                                                    <span className="text-red-600">-</span>
                                                )}
                                            </td>
                                            <td className="p-3 text-center">
                                                {order.quote ? (
                                                    <span className="shipment-badge-success px-2 py-1">
                                                        &#10003; Cotizada
                                                    </span>
                                                ) : (
                                                    <span className="shipment-badge-error px-2 py-1">
                                                        &#10007; Error
                                                    </span>
                                                )}
                                            </td>
                                            <td className="p-3 text-center">
                                                <button
                                                    onClick={() => setSelectedOrderForDetails(order)}
                                                    className="text-gray-400 hover:text-gray-600 dark:text-gray-300"
                                                >
                                                    &#8942;
                                                </button>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>

                        {walletBalance !== null && walletBalance < totalCost && (
                            <div className="shipment-alert shipment-alert-error">
                                &#9888;&#65039; Saldo insuficiente. Necesitas ${(totalCost - walletBalance).toLocaleString()} COP adicionales.
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
                                Generar Gu&#237;as ({orders.filter(o => o.quote).length})
                            </Button>
                        </div>
                    </div>
                )}

                {step === 'generate' && (
                    <div className="space-y-4">
                        <p className="text-center text-gray-600 dark:text-gray-300">Generando gu&#237;as...</p>
                        <div className="shipment-progress-bar rounded-full h-4">
                            <div
                                className="shipment-progress-fill"
                                style={{ width: `${generatingProgress}%` }}
                            />
                        </div>
                        <p className="text-center text-sm text-gray-500 dark:text-gray-400">
                            {Math.round(generatingProgress)}% completado
                        </p>
                        <div className="text-center text-sm">
                            <p className="text-green-600">Exitosas: {generatedCount}</p>
                            <p className="text-red-600">Fallidas: {failedCount}</p>
                        </div>
                    </div>
                )}

                {step === 'complete' && (
                    <div className="space-y-4">
                        <div className="shipment-alert shipment-alert-success">
                            <h3 className="font-semibold mb-2">&#9989; Proceso Completado</h3>
                            <div className="space-y-1 text-sm">
                                <p>Gu&#237;as generadas exitosamente: {generatedCount}</p>
                                <p>Gu&#237;as fallidas: {failedCount}</p>
                            </div>
                        </div>

                        {generationErrors.length > 0 && (
                            <div className="shipment-alert shipment-alert-error max-h-48 overflow-y-auto">
                                <h4 className="font-semibold mb-2">Errores:</h4>
                                <ul className="text-sm space-y-1">
                                    {generationErrors.map((err, idx) => (
                                        <li key={idx}>&#8226; {err}</li>
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
                    <div className="mt-4 shipment-alert shipment-alert-error">
                        {error}
                    </div>
                )}
            </div>

            {selectedOrderForDetails && (
                <div className="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-[60] p-4">
                    <div className="bg-white rounded-lg shadow-2xl max-w-lg w-full p-6 space-y-4">
                        <div className="flex justify-between items-center">
                            <h3 className="text-lg font-bold">Detalles de Cotizaci&#243;n - {selectedOrderForDetails.order_number}</h3>
                            <button onClick={() => setSelectedOrderForDetails(null)} className="text-gray-400 hover:text-gray-600 dark:text-gray-300 text-xl">&times;</button>
                        </div>

                        <div className="grid grid-cols-2 gap-4 text-sm">
                            <div className="space-y-1">
                                <p className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Origen</p>
                                <p>{selectedWarehouse ? `${selectedWarehouse.city}, ${selectedWarehouse.state}` : 'Sin bodega seleccionada'}</p>
                                <p className="text-xs text-gray-500 dark:text-gray-400">{selectedWarehouse?.street || selectedWarehouse?.address || 'Sin direcci\u00f3n'}</p>
                            </div>
                            <div className="space-y-1">
                                <p className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Destino</p>
                                <p>{selectedOrderForDetails.shipping_city}</p>
                                <p className="text-xs text-gray-500 dark:text-gray-400">{selectedOrderForDetails.shipping_street}</p>
                            </div>
                            <div className="space-y-1">
                                <p className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Paquete</p>
                                <p>{selectedOrderForDetails.weight}kg</p>
                                <p className="text-xs text-gray-500 dark:text-gray-400">{selectedOrderForDetails.height}x{selectedOrderForDetails.width}x{selectedOrderForDetails.length} cm</p>
                            </div>
                            <div className="space-y-1">
                                <p className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Valor Declarado</p>
                                <p>${selectedOrderForDetails.total_amount?.toLocaleString()}</p>
                            </div>
                        </div>

                        {selectedOrderForDetails.quote ? (
                            <div className="shipment-alert shipment-alert-warning p-3">
                                <p className="font-bold text-sm mb-2">Transportadora: {selectedOrderForDetails.quote.carrier}</p>
                                <div className="grid grid-cols-2 gap-2 text-xs">
                                    <p>Flete: ${selectedOrderForDetails.quote.flete.toLocaleString()}</p>
                                    <p>Seg. obligatorio: ${(selectedOrderForDetails.quote.minimumInsurance ?? 0).toLocaleString()}</p>
                                    <p>Seg. adicional: ${(selectedOrderForDetails.quote.extraInsurance ?? 0).toLocaleString()} <span className="text-emerald-700">(incluido)</span></p>
                                    <p className="font-bold">Total: ${(selectedOrderForDetails.quote.flete + (selectedOrderForDetails.quote.minimumInsurance ?? 0) + (selectedOrderForDetails.quote.extraInsurance ?? 0)).toLocaleString()}</p>
                                    <p>Entrega: {selectedOrderForDetails.quote.deliveryDays} d&#237;as</p>
                                </div>
                            </div>
                        ) : (
                            <div className="shipment-alert shipment-alert-error text-sm">
                                Error: {selectedOrderForDetails.quoteError || 'No se pudo cotizar'}
                            </div>
                        )}

                        <div className="flex justify-end gap-2 pt-2">
                            <Button variant="outline" size="sm" onClick={() => setSelectedOrderForDetails(null)}>Cerrar</Button>
                            <Button size="sm" onClick={() => {
                                setSelectedOrderForDetails(null);
                            }}>Cerrar</Button>
                        </div>
                    </div>
                </div>
            )}

            {selectedOrderForEdit && (
                <div className="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-[70] p-4">
                    <div className="bg-white rounded-lg shadow-2xl max-w-lg w-full p-6 space-y-4">
                        <div className="flex justify-between items-center">
                            <h3 className="text-lg font-bold">Editar Orden - {selectedOrderForEdit.order_number}</h3>
                            <button onClick={() => setSelectedOrderForEdit(null)} className="text-gray-400 hover:text-gray-600 dark:text-gray-300 text-xl">&times;</button>
                        </div>

                        <div className="grid grid-cols-2 gap-4 text-sm">
                            <div className="space-y-1 col-span-2">
                                <label className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Ciudad Destino</label>
                                <input type="text" className="w-full border p-2 rounded" value={editForm.shipping_city || ''} onChange={e => setEditForm({ ...editForm, shipping_city: e.target.value })} placeholder={'Ej. Bogot\u00e1'} />
                            </div>
                            <div className="space-y-1 col-span-2">
                                <label className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Departamento</label>
                                <input type="text" className="w-full border p-2 rounded" value={editForm.shipping_state || ''} onChange={e => setEditForm({ ...editForm, shipping_state: e.target.value })} placeholder="Ej. Cundinamarca" />
                            </div>
                            <div className="space-y-1 col-span-2">
                                <label className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Direcci&#243;n</label>
                                <input type="text" className="w-full border p-2 rounded" value={editForm.shipping_street || ''} onChange={e => setEditForm({ ...editForm, shipping_street: e.target.value })} />
                            </div>
                            <div className="space-y-1 col-span-2">
                                <label className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Tel&#233;fono Cliente</label>
                                <input type="text" className="w-full border p-2 rounded" value={editForm.customer_phone || ''} onChange={e => setEditForm({ ...editForm, customer_phone: e.target.value })} />
                            </div>
                            <div className="space-y-1">
                                <label className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Peso (kg)</label>
                                <input type="number" className="w-full border p-2 rounded" value={editForm.weight || 1} onChange={e => setEditForm({ ...editForm, weight: Number(e.target.value) })} />
                            </div>
                            <div className="space-y-1">
                                <label className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Largo (cm)</label>
                                <input type="number" className="w-full border p-2 rounded" value={editForm.length || 10} onChange={e => setEditForm({ ...editForm, length: Number(e.target.value) })} />
                            </div>
                            <div className="space-y-1">
                                <label className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Ancho (cm)</label>
                                <input type="number" className="w-full border p-2 rounded" value={editForm.width || 10} onChange={e => setEditForm({ ...editForm, width: Number(e.target.value) })} />
                            </div>
                            <div className="space-y-1">
                                <label className="font-bold text-gray-500 dark:text-gray-400 uppercase text-[10px]">Alto (cm)</label>
                                <input type="number" className="w-full border p-2 rounded" value={editForm.height || 10} onChange={e => setEditForm({ ...editForm, height: Number(e.target.value) })} />
                            </div>
                        </div>

                        <div className="flex justify-end gap-2 pt-4">
                            <Button variant="outline" size="sm" onClick={() => setSelectedOrderForEdit(null)}>Cancelar</Button>
                            <Button size="sm" onClick={handleSaveEdit} disabled={isSaving}>
                                {isSaving ? 'Guardando...' : 'Guardar Datos'}
                            </Button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
