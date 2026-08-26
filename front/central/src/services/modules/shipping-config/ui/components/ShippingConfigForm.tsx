'use client';

import { useState, useEffect, useCallback } from 'react';
import { Button, Input, Select } from '@/shared/ui';
import { useToast } from '@/shared/providers/toast-provider';
import {
    TruckIcon,
    CubeIcon,
    BuildingStorefrontIcon,
    InformationCircleIcon,
    PlusIcon,
    TrashIcon,
    StarIcon,
    LinkIcon,
} from '@heroicons/react/24/outline';
import {
    getShippingConfigAction,
    saveShippingConfigAction,
    setDefaultWarehouseAction,
} from '../../infra/actions';
import {
    ShippingConfigOverview,
    ShippingBox,
    CarrierSetting,
} from '../../domain/types';

interface ShippingConfigFormProps {
    businessId?: number;
    onClose?: () => void;
}

const emptyBox = (): ShippingBox => ({
    name: '',
    weight: 1,
    length: 20,
    width: 20,
    height: 20,
    max_items: 5,
});


interface ToggleProps {
    checked: boolean;
    onChange: (value: boolean) => void;
    disabled?: boolean;
    label: string;
    icon?: React.ReactNode;
    size?: 'sm' | 'md';
}

function Toggle({ checked, onChange, disabled, label, icon, size = 'sm' }: ToggleProps) {
    const track = size === 'md' ? 'h-6 w-11' : 'h-5 w-9';
    const knob = size === 'md' ? 'h-5 w-5' : 'h-3.5 w-3.5';
    const shift = size === 'md' ? 'translate-x-5' : 'translate-x-4';

    return (
        <button
            type="button"
            role="switch"
            aria-checked={checked}
            aria-label={label}
            disabled={disabled}
            onClick={() => onChange(!checked)}
            className={`inline-flex items-center gap-2 ${disabled ? 'opacity-40 cursor-not-allowed' : 'cursor-pointer'}`}
        >
            <span
                className={`relative inline-flex ${track} flex-shrink-0 items-center rounded-full transition-colors ${checked ? '' : 'bg-gray-300 dark:bg-gray-600'}`}
                style={checked ? { backgroundColor: 'var(--color-primary)' } : undefined}
            >
                <span
                    className={`inline-block ${knob} rounded-full bg-white shadow-sm transition-transform dark:bg-gray-100 ${checked ? shift : 'translate-x-0.5'}`}
                />
            </span>
            <span className="inline-flex items-center gap-1 text-xs text-gray-600 dark:text-gray-300">
                {icon}
                {label}
            </span>
        </button>
    );
}

export function ShippingConfigForm({ businessId, onClose }: ShippingConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [overview, setOverview] = useState<ShippingConfigOverview | null>(null);
    const [strategy, setStrategy] = useState<'product_dimensions' | 'standard_box'>('product_dimensions');
    const [boxes, setBoxes] = useState<ShippingBox[]>([]);
    const [carriers, setCarriers] = useState<CarrierSetting[]>([]);
    const [alwaysInsure, setAlwaysInsure] = useState(false);

    const load = useCallback(async () => {
        setLoading(true);
        const res = await getShippingConfigAction(businessId);
        if (res.success && res.data) {
            const data = res.data as ShippingConfigOverview;
            setOverview(data);
            setStrategy(data.business.package_strategy);
            setBoxes(data.business.boxes || []);
            setCarriers(data.business.carriers || []);
            setAlwaysInsure(data.business.always_insure || false);
        } else {
            showToast((res as { message?: string }).message || 'No se pudo cargar la configuracion', 'error');
        }
        setLoading(false);
    }, [businessId, showToast]);

    useEffect(() => {
        load();
    }, [load]);

    const updateCarrier = (code: string, patch: Partial<CarrierSetting>) => {
        setCarriers((prev) => prev.map((c) => (c.code === code ? { ...c, ...patch } : c)));
    };

    const updateDirect = (code: string, enabled: boolean) => {
        setCarriers((prev) =>
            prev.map((c) =>
                c.code === code ? { ...c, direct: { ...c.direct, enabled } } : c
            )
        );
    };

    const updateBox = (index: number, patch: Partial<ShippingBox>) => {
        setBoxes((prev) => prev.map((b, i) => (i === index ? { ...b, ...patch } : b)));
    };

    const handleSave = async () => {
        setSaving(true);
        const res = await saveShippingConfigAction({ package_strategy: strategy, boxes, carriers, always_insure: alwaysInsure }, businessId);
        if (res.success) {
            showToast('Configuracion de envios guardada', 'success');
            await load();
        } else {
            showToast((res as { message?: string }).message || 'Error al guardar', 'error');
        }
        setSaving(false);
    };

    const handleSetDefault = async (warehouseId: number) => {
        const res = await setDefaultWarehouseAction(warehouseId, businessId);
        if (res.success) {
            showToast('Bodega predeterminada actualizada', 'success');
            await load();
        } else {
            showToast((res as { message?: string }).message || 'Error al cambiar la bodega', 'error');
        }
    };

    const carrierName = (code: string) =>
        overview?.carriers.find((c) => c.code === code)?.name || code;

    if (loading) {
        return (
            <div className="flex items-center justify-center py-16">
                <div className="w-10 h-10 rounded-full animate-spin" style={{ borderWidth: '3px', borderStyle: 'solid', borderColor: 'rgba(0,0,0,0.1)', borderTopColor: 'var(--color-primary)' }} />
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div className="bg-gray-50 dark:bg-gray-700 rounded-xl p-6 space-y-4">
                <div className="flex items-center gap-2 mb-2">
                    <TruckIcon className="w-5 h-5 text-gray-700 dark:text-gray-200" />
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Transportadoras</h3>
                </div>
                <p className="text-xs text-gray-500 dark:text-gray-400 flex items-start gap-1">
                    <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                    <span>Elige que transportadoras se ofrecen en tus cotizaciones y guias. Si tienes convenio directo con alguna, activa su integracion propia.</span>
                </p>

                <div className="bg-white dark:bg-gray-800 rounded-lg border dark:border-gray-600 p-3 flex items-center justify-between gap-3 flex-wrap">
                    <div>
                        <span className="text-sm font-medium text-gray-800 dark:text-gray-100">Cotizar con seguro adicional en el checkout de WooCommerce</span>
                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
                            El seguro minimo obligatorio siempre va incluido en el precio, sin importar esta opcion. Activarla suma el seguro adicional (cobertura sobre el valor declarado) a todas las cotizaciones que ve el cliente en el checkout.
                        </p>
                    </div>
                    <Toggle
                        size="md"
                        checked={alwaysInsure}
                        onChange={setAlwaysInsure}
                        label="Seguro adicional en checkout"
                    />
                </div>

                <div className="space-y-2">
                    {carriers.map((c) => (
                        <div key={c.code} className="bg-white dark:bg-gray-800 rounded-lg border dark:border-gray-600 p-3">
                            <div className="flex items-center justify-between gap-3 flex-wrap">
                                <div className="flex items-center gap-2">
                                    <Toggle
                                        size="md"
                                        checked={c.enabled}
                                        onChange={(value) => updateCarrier(c.code, { enabled: value })}
                                        label=""
                                    />
                                    <span className="text-sm font-medium text-gray-800 dark:text-gray-100">{carrierName(c.code)}</span>
                                </div>

                                <div className="flex items-center gap-5">
                                    <Toggle
                                        checked={c.allow_prepaid}
                                        disabled={!c.enabled}
                                        onChange={(value) => updateCarrier(c.code, { allow_prepaid: value })}
                                        label="Prepago"
                                    />
                                    <Toggle
                                        checked={c.allow_cod}
                                        disabled={!c.enabled}
                                        onChange={(value) => updateCarrier(c.code, { allow_cod: value })}
                                        label="Contra entrega"
                                    />
                                    <Toggle
                                        checked={c.direct.enabled}
                                        disabled={!c.enabled}
                                        onChange={(value) => updateDirect(c.code, value)}
                                        label="Integracion propia"
                                        icon={<LinkIcon className="w-3.5 h-3.5" />}
                                    />
                                </div>
                            </div>

                            {c.direct.enabled && c.direct.status === 'pending' && (
                                <p className="mt-2 text-xs text-amber-700 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 rounded px-2 py-1.5">
                                    Tu convenio directo con {carrierName(c.code)} quedo registrado. La conexion estara disponible proximamente; mientras tanto seguimos despachando con la cuenta de Probability.
                                </p>
                            )}
                        </div>
                    ))}
                </div>
            </div>

            <div className="bg-gray-50 dark:bg-gray-700 rounded-xl p-6 space-y-4">
                <div className="flex items-center gap-2 mb-2">
                    <CubeIcon className="w-5 h-5 text-gray-700 dark:text-gray-200" />
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Empaque</h3>
                </div>
                <p className="text-xs text-gray-500 dark:text-gray-400 flex items-start gap-1">
                    <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                    <span>
                        Con estos datos declaramos el paquete a la transportadora al cotizar y al generar la guia.
                        La tarifa se calcula con el peso volumetrico, asi que declarar de menos termina en recobros.
                    </span>
                </p>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">Como calculamos el paquete</label>
                    <Select
                        value={strategy}
                        onChange={(e) => setStrategy(e.target.value as 'product_dimensions' | 'standard_box')}
                        options={[
                            { value: 'product_dimensions', label: 'Dimensiones del producto' },
                            { value: 'standard_box', label: 'Cajas estandar' },
                        ]}
                        className="bg-white dark:bg-gray-800"
                    />
                    <div className="mt-2 rounded-lg bg-white dark:bg-gray-800 border dark:border-gray-600 p-3 text-xs text-gray-600 dark:text-gray-300 space-y-1.5">
                        {strategy === 'standard_box' ? (
                            <>
                                <p><span className="font-semibold">Cajas estandar:</span> empacas siempre en las mismas cajas. Elegimos la caja mas pequena en la que quepa el pedido.</p>
                                <p>1. Descartamos las cajas cuyo tope de items sea menor a la cantidad del pedido.</p>
                                <p>2. Descartamos las cajas donde el producto mas grande no quepa (comparamos los tres lados, rotando la caja).</p>
                                <p>3. De las que quedan usamos la mas ajustada.</p>
                                <p>El peso declarado es el real del pedido; si es menor al de la caja, usamos el de la caja.</p>
                            </>
                        ) : (
                            <>
                                <p><span className="font-semibold">Dimensiones del producto:</span> tomamos el peso y las medidas del catalogo.</p>
                                <p>El peso se suma por cantidad; de las medidas tomamos el lado mas grande entre los productos del pedido.</p>
                                <p>Si un producto no tiene medidas cargadas, ese lado se declara en 10 cm y el peso en 1 kg. Revisa tu catalogo para evitarlo.</p>
                            </>
                        )}
                    </div>
                </div>

                {strategy === 'standard_box' && (
                    <div className="space-y-3">
                        <p className="text-xs text-gray-500 dark:text-gray-400 flex items-start gap-1">
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>
                                Carga una caja por cada empaque que uses. <span className="font-medium">Max items</span> es cuantos productos caben;
                                el <span className="font-medium">peso</span> es el del empaque vacio con relleno.
                                Ejemplo: Chica 20x20x20, 1 kg, hasta 5 items; Mediana 30x40x30, 3 kg, hasta 10.
                            </span>
                        </p>
                        {boxes.map((box, index) => (
                            <div key={index} className="bg-white dark:bg-gray-800 rounded-lg border dark:border-gray-600 p-3">
                                <div className="grid grid-cols-2 md:grid-cols-6 gap-2 items-end">
                                    <div className="col-span-2">
                                        <label className="block text-xs text-gray-600 dark:text-gray-300 mb-1">Nombre</label>
                                        <Input
                                            type="text"
                                            value={box.name}
                                            onChange={(e) => updateBox(index, { name: e.target.value })}
                                            placeholder="Caja Chica"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-xs text-gray-600 dark:text-gray-300 mb-1">Largo (cm)</label>
                                        <Input type="number" value={box.length ?? ''} onChange={(e) => updateBox(index, { length: Number(e.target.value) })} />
                                    </div>
                                    <div>
                                        <label className="block text-xs text-gray-600 dark:text-gray-300 mb-1">Ancho (cm)</label>
                                        <Input type="number" value={box.width ?? ''} onChange={(e) => updateBox(index, { width: Number(e.target.value) })} />
                                    </div>
                                    <div>
                                        <label className="block text-xs text-gray-600 dark:text-gray-300 mb-1">Alto (cm)</label>
                                        <Input type="number" value={box.height ?? ''} onChange={(e) => updateBox(index, { height: Number(e.target.value) })} />
                                    </div>
                                    <div className="flex gap-2 items-end">
                                        <div className="flex-1">
                                            <label className="block text-xs text-gray-600 dark:text-gray-300 mb-1">Max items</label>
                                            <Input type="number" value={box.max_items} onChange={(e) => updateBox(index, { max_items: Number(e.target.value) })} />
                                        </div>
                                        <button
                                            type="button"
                                            onClick={() => setBoxes((prev) => prev.filter((_, i) => i !== index))}
                                            className="p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded"
                                            aria-label="Eliminar caja"
                                        >
                                            <TrashIcon className="w-4 h-4" />
                                        </button>
                                    </div>
                                </div>
                                <div className="mt-2 w-32">
                                    <label className="block text-xs text-gray-600 dark:text-gray-300 mb-1">Peso (kg)</label>
                                    <Input type="number" value={box.weight ?? ''} onChange={(e) => updateBox(index, { weight: Number(e.target.value) })} />
                                </div>
                            </div>
                        ))}

                        <Button type="button" variant="secondary" onClick={() => setBoxes((prev) => [...prev, emptyBox()])}>
                            <PlusIcon className="w-4 h-4 mr-1" />
                            Agregar caja
                        </Button>
                    </div>
                )}
            </div>

            <div className="bg-gray-50 dark:bg-gray-700 rounded-xl p-6 space-y-3">
                <div className="flex items-center gap-2 mb-2">
                    <BuildingStorefrontIcon className="w-5 h-5 text-gray-700 dark:text-gray-200" />
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Direcciones de origen</h3>
                </div>
                <p className="text-xs text-gray-500 dark:text-gray-400 flex items-start gap-1">
                    <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                    <span>Se configuran en el modulo de bodegas. Aqui solo eliges desde cual despachas por defecto.</span>
                </p>

                <div className="space-y-2">
                    {(overview?.warehouses || []).map((w) => (
                        <div key={w.id} className="flex items-center justify-between bg-white dark:bg-gray-800 rounded-lg border dark:border-gray-600 p-3">
                            <div>
                                <p className="text-sm font-medium text-gray-800 dark:text-gray-100">
                                    {w.name}
                                    {w.is_default && (
                                        <span
                                            className="ml-2 text-xs px-2 py-0.5 rounded-full text-white"
                                            style={{ backgroundColor: 'var(--color-primary)' }}
                                        >
                                            Predeterminada
                                        </span>
                                    )}
                                </p>
                                <p className="text-xs text-gray-500 dark:text-gray-400">
                                    {[w.address, w.city, w.state].filter(Boolean).join(', ') || 'Sin direccion registrada'}
                                </p>
                            </div>
                            {!w.is_default && (
                                <button
                                    type="button"
                                    onClick={() => handleSetDefault(w.id)}
                                    className="text-xs flex items-center gap-1 px-2 py-1 rounded border transition-colors hover:opacity-80"
                                    style={{ borderColor: 'var(--color-primary)', color: 'var(--color-primary)' }}
                                >
                                    <StarIcon className="w-3.5 h-3.5" />
                                    Usar por defecto
                                </button>
                            )}
                        </div>
                    ))}
                    {(overview?.warehouses || []).length === 0 && (
                        <p className="text-sm text-gray-500 dark:text-gray-400">Este negocio aun no tiene bodegas registradas.</p>
                    )}
                </div>
            </div>

            <div className="flex justify-end gap-2">
                {onClose && (
                    <Button type="button" variant="secondary" onClick={onClose} disabled={saving}>
                        Cerrar
                    </Button>
                )}
                <Button type="button" onClick={handleSave} disabled={saving}>
                    {saving ? 'Guardando...' : 'Guardar configuracion'}
                </Button>
            </div>
        </div>
    );
}
