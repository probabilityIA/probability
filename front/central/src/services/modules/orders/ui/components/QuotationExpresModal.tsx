'use client';

import { useState, useEffect, useRef } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Input, Button, Stepper } from "@/shared/ui";
import { EnvioClickQuoteRequest, EnvioClickRate } from "@/services/modules/shipments/domain/types";
import { quoteShipmentAction } from "@/services/modules/shipments/infra/actions";
import { getWarehousesAction } from "@/services/modules/warehouses/infra/actions";
import { Warehouse } from "@/services/modules/warehouses/domain/types";
import { usePermissions } from "@/shared/contexts/permissions-context";
import danes from "@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json";
import { getActionError } from '@/shared/utils/action-result';
import { CookieStorage } from "@/shared/config";
import { lookupGeozoneAction, getDeliveryProbabilityAction } from "@/services/modules/geozones/infra/actions";
import type { Geozone, ProbabilityResult } from "@/services/modules/geozones/domain/types";
import '@/shared/ui/styles/shipment-modals.css';
import dynamic from 'next/dynamic';
import AddressAutocomplete from './AddressAutocomplete';

const GeozoneMiniMap = dynamic(
    () => import('@/services/modules/geozones/ui/components/GeozoneMiniMap').then(m => m.GeozoneMiniMap),
    { ssr: false }
);

const CarrierEffectivenessRates = dynamic(
    () => import('@/services/modules/geozones/ui/components/CarrierEffectivenessRates').then(m => m.CarrierEffectivenessRates),
    { ssr: false }
);

const normalizeLocationName = (str: string) => {
    if (!str) return "";
    let s = str.normalize("NFD").replace(/[̀-ͯ]/g, "").toUpperCase().trim();
    s = s.replace(/,\s*D\.C\./g, "").replace(/\sD\.C\./g, "").replace(/\sDC\b/g, "").trim();
    return s;
};


const getCarrierLogo = (carrierName: string): string => {
    const carrierLogos: { [key: string]: string } = {
        'SERVIENTREGA': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_servientrega.png',
        'COORDINADORA': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_coordinadora.png',
        'DHLEXPRESS': 'https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png',
        'DHL': 'https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png',
        'FEDEX': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_fedex.png',
        'INTERRAPIDISIMO': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_inerapidisimo.png',
        '472LOGISTICA': 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTnDF0ozRHf3s5BPqLsr7Vg-X8JRzECvFvwBQ&s',
        'SPEED': 'https://speedcargopa.com/wp-content/uploads/2021/03/Logo-mejorado-transparencia.png',
        'SPEEDCARGO': 'https://speedcargopa.com/wp-content/uploads/2021/03/Logo-mejorado-transparencia.png',
        'ENVIA': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_envia.png',
        'PIBOX': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_pibox.png',
        'TCC': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_TCC.png',
        'TRANSPORTADORADECARACOLOMBIA': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_TCC.png',
        '99MINUTOS': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_99minutos.webp',
        'DEPRISA': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_deprisa.png',
        'MENSAJERIAUBANA': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_mensajerosUrbanos.png',
        'MENSAJERIA URBANA': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_mensajerosUrbanos.png',
        'MENSAJEROS_URBANOS_EXPRESS': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_mensajerosUrbanos.png',
        'MENSAJEROSURBANOSEXPRESS': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_mensajerosUrbanos.png',
    };

    const trimmedName = carrierName.trim();
    if (carrierLogos[trimmedName]) {
        return carrierLogos[trimmedName];
    }

    const normalizedCarrier = normalizeLocationName(trimmedName).replace(/[_\s]/g, '');
    const normalizedLogos = Object.keys(carrierLogos).reduce((acc, key) => {
        acc[normalizeLocationName(key).replace(/[_\s]/g, '')] = carrierLogos[key];
        return acc;
    }, {} as Record<string, string>);
    return normalizedLogos[normalizedCarrier] || 'https://via.placeholder.com/56?text=' + encodeURIComponent(carrierName.substring(0, 3));
};

const formatCarrierName = (name: string): string => {
    return name
        .replace(/_/g, ' ')
        .replace(/EXPRESS/g, '')
        .trim();
};

const findDaneCode = (city: string, state: string) => {
    const targetCity = normalizeLocationName(city);
    const targetState = normalizeLocationName(state);
    if (!targetCity) return null;
    const entries = Object.entries(danes);
    const exactMatch = entries.find(([_, data]: [string, any]) => {
        const dCity = normalizeLocationName(data.ciudad);
        const dState = normalizeLocationName(data.departamento);
        return dCity === targetCity && dState === targetState;
    });
    if (exactMatch) return exactMatch[0];
    const cityMatch = entries.find(([_, data]: [string, any]) => {
        const dCity = normalizeLocationName(data.ciudad);
        return dCity === targetCity;
    });
    if (cityMatch) return cityMatch[0];
    return null;
};

const formSchema = z.object({
    originDaneCode: z.string().min(1, "Código DANE de origen requerido"),
    originAddress: z.string().min(2, "Dirección de origen requerida").max(50),
    destDaneCode: z.string().min(1, "Código DANE de destino requerido"),
    destAddress: z.string().min(8, "Dirección de destino requerida").max(50),
    weight: z.number().min(1).max(1000),
    height: z.number().min(1).max(300),
    width: z.number().min(1).max(300),
    length: z.number().min(1).max(300),
    description: z.string().min(3).max(25),
    contentValue: z.number().min(0).max(3000000),
});

type FormValues = z.infer<typeof formSchema>;

interface QuotationExpresModalProps {
    isOpen: boolean;
    onClose: () => void;
}

const STEPS = [
    { id: 1, label: "Origen y Destino" },
    { id: 2, label: "Cotizaciones" },
];

export function QuotationExpresModal({ isOpen, onClose }: QuotationExpresModalProps) {
    const { permissions, isSuperAdmin } = usePermissions();
    const effectiveBusinessId = permissions?.business_id || undefined;

    const [currentStep, setCurrentStep] = useState(1);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [rates, setRates] = useState<EnvioClickRate[]>([]);
    const [originWarehouses, setOriginWarehouses] = useState<Warehouse[]>([]);
    const [selectedOriginWarehouse, setSelectedOriginWarehouse] = useState<Warehouse | null>(null);
    const [originSearch, setOriginSearch] = useState("");
    const [destSearch, setDestSearch] = useState("");
    const [destLat, setDestLat] = useState<number | null>(null);
    const [destLng, setDestLng] = useState<number | null>(null);
    const [destGeozone, setDestGeozone] = useState<Geozone | null>(null);
    const [selectedRateId, setSelectedRateId] = useState<string | null>(null);
    const [selectedCarrier, setSelectedCarrier] = useState<string | null>(null);
    const [selectedCarrierProb, setSelectedCarrierProb] = useState<ProbabilityResult | null>(null);
    const [showOriginResults, setShowOriginResults] = useState(false);
    const [showDestResults, setShowDestResults] = useState(false);

    const originRef = useRef<HTMLDivElement>(null);
    const destRef = useRef<HTMLDivElement>(null);

    const [businessColors, setBusinessColors] = useState({
        primary: '#0f172a',
        secondary: '#be185d',
        tertiary: '#06b6d4',
        quaternary: '#f59e0b',
    });

    useEffect(() => {
        const loadColors = () => {
            const colors = CookieStorage.getBusinessColors();
            if (colors) {
                setBusinessColors({
                    primary: colors.primary || '#0f172a',
                    secondary: colors.secondary || '#be185d',
                    tertiary: colors.tertiary || '#06b6d4',
                    quaternary: colors.quaternary || '#f59e0b',
                });
            }
        };
        loadColors();
        window.addEventListener('businessChanged', loadColors);
        return () => window.removeEventListener('businessChanged', loadColors);
    }, []);

    const daneOptions = Object.entries(danes).map(([code, data]: [string, any]) => ({
        value: code,
        label: `${data.ciudad} (${data.departamento})`
    })).sort((a, b) => a.label.localeCompare(b.label));

    const filteredOriginOptions = daneOptions.filter(opt =>
        opt.label.toLowerCase().includes(originSearch.toLowerCase())
    );

    const filteredDestOptions = daneOptions.filter(opt =>
        opt.label.toLowerCase().includes(destSearch.toLowerCase())
    );

    const form = useForm<FormValues>({
        resolver: zodResolver(formSchema),
        mode: 'onChange',
        defaultValues: {
            originDaneCode: "",
            originAddress: "",
            destDaneCode: "",
            destAddress: "",
            weight: 1,
            height: 10,
            width: 10,
            length: 10,
            description: "E-commerce Order",
            contentValue: 0,
        },
    });

    const handleWarehouseSelect = (wh: Warehouse) => {
        console.log("🏭 handleWarehouseSelect:", wh);
        setSelectedOriginWarehouse(wh);
        const daneCode = wh.city_dane_code || findDaneCode(wh.city || "", wh.state || "") || "";
        console.log("Setting originDaneCode:", daneCode);
        form.setValue("originDaneCode", daneCode, { shouldValidate: true });
        const address = wh.street || wh.address;
        console.log("Setting originAddress:", address);
        form.setValue("originAddress", address, { shouldValidate: true });
        setOriginSearch(`${wh.city} (${wh.state})`);
    };

    useEffect(() => {
        if (isOpen) {
            getWarehousesAction({
                business_id: effectiveBusinessId || undefined,
                is_active: true,
                page: 1,
                page_size: 100,
            }).then(res => {
                if (res.data && res.data.length > 0) {
                    setOriginWarehouses(res.data);
                    const defaultWh = res.data.find(w => w.is_default);
                    if (defaultWh) {
                        handleWarehouseSelect(defaultWh);
                    }
                }
            }).catch(err => {
                console.error("Error loading warehouses:", err);
            });
        }
    }, [isOpen]);

    const handleStepChange = (step: number) => {
        setCurrentStep(step);
        if (step === 1) {
            setSelectedRateId(null);
            setSelectedCarrier(null);
            setSelectedCarrierProb(null);
        }
    };

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (originRef.current && !originRef.current.contains(event.target as Node)) {
                setShowOriginResults(false);
            }
            if (destRef.current && !destRef.current.contains(event.target as Node)) {
                setShowDestResults(false);
            }
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    useEffect(() => {
        if (!isOpen) {
            setCurrentStep(1);
            setRates([]);
            setError(null);
            setSelectedOriginWarehouse(null);
            setDestGeozone(null);
            setSelectedRateId(null);
            setSelectedCarrier(null);
            setSelectedCarrierProb(null);
            form.reset();
        }
    }, [isOpen]);

    useEffect(() => {
        console.log('🔍 Geozone effect triggered - lat:', destLat, 'lng:', destLng, 'business:', effectiveBusinessId, 'currentStep:', currentStep);

        if (destLat === null || destLat === undefined || destLng === null || destLng === undefined) {
            console.log('⏭️ Coordinates not ready');
            return;
        }

        if (!effectiveBusinessId && effectiveBusinessId !== 0) {
            console.log('⏭️ Business ID not ready');
            return;
        }

        let cancelled = false;
        (async () => {
            try {
                console.log('📍 Calling lookupGeozoneAction with:', { lat: destLat, lng: destLng, business_id: effectiveBusinessId });
                const geozones = await lookupGeozoneAction({
                    lat: destLat,
                    lng: destLng,
                    business_id: effectiveBusinessId || 0,
                });
                console.log('✅ Lookup response:', geozones);
                if (cancelled) {
                    console.log('Request was cancelled');
                    return;
                }
                if (geozones && geozones.data && Array.isArray(geozones.data) && geozones.data.length > 0) {
                    const geozone = geozones.data[0];
                    console.log('🎯 Setting geozone:', geozone.id, geozone.name, 'hasGeometry:', !!geozone.geometry);
                    setDestGeozone(geozone);
                } else {
                    console.log('⚠️ No geozones in response. Data:', geozones?.data);
                    setDestGeozone(null);
                }
            } catch (err) {
                console.error('❌ Error getting geozone:', err, 'message:', (err as any)?.message);
                if (!cancelled) {
                    setDestGeozone(null);
                }
            }
        })();

        return () => { cancelled = true; };
    }, [destLat, destLng, effectiveBusinessId]);

    useEffect(() => {
        if (destGeozone) {
            console.log('✅ destGeozone state updated:', { id: destGeozone.id, name: destGeozone.name, type: destGeozone.type, hasGeometry: !!destGeozone.geometry });
        } else {
            console.log('❌ destGeozone is now null or undefined');
        }
    }, [destGeozone]);

    const handleSubmit = async (data: FormValues) => {
        if (!data.originDaneCode || !data.destDaneCode) {
            setError("⚠️ Por favor selecciona códigos DANE válidos para origen y destino");
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const quotePayload: EnvioClickQuoteRequest = {
                packages: [{
                    weight: data.weight,
                    height: data.height,
                    width: data.width,
                    length: data.length,
                }],
                description: data.description,
                contentValue: data.contentValue,
                codValue: 0,
                includeGuideCost: false,
                insurance: false,
                codPaymentMethod: "cash",
                origin: {
                    daneCode: data.originDaneCode,
                    address: data.originAddress,
                },
                destination: {
                    daneCode: data.destDaneCode,
                    address: data.destAddress,
                },
            };

            const response = await quoteShipmentAction(quotePayload);
            if (!response.success) {
                setError(response.message || "Error al cotizar");
                setLoading(false);
                return;
            }

            const syncRates: EnvioClickRate[] = response.data?.data?.rates || [];
            if (syncRates.length > 0) {
                setRates(syncRates);
                setCurrentStep(2);
            } else {
                setError("No hay transportadoras disponibles para esta ruta");
            }
        } catch (err: any) {
            setError(getActionError(err, "Error al cotizar envío"));
        } finally {
            setLoading(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black/20 backdrop-blur-sm flex items-center justify-center z-50 p-2">
            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl flex flex-col overflow-hidden" style={{ width: '95%', maxWidth: '1200px', maxHeight: '90vh' }}>
                <div className="bg-white dark:bg-gray-800 border-b px-3 py-3 flex-shrink-0">
                    <div className="flex justify-between items-center mb-2">
                        <h2 className="text-2xl font-bold" style={{ color: businessColors.primary }}>Cotizador Expres</h2>
                        <button
                            onClick={onClose}
                            className="text-gray-500 dark:text-gray-400 hover:text-gray-700 text-2xl"
                        >
                            ×
                        </button>
                    </div>
                    <Stepper steps={STEPS} currentStep={currentStep} />
                </div>

                <div className="p-3 flex flex-col flex-1 overflow-hidden min-h-0">
                    {error && (
                        <div className="mb-3 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg text-red-700 dark:text-red-400 text-sm">
                            {error}
                        </div>
                    )}

                    {currentStep === 1 && (
                        <form id="quotation-form" onSubmit={form.handleSubmit(handleSubmit)} className="flex flex-col h-full overflow-hidden min-h-0">
                            <div className="flex-1 overflow-y-auto min-h-0 pr-3">
                                <div className="space-y-4">
                                    <div className="grid grid-cols-2 gap-4">
                                        <div className="shipment-section-origin">
                                            <div className="flex items-center justify-between">
                                                <div className="flex items-center gap-2">
                                                    <div className="w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold" style={{ backgroundColor: businessColors.tertiary + '30' }}>A</div>
                                                    <h3 className="font-semibold">Origen</h3>
                                                </div>
                                                {originWarehouses.length > 0 && (
                                                    <select
                                                        className="text-[11px] px-1.5 py-0.5 border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-1"
                                                        onChange={(e) => {
                                                            const wh = originWarehouses.find(a => a.id === parseInt(e.target.value));
                                                            if (wh) {
                                                                handleWarehouseSelect(wh);
                                                            }
                                                        }}
                                                        defaultValue=""
                                                    >
                                                        <option value="" disabled>Mis direcciones...</option>
                                                        {originWarehouses.map(a => (
                                                            <option key={a.id} value={a.id}>{a.name}</option>
                                                        ))}
                                                    </select>
                                                )}
                                            </div>

                                            <div ref={originRef} className="relative">
                                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                                                    Ciudad remitente *
                                                </label>
                                                <input
                                                    type="text"
                                                    value={originSearch}
                                                    onChange={(e) => {
                                                        setOriginSearch(e.target.value);
                                                        setShowOriginResults(true);
                                                        if (!e.target.value) form.setValue("originDaneCode", "", { shouldValidate: true });
                                                    }}
                                                    onFocus={() => setShowOriginResults(true)}
                                                    className="shipment-input"
                                                    placeholder="Buscar ciudad..."
                                                />
                                                {showOriginResults && filteredOriginOptions.length > 0 && (
                                                    <div className="absolute z-10 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md shadow-lg max-h-60 overflow-y-auto">
                                                        {filteredOriginOptions.slice(0, 50).map((opt) => (
                                                            <div
                                                                key={opt.value}
                                                                onClick={() => {
                                                                    form.setValue("originDaneCode", opt.value, { shouldValidate: true });
                                                                    setOriginSearch(opt.label);
                                                                    setShowOriginResults(false);
                                                                }}
                                                                className="px-3 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer text-sm"
                                                            >
                                                                {opt.label}
                                                            </div>
                                                        ))}
                                                    </div>
                                                )}
                                            </div>

                                            <div>
                                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                                                    Calle y Número *
                                                </label>
                                                <AddressAutocomplete
                                                    value={form.watch("originAddress")}
                                                    onChange={(val) => form.setValue("originAddress", val, { shouldValidate: true })}
                                                    city={originSearch}
                                                    placeholder="Calle 98 62-37"
                                                />
                                            </div>
                                        </div>

                                        <div className="shipment-section-destination">
                                            <div className="flex items-center gap-2">
                                                <div className="w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold" style={{ backgroundColor: businessColors.secondary + '30' }}>B</div>
                                                <h3 className="font-semibold">Destino</h3>
                                            </div>

                                            <div ref={destRef} className="relative">
                                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                                                    Ciudad destinatario *
                                                </label>
                                                <input
                                                    type="text"
                                                    value={destSearch}
                                                    onChange={(e) => {
                                                        setDestSearch(e.target.value);
                                                        setShowDestResults(true);
                                                        if (!e.target.value) form.setValue("destDaneCode", "", { shouldValidate: true });
                                                    }}
                                                    onFocus={() => setShowDestResults(true)}
                                                    className="shipment-input"
                                                    placeholder="Buscar ciudad..."
                                                />
                                                {showDestResults && filteredDestOptions.length > 0 && (
                                                    <div className="absolute z-10 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md shadow-lg max-h-60 overflow-y-auto">
                                                        {filteredDestOptions.slice(0, 50).map((opt) => (
                                                            <div
                                                                key={opt.value}
                                                                onClick={() => {
                                                                    form.setValue("destDaneCode", opt.value, { shouldValidate: true });
                                                                    setDestSearch(opt.label);
                                                                    setShowDestResults(false);
                                                                }}
                                                                className="px-3 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer text-sm"
                                                            >
                                                                {opt.label}
                                                            </div>
                                                        ))}
                                                    </div>
                                                )}
                                            </div>

                                            <div>
                                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                                                    Calle y Número *
                                                </label>
                                                <AddressAutocomplete
                                                    value={form.watch("destAddress")}
                                                    onChange={(val) => form.setValue("destAddress", val, { shouldValidate: true })}
                                                    onSelect={(suggestion) => {
                                                        if (suggestion.city && suggestion.state) {
                                                            const daneCode = findDaneCode(suggestion.city, suggestion.state);
                                                            if (daneCode) {
                                                                form.setValue("destDaneCode", daneCode, { shouldValidate: true });
                                                                setDestSearch(`${suggestion.city} (${suggestion.state})`);
                                                            }
                                                        }
                                                        if (suggestion.lat && suggestion.lon) {
                                                            setDestLat(suggestion.lat);
                                                            setDestLng(suggestion.lon);
                                                        }
                                                    }}
                                                    city={destSearch}
                                                    placeholder="Carrera 46 # 93 - 45"
                                                />
                                            </div>
                                        </div>
                                    </div>

                                    <div className="bg-gray-50/80 dark:bg-gray-700/30 border border-gray-200 dark:border-gray-600/30 rounded-xl p-4">
                                        <h3 className="font-semibold text-base mb-3">Paquete</h3>
                                        <div className="grid grid-cols-4 gap-2">
                                            <Input
                                                compact
                                                label="Peso (kg) *"
                                                type="number"
                                                step="0.1"
                                                {...form.register("weight", { valueAsNumber: true })}
                                            />
                                            <Input
                                                compact
                                                label="Alto (cm) *"
                                                type="number"
                                                {...form.register("height", { valueAsNumber: true })}
                                            />
                                            <Input
                                                compact
                                                label="Ancho (cm) *"
                                                type="number"
                                                {...form.register("width", { valueAsNumber: true })}
                                            />
                                            <Input
                                                compact
                                                label="Largo (cm) *"
                                                type="number"
                                                {...form.register("length", { valueAsNumber: true })}
                                            />
                                        </div>

                                        <div className="grid grid-cols-2 gap-2 mt-3">
                                            <Input
                                                compact
                                                label="Descripción *"
                                                {...form.register("description")}
                                                placeholder="descripción"
                                            />
                                            <Input
                                                compact
                                                label="Valor factura *"
                                                type="number"
                                                {...form.register("contentValue", { valueAsNumber: true })}
                                            />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </form>
                    )}

                    {currentStep === 2 && (
                        <div className="flex flex-row h-full gap-3 overflow-hidden">
                            <div className="w-1/3 h-full flex-shrink-0 border border-gray-200 dark:border-gray-600 rounded-lg p-2 overflow-hidden">
                                <GeozoneMiniMap
                                    businessId={effectiveBusinessId || 0}
                                    geozone={destGeozone}
                                    lat={destLat || 4.5709}
                                    lng={destLng || -74.2973}
                                    height="400px"
                                    showHeader={false}
                                    origin={selectedOriginWarehouse ? {
                                        address: [selectedOriginWarehouse.street || selectedOriginWarehouse.address, selectedOriginWarehouse.city, selectedOriginWarehouse.state].filter(Boolean).join(', '),
                                        lat: selectedOriginWarehouse.latitude ?? null,
                                        lng: selectedOriginWarehouse.longitude ?? null,
                                    } : null}
                                    destination={{
                                        address: destSearch ? `${destSearch}, ${form.watch("destAddress")}` : form.watch("destAddress"),
                                    }}
                                    carrierRate={selectedCarrierProb?.delivery_rate ?? null}
                                    carrierName={selectedCarrier || null}
                                    carrierEstimated={selectedCarrierProb?.is_estimated || false}
                                />
                            </div>

                            <div className="w-2/3 flex flex-col overflow-y-auto">
                                <div className="pb-2">
                                    <h3 className="font-semibold text-lg text-gray-700 dark:text-gray-200 mb-2">
                                        Filtra por servicio / Transportadora
                                    </h3>
                                </div>

                                <div className="overflow-y-auto border border-gray-200 dark:border-gray-600 rounded-lg p-3 bg-white dark:bg-gray-800 flex-1">
                                    {rates.length === 0 ? (
                                        <div className="flex items-center justify-center gap-3 py-10">
                                            <div className="shipment-spinner" style={{ width: 28, height: 28 }} />
                                            <span className="text-sm font-medium text-gray-600 dark:text-gray-300">Cargando cotizaciones...</span>
                                        </div>
                                    ) : (
                                        <div className="grid grid-cols-2 gap-4 auto-rows-max">
                                            {rates.map((rate) => {
                                                const minIns = rate.minimumInsurance ?? 0;
                                                const totalCost = rate.flete + minIns;

                                                const allCosts = rates.map(r => r.flete + (r.minimumInsurance ?? 0));
                                                const minCost = Math.min(...allCosts);
                                                const minDays = Math.min(...rates.map(r => r.deliveryDays || 999));

                                                const isFastest = rate.deliveryDays === minDays;
                                                const isCheapest = totalCost === minCost;
                                                const isSameDay = rate.deliveryDays === 0;
                                                const hasSpecialBadge = isCheapest || isFastest || isSameDay;

                                                let badgeColor = businessColors.quaternary;
                                                let badgeLabel = '';
                                                let borderColor = businessColors.tertiary;

                                                if (isCheapest) {
                                                    badgeColor = businessColors.secondary;
                                                    badgeLabel = 'MÁS ECONÓMICA';
                                                } else if (isSameDay) {
                                                    badgeColor = businessColors.quaternary;
                                                    badgeLabel = 'MISMO DÍA';
                                                } else if (isFastest) {
                                                    badgeColor = businessColors.tertiary;
                                                    badgeLabel = 'MÁS RÁPIDA';
                                                }

                                                const isSelected = selectedRateId === rate.idRate;

                                                return (
                                                    <div
                                                        key={rate.idRate}
                                                        onClick={async () => {
                                                            setSelectedRateId(rate.idRate);
                                                            setSelectedCarrier(rate.carrier);
                                                            try {
                                                                const prob = await getDeliveryProbabilityAction({
                                                                    business_id: effectiveBusinessId || 0,
                                                                    lat: destLat || undefined,
                                                                    lng: destLng || undefined,
                                                                    carrier: rate.carrier,
                                                                });
                                                                setSelectedCarrierProb(prob);
                                                            } catch (err) {
                                                                console.error('Error getting probability:', err);
                                                            }
                                                        }}
                                                        className={`relative border-2 rounded-3xl transition-all cursor-pointer ${
                                                            isSelected ? 'ring-2 ring-offset-2' : ''
                                                        } ${
                                                            hasSpecialBadge ? 'shadow-sm hover:shadow-lg hover:-translate-y-0.5' : 'shadow-sm hover:shadow-lg hover:-translate-y-0.5'
                                                        }`}
                                                        style={{
                                                            backgroundColor: `${businessColors.tertiary}08`,
                                                            borderColor: isSelected ? businessColors.tertiary : (hasSpecialBadge ? borderColor : '#d1d5db'),
                                                            borderWidth: isSelected ? '2px' : '2px',
                                                            padding: '18px',
                                                            display: 'flex',
                                                            flexDirection: 'column',
                                                            gap: '0',
                                                            ringColor: businessColors.tertiary,
                                                        }}
                                                    >
                                                        {hasSpecialBadge && (
                                                            <div
                                                                style={{
                                                                    position: 'absolute',
                                                                    top: '-14px',
                                                                    left: '50%',
                                                                    transform: 'translateX(-50%)',
                                                                    display: 'inline-flex',
                                                                    alignItems: 'center',
                                                                    justifyContent: 'center',
                                                                    padding: '6px 20px',
                                                                    borderRadius: '999px',
                                                                    fontSize: '13px',
                                                                    fontWeight: 700,
                                                                    color: '#fff',
                                                                    backgroundColor: badgeColor,
                                                                    border: '3px solid #fff',
                                                                    boxShadow: `0 4px 12px ${badgeColor}66`,
                                                                    whiteSpace: 'nowrap',
                                                                }}
                                                            >
                                                                {badgeLabel}
                                                            </div>
                                                        )}

                                                        <div style={{ display: 'flex', gap: '0', flex: 1 }}>
                                                            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', flex: '0 0 50%', textAlign: 'center' }}>
                                                                <div
                                                                    className="flex-shrink-0 rounded-[14px] overflow-hidden flex items-center justify-center bg-gradient-to-br from-purple-100 to-pink-100"
                                                                    style={{
                                                                        width: '65px',
                                                                        height: '65px',
                                                                        margin: '0 auto',
                                                                    }}
                                                                >
                                                                    <img
                                                                        src={getCarrierLogo(rate.carrier)}
                                                                        alt={rate.carrier}
                                                                        className="w-full h-full object-contain"
                                                                        onError={(e) => {
                                                                            e.currentTarget.style.display = 'none';
                                                                            e.currentTarget.parentElement!.textContent = rate.carrier.substring(0, 2);
                                                                            e.currentTarget.parentElement!.style.fontSize = '20px';
                                                                            e.currentTarget.parentElement!.style.fontWeight = '700';
                                                                        }}
                                                                    />
                                                                </div>

                                                                <div>
                                                                    <div style={{
                                                                        fontSize: '12px',
                                                                        fontWeight: 800,
                                                                        color: '#0f1417',
                                                                        letterSpacing: '.01em',
                                                                        lineHeight: 1.3,
                                                                        wordBreak: 'break-word',
                                                                        maxWidth: '100%',
                                                                    }}>
                                                                        {formatCarrierName(rate.carrier)}
                                                                    </div>
                                                                    <div style={{ fontSize: '11px', color: '#6b757c', marginTop: '2px' }}>
                                                                        {rate.product}
                                                                    </div>
                                                                </div>

                                                                <div
                                                                    style={{
                                                                        display: 'flex',
                                                                        alignItems: 'baseline',
                                                                        justifyContent: 'center',
                                                                        gap: '4px',
                                                                        fontSize: '28px',
                                                                        fontWeight: 700,
                                                                        color: '#0f1417',
                                                                        letterSpacing: '-.02em',
                                                                        lineHeight: 1,
                                                                        fontVariantNumeric: 'tabular-nums',
                                                                        margin: '4px 0 0',
                                                                    }}
                                                                >
                                                                    <span>${totalCost.toLocaleString()}</span>
                                                                    <span style={{ fontSize: '10px', color: '#6b757c', fontWeight: 500 }}>COP</span>
                                                                </div>

                                                                <div style={{ display: 'flex', flexDirection: 'column', gap: '2px', fontSize: '10px', color: '#3b4248', marginTop: '4px', textAlign: 'center', fontWeight: 700 }}>
                                                                    <div>Costo: ${totalCost.toLocaleString()}</div>
                                                                    <div>Guía: ${rate.flete.toLocaleString()}</div>
                                                                    {minIns > 0 ? (
                                                                        <div style={{ color: '#059669', fontSize: '9px' }}>
                                                                            (Seguro: ${minIns.toLocaleString()})
                                                                        </div>
                                                                    ) : (
                                                                        <div style={{ color: '#6b757c', fontSize: '9px' }}>(Sin asegurar)</div>
                                                                    )}
                                                                </div>

                                                                <div style={{
                                                                    display: 'inline-flex',
                                                                    alignItems: 'center',
                                                                    justifyContent: 'center',
                                                                    gap: '0',
                                                                    padding: '0',
                                                                    borderRadius: '0',
                                                                    fontSize: '11px',
                                                                    fontWeight: 600,
                                                                    backgroundColor: 'transparent',
                                                                    color: rate.deliveryDays === 0 || rate.deliveryDays <= 1 ? '#0891b2' : '#3b4248',
                                                                    whiteSpace: 'nowrap',
                                                                    border: 'none',
                                                                    margin: '2px auto',
                                                                }}>
                                                                    {rate.deliveryDays === 0 ? 'Mismo día' : rate.deliveryDays === 1 ? '1 día' : `${rate.deliveryDays} días`}
                                                                </div>
                                                            </div>

                                                            <div style={{
                                                                width: '1px',
                                                                backgroundColor: '#e5e7eb',
                                                                margin: '12px 0',
                                                                flexShrink: 0,
                                                            }}></div>

                                                            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', flex: '0 0 50%', padding: '0 12px', justifyContent: 'center' }}>
                                                                <CarrierEffectivenessRates
                                                                    businessId={effectiveBusinessId || 0}
                                                                    carrier={rate.carrier}
                                                                />
                                                            </div>
                                                        </div>
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                    )}
                </div>

                <div className="bg-white dark:bg-gray-800 border-t px-3 py-3 flex-shrink-0">
                    <div className="flex gap-2">
                        {currentStep === 2 && (
                            <Button
                                variant="secondary"
                                onClick={() => handleStepChange(1)}
                                className="flex-1"
                            >
                                Atrás
                            </Button>
                        )}
                        <Button
                            variant="secondary"
                            onClick={onClose}
                            className="flex-1"
                        >
                            {currentStep === 1 ? 'Cancelar' : 'Cerrar'}
                        </Button>
                        {currentStep === 1 && (
                            <Button
                                type="submit"
                                form="quotation-form"
                                loading={loading}
                                className="flex-1"
                                style={{ backgroundColor: businessColors.tertiary }}
                            >
                                Cotizar
                            </Button>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}
