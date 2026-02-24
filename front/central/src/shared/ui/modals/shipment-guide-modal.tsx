import { useState, useEffect, useRef } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Input, Button, Stepper } from "@/shared/ui";
import { ShipmentApiRepository } from "@/services/modules/shipments/infra/repository/api-repository";
import { EnvioClickQuoteRequest, EnvioClickRate } from "@/services/modules/shipments/domain/types";
import { Order } from "@/services/modules/orders/domain/types";
import { getWalletBalanceAction } from "@/services/modules/wallet/infra/actions";
import { getOriginAddressesAction, quoteShipmentAction, generateGuideAction } from "@/services/modules/shipments/infra/actions";
import { OriginAddress } from "@/services/modules/shipments/domain/types";
import danes from "@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json";

const normalizeLocationName = (str: string) => {
    if (!str) return "";
    let s = str.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toUpperCase().trim();
    // Remove variations of D.C. to improve matching
    s = s.replace(/,\s*D\.C\./g, "").replace(/\sD\.C\./g, "").replace(/\sDC\b/g, "").trim();
    return s;
};

const findDaneCode = (city: string, state: string) => {
    const targetCity = normalizeLocationName(city);
    const targetState = normalizeLocationName(state);

    if (!targetCity) return null;

    const entries = Object.entries(danes);

    // 1. Try exact match with city and state
    const exactMatch = entries.find(([_, data]: [string, any]) => {
        const dCity = normalizeLocationName(data.ciudad);
        const dState = normalizeLocationName(data.departamento);
        return dCity === targetCity && dState === targetState;
    });
    if (exactMatch) return exactMatch[0];

    // 2. Try match with city only
    const cityMatch = entries.find(([_, data]: [string, any]) => {
        const dCity = normalizeLocationName(data.ciudad);
        return dCity === targetCity;
    });
    if (cityMatch) return cityMatch[0];

    // 3. Try partial match
    const partialMatch = entries.find(([_, data]: [string, any]) => {
        const dCity = normalizeLocationName(data.ciudad);
        return dCity.includes(targetCity) || targetCity.includes(dCity);
    });
    if (partialMatch) return partialMatch[0];

    return null;
};

interface ShipmentGuideModalProps {
    isOpen: boolean;
    onClose: () => void;
    order?: Order;
    onGuideGenerated?: (trackingNumber: string) => void;
    recommendedCarrier?: string;
}

// Step 1: Origin/Destination/Package Schema
const step1Schema = z.object({
    originDaneCode: z.string().min(8, "Código DANE de origen requerido"),
    originAddress: z.string().min(2, "Dirección de origen requerida").max(50),
    destDaneCode: z.string().min(8, "Código DANE de destino requerido"),
    destAddress: z.string().min(8, "Dirección de destino requerida").max(50),
    weight: z.number().min(1).max(1000),
    height: z.number().min(1).max(300),
    width: z.number().min(1).max(300),
    length: z.number().min(1).max(300),
    description: z.string().min(3).max(25),
    contentValue: z.number().min(0).max(3000000),
    codValue: z.number().min(0).max(3000000).optional(),
    includeGuideCost: z.boolean(),
    codPaymentMethod: z.enum(["cash", "data_phone"]),
});

// Step 3: Detailed Contact Info Schema
const step3Schema = z.object({
    originCompany: z.string().min(2).max(28),
    originFirstName: z.string().min(2).max(14),
    originLastName: z.string().min(2).max(14),
    originEmail: z.string().email().min(8).max(60),
    originPhone: z.string().length(10),
    originSuburb: z.string().min(2).max(30),
    originCrossStreet: z.string().min(2).max(35),
    originReference: z.string().min(2).max(25),
    destCompany: z.string().min(2).max(28),
    destFirstName: z.string().min(2).max(14),
    destLastName: z.string().min(2).max(14),
    destEmail: z.string().email().min(8).max(60),
    destPhone: z.string().length(10),
    destSuburb: z.string().min(2).max(30),
    destCrossStreet: z.string().min(2).max(35),
    destReference: z.string().min(2).max(25),
    requestPickup: z.boolean(),
    insurance: z.boolean(),
    myShipmentReference: z.string().min(2).max(28),
    external_order_id: z.string().min(1).max(28).optional(),
});

type Step1Values = z.infer<typeof step1Schema>;
type Step3Values = z.infer<typeof step3Schema>;

const STEPS = [
    { id: 1, label: "Origen y Destino" },
    { id: 2, label: "Cotización" },
    { id: 3, label: "Detalles" },
    { id: 4, label: "Pago" },
];

export default function ShipmentGuideModal({ isOpen, onClose, order, onGuideGenerated, recommendedCarrier }: ShipmentGuideModalProps) {
    const [currentStep, setCurrentStep] = useState(1);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    // Step 1 data
    const [step1Data, setStep1Data] = useState<Step1Values | null>(null);

    // Step 2 data (quotes)
    const [rates, setRates] = useState<EnvioClickRate[]>([]);
    const [selectedRate, setSelectedRate] = useState<EnvioClickRate | null>(null);

    // Step 3 data
    const [step3Data, setStep3Data] = useState<Step3Values | null>(null);

    // Step 4 data
    const [walletBalance, setWalletBalance] = useState<number | null>(null);
    const [originAddresses, setOriginAddresses] = useState<OriginAddress[]>([]);
    const [generatedPdfUrl, setGeneratedPdfUrl] = useState<string | null>(null);
    const [trackingNumber, setTrackingNumber] = useState<string | null>(null);

    // DANE search states
    const [originSearch, setOriginSearch] = useState("");
    const [destSearch, setDestSearch] = useState("");
    const [showOriginResults, setShowOriginResults] = useState(false);
    const [showDestResults, setShowDestResults] = useState(false);
    const originRef = useRef<HTMLDivElement>(null);
    const destRef = useRef<HTMLDivElement>(null);

    // const repo = new ShipmentApiRepository(); // Eliminado para usar Server Actions

    // DANE options
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

    // Step 1 Form
    const step1Form = useForm<Step1Values>({
        resolver: zodResolver(step1Schema),
        defaultValues: {
            originDaneCode: "11001000",
            originAddress: "",
            destDaneCode: "11001000",
            destAddress: "",
            weight: 1,
            height: 10,
            width: 10,
            length: 10,
            description: "E-commerce Order",
            contentValue: 0,
            codValue: 0,
            includeGuideCost: false,
            codPaymentMethod: "cash",
        },
    });

    // Step 3 Form
    const step3Form = useForm<Step3Values>({
        resolver: zodResolver(step3Schema),
        defaultValues: {
            originCompany: "Mi Empresa",
            originFirstName: "",
            originLastName: "",
            originEmail: "",
            originPhone: "",
            originSuburb: "",
            originCrossStreet: "",
            originReference: "",
            destCompany: "NA",
            destFirstName: "",
            destLastName: "",
            destEmail: "",
            destPhone: "",
            destSuburb: "",
            destCrossStreet: "",
            destReference: "",
            requestPickup: false,
            insurance: true,
            myShipmentReference: "",
            external_order_id: "",
        },
    });

    // Fetch initial data on open
    useEffect(() => {
        if (isOpen) {
            getWalletBalanceAction().then(res => {
                if (res.success && res.data) setWalletBalance(res.data.Balance);
            });
            getOriginAddressesAction().then(res => {
                if (res.success && res.data) {
                    setOriginAddresses(res.data);
                    // Si hay una predeterminada, seleccionarla automáticamente
                    const defaultAddr = res.data.find(a => a.is_default);
                    if (defaultAddr) {
                        handleOriginAddressSelect(defaultAddr);
                    }
                }
            });
        }
    }, [isOpen]);

    const handleOriginAddressSelect = (addr: OriginAddress) => {
        // Step 1
        step1Form.setValue("originDaneCode", addr.city_dane_code);
        step1Form.setValue("originAddress", addr.street);
        setOriginSearch(`${addr.city} (${addr.state})`);

        // Step 3
        step3Form.setValue("originCompany", addr.company);
        step3Form.setValue("originFirstName", addr.first_name);
        step3Form.setValue("originLastName", addr.last_name);
        step3Form.setValue("originEmail", addr.email);
        step3Form.setValue("originPhone", addr.phone);
        step3Form.setValue("originSuburb", addr.suburb || "");
        step3Form.setValue("originCrossStreet", addr.street);
        step3Form.setValue("originReference", ""); // Opcional
    };

    // Pre-fill from order
    useEffect(() => {
        if (isOpen && order) {
            // Step 1
            step1Form.setValue("contentValue", order.total_amount);
            step1Form.setValue("codValue", order.total_amount);
            step1Form.setValue("description", `Order ${order.order_number}`);
            step1Form.setValue("destAddress", order.shipping_street);

            if (order.weight && order.weight > 0) {
                step1Form.setValue("weight", order.weight);
                step1Form.setValue("height", order.height || 10);
                step1Form.setValue("width", order.width || 10);
                step1Form.setValue("length", order.length || 10);
            }

            // Try to find DANE code by city
            const mappedDane = findDaneCode(order.shipping_city || "", order.shipping_state || "");
            const finalDane = mappedDane || "11001000"; // Fallback to Bogota

            step1Form.setValue("destDaneCode", finalDane);
            const cityData = danes[finalDane as keyof typeof danes];
            if (cityData) {
                setDestSearch(`${(cityData as any).ciudad} (${(cityData as any).departamento})`);
            }

            // Step 3
            step3Form.setValue("destCompany", order.customer_name);
            step3Form.setValue("destFirstName", order.customer_name.split(" ")[0] || "");
            step3Form.setValue("destLastName", order.customer_name.split(" ").slice(1).join(" ") || ".");
            step3Form.setValue("destEmail", order.customer_email);
            step3Form.setValue("destPhone", order.customer_phone);
            step3Form.setValue("destSuburb", order.shipping_state || "");
            step3Form.setValue("myShipmentReference", "Orden " + (order.internal_number || order.order_number));
            step3Form.setValue("external_order_id", order.order_number);
        }
    }, [isOpen, order]);

    // Click outside to close DANE dropdowns
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

    // Reset on close
    useEffect(() => {
        if (!isOpen) {
            setCurrentStep(1);
            setStep1Data(null);
            setRates([]);
            setSelectedRate(null);
            setStep3Data(null);
            setGeneratedPdfUrl(null);
            setTrackingNumber(null);
            setError(null);
            setSuccess(null);
            step1Form.reset();
            step3Form.reset();
        }
    }, [isOpen]);

    // Step 1: Quote
    const handleStep1Submit = async (data: Step1Values) => {
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
                codValue: data.codValue,
                includeGuideCost: data.includeGuideCost,
                codPaymentMethod: data.codPaymentMethod,
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
            if (response.success && response.data?.data?.rates && response.data.data.rates.length > 0) {
                setRates(response.data.data.rates);
                setStep1Data(data);
                setCurrentStep(2);
            } else {
                setError("No se encontraron tarifas disponibles");
            }
        } catch (err: any) {
            setError(err.message || "Error al cotizar envío");
        } finally {
            setLoading(false);
        }
    };

    // Step 2: Select Rate
    const handleRateSelection = (rate: EnvioClickRate) => {
        setSelectedRate(rate);
        setCurrentStep(3);
    };

    // Step 3: Details
    const handleStep3Submit = async (data: Step3Values) => {
        setStep3Data(data);
        setCurrentStep(4);
    };

    // Step 4: Generate Guide
    const handleFinalGenerate = async () => {
        if (!step1Data || !selectedRate || !step3Data) {
            setError("Faltan datos para generar la guía");
            return;
        }

        // Check wallet balance
        const totalCost = selectedRate.flete + (selectedRate.minimumInsurance ?? 0) + (selectedRate.extraInsurance ?? 0);
        if (walletBalance !== null && walletBalance < totalCost) {
            setError(`Saldo insuficiente. Necesitas $${totalCost.toLocaleString()} pero tienes $${walletBalance.toLocaleString()}`);
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const generatePayload: EnvioClickQuoteRequest = {
                idRate: selectedRate.idRate,
                myShipmentReference: step3Data.myShipmentReference,
                external_order_id: step3Data.external_order_id,
                order_uuid: order?.id,
                requestPickup: step3Data.requestPickup,
                pickupDate: new Date().toISOString().split("T")[0],
                insurance: step3Data.insurance,
                description: step1Data.description,
                contentValue: step1Data.contentValue,
                codValue: step1Data.codValue,
                includeGuideCost: step1Data.includeGuideCost,
                codPaymentMethod: step1Data.codPaymentMethod,
                totalCost: totalCost,
                packages: [{
                    weight: step1Data.weight,
                    height: step1Data.height,
                    width: step1Data.width,
                    length: step1Data.length,
                }],
                origin: {
                    daneCode: step1Data.originDaneCode,
                    address: step1Data.originAddress,
                    company: step3Data.originCompany,
                    firstName: step3Data.originFirstName,
                    lastName: step3Data.originLastName,
                    email: step3Data.originEmail,
                    phone: step3Data.originPhone,
                    suburb: step3Data.originSuburb,
                    crossStreet: step3Data.originCrossStreet,
                    reference: step3Data.originReference,
                },
                destination: {
                    daneCode: step1Data.destDaneCode,
                    address: step1Data.destAddress,
                    company: step3Data.destCompany,
                    firstName: step3Data.destFirstName,
                    lastName: step3Data.destLastName,
                    email: step3Data.destEmail,
                    phone: step3Data.destPhone,
                    suburb: step3Data.destSuburb,
                    crossStreet: step3Data.destCrossStreet,
                    reference: step3Data.destReference,
                },
            };

            const response = await generateGuideAction(generatePayload);
            if (response.success && response.data?.data) {
                setGeneratedPdfUrl(response.data.data.url);
                setTrackingNumber(response.data.data.tracker);
                if (onGuideGenerated && response.data.data.tracker) {
                    onGuideGenerated(response.data.data.tracker);
                }
                setSuccess("¡Guía generada exitosamente!");
            }
        } catch (err: any) {
            setError(err.message || "Error al generar guía");
        } finally {
            setLoading(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
                {/* Header */}
                <div className="sticky top-0 bg-white border-b px-6 py-4 z-10">
                    <div className="flex justify-between items-center mb-4">
                        <h2 className="text-2xl font-bold text-gray-800">Generar Guía de Envío</h2>
                        <button
                            onClick={onClose}
                            className="text-gray-500 hover:text-gray-700 text-2xl"
                        >
                            ×
                        </button>
                    </div>
                    <Stepper steps={STEPS} currentStep={currentStep} />
                </div>

                {/* Content */}
                <div className="p-6">
                    {error && (
                        <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
                            {error}
                        </div>
                    )}

                    {success && (
                        <div className="mb-4 p-4 bg-green-50 border border-green-200 rounded-lg text-green-700">
                            {success}
                        </div>
                    )}

                    {/* Step 1: Origin/Destination/Package */}
                    {currentStep === 1 && (
                        <form onSubmit={step1Form.handleSubmit(handleStep1Submit)} className="space-y-6">
                            <div className="grid grid-cols-2 gap-6">
                                {/* Origin */}
                                <div className="space-y-4">
                                    <div className="flex items-center justify-between">
                                        <h3 className="font-semibold text-lg text-gray-700">Origen</h3>
                                        {originAddresses.length > 0 && (
                                            <select
                                                className="text-xs border border-gray-200 rounded px-2 py-1 bg-white focus:outline-none focus:ring-1 focus:ring-orange-500"
                                                onChange={(e) => {
                                                    const addr = originAddresses.find(a => a.id === parseInt(e.target.value));
                                                    if (addr) handleOriginAddressSelect(addr);
                                                }}
                                                defaultValue=""
                                            >
                                                <option value="" disabled>Mis direcciones...</option>
                                                {originAddresses.map(a => (
                                                    <option key={a.id} value={a.id}>{a.alias}</option>
                                                ))}
                                            </select>
                                        )}
                                    </div>

                                    <div ref={originRef} className="relative">
                                        <label className="block text-sm font-medium text-gray-700 mb-1">
                                            Ciudad remitente *
                                        </label>
                                        <input
                                            type="text"
                                            value={originSearch}
                                            onChange={(e) => {
                                                setOriginSearch(e.target.value);
                                                setShowOriginResults(true);
                                            }}
                                            onFocus={() => setShowOriginResults(true)}
                                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
                                            placeholder="Buscar ciudad..."
                                        />
                                        {showOriginResults && filteredOriginOptions.length > 0 && (
                                            <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-y-auto">
                                                {filteredOriginOptions.slice(0, 50).map((opt) => (
                                                    <div
                                                        key={opt.value}
                                                        onClick={() => {
                                                            step1Form.setValue("originDaneCode", opt.value);
                                                            setOriginSearch(opt.label);
                                                            setShowOriginResults(false);
                                                        }}
                                                        className="px-3 py-2 hover:bg-gray-100 cursor-pointer"
                                                    >
                                                        {opt.label}
                                                    </div>
                                                ))}
                                            </div>
                                        )}
                                    </div>

                                    <Input
                                        label="Calle y Número *"
                                        {...step1Form.register("originAddress")}
                                        error={step1Form.formState.errors.originAddress?.message}
                                        placeholder="Calle 98 62-37"
                                    />
                                </div>

                                {/* Destination */}
                                <div className="space-y-4">
                                    <h3 className="font-semibold text-lg text-gray-700">Destino</h3>

                                    <div ref={destRef} className="relative">
                                        <label className="block text-sm font-medium text-gray-700 mb-1">
                                            Ciudad destinatario *
                                        </label>
                                        <input
                                            type="text"
                                            value={destSearch}
                                            onChange={(e) => {
                                                setDestSearch(e.target.value);
                                                setShowDestResults(true);
                                            }}
                                            onFocus={() => setShowDestResults(true)}
                                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
                                            placeholder="Buscar ciudad..."
                                        />
                                        {showDestResults && filteredDestOptions.length > 0 && (
                                            <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-y-auto">
                                                {filteredDestOptions.slice(0, 50).map((opt) => (
                                                    <div
                                                        key={opt.value}
                                                        onClick={() => {
                                                            step1Form.setValue("destDaneCode", opt.value);
                                                            setDestSearch(opt.label);
                                                            setShowDestResults(false);
                                                        }}
                                                        className="px-3 py-2 hover:bg-gray-100 cursor-pointer"
                                                    >
                                                        {opt.label}
                                                    </div>
                                                ))}
                                            </div>
                                        )}
                                    </div>

                                    <Input
                                        label="Calle y Número *"
                                        {...step1Form.register("destAddress")}
                                        error={step1Form.formState.errors.destAddress?.message}
                                        placeholder="Carrera 46 # 93 - 45"
                                    />
                                </div>
                            </div>

                            {/* Package Details */}
                            <div className="border-t pt-6">
                                <h3 className="font-semibold text-lg text-gray-700 mb-4">Características del paquete</h3>
                                <div className="grid grid-cols-4 gap-4">
                                    <Input
                                        label="Peso (kg) *"
                                        type="number"
                                        step="0.1"
                                        {...step1Form.register("weight", { valueAsNumber: true })}
                                        error={step1Form.formState.errors.weight?.message}
                                    />
                                    <Input
                                        label="Alto (cm) *"
                                        type="number"
                                        {...step1Form.register("height", { valueAsNumber: true })}
                                        error={step1Form.formState.errors.height?.message}
                                    />
                                    <Input
                                        label="Ancho (cm) *"
                                        type="number"
                                        {...step1Form.register("width", { valueAsNumber: true })}
                                        error={step1Form.formState.errors.width?.message}
                                    />
                                    <Input
                                        label="Largo (cm) *"
                                        type="number"
                                        {...step1Form.register("length", { valueAsNumber: true })}
                                        error={step1Form.formState.errors.length?.message}
                                    />
                                </div>
                            </div>

                            {/* Additional Info */}
                            <div className="grid grid-cols-2 gap-4">
                                <Input
                                    label="Descripción *"
                                    {...step1Form.register("description")}
                                    error={step1Form.formState.errors.description?.message}
                                    placeholder="descripción"
                                />
                                <Input
                                    label="Valor factura declarado *"
                                    type="number"
                                    {...step1Form.register("contentValue", { valueAsNumber: true })}
                                    error={step1Form.formState.errors.contentValue?.message}
                                />
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="flex items-center space-x-2">
                                        <input
                                            type="checkbox"
                                            {...step1Form.register("includeGuideCost")}
                                            className="rounded"
                                        />
                                        <span className="text-sm">Incluir costo de guía en COD</span>
                                    </label>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        Método de pago COD
                                    </label>
                                    <select
                                        {...step1Form.register("codPaymentMethod")}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
                                    >
                                        <option value="cash">Efectivo</option>
                                        <option value="data_phone">Datáfono</option>
                                    </select>
                                    {step1Form.formState.errors.codPaymentMethod?.message && (
                                        <p className="text-sm text-red-500 mt-1">
                                            {step1Form.formState.errors.codPaymentMethod.message}
                                        </p>
                                    )}
                                </div>
                            </div>

                            <div className="flex justify-end">
                                <Button type="submit" disabled={loading}>
                                    {loading ? "Cotizando..." : "Siguiente"}
                                </Button>
                            </div>
                        </form>
                    )}

                    {/* Step 2: Quote Selection */}
                    {currentStep === 2 && (
                        <div className="space-y-4">
                            <h3 className="font-semibold text-lg text-gray-700 mb-4">
                                Filtra por servicio / Transportadora
                            </h3>
                            <p className="text-sm text-gray-600 mb-4">Todos los precios incluyen IVA</p>

                            {rates.map((rate) => {
                                const totalCost = rate.flete + (rate.minimumInsurance ?? 0) + (rate.extraInsurance ?? 0);
                                const isCOD = rate.cod;

                                return (
                                    <div
                                        key={rate.idRate}
                                        onClick={() => handleRateSelection(rate)}
                                        className="border rounded-lg p-4 hover:border-orange-500 cursor-pointer transition-colors"
                                    >
                                        <div className="flex items-center justify-between">
                                            <div className="flex items-center space-x-4">
                                                <div className="w-16 h-16 bg-gray-100 rounded flex items-center justify-center">
                                                    <span className="font-bold text-xs">{rate.carrier}</span>
                                                </div>
                                                <div>
                                                    <div className="font-semibold">{rate.carrier}</div>
                                                    <div className="text-sm text-gray-600">{rate.product}</div>
                                                    <div className="text-xs text-gray-500">
                                                        Entrega aprox. {rate.deliveryDays} días
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="text-right">
                                                <div className="text-2xl font-bold text-orange-600">
                                                    ${totalCost.toLocaleString()} COP
                                                </div>
                                                <div className="text-xs text-gray-500 mt-1">
                                                    Entrega: {rate.deliveryDays} días
                                                </div>
                                                {isCOD && (
                                                    <div className="text-xs text-blue-600 mt-1">
                                                        COD disponible
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                );
                            })}

                            <div className="flex justify-between mt-6">
                                <Button
                                    variant="outline"
                                    onClick={() => setCurrentStep(1)}
                                >
                                    Atrás
                                </Button>
                            </div>
                        </div>
                    )}

                    {/* Step 3: Details */}
                    {currentStep === 3 && (
                        <form onSubmit={step3Form.handleSubmit(handleStep3Submit)} className="space-y-6">
                            {/* Origin Details */}
                            <div>
                                <h3 className="font-semibold text-lg text-gray-700 mb-4">Dirección - Remitente</h3>
                                <div className="grid grid-cols-3 gap-4">
                                    <Input
                                        label="Calle *"
                                        {...step3Form.register("originCrossStreet")}
                                        error={step3Form.formState.errors.originCrossStreet?.message}
                                        placeholder="calle 75 sur n 42-97"
                                    />
                                    <Input
                                        label="Edificio/Interior/Apto *"
                                        {...step3Form.register("originReference")}
                                        error={step3Form.formState.errors.originReference?.message}
                                        placeholder="apt 801"
                                    />
                                    <Input
                                        label="Barrio *"
                                        {...step3Form.register("originSuburb")}
                                        error={step3Form.formState.errors.originSuburb?.message}
                                        placeholder="sector Aves María"
                                    />
                                </div>

                                <h4 className="font-medium text-gray-700 mt-6 mb-3">Referencias - Empresa</h4>
                                <Input
                                    label="Empresa"
                                    {...step3Form.register("originCompany")}
                                    error={step3Form.formState.errors.originCompany?.message}
                                    placeholder="ProbabilityIA"
                                />

                                <h4 className="font-medium text-gray-700 mt-6 mb-3">Datos de contacto</h4>
                                <div className="grid grid-cols-2 gap-4">
                                    <Input
                                        label="Nombre *"
                                        {...step3Form.register("originFirstName")}
                                        error={step3Form.formState.errors.originFirstName?.message}
                                        placeholder="Luisa"
                                    />
                                    <Input
                                        label="Apellido *"
                                        {...step3Form.register("originLastName")}
                                        error={step3Form.formState.errors.originLastName?.message}
                                        placeholder="Muñoz"
                                    />
                                    <Input
                                        label="Teléfono *"
                                        {...step3Form.register("originPhone")}
                                        error={step3Form.formState.errors.originPhone?.message}
                                        placeholder="3224098631"
                                    />
                                    <Input
                                        label="Correo *"
                                        type="email"
                                        {...step3Form.register("originEmail")}
                                        error={step3Form.formState.errors.originEmail?.message}
                                        placeholder="probabilitysa@gmail.com"
                                    />
                                </div>
                            </div>

                            {/* Destination Details */}
                            <div className="border-t pt-6">
                                <h3 className="font-semibold text-lg text-gray-700 mb-4">Destinatario</h3>
                                <div className="grid grid-cols-3 gap-4">
                                    <Input
                                        label="Calle *"
                                        {...step3Form.register("destCrossStreet")}
                                        error={step3Form.formState.errors.destCrossStreet?.message}
                                        placeholder="calle 75 sur n 42-97"
                                    />
                                    <Input
                                        label="Edificio/Interior/Apto *"
                                        {...step3Form.register("destReference")}
                                        error={step3Form.formState.errors.destReference?.message}
                                        placeholder="apt 801"
                                    />
                                    <Input
                                        label="Barrio *"
                                        {...step3Form.register("destSuburb")}
                                        error={step3Form.formState.errors.destSuburb?.message}
                                        placeholder="sector Aves María"
                                    />
                                </div>

                                <h4 className="font-medium text-gray-700 mt-6 mb-3">Referencias - Empresa</h4>
                                <Input
                                    label="Empresa"
                                    {...step3Form.register("destCompany")}
                                    error={step3Form.formState.errors.destCompany?.message}
                                    placeholder="ProbabilityIA"
                                />

                                <h4 className="font-medium text-gray-700 mt-6 mb-3">Datos de contacto</h4>
                                <div className="grid grid-cols-2 gap-4">
                                    <Input
                                        label="Nombre *"
                                        {...step3Form.register("destFirstName")}
                                        error={step3Form.formState.errors.destFirstName?.message}
                                        placeholder="Luisa"
                                    />
                                    <Input
                                        label="Apellido *"
                                        {...step3Form.register("destLastName")}
                                        error={step3Form.formState.errors.destLastName?.message}
                                        placeholder="Muñoz"
                                    />
                                    <Input
                                        label="Teléfono *"
                                        {...step3Form.register("destPhone")}
                                        error={step3Form.formState.errors.destPhone?.message}
                                        placeholder="3224098631"
                                    />
                                    <Input
                                        label="Correo *"
                                        type="email"
                                        {...step3Form.register("destEmail")}
                                        error={step3Form.formState.errors.destEmail?.message}
                                        placeholder="probabilitysa@gmail.com"
                                    />
                                </div>
                            </div>

                            {/* Additional Options */}
                            <div className="border-t pt-6">
                                <div className="grid grid-cols-2 gap-4">
                                    <Input
                                        label="Mi referencia de envío"
                                        {...step3Form.register("myShipmentReference")}
                                        error={step3Form.formState.errors.myShipmentReference?.message}
                                        placeholder="Orden 5649"
                                    />
                                    <Input
                                        label="Número de orden externo"
                                        {...step3Form.register("external_order_id")}
                                        error={step3Form.formState.errors.external_order_id?.message}
                                        placeholder="ORD345678"
                                    />
                                </div>

                                <div className="mt-4 space-y-2">
                                    <label className="flex items-center space-x-2">
                                        <input
                                            type="checkbox"
                                            {...step3Form.register("requestPickup")}
                                            className="rounded"
                                        />
                                        <span className="text-sm">Solicitar recolección</span>
                                    </label>
                                    <label className="flex items-center space-x-2">
                                        <input
                                            type="checkbox"
                                            {...step3Form.register("insurance")}
                                            className="rounded"
                                        />
                                        <span className="text-sm">Asegurar envío</span>
                                    </label>
                                </div>
                            </div>

                            <div className="flex justify-between mt-6">
                                <Button
                                    variant="outline"
                                    onClick={() => setCurrentStep(2)}
                                    type="button"
                                >
                                    Atrás
                                </Button>
                                <Button type="submit">
                                    Siguiente
                                </Button>
                            </div>
                        </form>
                    )}

                    {/* Step 4: Payment & Confirmation */}
                    {currentStep === 4 && selectedRate && (
                        <div className="space-y-6">
                            <h3 className="font-semibold text-lg text-gray-700">Resumen de tu envío</h3>

                            <div className="bg-gray-50 p-4 rounded-lg">
                                <div className="flex items-center justify-between mb-4">
                                    <div className="flex items-center space-x-2">
                                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                                        </svg>
                                        <span className="font-medium">1 Envíos</span>
                                    </div>
                                    <div className="text-right">
                                        <div className="text-sm text-gray-600">TOTAL:</div>
                                        <div className="text-2xl font-bold text-orange-600">
                                            ${(selectedRate.flete + (selectedRate.minimumInsurance ?? 0)).toLocaleString()}
                                        </div>
                                    </div>
                                </div>

                                <div className="border-t pt-4">
                                    <div className="text-sm text-gray-600 mb-2">Carrier: {selectedRate.carrier}</div>
                                    <div className="text-sm text-gray-600">Producto: {selectedRate.product}</div>
                                </div>
                            </div>

                            <div>
                                <h4 className="font-medium text-gray-700 mb-3">Selecciona tu método de pago</h4>
                                <div className="grid grid-cols-2 gap-4">
                                    <div className="border-2 border-orange-500 rounded-lg p-4 bg-orange-50">
                                        <div className="flex items-center justify-center mb-2">
                                            <svg className="w-8 h-8 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
                                            </svg>
                                        </div>
                                        <div className="text-center font-semibold">Monedero</div>
                                        <div className="text-center text-sm text-gray-600">
                                            ${walletBalance?.toLocaleString() || 0}
                                        </div>
                                    </div>
                                </div>
                            </div>

                            {generatedPdfUrl && trackingNumber ? (
                                <div className="bg-green-50 border border-green-200 rounded-lg p-6">
                                    <h4 className="font-semibold text-green-800 mb-4">¡Guía generada exitosamente!</h4>
                                    <div className="space-y-2">
                                        <p className="text-sm"><strong>Tracking:</strong> {trackingNumber}</p>
                                        <a
                                            href={generatedPdfUrl}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="inline-block px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
                                        >
                                            Descargar Guía PDF
                                        </a>
                                    </div>
                                </div>
                            ) : (
                                <div className="flex justify-between mt-6">
                                    <Button
                                        variant="outline"
                                        onClick={() => setCurrentStep(3)}
                                        disabled={loading}
                                    >
                                        Atrás
                                    </Button>
                                    <Button
                                        onClick={handleFinalGenerate}
                                        disabled={loading}
                                        className="bg-green-600 hover:bg-green-700"
                                    >
                                        {loading ? "Generando..." : "Pagar guías"}
                                    </Button>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
