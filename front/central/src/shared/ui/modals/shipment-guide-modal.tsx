'use client';

import { useState, useEffect, useRef } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Input, Button, Stepper, BrandLoaderOverlay } from "@/shared/ui";
import { ShipmentApiRepository } from "@/services/modules/shipments/infra/repository/api-repository";
import { EnvioClickQuoteRequest, EnvioClickRate } from "@/services/modules/shipments/domain/types";
import { Order } from "@/services/modules/orders/domain/types";
import { getWalletBalanceAction } from "@/services/modules/wallet/infra/actions";
import { quoteShipmentAction, generateGuideAction, getShipmentByIdAction } from "@/services/modules/shipments/infra/actions";
import { getWarehousesAction } from "@/services/modules/warehouses/infra/actions";
import { Warehouse } from "@/services/modules/warehouses/domain/types";
import danes from "@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json";
import { findDaneCode } from "@/shared/utils/dane-lookup";
import { CarrierEffectivenessRates } from "@/services/modules/geozones/ui/components/CarrierEffectivenessRates";
import { getProbabilityByCarrierAction } from "@/services/modules/geozones/infra/actions";
import type { ProbabilityResult } from "@/services/modules/geozones/domain/types";
import { useShipmentSSE } from "@/services/modules/shipments/ui/hooks/useShipmentSSE";
import { usePermissions } from "@/shared/contexts/permissions-context";
import { getActionError } from '@/shared/utils/action-result';
import { CarrierOfficeSelector } from "@/services/modules/shipments/ui/components/CarrierOfficeSelector";
import { CookieStorage } from "@/shared/config";
import '@/shared/ui/styles/shipment-modals.css';
import dynamic from 'next/dynamic';

const GUIDE_SSE_GRACE_MS = 45000;
const GUIDE_POLL_INTERVAL_MS = 5000;
const GUIDE_POLL_MAX_ATTEMPTS = 24;

const GeozoneMiniMap = dynamic(
    () => import('@/services/modules/geozones/ui/components/GeozoneMiniMap').then(m => m.GeozoneMiniMap),
    { ssr: false }
);

const normalizeLocationName = (str: string) => {
    if (!str) return "";
    let s = str.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toUpperCase().trim();
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

const formatCarrierName = (name: string): string => {
    return name
        .replace(/_/g, ' ')
        .replace(/EXPRESS/g, '')
        .trim();
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
        'MENSAJER\u00cdA URBANA': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_mensajerosUrbanos.png',
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

interface ShipmentGuideModalProps {
    isOpen: boolean;
    onClose: () => void;
    order?: Order;
    onGuideGenerated?: (data: { tracking_number: string; carrier?: string; label_url?: string }) => void;
    recommendedCarrier?: string;
}

function normalizeColombianPhone(raw: string | undefined | null): string {
    if (!raw) return "";
    let digits = raw.replace(/\D/g, "");
    if (digits.length === 12 && digits.startsWith("57")) {
        digits = digits.slice(2);
    } else if (digits.length === 13 && digits.startsWith("0057")) {
        digits = digits.slice(4);
    } else if (digits.length === 11 && digits.startsWith("0")) {
        digits = digits.slice(1);
    }
    return digits;
}

const step1Schema = z.object({
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
    codValue: z.number().min(0).max(3000000).optional(),
    includeGuideCost: z.boolean(),
    insurance: z.boolean(),
    codPaymentMethod: z.enum(["cash", "data_phone"]),
});

const step3Schema = z.object({
    originCompany: z.string().min(2, "Min 2 caracteres").max(28, "Max 28 caracteres"),
    originFirstName: z.string().min(2, "Min 2 caracteres").max(14, "Max 14 caracteres (l\u00edmite del transportador)"),
    originLastName: z.string().min(2, "Min 2 caracteres").max(14, "Max 14 caracteres (l\u00edmite del transportador)"),
    originEmail: z.string().email("Email invalido").min(8, "Min 8 caracteres").max(60, "Max 60 caracteres"),
    originPhone: z.string().length(10, "Debe tener 10 digitos"),
    originSuburb: z.string().min(2, "Min 2 caracteres").max(30, "Max 30 caracteres"),
    originCrossStreet: z.string().min(2, "Min 2 caracteres").max(35, "Max 35 caracteres"),
    originReference: z.string().min(2, "Min 2 caracteres").max(25, "Max 25 caracteres"),
    destCompany: z.string().min(2, "Min 2 caracteres").max(28, "Max 28 caracteres").optional(),
    destFirstName: z.string().min(2, "Min 2 caracteres").max(14, "Max 14 caracteres (l\u00edmite del transportador)"),
    destLastName: z.string().min(2, "Min 2 caracteres").max(14, "Max 14 caracteres (l\u00edmite del transportador)"),
    destEmail: z.string().email("Email invalido").min(8, "Min 8 caracteres").max(60, "Max 60 caracteres"),
    destPhone: z.string().length(10, "Debe tener 10 digitos"),
    destSuburb: z.string().max(30, "Max 30 caracteres").refine((v) => !v || v.length >= 2, "Min 2 caracteres").optional(),
    destCrossStreet: z.string().min(2, "Min 2 caracteres").max(35, "Max 35 caracteres"),
    destReference: z.string().max(25, "Max 25 caracteres").refine((v) => !v || v.length >= 2, "Min 2 caracteres").optional(),
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

    const [pendingCorrelationId, setPendingCorrelationId] = useState<string | null>(null);
    const pendingStep1DataRef = useRef<Step1Values | null>(null);

    const [pendingGuideCorrelationId, setPendingGuideCorrelationId] = useState<string | null>(null);
    const [guideGenerationRequested, setGuideGenerationRequested] = useState(false);

    const { permissions, isSuperAdmin } = usePermissions();
    const businessId = permissions?.business_id || 0;
    const effectiveBusinessId = isSuperAdmin ? (order?.business_id ?? 0) : 0;

    const [step1Data, setStep1Data] = useState<Step1Values | null>(null);

    const [rates, setRates] = useState<EnvioClickRate[]>([]);
    const [selectedRate, setSelectedRate] = useState<EnvioClickRate | null>(null);

    const [step3Data, setStep3Data] = useState<Step3Values | null>(null);

    const [walletBalance, setWalletBalance] = useState<number | null>(null);
    const [originWarehouses, setOriginWarehouses] = useState<Warehouse[]>([]);
    const [selectedOriginWarehouse, setSelectedOriginWarehouse] = useState<Warehouse | null>(null);
    const [carrierProbabilities, setCarrierProbabilities] = useState<ProbabilityResult[]>([]);

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

    const normalizeCarrierKey = (s: string) =>
        (s || '').normalize('NFD').replace(/[̀-ͯ]/g, '').toUpperCase().replace(/[^A-Z0-9]/g, '');

    const quotedCarrierKey = normalizeCarrierKey((order?.quoted_shipping?.carrier || '').split(' - ')[0]);

    useEffect(() => {
        if (currentStep !== 2 || !order?.id || !order?.business_id) return;
        let cancelled = false;
        getProbabilityByCarrierAction(order.id, order.business_id)
            .then((res) => { if (!cancelled) setCarrierProbabilities(Array.isArray(res) ? res : []); })
            .catch(() => { if (!cancelled) setCarrierProbabilities([]); });
        return () => { cancelled = true; };
    }, [currentStep, order?.id, order?.business_id]);

    const selectedCarrierKey = selectedRate ? normalizeCarrierKey(selectedRate.carrier || '') : '';
    let selectedCarrierProb = selectedCarrierKey
        ? carrierProbabilities.find(p => normalizeCarrierKey(p.carrier || '') === selectedCarrierKey)
        : undefined;

    if (!selectedCarrierProb && selectedRate?.carrier) {
        const seed = `${order?.business_id}|${selectedRate.carrier}|delivery`;
        let h = 0;
        for (let i = 0; i < seed.length; i++) h = (h * 31 + seed.charCodeAt(i)) >>> 0;
        const deliveryRate = 0.80 + ((h % 1001) / 1000) * 0.15;
        console.log('Generating fallback for carrier:', selectedRate.carrier, 'rate:', deliveryRate);
        selectedCarrierProb = {
            found: false,
            delivery_rate: deliveryRate,
            carrier: selectedRate.carrier,
            is_estimated: true,
            estimate_source: 'carrier_baseline',
        } as ProbabilityResult;
    } else if (selectedCarrierProb) {
        console.log('Using real carrier prob:', selectedCarrierProb.carrier, 'rate:', selectedCarrierProb.delivery_rate);
    } else {
        console.log('No selected rate or carrier:', selectedRate, carrierProbabilities);
    }
    const [generatedPdfUrl, setGeneratedPdfUrl] = useState<string | null>(null);
    const [generatedShipmentId, setGeneratedShipmentId] = useState<number | null>(null);
    const pendingGuideShipmentIdRef = useRef<number | null>(null);
    const [guideBlocked, setGuideBlocked] = useState(false);
    const [trackingNumber, setTrackingNumber] = useState<string | null>(null);
    const [selectedCarrier, setSelectedCarrier] = useState<string | null>(null);

    const [guideFormats, setGuideFormats] = useState<any[]>([]);
    const [showGuideFormatDropdown, setShowGuideFormatDropdown] = useState(false);
    const [selectedGuideFormat, setSelectedGuideFormat] = useState<string | null>(null);
    const formatDropdownRef = useRef<HTMLDivElement>(null);

    const [originSearch, setOriginSearch] = useState("");
    const [destSearch, setDestSearch] = useState("");
    const [showOriginResults, setShowOriginResults] = useState(false);
    const [showDestResults, setShowDestResults] = useState(false);
    
    const [showOriginOffices, setShowOriginOffices] = useState(false);
    const [showDestOffices, setShowDestOffices] = useState(false);
    const [officeCarrier, setOfficeCarrier] = useState<string | null>(null);
    const [mapViewMode, setMapViewMode] = useState<'origin-destination' | 'destination-only'>('origin-destination');

    const originRef = useRef<HTMLDivElement>(null);
    const destRef = useRef<HTMLDivElement>(null);

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
            codValue: orderIsCOD ? order!.cod_total! : 0,
            includeGuideCost: false,
            insurance: false,
            codPaymentMethod: "cash",
        },
    });

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

    const handleWarehouseSelect = (wh: Warehouse) => {
        setSelectedOriginWarehouse(wh);
        const daneCode = wh.city_dane_code || findDaneCode(wh.city || "", wh.state || "") || "";
        step1Form.setValue("originDaneCode", daneCode, { shouldValidate: true });
        step1Form.setValue("originAddress", wh.street || wh.address, { shouldValidate: true });
        setOriginSearch(`${wh.city} (${wh.state})`);

        let buildingRef = "";
        if (wh.id) {
            try {
                buildingRef = localStorage.getItem(`wh_building_${wh.id}`) || "";
            } catch {}
        }

        step3Form.setValue("originCompany", wh.company || wh.name);
        step3Form.setValue("originFirstName", wh.first_name || wh.contact_name?.split(' ')[0] || "");
        step3Form.setValue("originLastName", wh.last_name || wh.contact_name?.split(' ').slice(1).join(' ') || "");
        step3Form.setValue("originEmail", wh.email || wh.contact_email || "");
        step3Form.setValue("originPhone", normalizeColombianPhone(wh.phone));
        step3Form.setValue("originSuburb", wh.suburb || "");
        step3Form.setValue("originCrossStreet", wh.street || wh.address || "");
        step3Form.setValue("originReference", buildingRef);
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

    useEffect(() => {
        if (isOpen && order) {
            step1Form.setValue("contentValue", order.total_amount);
            step1Form.setValue("description", `Order ${order.order_number}`);
            step1Form.setValue("destAddress", (order.shipping_street || "").split(" | ")[0], { shouldValidate: true });

            if (order.weight && order.weight > 0) {
                step1Form.setValue("weight", order.weight, { shouldValidate: true });
                step1Form.setValue("height", order.height || 10, { shouldValidate: true });
                step1Form.setValue("width", order.width || 10, { shouldValidate: true });
                step1Form.setValue("length", order.length || 10, { shouldValidate: true });
            }

            const verifiedDaneKey = order.destination_dane_code ? `${order.destination_dane_code}000` : null;
            const verifiedDane = verifiedDaneKey && danes[verifiedDaneKey as keyof typeof danes]
                ? verifiedDaneKey
                : null;
            const mappedDane = verifiedDane || findDaneCode(order.shipping_city || "", order.shipping_state || "");
            if (mappedDane) {
                step1Form.setValue("destDaneCode", mappedDane, { shouldValidate: true });
                const cityData = danes[mappedDane as keyof typeof danes];
                if (cityData) {
                    setDestSearch(`${(cityData as any).ciudad} (${(cityData as any).departamento})`);
                }
            }

            if (order.cod_total && order.cod_total > 0) {
                step1Form.setValue("codValue", order.cod_total, { shouldValidate: true });
                step1Form.setValue("codPaymentMethod", "cash");
            }

            step3Form.setValue("destCompany", order.customer_name);
            step3Form.setValue("destFirstName", order.customer_name.split(" ")[0] || "");
            step3Form.setValue("destLastName", order.customer_name.split(" ").slice(1).join(" ") || ".");
            step3Form.setValue("destEmail", order.customer_email);
            step3Form.setValue("destPhone", normalizeColombianPhone(order.customer_phone));
            const streetParts = (order.shipping_street || "").split(" | ");
            step3Form.setValue("destCrossStreet", (streetParts[0] || "").substring(0, 35));
            if (streetParts[1]) step3Form.setValue("destReference", streetParts[1].substring(0, 25));
            if (streetParts[2]) step3Form.setValue("destSuburb", streetParts[2].substring(0, 30));
            step3Form.setValue("myShipmentReference", "Orden " + (order.internal_number || order.order_number));
            step3Form.setValue("external_order_id", order.order_number);
        }
    }, [isOpen, order]);

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

    useShipmentSSE({
        businessId,
        onQuoteReceived: (data) => {
            if (!pendingCorrelationId || data.correlation_id !== pendingCorrelationId) return;
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
            if (!pendingCorrelationId || data.correlation_id !== pendingCorrelationId) return;
            setPendingCorrelationId(null);
            pendingStep1DataRef.current = null;
            setError(data.error_message || "Error al cotizar envío");
            setLoading(false);
        },
        onGuideGenerated: async (data) => {
            if (!pendingGuideCorrelationId) return;
            if (data.correlation_id && data.correlation_id !== pendingGuideCorrelationId) return;
            setPendingGuideCorrelationId(null);
            if (data.label_url) setGeneratedPdfUrl(data.label_url);
            if (data.shipment_id) setGeneratedShipmentId(data.shipment_id);
            if (data.tracking_number) {
                setTrackingNumber(data.tracking_number);
                if (data.carrier) setSelectedCarrier(data.carrier);

                if (selectedRate) {
                    const insuranceCost = (selectedRate.minimumInsurance ?? 0) + (step1Data?.insurance ? (selectedRate.extraInsurance ?? 0) : 0);
                    const codMargin = selectedRate.cod ? (selectedRate.codProbabilityMargin ?? 0) : 0;
                    const totalCost = selectedRate.flete + insuranceCost + codMargin;
                    const balanceResponse = await getWalletBalanceAction();
                    if (balanceResponse.success && balanceResponse.data) {
                        setWalletBalance(balanceResponse.data.Balance);
                    }
                    const carrierName = (selectedRate?.carrier || data.carrier) ?? null;
                    if (carrierName) setSelectedCarrier(carrierName);
                    const carrierText = carrierName ? ` con ${carrierName}` : '';
                    setSuccess(`✅ Guía generada exitosamente. Se descontaron $${totalCost.toLocaleString()} de tu billetera${carrierText}.`);
                }

                const carrier = data.carrier || selectedRate?.carrier || '';
                if (onGuideGenerated) onGuideGenerated({
                    tracking_number: data.tracking_number,
                    carrier: carrier,
                    label_url: data.label_url
                });
            }
            setLoading(false);
        },
        onGuideFailed: async (data) => {
            if (pendingGuideCorrelationId && data.correlation_id !== pendingGuideCorrelationId) return;
            setPendingGuideCorrelationId(null);
            setError(data.error_message || "Error al generar la guía");
            setLoading(false);
            const id = pendingGuideShipmentIdRef.current;
            if (!id) return;
            const res = await getShipmentByIdAction(id);
            if (res.success && (res.data as any)?.status === 'needs_verification') {
                setGuideBlocked(true);
            }
        },
    });

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

    useEffect(() => {
        if (!pendingGuideCorrelationId) return;
        let cancelled = false;
        let polls = 0;

        const finishWithShipment = (shipment: any) => {
            setPendingGuideCorrelationId(null);
            setGeneratedShipmentId(shipment.id ?? null);
            if (shipment.guide_url) setGeneratedPdfUrl(shipment.guide_url);
            if (shipment.tracking_number) setTrackingNumber(shipment.tracking_number);
            if (shipment.carrier) setSelectedCarrier(shipment.carrier);
            setSuccess(`Gu\u00eda generada exitosamente${shipment.carrier ? ` con ${shipment.carrier}` : ''}.`);
            setLoading(false);
            if (onGuideGenerated && shipment.tracking_number) {
                onGuideGenerated({
                    tracking_number: shipment.tracking_number,
                    carrier: shipment.carrier || '',
                    label_url: shipment.guide_url || undefined,
                });
            }
        };

        const poll = async () => {
            if (cancelled) return;
            polls += 1;
            const id = pendingGuideShipmentIdRef.current;
            if (!id) {
                setPendingGuideCorrelationId(null);
                setError("No pudimos confirmar el estado de la gu\u00eda. Rev\u00edsala en la lista de env\u00edos antes de intentar de nuevo.");
                setLoading(false);
                return;
            }
            const res = await getShipmentByIdAction(id);
            if (cancelled) return;
            const shipment: any = res.success ? res.data : null;
            if (shipment?.tracking_number || shipment?.guide_url) {
                finishWithShipment(shipment);
                return;
            }
            if (shipment?.status === 'failed') {
                setPendingGuideCorrelationId(null);
                setError("La transportadora rechaz\u00f3 la gu\u00eda. Corrige los datos e intenta de nuevo.");
                setLoading(false);
                return;
            }
            if (shipment?.status === 'needs_verification') {
                setPendingGuideCorrelationId(null);
                setGuideBlocked(true);
                setError("No pudimos confirmar la respuesta de la transportadora. La gu\u00eda puede haberse creado: verif\u00edcala antes de generar otra para no duplicarla.");
                setLoading(false);
                return;
            }
            if (polls >= GUIDE_POLL_MAX_ATTEMPTS) {
                setPendingGuideCorrelationId(null);
                setGuideBlocked(true);
                setError("La gu\u00eda sigue gener\u00e1ndose. NO vuelvas a generarla: revisa la lista de env\u00edos en unos minutos para no crear una gu\u00eda duplicada.");
                setLoading(false);
                return;
            }
            setError("La gu\u00eda est\u00e1 tardando m\u00e1s de lo normal. Estamos verificando con la transportadora, no cierres esta ventana.");
            timer = setTimeout(poll, GUIDE_POLL_INTERVAL_MS);
        };

        let timer = setTimeout(poll, GUIDE_SSE_GRACE_MS);
        return () => {
            cancelled = true;
            clearTimeout(timer);
        };
    }, [pendingGuideCorrelationId]);

    useEffect(() => {
        if (!trackingNumber || !selectedCarrier) return;

        let cancelled = false;
        const loadFormats = async () => {
            try {
                const params = new URLSearchParams({ carrier: selectedCarrier });
                const response = await fetch(`/api/v1/shipments/guide-formats?${params}`, {
                    headers: { 'Content-Type': 'application/json' }
                });
                if (!cancelled && response.ok) {
                    const data = await response.json();
                    const allFormats = data.data || [];
                    const customFormats = allFormats.filter((f: any) => f.strategy === 'rebuild');
                    setGuideFormats(customFormats);
                    if (customFormats.length > 0) {
                        setSelectedGuideFormat(customFormats[0].code);
                    }
                }
            } catch (err) {
                console.error('Failed to load guide formats:', err);
            }
        };

        loadFormats();
        return () => { cancelled = true; };
    }, [trackingNumber, selectedCarrier]);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (formatDropdownRef.current && !formatDropdownRef.current.contains(event.target as Node)) {
                setShowGuideFormatDropdown(false);
            }
        };
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    useEffect(() => {
        if (!isOpen) {
            setCurrentStep(1);
            setStep1Data(null);
            setRates([]);
            setSelectedRate(null);
            setStep3Data(null);
            setGeneratedPdfUrl(null);
            setGeneratedShipmentId(null);
            setTrackingNumber(null);
            setSelectedCarrier(null);
            setOfficeCarrier(null);
            setError(null);
            setSuccess(null);
            setPendingCorrelationId(null);
            setPendingGuideCorrelationId(null);
            pendingStep1DataRef.current = null;
            setGuideFormats([]);
            setShowGuideFormatDropdown(false);
            setSelectedGuideFormat(null);
            step1Form.reset();
            step3Form.reset();
        }
    }, [isOpen]);

    const handleStep1Submit = async (data: Step1Values) => {
        if (!data.originDaneCode || !data.destDaneCode) {
            setError("⚠️ Por favor selecciona códigos DANE válidos para origen y destino");
            setLoading(false);
            return;
        }

        step3Form.setValue("originCrossStreet", data.originAddress);
        step3Form.setValue("destCrossStreet", data.destAddress);

        const errors = step1Form.formState.errors;
        if (Object.keys(errors).length > 0) {
            const fieldLabels: { [key: string]: string } = {
                originDaneCode: "Origen",
                originAddress: "dirección de origen",
                destDaneCode: "Destino",
                destAddress: "dirección de destino",
                weight: "peso",
                height: "altura",
                width: "ancho",
                length: "largo",
                description: "descripción",
                contentValue: "valor declarado",
                codPaymentMethod: "método de pago",
            };

            const errorFields: string[] = [];
            Object.entries(errors).forEach(([field, error]) => {
                const label = fieldLabels[field] || field;
                errorFields.push(label);
            });

            setError(`Completa: ${errorFields.join(', ')}`);
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
                const msg = response.message || "Error al cotizar";
                if (msg.toLowerCase().includes('no hay') || msg.toLowerCase().includes('sin transportador') || msg.toLowerCase().includes('integraci')) {
                    setError("No hay transportadoras disponibles para esta ruta. Verifica la integración con tus proveedores logísticos.");
                } else {
                    setError(msg);
                }
                setLoading(false);
                return;
            }

            const syncRates: EnvioClickRate[] = response.data?.rates || [];
            if (syncRates.length > 0) {
                setRates(syncRates);
                setStep1Data(data);
                setCurrentStep(2);
                setLoading(false);
                return;
            }

            setError("No hay transportadoras disponibles para esta ruta. Verifica la integración con tus proveedores logísticos.");
            setLoading(false);
            return;
        } catch (err: any) {
            setError(getActionError(err, "Error al cotizar envío"));
            setLoading(false);
        }
    };

    const handleRateSelection = (rate: EnvioClickRate) => {
        setSelectedRate(rate);
    };

    const handleStep2Continue = () => {
        if (selectedRate) setCurrentStep(3);
    };

    const handleStep3Submit = async (data: Step3Values) => {
        const errors = step3Form.formState.errors;

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

            const originErrors: string[] = [];
            const destErrors: string[] = [];

            Object.entries(errors).forEach(([field]) => {
                const label = fieldLabels[field] || field;
                if (field.startsWith('origin')) {
                    originErrors.push(label);
                } else if (field.startsWith('dest')) {
                    destErrors.push(label);
                }
            });

            const sections: string[] = [];
            if (originErrors.length > 0) {
                sections.push(`Remitente: ${originErrors.join(', ')}`);
            }
            if (destErrors.length > 0) {
                sections.push(`Destinatario: ${destErrors.join(', ')}`);
            }

            setError(`Completa: ${sections.join(' | ')}`);

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

        if (!selectedRate || !step3Data || !step1Data) return;
        const insuranceCost = (selectedRate.minimumInsurance ?? 0) + (step1Data.insurance ? (selectedRate.extraInsurance ?? 0) : 0);
        const codCarrierFee = selectedRate.cod ? (selectedRate.codCarrierFee ?? 0) : 0;
        const codMargin = selectedRate.cod ? (selectedRate.codProbabilityMargin ?? 0) : 0;
        const totalCost = selectedRate.flete + insuranceCost + codMargin;
        if (walletBalance !== null && walletBalance < totalCost) {
            setError(`Saldo insuficiente. Necesitas $${totalCost.toLocaleString()} pero tienes $${walletBalance.toLocaleString()}`);
            return;
        }

        setLoading(true);
        setError(null);

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
                codValue: (step1Data.codValue ?? 0) + codCarrierFee,
                includeGuideCost: step1Data.includeGuideCost,
                codPaymentMethod: step1Data.codPaymentMethod,
                totalCost: totalCost,
                codCarrierFee: codCarrierFee > 0 ? codCarrierFee : undefined,
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

            if (response.data?.data?.url) {
                const tracker = response.data.data.tracker;
                const carrier = (response.data?.data as any)?.carrier;
                const shipmentIdFromResponse = (response.data?.data as any)?.shipment_id;
                setGeneratedPdfUrl(response.data.data.url);
                if (shipmentIdFromResponse) setGeneratedShipmentId(shipmentIdFromResponse);
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

            pendingGuideShipmentIdRef.current = (response.data as any)?.shipment_id ?? null;
            setPendingGuideCorrelationId((response.data as any)?.correlation_id || null);
            setGuideGenerationRequested(true);
            setCurrentStep(4);
        } catch (err: any) {
            setError(getActionError(err, "Error al generar guía"));
            setLoading(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black/20 backdrop-blur-sm flex items-center justify-center z-50 p-2">
            <div className="relative bg-white dark:bg-gray-800 rounded-2xl shadow-xl flex flex-col overflow-hidden" style={{ width: '85%', maxHeight: '90vh' }}>
                {loading && currentStep === 4 && (
                    <BrandLoaderOverlay
                        title="Generando tu gu\u00eda..."
                        subtitle="Estamos confirmando con la transportadora. No cierres esta ventana ni vuelvas a generar."
                    />
                )}
                <div className="bg-white dark:bg-gray-800 border-b px-3 py-3 flex-shrink-0">
                    <div className="flex justify-between items-center mb-2">
                        <h2 className="text-2xl font-bold" style={{ color: 'var(--color-primary)' }}>Generar Guía de Envío</h2>
                        <button
                            onClick={onClose}
                            className="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:text-gray-200 dark:text-gray-200 text-2xl"
                        >
                            ×
                        </button>
                    </div>
                    <Stepper steps={STEPS} currentStep={currentStep} />
                </div>

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

                    {currentStep === 1 && (
                         
                        <form onSubmit={step1Form.handleSubmit(handleStep1Submit)} className="flex flex-col h-full overflow-hidden min-h-0" data-testid="step1-form">
                            <div className="flex-1 overflow-y-auto min-h-0 pr-3">
                                <div className="space-y-4">
                                    <div className="grid grid-cols-2 gap-4">
                                        <div className="shipment-section-origin">
                                            <div className="flex items-center justify-between">
                                                <div className="flex items-center gap-2">
                                                    <div className="shipment-section-origin-icon w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold">A</div>
                                                    <h3 className="shipment-section-origin-label">Origen</h3>
                                                </div>
                                                {originWarehouses.length > 0 && (
                                                    <select
                                                        className="text-[11px] px-1.5 py-0.5 border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-1"
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
                                                    className={`shipment-input ${step1Form.formState.errors.originDaneCode ? "shipment-input-error" : ""}`}
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

                                        <div className="shipment-section-destination">
                                            <div className="flex items-center gap-2">
                                                <div className="shipment-section-destination-icon w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold">B</div>
                                                <h3 className="shipment-section-destination-label">Destino</h3>
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
                                                    className={`shipment-input ${step1Form.formState.errors.destDaneCode ? "shipment-input-error" : ""}`}
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
                                                label="Descripci\u00f3n *"
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
                                            />
                                        </div>

                                        <div className="flex items-center gap-6 mt-3 pt-3 border-t border-gray-200 dark:border-gray-600/30">
                                            <label className="flex items-center gap-2 cursor-pointer">
                                                <input
                                                    type="checkbox"
                                                    {...step1Form.register("insurance")}
                                                    className="shipment-checkbox"
                                                />
                                                <span className="text-sm text-gray-700 dark:text-gray-300">Asegurar {'env\u00edo'}</span>
                                            </label>
                                            {(step1Form.watch("codValue") ?? 0) > 0 && (
                                                <>
                                                    <label className="flex items-center gap-2 cursor-pointer">
                                                        <input
                                                            type="checkbox"
                                                            {...step1Form.register("includeGuideCost")}
                                                            className="shipment-checkbox"
                                                        />
                                                        <span className="text-sm text-gray-700 dark:text-gray-300">Incluir costo {'gu\u00eda'} en contra entrega</span>
                                                    </label>
                                                    <div className="flex items-center gap-2">
                                                        <span className="text-sm text-gray-700 dark:text-gray-300">{'M\u00e9todo'} pago:</span>
                                                        <select
                                                            {...step1Form.register("codPaymentMethod")}
                                                            className="shipment-input px-2 py-1 text-sm"
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

                    {currentStep === 2 && (
                        <div className="flex flex-row h-full gap-3 overflow-hidden">
                            {order?.business_id && order.business_id > 0 && order.id && (
                                <div className="w-1/3 h-full flex-shrink-0 border border-gray-200 dark:border-gray-600 rounded-lg flex flex-col overflow-hidden">
                                    <div className="px-3 py-2 bg-gray-50 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-600 flex items-center gap-2">
                                        <button
                                            onClick={() => setMapViewMode('origin-destination')}
                                            className="px-2 py-1 text-xs rounded font-medium transition-colors text-white"
                                            style={{
                                                backgroundColor: mapViewMode === 'origin-destination' ? businessColors.primary : '#d1d5db',
                                                color: mapViewMode === 'origin-destination' ? '#ffffff' : '#374151'
                                            }}
                                        >
                                            Ruta
                                        </button>
                                        <button
                                            onClick={() => setMapViewMode('destination-only')}
                                            className="px-2 py-1 text-xs rounded font-medium transition-colors"
                                            style={{
                                                backgroundColor: mapViewMode === 'destination-only' ? businessColors.primary : '#d1d5db',
                                                color: mapViewMode === 'destination-only' ? '#ffffff' : '#374151'
                                            }}
                                        >
                                            Zona
                                        </button>
                                    </div>
                                    <div className="flex-1 overflow-hidden p-2 flex flex-col">
                                        <GeozoneMiniMap
                                            businessId={order.business_id}
                                            orderId={order.id}
                                            lat={order.shipping_lat ?? null}
                                            lng={order.shipping_lng ?? null}
                                            height="360px"
                                            origin={selectedOriginWarehouse ? {
                                                address: [selectedOriginWarehouse.street || selectedOriginWarehouse.address, selectedOriginWarehouse.city, selectedOriginWarehouse.state].filter(Boolean).join(', '),
                                                lat: selectedOriginWarehouse.latitude ?? null,
                                                lng: selectedOriginWarehouse.longitude ?? null,
                                            } : null}
                                            destination={{
                                                address: [order.shipping_street, order.shipping_city, order.shipping_state].filter(Boolean).join(', '),
                                            }}
                                            carrierRate={selectedCarrierProb?.delivery_rate ?? (selectedRate ? 0.85 : null)}
                                            carrierName={selectedRate?.carrier || null}
                                            carrierEstimated={selectedCarrierProb?.is_estimated || !selectedCarrierProb?.found || !selectedCarrierProb}
                                            viewMode={mapViewMode}
                                        />
                                    </div>
                                </div>
                            )}
                            <div className="w-2/3 flex flex-col overflow-y-auto">
                                <div className="pb-2">
                                    <h3 className="font-semibold text-lg text-gray-700 dark:text-gray-200 mb-2">
                                        Filtra por servicio / Transportadora
                                    </h3>
                                    <div className="flex items-center gap-2 flex-wrap">
                                        {(step1Data?.codValue ?? 0) > 0 && (
                                            <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-100 text-amber-700 border border-amber-300 dark:bg-amber-900/30 dark:text-amber-300 dark:border-amber-600">
                                                Contra Entrega - Solo opciones contra entrega
                                            </span>
                                        )}
                                        {officeCarrier && (
                                            <span className="shipment-badge-primary">
                                                Filtrado: {officeCarrier}
                                                <button
                                                    type="button"
                                                    onClick={() => setOfficeCarrier(null)}
                                                    className="ml-1 hover:opacity-75 font-bold leading-none"
                                                    title="Quitar filtro"
                                                >
                                                    x
                                                </button>
                                            </span>
                                        )}
                                    </div>
                                </div>

                                <div className="overflow-y-auto border border-gray-200 dark:border-gray-600 rounded-lg p-3 bg-white dark:bg-gray-800 flex-1">
                                {rates.length === 0 ? (
                                    <div className="flex items-center justify-center gap-3 py-10">
                                        <div className="shipment-spinner" style={{ width: 28, height: 28 }} />
                                        <span className="text-sm font-medium text-gray-600 dark:text-gray-300">Cargando cotizaciones...</span>
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
                                                        ? "No hay transportadoras disponibles con opci\u00f3n contra entrega para esta ruta"
                                                        : "No se encontraron cotizaciones para esta ruta"}
                                                </span>
                                            </div>
                                        );
                                    }
                                    return (
                                    <>
                                    {order?.quoted_shipping && (
                                        <div className="mb-3 flex items-center gap-2 rounded-xl border border-indigo-200 bg-indigo-50 px-3 py-2 dark:border-indigo-800 dark:bg-indigo-900/20">
                                            <span className="rounded-full bg-indigo-600 px-2 py-0.5 text-[10px] font-bold uppercase tracking-wide text-white">
                                                Cotizado
                                            </span>
                                            <span className="text-[12px] text-indigo-900 dark:text-indigo-200">
                                                El cliente eligio <strong>{order.quoted_shipping.title}</strong> por{' '}
                                                <strong>${order.quoted_shipping.price.toLocaleString('es-CO')}</strong> en el checkout.
                                                Genera la guia sobre esa transportadora.
                                            </span>
                                        </div>
                                    )}
                                    <div className="grid grid-cols-2 gap-4 auto-rows-max">
                                        {filteredRates.map((rate) => {
                                            const minIns = rate.minimumInsurance ?? 0;
                                            const extraIns = rate.extraInsurance ?? 0;
                                            const insuranceCost = minIns + (step1Data?.insurance ? extraIns : 0);
                                            const codMargin = rate.cod ? (rate.codProbabilityMargin ?? 0) : 0;
                                            const codCarrierFee = rate.cod ? (rate.codCarrierFee ?? 0) : 0;
                                            const totalCost = rate.flete + insuranceCost + codMargin;
                                            const isInsured = step1Data?.insurance === true && insuranceCost > 0;
                                            const isCOD = rate.cod;

                                            const allCosts = filteredRates.map(r => r.flete + (r.minimumInsurance ?? 0) + (step1Data?.insurance ? (r.extraInsurance ?? 0) : 0) + (r.cod ? (r.codProbabilityMargin ?? 0) : 0));
                                            const minCost = Math.min(...allCosts);
                                            const minDays = Math.min(...filteredRates.map(r => r.deliveryDays || 999));

                                            const getCarrierProb = (carrierName: string) =>
                                                carrierProbabilities.find(p => normalizeCarrierKey(p.carrier || '') === normalizeCarrierKey(carrierName));

                                            const allEffs = filteredRates.map(r => getCarrierProb(r.carrier)?.delivery_rate ?? 0);
                                            const maxEff = Math.max(...allEffs);

                                            const isFastest = rate.deliveryDays === minDays;
                                            const isCheapest = totalCost === minCost;
                                            const eff = getCarrierProb(rate.carrier)?.delivery_rate ?? 0;
                                            const isMostEffective = eff === maxEff && maxEff > 0 && eff > 0;

                                            const badges = [];
                                            if (isFastest) badges.push('fastest');
                                            if (isCheapest) badges.push('cheapest');
                                            if (isMostEffective) badges.push('mostEffective');

                                            const isSameDay = rate.deliveryDays === 0;
                                            const isQuotedInCheckout = !!quotedCarrierKey &&
                                                normalizeCarrierKey(rate.carrier) === quotedCarrierKey;
                                            const hasSpecialBadge = isQuotedInCheckout || isMostEffective || isCheapest || isFastest || isSameDay;

                                            let badgeColor = businessColors.quaternary;
                                            let badgeLabel = '';
                                            let borderColor = businessColors.tertiary;

                                            if (isQuotedInCheckout) {
                                                badgeColor = '#4f46e5';
                                                badgeLabel = 'COTIZADA EN EL CHECKOUT';
                                                borderColor = '#4f46e5';
                                            } else if (isMostEffective) {
                                                badgeColor = businessColors.primary;
                                                badgeLabel = 'RECOMENDADO';
                                            } else if (isCheapest) {
                                                badgeColor = businessColors.secondary;
                                                badgeLabel = 'MÁS ECONÓMICA';
                                            } else if (isSameDay) {
                                                badgeColor = businessColors.quaternary;
                                                badgeLabel = 'MISMO DÍA';
                                            } else if (isFastest) {
                                                badgeColor = businessColors.tertiary;
                                                badgeLabel = 'MÁS RÁPIDA';
                                            }

                                            const isSelected = selectedRate?.idRate === rate.idRate;

                                            return (
                                                <div
                                                    key={rate.idRate}
                                                    onClick={() => handleRateSelection(rate)}
                                                    className={`relative border-2 rounded-3xl transition-all cursor-pointer ${
                                                        hasSpecialBadge ? 'shadow-sm hover:shadow-lg hover:-translate-y-0.5' : 'shadow-sm hover:shadow-lg hover:-translate-y-0.5'
                                                    }`}
                                                    style={{
                                                        backgroundColor: isSelected ? `${businessColors.quaternary}40` : `${businessColors.tertiary}08`,
                                                        borderColor: hasSpecialBadge ? borderColor : selectedRate?.idRate === rate.idRate ? businessColors.tertiary : '#d1d5db',
                                                        padding: '18px',
                                                        display: 'flex',
                                                        flexDirection: 'column',
                                                        gap: '0',
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

                                                            <div style={{ display: 'flex', flexDirection: 'column', gap: '4px', fontSize: '14px', color: '#3b4248', marginTop: '8px', textAlign: 'center', fontWeight: 700 }}>
                                                                <div>Costo Guía: ${totalCost.toLocaleString()}</div>
                                                                <div>Guía + Comisión: ${(totalCost + (codCarrierFee || 0)).toLocaleString()}</div>
                                                                {isInsured ? (
                                                                    <div style={{ color: '#059669', fontSize: '13px' }}>
                                                                        (Seguro: ${insuranceCost.toLocaleString()})
                                                                    </div>
                                                                ) : (
                                                                    <div style={{ color: '#6b757c', fontSize: '13px' }}>(Sin asegurar)</div>
                                                                )}
                                                                {isCOD && codCarrierFee > 0 && (
                                                                    <div style={{ color: '#0891b2', fontSize: '13px' }}>
                                                                        Comisión carrier: ${codCarrierFee.toLocaleString()}
                                                                    </div>
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

                                                            {isCOD && (
                                                                <div style={{
                                                                    display: 'inline-flex',
                                                                    alignItems: 'center',
                                                                    justifyContent: 'center',
                                                                    padding: '4px 10px',
                                                                    borderRadius: '8px',
                                                                    fontSize: '10px',
                                                                    fontWeight: 600,
                                                                    color: '#6b21a8',
                                                                    backgroundColor: '#f3e8ff',
                                                                    border: '1px solid #e9d5ff',
                                                                    margin: '2px auto 0',
                                                                }}>
                                                                    Contra Entrega
                                                                </div>
                                                            )}
                                                        </div>

                                                        <div style={{
                                                            width: '1px',
                                                            backgroundColor: '#e5e7eb',
                                                            margin: '12px 0',
                                                            flexShrink: 0,
                                                        }}></div>

                                                        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', flex: '0 0 50%', padding: '0 12px', borderLeft: 'none', justifyContent: 'center', alignItems: 'stretch' }}>
                                                            {order?.business_id && order.business_id > 0 && order.id && (
                                                                <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                                                                    <CarrierEffectivenessRates
                                                                        businessId={order.business_id}
                                                                        orderId={order.id}
                                                                        carrier={rate.carrier}
                                                                    />
                                                                </div>
                                                            )}

                                                        </div>
                                                    </div>
                                                </div>
                                            );
                                        })}
                                    </div>
                                    </>
                                    );
                                })()}
                                </div>
                            </div>
                        </div>
                    )}

                    {currentStep === 3 && (
                        <form onSubmit={step3Form.handleSubmit(handleStep3Submit)} className="flex flex-col h-full overflow-hidden">
                            <div className="overflow-y-auto flex-1 space-y-3 pr-1">
                                <div className="grid grid-cols-1 lg:grid-cols-2 gap-3">
                                    <div className="shipment-section-origin rounded-xl p-4">
                                        <div className="flex items-center gap-2 mb-3 pb-2 border-b" style={{ borderColor: 'color-mix(in oklab, var(--color-primary) 30%, transparent)' }}>
                                            <div className="shipment-section-origin-icon w-7 h-7 rounded-lg flex items-center justify-center">
                                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
                                            </div>
                                            <h3 className="shipment-section-origin-label text-sm">Remitente (Origen)</h3>
                                        </div>

                                        <p className="text-[10px] text-gray-500 dark:text-gray-400 uppercase tracking-wider font-semibold mb-1">{'Direcci\u00f3n'}</p>
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
                                            <Input compact label="Tel\u00e9fono *" {...step3Form.register("originPhone")} error={step3Form.formState.errors.originPhone?.message} placeholder="3224098631" />
                                            <Input compact label="Correo *" type="email" {...step3Form.register("originEmail")} error={step3Form.formState.errors.originEmail?.message} placeholder="correo@ejemplo.com" />
                                        </div>
                                    </div>

                                    <div className="shipment-section-destination rounded-xl p-4">
                                        <div className="flex items-center gap-2 mb-3 pb-2 border-b" style={{ borderColor: 'color-mix(in oklab, var(--color-secondary) 30%, transparent)' }}>
                                            <div className="shipment-section-destination-icon w-7 h-7 rounded-lg flex items-center justify-center">
                                                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="1" y="3" width="15" height="13"/><polygon points="16 8 20 8 23 11 23 16 16 16 16 8"/><circle cx="5.5" cy="18.5" r="2.5"/><circle cx="18.5" cy="18.5" r="2.5"/></svg>
                                            </div>
                                            <h3 className="shipment-section-destination-label text-sm">Destinatario (Destino)</h3>
                                        </div>

                                        <p className="text-[10px] text-gray-500 dark:text-gray-400 uppercase tracking-wider font-semibold mb-1">{'Direcci\u00f3n'}</p>
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
                                            <Input compact label="Tel\u00e9fono *" {...step3Form.register("destPhone")} error={step3Form.formState.errors.destPhone?.message} placeholder="3224098631" />
                                            <Input compact label="Correo *" type="email" {...step3Form.register("destEmail")} error={step3Form.formState.errors.destEmail?.message} placeholder="correo@ejemplo.com" />
                                        </div>
                                    </div>
                                </div>

                                <div className="border border-gray-200 dark:border-gray-600 rounded-xl bg-gray-50/60 dark:bg-gray-700/30 p-4">
                                    <p className="text-[10px] text-gray-500 dark:text-gray-400 uppercase tracking-wider font-semibold mb-2">Opciones adicionales</p>
                                    <div className="grid grid-cols-2 gap-2">
                                        <Input compact label="Mi referencia de env\u00edo" {...step3Form.register("myShipmentReference")} error={step3Form.formState.errors.myShipmentReference?.message} placeholder="Orden 5649" />
                                        <Input compact label="N\u00famero de orden externo" {...step3Form.register("external_order_id")} error={step3Form.formState.errors.external_order_id?.message} placeholder="ORD345678" />
                                    </div>
                                    <label className="flex items-center space-x-2 mt-2">
                                        <input type="checkbox" {...step3Form.register("requestPickup")} className="rounded w-5 h-5" />
                                        <span className="text-sm font-medium">Solicitar {'recolecci\u00f3n'}</span>
                                    </label>
                                </div>
                            </div>
                        </form>
                    )}

                    {currentStep === 4 && selectedRate && (
                        <div className="flex flex-col h-full w-full overflow-hidden gap-3">
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
                                            <div className="shipment-cost-amount text-2xl">
                                                ${(selectedRate.flete + (selectedRate.minimumInsurance ?? 0) + (step1Data?.insurance ? (selectedRate.extraInsurance ?? 0) : 0) + (selectedRate.cod ? (selectedRate.codProbabilityMargin ?? 0) : 0)).toLocaleString()}
                                            </div>
                                            <div className="text-xs text-gray-500 dark:text-gray-400 mt-1 text-right leading-tight">
                                                Guía: ${selectedRate.flete.toLocaleString()}<br />
                                                Seg. obligatorio: ${(selectedRate.minimumInsurance ?? 0).toLocaleString()}<br />
                                                Seg. adicional: ${(selectedRate.extraInsurance ?? 0).toLocaleString()} <span className={step1Data?.insurance ? 'text-emerald-600' : 'text-gray-400'}>{step1Data?.insurance ? '(incluido)' : '(no incluido)'}</span>
                                                {selectedRate.cod && (selectedRate.codCarrierFee ?? 0) > 0 && (<><br /><span className="text-cyan-700 dark:text-cyan-400">Comisión carrier: ${(selectedRate.codCarrierFee ?? 0).toLocaleString()}</span></>)}
                                            </div>
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

                            <div className="flex-1 min-h-0 overflow-y-auto overflow-x-hidden space-y-3 pr-3" style={{ maxHeight: 'calc(85vh - 280px)' }}>
                                <div>
                                    <h4 className="font-medium text-gray-700 dark:text-gray-200 dark:text-gray-200 mb-3">Selecciona tu método de pago</h4>
                                    <div className={`grid gap-2 ${generatedPdfUrl ? 'grid-cols-2' : 'grid-cols-1'}`}>
                                        <div className="border-2 rounded-lg p-2" style={{ borderColor: 'var(--color-primary)', background: 'color-mix(in oklab, var(--color-primary) 10%, transparent)' }}>
                                            <div className="flex items-center justify-center mb-2">
                                                <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24" style={{ color: 'var(--color-primary)' }}>
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
                                                </svg>
                                            </div>
                                            <div className="text-center font-semibold">Monedero</div>
                                            <div className="text-center text-sm text-gray-600 dark:text-gray-300">
                                                ${walletBalance?.toLocaleString() || 0}
                                            </div>
                                        </div>

                                        {generatedPdfUrl && (
                                            <div className="shipment-success-container">
                                                <div className="shipment-success-icon">
                                                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
                                                    </svg>
                                                </div>
                                                <div className="text-center min-w-0">
                                                    <p className="shipment-success-text">¡Guía generada exitosamente!</p>
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
                                                    <div ref={formatDropdownRef} className="relative">
                                                        {guideFormats.length > 0 && (
                                                        <button
                                                            onClick={() => setShowGuideFormatDropdown(!showGuideFormatDropdown)}
                                                            className="flex items-center justify-center gap-1 w-full py-1.5 px-2 bg-green-600 hover:bg-green-700 text-white font-semibold rounded-lg transition-colors text-xs"
                                                        >
                                                            <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                                                            </svg>
                                                            Guías personalizadas
                                                            <svg className={`w-3 h-3 transition-transform ${showGuideFormatDropdown ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                                                            </svg>
                                                        </button>
                                                        )}
                                                        {showGuideFormatDropdown && guideFormats.length > 0 && (
                                                            <div className="absolute z-20 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md shadow-lg overflow-hidden">
                                                                {guideFormats.map((format) => (
                                                                    <button
                                                                        key={format.code}
                                                                        type="button"
                                                                        disabled={!generatedShipmentId}
                                                                        onClick={() => {
                                                                            if (!generatedShipmentId) return;
                                                                            window.open(`/internal/shipment-guide/${generatedShipmentId}?format=${encodeURIComponent(format.code)}`, '_blank');
                                                                            setSelectedGuideFormat(format.code);
                                                                            setShowGuideFormatDropdown(false);
                                                                        }}
                                                                        className="w-full flex items-center justify-between px-3 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300 text-xs border-b dark:border-gray-700 last:border-b-0 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                                                                    >
                                                                        <div>
                                                                            <div className="font-medium">{format.code.toUpperCase()}</div>
                                                                            {format.width_cm && format.height_cm && (
                                                                                <div className="text-gray-500 dark:text-gray-400 text-[10px]">
                                                                                    {format.width_cm} x {format.height_cm} cm
                                                                                </div>
                                                                            )}
                                                                        </div>
                                                                        {selectedGuideFormat === format.code && (
                                                                            <svg className="w-3 h-3 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                                                                                <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                                                                            </svg>
                                                                        )}
                                                                    </button>
                                                                ))}
                                                            </div>
                                                        )}
                                                    </div>
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

                <div className="bg-white dark:bg-gray-800 border-t px-3 py-3 flex-shrink-0 flex justify-between items-center gap-3">
                    {(currentStep === 2 || currentStep === 3 || currentStep === 4) && !generatedPdfUrl && (
                        <Button
                            variant="outline"
                            onClick={() => setCurrentStep(currentStep - 1)}
                            disabled={loading}
                        >
                            Atrás
                        </Button>
                    )}

                    {(currentStep === 1 || (currentStep === 4 && generatedPdfUrl)) && (
                        <div />
                    )}

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
                                        ([field]) => fieldLabels[field] || field
                                    );
                                    setError(`Completa: ${errorFields.join(', ')}`);
                                })();
                            }}
                            disabled={loading}
                            className="shipment-btn-primary"
                        >
                            {loading ? "Cotizando..." : "Siguiente"}
                        </Button>
                    )}

                    {currentStep === 2 && (
                        <div className="flex items-center gap-3">
                            {!selectedRate && (
                                <span className="text-sm text-gray-600 dark:text-gray-300 italic">
                                    Selecciona una transportadora
                                </span>
                            )}
                            {selectedRate && (
                                <span className="text-sm text-gray-700 dark:text-gray-200">
                                    Transportadora: <strong>{selectedRate.carrier}</strong>
                                </span>
                            )}
                            <Button
                                className="shipment-btn-primary"
                                onClick={handleStep2Continue}
                                disabled={!selectedRate}
                            >
                                Continuar
                            </Button>
                        </div>
                    )}

                    {currentStep === 3 && (
                        <Button
                            className="shipment-btn-primary"
                            onClick={async () => {
                                const isValid = await step3Form.trigger();
                                if (isValid) {
                                    const data = step3Form.getValues();
                                    setStep3Data(data);
                                    setCurrentStep(4);
                                }
                            }}
                            disabled={loading || Object.keys(step3Form.formState.errors).length > 0}
                            title={Object.keys(step3Form.formState.errors).length > 0 ? "Completa todos los campos requeridos" : ""}
                        >
                            {Object.keys(step3Form.formState.errors).length > 0
                                ? `⚠️ ${Object.keys(step3Form.formState.errors).length} campo(s) incompleto(s)`
                                : "Siguiente"
                            }
                        </Button>
                    )}

                    {currentStep === 4 && !generatedPdfUrl && (
                        <Button
                            onClick={handleFinalGenerate}
                            disabled={loading || guideBlocked}
                            className="shipment-btn-secondary"
                        >
                            {loading ? "Generando..." : guideBlocked ? "Verifica la gu\u00eda antes de reintentar" : "Pagar gu\u00edas"}
                        </Button>
                    )}
                </div>
            </div>
        </div>
    );
}
