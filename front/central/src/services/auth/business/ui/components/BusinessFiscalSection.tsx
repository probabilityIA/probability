'use client';

import React, { useState } from 'react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Select } from '@/shared/ui/select';
import { Alert } from '@/shared/ui/alert';
import { Spinner } from '@/shared/ui/spinner';
import { Checkbox } from '@/shared/ui/checkbox';
import { getBusinessFiscalProfileAction, updateBusinessFiscalProfileAction } from '../../infra/actions';
import { getActionError } from '@/shared/utils/action-result';

interface FiscalFormState {
    document_type: string;
    document_number: string;
    dv: string;
    person_type: string;
    tax_regime: string;
    municipality_id: string;
    address: string;
    phone: string;
    billing_email: string;
    applies_iva: boolean;
    iva_rate: string;
    applies_retefuente: boolean;
    retefuente_rate: string;
    applies_reteica: boolean;
    reteica_rate: string;
    notes: string;
}

const EMPTY_FORM: FiscalFormState = {
    document_type: 'NIT',
    document_number: '',
    dv: '',
    person_type: 'JURÍDICA',
    tax_regime: '',
    municipality_id: '',
    address: '',
    phone: '',
    billing_email: '',
    applies_iva: false,
    iva_rate: '19',
    applies_retefuente: false,
    retefuente_rate: '11',
    applies_reteica: false,
    reteica_rate: '0.966',
    notes: '',
};

interface BusinessFiscalSectionProps {
    businessId: number;
}

interface TaxRowProps {
    label: string;
    checked: boolean;
    rate: string;
    onCheckedChange: (checked: boolean) => void;
    onRateChange: (rate: string) => void;
}

const TaxRow: React.FC<TaxRowProps> = ({ label, checked, rate, onCheckedChange, onRateChange }) => (
    <div className="flex items-center justify-between gap-4 p-3 border border-gray-200 dark:border-gray-700 rounded-lg">
        <label className="flex items-center gap-2 cursor-pointer text-sm text-gray-700 dark:text-gray-200">
            <Checkbox checked={checked} onCheckedChange={onCheckedChange} />
            {label}
        </label>
        <div className="flex items-center gap-2">
            <input
                type="number"
                min={0}
                max={100}
                step="0.001"
                value={rate}
                disabled={!checked}
                onChange={(e) => onRateChange(e.target.value)}
                aria-label={`Tarifa de ${label}`}
                className="input w-24 text-right disabled:opacity-50 disabled:cursor-not-allowed"
            />
            <span className="text-sm text-gray-500 dark:text-gray-400">%</span>
        </div>
    </div>
);

export const BusinessFiscalSection: React.FC<BusinessFiscalSectionProps> = ({ businessId }) => {
    const [isOpen, setIsOpen] = useState(false);
    const [loaded, setLoaded] = useState(false);
    const [loading, setLoading] = useState(false);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [form, setForm] = useState<FiscalFormState>(EMPTY_FORM);

    const handleChange = <K extends keyof FiscalFormState>(field: K, value: FiscalFormState[K]) => {
        setSuccess(null);
        setForm(prev => ({ ...prev, [field]: value }));
    };

    const loadProfile = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getBusinessFiscalProfileAction(businessId);
            const data = response.data;
            if (data) {
                setForm({
                    document_type: data.document_type || 'NIT',
                    document_number: data.document_number || '',
                    dv: data.dv || '',
                    person_type: data.person_type || 'JURÍDICA',
                    tax_regime: data.tax_regime || '',
                    municipality_id: data.municipality_id || '',
                    address: data.address || '',
                    phone: data.phone || '',
                    billing_email: data.billing_email || '',
                    applies_iva: data.applies_iva,
                    iva_rate: String(data.iva_rate ?? 19),
                    applies_retefuente: data.applies_retefuente,
                    retefuente_rate: String(data.retefuente_rate ?? 11),
                    applies_reteica: data.applies_reteica,
                    reteica_rate: String(data.reteica_rate ?? 0.966),
                    notes: data.notes || '',
                });
            }
            setLoaded(true);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar los datos fiscales'));
        } finally {
            setLoading(false);
        }
    };

    const handleToggle = () => {
        const next = !isOpen;
        setIsOpen(next);
        if (next && !loaded && !loading) {
            loadProfile();
        }
    };

    const handleSave = async () => {
        setSaving(true);
        setError(null);
        setSuccess(null);
        try {
            await updateBusinessFiscalProfileAction(businessId, {
                document_type: form.document_type,
                document_number: form.document_number,
                dv: form.dv,
                person_type: form.person_type,
                tax_regime: form.tax_regime,
                municipality_id: form.municipality_id,
                address: form.address,
                phone: form.phone,
                billing_email: form.billing_email,
                applies_iva: form.applies_iva,
                iva_rate: parseFloat(form.iva_rate) || 0,
                applies_retefuente: form.applies_retefuente,
                retefuente_rate: parseFloat(form.retefuente_rate) || 0,
                applies_reteica: form.applies_reteica,
                reteica_rate: parseFloat(form.reteica_rate) || 0,
                notes: form.notes,
            });
            setSuccess('Datos fiscales guardados correctamente');
        } catch (err: any) {
            setError(getActionError(err, 'Error al guardar los datos fiscales'));
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
            <button
                type="button"
                onClick={handleToggle}
                className="w-full flex items-center justify-between px-4 py-3 bg-gray-50 dark:bg-gray-700/50 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            >
                <span className="text-sm font-semibold text-gray-900 dark:text-white uppercase tracking-wide">
                    Datos fiscales y facturacion
                </span>
                <svg
                    className={`w-5 h-5 text-gray-500 dark:text-gray-400 transition-transform ${isOpen ? 'rotate-180' : ''}`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
            </button>

            {isOpen && (
                <div className="p-4 space-y-4">
                    {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}
                    {success && <Alert type="success" onClose={() => setSuccess(null)}>{success}</Alert>}

                    {loading ? (
                        <div className="flex justify-center py-8">
                            <Spinner />
                        </div>
                    ) : (
                        <>
                            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                                <Select
                                    label="Tipo de documento"
                                    value={form.document_type}
                                    onChange={(e) => handleChange('document_type', e.target.value)}
                                    options={[
                                        { value: 'NIT', label: 'NIT' },
                                        { value: 'CC', label: 'Cédula de ciudadanía' },
                                    ]}
                                />
                                <Input
                                    label="Número de documento"
                                    value={form.document_number}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('document_number', e.target.value)}
                                />
                                <Input
                                    label="DV"
                                    value={form.dv}
                                    maxLength={2}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('dv', e.target.value)}
                                    helperText="Dígito de verificación"
                                />
                            </div>

                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                                <Select
                                    label="Tipo de persona"
                                    value={form.person_type}
                                    onChange={(e) => handleChange('person_type', e.target.value)}
                                    options={[
                                        { value: 'JURÍDICA', label: 'Jurídica' },
                                        { value: 'NATURAL', label: 'Natural' },
                                    ]}
                                />
                                <Input
                                    label="Régimen tributario"
                                    value={form.tax_regime}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('tax_regime', e.target.value)}
                                />
                            </div>

                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                                <Input
                                    label="Municipio"
                                    value={form.municipality_id}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('municipality_id', e.target.value)}
                                    helperText="Código DIAN/Factus, ej 980 = Bogotá"
                                />
                                <Input
                                    label="Dirección"
                                    value={form.address}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('address', e.target.value)}
                                />
                            </div>

                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                                <Input
                                    label="Teléfono"
                                    value={form.phone}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('phone', e.target.value)}
                                />
                                <Input
                                    label="Correo de facturación"
                                    type="email"
                                    value={form.billing_email}
                                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleChange('billing_email', e.target.value)}
                                />
                            </div>

                            <div className="space-y-2">
                                <span className="block text-sm font-medium text-gray-700 dark:text-gray-200">
                                    Impuestos del cliente
                                </span>
                                <p className="text-xs text-gray-500 dark:text-gray-400">
                                    Define que impuestos y retenciones se aplican al facturarle a este negocio.
                                </p>
                                <div className="space-y-2">
                                    <TaxRow
                                        label="Responsable de IVA"
                                        checked={form.applies_iva}
                                        rate={form.iva_rate}
                                        onCheckedChange={(checked) => handleChange('applies_iva', checked)}
                                        onRateChange={(rate) => handleChange('iva_rate', rate)}
                                    />
                                    <TaxRow
                                        label="Practica retención en la fuente"
                                        checked={form.applies_retefuente}
                                        rate={form.retefuente_rate}
                                        onCheckedChange={(checked) => handleChange('applies_retefuente', checked)}
                                        onRateChange={(rate) => handleChange('retefuente_rate', rate)}
                                    />
                                    <TaxRow
                                        label="Practica ReteICA"
                                        checked={form.applies_reteica}
                                        rate={form.reteica_rate}
                                        onCheckedChange={(checked) => handleChange('applies_reteica', checked)}
                                        onRateChange={(rate) => handleChange('reteica_rate', rate)}
                                    />
                                </div>
                            </div>

                            <div className="space-y-2">
                                <label
                                    htmlFor={`fiscal-notes-${businessId}`}
                                    className="block text-sm font-medium text-gray-700 dark:text-gray-200"
                                >
                                    Notas
                                </label>
                                <textarea
                                    id={`fiscal-notes-${businessId}`}
                                    value={form.notes}
                                    onChange={(e) => handleChange('notes', e.target.value)}
                                    rows={3}
                                    className="input resize-y"
                                />
                            </div>

                            <div className="flex justify-end pt-2 border-t border-gray-200 dark:border-gray-700">
                                <Button type="button" onClick={handleSave} disabled={saving || loading}>
                                    {saving ? <Spinner size="sm" /> : 'Guardar datos fiscales'}
                                </Button>
                            </div>
                        </>
                    )}
                </div>
            )}
        </div>
    );
};
