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
import '@/shared/ui/styles/shipment-modals.css';

const normalizeLocationName = (str: string) => {
    if (!str) return "";
    let s = str.normalize("NFD").replace(/[̀-ͯ]/g, "").toUpperCase().trim();
    s = s.replace(/,\s*D\.C\./g, "").replace(/\sD\.C\./g, "").replace(/\sDC\b/g, "").trim();
    return s;
};

const buildWarehouseAddress = (warehouse: Warehouse | null): string => {
    if (!warehouse) return "";
    const parts = [];
    if (warehouse.street) parts.push(warehouse.street);
    if (warehouse.suburb) parts.push(warehouse.suburb);
    if (warehouse.address) parts.push(warehouse.address);
    return parts.filter(p => p && p.trim()).join(", ") || "";
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
    };
    const trimmedName = carrierName.trim();
    if (carrierLogos[trimmedName]) return carrierLogos[trimmedName];
    return 'https://via.placeholder.com/56?text=' + encodeURIComponent(carrierName.substring(0, 3));
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

    useEffect(() => {
        if (isOpen) {
            getWarehousesAction({
                is_active: true,
                page: 1,
                page_size: 100,
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
                        <form onSubmit={form.handleSubmit(handleSubmit)} className="flex flex-col h-full overflow-hidden min-h-0">
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

                                            <Input
                                                compact
                                                label="Calle y Número *"
                                                {...form.register("destAddress")}
                                                placeholder="Carrera 46 # 93 - 45"
                                            />
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
                        <div className="flex flex-col h-full">
                            <div className="flex-1 overflow-y-auto">
                                <h3 className="font-semibold text-lg mb-4">Tarifas Disponibles</h3>
                                <div className="grid grid-cols-2 gap-4">
                                    {rates.map((rate) => (
                                        <div
                                            key={rate.idRate}
                                            className="border-2 rounded-xl p-4 cursor-pointer transition-all hover:shadow-lg"
                                            style={{
                                                borderColor: businessColors.tertiary,
                                                backgroundColor: businessColors.tertiary + '08'
                                            }}
                                        >
                                            <div className="flex items-center gap-3 mb-3">
                                                <img
                                                    src={getCarrierLogo(rate.carrier)}
                                                    alt={rate.carrier}
                                                    className="w-12 h-12 object-contain"
                                                    onError={(e) => {
                                                        e.currentTarget.style.display = 'none';
                                                    }}
                                                />
                                                <div>
                                                    <div className="font-bold text-sm">{rate.carrier}</div>
                                                    <div className="text-xs text-gray-500">{rate.product}</div>
                                                </div>
                                            </div>

                                            <div className="space-y-2">
                                                <div className="flex justify-between items-center">
                                                    <span className="text-xs text-gray-600">Días de entrega:</span>
                                                    <span className="font-semibold">{rate.deliveryDays} días</span>
                                                </div>
                                                <div className="flex justify-between items-center">
                                                    <span className="text-xs text-gray-600">Valor:</span>
                                                    <span className="font-bold text-lg" style={{ color: businessColors.primary }}>
                                                        ${rate.flete.toLocaleString()}
                                                    </span>
                                                </div>
                                                {rate.minimumInsurance && (
                                                    <div className="text-xs text-gray-500">
                                                        Seguro mínimo: ${rate.minimumInsurance.toLocaleString()}
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="flex gap-2 mt-4 pt-4 border-t border-gray-200 dark:border-gray-600">
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
                                    Cerrar
                                </Button>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
