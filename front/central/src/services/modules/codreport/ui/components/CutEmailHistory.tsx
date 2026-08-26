'use client';

import { useEffect, useState } from 'react';
import { Mail, RefreshCw, CheckCircle2, XCircle } from 'lucide-react';
import { getCodCutEmailHistoryAction } from '../../infra/actions';
import { CutEmailLog } from '../../domain/types';
import { formatDateTime } from './helpers';

interface Props {
    cutId: number;
    businessId?: number | null;
    refreshKey?: number;
    compact?: boolean;
}

export function CutEmailHistory({ cutId, businessId, refreshKey, compact }: Props) {
    const [logs, setLogs] = useState<CutEmailLog[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        let cancelled = false;
        setLoading(true);
        getCodCutEmailHistoryAction(cutId, businessId || undefined).then(res => {
            if (cancelled) return;
            setLogs(res.success ? ((res.data || []) as CutEmailLog[]) : []);
            setLoading(false);
        });
        return () => { cancelled = true; };
    }, [cutId, businessId, refreshKey]);

    if (loading) {
        return (
            <div className="flex items-center gap-1.5 text-[11px] text-gray-400 py-1">
                <RefreshCw size={11} className="animate-spin" /> Cargando correos enviados...
            </div>
        );
    }

    if (logs.length === 0) {
        return (
            <div className="flex items-center gap-1.5 text-[11px] text-gray-400 py-1">
                <Mail size={11} /> Este corte aun no se ha enviado por correo.
            </div>
        );
    }

    return (
        <div className={compact ? '' : 'mb-3'}>
            <div className="flex items-center gap-1.5 text-[11px] font-semibold text-gray-500 dark:text-gray-400 mb-1">
                <Mail size={11} /> Correos enviados ({logs.length})
            </div>
            <ul className="space-y-0.5">
                {logs.map(l => (
                    <li key={l.id} className="flex items-center gap-2 text-[11px] text-gray-600 dark:text-gray-300 flex-wrap" title={l.status === 'failed' ? l.error_message : l.subject}>
                        {l.status === 'sent'
                            ? <CheckCircle2 size={12} className="text-emerald-600 shrink-0" />
                            : <XCircle size={12} className="text-red-600 shrink-0" />}
                        <span className="font-medium text-gray-800 dark:text-gray-100">{l.recipient}</span>
                        <span className="text-gray-400">{formatDateTime(l.sent_at)}</span>
                        {l.sent_by_name && <span className="text-gray-400">por {l.sent_by_name}</span>}
                        {l.status === 'failed' && <span className="text-red-600">fallo</span>}
                    </li>
                ))}
            </ul>
        </div>
    );
}
