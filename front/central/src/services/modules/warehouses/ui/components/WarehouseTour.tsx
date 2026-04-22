'use client';

import { useState, useEffect } from 'react';
import {
    XMarkIcon,
    ChevronLeftIcon,
    ChevronRightIcon,
} from '@heroicons/react/24/outline';

interface TourStep {
    id: string;
    icon: string;
    title: string;
    subtitle: string;
    whatIs: string;
    whenToUse: string;
    example: string;
    highlight?: string;
    visual?: React.ReactNode;
}

const STEPS: TourStep[] = [
    {
        id: 'welcome',
        icon: '🏭',
        title: 'Bodegas e infraestructura física',
        subtitle: 'Cómo organizamos el espacio',
        whatIs: 'Una bodega no es solo "un cuarto con productos". Es una estructura física con múltiples niveles de organización que permiten saber exactamente dónde está cada producto.',
        whenToUse: 'Cada vez que creas una bodega nueva, necesitas definir esta jerarquía para que el sistema sepa dónde almacenar, mover y contar los productos.',
        example: 'Bodega "Principal" → Zona "Picking" → Pasillo "A" → Rack "R-01" → Nivel "L-02" → Posición "POS-A0102"',
        highlight: 'Los 5 niveles de la jerarquía (Zona → Pasillo → Rack → Nivel → Posición) existen porque cada uno cumple un propósito distinto. Recorramos cada uno.',
    },
    {
        id: 'zone',
        icon: '🗺️',
        title: 'Zona',
        subtitle: 'Área funcional de la bodega',
        whatIs: 'Una zona divide la bodega en áreas según el USO que tienen. Cada zona tiene un propósito distinto y puede tener color propio para identificarla visualmente.',
        whenToUse: 'Al organizar una bodega nueva, piensa en qué actividades haces y crea una zona por cada una.',
        example: 'Tipos típicos de zona:\n• Picking (preparación de pedidos)\n• Bulk (almacenamiento grande)\n• Recibo (entrada de mercancía)\n• Despacho (salida)\n• Cross-dock (paso directo)\n• Devoluciones\n• Cuarentena (revisión)\n• Averiado (producto dañado)',
        highlight: 'El propósito de la zona se usa luego en reglas de put-away ("pon estos productos en zona Picking") y cross-dock.',
    },
    {
        id: 'aisle',
        icon: '📏',
        title: 'Pasillo',
        subtitle: 'Corredor dentro de una zona',
        whatIs: 'Un pasillo es el corredor físico por donde camina la persona que opera la bodega. Separa columnas de racks.',
        whenToUse: 'Dentro de cada zona, crea un pasillo por cada corredor físico que tengas en la bodega real.',
        example: 'Zona Picking tiene 3 pasillos:\n• A-01: pasillo frontal (cerca del despacho)\n• A-02: pasillo central\n• A-03: pasillo trasero\n\nEsto ayuda a optimizar rutas de picking (empezar por los más cercanos).',
        highlight: 'El pasillo no almacena productos — es solo un agrupador. Los productos se guardan en los racks que están en el pasillo.',
    },
    {
        id: 'rack',
        icon: '🗄️',
        title: 'Rack',
        subtitle: 'Estantería física',
        whatIs: 'Un rack es la estantería (o anaquel) dentro de un pasillo. Tiene múltiples niveles verticales donde se colocan los productos.',
        whenToUse: 'Crea un rack por cada estantería física. Define cuántos niveles tiene (típico: 3-6 niveles).',
        example: 'Pasillo A-01 tiene 2 racks:\n• R-01 "Rack Principal" (5 niveles, metálico)\n• R-02 "Rack Secundario" (3 niveles, pallet)\n\nEl rack tiene medidas físicas (alto, ancho, profundidad) y un código visible para que el operario lo identifique.',
        highlight: 'Un rack puede tener niveles de diferente altura — usa esto para productos voluminosos vs pequeños.',
    },
    {
        id: 'level',
        icon: '📐',
        title: 'Nivel (piso del rack)',
        subtitle: 'Fila horizontal dentro del rack',
        whatIs: 'Un nivel es un piso horizontal del rack. Se numera con un ordinal (1=abajo, 2=medio, etc.). Define cuántas posiciones caben en él.',
        whenToUse: 'Un rack se compone de niveles. Típicamente 3-5 niveles verticales por rack. El nivel 1 (más bajo) suele ser para productos pesados, los superiores para livianos.',
        example: 'Rack R-01 tiene:\n• Nivel L-01 (ordinal 1): suelo, productos pesados\n• Nivel L-02 (ordinal 2): altura media, picking\n• Nivel L-03 (ordinal 3): parte alta, stock buffer\n\nUn nivel puede tener 1 o múltiples posiciones (depende del ancho).',
        highlight: 'La ordinal es importante para slotting ABC: Clase A suele ir en nivel 2 (a la altura de las manos, más rápido de tomar).',
    },
    {
        id: 'position',
        icon: '📍',
        title: 'Posición (ubicación final)',
        subtitle: 'La celda exacta donde va el producto',
        whatIs: 'Una posición es la celda concreta y única donde se guarda un producto. Es la unidad atómica del almacenamiento: tiene dimensiones, capacidad de peso y puede tener flags especiales.',
        whenToUse: 'Crea una posición por cada "slot" donde puedas colocar un producto. Es lo que el operario escanea cuando guarda o saca.',
        example: 'Posición POS-A0102 (Pasillo A, Rack 01, Nivel 02, ordinal 2):\n• Dimensiones: 120 cm × 80 cm × 50 cm\n• Peso máximo: 50 kg\n• Tipo: picking\n• Flags: is_picking=true\n\nAquí van 6 unidades del SKU-123 lote LOT-005.',
        highlight: 'La posición es el único nivel que REALMENTE almacena productos. Los demás (zona, pasillo, rack, nivel) son solo organizadores.',
    },
    {
        id: 'flags',
        icon: '🏷️',
        title: 'Flags especiales de la posición',
        subtitle: 'Propósitos específicos',
        whatIs: 'Cada posición puede tener flags que indican para qué sirve. Esto permite al sistema decidir automáticamente dónde poner un producto.',
        whenToUse: 'Al crear posiciones, marca los flags que apliquen según el uso real de esa celda en la bodega.',
        example: 'Flags disponibles:\n• is_picking: posición para picking rápido\n• is_bulk: almacenamiento grande\n• is_quarantine: cuarentena de revisión\n• is_damaged: producto averiado\n• is_returns: devoluciones\n• is_cross_dock: paso directo sin almacenar\n• is_hazmat: material peligroso',
        highlight: 'El put-away automático usa estos flags: si un producto llega con el tag "frágil", el sistema no lo pone en is_bulk (posiciones pesadas).',
    },
    {
        id: 'lpn',
        icon: '📦',
        title: 'LPN (opcional)',
        subtitle: 'Contenedor lógico sobre la posición',
        whatIs: 'Una License Plate Number (LPN) es un pallet/caja/tote con código único que agrupa productos. Está EN una posición pero se mueve como unidad.',
        whenToUse: 'Cuando recibes un pallet completo del proveedor o preparas un contenedor para despacho. Mover un LPN mueve todo su contenido.',
        example: 'Posición POS-A0102 contiene el pallet LPN-PAL-042 que tiene:\n• 3 cajas de SKU-A\n• 2 cajas de SKU-B (lote L-001)\n\nAl escanear LPN-PAL-042 y moverlo a POS-B0203, todo se actualiza en un paso.',
        highlight: 'LPN es opcional. Los productos pueden estar "sueltos" en una posición o dentro de un LPN. Ambos flujos coexisten.',
    },
    {
        id: 'hierarchy-summary',
        icon: '🧩',
        title: 'Todo junto',
        subtitle: 'La ruta completa',
        whatIs: 'La dirección física completa de un producto se arma con los 5 (o 6 con LPN) niveles.',
        whenToUse: 'Cuando el operario escanea un código, el sistema puede decirle exactamente "camina al pasillo A-01, busca el rack R-02, sube al nivel 3, toma el producto de la posición POS-A0103".',
        example: 'Dirección completa de 10 camisetas rojas:\nBodega Principal > Zona Picking > Pasillo A-01 > Rack R-02 > Nivel L-03 (1.8m altura) > Posición POS-A0103 > LPN-TOT-015',
        highlight: 'Esta estructura es lo que permite un WMS real: conteos cíclicos precisos, slotting ABC, put-away automático, rutas de picking optimizadas.',
    },
    {
        id: 'how-to-create',
        icon: '✨',
        title: 'Cómo crear desde esta tabla',
        subtitle: 'Flujo recomendado',
        whatIs: 'Esta tabla te permite crear y editar toda la jerarquía sin salir. Expande cada bodega con el chevron → aparece el árbol con botones inline.',
        whenToUse: 'Orden recomendado al montar una bodega nueva: Bodega → Zonas → Pasillos → Racks → Niveles → Posiciones.',
        example: '1. Click "+ Nueva bodega" para crear la estructura raíz.\n2. Expande la bodega (chevron).\n3. Click "+ Nueva zona" dentro del árbol.\n4. Hover sobre la zona creada → click "+" para agregar pasillo.\n5. Hover sobre el pasillo → click "+" para agregar rack.\n6. Y así hasta posiciones.\n\nCada fila tiene botones de editar (amarillo) y eliminar (rojo) al pasar el mouse.',
        highlight: 'Puedes editar el código, nombre o propósito de cualquier nivel sin romper los hijos. Los productos y movimientos siguen apuntando a los IDs internos.',
    },
    {
        id: 'integration',
        icon: '🔌',
        title: 'Cómo se conecta con los demás módulos',
        subtitle: 'No es una tabla aislada',
        whatIs: 'La jerarquía que construyas aquí alimenta a todos los módulos operativos. Los códigos de posiciones se usan en ajustes, transferencias, kardex, put-away, conteos cíclicos y scanning.',
        whenToUse: 'Cada vez que haces un ajuste de stock, el sistema te pide el producto + ubicación (posición). Sin esta jerarquía, el stock queda "flotando" sin lugar físico.',
        example: '• Ajustar stock → selecciona la posición donde aplicas el +30\n• Transferir stock → indica posición origen y destino\n• Kardex → muestra en qué posición ocurrió cada movimiento\n• Conteo cíclico → genera lista de qué contar en qué posiciones\n• Put-away → sugiere posición según reglas\n• Scan → resuelve el código de la posición al buscarla',
        highlight: 'El tiempo invertido en configurar bien la jerarquía se recupera en minutos al día de picking y conteos bien hechos.',
    },
];

const STORAGE_KEY = 'warehouse_tour_seen_v1';

interface Props {
    isOpen: boolean;
    onClose: () => void;
}

export default function WarehouseTour({ isOpen, onClose }: Props) {
    const [currentStep, setCurrentStep] = useState(0);

    useEffect(() => {
        if (isOpen) setCurrentStep(0);
    }, [isOpen]);

    useEffect(() => {
        const handleEsc = (e: KeyboardEvent) => {
            if (!isOpen) return;
            if (e.key === 'Escape') handleClose();
            if (e.key === 'ArrowRight') next();
            if (e.key === 'ArrowLeft') prev();
        };
        window.addEventListener('keydown', handleEsc);
        return () => window.removeEventListener('keydown', handleEsc);
    }, [isOpen, currentStep]);

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
                            <p className="text-sm text-purple-900 dark:text-purple-100 font-medium">💡 {step.highlight}</p>
                        </div>
                    )}

                    {step.id === 'hierarchy-summary' && <VisualHierarchy />}
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

function VisualHierarchy() {
    const levels = [
        { icon: '🏭', label: 'Bodega', ex: 'Bodega Principal', color: 'bg-slate-600' },
        { icon: '🗺️', label: 'Zona', ex: 'Z-PICK (Picking)', color: 'bg-indigo-600' },
        { icon: '📏', label: 'Pasillo', ex: 'A-01 (Pasillo Frontal)', color: 'bg-emerald-600' },
        { icon: '🗄️', label: 'Rack', ex: 'R-02 (Rack Secundario)', color: 'bg-purple-600' },
        { icon: '📐', label: 'Nivel', ex: 'L-03 (Ordinal 3)', color: 'bg-amber-600' },
        { icon: '📍', label: 'Posición', ex: 'POS-A0103 (shelf)', color: 'bg-rose-600' },
    ];
    return (
        <div className="rounded-xl bg-gray-50 dark:bg-gray-900/40 border border-gray-200 dark:border-gray-700 p-4">
            <h4 className="text-xs uppercase tracking-wider font-bold text-gray-500 dark:text-gray-400 mb-3">Ruta visual</h4>
            <div className="space-y-2">
                {levels.map((l, i) => (
                    <div key={l.label} className="flex items-center gap-3" style={{ paddingLeft: `${i * 20}px` }}>
                        <span className={`${l.color} text-white w-8 h-8 rounded-lg flex items-center justify-center text-base shadow-sm`}>{l.icon}</span>
                        <div className="flex-1">
                            <p className="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider">{l.label}</p>
                            <p className="text-sm font-medium text-gray-900 dark:text-white font-mono">{l.ex}</p>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
