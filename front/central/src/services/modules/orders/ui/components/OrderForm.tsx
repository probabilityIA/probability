'use client';

import { useState, useRef, useEffect, useCallback } from 'react';
import { Order, CreateOrderDTO, UpdateOrderDTO } from '../../domain/types';
import { Product } from '../../../products/domain/types';
import { Button, Input, Alert, Modal } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useToast } from '@/shared/providers/toast-provider';
import ProductSelector from '../../../products/ui/components/ProductSelector';
import ProductForm from '../../../products/ui/components/ProductForm';
import { createOrderAction, updateOrderAction } from '../../infra/actions';
import danes from '@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json';
import { useClientSearch } from '../hooks/useClientSearch';
import { useWarehouses } from '../hooks/useWarehouses';
import ClientAutocomplete from './ClientAutocomplete';
import AddressAutocomplete, { AddressSuggestion } from './AddressAutocomplete';
import dynamic from 'next/dynamic';

const MapComponent = dynamic(() => import('@/shared/ui/MapComponent'), { ssr: false });
import { CustomerInfo } from '../../../customers/domain/types';
import { getCustomerAddressesAction } from '../../../customers/infra/actions';
import { getActionError } from '@/shared/utils/action-result';

interface OrderFormProps {
    order?: Order;
    onSuccess?: () => void;
    onCancel?: () => void;
    /** For super admin: the business selected in the page-level selector */
    selectedBusinessId?: number | null;
}

export default function OrderForm({ order, onSuccess, onCancel, selectedBusinessId }: OrderFormProps) {
    const isEdit = !!order;
    const { permissions } = usePermissions();
    const defaultBusinessId = selectedBusinessId || permissions?.business_id || 0;
    const { showToast } = useToast();
    const [cachedGuideCarrier, setCachedGuideCarrier] = useState<string | null>(null);

    // Load cached guide data from sessionStorage on mount
    useEffect(() => {
        console.log('📋 OrderForm mounted. isEdit:', isEdit, 'order_number:', order?.order_number);
        if (isEdit && order?.order_number) {
            const key = `guide_${order.order_number}`;
            const cached = sessionStorage.getItem(key);
            console.log('🔍 Looking for cache key:', key, 'Found:', cached);
            if (cached) {
                try {
                    const guideData = JSON.parse(cached);
                    console.log('✅ Parsed cache data:', guideData);
                    if (guideData.carrier) {
                        setCachedGuideCarrier(guideData.carrier);
                        console.log('🚚 Set carrier from cache:', guideData.carrier);
                    }
                } catch (e) {
                    console.error('❌ Error parsing cache:', e);
                }
            } else {
                console.log('⚠️ No cache found for this order');
            }
        }
    }, [isEdit, order?.order_number]);

    const [formData, setFormData] = useState({
        // Integration
        integration_id: order?.integration_id || 0,
        platform: order?.platform || 'manual',
        business_id: order?.business_id || defaultBusinessId,

        // Customer — fallback: split customer_name if first/last aren't set separately
        customer_name: order?.customer_name || '',
        customer_first_name: order?.customer_first_name || (order?.customer_name ? order.customer_name.split(' ')[0] : ''),
        customer_last_name: order?.customer_last_name || (order?.customer_name ? order.customer_name.split(' ').slice(1).join(' ') : ''),
        customer_email: order?.customer_email || '',
        customer_phone: order?.customer_phone || '',
        customer_dni: order?.customer_dni || '',

        shipping_street: order?.shipping_street ? order.shipping_street.split(' | ')[0] : '',
        shipping_city: order?.shipping_city || '',
        shipping_state: order?.shipping_state || '',
        shipping_country: order?.shipping_country || 'Colombia',
        shipping_postal_code: order?.shipping_postal_code || '',

        subtotal: order?.subtotal || 0,
        tax: order?.tax || 0,
        discount: order?.discount || 0,
        shipping_cost: order?.shipping_cost || 0,
        total_amount: order?.total_amount || 0,
        currency: order?.currency || 'COP',
        cod_total: order?.cod_total || 0,

        // Payment
        payment_method_id: order?.payment_method_id || 1,
        is_paid: order?.is_paid || false,

        // Status
        status: order?.status || 'pending',

        // Logistics (preserved on update)
        tracking_number: order?.tracking_number || '',
        tracking_link: order?.tracking_link || '',
        guide_id: order?.guide_id || '',
        warehouse_id: order?.warehouse_id || 0,
        warehouse_name: order?.warehouse_name || '',
        driver_name: order?.driver_name || '',
        is_last_mile: order?.is_last_mile || false,

        // Additional
        notes: order?.notes || '',
        invoiceable: order?.invoiceable ?? false,
        is_confirmed: order?.is_confirmed ?? null,
        novelty: order?.novelty || '',

        // Items
        items: order?.order_items || order?.items || [],

        // Extra
        integration_type: order?.integration_type || '',
        external_id: order?.external_id || '',
    });

    const [selectedProducts, setSelectedProducts] = useState<Product[]>(() => {
        // Try order_items first (structured items from backend), then fallback to items
        const items = order?.order_items ?? order?.items;
        if (!items || !Array.isArray(items)) return [];
        return (items as any[])
            .map((item: any) => ({
                ...item,
                // Map order_item fields to Product interface
                id: item.id?.toString() || item.product_id || '',
                sku: item.sku || item.product_sku || '',
                name: item.name || item.product_name || item.product_title || '',
                price: item.price ?? item.unit_price ?? 0,
                quantity: item.quantity ?? 1,
                stock: item.stock ?? item.stock_quantity ?? 0,
                manage_stock: item.manage_stock ?? item.track_inventory ?? false,
                thumbnail: item.thumbnail || item.image_url || undefined,
                currency: item.currency || order?.currency || 'COP',
            } as Product))
            .filter((p: any) => p.id);
    });

    const [shippingComplement, setShippingComplement] = useState('');
    const [showProductModal, setShowProductModal] = useState(false);

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // DANE search states
    const [citySearch, setCitySearch] = useState('');
    const [showCityResults, setShowCityResults] = useState(false);
    const cityRef = useRef<HTMLDivElement>(null);

    const [isCOD, setIsCOD] = useState(() => (order?.cod_total || 0) > 0);

    const [house, setHouse] = useState(() => {
        if (!order?.shipping_street) return '';
        const parts = order.shipping_street.split(' | ');
        return parts.length >= 2 ? parts[1] : '';
    });
    const [barrio, setBarrio] = useState(() => {
        if (!order?.shipping_street) return '';
        const parts = order.shipping_street.split(' | ');
        return parts.length >= 3 ? parts[2] : '';
    });

    const [addressCoords, setAddressCoords] = useState<{ lat: number; lon: number } | null>(null);
    const [addressAutofilled, setAddressAutofilled] = useState(false);


    // Client autocomplete (triggered from DNI, name, or email fields)
    const { results: clientResults, loading: clientLoading, searched: clientSearched, search: searchClients, clear: clearClients } = useClientSearch({
        businessId: formData.business_id,
        minChars: 3,
    });
    const { warehouses, loading: warehousesLoading } = useWarehouses({ businessId: formData.business_id });
    const [showClientDropdown, setShowClientDropdown] = useState(false);
    const [activeSearchField, setActiveSearchField] = useState<'dni' | 'name' | 'lastname' | 'email' | null>(null);

    const handleClientFieldChange = useCallback((field: 'dni' | 'name' | 'lastname' | 'email', value: string) => {
        if (field === 'dni') {
            setFormData(prev => ({ ...prev, customer_dni: value }));
        } else if (field === 'name') {
            setFormData(prev => ({ ...prev, customer_first_name: value }));
        } else if (field === 'lastname') {
            setFormData(prev => ({ ...prev, customer_last_name: value }));
        } else {
            setFormData(prev => ({ ...prev, customer_email: value }));
        }

        if (value.length >= 3) {
            searchClients(value);
            setShowClientDropdown(true);
            setActiveSearchField(field);
        } else {
            clearClients();
            setShowClientDropdown(false);
            setActiveSearchField(null);
        }
    }, [searchClients, clearClients]);

    const handleClientSelect = useCallback(async (client: CustomerInfo) => {
        const nameParts = client.name.split(' ');
        const firstName = nameParts[0] || '';
        const lastName = nameParts.slice(1).join(' ') || '';

        let phone = client.phone || '';
        if (phone.startsWith('+57')) phone = phone.slice(3);
        if (phone.startsWith('57') && phone.length > 10) phone = phone.slice(2);

        setFormData(prev => ({
            ...prev,
            customer_first_name: firstName,
            customer_last_name: lastName,
            customer_name: client.name,
            customer_email: client.email || '',
            customer_phone: phone,
            customer_dni: client.dni || prev.customer_dni,
        }));
        setShowClientDropdown(false);
        setActiveSearchField(null);
        clearClients();

        try {
            const addressRes = await getCustomerAddressesAction(client.id, {
                page: 1,
                page_size: 1,
                business_id: formData.business_id,
            });
            if (addressRes.data && addressRes.data.length > 0) {
                const addr = addressRes.data[0];
                const streetParts = addr.street ? addr.street.split(' | ') : [''];
                const mainStreet = streetParts[0] || '';
                const addrHouse = streetParts.length >= 2 ? streetParts[1] : '';
                const addrBarrio = streetParts.length >= 3 ? streetParts[2] : '';

                setFormData(prev => ({
                    ...prev,
                    shipping_street: mainStreet,
                    shipping_city: addr.city || '',
                    shipping_state: addr.state || '',
                    shipping_country: addr.country || 'Colombia',
                    shipping_postal_code: addr.postal_code || '',
                }));
                setHouse(addrHouse);
                setBarrio(addrBarrio);

                if (addr.city && addr.state) {
                    setCitySearch(`${addr.city} (${addr.state})`);
                }

                if (addr.latitude && addr.longitude) {
                    setAddressCoords({ lat: addr.latitude, lon: addr.longitude });
                }

                setAddressAutofilled(true);
            }
        } catch {
            // silently ignore - address autocomplete is best-effort
        }
    }, [clearClients, formData.business_id]);

    // DANE options
    const daneOptions = Object.entries(danes).map(([code, data]: [string, any]) => ({
        value: code,
        label: `${data.ciudad} (${data.departamento})`,
        ciudad: data.ciudad,
        departamento: data.departamento
    })).sort((a, b) => a.label.localeCompare(b.label));

    const filteredCityOptions = daneOptions.filter(opt =>
        opt.label.toLowerCase().includes(citySearch.toLowerCase())
    );

    // Auto-select warehouse: if editing use existing, else pick default or single warehouse
    useEffect(() => {
        if (warehouses.length === 0 || formData.warehouse_id) return;
        const defaultW = warehouses.find(w => w.is_default) || (warehouses.length === 1 ? warehouses[0] : null);
        if (defaultW) {
            setFormData(prev => ({ ...prev, warehouse_id: defaultW.id, warehouse_name: defaultW.name }));
        }
    }, [warehouses]);

    // Initialize citySearch when order is loaded
    useEffect(() => {
        if (order?.shipping_city && order?.shipping_state) {
            setCitySearch(`${order.shipping_city} (${order.shipping_state})`);
        }
    }, [order]);

    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (cityRef.current && !cityRef.current.contains(event.target as Node)) {
                setShowCityResults(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleCitySelect = (option: any) => {
        setFormData({
            ...formData,
            shipping_city: option.ciudad,
            shipping_state: option.departamento
        });
        setCitySearch(option.label);
        setShowCityResults(false);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            // Validation
            if ((!formData.customer_name && !formData.customer_first_name) || !formData.total_amount) {
                throw new Error('Por favor completa los campos requeridos');
            }

            const parts = [formData.shipping_street || ''];
            if (house.trim()) parts.push(house.trim());
            if (barrio.trim()) parts.push(barrio.trim());
            const fullShippingStreet = parts.join(' | ');

            const baseData = {
                ...formData,
                shipping_street: fullShippingStreet,
                shipping_lat: addressCoords?.lat,
                shipping_lng: addressCoords?.lon,
                items: selectedProducts.length > 0 ? selectedProducts : formData.items,
                customer_name: formData.customer_name || `${formData.customer_first_name} ${formData.customer_last_name}`.trim()
            };

            let response;
            if (isEdit && order) {
                const updateData: UpdateOrderDTO = {
                    ...baseData,
                    // Map is_confirmed → confirmation_status so clearing to null works
                    confirmation_status: formData.is_confirmed === true ? 'yes' : formData.is_confirmed === false ? 'no' : 'pending',
                    is_confirmed: undefined,
                };
                response = await updateOrderAction(order.id, updateData);
            } else {
                response = await createOrderAction(baseData as CreateOrderDTO);
            }

            if (response.success) {
                showToast(isEdit ? 'Orden actualizada exitosamente' : 'Orden creada exitosamente', 'success');
                if (onSuccess) onSuccess();
            } else {
                setError(response.message || 'Error al guardar la orden');
                showToast(response.message || 'Error al guardar la orden', 'error');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al guardar la orden'));
            showToast(err.message || 'Error al guardar la orden', 'error');
        } finally {
            setLoading(false);
        }
    };

    const handleProductsChange = (products: Product[]) => {
        setSelectedProducts(products);
        const subtotal = products.reduce((acc, p) => acc + p.price, 0);
        setFormData(prev => ({
            ...prev,
            subtotal,
            total_amount: subtotal + prev.tax - prev.discount + prev.shipping_cost,
        }));
    };

    // Auto-calculate total
    const calculateTotal = () => {
        setFormData(prev => ({
            ...prev,
            total_amount: prev.subtotal + prev.tax - prev.discount + prev.shipping_cost,
        }));
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4 bg-white dark:bg-gray-800 p-6 rounded-lg">
            {isEdit && order && (
                <div className="flex items-center gap-3 px-4 py-3 mb-4 bg-gradient-to-r from-purple-600 to-purple-500 rounded-lg shadow-sm">
                    {order.integration_logo_url && (
                        <img src={order.integration_logo_url} alt="" className="h-8 w-8 object-contain rounded" />
                    )}
                    <div>
                        <p className="text-xs text-white font-medium uppercase tracking-wide">
                            {order.integration_name || order.integration_type || 'Integración'}
                        </p>
                        <p className="text-sm font-bold text-white">
                            {order.order_number || order.internal_number || order.id}
                        </p>
                    </div>
                </div>
            )}
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
                <div>
                    <div className="bg-gray-50/80 dark:bg-gray-700/20 border border-gray-200 dark:border-gray-600/30 p-5 rounded-xl">
                        <div className="flex items-center gap-2 mb-4 pb-3 border-b border-gray-200/50 dark:border-gray-600/30">
                            <div className="w-7 h-7 rounded-lg bg-purple-100 dark:bg-purple-800/40 flex items-center justify-center text-purple-600 dark:text-purple-400 text-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                            </div>
                            <h3 className="text-base font-semibold text-purple-700 dark:text-purple-400">Cliente</h3>
                        </div>

                        {/* Warning banner */}
                        {!formData.customer_email && (formData.customer_first_name || formData.customer_name) && (
                            <div className="mb-4 px-3 py-2 bg-amber-50 border border-amber-200 rounded-lg">
                                <p className="text-xs text-amber-700">
                                    Completar el email del cliente mejora la analítica, seguimiento y facturación de tus órdenes.
                                </p>
                            </div>
                        )}

                        <div className="space-y-4">
                            {/* DNI — with autocomplete */}
                            <div className="relative">
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    DNI / Cédula
                                </label>
                                <div className="relative">
                                    <Input
                                        type="text"
                                        value={formData.customer_dni}
                                        onChange={(e) => handleClientFieldChange('dni', e.target.value)}
                                        onFocus={() => {
                                            if (formData.customer_dni.length >= 3 && (clientResults.length > 0 || clientLoading || clientSearched)) {
                                                setShowClientDropdown(true);
                                                setActiveSearchField('dni');
                                            }
                                        }}
                                        placeholder="Buscar cliente por cédula..."
                                        autoComplete="off"
                                        className={clientLoading && activeSearchField === 'dni' ? 'pr-10' : ''}
                                    />
                                    {clientLoading && activeSearchField === 'dni' && (
                                        <div className="absolute right-3 top-1/2 -translate-y-1/2">
                                            <div className="w-4 h-4 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" />
                                        </div>
                                    )}
                                </div>
                                {activeSearchField === 'dni' && (
                                    <ClientAutocomplete
                                        results={clientResults}
                                        loading={clientLoading}
                                        searched={clientSearched}
                                        visible={showClientDropdown}
                                        searchTerm={formData.customer_dni}
                                        onSelect={handleClientSelect}
                                        onClose={() => setShowClientDropdown(false)}
                                    />
                                )}
                            </div>
                            {/* Nombre — with autocomplete */}
                            <div className="relative">
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Nombre *
                                </label>
                                <div className="relative">
                                    <Input
                                        type="text"
                                        required
                                        value={formData.customer_first_name}
                                        onChange={(e) => handleClientFieldChange('name', e.target.value)}
                                        onFocus={() => {
                                            if (formData.customer_first_name.length >= 3 && (clientResults.length > 0 || clientLoading || clientSearched)) {
                                                setShowClientDropdown(true);
                                                setActiveSearchField('name');
                                            }
                                        }}
                                        placeholder="Buscar por nombre..."
                                        autoComplete="off"
                                        className={clientLoading && activeSearchField === 'name' ? 'pr-10' : ''}
                                    />
                                    {clientLoading && activeSearchField === 'name' && (
                                        <div className="absolute right-3 top-1/2 -translate-y-1/2">
                                            <div className="w-4 h-4 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" />
                                        </div>
                                    )}
                                </div>
                                {activeSearchField === 'name' && (
                                    <ClientAutocomplete
                                        results={clientResults}
                                        loading={clientLoading}
                                        searched={clientSearched}
                                        visible={showClientDropdown}
                                        searchTerm={formData.customer_first_name}
                                        onSelect={handleClientSelect}
                                        onClose={() => setShowClientDropdown(false)}
                                    />
                                )}
                            </div>
                            {/* Apellido — with autocomplete */}
                            <div className="relative">
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Apellido *
                                </label>
                                <div className="relative">
                                    <Input
                                        type="text"
                                        required
                                        value={formData.customer_last_name}
                                        onChange={(e) => handleClientFieldChange('lastname', e.target.value)}
                                        onFocus={() => {
                                            if (formData.customer_last_name.length >= 3 && (clientResults.length > 0 || clientLoading || clientSearched)) {
                                                setShowClientDropdown(true);
                                                setActiveSearchField('lastname');
                                            }
                                        }}
                                        placeholder="Buscar por apellido..."
                                        autoComplete="off"
                                        className={clientLoading && activeSearchField === 'lastname' ? 'pr-10' : ''}
                                    />
                                    {clientLoading && activeSearchField === 'lastname' && (
                                        <div className="absolute right-3 top-1/2 -translate-y-1/2">
                                            <div className="w-4 h-4 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" />
                                        </div>
                                    )}
                                </div>
                                {activeSearchField === 'lastname' && (
                                    <ClientAutocomplete
                                        results={clientResults}
                                        loading={clientLoading}
                                        searched={clientSearched}
                                        visible={showClientDropdown}
                                        searchTerm={formData.customer_last_name}
                                        onSelect={handleClientSelect}
                                        onClose={() => setShowClientDropdown(false)}
                                    />
                                )}
                            </div>
                            {/* Email — with autocomplete */}
                            <div className="relative">
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Email
                                </label>
                                <div className="relative">
                                    <Input
                                        type="email"
                                        value={formData.customer_email}
                                        onChange={(e) => handleClientFieldChange('email', e.target.value)}
                                        onFocus={() => {
                                            if (formData.customer_email.length >= 3 && (clientResults.length > 0 || clientLoading || clientSearched)) {
                                                setShowClientDropdown(true);
                                                setActiveSearchField('email');
                                            }
                                        }}
                                        placeholder="Buscar por email..."
                                        autoComplete="off"
                                        className={clientLoading && activeSearchField === 'email' ? 'pr-10' : ''}
                                    />
                                    {clientLoading && activeSearchField === 'email' && (
                                        <div className="absolute right-3 top-1/2 -translate-y-1/2">
                                            <div className="w-4 h-4 border-2 border-purple-500 border-t-transparent rounded-full animate-spin" />
                                        </div>
                                    )}
                                </div>
                                {activeSearchField === 'email' && (
                                    <ClientAutocomplete
                                        results={clientResults}
                                        loading={clientLoading}
                                        searched={clientSearched}
                                        visible={showClientDropdown}
                                        searchTerm={formData.customer_email}
                                        onSelect={handleClientSelect}
                                        onClose={() => setShowClientDropdown(false)}
                                    />
                                )}
                            </div>
                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Teléfono *
                                </label>
                                <div className="flex items-center w-full">
                                    <span className="px-3 py-2.5 bg-purple-700 text-white font-semibold rounded-l-lg border border-r-0 border-purple-700">
                                        +57
                                    </span>
                                    <div className="flex-1">
                                        <Input
                                            type="tel"
                                            value={formData.customer_phone}
                                            onChange={(e) => setFormData({ ...formData, customer_phone: e.target.value })}
                                            placeholder="300 1234567"
                                            className="rounded-l-none w-full"
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div>
                    <div className="bg-gray-50/80 dark:bg-gray-700/20 border border-gray-200 dark:border-gray-600/30 p-5 rounded-xl">
                        <div className="flex items-center gap-2 mb-4 pb-3 border-b border-gray-200/50 dark:border-gray-600/30">
                            <div className="w-7 h-7 rounded-lg bg-purple-100 dark:bg-purple-800/40 flex items-center justify-center text-purple-600 dark:text-purple-400 text-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
                            </div>
                            <h3 className="text-base font-semibold text-purple-700 dark:text-purple-400">Direccion de Envio</h3>
                        </div>

                        {addressAutofilled && (
                            <div className="mb-4 px-3 py-2 bg-green-50 border border-green-200 rounded-lg flex items-center justify-between">
                                <p className="text-xs text-green-700">
                                    Direccion autocompletada a partir del historial del cliente.
                                </p>
                                <button
                                    type="button"
                                    onClick={() => setAddressAutofilled(false)}
                                    className="text-green-500 hover:text-green-700 ml-2"
                                >
                                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                                </button>
                            </div>
                        )}

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="md:col-span-2">
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Dirección
                                </label>
                                <AddressAutocomplete
                                    value={formData.shipping_street}
                                    onChange={(val) => setFormData({ ...formData, shipping_street: val })}
                                    city={formData.shipping_city}
                                    onSelect={(s: AddressSuggestion) => {
                                        setAddressAutofilled(false);
                                        if (s.lat && s.lon) setAddressCoords({ lat: s.lat, lon: s.lon });
                                        if (s.neighbourhood) setBarrio(s.neighbourhood);
                                        if (s.postcode) setFormData(prev => ({ ...prev, shipping_postal_code: s.postcode }));
                                        if (s.city) {
                                            const match = daneOptions.find(
                                                (opt) => opt.ciudad.toLowerCase() === s.city.toLowerCase()
                                            ) || daneOptions.find(
                                                (opt) => opt.label.toLowerCase().includes(s.city.toLowerCase())
                                            );
                                            if (match) {
                                                handleCitySelect(match);
                                                setCitySearch(match.label);
                                            }
                                        }
                                    }}
                                />
                            </div>

                            {/* City with autocomplete */}
                            <div ref={cityRef} className="relative md:col-span-2">
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Ciudad y Departamento
                                </label>
                                <input
                                    type="text"
                                    value={citySearch}
                                    onChange={(e) => {
                                        setCitySearch(e.target.value);
                                        setShowCityResults(true);
                                    }}
                                    onFocus={() => setShowCityResults(true)}
                                    className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent text-black dark:text-white"
                                    placeholder="Buscar ciudad..."
                                />
                                {showCityResults && filteredCityOptions.length > 0 && (
                                    <div className="absolute z-10 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-300 rounded-md shadow-lg max-h-60 overflow-y-auto">
                                        {filteredCityOptions.slice(0, 50).map((opt) => (
                                            <div
                                                key={opt.value}
                                                onClick={() => handleCitySelect(opt)}
                                                className="px-3 py-2 hover:bg-purple-100 cursor-pointer text-black dark:text-white"
                                            >
                                                {opt.label}
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </div>

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Casa
                                </label>
                                <Input
                                    type="text"
                                    value={house}
                                    onChange={(e) => setHouse(e.target.value)}
                                    placeholder="Número de casa"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Barrio
                                </label>
                                <Input
                                    type="text"
                                    value={barrio}
                                    onChange={(e) => setBarrio(e.target.value)}
                                    placeholder="Nombre del barrio"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    País
                                </label>
                                <Input
                                    type="text"
                                    value={formData.shipping_country}
                                    disabled
                                    className="bg-gray-100 cursor-not-allowed text-gray-600 dark:text-gray-300"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Código Postal
                                </label>
                                <Input
                                    type="text"
                                    value={formData.shipping_postal_code}
                                    onChange={(e) => setFormData({ ...formData, shipping_postal_code: e.target.value })}
                                    placeholder="Código postal"
                                />
                            </div>

                        </div>

                        {/* Mini map preview */}
                        {addressCoords && (
                            <div className="mt-4">
                                <MapComponent
                                    address={formData.shipping_street}
                                    city={formData.shipping_city}
                                    latitude={addressCoords.lat}
                                    longitude={addressCoords.lon}
                                    height="180px"
                                />
                            </div>
                        )}
                    </div>
                </div>

                <div className="space-y-4">
                    <div className="bg-gray-50/80 dark:bg-gray-700/20 border border-gray-200 dark:border-gray-600/30 p-5 rounded-xl">
                        <div className="flex items-center gap-2 mb-4 pb-3 border-b border-gray-200/50 dark:border-gray-600/30">
                            <div className="w-7 h-7 rounded-lg bg-purple-100 dark:bg-purple-800/40 flex items-center justify-center text-purple-600 dark:text-purple-400 text-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
                            </div>
                            <h3 className="text-base font-semibold text-purple-700 dark:text-purple-400">Financiera</h3>
                        </div>
                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Total *
                                </label>
                                <div className="flex items-center w-full">
                                    <span className="px-3 py-2.5 bg-purple-700 text-white font-semibold rounded-l-lg border border-r-0 border-purple-700">
                                        $
                                    </span>
                                    <div className="flex-1 relative">
                                        <Input
                                            type="number"
                                            step="1"
                                            required
                                            value={formData.total_amount}
                                            onChange={(e) => setFormData({ ...formData, total_amount: parseFloat(e.target.value) || 0 })}
                                            className="rounded-l-none rounded-r-none w-full pr-16"
                                            placeholder="0"
                                        />
                                    </div>
                                    <span className="px-3 py-2.5 bg-purple-700 text-white font-semibold rounded-r-lg border border-l-0 border-purple-700">
                                        COP
                                    </span>
                                </div>
                            </div>

                            <div className="mt-3 pt-3 border-t border-gray-200/50 dark:border-gray-600/30">
                                <div className="flex items-center justify-between">
                                    <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Contra Entrega</span>
                                    <button
                                        type="button"
                                        role="switch"
                                        aria-checked={isCOD}
                                        onClick={() => {
                                            setIsCOD(!isCOD);
                                            setFormData(prev => ({ ...prev, cod_total: !isCOD ? prev.total_amount : 0 }));
                                        }}
                                        className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${isCOD ? 'bg-purple-600' : 'bg-gray-300 dark:bg-gray-600'}`}
                                    >
                                        <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${isCOD ? 'translate-x-6' : 'translate-x-1'}`} />
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div className="bg-gray-50/80 dark:bg-gray-700/20 border border-gray-200 dark:border-gray-600/30 p-5 rounded-xl">
                        <div className="flex items-center gap-2 mb-4 pb-3 border-b border-gray-200/50 dark:border-gray-600/30">
                            <div className="w-7 h-7 rounded-lg bg-purple-100 dark:bg-purple-800/40 flex items-center justify-center text-purple-600 dark:text-purple-400 text-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="1" y="4" width="22" height="16" rx="2" ry="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg>
                            </div>
                            <h3 className="text-base font-semibold text-purple-700 dark:text-purple-400">Pago y Estado</h3>
                        </div>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <label className="flex items-center">
                                    <input
                                        type="checkbox"
                                        checked={formData.is_paid}
                                        onChange={(e) => setFormData({ ...formData, is_paid: e.target.checked })}
                                        className="mr-2"
                                    />
                                    <span className="text-sm font-medium text-gray-700 dark:text-gray-200">Orden Pagada</span>
                                </label>
                            </div>
                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Estado
                                </label>
                                <select
                                    value={formData.status}
                                    onChange={(e) => setFormData({ ...formData, status: e.target.value })}
                                    className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-600 focus:border-transparent bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-100"
                                >
                                    <option value="pending">Pendiente</option>
                                    <option value="processing">Procesando</option>
                                    <option value="shipped">Enviado</option>
                                    <option value="delivered">Entregado</option>
                                    <option value="cancelled">Cancelado</option>
                                </select>
                            </div>
                            <div>
                                <label className="flex items-center">
                                    <input
                                        type="checkbox"
                                        checked={formData.invoiceable}
                                        onChange={(e) => setFormData({ ...formData, invoiceable: e.target.checked })}
                                        className="mr-2"
                                    />
                                    <span className="text-sm font-medium text-gray-700 dark:text-gray-200">Facturable</span>
                                </label>
                            </div>
                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Confirmación
                                </label>
                                <select
                                    value={formData.is_confirmed === true ? 'yes' : formData.is_confirmed === false ? 'no' : 'pending'}
                                    onChange={(e) => {
                                        const v = e.target.value;
                                        setFormData({ ...formData, is_confirmed: v === 'yes' ? true : v === 'no' ? false : null });
                                    }}
                                    className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-600 focus:border-transparent bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-100"
                                >
                                    <option value="pending">Pendiente</option>
                                    <option value="yes">Confirmado</option>
                                    <option value="no">No confirmado</option>
                                </select>
                            </div>
                        </div>
                    </div>

                </div>

                <div className="lg:col-span-2">
                    <div className="bg-gray-50/80 dark:bg-gray-700/20 border border-gray-200 dark:border-gray-600/30 p-5 rounded-xl">
                        <div className="flex items-center gap-2 mb-4 pb-3 border-b border-gray-200/50 dark:border-gray-600/30">
                            <div className="w-7 h-7 rounded-lg bg-gray-200 dark:bg-gray-600 flex items-center justify-center text-gray-600 dark:text-gray-300 text-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
                            </div>
                            <h3 className="text-base font-semibold text-gray-700 dark:text-gray-200">Productos</h3>
                        </div>
                        <ProductSelector
                            businessId={formData.business_id || 0}
                            selectedProducts={selectedProducts}
                            onSelect={handleProductsChange}
                            onCreateNew={() => setShowProductModal(true)}
                        />
                    </div>
                </div>

                <div>
                    <div className="bg-gray-50/80 dark:bg-gray-700/20 border border-gray-200 dark:border-gray-600/30 p-5 rounded-xl h-full">
                        <div className="flex items-center gap-2 mb-4 pb-3 border-b border-gray-200/50 dark:border-gray-600/30">
                            <div className="w-7 h-7 rounded-lg bg-gray-200 dark:bg-gray-600 flex items-center justify-center text-gray-600 dark:text-gray-300 text-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
                            </div>
                            <h3 className="text-base font-semibold text-gray-700 dark:text-gray-200">Notas</h3>
                        </div>
                        <textarea
                            value={formData.notes}
                            onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                            rows={5}
                            placeholder="Notas internas sobre la orden..."
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-600 focus:border-transparent resize-none text-sm"
                        />
                    </div>
                </div>

                <div className="lg:col-span-3">
                    <div className="bg-gray-50/80 dark:bg-gray-700/20 border border-gray-200 dark:border-gray-600/30 p-5 rounded-xl">
                        <div className="flex items-center gap-2 mb-4 pb-3 border-b border-gray-200/50 dark:border-gray-600/30">
                            <div className="w-7 h-7 rounded-lg bg-purple-100 dark:bg-purple-800/40 flex items-center justify-center text-purple-600 dark:text-purple-400 text-sm">
                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="1" y="3" width="15" height="13"/><polygon points="16 8 20 8 23 11 23 16 16 16 16 8"/><circle cx="5.5" cy="18.5" r="2.5"/><circle cx="18.5" cy="18.5" r="2.5"/></svg>
                            </div>
                            <h3 className="text-base font-semibold text-purple-700 dark:text-purple-400">Logistica</h3>
                        </div>
                        
                        {/* 24-hour processing notice */}
                        <div className="mb-4 bg-blue-50 border-l-4 border-blue-400 p-3 rounded-r-lg">
                            <p className="text-[11px] leading-tight text-blue-800">
                                <span className="font-bold uppercase tracking-wider block mb-1">ℹ️ Aviso de Procesamiento</span>
                                Recuerda que las transportadoras pueden demorar hasta <span className="font-bold">24 horas hábiles</span> en procesar y recolectar los pedidos después de generada la guía.
                            </p>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Número de Guía
                                </label>
                                <Input
                                    type="text"
                                    value={formData.tracking_number}
                                    onChange={(e) => setFormData({ ...formData, tracking_number: e.target.value })}
                                    placeholder="ej: 123456789"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    ID Guía
                                </label>
                                <Input
                                    type="text"
                                    value={formData.guide_id}
                                    onChange={(e) => setFormData({ ...formData, guide_id: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Bodega
                                </label>
                                {warehouses.length > 0 ? (
                                    <select
                                        value={formData.warehouse_id || ''}
                                        onChange={(e) => {
                                            const id = parseInt(e.target.value);
                                            const w = warehouses.find(w => w.id === id);
                                            setFormData({ ...formData, warehouse_id: id || 0, warehouse_name: w?.name || '' });
                                        }}
                                        className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-600 focus:border-transparent bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-100"
                                    >
                                        <option value="">Seleccionar bodega...</option>
                                        {warehouses.map(w => (
                                            <option key={w.id} value={w.id}>
                                                {w.name} ({w.code}){w.is_default ? ' - Default' : ''}
                                            </option>
                                        ))}
                                    </select>
                                ) : (
                                    <Input
                                        type="text"
                                        value={formData.warehouse_name}
                                        onChange={(e) => setFormData({ ...formData, warehouse_name: e.target.value })}
                                        placeholder={warehousesLoading ? 'Cargando bodegas...' : 'Nombre de bodega'}
                                    />
                                )}
                            </div>
                            <div>
                                <label className="block text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">
                                    Transportadora
                                </label>
                                <Input
                                    type="text"
                                    value={cachedGuideCarrier || (order as any)?.shipment?.carrier || 'Sin asignar'}
                                    disabled
                                    className="bg-gray-100 cursor-not-allowed"
                                />
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Actions */}
            <div className="flex justify-end space-x-3 pt-6 border-t border-gray-200 dark:border-gray-700">
                {onCancel && (
                    <button
                        type="button"
                        onClick={onCancel}
                        className="px-6 py-3 text-base text-white font-semibold rounded-lg bg-purple-700 hover:bg-purple-800 shadow-sm hover:shadow-md transition-all"
                    >
                        Cancelar
                    </button>
                )}
                <button
                    type="submit"
                    disabled={loading}
                    className="px-6 py-3 text-base text-white font-semibold rounded-lg bg-purple-700 hover:bg-purple-800 shadow-sm hover:shadow-md transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                    {loading && <div className="spinner w-4 h-4" />}
                    {isEdit ? 'Actualizar Orden' : 'Crear Orden'}
                </button>
            </div>
            <Modal
                isOpen={showProductModal}
                onClose={() => setShowProductModal(false)}
                title="Crear Nuevo Producto"
            >
                <ProductForm
                    onSuccess={() => {
                        setShowProductModal(false);
                        // Optional: Refresh products list if needed, but ProductSelector handles it
                    }}
                    onCancel={() => setShowProductModal(false)}
                />
            </Modal>
        </form>
    );
}
