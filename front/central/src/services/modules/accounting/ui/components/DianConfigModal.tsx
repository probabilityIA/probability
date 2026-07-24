'use client';

import { useEffect, useState } from 'react';
import { Modal } from '@/shared/ui';
import { DianConfig, UpdateDianConfigDTO } from '../../domain/types';
import { updateAccountingDianConfigAction } from '../../infra/actions';
import { useToast } from '@/shared/providers/toast-provider';

interface DianConfigModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSaved: () => void;
    config: DianConfig | null;
}

const inputClass = 'w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500';

export function DianConfigModal({ isOpen, onClose, onSaved, config }: DianConfigModalProps) {
    const { showToast } = useToast();
    const [apiUrl, setApiUrl] = useState('');
    const [clientId, setClientId] = useState('');
    const [clientSecret, setClientSecret] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [numberingRangeId, setNumberingRangeId] = useState('');
    const [municipalityId, setMunicipalityId] = useState('');
    const [identificationDocumentId, setIdentificationDocumentId] = useState('');
    const [legalOrganizationId, setLegalOrganizationId] = useState('');
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!isOpen) return;
        setApiUrl(config?.api_url || '');
        setClientId('');
        setClientSecret('');
        setUsername(config?.username || '');
        setPassword('');
        setNumberingRangeId(config?.numbering_range_id ? String(config.numbering_range_id) : '');
        setMunicipalityId(config?.municipality_id ? String(config.municipality_id) : '');
        setIdentificationDocumentId(config?.identification_document_id ? String(config.identification_document_id) : '');
        setLegalOrganizationId(config?.legal_organization_id ? String(config.legal_organization_id) : '');
        setError(null);
    }, [isOpen, config]);

    const optionalNumber = (value: string, label: string): { ok: true; value?: number } | { ok: false } => {
        if (!value.trim()) return { ok: true };
        const parsed = Number(value);
        if (isNaN(parsed) || parsed <= 0) {
            showToast(`${label} debe ser un numero valido`, 'error');
            return { ok: false };
        }
        return { ok: true, value: parsed };
    };

    const handleSubmit = async () => {
        if (!clientId.trim() || !clientSecret.trim() || !username.trim() || !password.trim()) {
            showToast('client_id, client_secret, usuario y contrasena son requeridos', 'error');
            return;
        }
        const numberingRange = optionalNumber(numberingRangeId, 'El rango de numeracion');
        if (!numberingRange.ok) return;
        const municipality = optionalNumber(municipalityId, 'El municipio');
        if (!municipality.ok) return;
        const identificationDocument = optionalNumber(identificationDocumentId, 'El tipo de documento');
        if (!identificationDocument.ok) return;
        const legalOrganization = optionalNumber(legalOrganizationId, 'El tipo de organizacion');
        if (!legalOrganization.ok) return;

        const payload: UpdateDianConfigDTO = {
            api_url: apiUrl.trim() || undefined,
            client_id: clientId.trim(),
            client_secret: clientSecret.trim(),
            username: username.trim(),
            password: password,
            numbering_range_id: numberingRange.value,
            municipality_id: municipality.value,
            identification_document_id: identificationDocument.value,
            legal_organization_id: legalOrganization.value,
        };

        setSaving(true);
        setError(null);
        const result = await updateAccountingDianConfigAction(payload);
        setSaving(false);
        if (result.success) {
            showToast('Configuracion DIAN guardada', 'success');
            onSaved();
        } else {
            setError(result.error);
        }
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title={config?.configured ? 'Editar configuracion DIAN' : 'Configurar facturacion electronica'}
            size="lg"
        >
            <div className="space-y-4">
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Credenciales de Factus para emitir facturas electronicas ante la DIAN. Se validan contra Factus antes de guardarse.
                </p>

                {error && (
                    <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 text-sm rounded-lg px-4 py-3">
                        {error}
                    </div>
                )}

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">URL del API</label>
                    <input
                        type="text"
                        autoComplete="off"
                        value={apiUrl}
                        onChange={(e) => setApiUrl(e.target.value)}
                        placeholder="https://api-sandbox.factus.com.co"
                        className={inputClass}
                    />
                    <p className="mt-1 text-xs text-gray-400 dark:text-gray-500">Vacio = produccion</p>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Client ID *</label>
                        <input
                            type="text"
                            autoComplete="off"
                            data-1p-ignore="true"
                            value={clientId}
                            onChange={(e) => setClientId(e.target.value)}
                            className={inputClass}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Client Secret *</label>
                        <input
                            type="password"
                            autoComplete="new-password"
                            data-1p-ignore="true"
                            value={clientSecret}
                            onChange={(e) => setClientSecret(e.target.value)}
                            className={inputClass}
                        />
                    </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Usuario *</label>
                        <input
                            type="text"
                            autoComplete="off"
                            data-1p-ignore="true"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            placeholder="correo@empresa.com"
                            className={inputClass}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Contrasena *</label>
                        <input
                            type="password"
                            autoComplete="new-password"
                            data-1p-ignore="true"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            className={inputClass}
                        />
                    </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Rango de numeracion (opcional)</label>
                        <input
                            type="number"
                            min="1"
                            step="1"
                            value={numberingRangeId}
                            onChange={(e) => setNumberingRangeId(e.target.value)}
                            className={inputClass}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Municipio (opcional)</label>
                        <input
                            type="number"
                            min="1"
                            step="1"
                            value={municipalityId}
                            onChange={(e) => setMunicipalityId(e.target.value)}
                            className={inputClass}
                        />
                    </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Tipo de documento (opcional)</label>
                        <input
                            type="number"
                            min="1"
                            step="1"
                            value={identificationDocumentId}
                            onChange={(e) => setIdentificationDocumentId(e.target.value)}
                            placeholder="3"
                            className={inputClass}
                        />
                        <p className="mt-1 text-xs text-gray-400 dark:text-gray-500">Default 3 = NIT segun Factus</p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Tipo de organizacion (opcional)</label>
                        <input
                            type="number"
                            min="1"
                            step="1"
                            value={legalOrganizationId}
                            onChange={(e) => setLegalOrganizationId(e.target.value)}
                            placeholder="2"
                            className={inputClass}
                        />
                        <p className="mt-1 text-xs text-gray-400 dark:text-gray-500">Default 2 = persona natural, 1 = juridica</p>
                    </div>
                </div>

                <div className="flex justify-end gap-3 pt-2">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 text-sm font-medium rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                    >
                        Cancelar
                    </button>
                    <button
                        onClick={handleSubmit}
                        disabled={saving}
                        className="px-4 py-2 text-sm font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-colors disabled:opacity-50"
                    >
                        {saving ? 'Validando...' : 'Guardar'}
                    </button>
                </div>
            </div>
        </Modal>
    );
}
