'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Button, Alert, Select } from '@/shared/ui';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import { getActionError } from '@/shared/utils/action-result';
import { CreditCardIcon, ArrowLeftIcon, BeakerIcon, CheckCircleIcon } from '@heroicons/react/24/outline';

interface BoldConfigFormProps {
    integrationTypeId: number;
    onSuccess?: () => void;
    onCancel?: () => void;
    integrationTypeBaseURLTest?: string;
}

const BOLD_DEFAULT_NAME = 'Bold';

interface SwitchProps {
    checked: boolean;
    onChange: (v: boolean) => void;
    titleOn: string;
    titleOff: string;
    description?: string;
    icon?: React.ReactNode;
    color?: 'green' | 'amber';
}

function Switch({ checked, onChange, titleOn, titleOff, description, icon, color = 'green' }: SwitchProps) {
    const ringColor = color === 'green' ? 'bg-green-600' : 'bg-amber-500';
    const titleClass = color === 'green'
        ? checked ? 'text-green-700 dark:text-green-300' : 'text-gray-500 dark:text-gray-400'
        : checked ? 'text-amber-700 dark:text-amber-300' : 'text-gray-500 dark:text-gray-400';
    return (
        <div className="flex items-start gap-4 p-4 bg-gray-50 dark:bg-gray-700 rounded-xl border border-gray-200 dark:border-gray-600">
            <button
                type="button"
                role="switch"
                aria-checked={checked}
                onClick={() => onChange(!checked)}
                className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 mt-1 ${
                    checked ? ringColor : 'bg-gray-300 dark:bg-gray-500'
                }`}
            >
                <span
                    className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ${
                        checked ? 'translate-x-5' : 'translate-x-0'
                    }`}
                />
            </button>
            <div className="flex-1">
                <div className="flex items-center gap-2">
                    {icon}
                    <span className={`font-semibold ${titleClass}`}>
                        {checked ? titleOn : titleOff}
                    </span>
                </div>
                {description && (
                    <p className="text-sm text-gray-600 dark:text-gray-300 mt-0.5">{description}</p>
                )}
            </div>
        </div>
    );
}

export function BoldConfigForm({ integrationTypeId, onSuccess, onCancel }: BoldConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [isActive, setIsActive] = useState(true);
    const [isTesting, setIsTesting] = useState(false);

    useEffect(() => {
        const checkUserAndLoadBusinesses = async () => {
            const permissions = TokenStorage.getPermissions();
            const isSuperUser = permissions?.is_super || false;
            setIsSuperAdmin(isSuperUser);

            if (isSuperUser) {
                setLoadingBusinesses(true);
                try {
                    const response = await getBusinessesSimpleAction();
                    if (response.success && response.data) {
                        setBusinesses(response.data);
                    }
                } catch (err) {
                    console.error('Error loading businesses:', err);
                    showToast('Error al cargar la lista de negocios', 'error');
                } finally {
                    setLoadingBusinesses(false);
                }
            } else if (permissions?.business_id) {
                setSelectedBusinessId(permissions.business_id);
            }
        };

        checkUserAndLoadBusinesses();
    }, []);

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        if (!selectedBusinessId) {
            setError('Selecciona un negocio');
            return;
        }
        setLoading(true);
        setError(null);

        try {
            const response = await createIntegrationAction({
                name: BOLD_DEFAULT_NAME,
                code: `bold_${Date.now()}`,
                integration_type_id: integrationTypeId,
                category: 'pay',
                business_id: selectedBusinessId,
                config: { use_platform_token: true } as any,
                is_active: isActive,
                is_default: false,
                is_testing: isTesting,
            });

            if (response.success) {
                showToast('Bold guardado para el negocio', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al guardar integración');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al guardar integración'));
            showToast('Error al guardar integración Bold', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-6" autoComplete="off">
            <div className="border-b border-gray-200 dark:border-gray-700 pb-6">
                <div className="flex items-center gap-3 mb-2">
                    <div className="p-2 bg-blue-50 rounded-lg">
                        <CreditCardIcon className="w-6 h-6 text-blue-600" />
                    </div>
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Bold Pagos</h2>
                </div>
                <p className="text-sm text-gray-600 dark:text-gray-300 ml-14">
                    Configura Bold como pasarela de pagos para este negocio.
                </p>
            </div>

            {error && <Alert type="error">{error}</Alert>}

            {isSuperAdmin && (
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                        Negocio <span className="text-red-500">*</span>
                    </label>
                    {loadingBusinesses ? (
                        <div className="text-sm text-gray-600 dark:text-gray-300 p-3 bg-gray-100 dark:bg-gray-700 rounded-lg">
                            Cargando negocios...
                        </div>
                    ) : (
                        <Select
                            value={selectedBusinessId?.toString() || ''}
                            onChange={(e) => setSelectedBusinessId(Number(e.target.value))}
                            options={[
                                { value: '', label: '-- Selecciona un negocio --' },
                                ...businesses.map((b) => ({ value: b.id.toString(), label: b.name })),
                            ]}
                            required
                            className="bg-white dark:bg-gray-800"
                        />
                    )}
                </div>
            )}

            <Switch
                checked={isActive}
                onChange={setIsActive}
                titleOn="Activo"
                titleOff="Desactivado"
                description={isActive ? 'Bold procesara pagos para este negocio.' : 'Bold no procesara pagos para este negocio.'}
                icon={<CheckCircleIcon className="w-5 h-5 text-green-600" />}
                color="green"
            />

            <Switch
                checked={isTesting}
                onChange={setIsTesting}
                titleOn="Modo Sandbox activo"
                titleOff="Modo Sandbox inactivo"
                description={isTesting ? 'Las transacciones iran al ambiente de pruebas (sin cobros reales).' : 'Las transacciones iran al ambiente real.'}
                icon={<BeakerIcon className="w-5 h-5 text-amber-600" />}
                color="amber"
            />

            <div className="flex justify-between items-center gap-3 pt-6 border-t border-gray-200 dark:border-gray-700">
                {onCancel && (
                    <Button
                        type="button"
                        onClick={onCancel}
                        disabled={loading}
                        className="min-w-[140px] bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 text-gray-700 dark:text-gray-200 border border-gray-300 dark:border-gray-600"
                    >
                        <ArrowLeftIcon className="w-4 h-4 mr-2" />
                        Cancelar
                    </Button>
                )}
                <Button
                    type="submit"
                    variant="primary"
                    disabled={loading || !selectedBusinessId}
                    className="min-w-[160px] bg-blue-600 hover:bg-blue-700 text-white font-semibold"
                >
                    {loading ? 'Guardando...' : 'Guardar'}
                </Button>
            </div>
        </form>
    );
}
