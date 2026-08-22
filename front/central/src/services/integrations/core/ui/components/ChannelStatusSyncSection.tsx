'use client';

export interface ChannelStatusSyncConfig {
    inbound: boolean;
    outbound: boolean;
}

interface ChannelStatusSyncSectionProps {
    channelName: string;
    value: ChannelStatusSyncConfig;
    onChange: (v: ChannelStatusSyncConfig) => void;
    accentColor?: string;
}

const DEFAULT_ACCENT = 'var(--color-primary)';
const BORDER = '#e9e9f0';

export function readChannelStatusSyncConfig(config: any): ChannelStatusSyncConfig {
    return {
        inbound: config?.status_inbound_enabled !== false,
        outbound: config?.status_sync_enabled !== false,
    };
}

export function writeChannelStatusSyncConfig(value: ChannelStatusSyncConfig) {
    return {
        status_inbound_enabled: value.inbound,
        status_sync_enabled: value.outbound,
    };
}

function Switch({ checked, onClick, accent }: { checked: boolean; onClick: () => void; accent: string }) {
    return (
        <button
            type="button"
            role="switch"
            aria-checked={checked}
            onClick={onClick}
            className={`relative inline-flex h-6 w-11 flex-shrink-0 items-center rounded-full transition-colors ${checked ? '' : 'bg-gray-300 dark:bg-gray-600'}`}
            style={checked ? { backgroundColor: accent } : undefined}
        >
            <span className={`inline-block h-5 w-5 transform rounded-full bg-white shadow transition-transform ${checked ? 'translate-x-5' : 'translate-x-0.5'}`} />
        </button>
    );
}

export function ChannelStatusSyncSection({ channelName, value, onChange, accentColor }: ChannelStatusSyncSectionProps) {
    const accent = accentColor || DEFAULT_ACCENT;

    return (
        <div className="space-y-2">
            <div
                className="flex items-start justify-between gap-3 rounded-lg bg-white dark:bg-gray-800 px-3 py-2.5"
                style={{ border: `1px solid ${BORDER}` }}
            >
                <div>
                    <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100">
                        Recibir estados desde {channelName}
                    </h4>
                    <p className="mt-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                        Cuando el pedido cambia de estado en el canal, la orden se actualiza aca. Apagalo si la
                        operacion se lleva en Probability y no quieres que el canal la mueva. Las cancelaciones y los
                        reembolsos entran igual, para no despachar algo que el cliente ya cancelo.
                    </p>
                </div>
                <Switch checked={value.inbound} accent={accent} onClick={() => onChange({ ...value, inbound: !value.inbound })} />
            </div>

            <div
                className="flex items-start justify-between gap-3 rounded-lg bg-white dark:bg-gray-800 px-3 py-2.5"
                style={{ border: `1px solid ${BORDER}` }}
            >
                <div>
                    <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100">
                        Enviar estados y guia a {channelName}
                    </h4>
                    <p className="mt-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                        Cuando la orden avanza en Probability se actualiza el pedido en el canal, con el numero de guia
                        y el seguimiento del envio. Apagalo si el canal no debe recibir nada desde aca.
                    </p>
                </div>
                <Switch checked={value.outbound} accent={accent} onClick={() => onChange({ ...value, outbound: !value.outbound })} />
            </div>
        </div>
    );
}
