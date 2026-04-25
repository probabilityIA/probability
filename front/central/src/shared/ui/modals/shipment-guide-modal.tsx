'use client';

import { useState, useEffect, useRef } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Input, Button, Stepper } from "@/shared/ui";
import { ShipmentApiRepository } from "@/services/modules/shipments/infra/repository/api-repository";
import { EnvioClickQuoteRequest, EnvioClickRate } from "@/services/modules/shipments/domain/types";
import { Order } from "@/services/modules/orders/domain/types";
import { getWalletBalanceAction } from "@/services/modules/wallet/infra/actions";
import { quoteShipmentAction, generateGuideAction } from "@/services/modules/shipments/infra/actions";
import { getWarehousesAction } from "@/services/modules/warehouses/infra/actions";
import { Warehouse } from "@/services/modules/warehouses/domain/types";
import danes from "@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json";
import { useShipmentSSE } from "@/services/modules/shipments/ui/hooks/useShipmentSSE";
import { usePermissions } from "@/shared/contexts/permissions-context";
import { getActionError } from '@/shared/utils/action-result';
import { CarrierOfficeSelector } from "@/services/modules/shipments/ui/components/CarrierOfficeSelector";

const normalizeLocationName = (str: string) => {
    if (!str) return "";
    let s = str.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toUpperCase().trim();
    // Remove variations of D.C. to improve matching
    s = s.replace(/,\s*D\.C\./g, "").replace(/\sD\.C\./g, "").replace(/\sDC\b/g, "").trim();
    return s;
};

const normalizeString = (str: string) =>
    str.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toUpperCase().trim();

const getCarrierLogoSize = (carrierName: string): { container: string; image: string } => {
    const largeLogoCarriers = ['COORDINADORA', '99MINUTOS', 'PIBOX', 'DEPRISA'];
    const normalizedCarrier = normalizeLocationName(carrierName);

    if (largeLogoCarriers.includes(normalizedCarrier)) {
        return { container: 'w-24 h-24', image: 'w-20 h-20' };
    }

    return { container: 'w-20 h-20', image: 'w-18 h-18' };
};

const getCarrierLogo = (carrierName: string): string => {
    const carrierLogos: { [key: string]: string } = {
        'SERVIENTREGA': 'https://i.revistapym.com.co/old/2021/09/WhatsApp-Image-2021-09-25-at-1.08.55-PM.jpeg?w=400&r=1_1',
        'COORDINADORA': 'https://olartemoure.com/wp-content/uploads/2023/05/coordinadora-logo.png',
        'DHLEXPRESS': 'https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png',
        'DHL': 'https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png',
        'FEDEX': 'https://upload.wikimedia.org/wikipedia/commons/thumb/9/9d/FedEx_Express.svg/960px-FedEx_Express.svg.png',
        'INTERRAPIDISIMO': 'https://interrapidisimo.com/wp-content/uploads/Logo-Inter-Rapidisimo-Vv-400x431-1.png',
        '472LOGISTICA': 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTnDF0ozRHf3s5BPqLsr7Vg-X8JRzECvFvwBQ&s',
        'SPEED': 'https://speedcargopa.com/wp-content/uploads/2021/03/Logo-mejorado-transparencia.png',
        'SPEEDCARGO': 'https://speedcargopa.com/wp-content/uploads/2021/03/Logo-mejorado-transparencia.png',
        'ENVIA': 'https://images.seeklogo.com/logo-png/31/1/envia-mensajeria-logo-png_seeklogo-311137.png',
        'PIBOX': 'https://play-lh.googleusercontent.com/r_zPLkaHZK4Odu1yp6dqIdUnVAmIiLc3s18F9gUFqcz8IyHqCb_aGHP4iJSesXxnUyU',
        'TCC': 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Logo_TCC.svg/1280px-Logo_TCC.svg.png',
        'TRANSPORTADORADECARACOLOMBIA': 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Logo_TCC.svg/1280px-Logo_TCC.svg.png',
        '99MINUTOS': 'https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Logo-99minutos.svg/3840px-Logo-99minutos.svg.png',
        'DEPRISA': 'https://www.specialcolombia.com/wp-content/uploads/2023/05/Logo_azul_concepto_azul-deprisa.png',
    };

    const normalizedCarrier = normalizeLocationName(carrierName);
    return carrierLogos[normalizedCarrier] || 'https://via.placeholder.com/56?text=' + encodeURIComponent(carrierName.substring(0, 3));
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
    onGuideGenerated?: (data: { tracking_number: string; carrier?: string; label_url?: string }) => void;
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
    insurance: z.boolean(),
    codPaymentMethod: z.enum(["cash", "data_phone"]),
});

// Step 3: Detailed Contact Info Schema
const step3Schema = z.object({
    originCompany: z.string().min(2, "Min 2 caracteres").max(28, "Max 28 caracteres"),
    originFirstName: z.string().min(2, "Min 2 caracteres").max(14, "Max 14 caracteres"),
    originLastName: z.string().min(2, "Min 2 caracteres").max(14, "Max 14 caracteres"),
    originEmail: z.string().email("Email invalido").min(8, "Min 8 caracteres").max(60, "Max 60 caracteres"),
    originPhone: z.string().length(10, "Debe tener 10 digitos"),
    originSuburb: z.string().min(2, "Min 2 caracteres").max(30, "Max 30 caracteres"),
    originCrossStreet: z.string().min(2, "Min 2 caracteres").max(35, "Max 35 caracteres"),
    originReference: z.string().min(2, "Min 2 caracteres").max(25, "Max 25 caracteres"),
    destCompany: z.string().min(2, "Min 2 caracteres").max(28, "Max 28 caracteres").optional(),
    destFirstName: z.string().min(2, "Min 2 caracteres").max(14, "Max 14 caracteres"),
    destLastName: z.string().min(2, "Min 2 caracteres").max(14, "Max 14 caracteres"),
    destEmail: z.string().email("Email invalido").min(8, "Min 8 caracteres").max(60, "Max 60 caracteres"),
    destPhone: z.string().length(10, "Debe tener 10 digitos"),
    destSuburb: z.string().min(2, "Min 2 caracteres").max(30, "Max 30 caracteres").optional(),
    destCrossStreet: z.string().min(2, "Min 2 caracteres").max(35, "Max 35 caracteres"),
    destReference: z.string().min(2, "Min 2 caracteres").max(25, "Max 25 caracteres").optional(),
    requestPickup: z.boolean(),
    myShipmentReference: z.string().min(2, "Min 2 caracteres").max(28, "Max 28 caracteres"),
    external_order_id: z.string().min(1, "Requerido").max(28, "Max 28 caracteres").optional(),
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

    // Async quote tracking: saves correlation_id while waiting for SSE response
    const [pendingCorrelationId, setPendingCorrelationId] = useState<string | null>(null);
    const pendingStep1DataRef = useRef<Step1Values | null>(null);

    // Async guide generation tracking
    const [pendingGuideCorrelationId, setPendingGuideCorrelationId] = useState<string | null>(null);
    const [guideGenerationRequested, setGuideGenerationRequested] = useState(false);

    // Get businessId from permissions for SSE connection
    const { permissions, isSuperAdmin } = usePermissions();
    const businessId = permissions?.business_id || 0;
    // For wallet queries: super admin acts on behalf of the order's business
    const effectiveBusinessId = isSuperAdmin ? (order?.business_id ?? 0) : 0;

    // Step 1 data
    const [step1Data, setStep1Data] = useState<Step1Values | null>(null);

    // Step 2 data (quotes)
    const [rates, setRates] = useState<EnvioClickRate[]>([]);
    const [selectedRate, setSelectedRate] = useState<EnvioClickRate | null>(null);

    // Step 3 data
    const [step3Data, setStep3Data] = useState<Step3Values | null>(null);

    // Step 4 data
    const [walletBalance, setWalletBalance] = useState<number | null>(null);
    const [originWarehouses, setOriginWarehouses] = useState<Warehouse[]>([]);
    const [generatedPdfUrl, setGeneratedPdfUrl] = useState<string | null>(null);
    const [trackingNumber, setTrackingNumber] = useState<string | null>(null);
    const [selectedCarrier, setSelectedCarrier] = useState<string | null>(null);

    // DANE search states
    const [originSearch, setOriginSearch] = useState("");
    const [destSearch, setDestSearch] = useState("");
    const [showOriginResults, setShowOriginResults] = useState(false);
    const [showDestResults, setShowDestResults] = useState(false);
    
    // Carrier Offices search states
    const [showOriginOffices, setShowOriginOffices] = useState(false);
    const [showDestOffices, setShowDestOffices] = useState(false);
    const [officeCarrier, setOfficeCarrier] = useState<string | null>(null);
    
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

    const orderIsCOD = !!(order?.cod_total && order.cod_total > 0);

    const step1Form = useForm<Step1Values>({
        resolver: zodResolver(step1Schema),
        mode: 'onChange',
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
            codValue: orderIsCOD ? order!.cod_total! : 0,
            includeGuideCost: false,
            insurance: false,
            codPaymentMethod: "cash",
        },
    });

    // Step 3 Form
    const step3Form = useForm<Step3Values>({
        resolver: zodResolver(step3Schema),
        mode: 'onChange',
        defaultValues: {
            originCompany: "Mi Empresa",
            originFirstName: "",
            originLastName: "",
            originEmail: "",
            originPhone: "",
            originSuburb: "",
            originCrossStreet: "",
            originReference: "",
            destCompany: "",
            destFirstName: "",
            destLastName: "",
            destEmail: "",
            destPhone: "",
            destSuburb: "",
            destCrossStreet: "",
            destReference: "",
            requestPickup: true,
            myShipmentReference: "",
            external_order_id: "",
        },
    });

    // Fetch initial data on open
    const handleWarehouseSelect = (wh: Warehouse) => {
        // Step 1
        step1Form.setValue("originDaneCode", wh.city_dane_code || "11001000", { shouldValidate: true });
        step1Form.setValue("originAddress", wh.street || wh.address, { shouldValidate: true });
        setOriginSearch(`${wh.city} (${wh.state})`);

        // Step 3
        step3Form.setValue("originCompany", wh.company || wh.name);
        step3Form.setValue("originFirstName", wh.first_name || wh.contact_name?.split(' ')[0] || "");
        step3Form.setValue("originLastName", wh.last_name || wh.contact_name?.split(' ').slice(1).join(' ') || "");
        step3Form.setValue("originEmail", wh.email || wh.contact_email || "");
        step3Form.setValue("originPhone", wh.phone || "");
        step3Form.setValue("originSuburb", wh.suburb || "");
        step3Form.setValue("originCrossStreet", wh.street || wh.address || "");
        step3Form.setValue("originReference", wh.city || wh.state || "");
    };

    useEffect(() => {
        if (isOpen) {
            const balanceBusinessId = effectiveBusinessId || undefined;
            getWalletBalanceAction(balanceBusinessId).then(res => {
                if (res.success && res.data) setWalletBalance(res.data.Balance);
            });
            getWarehousesAction({
                business_id: effectiveBusinessId || undefined,
                is_active: true,
                page: 1,
                page_size: 100,
            }).then(res => {
                if (res.data) {
                    setOriginWarehouses(res.data);
                    // Si la orden ya tiene warehouse_id, pre-seleccionar esa bodega
                    const preselect = order?.warehouse_id
                        ? res.data.find(w => w.id === order.warehouse_id)
                        : res.data.find(w => w.is_default);
                    if (preselect) {
                        handleWarehouseSelect(preselect);
                    }
                }
            }).catch(() => { });
        }
    }, [isOpen]);

    // Pre-fill from order
    useEffect(() => {
        if (isOpen && order) {
            // Step 1
            step1Form.setValue("contentValue", order.total_amount);
            step1Form.setValue("description", `Order ${order.order_number}`);
            step1Form.setValue("destAddress", (order.shipping_street || "").split(" | ")[0], { shouldValidate: true });

            if (order.weight && order.weight > 0) {
                step1Form.setValue("weight", order.weight, { shouldValidate: true });
                step1Form.setValue("height", order.height || 10, { shouldValidate: true });
                step1Form.setValue("width", order.width || 10, { shouldValidate: true });
                step1Form.setValue("length", order.length || 10, { shouldValidate: true });
            }

            // Try to find DANE code by city
            const mappedDane = findDaneCode(order.shipping_city || "", order.shipping_state || "");
            const finalDane = mappedDane || "11001000"; // Fallback to Bogota

            step1Form.setValue("destDaneCode", finalDane, { shouldValidate: true });
            const cityData = danes[finalDane as keyof typeof danes];
            if (cityData) {
                setDestSearch(`${(cityData as any).ciudad} (${(cityData as any).departamento})`);
            }

            if (order.cod_total && order.cod_total > 0) {
                step1Form.setValue("codValue", order.cod_total, { shouldValidate: true });
                step1Form.setValue("codPaymentMethod", "cash");
            }

            step3Form.setValue("destCompany", order.customer_name);
            step3Form.setValue("destFirstName", order.customer_name.split(" ")[0] || "");
            step3Form.setValue("destLastName", order.customer_name.split(" ").slice(1).join(" ") || ".");
            step3Form.setValue("destEmail", order.customer_email);
            step3Form.setValue("destPhone", order.customer_phone);
            const streetParts = (order.shipping_street || "").split(" | ");
            step3Form.setValue("destCrossStreet", (streetParts[0] || "").substring(0, 35));
            if (streetParts[1]) step3Form.setValue("destReference", streetParts[1].substring(0, 25));
            if (streetParts[2]) step3Form.setValue("destSuburb", streetParts[2].substring(0, 30));
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

    // SSE: listen for async quote/guide results
    useShipmentSSE({
        businessId,
        onQuoteReceived: (data) => {
            if (pendingCorrelationId && data.correlation_id !== pendingCorrelationId) return;
            setPendingCorrelationId(null);
            const quotes = data.quotes as any;
            const rates: EnvioClickRate[] = quotes?.data?.rates || quotes?.rates || [];
            if (rates.length > 0) {
                setRates(rates);
                if (pendingStep1DataRef.current) {
                    setStep1Data(pendingStep1DataRef.current);
                    pendingStep1DataRef.current = null;
                }
                setCurrentStep(2);
            } else {
                setError("No se encontraron tarifas disponibles");
            }
            setLoading(false);
        },
        onQuoteFailed: (data) => {
            if (pendingCorrelationId && data.correlation_id !== pendingCorrelationId) return;
            setPendingCorrelationId(null);
            pendingStep1DataRef.current = null;
            setError(data.error_message || "Error al cotizar envío");
            setLoading(false);
        },
        onGuideGenerated: async (data) => {
            if (!pendingGuideCorrelationId) return; // no pending guide request in this modal
            // Only reject if correlation_id is present AND doesn't match (backend now includes it)
            if (data.correlation_id && data.correlation_id !== pendingGuideCorrelationId) return;
            setPendingGuideCorrelationId(null);
            if (data.label_url) setGeneratedPdfUrl(data.label_url);
            if (data.tracking_number) {
                setTrackingNumber(data.tracking_number);
                if (data.carrier) setSelectedCarrier(data.carrier);

                if (selectedRate) {
                    const insuranceCost = step1Data?.insurance ? ((selectedRate.minimumInsurance ?? 0) + (selectedRate.extraInsurance ?? 0)) : 0;
                    const totalCost = selectedRate.flete + insuranceCost;
                    const balanceResponse = await getWalletBalanceAction();
                    if (balanceResponse.success && balanceResponse.data) {
                        setWalletBalance(balanceResponse.data.Balance);
                    }
                    const carrierName = (selectedRate?.carrier || data.carrier) ?? null;
                    if (carrierName) setSelectedCarrier(carrierName);
                    const carrierText = carrierName ? ` con ${carrierName}` : '';
                    setSuccess(`✅ Guía generada exitosamente. Se descontaron $${totalCost.toLocaleString()} de tu billetera${carrierText}.`);
                }

                // Fallback to selectedRate.carrier if data.carrier is empty
                const carrier = data.carrier || selectedRate?.carrier || '';
                if (onGuideGenerated) onGuideGenerated({
                    tracking_number: data.tracking_number,
                    carrier: carrier,
                    label_url: data.label_url
                });
            }
            setLoading(false);
        },
        onGuideFailed: (data) => {
            if (pendingGuideCorrelationId && data.correlation_id !== pendingGuideCorrelationId) return;
            setPendingGuideCorrelationId(null);
            setError(data.error_message || "Error al generar la guía");
            setLoading(false);
        },
    });

    // Timeout: if quote SSE never arrives, stop loading after 30s
    useEffect(() => {
        if (!pendingCorrelationId) return;
        const timeout = setTimeout(() => {
            setPendingCorrelationId(null);
            pendingStep1DataRef.current = null;
            setError("Tiempo de espera agotado. Verifica tu conexión e intenta de nuevo.");
            setLoading(false);
        }, 30000);
        return () => clearTimeout(timeout);
    }, [pendingCorrelationId]);

    // Timeout: if guide SSE never arrives, stop loading after 45s
    useEffect(() => {
        if (!pendingGuideCorrelationId) return;
        const timeout = setTimeout(() => {
            setPendingGuideCorrelationId(null);
            setError("Tiempo de espera agotado al generar la guía. Verifica en la lista de envíos.");
            setLoading(false);
        }, 45000);
        return () => clearTimeout(timeout);
    }, [pendingGuideCorrelationId]);

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
            setSelectedCarrier(null);
            setOfficeCarrier(null);
            setError(null);
            setSuccess(null);
            setPendingCorrelationId(null);
            setPendingGuideCorrelationId(null);
            pendingStep1DataRef.current = null;
            step1Form.reset();
            step3Form.reset();
        }
    }, [isOpen]);

    // Step 1: Quote (async - sends to queue, result arrives via SSE)
    const handleStep1Submit = async (data: Step1Values) => {
        // Autocompletar el paso 3 con las direcciones del paso 1
        step3Form.setValue("originCrossStreet", data.originAddress);
        step3Form.setValue("destCrossStreet", data.destAddress);

        // Check for validation errors
        const errors = step1Form.formState.errors;
        if (Object.keys(errors).length > 0) {
            const fieldLabels: { [key: string]: string } = {
                originDaneCode: "Ciudad de Origen",
                originAddress: "Dirección de Origen",
                destDaneCode: "Ciudad de Destino",
                destAddress: "Dirección de Destino",
                weight: "Peso del paquete",
                height: "Alto del paquete",
                width: "Ancho del paquete",
                length: "Largo del paquete",
                description: "Descripción del contenido",
                contentValue: "Valor de la mercancía",
                codPaymentMethod: "Método de pago COD",
            };

            const errorFields: string[] = [];
            Object.entries(errors).forEach(([field, error]) => {
                const label = fieldLabels[field] || field;
                errorFields.push(`  • ${label}`);
            });

            setError(`⚠️ Por favor completa los siguientes campos:\n${errorFields.join('\n')}`);
            setLoading(false);
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const quotePayload: EnvioClickQuoteRequest = {
                order_uuid: order?.id,
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
                insurance: data.insurance,
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
            if (!response.success) {
                setError(response.message || "Error al enviar solicitud de cotización");
                setLoading(false);
                return;
            }

            // Synchronous path: backend polled Redis and returned rates directly
            const syncRates: EnvioClickRate[] = response.data?.data?.rates || [];
            if (syncRates.length > 0) {
                setRates(syncRates);
                setStep1Data(data);
                setCurrentStep(2);
                setLoading(false);
                return;
            }

            // Asynchronous path: wait for SSE response
            pendingStep1DataRef.current = data;
            setPendingCorrelationId(response.data?.correlation_id || null);
            // loading stays true until SSE response arrives
        } catch (err: any) {
            setError(getActionError(err, "Error al cotizar envío"));
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
        // Collect all validation errors
        const errors = step3Form.formState.errors;

        // Debug: Log all data and errors
        console.log('📋 Step 3 Data:', data);
        console.log('❌ Step 3 Errors:', errors);
        console.log('📊 Error Count:', Object.keys(errors).length);

        if (Object.keys(errors).length > 0) {
            const fieldLabels: { [key: string]: string } = {
                originCrossStreet: "Calle",
                originReference: "Referencia",
                originSuburb: "Barrio",
                originCompany: "Empresa",
                originFirstName: "Nombre",
                originLastName: "Apellido",
                originPhone: "Teléfono",
                originEmail: "Email",
                destCrossStreet: "Calle",
                destReference: "Edificio/Interior/Apto",
                destSuburb: "Barrio",
                destCompany: "Empresa",
                destFirstName: "Nombre",
                destLastName: "Apellido",
                destPhone: "Teléfono",
                destEmail: "Email",
                myShipmentReference: "Mi Referencia de Envío",
            };

            // Group errors by section
            const originErrors: string[] = [];
            const destErrors: string[] = [];
            const otherErrors: string[] = [];

            Object.entries(errors).forEach(([field, error]) => {
                const label = fieldLabels[field] || field;
                const value = (data as any)[field];
                const valueStr = typeof value === 'string' ? `"${value}"` : String(value);
                const errorMsg = error?.message
                    ? `${label}: ${error.message} (Valor actual: ${valueStr})`
                    : `${label}: Campo inválido (Valor actual: ${valueStr})`;

                // Log individual error
                console.log(`Field: ${field}, Value: ${valueStr}, Error: ${error?.message}`);

                if (field.startsWith('origin')) {
                    originErrors.push(errorMsg);
                } else if (field.startsWith('dest')) {
                    destErrors.push(errorMsg);
                } else {
                    otherErrors.push(errorMsg);
                }
            });

            const sections: string[] = [];
            if (originErrors.length > 0) {
                sections.push(`📍 REMITENTE (Origen):\n${originErrors.map(e => `  ❌ ${e}`).join('\n')}`);
            }
            if (destErrors.length > 0) {
                sections.push(`📦 DESTINATARIO:\n${destErrors.map(e => `  ❌ ${e}`).join('\n')}`);
            }
            if (otherErrors.length > 0) {
                sections.push(`📋 INFORMACIÓN:\n${otherErrors.map(e => `  ❌ ${e}`).join('\n')}`);
            }

            setError(`⚠️ Errores encontrados - Por favor corrige lo siguiente:\n\n${sections.join('\n\n')}`);

            // Scroll al campo con error
            setTimeout(() => {
                const firstErrorField = Object.keys(errors)[0];
                const input = document.querySelector(`[name="${firstErrorField}"]`) as HTMLElement;
                if (input) {
                    input.scrollIntoView({ behavior: 'smooth', block: 'center' });
                    input.focus();
                }
            }, 100);
            return;
        }
        setStep3Data(data);
        setCurrentStep(4);
    };

    // Step 4: Generate Guide
    const handleFinalGenerate = async () => {
        const missingFields: string[] = [];

        if (!step1Data) {
            missingFields.push("⚠️ Paso 1: No completaste Origen, Destino o Paquete");
        }
        if (!selectedRate) {
            missingFields.push("⚠️ Paso 2: No seleccionaste una transportadora o tarifa");
        }
        if (!step3Data) {
            missingFields.push("⚠️ Paso 3: No completaste los detalles de dirección");
        }

        if (missingFields.length > 0) {
            setError(missingFields.join("\n"));
            return;
        }

        // Check wallet balance
        if (!selectedRate || !step3Data || !step1Data) return;
        const insuranceCost = step1Data.insurance ? ((selectedRate.minimumInsurance ?? 0) + (selectedRate.extraInsurance ?? 0)) : 0;
        const totalCost = selectedRate.flete + insuranceCost;
        if (walletBalance !== null && walletBalance < totalCost) {
            setError(`Saldo insuficiente. Necesitas $${totalCost.toLocaleString()} pero tienes $${walletBalance.toLocaleString()}`);
            return;
        }

        setLoading(true);
        setError(null);

        // DEBUG: Check carrier data
        console.log('DEBUG: About to generate guide with selectedRate=', selectedRate);
        console.log('DEBUG: selectedRate.carrier=', selectedRate.carrier);

        try {
            const generatePayload: EnvioClickQuoteRequest = {
                idRate: selectedRate.idRate,
                carrier: selectedRate.carrier,
                myShipmentReference: step3Data.myShipmentReference,
                external_order_id: step3Data.external_order_id,
                order_uuid: order?.id,
                requestPickup: step3Data.requestPickup,
                pickupDate: new Date().toISOString().split("T")[0],
                insurance: step1Data.insurance,
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
            if (!response.success) {
                setError(response.message || "Error al enviar solicitud de generación de guía");
                setLoading(false);
                return;
            }

            // Sync path: backend returned guide data directly (legacy)
            if (response.data?.data?.url) {
                const tracker = response.data.data.tracker;
                const carrier = (response.data?.data as any)?.carrier;
                setGeneratedPdfUrl(response.data.data.url);
                setTrackingNumber(tracker);
                if (carrier) setSelectedCarrier(carrier);

                const balanceResponse = await getWalletBalanceAction();
                if (balanceResponse.success && balanceResponse.data) {
                    setWalletBalance(balanceResponse.data.Balance);
                }
                const syncCarrier = (response.data?.data as any)?.carrier;
                const carrierText = syncCarrier ? ` con ${syncCarrier}` : '';
                if (syncCarrier) setSelectedCarrier(syncCarrier);
                setSuccess(`✅ Guía generada exitosamente. Se descontaron $${totalCost.toLocaleString()} de tu billetera${carrierText}.`);

                if (onGuideGenerated && tracker) {
                    onGuideGenerated({
                        tracking_number: tracker,
                        carrier: carrier,
                        label_url: generatedPdfUrl || undefined
                    });
                }
                setLoading(false);
                return;
            }

            // Async path (202 Accepted): wait for SSE event shipment.guide_generated
            setPendingGuideCorrelationId(response.data?.correlation_id || null);
            setGuideGenerationRequested(true);
            setCurrentStep(4);
            // loading stays true until SSE arrives or timeout fires
        } catch (err: any) {
            setError(getActionError(err, "Error al generar guía"));
            setLoading(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black/20 backdrop-blur-sm flex items-center justify-center z-50 p-2">
            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl flex flex-col overflow-hidden" style={{ width: '85%', maxHeight: '90vh' }}>
                <div className="bg-white dark:bg-gray-800 border-b px-3 py-3 flex-shrink-0">
                    <div className="flex justify-between items-center mb-2">
                        <h2 className="text-2xl font-bold text-purple-700 dark:text-purple-400">Generar Guía de Envío</h2>
                        <button
                            onClick={onClose}
                            className="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:text-gray-200 dark:text-gray-200 text-2xl"
                        >
                            ×
                        </button>
                    </div>
                    <Stepper steps={STEPS} currentStep={currentStep} />
                </div>

                {/* Content */}
                <div className="p-3 flex flex-col flex-1 overflow-hidden min-h-0">
                    {error && (
                        <div className="mb-3 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg text-red-700 dark:text-red-400 text-sm">
                            {error.includes('\n') ? (
                                <div>
                                    <div className="font-semibold mb-2">⚠️ Por favor corrige los siguientes errores:</div>
                                    <ul className="list-disc list-inside space-y-1">
                                        {error.split('\n').filter(line => line.trim()).map((line, idx) => (
                                            <li key={idx}>{line}</li>
                                        ))}
                                    </ul>
                                </div>
                            ) : (
                                error
                            )}
                        </div>
                    )}

                    {success && (
                        <div className="mb-2 p-2 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-700 rounded-lg text-green-700 dark:text-green-400">
                            {success}
                        </div>
                    )}

                    {/* Step 1: Origin/Destination/Package */}
                    {currentStep === 1 && (
                         
                        <form onSubmit={step1Form.handleSubmit(handleStep1Submit)} className="flex flex-col h-full overflow-hidden min-h-0" data-testid="step1-form">
                            <div className="flex-1 overflow-y-auto min-h-0 pr-3">
                                <div className="space-y-4">
                                    <div className="grid grid-cols-2 gap-4">
                                        <div className="bg-purple-50/50 dark:bg-purple-900/10 border border-purple-100 dark:border-purple-800/30 rounded-xl p-4 space-y-2">
                                            <div className="flex items-center justify-between">
                                                <div className="flex items-center gap-2">
                                                    <div className="w-8 h-8 rounded-lg bg-purple-100 dark:bg-purple-800/40 flex items-center justify-center text-purple-600 dark:text-purple-400 text-sm font-bold">A</div>
                                                    <h3 className="font-semibold text-base text-purple-700 dark:text-purple-400">Origen</h3>
                                                </div>
                                                {originWarehouses.length > 0 && (
                                                    <select
                                                        className="text-xs border border-gray-200 dark:border-gray-600 rounded px-2 py-1 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-1 focus:ring-purple-500"
                                                        onChange={(e) => {
                                                            const addr = originWarehouses.find(a => a.id === parseInt(e.target.value));
                                                            if (addr) handleWarehouseSelect(addr);
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
                                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-200 mb-1">
                                                    Ciudad remitente *
                                                </label>
                                                <input
                                                    type="text"
                                                    value={originSearch}
                                                    onChange={(e) => {
                                                        setOriginSearch(e.target.value);
                                                        setShowOriginResults(true);
                                                        if (!e.target.value) step1Form.setValue("originDaneCode", "", { shouldValidate: true });
                                                    }}
                                                    onFocus={() => setShowOriginResults(true)}
                                                    className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white ${step1Form.formState.errors.originDaneCode ? "border-red-500 bg-red-50 dark:bg-red-900/20" : "border-gray-300 dark:border-gray-600"}`}
                                                    placeholder="Buscar ciudad..."
                                                />
                                                {step1Form.formState.errors.originDaneCode && (
                                                    <p className="mt-1 text-xs text-red-600 dark:text-red-400">Selecciona una ciudad de origen de la lista</p>
                                                )}
                                                {showOriginResults && filteredOriginOptions.length > 0 && (
                                                    <div className="absolute z-10 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md shadow-lg max-h-60 overflow-y-auto">
                                                        {filteredOriginOptions.slice(0, 50).map((opt) => (
                                                            <div
                                                                key={opt.value}
                                                                onClick={() => {
                                                                    step1Form.setValue("originDaneCode", opt.value, { shouldValidate: true });
                                                                    setOriginSearch(opt.label);
                                                                    setShowOriginResults(false);
                                                                }}
                                                                className="px-3 py-2 hover:bg-gray-100 dark:bg-gray-700 dark:hover:bg-gray-700 cursor-pointer"
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
                                                {...step1Form.register("originAddress")}
                                                error={step1Form.formState.errors.originAddress?.message}
                                                placeholder="Calle 98 62-37"
                                            />
                                            {originSearch && (
                                                <div className="mt-1">
                                                    <button 
                                                        type="button" 
                                                        onClick={() => setShowOriginOffices(!showOriginOffices)}
                                                        className="text-xs text-purple-600 dark:text-purple-400 hover:underline flex items-center gap-1 font-medium"
                                                    >
                                                        📍 ¿Recoger en oficina principal?
                                                    </button>
                                                    {showOriginOffices && (
                                                        <CarrierOfficeSelector 
                                                            city={originSearch}
                                                            onSelectAddress={(addr, carrierId) => {
                                                                step1Form.setValue("originAddress", addr, { shouldValidate: true });
                                                                setOfficeCarrier(carrierId);
                                                                setShowOriginOffices(false);
                                                            }}
                                                            onClose={() => setShowOriginOffices(false)}
                                                        />
                                                    )}
                                                </div>
                                            )}
                                        </div>

                                        <div className="bg-blue-50/50 dark:bg-blue-900/10 border border-blue-100 dark:border-blue-800/30 rounded-xl p-4 space-y-2">
                                            <div className="flex items-center gap-2">
                                                <div className="w-8 h-8 rounded-lg bg-blue-100 dark:bg-blue-800/40 flex items-center justify-center text-blue-600 dark:text-blue-400 text-sm font-bold">B</div>
                                                <h3 className="font-semibold text-base text-blue-700 dark:text-blue-400">Destino</h3>
                                            </div>

                                            <div ref={destRef} className="relative">
                                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-200 mb-1">
                                                    Ciudad destinatario *
                                                </label>
                                                <input
                                                    type="text"
                                                    value={destSearch}
                                                    onChange={(e) => {
                                                        setDestSearch(e.target.value);
                                                        setShowDestResults(true);
                                                        if (!e.target.value) step1Form.setValue("destDaneCode", "", { shouldValidate: true });
                                                    }}
                                                    onFocus={() => setShowDestResults(true)}
                                                    className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white ${step1Form.formState.errors.destDaneCode ? "border-red-500 bg-red-50 dark:bg-red-900/20" : "border-gray-300 dark:border-gray-600"}`}
                                                    placeholder="Buscar ciudad..."
                                                />
                                                {step1Form.formState.errors.destDaneCode && (
                                                    <p className="mt-1 text-xs text-red-600 dark:text-red-400">Selecciona una ciudad de destino de la lista</p>
                                                )}
                                                {showDestResults && filteredDestOptions.length > 0 && (
                                                    <div className="absolute z-10 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md shadow-lg max-h-60 overflow-y-auto">
                                                        {filteredDestOptions.slice(0, 50).map((opt) => (
                                                            <div
                                                                key={opt.value}
                                                                onClick={() => {
                                                                    step1Form.setValue("destDaneCode", opt.value, { shouldValidate: true });
                                                                    setDestSearch(opt.label);
                                                                    setShowDestResults(false);
                                                                }}
                                                                className="px-3 py-2 hover:bg-gray-100 dark:bg-gray-700 dark:hover:bg-gray-700 cursor-pointer"
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
                                                {...step1Form.register("destAddress")}
                                                error={step1Form.formState.errors.destAddress?.message}
                                                placeholder="Carrera 46 # 93 - 45"
                                            />
                                            {destSearch && (
                                                <div className="mt-1">
                                                    <button 
                                                        type="button" 
                                                        onClick={() => setShowDestOffices(!showDestOffices)}
                                                        className="text-xs text-purple-600 dark:text-purple-400 hover:underline flex items-center gap-1 font-medium"
                                                    >
                                                        📍 ¿Enviar a oficina principal?
                                                    </button>
                                                    {showDestOffices && (
                                                        <CarrierOfficeSelector 
                                                            city={destSearch}
                                                            onSelectAddress={(addr, carrierId) => {
                                                                step1Form.setValue("destAddress", addr, { shouldValidate: true });
                                                                setOfficeCarrier(carrierId);
                                                                setShowDestOffices(false);
                                                            }}
                                                            onClose={() => setShowDestOffices(false)}
                                                        />
                                                    )}
                                                </div>
                                            )}
                                        </div>
                                    </div>

                                    <div className="bg-gray-50/80 dark:bg-gray-700/30 border border-gray-200 dark:border-gray-600/30 rounded-xl p-4">
                                        <div className="flex items-center gap-2 mb-3">
                                            <div className="w-8 h-8 rounded-lg bg-gray-200 dark:bg-gray-600 flex items-center justify-center text-gray-600 dark:text-gray-300 text-lg">
                                                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>
                                            </div>
                                            <h3 className="font-semibold text-base text-gray-700 dark:text-gray-200">Paquete</h3>
                                        </div>
                                        <div className="grid grid-cols-4 gap-2">
                                            <Input
                                                compact
                                                label="Peso (kg) *"
                                                type="number"
                                                step="0.1"
                                                {...step1Form.register("weight", { valueAsNumber: true })}
                                                error={step1Form.formState.errors.weight?.message}
                                            />
                                            <Input
                                                compact
                                                label="Alto (cm) *"
                                                type="number"
                                                {...step1Form.register("height", { valueAsNumber: true })}
                                                error={step1Form.formState.errors.height?.message}
                                            />
                                            <Input
                                                compact
                                                label="Ancho (cm) *"
                                                type="number"
                                                {...step1Form.register("width", { valueAsNumber: true })}
                                                error={step1Form.formState.errors.width?.message}
                                            />
                                            <Input
                                                compact
                                                label="Largo (cm) *"
                                                type="number"
                                                {...step1Form.register("length", { valueAsNumber: true })}
                                                error={step1Form.formState.errors.length?.message}
                                            />
                                        </div>

                                    {orderIsCOD && (
                                        <div className="flex items-center gap-2 px-3 py-2 rounded-lg bg-amber-50 border border-amber-300 dark:bg-amber-900/30 dark:border-amber-600">
                                            <span className="text-amber-700 dark:text-amber-300 font-semibold text-sm">
                                                Orden Contra Entrega - ${order!.cod_total!.toLocaleString()} COP
                                            </span>
                                        </div>
                                    )}

                                        <div className="grid grid-cols-3 gap-2 mt-3">
                                            <Input
                                                compact
                                                label="Descripcion *"
                                                {...step1Form.register("description")}
                                                error={step1Form.formState.errors.description?.message}
                                                placeholder="descripcion"
                                            />
                                            <Input
                                                compact
                                                label="Valor factura declarado *"
                                                type="number"
                                                {...step1Form.register("contentValue", { valueAsNumber: true })}
                                                error={step1Form.formState.errors.contentValue?.message}
                                            />
                                            <Input
                                                compact
                                                label="Valor contra entrega"
                                                type="number"
                                                {...step1Form.register("codValue", { valueAsNumber: true })}
                                                error={step1Form.formState.errors.codValue?.message}
                                                readOnly={orderIsCOD}
                                            />
                                        </div>

                                        <div className="flex items-center gap-6 mt-3 pt-3 border-t border-gray-200 dark:border-gray-600/30">
                                            <label className="flex items-center gap-2 cursor-pointer">
                                                <input
                                                    type="checkbox"
                                                    {...step1Form.register("insurance")}
                                                    className="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
                                                />
                                                <span className="text-sm text-gray-700 dark:text-gray-300">Asegurar envio</span>
                                            </label>
                                            {(step1Form.watch("codValue") ?? 0) > 0 && (
                                                <>
                                                    <label className="flex items-center gap-2 cursor-pointer">
                                                        <input
                                                            type="checkbox"
                                                            {...step1Form.register("includeGuideCost")}
                                                            className="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
                                                        />
                                                        <span className="text-sm text-gray-700 dark:text-gray-300">Incluir costo guia en contra entrega</span>
                                                    </label>
                                                    <div className="flex items-center gap-2">
                                                        <span className="text-sm text-gray-700 dark:text-gray-300">Metodo pago:</span>
                                                        <select
                                                            {...step1Form.register("codPaymentMethod")}
                                                            className="px-2 py-1 text-sm border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-1 focus:ring-purple-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                                                        >
                                                            <option value="cash">Efectivo</option>
                                                            <option value="data_phone">Datafono</option>
                                                        </select>
                                                    </div>
                                                </>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </form>
                    )}

                    {/* Step 2: Quote Selection */}
                    {currentStep === 2 && (
                        <div className="flex flex-col h-full overflow-y-auto">
                            <div className="pb-2">
                                <h3 className="font-semibold text-lg text-gray-700 dark:text-gray-200 mb-2">
                                    Filtra por servicio / Transportadora
                                </h3>
                                <div className="flex items-center gap-2 flex-wrap">
                                    <p className="text-sm text-gray-600 dark:text-gray-300">Todos los precios incluyen IVA</p>
                                    {(step1Data?.codValue ?? 0) > 0 && (
                                        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-100 text-amber-700 border border-amber-300 dark:bg-amber-900/30 dark:text-amber-300 dark:border-amber-600">
                                            Contra Entrega - Solo opciones contra entrega
                                        </span>
                                    )}
                                    {officeCarrier && (
                                        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-purple-100 text-purple-700 border border-purple-300 dark:bg-purple-900/30 dark:text-purple-300 dark:border-purple-600">
                                            Filtrado: {officeCarrier}
                                            <button
                                                type="button"
                                                onClick={() => setOfficeCarrier(null)}
                                                className="ml-1 text-purple-500 hover:text-purple-800 dark:hover:text-purple-100 font-bold leading-none"
                                                title="Quitar filtro"
                                            >
                                                x
                                            </button>
                                        </span>
                                    )}
                                </div>
                            </div>

                            <div className="overflow-y-auto border border-purple-200 rounded-lg p-3 bg-purple-50" style={{ maxHeight: 'calc(85vh - 350px)' }}>
                                {rates.length === 0 ? (
                                    <div className="flex items-center justify-center gap-3 py-10 text-purple-400">
                                        <div style={{ width: 28, height: 28, border: '3px solid #a855f7', borderTopColor: 'transparent', borderRadius: '50%', animation: 'spin 0.8s linear infinite' }} />
                                        <span className="text-sm font-medium">Cargando cotizaciones...</span>
                                        <style>{`@keyframes spin { to { transform: rotate(360deg); } }`}</style>
                                    </div>
                                ) : (() => {
                                    const filteredRates = rates.filter(rate => {
                                        const isCodRequest = (step1Data?.codValue ?? 0) > 0;
                                        if (isCodRequest && !rate.cod) return false;
                                        if (officeCarrier && !rate.carrier.toLowerCase().includes(officeCarrier.toLowerCase())) return false;
                                        return true;
                                    });
                                    if (filteredRates.length === 0) {
                                        return (
                                            <div className="flex flex-col items-center justify-center gap-2 py-10 text-amber-600">
                                                <span className="text-sm font-medium">
                                                    {(step1Data?.codValue ?? 0) > 0
                                                        ? "No hay transportadoras disponibles con opcion contra entrega para esta ruta"
                                                        : "No se encontraron cotizaciones para esta ruta"}
                                                </span>
                                            </div>
                                        );
                                    }
                                    return (
                                    <div className="grid grid-cols-4 gap-3 auto-rows-max">
                                        {filteredRates.map((rate) => {
                                            const insuranceCost = step1Data?.insurance ? ((rate.minimumInsurance ?? 0) + (rate.extraInsurance ?? 0)) : 0;
                                            const totalCost = rate.flete + insuranceCost;
                                            const isCOD = rate.cod;

                                            return (
                                                <div
                                                    key={rate.idRate}
                                                    onClick={() => handleRateSelection(rate)}
                                                    className="border border-gray-200 dark:border-gray-600 rounded-lg p-3 hover:border-purple-500 hover:shadow-md cursor-pointer transition-all bg-white dark:bg-gray-800"
                                                >
                                                    <div className="flex flex-col h-full">
                                                        <div className="flex flex-col items-center mb-2">
                                                            <div className={`${getCarrierLogoSize(rate.carrier).container} bg-purple-50 rounded-lg flex items-center justify-center mb-2 overflow-hidden`}>
                                                                <img
                                                                    src={getCarrierLogo(rate.carrier)}
                                                                    alt={rate.carrier}
                                                                    className={`${getCarrierLogoSize(rate.carrier).image} object-contain`}
                                                                    onError={(e) => {
                                                                        e.currentTarget.style.display = 'none';
                                                                        e.currentTarget.parentElement!.innerHTML = `<span class="font-bold text-xs text-center text-purple-600">${rate.carrier.substring(0, 3)}</span>`;
                                                                    }}
                                                                />
                                                            </div>
                                                            <div className="text-center">
                                                                <div className="font-semibold text-sm">{rate.carrier}</div>
                                                                <div className="text-xs text-gray-600 dark:text-gray-300">{rate.product}</div>
                                                            </div>
                                                        </div>

                                                        <div className="border-t pt-2 mt-2 flex-1">
                                                            <div className="text-center mb-1">
                                                                <div className="text-xl font-bold text-purple-600">
                                                                    ${totalCost.toLocaleString()}
                                                                </div>
                                                                <div className="text-xs text-gray-500 dark:text-gray-400">COP</div>
                                                                {step1Data?.insurance ? (
                                                                    <div className="mt-1 text-[10px] text-gray-500 dark:text-gray-400 leading-tight">
                                                                        Guía: ${rate.flete.toLocaleString()}<br />
                                                                        Seguro: ${insuranceCost.toLocaleString()}
                                                                    </div>
                                                                ) : (
                                                                    <div className="mt-1 text-[10px] text-gray-500 dark:text-gray-400 leading-tight">
                                                                        Guía: ${rate.flete.toLocaleString()}<br />
                                                                        Seguro: No asegurado
                                                                    </div>
                                                                )}
                                                            </div>
                                                            <div className="text-center">
                                                                <div className="text-xs text-gray-700 dark:text-gray-200 dark:text-gray-200 font-medium">
                                                                    {rate.deliveryDays} días
                                                                </div>
                                                            </div>
                                                            {isCOD && (
                                                                <div className="mt-1 text-center">
                                                                    <span className="inline-block px-2 py-0.5 rounded-full text-[10px] font-bold bg-amber-100 text-amber-700 border border-amber-300">
                                                                        Contra Entrega
                                                                    </span>
                                                                </div>
                                                            )}
                                                        </div>
                                                    </div>
                                                </div>
                                            );
                                        })}
                                    </div>
                                    );
                                })()}
                            </div>
                        </div>
                    )}

                    {/* Step 3: Details */}
                    {currentStep === 3 && (
                        <form onSubmit={step3Form.handleSubmit(handleStep3Submit)} className="flex flex-col h-full overflow-hidden">
                            <div className="overflow-y-auto flex-1 space-y-3 pr-1">
                                <div className="grid grid-cols-1 lg:grid-cols-2 gap-3">
                                    <div className="border border-purple-200 dark:border-purple-700 rounded-xl bg-purple-50/40 dark:bg-purple-900/10 p-4">
                                        <div className="flex items-center gap-2 mb-3 pb-2 border-b border-purple-200/60 dark:border-purple-700/40">
                                            <div className="w-7 h-7 rounded-lg bg-purple-100 dark:bg-purple-800/40 flex items-center justify-center">
                                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-purple-600 dark:text-purple-400"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
                                            </div>
                                            <h3 className="text-sm font-semibold text-purple-700 dark:text-purple-400">Remitente (Origen)</h3>
                                        </div>

                                        <p className="text-[10px] text-gray-500 dark:text-gray-400 uppercase tracking-wider font-semibold mb-1">Direccion</p>
                                        <div className="grid grid-cols-3 gap-1.5">
                                            <Input compact label="Calle *" {...step3Form.register("originCrossStreet")} error={step3Form.formState.errors.originCrossStreet?.message} placeholder="calle 75 sur n 42-97" />
                                            <Input compact label="Edificio/Apto *" {...step3Form.register("originReference")} error={step3Form.formState.errors.originReference?.message} placeholder="apt 801" />
                                            <Input compact label="Barrio *" {...step3Form.register("originSuburb")} error={step3Form.formState.errors.originSuburb?.message} placeholder="sector Aves Maria" />
                                        </div>

                                        <Input compact label="Empresa" {...step3Form.register("originCompany")} error={step3Form.formState.errors.originCompany?.message} placeholder="ProbabilityIA" className="mt-1.5" />

                                        <p className="text-[10px] text-gray-500 dark:text-gray-400 uppercase tracking-wider font-semibold mt-2 mb-1">Contacto</p>
                                        <div className="grid grid-cols-2 gap-1.5">
                                            <Input compact label="Nombre *" {...step3Form.register("originFirstName")} error={step3Form.formState.errors.originFirstName?.message} placeholder="Luisa" />
                                            <Input compact label="Apellido *" {...step3Form.register("originLastName")} error={step3Form.formState.errors.originLastName?.message} placeholder="Munoz" />
                                            <Input compact label="Telefono *" {...step3Form.register("originPhone")} error={step3Form.formState.errors.originPhone?.message} placeholder="3224098631" />
                                            <Input compact label="Correo *" type="email" {...step3Form.register("originEmail")} error={step3Form.formState.errors.originEmail?.message} placeholder="correo@ejemplo.com" />
                                        </div>
                                    </div>

                                    <div className="border border-blue-200 dark:border-blue-700 rounded-xl bg-blue-50/40 dark:bg-blue-900/10 p-4">
                                        <div className="flex items-center gap-2 mb-3 pb-2 border-b border-blue-200/60 dark:border-blue-700/40">
                                            <div className="w-7 h-7 rounded-lg bg-blue-100 dark:bg-blue-800/40 flex items-center justify-center">
                                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-blue-600 dark:text-blue-400"><rect x="1" y="3" width="15" height="13"/><polygon points="16 8 20 8 23 11 23 16 16 16 16 8"/><circle cx="5.5" cy="18.5" r="2.5"/><circle cx="18.5" cy="18.5" r="2.5"/></svg>
                                            </div>
                                            <h3 className="text-sm font-semibold text-blue-700 dark:text-blue-400">Destinatario (Destino)</h3>
                                        </div>

                                        <p className="text-[10px] text-gray-500 dark:text-gray-400 uppercase tracking-wider font-semibold mb-1">Direccion</p>
                                        <div className="grid grid-cols-3 gap-1.5">
                                            <Input compact label="Calle *" {...step3Form.register("destCrossStreet")} error={step3Form.formState.errors.destCrossStreet?.message} placeholder="calle 75 sur n 42-97" />
                                            <Input compact label="Edificio/Apto" {...step3Form.register("destReference")} error={step3Form.formState.errors.destReference?.message} placeholder="casa #" />
                                            <Input compact label="Barrio" {...step3Form.register("destSuburb")} error={step3Form.formState.errors.destSuburb?.message} placeholder="Nombre barrio" />
                                        </div>

                                        <Input compact label="Empresa" {...step3Form.register("destCompany")} error={step3Form.formState.errors.destCompany?.message} placeholder="Empresa (opcional)" className="mt-1.5" />

                                        <p className="text-[10px] text-gray-500 dark:text-gray-400 uppercase tracking-wider font-semibold mt-2 mb-1">Contacto</p>
                                        <div className="grid grid-cols-2 gap-1.5">
                                            <Input compact label="Nombre *" {...step3Form.register("destFirstName")} error={step3Form.formState.errors.destFirstName?.message} placeholder="Luisa" />
                                            <Input compact label="Apellido *" {...step3Form.register("destLastName")} error={step3Form.formState.errors.destLastName?.message} placeholder="Munoz" />
                                            <Input compact label="Telefono *" {...step3Form.register("destPhone")} error={step3Form.formState.errors.destPhone?.message} placeholder="3224098631" />
                                            <Input compact label="Correo *" type="email" {...step3Form.register("destEmail")} error={step3Form.formState.errors.destEmail?.message} placeholder="correo@ejemplo.com" />
                                        </div>
                                    </div>
                                </div>

                                <div className="border border-gray-200 dark:border-gray-600 rounded-xl bg-gray-50/60 dark:bg-gray-700/30 p-4">
                                    <p className="text-[10px] text-gray-500 dark:text-gray-400 uppercase tracking-wider font-semibold mb-2">Opciones adicionales</p>
                                    <div className="grid grid-cols-2 gap-2">
                                        <Input compact label="Mi referencia de envio" {...step3Form.register("myShipmentReference")} error={step3Form.formState.errors.myShipmentReference?.message} placeholder="Orden 5649" />
                                        <Input compact label="Numero de orden externo" {...step3Form.register("external_order_id")} error={step3Form.formState.errors.external_order_id?.message} placeholder="ORD345678" />
                                    </div>
                                    <label className="flex items-center space-x-2 mt-2">
                                        <input type="checkbox" {...step3Form.register("requestPickup")} className="rounded w-5 h-5" />
                                        <span className="text-sm font-medium">Solicitar recoleccion</span>
                                    </label>
                                </div>
                            </div>
                        </form>
                    )}

                    {/* Step 4: Payment & Confirmation */}
                    {currentStep === 4 && selectedRate && (
                        <div className="flex flex-col h-full w-full overflow-hidden gap-3">
                            {/* Resumen de Envío - No Scrolleable */}
                            <div className="flex-shrink-0 space-y-3">
                                <h3 className="font-semibold text-lg text-gray-700 dark:text-gray-200 dark:text-gray-200">Resumen de tu envío</h3>

                                <div className="bg-gray-50 dark:bg-gray-700 p-2 rounded-lg">
                                    <div className="flex items-center justify-between mb-2">
                                        <div className="flex items-center space-x-2">
                                            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                                            </svg>
                                            <span className="font-medium">1 Envíos</span>
                                        </div>
                                        <div className="text-right">
                                            <div className="text-sm text-gray-600 dark:text-gray-300">TOTAL:</div>
                                            <div className="text-2xl font-bold text-purple-600">
                                                ${(selectedRate.flete + (step1Data?.insurance ? ((selectedRate.minimumInsurance ?? 0) + (selectedRate.extraInsurance ?? 0)) : 0)).toLocaleString()}
                                            </div>
                                            {step1Data?.insurance ? (
                                                <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                                                    Guía: ${selectedRate.flete.toLocaleString()} | Seguro: ${((selectedRate.minimumInsurance ?? 0) + (selectedRate.extraInsurance ?? 0)).toLocaleString()}
                                                </div>
                                            ) : (
                                                <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                                                    Guía: ${selectedRate.flete.toLocaleString()} | Seguro: No asegurado
                                                </div>
                                            )}
                                        </div>
                                    </div>

                                    <div className="border-t pt-4 flex items-center gap-4">
                                        <img
                                            src={getCarrierLogo(selectedRate.carrier)}
                                            alt={selectedRate.carrier}
                                            className="w-16 h-16 object-contain rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 p-1 flex-shrink-0"
                                            onError={(e) => { (e.target as HTMLImageElement).style.display = 'none'; }}
                                        />
                                        <div>
                                            <div className="font-medium text-gray-800 dark:text-gray-100">{selectedRate.carrier}</div>
                                            <div className="text-sm text-gray-500 dark:text-gray-400">{selectedRate.product}</div>
                                            {selectedRate.deliveryDays > 0 && (
                                                <div className="text-xs text-gray-400 mt-1">{selectedRate.deliveryDays} día{selectedRate.deliveryDays !== 1 ? 's' : ''} hábil{selectedRate.deliveryDays !== 1 ? 'es' : ''}</div>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </div>

                            {/* Información de Pago - Con Scroll */}
                            <div className="flex-1 min-h-0 overflow-y-auto overflow-x-hidden space-y-3 pr-3" style={{ maxHeight: 'calc(85vh - 280px)' }}>
                                <div>
                                    <h4 className="font-medium text-gray-700 dark:text-gray-200 dark:text-gray-200 mb-3">Selecciona tu método de pago</h4>
                                    <div className={`grid gap-2 ${generatedPdfUrl ? 'grid-cols-2' : 'grid-cols-1'}`}>
                                        <div className="border-2 border-purple-500 rounded-lg p-2 bg-purple-50">
                                            <div className="flex items-center justify-center mb-2">
                                                <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
                                                </svg>
                                            </div>
                                            <div className="text-center font-semibold">Monedero</div>
                                            <div className="text-center text-sm text-gray-600 dark:text-gray-300">
                                                ${walletBalance?.toLocaleString() || 0}
                                            </div>
                                        </div>

                                        {generatedPdfUrl && (
                                            /* ── Success state ── */
                                            <div className="rounded-xl border-2 border-emerald-200 bg-emerald-50 p-4 flex flex-col items-center gap-3">
                                                <div className="w-12 h-12 rounded-full bg-emerald-100 flex items-center justify-center flex-shrink-0">
                                                    <svg className="w-6 h-6 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
                                                    </svg>
                                                </div>
                                                <div className="text-center min-w-0">
                                                    <p className="font-bold text-emerald-800 text-sm">¡Guía generada exitosamente!</p>
                                                    {selectedCarrier && (
                                                        <div className="flex items-center justify-center gap-2 mt-1.5">
                                                            <img
                                                                src={getCarrierLogo(selectedCarrier)}
                                                                alt={selectedCarrier}
                                                                className="w-5 h-5 object-contain"
                                                            />
                                                            <span className="text-xs text-emerald-700 font-semibold">{selectedCarrier}</span>
                                                        </div>
                                                    )}
                                                    {trackingNumber && (
                                                        <p className="text-xs text-emerald-700 mt-1.5 font-mono bg-emerald-100 px-2 py-0.5 rounded-full inline-block">
                                                            {trackingNumber}
                                                        </p>
                                                    )}
                                                </div>
                                                <div className="flex flex-col gap-1.5 w-full">
                                                    <a
                                                        href={generatedPdfUrl}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        className="flex items-center justify-center gap-1 w-full py-1.5 px-2 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition-colors text-xs"
                                                    >
                                                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                                                        </svg>
                                                        Abrir
                                                    </a>
                                                    <a
                                                        href={generatedPdfUrl}
                                                        download
                                                        className="flex items-center justify-center gap-1 w-full py-1.5 px-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-200 dark:text-gray-200 font-semibold rounded-lg transition-colors text-xs"
                                                    >
                                                        <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                                                        </svg>
                                                        Descargar
                                                    </a>
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                </div>
                                <div className="pb-2" />
                            </div>
                        </div>
                    )}
                </div>

                {/* Footer with Buttons */}
                <div className="bg-white dark:bg-gray-800 border-t px-3 py-3 flex-shrink-0 flex justify-between items-center gap-3">
                    {/* Back Button */}
                    {(currentStep === 2 || currentStep === 3 || currentStep === 4) && !generatedPdfUrl && (
                        <Button
                            variant="outline"
                            onClick={() => setCurrentStep(currentStep - 1)}
                            disabled={loading}
                        >
                            Atrás
                        </Button>
                    )}

                    {/* Spacer when no back button */}
                    {(currentStep === 1 || (currentStep === 4 && generatedPdfUrl)) && (
                        <div />
                    )}

                    {/* Step 1: Next Button */}
                    {currentStep === 1 && (
                        <Button
                            variant="primary"
                            onClick={() => {
                                const fieldLabels: { [key: string]: string } = {
                                    originDaneCode: "Ciudad de Origen",
                                    originAddress: "Dirección de Origen",
                                    destDaneCode: "Ciudad de Destino",
                                    destAddress: "Dirección de Destino",
                                    weight: "Peso del paquete",
                                    height: "Alto del paquete",
                                    width: "Ancho del paquete",
                                    length: "Largo del paquete",
                                    description: "Descripción del contenido",
                                    contentValue: "Valor de la mercancía",
                                    codPaymentMethod: "Método de pago COD",
                                };
                                step1Form.handleSubmit(handleStep1Submit, (errors) => {
                                    const errorFields = Object.entries(errors).map(
                                        ([field, err]) => `  • ${fieldLabels[field] || field}: ${(err as any)?.message || "inválido"}`
                                    );
                                    setError(`⚠️ Por favor completa los siguientes campos:\n${errorFields.join('\n')}`);
                                })();
                            }}
                            disabled={loading}
                            style={{ background: '#7c3aed' }}
                        >
                            {loading ? "Cotizando..." : "Siguiente"}
                        </Button>
                    )}

                    {/* Step 2: NO "Siguiente" Button - User selects a rate to advance */}
                    {currentStep === 2 && (
                        <div className="text-sm text-gray-600 dark:text-gray-300 italic">
                            📌 Selecciona una transportadora para continuar
                        </div>
                    )}

                    {/* Step 3: Next Button - DISABLED if form has errors */}
                    {currentStep === 3 && (
                        <Button
                            variant="primary"
                            onClick={async () => {
                                // Trigger validation
                                const isValid = await step3Form.trigger();
                                if (isValid) {
                                    const data = step3Form.getValues();
                                    setStep3Data(data);
                                    setCurrentStep(4);
                                }
                            }}
                            disabled={loading || Object.keys(step3Form.formState.errors).length > 0}
                            title={Object.keys(step3Form.formState.errors).length > 0 ? "Completa todos los campos requeridos" : ""}
                            style={{
                                background: Object.keys(step3Form.formState.errors).length > 0 ? '#ccc' : '#7c3aed',
                                cursor: Object.keys(step3Form.formState.errors).length > 0 ? 'not-allowed' : 'pointer'
                            }}
                        >
                            {Object.keys(step3Form.formState.errors).length > 0
                                ? `⚠️ ${Object.keys(step3Form.formState.errors).length} campo(s) incompleto(s)`
                                : "Siguiente"
                            }
                        </Button>
                    )}

                    {/* Step 4: Pay Button */}
                    {currentStep === 4 && !generatedPdfUrl && (
                        <Button
                            onClick={handleFinalGenerate}
                            disabled={loading}
                            className="bg-green-600 hover:bg-green-700"
                        >
                            {loading ? "Generando..." : "Pagar guías"}
                        </Button>
                    )}
                </div>
            </div>
        </div>
    );
}
