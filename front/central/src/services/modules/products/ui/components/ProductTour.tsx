'use client';

import { useState, useEffect } from 'react';
import { XMarkIcon, ChevronLeftIcon, ChevronRightIcon } from '@heroicons/react/24/outline';

interface TourStep {
    id: string;
    icon: string;
    title: string;
    subtitle: string;
    whatIs: string;
    whenToUse: string;
    example: string;
    highlight?: string;
}

const STEPS: TourStep[] = [
    {
        id: 'welcome',
        icon: '📦',
        title: 'Productos y Catálogo',
        subtitle: 'Cómo funciona el sistema de productos',
        whatIs: 'El módulo de productos es el catálogo central de tu negocio. Cada fila es un SKU único: un producto con su código, precio, stock e integraciones. Las "Familias de variantes" son una capa opcional para agrupar SKUs relacionados (mismo producto en distintos colores, tallas, sabores).',
        whenToUse: 'Empieza aquí cuando quieras agregar un producto nuevo, editar su precio, conectarlo a una integración (Shopify, MercadoLibre) o revisar su inventario.',
        example: 'Tienda de suplementos:\n• SKU PT01001: Proteína Whey Vainilla 1kg\n• SKU PT01002: Proteína Whey Chocolate 1kg\n• Ambos pertenecen a la familia "Proteína Whey" con eje variante: sabor',
        highlight: 'Cada SKU es independiente: tiene su propio stock, precio e integraciones. La familia solo agrupa visualmente.',
    },
    {
        id: 'sku',
        icon: '🏷️',
        title: 'SKU (Stock Keeping Unit)',
        subtitle: 'La unidad atómica del catálogo',
        whatIs: 'Un SKU es un código único que identifica una variante específica de un producto. Si vendes una camisa en 3 tallas y 2 colores, tienes 6 SKUs distintos — aunque sean "la misma camisa".',
        whenToUse: 'Cada vez que el precio, el stock o el empaque es diferente, necesitas un SKU diferente. No uses el mismo SKU para cosas distintas.',
        example: 'Correcto:\n• CAM-ROJA-S (Camisa roja talla S)\n• CAM-ROJA-M (Camisa roja talla M)\n• CAM-AZUL-S (Camisa azul talla S)\n\nIncorrecto:\n• CAM-001 (para todas las camisas, sin distinguir)',
        highlight: 'El barcode (EAN/UPC) es opcional pero útil para scanning en bodega. Si no tienes código de barras, el SKU sirve como identificador de escaneo.',
    },
    {
        id: 'stock',
        icon: '📊',
        title: 'Stock e Inventario',
        subtitle: 'Gestión de cantidad disponible',
        whatIs: 'Cada SKU puede tener "Track inventory" activado o no. Si está activado, el sistema lleva cuenta del stock y puede bloquearte cuando se agota. Si no está activado, aparece "∞" y nunca se agota.',
        whenToUse: 'Activa el control de inventario para productos físicos con stock limitado. Desáctivalo para servicios, suscripciones o productos digitales.',
        example: 'Con control de inventario:\n• Proteína Whey: 255 unidades disponibles\n• Al procesar un pedido: stock baja automáticamente\n• Al llegar mercancía: stock sube con ajuste o recibo\n\nSin control:\n• Servicio de instalación: siempre "disponible", nunca se agota',
        highlight: 'El stock que ves aquí viene del módulo de Bodegas. Para ajustar stock físicamente, usa el módulo de Inventario → Ajustes.',
    },
    {
        id: 'families',
        icon: '🗂️',
        title: 'Familias de variantes',
        subtitle: 'Agrupa SKUs relacionados',
        whatIs: 'Una familia es un agrupador lógico de SKUs que son variantes del mismo producto base. No tiene stock ni precio propio — solo agrupa. Define los "ejes de variación" (color, talla, sabor) para que el sistema sepa cómo clasificar cada SKU dentro de ella.',
        whenToUse: 'Crea una familia cuando tengas 2 o más SKUs que son variantes de un mismo producto base y quieras verlos agrupados en reportes, catálogos o marketplaces.',
        example: 'Familia: "Proteina Whey Test"\n• Ejes: [{key: "sabor", label: "Sabor"}, {key: "talla", label: "Tamaño"}]\n• Variantes:\n  - PT01001: Vainilla 1kg (sabor=vainilla, talla=1kg)\n  - PT01002: Chocolate 1kg (sabor=chocolate, talla=1kg)\n  - PT01003: Fresa 2kg (sabor=fresa, talla=2kg)',
        highlight: 'Los SKUs existen de forma independiente aunque estén en una familia. Puedes tener SKUs sin familia y SKUs en familia en el mismo catálogo.',
    },
    {
        id: 'variant-axes',
        icon: '🔁',
        title: 'Ejes de variante',
        subtitle: 'Cómo definir variaciones',
        whatIs: 'Un eje de variante es una dimensión por la que los SKUs de una familia difieren. Se define en la familia como un JSON con "key" (identificador interno) y "label" (nombre visible). Cuando asignas un SKU a la familia, le das su "variant_attributes" con el valor en cada eje.',
        whenToUse: 'Define los ejes al crear la familia, antes de asignar SKUs. Usa claves simples en minuscula sin espacios (ej: "color", "talla", "sabor", "gramaje").',
        example: 'Familia Camisetas:\nEjes: [{"key":"color","label":"Color"},{"key":"talla","label":"Talla"}]\n\nSKU CAM-ROJA-S:\nVariant attributes: {"color":"Rojo","talla":"S"}\n→ Variant label automático: "Rojo - S"\n\nSKU CAM-AZUL-M:\nVariant attributes: {"color":"Azul","talla":"M"}\n→ Variant label: "Azul - M"',
        highlight: 'La "variant_signature" se genera automáticamente al canonicalizar los atributos. Dos SKUs en la misma familia no pueden tener exactamente los mismos atributos — eso sería un duplicado.',
    },
    {
        id: 'integrations',
        icon: '🔌',
        title: 'Integraciones por SKU',
        subtitle: 'Conectar con Shopify, MercadoLibre, etc.',
        whatIs: 'Cada SKU puede conectarse a múltiples plataformas externas. La integración guarda el ID externo (ej: el product_id de Shopify) para sincronizar precios, stock e información automáticamente cuando llegan pedidos.',
        whenToUse: 'Usa el botón "Integraciones" en cada fila para ver o editar las conexiones externas de ese SKU. Esto se configura una vez y luego es automático.',
        example: 'SKU PT01001 tiene 2 integraciones:\n• shopify: external_id="7823456789", tienda "mitienda.myshopify.com"\n• amazon: external_id="B0XXXXXX", ASIN marketplace\n\nCuando llega un pedido de Shopify con product_id=7823456789, el sistema lo mapea automáticamente al SKU PT01001.',
        highlight: 'Sin integración configurada, los pedidos de plataformas externas no pueden asociarse al SKU correcto. Los pedidos llegarían pero el producto quedaría sin mapear.',
    },
    {
        id: 'workflow',
        icon: '✨',
        title: 'Flujo recomendado',
        subtitle: 'Por dónde empezar',
        whatIs: 'El orden importa. Crear bien la base desde el principio ahorra problemas de mapeo de pedidos y control de stock.',
        whenToUse: 'Sigue este flujo al montar un nuevo catálogo o al agregar una familia de productos.',
        example: '1. (Opcional) Crea la familia en la pestaña "Familias de variantes"\n   → Define nombre, marca, categoría y ejes de variante\n\n2. Crea cada SKU en "SKUs / Productos"\n   → Código, nombre, precio, stock, imagen\n   → Si aplica: asigna la familia y sus atributos de variante\n\n3. Configura integraciones\n   → Click "Integraciones" en cada SKU\n   → Agrega el ID externo de cada plataforma\n\n4. Activa el SKU\n   → El badge "Estado" cambia a Activo\n   → Ya puede recibir pedidos y manejar stock',
        highlight: 'Puedes crear SKUs sin familia primero y luego agruparlos. La familia se puede asignar o cambiar editando el SKU.',
    },
];

const STORAGE_KEY = 'products_tour_seen_v1';

interface Props {
    isOpen: boolean;
    onClose: () => void;
}

export default function ProductTour({ isOpen, onClose }: Props) {
    const [currentStep, setCurrentStep] = useState(0);

    useEffect(() => {
        if (isOpen) setCurrentStep(0);
    }, [isOpen]);

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
        const handleKey = (e: KeyboardEvent) => {
            if (!isOpen) return;
            if (e.key === 'Escape') handleClose();
            if (e.key === 'ArrowRight') next();
            if (e.key === 'ArrowLeft') prev();
        };
        window.addEventListener('keydown', handleKey);
        return () => window.removeEventListener('keydown', handleKey);
    }, [isOpen, currentStep]);

    if (!isOpen) return null;

    const step = STEPS[currentStep];
    const progress = ((currentStep + 1) / STEPS.length) * 100;

    return (
        <div className="fixed inset-0 z-[60] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-3xl max-h-[90vh] overflow-hidden flex flex-col">
                <div className="h-1 bg-gray-100 dark:bg-gray-700">
                    <div
                        className="h-full transition-all duration-300"
                        style={{ width: `${progress}%`, background: 'var(--color-primary, #7c3aed)' }}
                    />
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
                    <Section title="Ejemplo real" body={step.example} accent="emerald" multiline />

                    {step.highlight && (
                        <div className="p-4 rounded-xl bg-gradient-to-br from-purple-50 to-indigo-50 dark:from-purple-900/20 dark:to-indigo-900/20 border border-purple-200 dark:border-purple-800">
                            <p className="text-sm text-purple-900 dark:text-purple-100 font-medium">&#128161; {step.highlight}</p>
                        </div>
                    )}

                    {step.id === 'families' && <VisualFamilyDiagram />}
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
                                className={`h-2 rounded-full transition-all ${
                                    i === currentStep
                                        ? 'w-6'
                                        : i < currentStep
                                            ? 'w-2 opacity-60'
                                            : 'w-2 bg-gray-300 dark:bg-gray-600'
                                }`}
                                style={i <= currentStep ? { background: 'var(--color-primary, #7c3aed)' } : {}}
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

function Section({ title, body, accent, multiline }: { title: string; body: string; accent: 'purple' | 'blue' | 'emerald'; multiline?: boolean }) {
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
            <p className={`text-sm text-gray-700 dark:text-gray-200 leading-relaxed ${multiline ? 'whitespace-pre-line' : ''}`}>{body}</p>
        </div>
    );
}

function VisualFamilyDiagram() {
    return (
        <div className="rounded-xl bg-gray-50 dark:bg-gray-900/40 border border-gray-200 dark:border-gray-700 p-4">
            <h4 className="text-xs uppercase tracking-wider font-bold text-gray-500 dark:text-gray-400 mb-3">Relacion familia - SKU</h4>
            <div className="flex flex-col gap-2">
                <div className="flex items-center gap-3 p-3 bg-purple-100 dark:bg-purple-900/30 rounded-lg">
                    <span className="text-xl">🗂️</span>
                    <div>
                        <p className="text-xs text-purple-600 dark:text-purple-300 uppercase font-bold">Familia</p>
                        <p className="text-sm font-mono font-semibold text-gray-900 dark:text-white">Proteina Whey</p>
                        <p className="text-xs text-gray-500">ejes: sabor, talla</p>
                    </div>
                </div>
                <div className="ml-6 flex flex-col gap-1.5">
                    {[
                        { sku: 'PT01001', label: 'Vainilla · 1kg', stock: 50 },
                        { sku: 'PT01002', label: 'Chocolate · 1kg', stock: 30 },
                        { sku: 'PT01003', label: 'Fresa · 2kg', stock: 0 },
                    ].map(v => (
                        <div key={v.sku} className="flex items-center gap-3 p-2.5 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
                            <div className="w-2 h-2 rounded-full bg-purple-400 flex-shrink-0" />
                            <span className="text-xs font-mono text-gray-500 w-16">{v.sku}</span>
                            <span className="text-sm font-medium text-gray-900 dark:text-white flex-1">{v.label}</span>
                            <span className={`text-xs px-2 py-0.5 rounded-full font-semibold ${v.stock === 0 ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
                                {v.stock === 0 ? 'Agotado' : `${v.stock} uds`}
                            </span>
                        </div>
                    ))}
                </div>
                <p className="text-xs text-gray-400 dark:text-gray-500 mt-1 pl-1">Cada SKU tiene su propio stock y precio. La familia no tiene stock.</p>
            </div>
        </div>
    );
}
