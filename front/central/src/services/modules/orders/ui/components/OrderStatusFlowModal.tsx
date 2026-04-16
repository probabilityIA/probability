'use client';

import { useState } from 'react';
import { Modal } from '@/shared/ui/modal';
import {
    STATUS_FLOW,
    CATEGORY_LABELS,
    CATEGORY_ORDER,
    type StatusStep,
} from '../../domain/order-status-transitions';

interface OrderStatusFlowModalProps {
    isOpen: boolean;
    onClose: () => void;
}

export function OrderStatusFlowModal({ isOpen, onClose }: OrderStatusFlowModalProps) {
    const [selectedStatus, setSelectedStatus] = useState<string | null>(null);

    const mainFlow = STATUS_FLOW.filter(s => s.category !== 'issue');
    const issueStatuses = STATUS_FLOW.filter(s => s.category === 'issue');

    const groupedByCategory = CATEGORY_ORDER.reduce((acc, cat) => {
        const items = STATUS_FLOW.filter(s => s.category === cat);
        if (items.length > 0) acc[cat] = items;
        return acc;
    }, {} as Record<string, StatusStep[]>);

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Flujo de Estados de Orden" size="4xl">
            <div className="space-y-6">
                {/* Descripcion general */}
                <div className="bg-purple-50 dark:bg-purple-900/20 border border-purple-200 dark:border-purple-700 rounded-xl p-4">
                    <p className="text-sm text-purple-800 dark:text-purple-200">
                        Las ordenes siguen un flujo logistico secuencial. Cada estado representa una etapa del proceso,
                        desde la recepcion del pedido hasta la entrega o devolucion. Los estados terminales
                        (<strong>Cancelada</strong> y <strong>Reembolsada</strong>) no permiten mas cambios.
                    </p>
                </div>

                {/* Flujo principal visual */}
                <div>
                    <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3 uppercase tracking-wide">Flujo Principal</h4>
                    <div className="flex flex-wrap items-center gap-1">
                        {mainFlow.map((step, i) => {
                            const isSelected = selectedStatus === step.code;
                            const isTerminal = step.category === 'final';
                            return (
                                <div key={step.code} className="flex items-center">
                                    <button
                                        onClick={() => setSelectedStatus(isSelected ? null : step.code)}
                                        className={`relative px-3 py-1.5 text-xs font-medium rounded-full transition-all cursor-pointer border-2 ${
                                            isSelected
                                                ? 'ring-2 ring-offset-2 ring-purple-500 scale-110'
                                                : 'hover:scale-105'
                                        } ${isTerminal ? 'border-dashed' : 'border-transparent'}`}
                                        style={{
                                            backgroundColor: step.color + '20',
                                            color: step.color,
                                            borderColor: isSelected ? step.color : isTerminal ? step.color : 'transparent',
                                        }}
                                    >
                                        {step.name}
                                    </button>
                                    {i < mainFlow.length - 1 && (
                                        <svg className="w-4 h-4 text-gray-300 dark:text-gray-600 mx-0.5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                        </svg>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                </div>

                {/* Novedades / ramificaciones */}
                <div>
                    <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3 uppercase tracking-wide">Novedades (ramificaciones del flujo)</h4>
                    <div className="flex flex-wrap gap-2">
                        {issueStatuses.map((step) => {
                            const isSelected = selectedStatus === step.code;
                            return (
                                <button
                                    key={step.code}
                                    onClick={() => setSelectedStatus(isSelected ? null : step.code)}
                                    className={`px-3 py-1.5 text-xs font-medium rounded-full transition-all cursor-pointer border-2 border-dashed ${
                                        isSelected ? 'ring-2 ring-offset-2 ring-orange-400 scale-110' : 'hover:scale-105'
                                    }`}
                                    style={{
                                        backgroundColor: step.color + '20',
                                        color: step.color,
                                        borderColor: step.color,
                                    }}
                                >
                                    {step.name}
                                </button>
                            );
                        })}
                    </div>
                </div>

                {/* Detalle del estado seleccionado */}
                {selectedStatus && (
                    <div className="bg-gray-50 dark:bg-gray-700/50 rounded-xl p-4 border border-gray-200 dark:border-gray-600 animate-in fade-in duration-200">
                        {(() => {
                            const step = STATUS_FLOW.find(s => s.code === selectedStatus);
                            if (!step) return null;
                            return (
                                <div className="flex items-start gap-3">
                                    <div
                                        className="w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5"
                                        style={{ backgroundColor: step.color + '20' }}
                                    >
                                        <div className="w-4 h-4 rounded-full" style={{ backgroundColor: step.color }} />
                                    </div>
                                    <div>
                                        <p className="font-semibold text-gray-900 dark:text-white">{step.name}</p>
                                        <p className="text-xs text-gray-500 dark:text-gray-400 font-mono">{step.code}</p>
                                        <p className="text-sm text-gray-700 dark:text-gray-300 mt-1">{step.description}</p>
                                    </div>
                                </div>
                            );
                        })()}
                    </div>
                )}

                {/* Tabla por categorias */}
                <div>
                    <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3 uppercase tracking-wide">Detalle por Fase</h4>
                    <div className="grid gap-3">
                        {CATEGORY_ORDER.map((cat) => {
                            const items = groupedByCategory[cat];
                            if (!items) return null;
                            const catInfo = CATEGORY_LABELS[cat];
                            return (
                                <div key={cat} className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                                    <div className="bg-gray-100 dark:bg-gray-700 px-4 py-2 flex items-center gap-2">
                                        <span className="w-6 h-6 rounded-full bg-purple-600 text-white text-xs font-bold flex items-center justify-center">
                                            {catInfo.icon}
                                        </span>
                                        <span className="text-sm font-semibold text-gray-800 dark:text-gray-200">{catInfo.label}</span>
                                    </div>
                                    <div className="divide-y divide-gray-100 dark:divide-gray-700">
                                        {items.map((step) => (
                                            <div
                                                key={step.code}
                                                className={`flex items-center gap-3 px-4 py-2.5 transition-colors cursor-pointer ${
                                                    selectedStatus === step.code
                                                        ? 'bg-purple-50 dark:bg-purple-900/20'
                                                        : 'hover:bg-gray-50 dark:hover:bg-gray-800/50'
                                                }`}
                                                onClick={() => setSelectedStatus(selectedStatus === step.code ? null : step.code)}
                                            >
                                                <div
                                                    className="w-3 h-3 rounded-full flex-shrink-0"
                                                    style={{ backgroundColor: step.color }}
                                                />
                                                <div className="flex-1 min-w-0">
                                                    <div className="flex items-center gap-2">
                                                        <span className="text-sm font-medium text-gray-900 dark:text-white">{step.name}</span>
                                                        <span className="text-xs text-gray-400 font-mono">{step.code}</span>
                                                    </div>
                                                    <p className="text-xs text-gray-500 dark:text-gray-400 truncate">{step.description}</p>
                                                </div>
                                                <span
                                                    className="px-2 py-0.5 text-xs font-medium rounded-full flex-shrink-0"
                                                    style={{
                                                        backgroundColor: step.color + '20',
                                                        color: step.color,
                                                    }}
                                                >
                                                    {step.name}
                                                </span>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                </div>

                {/* Reglas especiales */}
                <div className="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-700 rounded-xl p-4">
                    <h4 className="text-sm font-semibold text-amber-800 dark:text-amber-200 mb-2">Reglas Especiales</h4>
                    <ul className="text-sm text-amber-700 dark:text-amber-300 space-y-1 list-disc list-inside">
                        <li><strong>Cancelada</strong> se puede aplicar desde cualquier estado no terminal</li>
                        <li><strong>Reembolsada</strong> es posible desde Entregada, Completada o Devuelto</li>
                        <li><strong>Novedad de entrega</strong> permite reintentar asignando un nuevo piloto</li>
                        <li><strong>Novedad de inventario</strong> permite volver a Picking cuando se resuelve</li>
                        <li><strong>En espera</strong> puede volver a Pendiente para continuar el flujo</li>
                        <li>Los estados <strong>Cancelada</strong> y <strong>Reembolsada</strong> son terminales (sin retorno)</li>
                    </ul>
                </div>
            </div>
        </Modal>
    );
}
