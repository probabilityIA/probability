'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { ArrowRight, Check, Loader2, RotateCcw } from 'lucide-react';
import { getProductMatchConfigAction, updateProductMatchConfigAction } from '../../infra/actions';
import type { ProductMatchConfig, ProductMatchField, ProductMatchRule } from '../../domain/types';

const FIELD_LABELS: Record<ProductMatchField, string> = {
    sku: 'SKU',
    barcode: 'Codigo de barras',
    external_id: 'ID externo',
    variant_id: 'ID de variante',
    name: 'Nombre',
};

interface ProductMatchRulesCardProps {
    integrationId: number;
    businessId?: number;
    channelName?: string;
}

function ruleKey(rule: ProductMatchRule) {
    return `${rule.probability}->${rule.channel}`;
}

export function ProductMatchRulesCard({ integrationId, businessId, channelName }: ProductMatchRulesCardProps) {
    const [config, setConfig] = useState<ProductMatchConfig | null>(null);
    const [rule, setRule] = useState<ProductMatchRule | null>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [message, setMessage] = useState<{ tone: 'ok' | 'error'; text: string } | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        const res = await getProductMatchConfigAction(integrationId, businessId) as {
            success?: boolean;
            message?: string;
            data?: ProductMatchConfig;
        };
        if (res?.success && res.data) {
            setConfig(res.data);
            setRule(res.data.rules[0] ?? null);
        } else {
            setMessage({ tone: 'error', text: res?.message || 'No se pudo cargar la configuracion de match' });
        }
        setLoading(false);
    }, [integrationId, businessId]);

    useEffect(() => {
        void load();
    }, [load]);

    const savedRule = config?.rules[0] ?? null;

    const dirty = useMemo(() => {
        if (!rule || !savedRule) return false;
        return ruleKey(rule) !== ruleKey(savedRule);
    }, [rule, savedRule]);

    const availableRules = useMemo(() => {
        if (!config) return [] as ProductMatchRule[];
        const candidates: ProductMatchRule[] = [];
        for (const probability of config.options.probability) {
            for (const channel of config.options.channel) {
                candidates.push({ probability, channel });
            }
        }
        return candidates;
    }, [config]);

    const save = async () => {
        if (!dirty || !rule) return;
        setSaving(true);
        setMessage(null);
        const res = await updateProductMatchConfigAction(integrationId, [rule], businessId) as {
            success?: boolean;
            message?: string;
            data?: ProductMatchConfig;
        };
        if (res?.success && res.data) {
            setConfig(res.data);
            setRule(res.data.rules[0] ?? null);
            setMessage({ tone: 'ok', text: 'Regla de match guardada' });
        } else {
            setMessage({ tone: 'error', text: res?.message || 'No se pudo guardar la regla' });
        }
        setSaving(false);
    };

    const restoreDefault = () => {
        if (!config) return;
        setRule(config.default_rules[0] ?? null);
    };

    if (loading) {
        return (
            <div className="flex items-center gap-2 rounded-lg border border-gray-200 bg-white p-4 text-sm text-gray-500 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400">
                <Loader2 size={14} className="animate-spin" />
                Cargando regla de match de producto...
            </div>
        );
    }

    if (!config) {
        return (
            <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/30 dark:text-red-300">
                {message?.text || 'No se pudo cargar la configuracion de match'}
            </div>
        );
    }

    return (
        <div className="space-y-3 rounded-lg border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800">
            <div>
                <h4 className="text-sm font-semibold text-gray-900 dark:text-white">Match de productos</h4>
                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    Elige el unico campo con el que se reconoce que un producto de Probability y uno de
                    {channelName ? ` ${channelName}` : ' este canal'} son el mismo. Solo se permite una regla: si no
                    coincide por ese campo, el producto no se considera el mismo.
                </p>
            </div>

            <div className="grid gap-1.5 sm:grid-cols-2">
                {availableRules.map(option => {
                    const selected = rule ? ruleKey(rule) === ruleKey(option) : false;
                    return (
                        <button
                            key={ruleKey(option)}
                            type="button"
                            role="radio"
                            aria-checked={selected}
                            disabled={saving}
                            onClick={() => setRule(option)}
                            className={`flex items-center gap-2 rounded-lg border px-3 py-2 text-left transition-colors disabled:opacity-50 ${
                                selected
                                    ? 'border-indigo-500 bg-indigo-50 dark:border-indigo-400 dark:bg-indigo-900/30'
                                    : 'border-gray-200 bg-gray-50 hover:border-indigo-300 dark:border-gray-600 dark:bg-gray-700/50'
                            }`}
                        >
                            <span
                                className={`flex h-4 w-4 shrink-0 items-center justify-center rounded-full border ${
                                    selected ? 'border-indigo-600 bg-indigo-600' : 'border-gray-400 dark:border-gray-500'
                                }`}
                            >
                                {selected && <Check size={10} className="text-white" />}
                            </span>
                            <span className="flex flex-1 items-center gap-2 text-sm text-gray-800 dark:text-gray-100">
                                <span className="rounded-md bg-indigo-50 px-2 py-0.5 text-xs font-semibold text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300">
                                    {FIELD_LABELS[option.probability]}
                                </span>
                                <ArrowRight size={12} className="text-gray-400" />
                                <span className="rounded-md bg-emerald-50 px-2 py-0.5 text-xs font-semibold text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300">
                                    {FIELD_LABELS[option.channel]}
                                </span>
                            </span>
                        </button>
                    );
                })}
            </div>

            {message && (
                <p className={`text-xs ${message.tone === 'ok' ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-600 dark:text-red-400'}`}>
                    {message.text}
                </p>
            )}

            <div className="flex items-center justify-between gap-2">
                <span className="text-[11px] text-gray-400 dark:text-gray-500">
                    {config.is_override ? 'Configuracion propia de esta integracion' : 'Usando el valor por defecto del canal'}
                </span>
                <div className="flex gap-2">
                    <button
                        type="button"
                        onClick={restoreDefault}
                        disabled={saving}
                        className="flex items-center gap-1 rounded-lg border border-gray-300 px-3 py-1.5 text-xs font-medium text-gray-600 transition-colors hover:bg-gray-50 disabled:opacity-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
                    >
                        <RotateCcw size={12} />
                        Valor por defecto
                    </button>
                    <button
                        type="button"
                        onClick={save}
                        disabled={saving || !dirty}
                        className="flex items-center gap-1 rounded-lg bg-indigo-600 px-3 py-1.5 text-xs font-semibold text-white transition-colors hover:bg-indigo-700 disabled:cursor-not-allowed disabled:opacity-50"
                    >
                        {saving ? <Loader2 size={12} className="animate-spin" /> : <Check size={12} />}
                        Guardar
                    </button>
                </div>
            </div>
        </div>
    );
}
