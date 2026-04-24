'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import {
    XMarkIcon,
    ChevronLeftIcon,
    ChevronRightIcon,
    ArrowRightIcon,
} from '@heroicons/react/24/outline';

interface TourStep {
    id: string;
    icon: string;
    title: string;
    subtitle: string;
    whatIs: string;
    whenToUse: string;
    example: string;
    href?: string;
    highlight?: string;
}

const STEPS: TourStep[] = [
    {
        id: 'welcome',
        icon: '📦',
        title: 'Bienvenido al módulo de Inventario',
        subtitle: 'Tu sistema WMS completo',
        whatIs: 'Este módulo maneja todo lo relacionado con tu inventario físico: dónde está cada producto, cuánto tienes, quién lo movió, cuándo vence, y mucho más.',
        whenToUse: 'Siempre que necesites saber o actualizar el estado del inventario de uno o varios productos.',
        example: 'Desde "tengo 50 camisetas" hasta "la camiseta SKU-123 está en el pasillo A-01, rack 2, nivel 3, en el lote LOT-005 que vence en 60 días".',
        highlight: 'Recorramos los 11 sub-módulos para entender para qué sirve cada uno.',
    },
    {
        id: 'warehouses',
        icon: '🏭',
        title: 'Bodegas',
        subtitle: 'Dónde se guarda tu inventario',
        whatIs: 'Define la jerarquía física de cada bodega: Zonas → Pasillos → Racks → Niveles → Posiciones. Cada posición puede tener dimensiones (alto/ancho/largo), peso máximo y flags (picking, bulk, cross-dock, averiado, etc.).',
        whenToUse: 'Cuando creas una bodega nueva, abres una sección, o reorganizas el layout.',
        example: 'Bodega "Principal" → Zona "Picking" → Pasillo "A" → Rack "R-01" → Nivel "L-02" → Posición "P-03".',
        href: '/warehouses',
    },
    {
        id: 'stock',
        icon: '📊',
        title: 'Stock',
        subtitle: 'Inventario actual por bodega',
        whatIs: 'Vista consolidada de cuánto tienes de cada producto en cada bodega, con disponible / reservado / mínimo / máximo. Desde aquí puedes ajustar stock o transferir entre bodegas.',
        whenToUse: 'Consulta diaria de lo que hay en cada bodega o ajustes manuales (robo, daño, corrección).',
        example: 'Producto SKU-123 en bodega Principal: 50 totales, 10 reservados, 40 disponibles, mínimo 20 (está OK).',
        href: '/inventory',
    },
    {
        id: 'movements',
        icon: '🔄',
        title: 'Movimientos',
        subtitle: 'Historial crudo de entradas y salidas',
        whatIs: 'Log completo de cada movimiento de inventario: entradas, salidas, transferencias, ajustes, reservas, devoluciones. Cada uno con su tipo, producto, bodega, cantidad, usuario y fecha.',
        whenToUse: 'Para rastrear "qué pasó con este producto" o auditar quién movió qué.',
        example: 'Filtrar por producto SKU-123 y ver 12 entradas de proveedor y 47 salidas por ventas en los últimos 30 días.',
        href: '/inventory/movements',
    },
    {
        id: 'traceability',
        icon: '🏷️',
        title: 'Trazabilidad (Lotes + Series + UoM)',
        subtitle: 'Identificación granular',
        whatIs: '3 capas unificadas: (1) Lotes para grupos con fecha de fabricación/vencimiento (FEFO), (2) Series para unidades individuales con número serial único, (3) UoM para convertir entre unidades (caja/pallet/docena/etc).',
        whenToUse: 'Lotes → alimentos, medicamentos, químicos. Series → electrónicos, herramientas. UoM → productos que se compran en una unidad y se venden en otra.',
        example: 'Lote LOT-005 vence en 30 días → prioridad FEFO. Serie SN-42 de la TV se vendió el 15/mar. 1 caja = 12 unidades.',
        href: '/inventory/traceability',
    },
    {
        id: 'kardex',
        icon: '📑',
        title: 'Kardex',
        subtitle: 'Reporte contable por producto/bodega',
        whatIs: 'Report con saldo acumulado (running balance) para un producto en una bodega entre dos fechas. Muestra todos los movimientos, totales de entradas/salidas y saldo final. Exportable a CSV.',
        whenToUse: 'Cierre contable mensual, conciliación con compras/ventas, informes a la DIAN.',
        example: 'Kardex del producto SKU-123 en Principal durante marzo: 150 entradas, 80 salidas, saldo final 70.',
        href: '/inventory/kardex',
    },
    {
        id: 'operations',
        icon: '📥',
        title: 'Operaciones (Put-away + Reposición + Cross-dock)',
        subtitle: 'Flujos operativos WMS',
        whatIs: '3 operaciones unificadas: (1) Put-away: dónde poner la mercancía nueva según reglas y prioridades. (2) Reposición: mover stock de bulk a picking cuando baja del mínimo. (3) Cross-dock: despachar sin almacenar cuando llega ya vendida.',
        whenToUse: 'Put-away al recibir · Reposición al detectar bajo stock en picking · Cross-dock cuando inbound coincide con outbound pendiente.',
        example: '100 cajas llegan → sistema sugiere ponerlas en Zona A (put-away) · Picking baja a 5 uds → se crea tarea de reposición desde Bulk (100 uds) · Orden de cliente llega justo cuando entran sus productos → cross-dock sin almacenar.',
        href: '/inventory/operations',
    },
    {
        id: 'slotting',
        icon: '📈',
        title: 'Slotting ABC',
        subtitle: 'Clasificación por rotación',
        whatIs: 'Analiza movimientos de los últimos N días y clasifica productos por velocidad: A (80% del volumen, alta rotación), B (15%), C (5%, baja rotación). Ayuda a decidir dónde ubicar cada producto.',
        whenToUse: 'Rediseño del layout · optimizar tiempos de picking · identificar obsoletos.',
        example: 'Ejecutar slotting 90 días → 20 SKUs Clase A (ponerlos en Zona Picking cerca del despacho), 50 Clase B (zona media), 230 Clase C (bulk profundo).',
        href: '/inventory/analytics/slotting',
    },
    {
        id: 'audit',
        icon: '✅',
        title: 'Auditoría (Planes + Tareas + Discrepancias)',
        subtitle: 'Conteos cíclicos',
        whatIs: '3 pasos unificados: (1) Plan define la estrategia (ABC, zona, random, total) y frecuencia. (2) Tarea genera las líneas a contar cuando llega el momento. (3) Si lo contado difiere del sistema → Discrepancia que se aprueba (ajusta stock) o rechaza.',
        whenToUse: 'En lugar de parar la bodega un día al año, cuentas un subset cada semana/mes. Los ajustes se aplican automáticamente al aprobar.',
        example: 'Plan ABC semanal → tarea cuenta los 20 SKUs tipo A → contador marca 48 en lugar de 50 esperado → discrepancia abierta → supervisor aprueba con nota "merma picking" → stock se ajusta a 48 y queda el movimiento tipo count_adjustment.',
        href: '/inventory/audit',
    },
    {
        id: 'lpn',
        icon: '📦',
        title: 'LPN (License Plate Numbers)',
        subtitle: 'Contenedores lógicos',
        whatIs: 'Pallet / caja / tote con un código único que agrupa productos para moverlos juntos. Un LPN puede tener múltiples productos/lotes/series adentro.',
        whenToUse: 'Al recibir un pallet completo del proveedor, al preparar un contenedor para despacho, al transferir entre bodegas.',
        example: 'Pallet LPN-001 contiene 5 cajas de SKU-A (lote L-001) + 3 cajas de SKU-B (lote L-002). Al mover el LPN, todos los productos se mueven con él.',
        href: '/inventory/lpn',
    },
    {
        id: 'scan',
        icon: '📱',
        title: 'Scan',
        subtitle: 'Resolución universal de códigos',
        whatIs: 'Ingresa cualquier código (barra, QR, LPN, número de serie, lote, código de ubicación) y el sistema resuelve qué es con búsqueda 6-way: LPN → Ubicación → Serie → Lote → UoM barcode → Producto.',
        whenToUse: 'Terminal móvil del operario de bodega · validación rápida · búsqueda sin saber el tipo.',
        example: 'Escanear "LOT-005" → sistema detecta que es un lote y muestra producto + ubicación + cantidad disponible.',
        href: '/inventory/mobile',
    },
    {
        id: 'sync',
        icon: '🔄',
        title: 'Sync Logs',
        subtitle: 'Bitácora de integraciones',
        whatIs: 'Registro de cada sincronización inbound/outbound con integraciones (Shopify, Amazon, Meli, WhatsApp). Usa hash SHA-256 para idempotencia (evitar duplicados).',
        whenToUse: 'Debugging cuando Shopify reporta un stock que no cuadra con el nuestro · verificar qué deltas se enviaron.',
        example: 'Log #42: inbound de Shopify, hash abc123, status=success, synced_at=hace 5min.',
        href: '/inventory/sync/logs',
    },
    {
        id: 'summary',
        icon: '🎓',
        title: 'Resumen',
        subtitle: '¿Por dónde empezar?',
        whatIs: 'Flujo típico: 1) Configura Bodegas con jerarquía. 2) Carga Stock inicial. 3) Define reglas de Put-away. 4) Activa Trazabilidad si tus productos lo requieren. 5) Programa conteos en Auditoría. 6) Revisa Slotting cada trimestre.',
        whenToUse: 'Vuelve a este tour cuando quieras refrescar qué hace cada parte.',
        example: 'Puedes abrirlo otra vez desde el botón "📖 Guía" en la barra superior.',
        highlight: 'Todos los sub-módulos están conectados: un escaneo crea un evento, un conteo crea una discrepancia que al aprobarse crea un movimiento que alimenta el kardex, etc.',
    },
];

const STORAGE_KEY = 'inventory_tour_seen_v1';

interface Props {
    isOpen: boolean;
    onClose: () => void;
    initialStep?: number;
}

export default function InventoryTour({ isOpen, onClose, initialStep = 0 }: Props) {
    const [currentStep, setCurrentStep] = useState(initialStep);

    useEffect(() => {
        if (isOpen) setCurrentStep(initialStep);
    }, [isOpen, initialStep]);

    const handleClose = () => {
        try { localStorage.setItem(STORAGE_KEY, 'true'); } catch {}
        onClose();
    };

    const next = () => {
        if (currentStep < STEPS.length - 1) setCurrentStep(currentStep + 1);
        else handleClose();
    };

    const prev = () => {
        if (currentStep > 0) setCurrentStep(currentStep - 1);
    };

    useEffect(() => {
        const handleEsc = (e: KeyboardEvent) => {
            if (e.key === 'Escape' && isOpen) handleClose();
            if (e.key === 'ArrowRight' && isOpen) next();
            if (e.key === 'ArrowLeft' && isOpen) prev();
        };
        window.addEventListener('keydown', handleEsc);
        return () => window.removeEventListener('keydown', handleEsc);
    }, [isOpen, currentStep]);

    if (!isOpen) return null;

    const step = STEPS[currentStep];
    const progress = ((currentStep + 1) / STEPS.length) * 100;

    return (
        <div className="fixed inset-0 z-[60] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-3xl max-h-[90vh] overflow-hidden flex flex-col">
                <div className="h-1 bg-gray-100 dark:bg-gray-700">
                    <div className="h-full bg-gradient-to-r from-[#7c3aed] to-[#a855f7] transition-all duration-300" style={{ width: `${progress}%` }} />
                </div>

                <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                    <div className="flex items-center gap-3">
                        <span className="text-3xl">{step.icon}</span>
                        <div>
                            <h2 className="text-lg font-bold text-gray-900 dark:text-white">{step.title}</h2>
                            <p className="text-xs text-gray-500 dark:text-gray-400">{step.subtitle}</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-2">
                        <span className="text-xs text-gray-400">{currentStep + 1} / {STEPS.length}</span>
                        <button onClick={handleClose} className="p-1.5 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-500 dark:text-gray-400">
                            <XMarkIcon className="w-5 h-5" />
                        </button>
                    </div>
                </div>

                <div className="flex-1 overflow-y-auto p-6 space-y-5">
                    <Section title="¿Qué es?" body={step.whatIs} accent="purple" />
                    <Section title="¿Cuándo usarlo?" body={step.whenToUse} accent="blue" />
                    <Section title="Ejemplo real" body={step.example} accent="emerald" />

                    {step.highlight && (
                        <div className="p-4 rounded-xl bg-gradient-to-br from-purple-50 to-indigo-50 dark:from-purple-900/20 dark:to-indigo-900/20 border border-purple-200 dark:border-purple-800">
                            <p className="text-sm text-purple-900 dark:text-purple-100 font-medium">💡 {step.highlight}</p>
                        </div>
                    )}

                    {step.href && (
                        <Link
                            href={step.href}
                            onClick={handleClose}
                            className="inline-flex items-center gap-2 px-5 py-3 btn-business-primary text-white font-semibold rounded-lg shadow-lg hover:shadow-xl transition-all"
                        >
                            Ir a {step.title} <ArrowRightIcon className="w-4 h-4" />
                        </Link>
                    )}
                </div>

                <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between bg-gray-50 dark:bg-gray-900/40">
                    <button
                        onClick={prev}
                        disabled={currentStep === 0}
                        className="inline-flex items-center gap-1 px-3 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-md disabled:opacity-40 disabled:cursor-not-allowed"
                    >
                        <ChevronLeftIcon className="w-4 h-4" /> Atrás
                    </button>

                    <div className="flex items-center gap-1">
                        {STEPS.map((s, i) => (
                            <button
                                key={s.id}
                                onClick={() => setCurrentStep(i)}
                                className={`h-2 rounded-full transition-all ${i === currentStep ? 'w-6 bg-purple-600' : i < currentStep ? 'w-2 bg-purple-400' : 'w-2 bg-gray-300 dark:bg-gray-600'}`}
                                title={s.title}
                            />
                        ))}
                    </div>

                    <button
                        onClick={next}
                        className="inline-flex items-center gap-1 px-4 py-2 text-sm btn-business-primary text-white font-semibold rounded-md"
                    >
                        {currentStep === STEPS.length - 1 ? 'Finalizar' : 'Siguiente'}
                        <ChevronRightIcon className="w-4 h-4" />
                    </button>
                </div>
            </div>
        </div>
    );
}

function Section({ title, body, accent }: { title: string; body: string; accent: 'purple' | 'blue' | 'emerald' }) {
    const accentStyles: Record<string, string> = {
        purple: 'border-l-purple-500 bg-purple-50/40 dark:bg-purple-900/10',
        blue: 'border-l-blue-500 bg-blue-50/40 dark:bg-blue-900/10',
        emerald: 'border-l-emerald-500 bg-emerald-50/40 dark:bg-emerald-900/10',
    };
    const titleStyles: Record<string, string> = {
        purple: 'text-purple-700 dark:text-purple-300',
        blue: 'text-blue-700 dark:text-blue-300',
        emerald: 'text-emerald-700 dark:text-emerald-300',
    };
    return (
        <div className={`border-l-4 ${accentStyles[accent]} px-4 py-3 rounded-r-md`}>
            <h4 className={`text-xs uppercase tracking-wider font-bold mb-1 ${titleStyles[accent]}`}>{title}</h4>
            <p className="text-sm text-gray-700 dark:text-gray-200 leading-relaxed">{body}</p>
        </div>
    );
}

export function useInventoryTour() {
    const [isOpen, setIsOpen] = useState(false);

    const open = () => setIsOpen(true);
    const close = () => setIsOpen(false);
    const hasSeen = () => {
        try { return localStorage.getItem(STORAGE_KEY) === 'true'; } catch { return false; }
    };

    return { isOpen, open, close, hasSeen };
}
