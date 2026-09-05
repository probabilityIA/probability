'use client';

import { AlertTriangle, Info, ShieldAlert } from 'lucide-react';

type Tone = 'danger' | 'warning' | 'info';

interface ActionConfirmDialogProps {
    isOpen: boolean;
    title: string;
    description: string;
    count: number;
    countLabel: string;
    warning?: string;
    tone?: Tone;
    confirmText: string;
    cancelText?: string;
    onConfirm: () => void;
    onCancel: () => void;
}

const TONES: Record<Tone, { icon: typeof Info; head: string; iconBox: string; iconColor: string; countBox: string; countText: string; button: string }> = {
    danger: {
        icon: ShieldAlert,
        head: 'from-red-50 to-orange-50 dark:from-red-950/40 dark:to-orange-950/40',
        iconBox: 'bg-red-500/10 dark:bg-red-400/10',
        iconColor: 'text-red-600 dark:text-red-400',
        countBox: 'border-red-200 dark:border-red-800 bg-red-50/60 dark:bg-red-900/20',
        countText: 'text-red-700 dark:text-red-300',
        button: 'bg-red-600 hover:bg-red-700',
    },
    warning: {
        icon: AlertTriangle,
        head: 'from-amber-50 to-yellow-50 dark:from-amber-950/40 dark:to-yellow-950/40',
        iconBox: 'bg-amber-500/10 dark:bg-amber-400/10',
        iconColor: 'text-amber-600 dark:text-amber-400',
        countBox: 'border-amber-200 dark:border-amber-800 bg-amber-50/60 dark:bg-amber-900/20',
        countText: 'text-amber-700 dark:text-amber-300',
        button: 'bg-amber-600 hover:bg-amber-700',
    },
    info: {
        icon: Info,
        head: 'from-violet-50 to-indigo-50 dark:from-violet-950/40 dark:to-indigo-950/40',
        iconBox: 'bg-violet-500/10 dark:bg-violet-400/10',
        iconColor: 'text-violet-600 dark:text-violet-400',
        countBox: 'border-violet-200 dark:border-violet-800 bg-violet-50/60 dark:bg-violet-900/20',
        countText: 'text-violet-700 dark:text-violet-300',
        button: 'bg-violet-600 hover:bg-violet-700',
    },
};

export function ActionConfirmDialog({
    isOpen,
    title,
    description,
    count,
    countLabel,
    warning,
    tone = 'info',
    confirmText,
    cancelText,
    onConfirm,
    onCancel,
}: ActionConfirmDialogProps) {
    if (!isOpen) return null;

    const t = TONES[tone];
    const Icon = t.icon;
    const nothingToDo = count === 0;

    return (
        <div className="fixed inset-0 z-[1200] flex items-center justify-center bg-black/70 backdrop-blur-sm p-4">
            <div className="bg-white dark:bg-gray-900 rounded-2xl shadow-2xl w-full max-w-md flex flex-col overflow-hidden border border-gray-200 dark:border-gray-700">
                <div className={`flex items-center gap-3 px-6 py-5 border-b border-gray-100 dark:border-gray-800 bg-gradient-to-r ${t.head}`}>
                    <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${t.iconBox}`}>
                        <Icon size={20} className={t.iconColor} />
                    </div>
                    <h3 className="text-base font-bold text-gray-900 dark:text-white">{title}</h3>
                </div>

                <div className="px-6 py-5 space-y-4">
                    <p className="text-sm text-gray-700 dark:text-gray-200 leading-relaxed">{description}</p>

                    <div className={`rounded-xl border px-4 py-3 flex items-baseline gap-2 ${t.countBox}`}>
                        <span className={`text-2xl font-bold tabular-nums ${t.countText}`}>{count}</span>
                        <span className="text-sm text-gray-600 dark:text-gray-300">{countLabel}</span>
                    </div>

                    {warning && (
                        <div className="flex items-start gap-2 rounded-lg bg-gray-50 dark:bg-gray-800/60 border border-gray-200 dark:border-gray-700 px-3 py-2">
                            <AlertTriangle size={15} className="text-amber-500 mt-0.5 shrink-0" />
                            <p className="text-xs text-gray-600 dark:text-gray-300 leading-relaxed">{warning}</p>
                        </div>
                    )}

                    {nothingToDo && (
                        <p className="text-xs text-gray-500 dark:text-gray-400">
                            {'No hay productos que procesar con esta acción.'}
                        </p>
                    )}
                </div>

                <div className="flex justify-end gap-3 px-6 py-4 border-t border-gray-100 dark:border-gray-800 bg-gray-50/60 dark:bg-gray-800/30">
                    <button
                        onClick={onCancel}
                        className="px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 text-sm font-semibold text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                    >
                        {cancelText || 'Cancelar'}
                    </button>
                    <button
                        onClick={onConfirm}
                        disabled={nothingToDo}
                        className={`px-5 py-2 rounded-lg text-sm font-semibold text-white transition-colors disabled:opacity-40 disabled:cursor-not-allowed ${t.button}`}
                    >
                        {confirmText}
                    </button>
                </div>
            </div>
        </div>
    );
}
