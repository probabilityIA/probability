'use client';

import { useEffect, useState } from 'react';
import { PlayIcon, StopIcon, ArrowPathIcon } from '@heroicons/react/24/outline';
import { TokenStorage } from '@/shared/utils/token-storage';
import { getWooStoreStatusAction, startWooStoreAction, stopWooStoreAction } from '@/services/woostore/infra/actions';
import { WooStoreState } from '@/services/woostore/domain/types';

const STATE_LABEL: Record<string, string> = {
    running: 'Encendida',
    stopped: 'Apagada',
    stopping: 'Apagando...',
    pending: 'Encendiendo...',
    'shutting-down': 'Apagando...',
};

export function WooStorePowerWidget() {
    const [isSuper, setIsSuper] = useState(false);
    const [st, setSt] = useState<WooStoreState | null>(null);
    const [loading, setLoading] = useState(false);

    const refresh = async () => {
        const res = await getWooStoreStatusAction();
        setSt(res);
    };

    useEffect(() => {
        const perms = TokenStorage.getPermissions();
        if (!perms?.is_super) return;
        setIsSuper(true);
        refresh();
    }, []);

    if (!isSuper) return null;

    const state = st?.state || 'desconocido';
    const isRunning = state === 'running';
    const isBusy = state === 'pending' || state === 'stopping' || state === 'shutting-down';

    const doAction = async (action: 'start' | 'stop') => {
        setLoading(true);
        try {
            const res = action === 'start' ? await startWooStoreAction() : await stopWooStoreAction();
            setSt(res);
            setTimeout(refresh, 4000);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="rounded-xl p-4 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between"
            style={{ backgroundColor: '#ffffff', border: '1px solid #eceaf3' }}>
            <div className="flex items-center gap-3">
                <span className="w-2.5 h-2.5 rounded-full" style={{ backgroundColor: isRunning ? '#16a34a' : isBusy ? '#f59e0b' : '#9ca3af' }} />
                <div>
                    <p className="text-[13px] font-bold text-gray-900 dark:text-gray-100">Tienda de pruebas WooCommerce</p>
                    <p className="text-[11px] text-gray-500">
                        Estado: {STATE_LABEL[state] || state}
                        {st?.store_url && isRunning && (
                            <> · <a href={st.store_url} target="_blank" rel="noreferrer" className="underline">abrir</a></>
                        )}
                        {st?.error && <span className="text-red-600"> · {st.error}</span>}
                    </p>
                </div>
            </div>
            <div className="flex items-center gap-2">
                <button type="button" onClick={refresh} disabled={loading}
                    className="p-2 rounded-lg border text-gray-500 disabled:opacity-50" style={{ borderColor: '#e9e9f0' }} title="Actualizar">
                    <ArrowPathIcon className="w-4 h-4" />
                </button>
                {isRunning ? (
                    <button type="button" onClick={() => doAction('stop')} disabled={loading || isBusy}
                        className="px-3 py-1.5 text-[12px] font-semibold rounded-lg flex items-center gap-1.5 disabled:opacity-50"
                        style={{ backgroundColor: '#ffffff', color: '#b42318', border: '1px solid #f3c9c9' }}>
                        <StopIcon className="w-4 h-4" /> Apagar
                    </button>
                ) : (
                    <button type="button" onClick={() => doAction('start')} disabled={loading || isBusy}
                        className="px-3 py-1.5 text-[12px] font-semibold rounded-lg text-white flex items-center gap-1.5 disabled:opacity-50"
                        style={{ backgroundColor: 'var(--color-primary)' }}>
                        <PlayIcon className="w-4 h-4" /> Encender
                    </button>
                )}
            </div>
        </div>
    );
}
