'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Concept, DianConfig, Tax } from '../../domain/types';
import {
    setAccountingConceptTaxAction,
    updateAccountingConceptAction,
    updateAccountingTaxAction,
} from '../../infra/actions';
import { kindLabel, taxKindBadgeClass, taxKindLabel } from '../format';
import { useToast } from '@/shared/providers/toast-provider';
import { ConceptFormModal } from './ConceptFormModal';
import { TaxFormModal } from './TaxFormModal';
import { DianConfigModal } from './DianConfigModal';

interface ConfigViewProps {
    concepts: Concept[];
    taxes: Tax[];
    dianConfig: DianConfig | null;
    error: string | null;
}

function Switch({ checked, disabled, onChange, label }: { checked: boolean; disabled?: boolean; onChange: () => void; label: string }) {
    return (
        <button
            type="button"
            role="switch"
            aria-checked={checked}
            aria-label={label}
            title={label}
            disabled={disabled}
            onClick={onChange}
            className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors disabled:opacity-50 ${checked ? 'bg-purple-600' : 'bg-gray-300 dark:bg-gray-600'}`}
        >
            <span className={`inline-block h-3.5 w-3.5 transform rounded-full bg-white transition-transform ${checked ? 'translate-x-5' : 'translate-x-1'}`} />
        </button>
    );
}

export function ConfigView({ concepts, taxes, dianConfig, error }: ConfigViewProps) {
    const router = useRouter();
    const { showToast } = useToast();
    const [conceptModalOpen, setConceptModalOpen] = useState(false);
    const [conceptToEdit, setConceptToEdit] = useState<Concept | null>(null);
    const [taxModalOpen, setTaxModalOpen] = useState(false);
    const [taxToEdit, setTaxToEdit] = useState<Tax | null>(null);
    const [dianModalOpen, setDianModalOpen] = useState(false);
    const [busyKey, setBusyKey] = useState<string | null>(null);

    const refresh = () => router.refresh();

    const toggleConcept = async (concept: Concept, field: 'is_active' | 'is_real_income') => {
        const key = `concept-${concept.id}-${field}`;
        setBusyKey(key);
        const result = await updateAccountingConceptAction(concept.id, {
            name: concept.name,
            description: concept.description,
            is_real_income: field === 'is_real_income' ? !concept.is_real_income : concept.is_real_income,
            is_active: field === 'is_active' ? !concept.is_active : concept.is_active,
        });
        setBusyKey(null);
        if (result.success) {
            refresh();
        } else {
            showToast(result.error, 'error');
        }
    };

    const toggleConceptTax = async (concept: Concept, tax: Tax) => {
        const key = `concept-${concept.id}-tax-${tax.id}`;
        const current = concept.taxes?.find((t) => t.tax_id === tax.id);
        const nextActive = !(current?.is_active ?? false);
        setBusyKey(key);
        const result = await setAccountingConceptTaxAction(concept.id, tax.id, nextActive);
        setBusyKey(null);
        if (result.success) {
            refresh();
        } else {
            showToast(result.error, 'error');
        }
    };

    const toggleTax = async (tax: Tax) => {
        const key = `tax-${tax.id}`;
        setBusyKey(key);
        const result = await updateAccountingTaxAction(tax.id, {
            name: tax.name,
            description: tax.description,
            rate_percent: tax.rate_percent,
            kind: tax.kind,
            is_active: !tax.is_active,
        });
        setBusyKey(null);
        if (result.success) {
            refresh();
        } else {
            showToast(result.error, 'error');
        }
    };

    return (
        <div className="space-y-8">
            <div>
                <h1 className="text-xl font-semibold text-gray-900 dark:text-white">Configuracion contable</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                    Conceptos e impuestos de la contabilidad interna
                </p>
            </div>

            {error && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 text-sm rounded-lg px-4 py-3">
                    {error}
                </div>
            )}

            <section className="space-y-3">
                <div className="flex items-center justify-between">
                    <h2 className="text-base font-semibold text-gray-900 dark:text-white">Conceptos</h2>
                    <button
                        onClick={() => {
                            setConceptToEdit(null);
                            setConceptModalOpen(true);
                        }}
                        className="inline-flex items-center gap-2 px-3 py-2 text-sm font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-colors"
                    >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                        </svg>
                        Nuevo concepto
                    </button>
                </div>

                <div className="space-y-3">
                    {concepts.length === 0 && (
                        <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-6 text-center text-sm text-gray-500 dark:text-gray-400">
                            No hay conceptos registrados
                        </div>
                    )}
                    {concepts.map((concept) => (
                        <div key={concept.id} className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm p-4">
                            <div className="flex flex-col lg:flex-row lg:items-center gap-4">
                                <div className="flex-1 min-w-0">
                                    <div className="flex flex-wrap items-center gap-2">
                                        <span className="font-medium text-gray-900 dark:text-white">{concept.name}</span>
                                        <span className="text-xs text-gray-400">{concept.code}</span>
                                        <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-semibold ${concept.kind === 'INCOME' ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300' : 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300'}`}>
                                            {kindLabel(concept.kind)}
                                        </span>
                                        {concept.is_automatic && (
                                            <span className="inline-block px-2 py-0.5 rounded-full text-xs font-semibold bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300">
                                                Automatico{concept.source_type ? `: ${concept.source_type}` : ''}
                                            </span>
                                        )}
                                    </div>
                                    {concept.description && (
                                        <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">{concept.description}</p>
                                    )}
                                    {taxes.length > 0 && (
                                        <div className="mt-3 flex flex-wrap items-center gap-2">
                                            <span className="text-xs font-medium text-gray-500 dark:text-gray-400">Impuestos:</span>
                                            {taxes.map((tax) => {
                                                const active = concept.taxes?.find((t) => t.tax_id === tax.id)?.is_active ?? false;
                                                const busy = busyKey === `concept-${concept.id}-tax-${tax.id}`;
                                                return (
                                                    <button
                                                        key={tax.id}
                                                        onClick={() => toggleConceptTax(concept, tax)}
                                                        disabled={busy}
                                                        title={`${tax.name} (${tax.rate_percent}%)`}
                                                        className={`px-2.5 py-1 rounded-full text-xs font-semibold border transition-colors disabled:opacity-50 ${active
                                                            ? 'bg-purple-600 border-purple-600 text-white'
                                                            : 'bg-transparent border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 hover:border-purple-400'
                                                        }`}
                                                    >
                                                        {tax.code} {tax.rate_percent}%
                                                    </button>
                                                );
                                            })}
                                        </div>
                                    )}
                                </div>
                                <div className="flex items-center gap-5 shrink-0">
                                    <div className="flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                                        <span>Ganancia real</span>
                                        <Switch
                                            checked={concept.is_real_income}
                                            disabled={busyKey === `concept-${concept.id}-is_real_income`}
                                            onChange={() => toggleConcept(concept, 'is_real_income')}
                                            label="Es ganancia real"
                                        />
                                    </div>
                                    <div className="flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                                        <span>Activo</span>
                                        <Switch
                                            checked={concept.is_active}
                                            disabled={busyKey === `concept-${concept.id}-is_active`}
                                            onChange={() => toggleConcept(concept, 'is_active')}
                                            label="Concepto activo"
                                        />
                                    </div>
                                    <button
                                        onClick={() => {
                                            setConceptToEdit(concept);
                                            setConceptModalOpen(true);
                                        }}
                                        className="p-2 rounded-lg text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                                        title="Editar concepto"
                                    >
                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                        </svg>
                                    </button>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </section>

            <section className="space-y-3">
                <div className="flex items-center justify-between">
                    <h2 className="text-base font-semibold text-gray-900 dark:text-white">Impuestos</h2>
                    <button
                        onClick={() => {
                            setTaxToEdit(null);
                            setTaxModalOpen(true);
                        }}
                        className="inline-flex items-center gap-2 px-3 py-2 text-sm font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-colors"
                    >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                        </svg>
                        Nuevo impuesto
                    </button>
                </div>

                <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm overflow-hidden">
                    <div className="overflow-x-auto">
                        <table className="w-full text-sm">
                            <thead>
                                <tr className="text-left text-xs text-gray-500 dark:text-gray-400 uppercase bg-gray-50 dark:bg-gray-900/40">
                                    <th className="px-4 py-2 font-medium">Codigo</th>
                                    <th className="px-4 py-2 font-medium">Nombre</th>
                                    <th className="px-4 py-2 font-medium">Tipo</th>
                                    <th className="px-4 py-2 font-medium text-right">Tarifa</th>
                                    <th className="px-4 py-2 font-medium text-center">Activo</th>
                                    <th className="px-4 py-2 font-medium text-right">Acciones</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                {taxes.length === 0 && (
                                    <tr>
                                        <td colSpan={6} className="px-4 py-6 text-center text-gray-500 dark:text-gray-400">
                                            No hay impuestos registrados
                                        </td>
                                    </tr>
                                )}
                                {taxes.map((tax) => (
                                    <tr key={tax.id} className="text-gray-800 dark:text-gray-200">
                                        <td className="px-4 py-2 font-mono text-xs">{tax.code}</td>
                                        <td className="px-4 py-2">
                                            <span className="font-medium">{tax.name}</span>
                                            {tax.description && (
                                                <p className="text-xs text-gray-500 dark:text-gray-400">{tax.description}</p>
                                            )}
                                        </td>
                                        <td className="px-4 py-2">
                                            <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-semibold ${taxKindBadgeClass(tax.kind)}`}>
                                                {taxKindLabel(tax.kind)}
                                            </span>
                                        </td>
                                        <td className="px-4 py-2 text-right">{tax.rate_percent}%</td>
                                        <td className="px-4 py-2 text-center">
                                            <Switch
                                                checked={tax.is_active}
                                                disabled={busyKey === `tax-${tax.id}`}
                                                onChange={() => toggleTax(tax)}
                                                label="Impuesto activo"
                                            />
                                        </td>
                                        <td className="px-4 py-2 text-right">
                                            <button
                                                onClick={() => {
                                                    setTaxToEdit(tax);
                                                    setTaxModalOpen(true);
                                                }}
                                                className="p-2 rounded-lg text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                                                title="Editar impuesto"
                                            >
                                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                                </svg>
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </section>

            <section className="space-y-3">
                <h2 className="text-base font-semibold text-gray-900 dark:text-white">Facturacion electronica (DIAN)</h2>
                <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm p-4">
                    <div className="flex flex-col sm:flex-row sm:items-center gap-4">
                        <div className="flex-1 min-w-0">
                            <div className="flex flex-wrap items-center gap-2">
                                <span className="font-medium text-gray-900 dark:text-white">Factus</span>
                                <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-semibold ${dianConfig?.configured
                                    ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300'
                                    : 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
                                }`}>
                                    {dianConfig?.configured ? 'Configurada' : 'Sin configurar'}
                                </span>
                            </div>
                            {dianConfig?.configured ? (
                                <div className="mt-2 space-y-0.5 text-xs text-gray-500 dark:text-gray-400">
                                    <p>API: {dianConfig.api_url || 'Produccion'}</p>
                                    <p>Usuario: {dianConfig.username || '-'}</p>
                                    <p>Rango de numeracion: {dianConfig.numbering_range_id || 'Automatico'}</p>
                                    <p>Medio de pago: {dianConfig.payment_method_code ? `Codigo ${dianConfig.payment_method_code}` : 'No definido'}</p>
                                </div>
                            ) : (
                                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                                    Conecta las credenciales de Factus para emitir facturas electronicas ante la DIAN
                                </p>
                            )}
                        </div>
                        <button
                            onClick={() => setDianModalOpen(true)}
                            className="shrink-0 inline-flex items-center gap-2 px-3 py-2 text-sm font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-colors"
                        >
                            {dianConfig?.configured ? 'Editar' : 'Configurar'}
                        </button>
                    </div>
                </div>
            </section>

            <ConceptFormModal
                isOpen={conceptModalOpen}
                onClose={() => setConceptModalOpen(false)}
                onSaved={() => {
                    setConceptModalOpen(false);
                    refresh();
                }}
                concept={conceptToEdit}
            />

            <TaxFormModal
                isOpen={taxModalOpen}
                onClose={() => setTaxModalOpen(false)}
                onSaved={() => {
                    setTaxModalOpen(false);
                    refresh();
                }}
                tax={taxToEdit}
            />

            <DianConfigModal
                isOpen={dianModalOpen}
                onClose={() => setDianModalOpen(false)}
                onSaved={() => {
                    setDianModalOpen(false);
                    refresh();
                }}
                config={dianConfig}
            />
        </div>
    );
}
