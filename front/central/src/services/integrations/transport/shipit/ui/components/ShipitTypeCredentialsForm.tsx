'use client';

import { KeyIcon, InformationCircleIcon } from '@heroicons/react/24/outline';
import { SecretInput } from '@/shared/ui';

export interface ShipitPlatformCredentials {
    email: string;
    access_token: string;
}

interface ShipitTypeCredentialsFormProps {
    credentials: ShipitPlatformCredentials;
    onChange: (credentials: ShipitPlatformCredentials) => void;
    isEditing?: boolean;
}

const fieldLabel = 'block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1';
const fieldHint = 'text-[11px] text-gray-400 dark:text-gray-500 mt-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)] font-mono';
const INPUT_BORDER = '#e9e9f0';

export default function ShipitTypeCredentialsForm({
    credentials,
    onChange,
    isEditing = false,
}: ShipitTypeCredentialsFormProps) {
    const set = (patch: Partial<ShipitPlatformCredentials>) => onChange({ ...credentials, ...patch });

    return (
        <div className="space-y-4">
            <div className="flex items-center gap-2">
                <KeyIcon className="w-5 h-5 text-[var(--color-primary)]" />
                <h3 className="text-sm font-bold text-gray-900 dark:text-white">Cuenta Shipit de la Plataforma</h3>
            </div>

            <p className="text-xs text-gray-600 dark:text-gray-300 flex items-start gap-1.5">
                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                <span>
                    Credenciales de la cuenta Shipit global (app.shipit.cl &gt; Configuracion &gt; API).
                    Todos los negocios generan envios con esta misma cuenta. Se guardan encriptadas (AES-256-GCM).
                </span>
            </p>

            <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                <div>
                    <label className={fieldLabel}>Email de la cuenta</label>
                    <input
                        type="email"
                        value={credentials.email}
                        onChange={(e) => set({ email: e.target.value })}
                        placeholder="cuenta@dominio.cl"
                        className={inputCls}
                        style={{ borderColor: INPUT_BORDER }}
                    />
                    <p className={fieldHint}>Header X-Shipit-Email</p>
                </div>
                <div>
                    <label className={fieldLabel}>Access Token</label>
                    <SecretInput
                        value={credentials.access_token}
                        onChange={(e) => set({ access_token: e.target.value })}
                        placeholder={isEditing ? 'Dejar vacio para mantener actual' : 'Access Token de la API'}
                        className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                    />
                    <p className={fieldHint}>Header X-Shipit-Access-Token</p>
                </div>
            </div>
        </div>
    );
}
