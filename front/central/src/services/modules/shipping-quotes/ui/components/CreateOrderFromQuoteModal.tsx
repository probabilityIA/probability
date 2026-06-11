'use client';

import { useEffect, useMemo, useState } from 'react';
import dynamic from 'next/dynamic';
import { X, Package, Truck, User, MapPin, DollarSign, AlertTriangle, Loader2 } from 'lucide-react';
import { useToast } from '@/shared/providers/toast-provider';
import { getCarrierLogo } from '@/shared/utils/carrier-logos';
import { createOrderAction } from '@/services/modules/orders/infra/actions';
import { generateGuideAction } from '@/services/modules/shipments/infra/actions';
import { lookupGeozoneAction, getDeliveryProbabilityAction } from '@/services/modules/geozones/infra/actions';
import type { Geozone, ProbabilityResult } from '@/services/modules/geozones/domain/types';
import danes from '@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json';
import { associateQuoteAction, SavedQuote, SavedQuoteRate } from '../../infra/actions';

const GeozoneMiniMap = dynamic(
    () => import('@/services/modules/geozones/ui/components/GeozoneMiniMap').then(m => m.GeozoneMiniMap),
    { ssr: false },
);

interface Props {
    quote: SavedQuote;
    businessId: number | null;
    onClose: () => void;
    onSuccess: () => void;
}

const money = (v: number) => `$${Math.round(v).toLocaleString('es-CO')}`;

const titleCase = (s: string) =>
    s.toLowerCase().replace(/(^|\s)\S/g, c => c.toUpperCase());

const daneInfo = (code?: string): { city: string; state: string } | null => {
    if (!code) return null;
    const map = danes as Record<string, { ciudad: string; departamento: string }>;
    const d = map[code] || map[code.padEnd(8, '0')];
    return d ? { city: titleCase(d.ciudad), state: titleCase(d.departamento) } : null;
};

const geocodeAddress = async (address: string, city: string): Promise<{ lat: number; lng: number } | null> => {
    if (!city.trim()) return null;
    try {
        const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
        const res = await fetch(`${apiBase}/geocode?address=${encodeURIComponent(address)}&city=${encodeURIComponent(city)}`);
        if (!res.ok) return null;
        const d = await res.json();
        if (d?.found && typeof d.lat === 'number' && typeof d.lon === 'number') {
            return { lat: d.lat, lng: d.lon };
        }
        return null;
    } catch {
        return null;
    }
};

export default function CreateOrderFromQuoteModal({ quote, businessId, onClose, onSuccess }: Props) {
    const { showToast } = useToast();
    const payload = quote.request_payload || {};
    const destination = (payload.destination || {}) as Record<string, any>;
    const pkg = ((payload.packages || [])[0] || {}) as Record<string, any>;

    const rates = useMemo(() => {
        const list = [...(quote.rates || [])] as SavedQuoteRate[];
        list.sort((a, b) => (a.flete || 0) - (b.flete || 0));
        return list;
    }, [quote.rates]);

    const initialCOD = Number(payload.codValue || 0) > 0;
    const [isCOD, setIsCOD] = useState(initialCOD);
    const [autoGuide, setAutoGuide] = useState(false);
    const [confirmGuide, setConfirmGuide] = useState(false);
    const [rateIdx, setRateIdx] = useState(() => {
        const wanted = (quote.selected_carrier || '').toUpperCase().trim();
        if (wanted) {
            const i = rates.findIndex(r => (r.carrier || '').toUpperCase().trim() === wanted);
            if (i >= 0) return i;
        }
        if (initialCOD) {
            const i = rates.findIndex(r => (r as any).cod === true);
            if (i >= 0) return i;
        }
        return rates.length > 0 ? 0 : -1;
    });

    const destDane = daneInfo(String(destination.daneCode || ''));
    const originDane = daneInfo(String(((payload.origin || {}) as Record<string, any>).daneCode || ''));
    const originAddress = String(((payload.origin || {}) as Record<string, any>).address || '');

    const [firstName, setFirstName] = useState('');
    const [lastName, setLastName] = useState('');
    const [phone, setPhone] = useState(String(destination.phone || ''));
    const [email, setEmail] = useState(String(destination.email || ''));
    const [dni, setDni] = useState('');
    const [address, setAddress] = useState(String(destination.address || ''));
    const [city, setCity] = useState(String(destination.city || destDane?.city || ''));
    const [state, setState] = useState(String(destination.state || destDane?.state || ''));
    const [productValue, setProductValue] = useState<number>(Number(payload.contentValue || 0));

    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
    const [missingModal, setMissingModal] = useState<string[] | null>(null);

    const [destCoords, setDestCoords] = useState<{ lat: number; lng: number } | null>(null);
    const [originCoords, setOriginCoords] = useState<{ lat: number; lng: number } | null>(null);
    const [destZone, setDestZone] = useState<Geozone | null>(null);
    const [prob, setProb] = useState<ProbabilityResult | null>(null);

    const selectedRate = rateIdx >= 0 ? rates[rateIdx] : null;
    const effectiveBusinessId = businessId || quote.business_id;

    useEffect(() => {
        let cancelled = false;
        if (!originAddress || !originDane?.city) return;
        geocodeAddress(originAddress, originDane.city).then(coords => {
            if (!cancelled) setOriginCoords(coords);
        });
        return () => { cancelled = true; };
    }, [originAddress, originDane?.city]);

    useEffect(() => {
        let cancelled = false;
        const t = setTimeout(async () => {
            const coords = await geocodeAddress(address, city);
            if (!cancelled) setDestCoords(coords);
        }, 700);
        return () => { cancelled = true; clearTimeout(t); };
    }, [address, city]);

    useEffect(() => {
        let cancelled = false;
        if (!destCoords) {
            setDestZone(null);
            return;
        }
        lookupGeozoneAction({ lat: destCoords.lat, lng: destCoords.lng, business_id: effectiveBusinessId })
            .then((res: any) => {
                if (!cancelled) setDestZone(res?.data?.[0] || null);
            })
            .catch(() => { if (!cancelled) setDestZone(null); });
        return () => { cancelled = true; };
    }, [destCoords, effectiveBusinessId]);

    useEffect(() => {
        let cancelled = false;
        if (!destCoords) {
            setProb(null);
            return;
        }
        getDeliveryProbabilityAction({
            business_id: effectiveBusinessId,
            lat: destCoords.lat,
            lng: destCoords.lng,
            carrier: selectedRate?.carrier || undefined,
        })
            .then((res: any) => {
                if (!cancelled) setProb(res?.data || res || null);
            })
            .catch(() => { if (!cancelled) setProb(null); });
        return () => { cancelled = true; };
    }, [destCoords, effectiveBusinessId, selectedRate?.carrier]);

    const rateIsCOD = (selectedRate as any)?.cod === true;
    const fleteEstimate = selectedRate ? (selectedRate.flete || 0) + Number((selectedRate as any).codProbabilityMargin || 0) : 0;
    const codFee = isCOD ? Number((selectedRate as any)?.codCarrierFee || 0) : 0;
    const codToCollect = productValue + fleteEstimate + codFee;

    const expired = quote.expires_at ? new Date(quote.expires_at).getTime() < Date.now() : false;

    const selectRate = (i: number) => {
        setRateIdx(i);
        clearFieldError('rate');
        if (isCOD && (rates[i] as any)?.cod !== true) setIsCOD(false);
    };

    const toggleCOD = () => {
        const next = !isCOD;
        if (next && !rateIsCOD) {
            const i = rates.findIndex(r => (r as any).cod === true);
            if (i < 0) {
                showToast('Ninguna tarifa de esta cotizacion soporta contra entrega', 'error');
                return;
            }
            setRateIdx(i);
        }
        setIsCOD(next);
    };

    const clearFieldError = (field: string) => {
        setFieldErrors(prev => {
            if (!prev[field]) return prev;
            const next = { ...prev };
            delete next[field];
            return next;
        });
    };

    const inputCls = (field: string) => {
        const base = 'px-3 py-2 text-sm rounded-md border bg-white dark:bg-gray-700 text-gray-900 dark:text-white';
        return fieldErrors[field]
            ? `${base} border-red-400 dark:border-red-500 bg-red-50 dark:bg-red-900/20 focus:outline-red-500`
            : `${base} border-gray-200 dark:border-gray-600`;
    };

    const fieldError = (field: string) =>
        fieldErrors[field] ? <p className="mt-1 text-[11px] text-red-600 dark:text-red-400">{fieldErrors[field]}</p> : null;

    const validate = (): Record<string, string> => {
        const errs: Record<string, string> = {};
        if (!firstName.trim()) errs.firstName = 'Escribe el nombre del cliente';
        if (!lastName.trim()) errs.lastName = 'Escribe el apellido del cliente';
        if (!phone.trim()) errs.phone = 'Escribe el telefono del cliente';
        else if (!/^[\d\s+\-()]{7,}$/.test(phone.trim())) errs.phone = 'El telefono no parece valido';
        if (email.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim())) errs.email = 'El email no parece valido';
        if (!address.trim()) errs.address = 'Escribe la direccion de entrega';
        if (!city.trim()) errs.city = 'Escribe la ciudad de entrega';
        if (!state.trim()) errs.state = 'Escribe el departamento';
        if (!productValue || productValue <= 0) errs.productValue = 'El valor de los productos debe ser mayor a 0';
        if (autoGuide && !selectedRate) errs.rate = 'Selecciona una tarifa para generar la guia';
        return errs;
    };

    const handleSubmit = async () => {
        setError(null);
        const errs = validate();
        setFieldErrors(errs);
        if (Object.keys(errs).length > 0) {
            setMissingModal(Object.values(errs));
            return;
        }

        setSubmitting(true);
        try {
            const orderRes: any = await createOrderAction({
                business_id: businessId || quote.business_id,
                integration_id: 0,
                integration_type: 'platform',
                platform: 'manual',
                external_id: '',
                order_number: '',
                subtotal: productValue,
                tax: 0,
                discount: 0,
                shipping_cost: fleteEstimate,
                total_amount: productValue,
                currency: 'COP',
                cod_total: isCOD ? productValue + fleteEstimate : 0,
                customer_name: `${firstName} ${lastName}`.trim(),
                customer_first_name: firstName.trim(),
                customer_last_name: lastName.trim(),
                customer_email: email.trim(),
                customer_phone: phone.trim(),
                customer_dni: dni.trim(),
                shipping_street: address.trim(),
                shipping_city: city.trim(),
                shipping_state: state.trim(),
                shipping_country: 'Colombia',
                payment_method_id: isCOD ? 6 : 1,
                is_paid: false,
                items: [],
            } as any);

            const createdOrder = orderRes?.data || {};
            const orderId: string = createdOrder.id || createdOrder.ID || '';
            if (!orderRes?.success || !orderId) {
                setError(orderRes?.message || 'No se pudo crear la orden');
                setSubmitting(false);
                return;
            }

            const orderNumber: string = createdOrder.order_number || createdOrder.OrderNumber || orderId.slice(0, 8);
            let guideOk = false;

            if (autoGuide && selectedRate) {
                const guidePayload: any = {
                    ...payload,
                    idRate: (selectedRate as any).idRate,
                    carrier: selectedRate.carrier,
                    order_uuid: orderId,
                    external_order_id: orderId,
                    myShipmentReference: orderId,
                    contentValue: productValue,
                    codValue: isCOD ? productValue + fleteEstimate + codFee : 0,
                    codPaymentMethod: isCOD ? 'cash' : '',
                    includeGuideCost: false,
                    totalCost: fleteEstimate,
                    destination: {
                        ...destination,
                        firstName: firstName.trim(),
                        lastName: lastName.trim() || 'Cliente',
                        phone: phone.trim(),
                        email: email.trim(),
                        address: address.trim(),
                        city: city.trim(),
                        state: state.trim(),
                    },
                };
                if (isCOD && codFee > 0) guidePayload.codCarrierFee = codFee;

                const guideRes: any = await generateGuideAction(guidePayload);
                guideOk = guideRes?.success !== false;
                if (!guideOk) {
                    showToast(guideRes?.message || 'La orden se creo pero la guia fallo; generala desde la orden', 'error');
                }
            }

            const assocRes = await associateQuoteAction(quote.id, {
                order_uuid: orderId,
                selected_carrier: selectedRate?.carrier,
                selected_id_rate: (selectedRate as any)?.idRate,
                guide_requested: guideOk,
            }, businessId);
            if (!assocRes.success) {
                showToast(assocRes.message || 'Orden creada, pero no se pudo asociar la cotizacion', 'error');
            }

            showToast(
                guideOk
                    ? `Orden ${orderNumber} creada y guia solicitada con ${selectedRate?.carrier}`
                    : `Orden ${orderNumber} creada`,
                'success'
            );
            onSuccess();
            onClose();
        } catch (e: any) {
            setError(e?.message || 'Error inesperado al crear la orden');
        } finally {
            setSubmitting(false);
        }
    };

    return (
        <div className="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center z-50 p-3">
            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl w-full max-w-3xl max-h-[92vh] flex flex-col overflow-hidden">
                <div className="flex items-center justify-between px-5 py-3.5 border-b border-gray-100 dark:border-gray-700">
                    <div className="flex items-center gap-2">
                        <Package size={18} className="text-purple-600" />
                        <h2 className="text-base font-bold text-gray-900 dark:text-white">Crear orden desde cotizacion</h2>
                        <span className="text-xs text-gray-400">#{quote.id}</span>
                    </div>
                    <button onClick={onClose} className="p-1 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-500">
                        <X size={18} />
                    </button>
                </div>

                <div className="flex-1 overflow-y-auto p-5 space-y-5">
                    {error && (
                        <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg text-red-700 dark:text-red-400 text-sm">
                            {error}
                        </div>
                    )}
                    {expired && autoGuide && (
                        <div className="flex items-start gap-2 p-3 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-700 rounded-lg text-amber-700 dark:text-amber-400 text-xs">
                            <AlertTriangle size={14} className="mt-0.5 shrink-0" />
                            Esta cotizacion expiro: la tarifa pudo cambiar y la generacion de guia puede fallar. Si falla, genera la guia desde la orden.
                        </div>
                    )}

                    <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-3 text-xs text-gray-600 dark:text-gray-300 flex flex-wrap gap-x-6 gap-y-1">
                        <span><MapPin size={12} className="inline mr-1" />Destino: <strong>{String(destination.address || '-')}</strong> (DANE {String(destination.daneCode || '-')})</span>
                        <span>Paquete: {Number(pkg.weight || 0)}kg {Number(pkg.length || 0)}x{Number(pkg.width || 0)}x{Number(pkg.height || 0)}cm</span>
                        <span>Declarado: {money(Number(payload.contentValue || 0))}</span>
                    </div>

                    <div>
                        <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2 flex items-center gap-1.5">
                            <Truck size={14} className="text-purple-600" /> Tarifa
                        </h3>
                        <div className="space-y-1.5 max-h-44 overflow-y-auto pr-1">
                            {rates.map((r, i) => {
                                const cod = (r as any).cod === true;
                                const disabled = isCOD && !cod;
                                const logo = getCarrierLogo(r.carrier);
                                return (
                                    <label
                                        key={i}
                                        className={`flex items-center gap-3 px-3 py-2 rounded-lg border cursor-pointer text-sm transition-colors ${rateIdx === i ? 'border-purple-500 bg-purple-50 dark:bg-purple-900/20' : 'border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/40'} ${disabled ? 'opacity-40 cursor-not-allowed' : ''}`}
                                    >
                                        <input
                                            type="radio"
                                            name="rate"
                                            checked={rateIdx === i}
                                            disabled={disabled}
                                            onChange={() => selectRate(i)}
                                            className="accent-purple-600"
                                        />
                                        <span className="w-8 h-8 rounded-md border border-gray-200 dark:border-gray-600 bg-white flex items-center justify-center overflow-hidden shrink-0">
                                            {logo ? (
                                                <img
                                                    src={logo}
                                                    alt={r.carrier || ''}
                                                    className="w-7 h-7 object-contain"
                                                    onError={(e) => {
                                                        const el = e.currentTarget as HTMLImageElement;
                                                        el.style.display = 'none';
                                                        if (el.parentElement) el.parentElement.textContent = (r.carrier || '?').charAt(0).toUpperCase();
                                                    }}
                                                />
                                            ) : (
                                                <span className="text-xs font-bold text-gray-400">{(r.carrier || '?').charAt(0).toUpperCase()}</span>
                                            )}
                                        </span>
                                        <span className="font-medium text-gray-800 dark:text-gray-100 w-36 truncate">{r.carrier}</span>
                                        <span className="text-xs text-gray-500 w-24 truncate">{r.product}</span>
                                        <span className="text-xs text-gray-500">{r.deliveryDays ? `${r.deliveryDays}d` : ''}</span>
                                        <span className="ml-auto font-semibold text-gray-900 dark:text-white">{money((r.flete || 0) + Number((r as any).codProbabilityMargin || 0))}</span>
                                        {cod && (
                                            <span className="text-[10px] font-semibold px-1.5 py-0.5 rounded-full bg-emerald-100 text-emerald-700">
                                                COD {Number((r as any).codCarrierFee || 0) > 0 ? `+${money(Number((r as any).codCarrierFee))}` : ''}
                                            </span>
                                        )}
                                    </label>
                                );
                            })}
                            {rates.length === 0 && <p className="text-xs text-gray-400">La cotizacion no tiene tarifas guardadas.</p>}
                        </div>
                        {fieldError('rate')}
                    </div>

                    {(destCoords || destZone) && (
                        <div>
                            <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2 flex items-center gap-1.5">
                                <MapPin size={14} className="text-purple-600" /> Zona y probabilidad de entrega
                            </h3>
                            <GeozoneMiniMap
                                businessId={effectiveBusinessId}
                                geozone={destZone}
                                lat={destCoords?.lat ?? null}
                                lng={destCoords?.lng ?? null}
                                origin={originAddress ? { address: originAddress, lat: originCoords?.lat ?? null, lng: originCoords?.lng ?? null } : null}
                                destination={{ address: `${address}${city ? `, ${city}` : ''}` }}
                                carrierRate={prob?.delivery_rate ?? null}
                                carrierName={selectedRate?.carrier || null}
                                carrierEstimated={prob?.is_estimated}
                                height="170px"
                            />
                            {prob?.stats?.geozone_name && (
                                <p className="mt-1.5 text-[11px] text-gray-500 dark:text-gray-400">
                                    Estadisticas de la zona <strong>{prob.stats.geozone_name}</strong> ({prob.level || prob.stats.geozone_type}):
                                    {' '}{prob.stats.delivered} entregadas de {prob.stats.total} ordenes
                                    {prob.is_estimated ? ' (estimado por historial del carrier)' : ''}
                                </p>
                            )}
                        </div>
                    )}

                    <div>
                        <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2 flex items-center gap-1.5">
                            <User size={14} className="text-purple-600" /> Cliente
                        </h3>
                        <div className="grid grid-cols-2 gap-3">
                            <div>
                                <input value={firstName} onChange={e => { setFirstName(e.target.value); clearFieldError('firstName'); }} placeholder="Nombres *" className={`w-full ${inputCls('firstName')}`} />
                                {fieldError('firstName')}
                            </div>
                            <div>
                                <input value={lastName} onChange={e => { setLastName(e.target.value); clearFieldError('lastName'); }} placeholder="Apellidos *" className={`w-full ${inputCls('lastName')}`} />
                                {fieldError('lastName')}
                            </div>
                            <div>
                                <input value={phone} onChange={e => { setPhone(e.target.value); clearFieldError('phone'); }} placeholder="Telefono *" className={`w-full ${inputCls('phone')}`} />
                                {fieldError('phone')}
                            </div>
                            <div>
                                <input value={email} onChange={e => { setEmail(e.target.value); clearFieldError('email'); }} placeholder="Email" className={`w-full ${inputCls('email')}`} />
                                {fieldError('email')}
                            </div>
                            <input value={dni} onChange={e => setDni(e.target.value)} placeholder="Cedula / DNI" className={`col-span-2 ${inputCls('dni')}`} />
                        </div>
                    </div>

                    <div>
                        <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2 flex items-center gap-1.5">
                            <MapPin size={14} className="text-purple-600" /> Destino
                        </h3>
                        <div className="grid grid-cols-2 gap-3">
                            <div className="col-span-2">
                                <input value={address} onChange={e => { setAddress(e.target.value); clearFieldError('address'); }} placeholder="Direccion *" className={`w-full ${inputCls('address')}`} />
                                {fieldError('address')}
                            </div>
                            <div>
                                <input value={city} onChange={e => { setCity(e.target.value); clearFieldError('city'); }} placeholder="Ciudad *" className={`w-full ${inputCls('city')}`} />
                                {fieldError('city')}
                            </div>
                            <div>
                                <input value={state} onChange={e => { setState(e.target.value); clearFieldError('state'); }} placeholder="Departamento *" className={`w-full ${inputCls('state')}`} />
                                {fieldError('state')}
                            </div>
                        </div>
                    </div>

                    <div>
                        <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-2 flex items-center gap-1.5">
                            <DollarSign size={14} className="text-purple-600" /> Financiera
                        </h3>
                        <div className="grid grid-cols-2 gap-3 items-center">
                            <div>
                                <label className="block text-[11px] uppercase text-gray-400 mb-1">Valor productos (sin envio) *</label>
                                <input
                                    type="number"
                                    value={productValue}
                                    onChange={e => { setProductValue(parseFloat(e.target.value) || 0); clearFieldError('productValue'); }}
                                    className={`w-full ${inputCls('productValue')}`}
                                />
                                {fieldError('productValue')}
                            </div>
                            <div className="text-xs text-gray-500 dark:text-gray-400 space-y-0.5 pt-4">
                                <p>Envio estimado: <strong>{money(fleteEstimate)}</strong></p>
                                {isCOD && codFee > 0 && <p>Cargo COD carrier: <strong>{money(codFee)}</strong></p>}
                            </div>
                        </div>

                        <div className="mt-3 flex items-center justify-between rounded-lg border border-gray-200 dark:border-gray-700 px-3 py-2.5">
                            <div>
                                <p className="text-sm font-medium text-gray-700 dark:text-gray-200">Contra entrega</p>
                                {isCOD && (
                                    <p className="text-xs text-gray-500 dark:text-gray-400">
                                        Se cobrara al cliente: <strong className="text-emerald-600">{money(codToCollect)}</strong> (producto + envio{codFee > 0 ? ' + cargo COD' : ''})
                                    </p>
                                )}
                            </div>
                            <button
                                type="button"
                                role="switch"
                                aria-checked={isCOD}
                                onClick={toggleCOD}
                                className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors shrink-0 ${isCOD ? 'bg-emerald-500' : 'bg-gray-300 dark:bg-gray-600'}`}
                            >
                                <span className={`inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow-sm transition-transform ${isCOD ? 'translate-x-5' : 'translate-x-0.5'}`} />
                            </button>
                        </div>

                        <div className="mt-2 flex items-center justify-between rounded-lg border border-gray-200 dark:border-gray-700 px-3 py-2.5">
                            <div>
                                <p className="text-sm font-medium text-gray-700 dark:text-gray-200">Generar guia automaticamente</p>
                                <p className="text-xs text-gray-400">Al crear la orden se solicita la guia con la tarifa seleccionada.</p>
                            </div>
                            <button
                                type="button"
                                role="switch"
                                aria-checked={autoGuide}
                                onClick={() => {
                                    if (autoGuide) {
                                        setAutoGuide(false);
                                    } else {
                                        setConfirmGuide(true);
                                    }
                                }}
                                className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors shrink-0 ${autoGuide ? 'bg-purple-600' : 'bg-gray-300 dark:bg-gray-600'}`}
                            >
                                <span className={`inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow-sm transition-transform ${autoGuide ? 'translate-x-5' : 'translate-x-0.5'}`} />
                            </button>
                        </div>
                    </div>
                </div>

                <div className="flex items-center justify-end gap-2 px-5 py-3 border-t border-gray-100 dark:border-gray-700">
                    <button
                        onClick={onClose}
                        disabled={submitting}
                        className="px-4 py-2 text-sm rounded-lg border border-gray-200 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                    >
                        Cancelar
                    </button>
                    <button
                        onClick={handleSubmit}
                        disabled={submitting}
                        className="px-4 py-2 text-sm rounded-lg bg-purple-600 hover:bg-purple-700 text-white font-semibold disabled:opacity-60 flex items-center gap-2"
                    >
                        {submitting && <Loader2 size={14} className="animate-spin" />}
                        {submitting ? 'Creando...' : autoGuide ? 'Crear orden + guia' : 'Crear orden'}
                    </button>
                </div>
            </div>

            {missingModal && (
                <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-[60] p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-md p-5">
                        <div className="flex items-start gap-3">
                            <span className="w-9 h-9 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center shrink-0">
                                <AlertTriangle size={18} className="text-red-600" />
                            </span>
                            <div className="min-w-0">
                                <h3 className="text-sm font-bold text-gray-900 dark:text-white mb-1">
                                    Faltan datos para crear la orden
                                </h3>
                                <p className="text-xs text-gray-500 dark:text-gray-400 mb-2">
                                    Completa los campos marcados en rojo:
                                </p>
                                <ul className="space-y-1">
                                    {missingModal.map((m, i) => (
                                        <li key={i} className="text-xs text-red-600 dark:text-red-400 flex items-start gap-1.5">
                                            <span className="mt-0.5 shrink-0">•</span> {m}
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        </div>
                        <div className="flex justify-end mt-4">
                            <button
                                onClick={() => setMissingModal(null)}
                                className="px-4 py-1.5 text-xs rounded-lg bg-red-600 hover:bg-red-700 text-white font-semibold"
                            >
                                Entendido
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {confirmGuide && (
                <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-[60] p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-md p-5">
                        <div className="flex items-start gap-3">
                            <span className="w-9 h-9 rounded-full bg-amber-100 dark:bg-amber-900/30 flex items-center justify-center shrink-0">
                                <AlertTriangle size={18} className="text-amber-600" />
                            </span>
                            <div>
                                <h3 className="text-sm font-bold text-gray-900 dark:text-white mb-1">Generar guia real</h3>
                                <p className="text-xs text-gray-600 dark:text-gray-300">
                                    Al crear la orden se generara una <strong>guia real</strong> con
                                    {selectedRate ? ` ${selectedRate.carrier}` : ' la transportadora seleccionada'} por un costo estimado
                                    de <strong>{money(fleteEstimate)}</strong>. Este valor se descuenta de la billetera del negocio y
                                    la transportadora programara la recoleccion. Esta accion no se puede deshacer desde aqui.
                                </p>
                            </div>
                        </div>
                        <div className="flex justify-end gap-2 mt-4">
                            <button
                                onClick={() => setConfirmGuide(false)}
                                className="px-3 py-1.5 text-xs rounded-lg border border-gray-200 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700"
                            >
                                Cancelar
                            </button>
                            <button
                                onClick={() => { setAutoGuide(true); setConfirmGuide(false); }}
                                className="px-3 py-1.5 text-xs rounded-lg bg-amber-500 hover:bg-amber-600 text-white font-semibold"
                            >
                                Si, generar guia
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
