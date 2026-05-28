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
import danes from "@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json";
import { getActionError } from '@/shared/utils/action-result';
import { CookieStorage } from "@/shared/config";
import AddressAutocomplete, { AddressSuggestion } from './AddressAutocomplete';
import MapComponent from '@/shared/ui/MapComponent';
import '@/shared/ui/styles/shipment-modals.css';

const normalizeLocationName = (str: string) => {
    if (!str) return "";
    let s = str.normalize("NFD").replace(/[̀-ͯ]/g, "").toUpperCase().trim();
    s = s.replace(/,\s*D\.C\./g, "").replace(/\sD\.C\./g, "").replace(/\sDC\b/g, "").trim();
    return s;
};

const normalizeCity = (s: string) =>
    s.normalize('NFD').replace(/[̀-ͯ]/g, '')
        .toLowerCase()
        .replace(/\s*[,(]?\s*d\.?\s*c\.?\s*\)?\s*$/g, '')
        .trim();

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

const buildWarehouseAddress = (warehouse: Warehouse | null): string => {
    if (!warehouse) return "";
    return warehouse.address?.trim() || "";
};

const formatCarrierName = (carrierName: string): string => {
    return carrierName
        .split('_')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
        .join(' ');
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
        'MENSAJERIAURBANOSEXPRESS': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_mensajerosUrbanos.png',
    };
    const normalizedName = carrierName.trim().toUpperCase().replace(/\s+/g, '');
    if (carrierLogos[normalizedName]) return carrierLogos[normalizedName];
    const nameWithSpaces = carrierName.trim().toUpperCase();
    if (carrierLogos[nameWithSpaces]) return carrierLogos[nameWithSpaces];
    return 'https://via.placeholder.com/56?text=' + encodeURIComponent(carrierName.substring(0, 3));
};

const formSchema = z.object({
    originDaneCode: z.string().min(1, "Código DANE de origen requerido"),
    originAddress: z.string().min(2, "Dirección de origen requerida").max(200),
    destDaneCode: z.string().min(1, "Código DANE de destino requerido"),
    destAddress: z.string().min(8, "Dirección de destino requerida").max(200),
    weight: z.number().min(1).max(1000),
    height: z.number().min(1).max(300),
    width: z.number().min(1).max(300),
    length: z.number().min(1).max(300),
    description: z.string().min(3).max(100),
    contentValue: z.number().min(1, "Valor a facturar es obligatorio").max(3000000),
    codValue: z.number().min(0).max(3000000).optional(),
    enableCod: z.boolean(),
    enableInsurance: z.boolean(),
});

type FormValues = z.infer<typeof formSchema>;

interface QuotationExpresModalProps {
    isOpen: boolean;
    onClose: () => void;
    business_id?: number;
}

const STEPS = [
    { id: 1, label: "Origen y Destino" },
    { id: 2, label: "Cotizaciones" },
];

export function QuotationExpresModal({ isOpen, onClose, business_id }: QuotationExpresModalProps) {
    const [currentStep, setCurrentStep] = useState(1);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [rates, setRates] = useState<EnvioClickRate[]>([]);
    const [originWarehouses, setOriginWarehouses] = useState<Warehouse[]>([]);
    const [selectedOriginWarehouse, setSelectedOriginWarehouse] = useState<Warehouse | null>(null);
    const [originSearch, setOriginSearch] = useState("");
    const [destSearch, setDestSearch] = useState("");
    const [showOriginResults, setShowOriginResults] = useState(false);
    const [showDestResults, setShowDestResults] = useState(false);
    const [selectedRate, setSelectedRate] = useState<number | null>(null);

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
        label: `${data.ciudad} (${data.departamento})`,
        ciudad: data.ciudad,
        departamento: data.departamento
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
            codValue: 0,
            enableCod: false,
            enableInsurance: false,
        },
    });

    useEffect(() => {
        if (isOpen) {
            getWarehousesAction({
                is_active: true,
                page: 1,
                page_size: 100,
                ...(business_id && { business_id }),
            }).then(res => {
                if (res.data) {
                    setOriginWarehouses(res.data);
                    const defaultWh = res.data.find(w => w.is_default);
                    if (defaultWh) {
                        setSelectedOriginWarehouse(defaultWh);
                        const daneCode = defaultWh.city_dane_code || findDaneCode(defaultWh.city || "", defaultWh.state || "") || "";
                        form.setValue("originDaneCode", daneCode, { shouldValidate: true });
                        const warehouseAddress = buildWarehouseAddress(defaultWh);
                        form.setValue("originAddress", warehouseAddress, { shouldValidate: true });
                        setOriginSearch(`${defaultWh.city} (${defaultWh.state})`);
                    }
                }
            }).catch(() => { });
        }
    }, [isOpen]);

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
            form.reset();
        }
    }, [isOpen]);

    const enrichRatesWithEffectivity = async (rates: EnvioClickRate[], destDaneCode: string): Promise<EnvioClickRate[]> => {
        try {
            const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
            const response = await fetch(`${apiBase}/geozones/probability-by-dane?dane_code=${destDaneCode}`, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' }
            });

            if (!response.ok) return rates;

            const result = await response.json();
            if (!result?.data) return rates;

            const probData = result.data;
            const deliveryRate = probData.delivery_rate !== null && probData.delivery_rate !== undefined ? probData.delivery_rate / 100 : undefined;
            const collectionRate = probData.collection_rate !== null && probData.collection_rate !== undefined ? probData.collection_rate / 100 : undefined;

            return rates.map(rate => ({
                ...rate,
                deliveryRate,
                collectionRate,
            }));
        } catch (err) {
            console.error('Error enriching rates with effectivity:', err);
            return rates;
        }
    };

    const handleSubmit = async (data: FormValues) => {
        if (!data.originDaneCode || !data.destDaneCode) {
            setError("⚠️ Por favor selecciona códigos DANE válidos para origen y destino");
            return;
        }

        setLoading(true);
        setError(null);
        try {
            console.log('Form data antes de crear payload:', {
                originDaneCode: data.originDaneCode,
                originAddress: data.originAddress,
                destDaneCode: data.destDaneCode,
                destAddress: data.destAddress,
                originSearch,
                destSearch
            });

            const destDaneData = daneOptions.find(opt => opt.value === data.destDaneCode);
            const originDaneData = daneOptions.find(opt => opt.value === data.originDaneCode);

            const contactNameParts = selectedOriginWarehouse?.contact_name?.split(" ") || ["", ""];
            const firstName = contactNameParts[0] || "Warehouse";
            const lastName = contactNameParts.slice(1).join(" ") || "Contact";

            const quotePayload: any = {
                ...(business_id && { business_id }),
                packages: [{
                    weight: data.weight,
                    height: data.height,
                    width: data.width,
                    length: data.length,
                }],
                description: data.description,
                contentValue: data.contentValue,
                codValue: data.enableCod ? (data.codValue || data.contentValue) : 0,
                includeGuideCost: false,
                insurance: data.enableInsurance,
                codPaymentMethod: "cash",
                myShipmentReference: `REF-${Date.now()}`,
                external_order_id: `EXT-${Date.now()}`,
                requestPickup: false,
                pickupDate: new Date().toISOString().split("T")[0],
                origin: {
                    daneCode: data.originDaneCode,
                    address: data.originAddress,
                    firstName,
                    lastName,
                    email: selectedOriginWarehouse?.contact_email || "",
                    phone: selectedOriginWarehouse?.phone || "",
                    company: "",
                    crossStreet: "",
                    reference: "",
                    city: originDaneData?.ciudad || "",
                    state: originDaneData?.departamento || "",
                },
                destination: {
                    daneCode: data.destDaneCode,
                    address: data.destAddress,
                    firstName: "Destinatario",
                    lastName: "Default",
                    email: "",
                    phone: "",
                    company: "",
                    crossStreet: "",
                    reference: "",
                    city: destDaneData?.ciudad || "",
                    state: destDaneData?.departamento || "",
                },
            };

            const response = await quoteShipmentAction(quotePayload);

            if (!response.success) {
                setError(response.message || "Error al cotizar");
                setLoading(false);
                return;
            }

            const syncRates: EnvioClickRate[] = response.data?.rates || [];

            if (syncRates.length > 0) {
                const enrichedRates = await enrichRatesWithEffectivity(syncRates, data.destDaneCode);
                setRates(enrichedRates);
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
            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl flex flex-col overflow-hidden" style={{ width: '85%', maxHeight: '90vh' }}>
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
                        <form
                            onSubmit={form.handleSubmit(
                                handleSubmit,
                                (errors) => {
                                    const errorMessages = Object.entries(errors)
                                        .map(([field, error]: [string, any]) => `${field}: ${error?.message}`)
                                        .join('\n');
                                    setError(`Por favor completa los campos:\n${errorMessages}`);
                                }
                            )}
                            className="flex flex-col h-full overflow-hidden min-h-0"
                        >
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
                                                        className="text-[11px] px-1.5 py-0.5 border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                                                        onChange={(e) => {
                                                            const addr = originWarehouses.find(a => a.id === parseInt(e.target.value));
                                                            if (addr) {
                                                                setSelectedOriginWarehouse(addr);
                                                                const daneCode = addr.city_dane_code || findDaneCode(addr.city || "", addr.state || "") || "";
                                                                form.setValue("originDaneCode", daneCode, { shouldValidate: true });
                                                                const warehouseAddress = buildWarehouseAddress(addr);
                                                                form.setValue("originAddress", warehouseAddress, { shouldValidate: true });
                                                                setOriginSearch(`${addr.city} (${addr.state})`);
                                                            }
                                                        }}
                                                        value={selectedOriginWarehouse?.id || ""}
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

                                            <div className="space-y-1">
                                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200">
                                                    Dirección del remitente *
                                                </label>
                                                <div className="px-3 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md text-sm text-gray-900 dark:text-gray-100">
                                                    {selectedOriginWarehouse ? buildWarehouseAddress(selectedOriginWarehouse) : "Selecciona una bodega"}
                                                </div>
                                                <input
                                                    type="hidden"
                                                    {...form.register("originAddress")}
                                                />
                                            </div>
                                        </div>

                                        <div className="shipment-section-destination">
                                            <div className="flex items-center justify-between mb-3">
                                                <div className="flex items-center gap-2">
                                                    <div className="w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold" style={{ backgroundColor: businessColors.secondary + '30' }}>B</div>
                                                    <h3 className="font-semibold">Destino</h3>
                                                </div>
                                                <div className="flex gap-4">
                                                    <label className="flex items-center gap-2 cursor-pointer">
                                                        <input
                                                            type="checkbox"
                                                            {...form.register("enableCod")}
                                                            className="w-4 h-4 rounded border-gray-300 text-gray-600"
                                                        />
                                                        <span className="text-sm text-gray-700 dark:text-gray-200">Contra entrega</span>
                                                    </label>
                                                    <label className="flex items-center gap-2 cursor-pointer">
                                                        <input
                                                            type="checkbox"
                                                            {...form.register("enableInsurance")}
                                                            className="w-4 h-4 rounded border-gray-300 text-gray-600"
                                                        />
                                                        <span className="text-sm text-gray-700 dark:text-gray-200">Asegurar</span>
                                                    </label>
                                                </div>
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
                                                    className={`shipment-input ${form.formState.errors.destDaneCode ? 'border-red-500 focus:ring-red-500' : ''}`}
                                                    placeholder="Buscar ciudad..."
                                                />
                                                {form.formState.errors.destDaneCode && (
                                                    <p className="text-xs text-red-600 dark:text-red-400 mt-1">{form.formState.errors.destDaneCode.message}</p>
                                                )}
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
                                                    city={destSearch}
                                                    onSelect={(s: AddressSuggestion) => {
                                                        form.setValue("destAddress", s.display_name, { shouldValidate: true });
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
                                                                form.setValue("destDaneCode", match.value, { shouldValidate: true });
                                                                setDestSearch(match.label);
                                                                setShowDestResults(false);
                                                            }
                                                        }
                                                    }}
                                                    placeholder="Carrera 46 # 93 - 45"
                                                />
                                                {form.formState.errors.destAddress && (
                                                    <p className="text-xs text-red-600 dark:text-red-400 mt-1">{form.formState.errors.destAddress.message}</p>
                                                )}
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
                                            <div/>
                                        </div>

                                        <div className="grid grid-cols-2 gap-2 mt-2">
                                            <div>
                                                <Input
                                                    compact
                                                    label="Valor factura *"
                                                    type="number"
                                                    {...form.register("contentValue", { valueAsNumber: true })}
                                                />
                                                {form.formState.errors.contentValue && (
                                                    <p className="text-xs text-red-600 dark:text-red-400 mt-1">{form.formState.errors.contentValue.message}</p>
                                                )}
                                            </div>
                                            {form.watch("enableCod") && (
                                                <div>
                                                    <Input
                                                        compact
                                                        label="Valor contra entrega"
                                                        type="number"
                                                        {...form.register("codValue", { valueAsNumber: true })}
                                                    />
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div className="flex gap-2 mt-4 pt-4 border-t border-gray-200 dark:border-gray-600">
                                <Button
                                    variant="secondary"
                                    onClick={onClose}
                                    className="flex-1"
                                >
                                    Cancelar
                                </Button>
                                <Button
                                    type="submit"
                                    loading={loading}
                                    className="flex-1"
                                    style={{ backgroundColor: businessColors.tertiary }}
                                >
                                    Cotizar
                                </Button>
                            </div>
                        </form>
                    )}

                    {currentStep === 2 && (
                        <div className="flex flex-col h-full min-h-0">
                            <div className="mb-6">
                                <div className="flex justify-between items-center">
                                    <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">Cotizaciones Disponibles</h2>
                                    <span className="text-xs font-semibold text-gray-600 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-3 py-1 rounded-full">
                                        {rates.filter(r => !form.watch("enableCod") || r.cod).length} transportadoras encontradas
                                    </span>
                                </div>
                            </div>

                            <div className="flex-1 overflow-y-auto min-h-0 pr-3">
                                {rates.filter(rate => !form.watch("enableCod") || rate.cod).length === 0 && (
                                    <div className="text-center py-12 text-gray-500 dark:text-gray-400">
                                        {form.watch("enableCod") ? (
                                            <p>No hay transportistas que soporten contra entrega para esta ruta</p>
                                        ) : (
                                            <p>No hay cotizaciones disponibles</p>
                                        )}
                                    </div>
                                )}

                                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 pb-4">
                                    {rates.filter(rate => {
                                        if (form.watch("enableCod") && !rate.cod) return false;
                                        return true;
                                    }).map((rate, index) => {
                                        if (index === 0 && selectedRate === null) {
                                            setSelectedRate(rate.idRate);
                                        }

                                        const getDisplayPrice = (r: EnvioClickRate) => {
                                            const basePrice = r.flete;
                                            const minimumIns = r.minimumInsurance ?? 0;
                                            const insuranceCost = form.watch("enableInsurance") ? (r.extraInsurance ?? 0) : 0;
                                            const probabilityMargin = form.watch("enableCod") ? (r.codProbabilityMargin ?? 0) : 0;
                                            return basePrice + minimumIns + insuranceCost + probabilityMargin;
                                        };

                                        const isSelected = selectedRate === rate.idRate;
                                        const displayPrice = getDisplayPrice(rate);

                                        return (
                                            <div
                                                key={rate.idRate}
                                                onClick={() => setSelectedRate(isSelected ? null : rate.idRate)}
                                                className={`relative border-2 rounded-2xl p-6 cursor-pointer transition-all duration-300 overflow-hidden
                                                    ${isSelected
                                                        ? 'shadow-lg dark:shadow-lg'
                                                        : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600 hover:shadow-md'
                                                    }
                                                `}
                                                style={{
                                                    borderColor: isSelected ? businessColors.tertiary : '#e5e7eb',
                                                    backgroundColor: isSelected ? businessColors.tertiary + '05' : 'transparent',
                                                    boxShadow: isSelected ? `0 0 0 1px ${businessColors.tertiary}, 0 12px 32px ${businessColors.tertiary}20` : undefined
                                                }}
                                            >
                                                {isSelected && (
                                                    <div
                                                        className="absolute inset-0 top-0 h-1 opacity-40"
                                                        style={{
                                                            background: `radial-gradient(ellipse at center, ${businessColors.tertiary} 0%, transparent 70%)`
                                                        }}
                                                    />
                                                )}

                                                <div className="flex gap-3 mb-5 items-center">
                                                    <div
                                                        className="w-12 h-12 rounded-xl flex-shrink-0 flex items-center justify-center"
                                                        style={{
                                                            background: `linear-gradient(135deg, ${businessColors.tertiary}14 0%, ${businessColors.tertiary}08 100%)`
                                                        }}
                                                    >
                                                        <img
                                                            src={getCarrierLogo(rate.carrier)}
                                                            alt={rate.carrier}
                                                            className="w-10 h-10 object-contain"
                                                            onError={(e) => {
                                                                (e.target as HTMLImageElement).style.display = 'none';
                                                            }}
                                                        />
                                                    </div>

                                                    <div>
                                                        <div className="flex items-center gap-2">
                                                            <h4 className="font-semibold text-gray-900 dark:text-gray-100">{formatCarrierName(rate.carrier)}</h4>
                                                            <span className="text-xs font-semibold text-gray-600 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded-full">
                                                                {rate.product || 'Normal'}
                                                            </span>
                                                        </div>
                                                    </div>
                                                </div>

                                                <div className="grid grid-cols-2 gap-6 mb-5">
                                                    <div className="space-y-3">
                                                        <h5 className="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-3">Desglose</h5>

                                                        <div className="flex justify-between items-center text-sm">
                                                            <div className="flex items-center gap-2">
                                                                <div className="w-2 h-2 rounded-full" style={{ backgroundColor: '#999' }}></div>
                                                                <span className="text-gray-700 dark:text-gray-300 font-medium">Flete</span>
                                                            </div>
                                                            <span className="text-gray-900 dark:text-gray-100 font-semibold" style={{ fontVariantNumeric: 'tabular-nums' }}>
                                                                ${rate.flete.toLocaleString('es-CO', { maximumFractionDigits: 0 })}
                                                            </span>
                                                        </div>

                                                        {(rate.codProbabilityMargin ?? 0) > 0 && (
                                                            <div className="flex justify-between items-center text-sm">
                                                                <div className="flex items-center gap-2">
                                                                    <div className="w-2 h-2 rounded-full" style={{ backgroundColor: businessColors.tertiary }}></div>
                                                                    <span className="text-gray-700 dark:text-gray-300 font-medium">Comisión Probability</span>
                                                                </div>
                                                                <span className="text-gray-900 dark:text-gray-100 font-semibold" style={{ fontVariantNumeric: 'tabular-nums' }}>
                                                                    ${(rate.codProbabilityMargin ?? 0).toLocaleString('es-CO', { maximumFractionDigits: 0 })}
                                                                </span>
                                                            </div>
                                                        )}

                                                        {(rate.minimumInsurance ?? 0) > 0 && (
                                                            <div className="flex justify-between items-center text-sm">
                                                                <div className="flex items-center gap-2">
                                                                    <div className="w-2 h-2 rounded-full" style={{ backgroundColor: '#ff9500' }}></div>
                                                                    <span className="text-gray-700 dark:text-gray-300 font-medium">Seguro mín.</span>
                                                                </div>
                                                                <span className="text-gray-900 dark:text-gray-100 font-semibold" style={{ fontVariantNumeric: 'tabular-nums' }}>
                                                                    ${(rate.minimumInsurance ?? 0).toLocaleString('es-CO', { maximumFractionDigits: 0 })}
                                                                </span>
                                                            </div>
                                                        )}

                                                        {form.watch("enableInsurance") && (rate.extraInsurance ?? 0) > 0 && (
                                                            <div className="flex justify-between items-center text-sm">
                                                                <div className="flex items-center gap-2">
                                                                    <div className="w-2 h-2 rounded-full" style={{ backgroundColor: '#22c55e' }}></div>
                                                                    <span className="text-gray-700 dark:text-gray-300 font-medium">Seguro</span>
                                                                </div>
                                                                <span className="text-gray-900 dark:text-gray-100 font-semibold" style={{ fontVariantNumeric: 'tabular-nums' }}>
                                                                    ${(rate.extraInsurance ?? 0).toLocaleString('es-CO', { maximumFractionDigits: 0 })}
                                                                </span>
                                                            </div>
                                                        )}

                                                        {(rate.codCarrierFee ?? 0) > 0 && (
                                                            <div className="flex justify-between items-center text-sm">
                                                                <div className="flex items-center gap-2">
                                                                    <div className="w-2 h-2 rounded-full" style={{ backgroundColor: businessColors.tertiary }}></div>
                                                                    <span className="text-gray-700 dark:text-gray-300 font-medium">Comisión carrier</span>
                                                                </div>
                                                                <span className="text-gray-900 dark:text-gray-100 font-semibold" style={{ fontVariantNumeric: 'tabular-nums' }}>
                                                                    ${(rate.codCarrierFee ?? 0).toLocaleString('es-CO', { maximumFractionDigits: 0 })}
                                                                </span>
                                                            </div>
                                                        )}
                                                    </div>

                                                    <div className="border border-gray-300 dark:border-gray-600 rounded-2xl p-4 bg-white dark:bg-gray-700/40">
                                                        <div className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wide mb-2">Precio {formatCarrierName(rate.carrier)}</div>
                                                        <div className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-3" style={{ fontVariantNumeric: 'tabular-nums' }}>
                                                            ${(rate.flete + (rate.codProbabilityMargin ?? 0) + (rate.minimumInsurance ?? 0) + (form.watch("enableInsurance") ? (rate.extraInsurance ?? 0) : 0)).toLocaleString('es-CO', { maximumFractionDigits: 0 })}
                                                        </div>
                                                        {(rate.codCarrierFee ?? 0) > 0 && (
                                                            <>
                                                                <div className="border-t border-gray-200 dark:border-gray-600 my-3"></div>
                                                                <div className="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wide mb-2">+ Comisión Carrier</div>
                                                                <div className="text-xl font-bold text-gray-900 dark:text-gray-100" style={{ fontVariantNumeric: 'tabular-nums' }}>
                                                                    ${((rate.flete + (rate.codProbabilityMargin ?? 0) + (rate.minimumInsurance ?? 0) + (form.watch("enableInsurance") ? (rate.extraInsurance ?? 0) : 0)) + (rate.codCarrierFee ?? 0)).toLocaleString('es-CO', { maximumFractionDigits: 0 })}
                                                                </div>
                                                            </>
                                                        )}
                                                    </div>
                                                </div>

                                                <div className="border-t border-gray-200 dark:border-gray-700 pt-4 flex justify-between items-center">
                                                    <div className="flex items-center gap-2 text-sm font-semibold bg-gray-50 dark:bg-gray-700/40 px-3 py-2 rounded-lg" style={{ color: businessColors.primary }}>
                                                        {rate.deliveryDays === 0 || rate.deliveryDays === 1 ? 'Mismo día' : `${rate.deliveryDays} día${rate.deliveryDays !== 1 ? 's' : ''}`}
                                                    </div>
                                                    <span className="text-xs font-bold uppercase tracking-wide cursor-pointer" style={{ color: businessColors.tertiary }}>
                                                        {isSelected ? '✓ Seleccionada' : 'Seleccionar'}
                                                    </span>
                                                </div>
                                            </div>
                                        );
                                    })}
                                </div>
                            </div>

                            <div className="flex gap-3 mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
                                <Button
                                    variant="secondary"
                                    onClick={() => setCurrentStep(1)}
                                    className="flex-1"
                                >
                                    Atrás
                                </Button>
                                <Button
                                    onClick={onClose}
                                    className="flex-1"
                                    style={{ backgroundColor: businessColors.tertiary }}
                                >
                                    Confirmar selección
                                </Button>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
