'use client';

import { useState } from 'react';
import { Mail, RefreshCw, Send, X } from 'lucide-react';
import { sendCodCutEmailAction } from '../../infra/actions';
import { PaymentCut } from '../../domain/types';
import { formatMoney, formatDateOnly } from './helpers';

interface Props {
    cut: PaymentCut;
    businessId?: number | null;
    justConfirmed?: boolean;
    onClose: () => void;
    onSent: (msg: string) => void;
}

function parseEmails(raw: string): string[] {
    return raw
        .split(/[\s,;]+/)
        .map(e => e.trim().toLowerCase())
        .filter(Boolean);
}

export function CutEmailModal({ cut, businessId, justConfirmed, onClose, onSent }: Props) {
    const [raw, setRaw] = useState('');
    const [sending, setSending] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const emails = parseEmails(raw);
    const invalid = emails.filter(e => !/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(e));

    const send = async () => {
        if (invalid.length > 0) {
            setError(`Correo invalido: ${invalid[0]}`);
            return;
        }
        setSending(true);
        setError(null);
        const res = await sendCodCutEmailAction(cut.id, emails, businessId || undefined);
        setSending(false);
        if (res.success) {
            const to = res.data?.sent_to?.join(', ');
            onSent(to ? `Corte enviado a ${to}` : 'Corte enviado por correo');
            onClose();
        } else {
            setError(res.message || 'Error al enviar el correo');
        }
    };

    return (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl max-w-md w-full p-5">
                <div className="flex items-start justify-between mb-1">
                    <h3 className="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
                        <Mail size={18} className="text-purple-600" /> Enviar corte por correo
                    </h3>
                    <button onClick={onClose} disabled={sending} className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700">
                        <X size={16} />
                    </button>
                </div>
                {justConfirmed && (
                    <p className="text-xs text-emerald-700 bg-emerald-50 border border-emerald-200 rounded-md px-2.5 py-1.5 mb-2">
                        Corte confirmado. Si quieres, envia ahora la relacion de lo consignado.
                    </p>
                )}
                <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                    Se enviara la relacion completa del corte del periodo{' '}
                    <strong>{formatDateOnly(cut.period_start)} - {formatDateOnly(cut.period_end)}</strong>:
                    fecha de confirmacion, resumen por transportadora y cada orden con su guia, fecha de entrega y valor consignado
                    ({cut.orders_count} ordenes, {formatMoney(cut.total_collected)}).
                </p>
                <label className="block text-xs font-semibold text-gray-500 dark:text-gray-400 mb-1">Destinatarios</label>
                <textarea
                    value={raw}
                    onChange={e => { setRaw(e.target.value); setError(null); }}
                    rows={3}
                    placeholder="correo@ejemplo.com, otro@ejemplo.com"
                    className="w-full px-3 py-2 text-sm rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                />
                <p className="text-[11px] text-gray-400 mt-1 mb-3">
                    Separa varios correos con coma. Si lo dejas vacio se envia a tu propio correo.
                </p>
                {error && (
                    <div className="bg-red-50 border border-red-200 rounded-md px-3 py-2 text-sm text-red-700 mb-3">{error}</div>
                )}
                <div className="flex justify-end gap-2">
                    <button
                        onClick={onClose}
                        disabled={sending}
                        className="px-3 py-2 text-sm rounded-md border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                    >
                        {justConfirmed ? 'Ahora no' : 'Cancelar'}
                    </button>
                    <button
                        onClick={send}
                        disabled={sending || invalid.length > 0}
                        className="px-3 py-2 text-sm rounded-md bg-purple-600 hover:bg-purple-700 text-white font-semibold inline-flex items-center gap-1.5 disabled:opacity-50"
                    >
                        {sending ? <RefreshCw size={14} className="animate-spin" /> : <Send size={14} />}
                        Enviar
                    </button>
                </div>
            </div>
        </div>
    );
}
