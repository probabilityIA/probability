"use client";

import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Input, Button } from "@/shared/ui";
import { ShipmentApiRepository } from "@/services/modules/shipments/infra/repository/api-repository";
import { EnvioClickQuoteRequest, EnvioClickRate } from "@/services/modules/shipments/domain/types";
import { OrderApiRepository } from "@/services/modules/orders/infra/repository/api-repository";
import { Order } from "@/services/modules/orders/domain/types";

// Zod Schema
const addressSchema = z.object({
    company: z.string().optional(), // Made optional as per some flows, but validation might be needed on backend
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

export const ShippingForm = () => {
    const [loading, setLoading] = useState(false);
    const [success, setSuccess] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [orders, setOrders] = useState<Order[]>([]);
    const [rates, setRates] = useState<EnvioClickRate[]>([]);
    const [selectedRate, setSelectedRate] = useState<EnvioClickRate | null>(null);
    const [hasQuoted, setHasQuoted] = useState(false);

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
                // Default origin for convenience (as per user example)
                company: "Company", firstName: "Santiago", lastName: "Muñoz", email: "santiago@test.com", phone: "3227684041",
                address: "Calle 98 62-37", suburb: "Gaitan", crossStreet: "Calle 39c #10-69-NA", reference: "Casa puertas negras", daneCode: "11001000"
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

    useEffect(() => {
        const fetchOrders = async () => {
            const repo = new OrderApiRepository();
            try {
                // Fetching first page of orders
                const res = await repo.getOrders({ page_size: 50 });
                setOrders(res.data);
            } catch (e) {
                console.error("Failed to fetch orders", e);
            }
        };
        fetchOrders();
    }, []);

    const handleOrderSelect = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const orderId = e.target.value;
        const order = orders.find(o => o.id === orderId);
        if (order) {
            // Populate destination from order
            setValue("destination.company", order.customer_name); // Or business name if available
            setValue("destination.firstName", order.customer_name.split(" ")[0] || "");
            setValue("destination.lastName", order.customer_name.split(" ").slice(1).join(" ") || ".");
            setValue("destination.email", order.customer_email);
            setValue("destination.phone", order.customer_phone);
            setValue("destination.address", order.shipping_street);
            setValue("destination.suburb", order.shipping_state || ""); // Mapping state to suburb as fallback, simplistic
            // Note: EnvioClick requires specific suburb/daneCode often. User might need to edit.
            setValue("destination.crossStreet", "");
            setValue("destination.reference", "");
            setValue("destination.daneCode", "11001000"); // Default Bogota for now, or from order if available

            setValue("contentValue", order.total_amount);
            setValue("description", "Order " + order.order_number);
            setValue("external_order_id", order.order_number);
            setValue("myShipmentReference", "Orden " + order.internal_number);

            // Try to map dimensions if available
            if (order.weight && order.weight > 0) {
                setValue("packageSize", "custom");
                setValue("customPackage.weight", order.weight);
                setValue("customPackage.height", order.height || 10);
                setValue("customPackage.width", order.width || 10);
                setValue("customPackage.length", order.length || 10);
            }
        }
    };

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
            codValue: Number(data.contentValue), // Warning: codValue same as contentValue? inferred logic
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
            } else {
                setError("No se encontraron cotizaciones.");
            }
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : "Error consultando cotizaciones";
            setError(errorMessage);
        } finally {
            setLoading(false);
        }
    };

    const handleGenerate = async (data: FormValues) => {
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
            setSuccess(`Guía generada exitosamente! Tracking: ${res.data.trackingNumber}`);
            // Reset logic if needed
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : "Error generando guía";
            setError(errorMessage);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="p-6 bg-white rounded-lg shadow-md max-w-4xl mx-auto">
            <h2 className="text-2xl font-bold mb-6 text-gray-800">Generar Guía Envioclick</h2>

            <div className="mb-6">
                <label className="block text-sm font-medium text-gray-700 mb-2">Seleccionar Orden</label>
                <select
                    className="w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 p-2 border"
                    onChange={handleOrderSelect}
                    defaultValue=""
                >
                    <option value="" disabled>-- Seleccione una orden --</option>
                    {orders.map(o => (
                        <option key={o.id} value={o.id}>
                            {o.order_number} - {o.customer_name} ({o.total_amount} {o.currency})
                        </option>
                    ))}
                </select>
            </div>

            {error && (
                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
                    <p className="font-bold">Error</p>
                    <p>{error}</p>
                </div>
            )}

            {success && (
                <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded mb-4">
                    <p className="font-bold">Éxito</p>
                    <p>{success}</p>
                </div>
            )}

            <form onSubmit={handleSubmit(hasQuoted ? handleGenerate : handleQuote)} className="space-y-8">
                {/* Origin Section */}
                <section>
                    <h3 className="text-xl font-semibold mb-4 text-gray-700 border-b pb-2">Origen</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <Input label="Empresa" {...register("origin.company")} error={errors.origin?.company?.message} />
                        <Input label="Nombre" {...register("origin.firstName")} error={errors.origin?.firstName?.message} />
                        <Input label="Apellido" {...register("origin.lastName")} error={errors.origin?.lastName?.message} />
                        <Input label="Email" {...register("origin.email")} error={errors.origin?.email?.message} />
                        <Input label="Teléfono" {...register("origin.phone")} error={errors.origin?.phone?.message} />
                        <Input label="Dirección" {...register("origin.address")} error={errors.origin?.address?.message} />
                        <Input label="Barrio" {...register("origin.suburb")} error={errors.origin?.suburb?.message} />
                        <Input label="Cruzamiento" {...register("origin.crossStreet")} error={errors.origin?.crossStreet?.message} />
                        <Input label="Referencia" {...register("origin.reference")} error={errors.origin?.reference?.message} />
                        <Input label="Código DANE" {...register("origin.daneCode")} error={errors.origin?.daneCode?.message} />
                    </div>
                </section>

                {/* Destination Section */}
                <section>
                    <h3 className="text-xl font-semibold mb-4 text-gray-700 border-b pb-2">Destino</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <Input label="Empresa" {...register("destination.company")} error={errors.destination?.company?.message} />
                        <Input label="Nombre" {...register("destination.firstName")} error={errors.destination?.firstName?.message} />
                        <Input label="Apellido" {...register("destination.lastName")} error={errors.destination?.lastName?.message} />
                        <Input label="Email" {...register("destination.email")} error={errors.destination?.email?.message} />
                        <Input label="Teléfono" {...register("destination.phone")} error={errors.destination?.phone?.message} />
                        <Input label="Dirección" {...register("destination.address")} error={errors.destination?.address?.message} />
                        <Input label="Barrio" {...register("destination.suburb")} error={errors.destination?.suburb?.message} />
                        <Input label="Cruzamiento" {...register("destination.crossStreet")} error={errors.destination?.crossStreet?.message} />
                        <Input label="Referencia" {...register("destination.reference")} error={errors.destination?.reference?.message} />
                        <Input label="Código DANE" {...register("destination.daneCode")} error={errors.destination?.daneCode?.message} />
                    </div>
                </section>

                {/* Package Section */}
                <section>
                    <h3 className="text-xl font-semibold mb-4 text-gray-700 border-b pb-2">Paquete y Detalles</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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
                        <h3 className="text-xl font-semibold mb-4 text-gray-700">Cotizaciones Disponibles</h3>
                        <div className="space-y-2 max-h-60 overflow-y-auto">
                            {rates.map((rate) => (
                                <div
                                    key={rate.idRate}
                                    className={`p-3 border rounded cursor-pointer flex justify-between items-center ${selectedRate?.idRate === rate.idRate ? "border-indigo-500 bg-indigo-50" : "border-gray-200 bg-white"}`}
                                    onClick={() => setSelectedRate(rate)}
                                >
                                    <div>
                                        <p className="font-bold text-gray-800">{rate.carrier} - {rate.product}</p>
                                        <p className="text-sm text-gray-600">{rate.deliveryDays} días de entrega</p>
                                    </div>
                                    <div className="text-right">
                                        <p className="font-bold text-lg text-indigo-600">${rate.flete.toLocaleString()}</p>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </section>
                )}

                <div className="flex justify-end pt-4 gap-4">
                    {/* Show "Cotizar" button if not yet quoted or if user wants to quote again (maybe always visible? no, context dependent) 
                        For simplicity: 
                        - If no rates, show "Cotizar"
                        - If rates exist, show "Cotizar de nuevo" AND "Generar Guía" (if rate selected)
                     */}

                    <Button
                        type="button"
                        variant="secondary"
                        onClick={handleSubmit(handleQuote)}
                        disabled={loading}
                    >
                        {rates.length > 0 ? "Actualizar Cotización" : "Cotizar"}
                    </Button>

                    {rates.length > 0 && (
                        <Button type="button" onClick={handleSubmit(handleGenerate)} disabled={loading || !selectedRate}>
                            {loading ? "Generando..." : "Generar Guía"}
                        </Button>
                    )}
                </div>
            </form>
        </div>
    );
};
