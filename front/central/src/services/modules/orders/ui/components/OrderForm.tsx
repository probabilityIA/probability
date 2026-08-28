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
import { resolveCityState } from '@/shared/utils/dane-lookup';
import { useClientSearch } from '../hooks/useClientSearch';
import { useWarehouses } from '../hooks/useWarehouses';
import { useDynamicBusinessColors } from '../hooks/useDynamicBusinessColors';
import ClientAutocomplete from './ClientAutocomplete';
import AddressAutocomplete, { AddressSuggestion } from './AddressAutocomplete';
import PaymentMethodSelect from '../../../paymentmethods/ui/components/PaymentMethodSelect';
import dynamic from 'next/dynamic';

const MapComponent = dynamic(() => import('@/shared/ui/MapComponent'), { ssr: false });
import { CustomerInfo } from '../../../customers/domain/types';
import { getCustomerAddressesAction } from '../../../customers/infra/actions';
import { getEffectivePriceAction, listClientGroupsAction, getCatalogPricesAction, listAvailableClientsAction } from '../../../pricing/infra/actions';
import { ClientGroup } from '../../../pricing/domain/types';
import { getActionError } from '@/shared/utils/action-result';
import { salesChannelLabel } from '../../domain/sales-channel-labels';

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
    const businessName = (permissions as any)?.business_name || 'Negocio';
    const { showToast } = useToast();
    const [cachedGuideCarrier, setCachedGuideCarrier] = useState<string | null>(null);
    const { colors: businessColors } = useDynamicBusinessColors(defaultBusinessId);

    const primaryColor = businessColors?.primary_color || '#5b21b6';
    const secondaryColor = businessColors?.secondary_color || '#7c3aed';
    const tertiaryColor = businessColors?.tertiary_color || '#c4b5fd';
    const quaternaryColor = businessColors?.quaternary_color || '#ede9fe';

    useEffect(() => {
        if (isEdit && order?.order_number) {
            const key = `guide_${order.order_number}`;
            const cached = sessionStorage.getItem(key);
            if (cached) {
                try {
                    const guideData = JSON.parse(cached);
                    if (guideData.carrier) {
                        setCachedGuideCarrier(guideData.carrier);
                    }
                } catch {
                    setCachedGuideCarrier(null);
                }
            }
        }
    }, [isEdit, order?.order_number]);

    const [formData, setFormData] = useState({
        integration_id: order?.integration_id || 0,
        platform: order?.platform || 'manual',
        business_id: order?.business_id || defaultBusinessId,

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
        is_cod: order?.is_cod ?? ((order?.cod_total || 0) > 0),

        payment_method_id: order?.payment_method_id || 0,
        is_paid: order?.is_paid || false,

        status: order?.status || 'pending',

        tracking_number: order?.tracking_number || '',
        tracking_link: order?.tracking_link || '',
        guide_id: order?.guide_id || '',
        warehouse_id: order?.warehouse_id || 0,
        warehouse_name: order?.warehouse_name || '',
        driver_name: order?.driver_name || '',
        is_last_mile: order?.is_last_mile || false,

        notes: order?.notes || '',
        invoiceable: order?.invoiceable ?? false,
        is_confirmed: order?.is_confirmed ?? null,
        novelty: order?.novelty || '',

        items: order?.order_items || order?.items || [],

        integration_type: order?.integration_type || '',
        external_id: order?.external_id || '',
    });

    const [selectedProducts, setSelectedProducts] = useState<Product[]>(() => {
        const items = order?.order_items ?? order?.items;
        if (!items || !Array.isArray(items)) return [];
        return (items as any[])
            .map((item: any) => ({
                ...item,
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
    const [selectedClientId, setSelectedClientId] = useState<number | null>(null);
    const [isCOD, setIsCOD] = useState(() => order?.is_cod ?? ((order?.cod_total || 0) > 0));
    const [paymentMethodError, setPaymentMethodError] = useState(false);
    const [groups, setGroups] = useState<ClientGroup[]>([]);
    const [selectedGroupId, setSelectedGroupId] = useState<number | null>(null);
    const [clientGroupName, setClientGroupName] = useState('');

    const applyProducts = (products: Product[]) => {
        setSelectedProducts(products);
        const subtotal = products.reduce((acc, p) => acc + (p.price * (p.quantity || 1)), 0);
        setFormData(prev => {
            const orderValue = subtotal + prev.tax - prev.discount;
            return {
                ...prev,
                items: products,
                subtotal,
                total_amount: orderValue,
                cod_total: isCOD ? orderValue + prev.shipping_cost : prev.cod_total,
            };
        });
    };

    const repriceProducts = async (products: Product[], clientId: number): Promise<Product[]> => {
        return Promise.all(products.map(async (p) => {
            if ((p as any)._personalized) return p;
            const res = await getEffectivePriceAction(formData.business_id, p.id, clientId);
            if (res && typeof res.final_price === 'number') {
                return {
                    ...p,
                    price: res.final_price,
                    _personalized: true,
                    _basePrice: res.base_price,
                    _priceSource: res.source,
                    _groupName: res.group_name,
                } as Product;
            }
            return { ...p, _personalized: true } as Product;
        }));
    };

    const repriceByGroup = async (products: Product[], groupId: number | null): Promise<Product[]> => {
        if (!groupId) {
            return products.map(p => {
                const base = (p as any)._basePrice ?? p.price;
                return { ...p, price: base, _basePrice: base, _priceSource: 'base', _groupName: '' } as Product;
            });
        }
        const result = await getCatalogPricesAction(formData.business_id, { client_group_id: groupId }, '', 1);
        const byId = new Map(result.data.map(r => [r.product_id, r]));
        const groupName = groups.find(g => g.id === groupId)?.name || '';
        return products.map(p => {
            const row = byId.get(p.id);
            if (row) {
                const final = row.custom_price ?? row.base_price;
                return {
                    ...p,
                    price: final,
                    _basePrice: row.base_price,
                    _priceSource: row.custom_price != null ? 'group' : 'base',
                    _groupName: groupName,
                } as Product;
            }
            const base = (p as any)._basePrice ?? p.price;
            return { ...p, price: base, _basePrice: base, _priceSource: 'base', _groupName: groupName } as Product;
        });
    };

    const pricingInfo = (() => {
        const withGroup = selectedProducts.find(p => (p as any)._groupName);
        const groupName = withGroup ? String((withGroup as any)._groupName) : '';
        const anyCustom = selectedProducts.some(p => {
            const src = (p as any)._priceSource;
            return src && src !== 'base';
        });
        return { groupName, anyCustom };
    })();

    useEffect(() => {
        if (order?.order_items && Array.isArray(order.order_items)) {
            const mapped = (order.order_items as any[])
                .map((item: any) => ({
                    ...item,
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
            setSelectedProducts(mapped);
        }
    }, [order?.id, order?.updated_at, order?.order_items]);

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [citySearch, setCitySearch] = useState('');
    const [showCityResults, setShowCityResults] = useState(false);
    const [citySelected, setCitySelected] = useState(false);
    const [cityError, setCityError] = useState(false);
    const cityRef = useRef<HTMLDivElement>(null);

    const normalizeCity = (s: string) =>
        s.normalize('NFD').replace(/[\u0300-\u036f]/g, '')
            .toLowerCase()
            .replace(/\s*[,(]?\s*d\.?\s*c\.?\s*\)?\s*$/g, '')
            .trim();

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

        setSelectedClientId(client.id);
        setSelectedGroupId(null);
        if (selectedProducts.length > 0) {
            const cleared = selectedProducts.map(p => ({ ...p, _personalized: false } as Product));
            const repriced = await repriceProducts(cleared, client.id);
            applyProducts(repriced);
            const changed = repriced.some((p, i) => p.price !== selectedProducts[i]?.price);
            if (changed) {
                showToast('Se aplicaron los precios personalizados de este cliente', 'success');
            }
        }

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
        }
    }, [clearClients, formData.business_id, selectedProducts, showToast]);

    const daneOptions = Object.entries(danes).map(([code, data]: [string, any]) => ({
        value: code,
        label: `${data.ciudad} (${data.departamento})`,
        ciudad: data.ciudad,
        departamento: data.departamento
    })).sort((a, b) => a.label.localeCompare(b.label));

    const filteredCityOptions = daneOptions.filter(opt =>
        opt.label.toLowerCase().includes(citySearch.toLowerCase())
    );

    useEffect(() => {
        if (warehouses.length === 0 || formData.warehouse_id) return;
        const defaultW = warehouses.find(w => w.is_default) || (warehouses.length === 1 ? warehouses[0] : null);
        if (defaultW) {
            setFormData(prev => ({ ...prev, warehouse_id: defaultW.id, warehouse_name: defaultW.name }));
        }
    }, [warehouses]);

    useEffect(() => {
        if (!formData.business_id) return;
        (async () => {
            const res = await listClientGroupsAction(formData.business_id, '', 1);
            setGroups(res.data.filter(g => g.is_active));
        })();
    }, [formData.business_id]);

    useEffect(() => {
        if (!isEdit || !formData.business_id) return;
        const term = order?.customer_dni || order?.customer_email || '';
        if (!term) return;
        (async () => {
            const res = await listAvailableClientsAction(formData.business_id, term, false, 1);
            const match = res.data.find(c =>
                (order?.customer_dni && c.dni === order.customer_dni) ||
                (order?.customer_email && c.email === order.customer_email)
            ) || res.data[0];
            if (match && match.group_name) {
                setClientGroupName(match.group_name);
            }
        })();
    }, [isEdit, formData.business_id]);

    useEffect(() => {
        if (!order?.shipping_city && !order?.shipping_state) return;

        const resolved = resolveCityState(order.shipping_city || '', order.shipping_state || '', order.destination_dane_code);
        setCitySearch(resolved
            ? `${resolved.ciudad} (${resolved.departamento})`
            : `${order.shipping_city} (${order.shipping_state})`);
        setCitySelected(true);
    }, [order]);

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
        setCitySelected(true);
        setCityError(false);
    };

    const handleCityBlur = () => {
        setTimeout(() => {
            setShowCityResults(false);
            if (!citySelected && citySearch.trim() !== '') {
                const searchNorm = normalizeCity(citySearch.trim());
                const exact = daneOptions.find(o =>
                    o.label.toLowerCase() === citySearch.trim().toLowerCase() ||
                    normalizeCity(`${o.ciudad} (${o.departamento})`) === searchNorm
                );
                if (exact) {
                    handleCitySelect(exact);
                } else {
                    setCityError(true);
                }
            }
        }, 150);
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setPaymentMethodError(false);

        try {
            if ((!formData.customer_name && !formData.customer_first_name) || !formData.total_amount) {
                throw new Error('Por favor completa los campos requeridos');
            }
            if (!formData.payment_method_id || formData.payment_method_id <= 0) {
                setPaymentMethodError(true);
                throw new Error('Selecciona el medio de pago (con que paga el cliente)');
            }
            if (!formData.shipping_street?.trim() || !formData.shipping_city?.trim() || !formData.shipping_state?.trim()) {
                setCityError(!formData.shipping_city?.trim());
                throw new Error('La direccion de envio (calle, ciudad y departamento) es obligatoria');
            }

            const parts = [formData.shipping_street || ''];
            if (house.trim()) parts.push(house.trim());
            if (barrio.trim()) parts.push(barrio.trim());
            const fullShippingStreet = parts.join(' | ');

            const itemsToSend = selectedProducts.length > 0 ? selectedProducts : formData.items;

            const baseData = {
                ...formData,
                is_cod: isCOD,
                // formData.cod_total ya viene correcto desde la orden (o desde
                // applyProducts al agregar productos); no se recalcula aqui para
                // no perder comisiones/ajustes que no son solo producto+envio
                // (p.ej. ordenes de WooCommerce donde el canal ya cobro el total).
                cod_total: isCOD ? formData.cod_total : 0,
                payment_method_id: formData.payment_method_id,
                shipping_street: fullShippingStreet,
                shipping_lat: addressCoords?.lat,
                shipping_lng: addressCoords?.lon,
                items: itemsToSend,
                customer_name: formData.customer_name || `${formData.customer_first_name} ${formData.customer_last_name}`.trim(),
                client_group_id: !isEdit && selectedGroupId ? selectedGroupId : undefined,
            };

            let response;
            if (isEdit && order) {
                const updateData: UpdateOrderDTO = {
                    ...baseData,
                    confirmation_status: formData.is_confirmed === true ? 'yes' : formData.is_confirmed === false ? 'no' : 'pending',
                    is_confirmed: undefined,
                };
                response = await updateOrderAction(order.id, updateData);
            } else {
                response = await createOrderAction(baseData as CreateOrderDTO);
            }

            if (response.success) {
                showToast(isEdit ? 'Orden actualizada exitosamente' : 'Orden creada exitosamente', 'success');
                if (onSuccess) {
                    onSuccess();
                }
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
        const filteredProducts = products.filter(p => (p.quantity || 0) > 0);
        if (selectedClientId) {
            repriceProducts(filteredProducts, selectedClientId).then(applyProducts);
        } else if (selectedGroupId) {
            repriceByGroup(filteredProducts, selectedGroupId).then(applyProducts);
        } else {
            applyProducts(filteredProducts);
        }
    };

    const handleSelectGroup = async (groupId: number | null) => {
        setSelectedGroupId(groupId);
        if (selectedProducts.length > 0) {
            const repriced = await repriceByGroup(selectedProducts, groupId);
            applyProducts(repriced);
        }
    };

    const calculateTotal = () => {
        setFormData(prev => ({
            ...prev,
            total_amount: prev.subtotal + prev.tax - prev.discount,
        }));
    };

    // Igual que "se cobrara contra entrega": si cod_total ya trae un valor distinto
    // a producto+envio (p.ej. la comision completa que cobro WooCommerce en el
    // checkout), se usa como el total real de la orden en vez de recalcularlo.
    const yaLiquidadaEnCorte = order?.cod_cut_confirmed === true;
    const displayTotal = isCOD && !yaLiquidadaEnCorte && formData.cod_total > 0
        ? formData.cod_total
        : formData.total_amount + formData.shipping_cost;

    return (
        <form onSubmit={handleSubmit} className="flex h-full flex-col bg-slate-100 dark:bg-gray-900" style={{ fontFamily: "'Plus Jakarta Sans', sans-serif" }}>
            {isEdit && order && (
                <header className="flex shrink-0 flex-wrap items-center gap-2.5 border-b border-slate-200 bg-white px-4 py-2.5 dark:border-gray-700 dark:bg-gray-800">
                    <span className="flex h-8 w-8 items-center justify-center rounded-lg text-white" style={{ backgroundColor: primaryColor }}>
                        <svg className="h-4 w-4" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                            <circle cx="8" cy="5" r="3" />
                            <path d="M2 14c0-3 2.5-5 6-5s6 2 6 5" />
                        </svg>
                    </span>
                    <h2 className="text-base font-bold tracking-tight text-slate-900 dark:text-white">Editar Orden</h2>
                    <span className="rounded-full bg-slate-100 px-2.5 py-0.5 font-mono text-xs font-semibold text-slate-600 dark:bg-gray-700 dark:text-gray-200">
                        #{order.order_number || order.internal_number || order.id}
                    </span>
                    {businessName && (
                        <span className="rounded-full border border-slate-200 px-2.5 py-0.5 text-[11px] font-semibold uppercase text-slate-500 dark:border-gray-600 dark:text-gray-300">
                            {businessName}
                        </span>
                    )}
                    {isCOD && (
                        <span className="inline-flex items-center gap-1 rounded-full border border-amber-300 bg-amber-100 px-2.5 py-0.5 text-[11px] font-bold uppercase text-amber-700">
                            <svg className="h-3 w-3" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M3 10h18M7 15h4m-7 4h16a2 2 0 002-2V7a2 2 0 00-2-2H4a2 2 0 00-2 2v10a2 2 0 002 2z" />
                            </svg>
                            Contra entrega
                        </span>
                    )}
                    <span className="ml-auto flex items-center gap-1.5">
                        {order.integration_logo_url && (
                            <img src={order.integration_logo_url} alt={order.platform} className="h-5 w-5 object-contain" />
                        )}
                        <span className="text-xs font-semibold text-slate-600 dark:text-gray-300">{salesChannelLabel(order.platform)}</span>
                    </span>
                    {onCancel && (
                        <button
                            type="button"
                            onClick={onCancel}
                            title="Cerrar"
                            className="rounded-full p-1.5 text-slate-400 transition hover:bg-slate-100 hover:text-slate-700 dark:hover:bg-gray-700"
                        >
                            <svg className="h-5 w-5" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    )}
                </header>
            )}

            <div className="flex-1 overflow-y-auto px-4 py-4">
                {error && (
                    <Alert type="error" onClose={() => setError(null)}>
                        {error}
                    </Alert>
                )}

                <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
                <div>
                    <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                        <div className="mb-3 flex items-center gap-2 border-b border-slate-100 pb-2.5 dark:border-gray-700">
                            <div className="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-lg text-white" style={{ backgroundColor: primaryColor }}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="white" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                                    <circle cx="8" cy="5" r="3"/>
                                    <path d="M2 14c0-3 2.5-5 6-5s6 2 6 5"/>
                                </svg>
                            </div>
                            <h3 className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">Cliente</h3>
                        </div>

                        {!formData.customer_email && (formData.customer_first_name || formData.customer_name) && (
                            <div className="mb-4 px-3 py-2 bg-amber-50 border border-amber-200 rounded-lg">
                                <p className="text-xs text-amber-700">
                                    Completar el email del cliente mejora la analitica, seguimiento y facturacion de tus ordenes.
                                </p>
                            </div>
                        )}

                        <div className="space-y-4">
                            <div className="relative">
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    DNI / Cedula
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
                                        placeholder="Buscar cliente por cedula..."
                                        autoComplete="off"
                                        className={`${clientLoading && activeSearchField === 'dni' ? 'pr-10' : ''}`}
                                        style={{ borderColor: '#e8e0f5', height: '38px' }}
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
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                <div className="relative">
                                    <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
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
                                <div className="relative">
                                    <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
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
                            </div>
                            <div className="relative">
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
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
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    Telefono *
                                </label>
                                <div className="flex items-center w-full">
                                    <span className="px-3 py-2.5 text-white font-semibold rounded-l-lg border border-r-0" style={{ backgroundColor: primaryColor, borderColor: primaryColor }}>
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
                    <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                        <div className="mb-3 flex items-center gap-2 border-b border-slate-100 pb-2.5 dark:border-gray-700">
                            <div className="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-lg text-white" style={{ backgroundColor: primaryColor }}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="white" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                                    <path d="M8 14s-5-4.5-5-8a5 5 0 0 1 10 0c0 3.5-5 8-5 8z"/>
                                    <circle cx="8" cy="6" r="1.5"/>
                                </svg>
                            </div>
                            <h3 className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">Direccion de Envio</h3>
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
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    Direccion
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
                                            const cityKey = normalizeCity(s.city);
                                            const stateKey = s.state ? normalizeCity(s.state) : '';
                                            const match =
                                                daneOptions.find(opt =>
                                                    normalizeCity(opt.ciudad) === cityKey &&
                                                    (!stateKey || normalizeCity(opt.departamento) === stateKey)
                                                ) ||
                                                daneOptions.find(opt => normalizeCity(opt.ciudad) === cityKey) ||
                                                daneOptions.find(opt => normalizeCity(opt.label).includes(cityKey));
                                            if (match) {
                                                handleCitySelect(match);
                                            }
                                        }
                                    }}
                                />
                            </div>

                            <div ref={cityRef} className="relative md:col-span-2">
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    Ciudad y Departamento
                                </label>
                                <input
                                    type="text"
                                    value={citySearch}
                                    onChange={(e) => {
                                        setCitySearch(e.target.value);
                                        setShowCityResults(true);
                                        setCitySelected(false);
                                        setCityError(false);
                                        setFormData(prev => ({ ...prev, shipping_city: '', shipping_state: '' }));
                                    }}
                                    onFocus={() => setShowCityResults(true)}
                                    onBlur={handleCityBlur}
                                    className={`w-full px-3 py-2 bg-white dark:bg-gray-800 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent text-black dark:text-white ${
                                        cityError ? 'border-red-500' : citySelected ? 'border-green-400' : 'border-gray-300'
                                    }`}
                                    placeholder="Buscar ciudad... (selecciona una opcion)"
                                />
                                {cityError && (
                                    <p className="mt-1 text-xs text-red-600">Selecciona una opcion del listado</p>
                                )}
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
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    Casa
                                </label>
                                <Input
                                    type="text"
                                    value={house}
                                    onChange={(e) => setHouse(e.target.value)}
                                    placeholder="Numero de casa"
                                />
                            </div>

                            <div>
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
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
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    Pais
                                </label>
                                <Input
                                    type="text"
                                    value={formData.shipping_country}
                                    disabled
                                    className="bg-gray-100 cursor-not-allowed text-gray-600 dark:text-gray-300"
                                />
                            </div>
                            <div>
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    Codigo Postal
                                </label>
                                <Input
                                    type="text"
                                    value={formData.shipping_postal_code}
                                    onChange={(e) => setFormData({ ...formData, shipping_postal_code: e.target.value })}
                                    placeholder="Codigo postal"
                                />
                            </div>

                        </div>

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
                    <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                        <div className="mb-3 flex items-center gap-2 border-b border-slate-100 pb-2.5 dark:border-gray-700">
                            <div className="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-lg text-sm font-bold text-white" style={{ backgroundColor: primaryColor }}>
                                $
                            </div>
                            <h3 className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">Financiera</h3>
                        </div>
                        <div className="space-y-4">
                            <div>
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    Valor productos (sin envio) *
                                </label>
                                <div className="relative w-full">
                                    <span className="pointer-events-none absolute left-3 top-1/2 z-10 -translate-y-1/2 text-sm font-medium text-slate-400">$</span>
                                    <div className="flex-1 relative">
                                        <Input
                                            type="number"
                                            step="1"
                                            required
                                            value={formData.total_amount}
                                            onChange={(e) => setFormData({ ...formData, total_amount: parseFloat(e.target.value) || 0 })}
                                            className="w-full tabular-nums" style={{ paddingLeft: '28px' }}
                                            placeholder="0"
                                        />
                                    </div>
                                </div>
                            </div>

                            <div>
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    Costo de Envio
                                </label>
                                <div className="relative w-full">
                                    <span className="pointer-events-none absolute left-3 top-1/2 z-10 -translate-y-1/2 text-sm font-medium text-slate-400">$</span>
                                    <div className="flex-1">
                                        <Input
                                            type="number"
                                            step="1"
                                            value={formData.shipping_cost}
                                            onChange={(e) => {
                                                const sc = parseFloat(e.target.value) || 0;
                                                setFormData(prev => ({
                                                    ...prev,
                                                    shipping_cost: sc,
                                                    cod_total: isCOD ? prev.total_amount + sc : prev.cod_total,
                                                }));
                                            }}
                                            className="w-full tabular-nums" style={{ paddingLeft: '28px' }}
                                            placeholder="0"
                                        />
                                    </div>
                                </div>
                                <p className="mt-1.5 text-[11px] leading-snug text-slate-400">
                                    Se actualiza con el costo de la guia al generarla, o ingresalo manualmente.
                                </p>
                                {(() => {
                                    if (!isCOD || yaLiquidadaEnCorte || formData.cod_total <= 0) return null;
                                    const comision = formData.cod_total - formData.total_amount - formData.shipping_cost;
                                    if (comision <= 0) return null;
                                    return (
                                        <p className="mt-1 text-[11px] leading-snug text-slate-400">
                                            + comision contra entrega: <strong>{formData.currency} {comision.toLocaleString('es-CO')}</strong>
                                        </p>
                                    );
                                })()}
                            </div>

                            <div className="mt-3 border-t border-slate-100 pt-3 dark:border-gray-700">
                                <div className="flex items-start justify-between gap-3">
                                    <div>
                                        <p className="text-sm font-semibold text-slate-700 dark:text-gray-200">Contra Entrega</p>
                                        <p className="text-[11px] leading-snug text-slate-400">Cuando se cobra: el dinero se recauda al entregar.</p>
                                    </div>
                                    <button
                                        type="button"
                                        role="switch"
                                        aria-checked={isCOD}
                                        onClick={() => {
                                            const next = !isCOD;
                                            setIsCOD(next);
                                            setFormData(prev => ({
                                                ...prev,
                                                is_cod: next,
                                                cod_total: next ? prev.total_amount + prev.shipping_cost : 0,
                                            }));
                                        }}
                                        className="relative inline-flex h-6 w-11 items-center rounded-full transition-colors"
                                        style={{ backgroundColor: isCOD ? primaryColor : '#d1d5db' }}
                                    >
                                        <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${isCOD ? 'translate-x-6' : 'translate-x-1'}`} />
                                    </button>
                                </div>
                                {isCOD && (() => {
                                    const sumaProductoEnvio = formData.total_amount + formData.shipping_cost;
                                    const yaLiquidadaEnCorte = order?.cod_cut_confirmed === true;
                                    const montoCOD = !yaLiquidadaEnCorte && formData.cod_total > 0 ? formData.cod_total : sumaProductoEnvio;
                                    const liquidado = Math.abs(montoCOD - sumaProductoEnvio) >= 1;
                                    return (
                                        <p className="mt-2 text-xs text-gray-600 dark:text-gray-300">
                                            Se cobrara contra entrega: <strong>{formData.currency} {montoCOD.toLocaleString()}</strong>
                                            <span className="text-gray-400">
                                                {liquidado ? ' (valor liquidado por la transportadora)' : ' (producto + envio)'}
                                            </span>
                                        </p>
                                    );
                                })()}
                            </div>

                            <div className="flex items-center justify-between rounded-xl px-4 py-3" style={{ backgroundColor: primaryColor }}>
                                <span className="flex items-center gap-2 text-[11px] font-semibold uppercase tracking-wide text-white/60">Total con envio</span>
                                <span className="text-lg font-bold tabular-nums text-white">
                                    {formData.currency} {displayTotal.toLocaleString('es-CO')}
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="col-span-3">
                    <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                        <div className="mb-3 flex items-center gap-2 border-b border-slate-100 pb-2.5 dark:border-gray-700">
                            <div className="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-lg text-white" style={{ backgroundColor: primaryColor }}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="white" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                                    <rect x="2" y="4" width="12" height="9" rx="2"/>
                                    <path d="M2 8h12M5 11h2M9 11h2"/>
                                </svg>
                            </div>
                            <h3 className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">Pago y Estado</h3>
                        </div>
                        <div className="grid grid-cols-1 md:grid-cols-6 gap-4 items-start">
                            <div className="md:col-span-2">
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    Medio de pago *
                                </label>
                                <PaymentMethodSelect
                                    value={formData.payment_method_id}
                                    onChange={(id) => {
                                        setFormData(prev => ({ ...prev, payment_method_id: id }));
                                        setPaymentMethodError(false);
                                    }}
                                    hasError={paymentMethodError}
                                />
                                <p className="mt-1 text-[11px] text-gray-400">
                                    Con que paga el cliente. Contra entrega define <strong>cuando</strong> se cobra; el medio de pago, <strong>con que</strong> se paga.
                                </p>
                                {paymentMethodError && (
                                    <p className="mt-1 text-xs text-red-600">Selecciona el medio de pago</p>
                                )}
                            </div>
                            <div>
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
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
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    Confirmacion
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
                            <div className="md:col-span-2 flex flex-wrap items-center gap-6 md:pt-7">
                                <label className="flex items-center">
                                    <div className="relative flex items-center justify-center">
                                        <input
                                            type="checkbox"
                                            checked={formData.is_paid}
                                            onChange={(e) => setFormData({ ...formData, is_paid: e.target.checked })}
                                            className="appearance-none w-5 h-5 border-2 rounded cursor-pointer checked:bg-[var(--primary-color)] checked:border-[var(--primary-color)]" style={{ borderColor: tertiaryColor, '--primary-color': primaryColor } as any}
                                        />
                                        {formData.is_paid && (
                                            <svg className="absolute w-3 h-3 text-white pointer-events-none" fill="currentColor" viewBox="0 0 20 20">
                                                <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                            </svg>
                                        )}
                                    </div>
                                    <span className="text-sm font-medium text-gray-700 dark:text-gray-200 ml-2">Orden Pagada</span>
                                </label>
                                <label className="flex items-center">
                                    <div className="relative flex items-center justify-center">
                                        <input
                                            type="checkbox"
                                            checked={formData.invoiceable}
                                            onChange={(e) => setFormData({ ...formData, invoiceable: e.target.checked })}
                                            className="appearance-none w-5 h-5 border-2 rounded cursor-pointer checked:bg-[var(--primary-color)] checked:border-[var(--primary-color)]" style={{ borderColor: tertiaryColor, '--primary-color': primaryColor } as any}
                                        />
                                        {formData.invoiceable && (
                                            <svg className="absolute w-3 h-3 text-white pointer-events-none" fill="currentColor" viewBox="0 0 20 20">
                                                <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                            </svg>
                                        )}
                                    </div>
                                    <span className="text-sm font-medium text-gray-700 dark:text-gray-200 ml-2">Facturable</span>
                                </label>
                            </div>
                        </div>
                    </div>
                </div>


                <div className="lg:col-span-2">
                    <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                        <div className="mb-3 flex items-center gap-2 border-b border-slate-100 pb-2.5 dark:border-gray-700">
                            <div className="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-lg text-white" style={{ backgroundColor: primaryColor }}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="white" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round">
                                    <path d="M2 6l6-4 6 4v8l-6 4-6-4V6z"/>
                                    <path d="M8 2v4M2 6l6 4 6-4"/>
                                </svg>
                            </div>
                            <h3 className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">Productos</h3>
                        </div>
                        {selectedClientId && (pricingInfo.groupName || pricingInfo.anyCustom) && (
                            <div className="mb-4 px-3 py-2 bg-indigo-50 border border-indigo-200 rounded-lg">
                                <p className="text-xs text-indigo-800">
                                    {pricingInfo.groupName ? (
                                        <>
                                            Este cliente pertenece al grupo <strong>{pricingInfo.groupName}</strong>.
                                            Los precios mostrados son los personalizados de su grupo; la diferencia frente al precio base se indica en cada producto.
                                        </>
                                    ) : (
                                        <>Este cliente tiene precios especiales asignados. La diferencia frente al precio base se indica en cada producto.</>
                                    )}
                                </p>
                            </div>
                        )}
                        {isEdit && clientGroupName && (
                            <div className="mb-4 px-3 py-2 bg-indigo-50 border border-indigo-200 rounded-lg">
                                <p className="text-xs text-indigo-800">
                                    Este cliente pertenece al grupo <strong>{clientGroupName}</strong> y tiene precios personalizados por su grupo.
                                </p>
                            </div>
                        )}
                        {!isEdit && !selectedClientId && groups.length > 0 && (
                            <div className="mb-4 px-3 py-3 bg-indigo-50 border border-indigo-200 rounded-lg">
                                <p className="text-xs text-indigo-800 mb-2">
                                    Cliente nuevo sin grupo. Asignale un grupo de precios y los precios se ajustaran; al crear la orden el cliente quedara en ese grupo.
                                </p>
                                <div className="flex flex-wrap gap-2">
                                    {groups.map((g) => {
                                        const active = selectedGroupId === g.id;
                                        return (
                                            <button
                                                type="button"
                                                key={g.id}
                                                onClick={() => handleSelectGroup(g.id)}
                                                className="px-3 py-1.5 rounded-full text-xs font-bold border-2 transition-all"
                                                style={active
                                                    ? { backgroundColor: g.color || '#6b7280', borderColor: g.color || '#6b7280', color: '#fff' }
                                                    : { backgroundColor: '#fff', borderColor: g.color || '#6b7280', color: g.color || '#6b7280' }}
                                            >
                                                {g.name}
                                            </button>
                                        );
                                    })}
                                    {selectedGroupId && (
                                        <button
                                            type="button"
                                            onClick={() => handleSelectGroup(null)}
                                            className="px-3 py-1.5 rounded-full text-xs font-bold border-2 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 bg-white dark:bg-gray-800"
                                        >
                                            Sin grupo
                                        </button>
                                    )}
                                </div>
                                {selectedGroupId && (
                                    <p className="text-xs text-indigo-700 mt-2">
                                        Grupo <strong>{groups.find(g => g.id === selectedGroupId)?.name}</strong> aplicado: los precios se ajustaron y el cliente quedara en este grupo al crear la orden.
                                    </p>
                                )}
                            </div>
                        )}
                        <ProductSelector
                            businessId={formData.business_id || 0}
                            selectedProducts={selectedProducts}
                            onSelect={handleProductsChange}
                        />
                    </div>
                </div>

                <div className="lg:col-span-1">
                    <div className="bg-white dark:bg-gray-800 rounded-[14px] p-5 h-full" style={{ boxShadow: '0 2px 12px rgba(124, 58, 237, 0.06)' }}>
                        <div className="mb-3 flex items-center gap-2 border-b border-slate-100 pb-2.5 dark:border-gray-700">
                            <div className="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-lg text-white" style={{ backgroundColor: primaryColor }}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="white" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                                    <rect x="3" y="2" width="10" height="13" rx="1.5"/>
                                    <path d="M5.5 6h5M5.5 9h5M5.5 12h3"/>
                                </svg>
                            </div>
                            <h3 className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">Notas</h3>
                        </div>
                        <textarea
                            value={formData.notes}
                            onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                            rows={5}
                            placeholder="Notas internas sobre la orden..."
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-600 focus:border-transparent resize-none text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                        />
                    </div>
                </div>

                <div className="lg:col-span-3">
                    <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm dark:border-gray-700 dark:bg-gray-800">
                        <div className="mb-3 flex items-center gap-2 border-b border-slate-100 pb-2.5 dark:border-gray-700">
                            <div className="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-lg text-white" style={{ backgroundColor: primaryColor }}>
                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="white" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                                    <rect x="1" y="5" width="10" height="8" rx="1"/>
                                    <path d="M11 9h2.5L15 12v1h-1"/>
                                    <circle cx="4" cy="13" r="1.5"/>
                                    <circle cx="12" cy="13" r="1.5"/>
                                </svg>
                            </div>
                            <h3 className="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-gray-300">Logistica</h3>
                        </div>

                        <div className="mb-4 bg-blue-50 border-l-4 border-blue-400 p-3 rounded-r-lg">
                            <p className="text-[11px] leading-tight text-blue-800">
                                <span className="font-bold uppercase tracking-wider block mb-1">Aviso de Procesamiento</span>
                                Recuerda que las transportadoras pueden demorar hasta <span className="font-bold">24 horas habiles</span> en procesar y recolectar los pedidos despues de generada la guia.
                            </p>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                            <div>
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    Numero de Guia
                                </label>
                                <Input
                                    type="text"
                                    value={formData.tracking_number}
                                    onChange={(e) => setFormData({ ...formData, tracking_number: e.target.value })}
                                    placeholder="ej: 123456789"
                                />
                            </div>
                            <div>
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
                                    ID Guia
                                </label>
                                <Input
                                    type="text"
                                    value={formData.guide_id}
                                    onChange={(e) => setFormData({ ...formData, guide_id: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
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
                                <label className="mb-1.5 block text-[11px] font-semibold uppercase tracking-wide text-slate-400">
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

            </div>

            <footer className="flex shrink-0 flex-wrap items-center gap-3 border-t border-slate-200 bg-white px-4 py-3 dark:border-gray-700 dark:bg-gray-800">
                <div className="mr-auto">
                    <p className="text-[11px] font-semibold uppercase tracking-wide text-slate-400">Total de la orden</p>
                    <p className="text-xl font-bold tabular-nums text-slate-900 dark:text-white">
                        {formData.currency} {displayTotal.toLocaleString('es-CO')}
                    </p>
                </div>
                {onCancel && (
                    <button
                        type="button"
                        onClick={onCancel}
                        className="rounded-lg border border-slate-200 px-5 py-2 text-sm font-semibold text-slate-600 transition hover:bg-slate-50 dark:border-gray-600 dark:text-gray-200 dark:hover:bg-gray-700"
                    >
                        Cancelar
                    </button>
                )}
                <button
                    type="submit"
                    disabled={loading}
                    className="flex items-center gap-2 rounded-lg px-5 py-2 text-sm font-semibold text-white transition disabled:opacity-50"
                    style={{ backgroundColor: primaryColor }}
                >
                    {loading && <div className="spinner h-4 w-4" />}
                    {isEdit ? 'Actualizar Orden' : 'Crear Orden'}
                </button>
            </footer>

            <Modal
                isOpen={showProductModal}
                onClose={() => setShowProductModal(false)}
                title="Crear Nuevo Producto"
            >
                <ProductForm
                    onSuccess={() => {
                        setShowProductModal(false);
                    }}
                    onCancel={() => setShowProductModal(false)}
                />
            </Modal>
        </form>
    );
}
