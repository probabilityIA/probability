'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { ArrowRight, Check, GripVertical, Loader2, Plus, RotateCcw, Trash2 } from 'lucide-react';
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

function sameRules(a: ProductMatchRule[], b: ProductMatchRule[]) {
    if (a.length !== b.length) return false;
    return a.every((rule, index) => ruleKey(rule) === ruleKey(b[index]));
}

export function ProductMatchRulesCard({ integrationId, businessId, channelName }: ProductMatchRulesCardProps) {
    const [config, setConfig] = useState<ProductMatchConfig | null>(null);
    const [rules, setRules] = useState<ProductMatchRule[]>([]);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [message, setMessage] = useState<{ tone: 'ok' | 'error'; text: string } | null>(null);
    const [dragIndex, setDragIndex] = useState<number | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        const res = await getProductMatchConfigAction(integrationId, businessId) as {
            success?: boolean;
            message?: string;
            data?: ProductMatchConfig;
        };
        if (res?.success && res.data) {
            setConfig(res.data);
            setRules(res.data.rules);
        } else {
            setMessage({ tone: 'error', text: res?.message || 'No se pudo cargar la configuracion de match' });
        }
        setLoading(false);
    }, [integrationId, businessId]);

    useEffect(() => {
        void load();
    }, [load]);

    const dirty = useMemo(() => (config ? !sameRules(rules, config.rules) : false), [config, rules]);

    const usedKeys = useMemo(() => new Set(rules.map(ruleKey)), [rules]);

    const availableToAdd = useMemo(() => {
        if (!config) return [] as ProductMatchRule[];
        const candidates: ProductMatchRule[] = [];
        for (const probability of config.options.probability) {
            for (const channel of config.options.channel) {
                const rule = { probability, channel };
                if (!usedKeys.has(ruleKey(rule))) candidates.push(rule);
            }
        }
        return candidates;
    }, [config, usedKeys]);

    const move = (from: number, to: number) => {
        if (to < 0 || to >= rules.length || from === to) return;
        setRules(prev => {
            const next = [...prev];
            const [item] = next.splice(from, 1);
            next.splice(to, 0, item);
            return next;
        });
    };

    const save = async () => {
        if (!dirty || rules.length === 0) return;
        setSaving(true);
        setMessage(null);
        const res = await updateProductMatchConfigAction(integrationId, rules, businessId) as {
            success?: boolean;
            message?: string;
            data?: ProductMatchConfig;
        };
        if (res?.success && res.data) {
            setConfig(res.data);
            setRules(res.data.rules);
            setMessage({ tone: 'ok', text: 'Reglas de match guardadas' });
        } else {
            setMessage({ tone: 'error', text: res?.message || 'No se pudieron guardar las reglas' });
        }
        setSaving(false);
    };

    const restoreDefaults = () => {
        if (!config) return;
        setRules(config.default_rules);
    };

    if (loading) {
        return (
            <div className="flex items-center gap-2 rounded-lg border border-gray-200 bg-white p-4 text-sm text-gray-500 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400">
                <Loader2 size={14} className="animate-spin" />
                Cargando reglas de match de producto...
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
                    Define por que campos se reconoce que un producto de Probability y uno de
                    {channelName ? ` ${channelName}` : ' este canal'} son el mismo. Se evaluan en orden: si la
                    primera regla no encuentra coincidencia, se intenta la siguiente.
                </p>
            </div>

            <ul className="space-y-2">
                {rules.map((rule, index) => (
                    <li
                        key={ruleKey(rule)}
                        draggable={!saving}
                        onDragStart={() => setDragIndex(index)}
                        onDragOver={event => event.preventDefault()}
                        onDrop={() => {
                            if (dragIndex !== null) move(dragIndex, index);
                            setDragIndex(null);
                        }}
                        className="flex items-center gap-2 rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 dark:border-gray-600 dark:bg-gray-700/50"
                    >
                        <GripVertical size={14} className="cursor-grab text-gray-400" />
                        <span className="w-5 text-xs font-semibold text-gray-400">{index + 1}</span>
                        <span className="flex flex-1 items-center gap-2 text-sm text-gray-800 dark:text-gray-100">
                            <span className="rounded-md bg-indigo-50 px-2 py-0.5 text-xs font-semibold text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300">
                                {FIELD_LABELS[rule.probability]}
                            </span>
                            <ArrowRight size={12} className="text-gray-400" />
                            <span className="rounded-md bg-emerald-50 px-2 py-0.5 text-xs font-semibold text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300">
                                {FIELD_LABELS[rule.channel]}
                            </span>
                        </span>
                        <button
                            type="button"
                            onClick={() => setRules(prev => prev.filter((_, i) => i !== index))}
                            disabled={saving || rules.length === 1}
                            title={rules.length === 1 ? 'Debe quedar al menos una regla' : 'Quitar regla'}
                            className="text-gray-400 transition-colors hover:text-red-600 disabled:cursor-not-allowed disabled:opacity-40"
                        >
                            <Trash2 size={14} />
                        </button>
                    </li>
                ))}
            </ul>

            {availableToAdd.length > 0 && (
                <div className="flex flex-wrap gap-1.5">
                    {availableToAdd.map(rule => (
                        <button
                            key={ruleKey(rule)}
                            type="button"
                            disabled={saving}
                            onClick={() => setRules(prev => [...prev, rule])}
                            className="flex items-center gap-1 rounded-full border border-dashed border-gray-300 px-2 py-1 text-[11px] text-gray-600 transition-colors hover:border-indigo-400 hover:text-indigo-700 disabled:opacity-50 dark:border-gray-600 dark:text-gray-300"
                        >
                            <Plus size={11} />
                            {FIELD_LABELS[rule.probability]} → {FIELD_LABELS[rule.channel]}
                        </button>
                    ))}
                </div>
            )}

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
                        onClick={restoreDefaults}
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
