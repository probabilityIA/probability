import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Input, Button, Select } from "@/shared/ui";
import { ShipmentApiRepository } from "@/services/modules/shipments/infra/repository/api-repository";
import { EnvioClickQuoteRequest, EnvioClickRate } from "@/services/modules/shipments/domain/types";
import { Order } from "@/services/modules/orders/domain/types";
import { getWalletBalanceAction } from "@/services/modules/wallet/infra/actions";
import { getAIRecommendationAction } from "@/services/modules/orders/infra/actions";
import danes from "@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json";

interface ShipmentGuideModalProps {
    isOpen: boolean;
    onClose: () => void;
    order?: Order;
    recommendedCarrier?: string;
    onGuideGenerated?: (trackingNumber: string) => void;
}

// Zod Schema (Same as ShippingForm but adapted if needed)
const addressSchema = z.object({
    company: z.string().optional(),
    firstName: z.string().min(1, "Nombre es requerido"),
    lastName: z.string().min(1, "Apellido es requerido"),
    email: z.string().email("Email inválido"),
    phone: z.string().min(1, "Teléfono es requerido"),
    address: z.string().min(1, "Dirección es requerida"),
    suburb: z.string().min(1, "Barrio es requerido"),
    crossStreet: z.string().optional(),
    reference: z.string().optional(),
    daneCode: z.string().min(1, "Código DANE es requerido"),
});

const formSchema = z.object({
    origin: addressSchema,
    destination: addressSchema,
    packageSize: z.enum(["small", "medium", "large", "custom"]),
    customPackage: z.object({
        weight: z.number().min(0.1),
        height: z.number().min(1),
        width: z.number().min(1),
        length: z.number().min(1),
    }).optional(),
    description: z.string().min(1, "Descripción es requerida"),
    contentValue: z.number().min(0, "Valor declarado debe ser positivo"),
    requestPickup: z.boolean(),
    insurance: z.boolean(),
    codPaymentMethod: z.string().min(1, "Método de pago requerido").max(10, "Máximo 10 caracteres"),
    external_order_id: z.string().optional(),
    myShipmentReference: z.string().optional(),
});

type FormValues = z.infer<typeof formSchema>;

const PACKAGE_SIZES = {
    small: { weight: 1, height: 10, width: 10, length: 10, label: "Pequeño (1kg - 10x10x10)" },
    medium: { weight: 5, height: 30, width: 30, length: 30, label: "Mediano (5kg - 30x30x30)" },
    large: { weight: 10, height: 50, width: 50, length: 50, label: "Grande (10kg - 50x50x50)" },
    custom: { weight: 0, height: 0, width: 0, length: 0, label: "Personalizado" },
};

export default function ShipmentGuideModal({ isOpen, onClose, order, recommendedCarrier, onGuideGenerated }: ShipmentGuideModalProps) {
    const [loading, setLoading] = useState(false);
    const [success, setSuccess] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [rates, setRates] = useState<EnvioClickRate[]>([]);
    const [selectedRate, setSelectedRate] = useState<EnvioClickRate | null>(null);
    const [hasQuoted, setHasQuoted] = useState(false);
    const [generatedPdfUrl, setGeneratedPdfUrl] = useState<string | null>(null);
    const [walletBalance, setWalletBalance] = useState<number | null>(null);
    const [aiAnalysis, setAiAnalysis] = useState<{ recommended_carrier: string; reasoning: string } | null>(null);
    const [loadingAI, setLoadingAI] = useState(false);
    const [showBalanceModal, setShowBalanceModal] = useState(false);
    const [insufficientBalanceInfo, setInsufficientBalanceInfo] = useState<{ balance: number; cost: number } | null>(null);

    // DANE options for select
    const daneOptions = Object.entries(danes).map(([code, data]: [string, any]) => ({
        value: code,
        label: `${data.ciudad} (${data.departamento})`
    })).sort((a, b) => a.label.localeCompare(b.label));

    const {
        register,
        handleSubmit,
        watch,
        setValue,
        formState: { errors },
    } = useForm<FormValues>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            origin: {
                company: "", firstName: "", lastName: "", email: "", phone: "",
                address: "", suburb: "", crossStreet: "", reference: "", daneCode: ""
            },
            destination: {
                company: "", firstName: "", lastName: "", email: "", phone: "", address: "", suburb: "", crossStreet: "", reference: "", daneCode: ""
            },
            packageSize: "medium",
            insurance: true,
            requestPickup: false,
            contentValue: 0,
            codPaymentMethod: "cash",
            description: "E-commerce Order",
        },
    });

    const packageSize = watch("packageSize");

    // Fetch wallet balance on open
    useEffect(() => {
        if (isOpen) {
            getWalletBalanceAction().then(res => {
                if (res.success && res.data) setWalletBalance(res.data.Balance);
            });
        }
    }, [isOpen]);

    // Pre-fill form with order data on open
    useEffect(() => {
        if (isOpen && order) {
            setValue("destination.company", order.customer_name);
            setValue("destination.firstName", order.customer_name.split(" ")[0] || "");
            setValue("destination.lastName", order.customer_name.split(" ").slice(1).join(" ") || ".");
            setValue("destination.email", order.customer_email);
            setValue("destination.phone", order.customer_phone);
            setValue("destination.address", order.shipping_street);
            setValue("destination.suburb", order.shipping_state || "");

            // Try to find DANE code by city name
            const city = (order.shipping_city || "").toUpperCase();
            const foundDane = Object.entries(danes).find(([_, data]: [string, any]) => data.ciudad.includes(city));
            if (foundDane) {
                setValue("destination.daneCode", foundDane[0]);
            } else {
                setValue("destination.daneCode", "11001000"); // Fallback Bogota
            }

            setValue("contentValue", order.total_amount);
            setValue("description", "Order " + order.order_number);
            setValue("external_order_id", order.order_number);
            setValue("myShipmentReference", "Orden " + order.internal_number);

            if (order.weight && order.weight > 0) {
                setValue("packageSize", "custom");
                setValue("customPackage.weight", order.weight);
                setValue("customPackage.height", order.height || 10);
                setValue("customPackage.width", order.width || 10);
                setValue("customPackage.length", order.length || 10);
            }
        }
    }, [isOpen, order, setValue]);

    const buildPayload = (data: FormValues, idRate: number = 0): EnvioClickQuoteRequest => {
        const pkg =
            data.packageSize === "custom" && data.customPackage
                ? data.customPackage
                : PACKAGE_SIZES[data.packageSize];

        return {
            idRate: idRate,
            myShipmentReference: data.myShipmentReference || `REF-${Date.now()}`,
            external_order_id: data.external_order_id || `EXT-${Date.now()}`,
            requestPickup: data.requestPickup,
            pickupDate: new Date().toISOString().split("T")[0],
            insurance: data.insurance,
            description: data.description,
            contentValue: Number(data.contentValue),
            codValue: Number(data.contentValue),
            includeGuideCost: false,
            codPaymentMethod: data.codPaymentMethod,
            packages: [
                {
                    weight: Number(pkg.weight),
                    height: Number(pkg.height),
                    width: Number(pkg.width),
                    length: Number(pkg.length),
                },
            ],
            origin: {
                ...data.origin,
                company: data.origin.company || "",
                crossStreet: data.origin.crossStreet || "",
                reference: data.origin.reference || ""
            },
            destination: {
                ...data.destination,
                company: data.destination.company || "",
                crossStreet: data.destination.crossStreet || "",
                reference: data.destination.reference || ""
            },
        };
    };

    const handleQuote = async (data: FormValues) => {
        setLoading(true);
        setError(null);
        setRates([]);
        setSelectedRate(null);
        setHasQuoted(false);

        try {
            const payload = buildPayload(data);
            const repo = new ShipmentApiRepository();
            const res = await repo.quoteShipment(payload);

            if (res.data && res.data.rates && res.data.rates.length > 0) {
                setRates(res.data.rates);
                setHasQuoted(true);

                // Perform AI Analysis on real rates
                const originDane = danes[data.origin.daneCode as keyof typeof danes];
                const destDane = danes[data.destination.daneCode as keyof typeof danes];

                if (originDane && destDane) {
                    setLoadingAI(true);
                    try {
                        const aiRes = await getAIRecommendationAction(destDane.ciudad, destDane.departamento);
                        if (aiRes) {
                            // Validate that the recommended carrier exists in the quotes
                            const carrierExists = res.data.rates.some((r: any) =>
                                r.carrier.toLowerCase() === aiRes.recommended_carrier.toLowerCase()
                            );

                            if (carrierExists) {
                                setAiAnalysis({
                                    recommended_carrier: aiRes.recommended_carrier,
                                    reasoning: aiRes.reasoning
                                });
                            } else {
                                console.warn(`AI recommended ${aiRes.recommended_carrier} but it's not in available quotes.`);
                                setAiAnalysis(null);
                            }
                        }
                    } catch (e) {
                        console.error("AI Error:", e);
                    } finally {
                        setLoadingAI(false);
                    }
                }

            } else {
                setError("No se encontraron cotizaciones.");
            }
        } catch (err: any) {
            const msg = err.message || "";
            if (msg.toLowerCase().includes("no tiene suficiente crédito")) {
                setError("Saldo insuficiente");
            } else {
                setError(msg || "Error consultando cotizaciones");
            }
        } finally {
            setLoading(false);
        }
    };

    const handleGenerate = async (data: FormValues) => {
        if (selectedRate && walletBalance !== null && walletBalance < selectedRate.flete) {
            setInsufficientBalanceInfo({ balance: walletBalance, cost: selectedRate.flete });
            setShowBalanceModal(true);
            return;
        }

        if (!selectedRate) {
            setError("Debes seleccionar una cotización");
            return;
        }

        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            const payload = buildPayload(data, selectedRate.idRate);
            const repo = new ShipmentApiRepository();
            const res = await repo.generateGuide(payload);

            // API returns 'tracker' and 'url'
            const tracking = res.data.tracker;
            const pdfUrl = res.data.url;

            setSuccess(`Guía generada exitosamente! Tracking: ${tracking}`);

            // We can store the URL in a local state if we want to show it in the success UI specifically
            // For now, we'll append it to the success message or handle it in the UI render
            setGeneratedPdfUrl(pdfUrl);

            if (onGuideGenerated) onGuideGenerated(tracking);
        } catch (err: any) {
            const msg = err.message || "";
            if (msg.toLowerCase().includes("no tiene suficiente crédito")) {
                setError("Saldo insuficiente");
            } else {
                setError(msg || "Error generando guía");
            }
        } finally {
            setLoading(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
            <div className="bg-white rounded-xl shadow-xl w-full max-w-4xl overflow-hidden max-h-[90vh] flex flex-col">
                <div className="p-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
                    <h3 className="text-lg font-bold text-gray-800">
                        {order ? `Generar Guía de Envío para Orden #${order.order_number}` : 'Generar Guía de Envío (Prueba)'}
                    </h3>
                    <button onClick={onClose} className="text-gray-400 hover:text-gray-600">✕</button>
                </div>

                <div className="p-6 overflow-y-auto flex-1">
                    {error && (
                        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
                            <p className="font-bold">Error</p>
                            <p>{error}</p>
                        </div>
                    )}

                    {showBalanceModal && insufficientBalanceInfo && (
                        <div className="fixed inset-0 z-[60] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
                            <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md overflow-hidden animate-in fade-in zoom-in duration-300">
                                <div className="bg-red-600 p-6 text-center">
                                    <div className="w-16 h-16 bg-white/20 rounded-full flex items-center justify-center mx-auto mb-4">
                                        <span className="text-3xl">⚠️</span>
                                    </div>
                                    <h3 className="text-xl font-bold text-white">Saldo Insuficiente</h3>
                                </div>
                                <div className="p-8 text-center">
                                    <p className="text-gray-600 mb-6">
                                        No cuentas con saldo suficiente en tu billetera local para generar esta guía.
                                    </p>
                                    <div className="bg-gray-50 rounded-xl p-4 mb-8 grid grid-cols-2 gap-4">
                                        <div className="text-left">
                                            <p className="text-xs text-gray-500 uppercase font-bold mb-1">Tu Saldo</p>
                                            <p className="text-lg font-bold text-gray-900">${insufficientBalanceInfo.balance.toLocaleString()}</p>
                                        </div>
                                        <div className="text-left border-l pl-4">
                                            <p className="text-xs text-gray-500 uppercase font-bold mb-1">Costo Guía</p>
                                            <p className="text-lg font-bold text-red-600">${insufficientBalanceInfo.cost.toLocaleString()}</p>
                                        </div>
                                    </div>
                                    <div className="flex flex-col gap-3">
                                        <Button
                                            className="w-full bg-red-600 hover:bg-red-700 h-12 text-lg"
                                            onClick={() => window.location.href = '/wallet'}
                                        >
                                            Recargar Billetera
                                        </Button>
                                        <Button
                                            variant="secondary"
                                            className="w-full h-12"
                                            onClick={() => setShowBalanceModal(false)}
                                        >
                                            Volver
                                        </Button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}

                    {success && (
                        <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded mb-4 text-center">
                            <p className="font-bold text-lg mb-2">¡Éxito!</p>
                            <p>{success}</p>

                            {generatedPdfUrl && (
                                <div className="mt-4">
                                    <a
                                        href={generatedPdfUrl}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="inline-block bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition-colors"
                                    >
                                        📄 Ver Guía PDF
                                    </a>
                                </div>
                            )}

                            <Button className="mt-4" onClick={onClose}>Cerrar</Button>
                        </div>
                    )}

                    {!success && (
                        <form onSubmit={handleSubmit(hasQuoted ? handleGenerate : handleQuote)} className="space-y-6">
                            {/* Origin Section */}
                            <section>
                                <h3 className="text-md font-semibold mb-3 text-gray-700 border-b pb-1">Origen</h3>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                    <Input label="Empresa" {...register("origin.company")} error={errors.origin?.company?.message} />
                                    <Input label="Nombre" {...register("origin.firstName")} error={errors.origin?.firstName?.message} />
                                    <Input label="Apellido" {...register("origin.lastName")} error={errors.origin?.lastName?.message} />
                                    <Input label="Email" {...register("origin.email")} error={errors.origin?.email?.message} />
                                    <Input label="Teléfono" {...register("origin.phone")} error={errors.origin?.phone?.message} />
                                    <Input label="Dirección" {...register("origin.address")} error={errors.origin?.address?.message} />
                                    <Input label="Barrio" {...register("origin.suburb")} error={errors.origin?.suburb?.message} />
                                    <div className="md:col-span-2 grid grid-cols-1 md:col-span-2 gap-3">
                                        <Input label="Cruzamiento" {...register("origin.crossStreet")} error={errors.origin?.crossStreet?.message} />
                                        <Input label="Referencia" {...register("origin.reference")} error={errors.origin?.reference?.message} />
                                        <div className="md:col-span-2">
                                            <label className="block text-sm font-medium text-gray-700 mb-1">Ciudad (Buscable)</label>
                                            <select
                                                {...register("origin.daneCode")}
                                                className="w-full border-gray-300 rounded-md shadow-sm p-2 border text-sm"
                                            >
                                                <option value="">Buscar ciudad...</option>
                                                {daneOptions.map(opt => <option key={opt.value} value={opt.value}>{opt.label}</option>)}
                                            </select>
                                        </div>
                                    </div>
                                </div>
                            </section>

                            {/* Destination Section */}
                            <section>
                                <h3 className="text-md font-semibold mb-3 text-gray-700 border-b pb-1">Destino</h3>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                    <Input label="Empresa" {...register("destination.company")} error={errors.destination?.company?.message} />
                                    <Input label="Nombre" {...register("destination.firstName")} error={errors.destination?.firstName?.message} />
                                    <Input label="Apellido" {...register("destination.lastName")} error={errors.destination?.lastName?.message} />
                                    <Input label="Email" {...register("destination.email")} error={errors.destination?.email?.message} />
                                    <Input label="Teléfono" {...register("destination.phone")} error={errors.destination?.phone?.message} />
                                    <Input label="Dirección" {...register("destination.address")} error={errors.destination?.address?.message} />
                                    <Input label="Barrio" {...register("destination.suburb")} error={errors.destination?.suburb?.message} />
                                    <div className="md:col-span-2 grid grid-cols-1 md:col-span-2 gap-3">
                                        <Input label="Cruzamiento" {...register("destination.crossStreet")} error={errors.destination?.crossStreet?.message} />
                                        <Input label="Referencia" {...register("destination.reference")} error={errors.destination?.reference?.message} />
                                        <div className="md:col-span-2">
                                            <label className="block text-sm font-medium text-gray-700 mb-1">Ciudad (Buscable)</label>
                                            <select
                                                {...register("destination.daneCode")}
                                                className="w-full border-gray-300 rounded-md shadow-sm p-2 border text-sm"
                                            >
                                                <option value="">Buscar ciudad...</option>
                                                {daneOptions.map(opt => <option key={opt.value} value={opt.value}>{opt.label}</option>)}
                                            </select>
                                        </div>
                                    </div>
                                </div>
                            </section>

                            {/* Package Section */}
                            <section>
                                <h3 className="text-md font-semibold mb-3 text-gray-700 border-b pb-1">Paquete y Detalles</h3>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                    <div className="col-span-2">
                                        <label className="block text-sm font-medium text-gray-700 mb-1">Tamaño del Paquete</label>
                                        <select
                                            {...register("packageSize")}
                                            className="w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm p-2 border"
                                        >
                                            {Object.entries(PACKAGE_SIZES).map(([key, val]) => (
                                                <option key={key} value={key}>{val.label}</option>
                                            ))}
                                        </select>
                                    </div>

                                    {packageSize === "custom" && (
                                        <>
                                            <Input label="Peso (kg)" type="number" step="0.1" {...register("customPackage.weight", { valueAsNumber: true })} error={errors.customPackage?.weight?.message} />
                                            <Input label="Alto (cm)" type="number" {...register("customPackage.height", { valueAsNumber: true })} error={errors.customPackage?.height?.message} />
                                            <Input label="Ancho (cm)" type="number" {...register("customPackage.width", { valueAsNumber: true })} error={errors.customPackage?.width?.message} />
                                            <Input label="Largo (cm)" type="number" {...register("customPackage.length", { valueAsNumber: true })} error={errors.customPackage?.length?.message} />
                                        </>
                                    )}

                                    <div className="col-span-2">
                                        <Input label="Descripción" {...register("description")} error={errors.description?.message} />
                                    </div>
                                    <Input label="Valor Declarado" type="number" {...register("contentValue", { valueAsNumber: true })} error={errors.contentValue?.message} />

                                    <div>
                                        <label className="block text-sm font-medium text-gray-700 mb-1">Método de Pago COD</label>
                                        <select {...register("codPaymentMethod")} className="w-full border-gray-300 rounded-md shadow-sm p-2 border">
                                            <option value="cash">Efectivo (cash)</option>
                                            <option value="data_phone">Datáfono</option>
                                        </select>
                                        {errors.codPaymentMethod && <p className="text-red-500 text-xs mt-1">{errors.codPaymentMethod.message}</p>}
                                    </div>
                                </div>
                            </section>

                            {/* Rates Section */}
                            {rates.length > 0 && (
                                <section className="bg-gray-50 p-4 rounded-md">
                                    {loadingAI && (
                                        <div className="text-center py-4 text-blue-600 animate-pulse text-sm font-medium">
                                            🤖 IA Analizando mejores opciones...
                                        </div>
                                    )}

                                    {aiAnalysis && !loadingAI && (
                                        <div className="bg-blue-50 border border-blue-200 p-4 rounded-xl mb-4">
                                            <h4 className="text-blue-900 font-bold flex items-center gap-2 mb-1">
                                                <span>🤖</span> Recomendación: {aiAnalysis.recommended_carrier}
                                            </h4>
                                            <p className="text-blue-800 text-xs leading-relaxed italic">
                                                "{aiAnalysis.reasoning}"
                                            </p>
                                        </div>
                                    )}

                                    <div className="space-y-2 max-h-60 overflow-y-auto">
                                        {rates.map((rate) => {
                                            const isRecommended = rate.carrier === aiAnalysis?.recommended_carrier;
                                            return (
                                                <div
                                                    key={rate.idRate}
                                                    className={`p-3 border rounded cursor-pointer flex justify-between items-center transition-all ${selectedRate?.idRate === rate.idRate ? "border-indigo-500 bg-indigo-50 ring-1 ring-indigo-200" : isRecommended ? "border-green-300 bg-green-50" : "border-gray-200 bg-white hover:bg-gray-50"}`}
                                                    onClick={() => setSelectedRate(rate)}
                                                >
                                                    <div>
                                                        <div className="flex items-center gap-2">
                                                            <p className="font-bold text-gray-800">{rate.carrier}</p>
                                                            {isRecommended && <span className="text-[10px] bg-green-500 text-white px-2 py-0.5 rounded-full">Recomendado</span>}
                                                        </div>
                                                        <p className="text-sm text-gray-600">{rate.product} - {rate.deliveryDays} días</p>
                                                    </div>
                                                    <div className="text-right">
                                                        <p className="font-bold text-lg text-indigo-600">${rate.flete.toLocaleString()}</p>
                                                    </div>
                                                </div>
                                            );
                                        })}
                                    </div>
                                </section>
                            )}

                            <div className="flex justify-end pt-4 gap-4 sticky bottom-0 bg-white p-4 border-t border-gray-100">
                                <Button type="button" variant="secondary" onClick={onClose} disabled={loading}>
                                    Cancelar
                                </Button>
                                <Button
                                    type="button"
                                    variant="secondary"
                                    onClick={handleSubmit(handleQuote)}
                                    disabled={loading}
                                >
                                    {rates.length > 0 ? "Cotizar de Nuevo" : "Cotizar"}
                                </Button>

                                {rates.length > 0 && (
                                    <Button type="button" onClick={handleSubmit(handleGenerate)} disabled={loading || !selectedRate}>
                                        {loading ? "Generando..." : "Generar Guía"}
                                    </Button>
                                )}
                            </div>
                        </form>
                    )}
                </div>
            </div>
        </div>
    );
}
