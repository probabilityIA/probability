'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { UserCircleIcon } from '@heroicons/react/24/outline';

import {
    getWhatsAppBusinessProfileAction,
    updateWhatsAppBusinessProfileAction,
    updateWhatsAppBusinessProfilePhotoAction,
} from '../../infra/actions';
import { WhatsAppBusinessProfile } from '../../domain/types';
import { ACCENT, ActionButton, Card, fieldHint, fieldLabel, inputCls, INPUT_BORDER } from './ui-kit';

const MAX_ABOUT = 139;
const MAX_DESCRIPTION = 512;
const MAX_FOTO_MB = 5;

const CATEGORIAS: { value: string; label: string }[] = [
    { value: '', label: 'Sin categoría' },
    { value: 'RETAIL', label: 'Tienda / comercio' },
    { value: 'APPAREL', label: 'Ropa y accesorios' },
    { value: 'BEAUTY', label: 'Belleza y cuidado personal' },
    { value: 'GROCERY', label: 'Supermercado y alimentos' },
    { value: 'RESTAURANT', label: 'Restaurante' },
    { value: 'HEALTH', label: 'Salud' },
    { value: 'EDU', label: 'Educación' },
    { value: 'PROF_SERVICES', label: 'Servicios profesionales' },
    { value: 'AUTO', label: 'Automóvil' },
    { value: 'TRAVEL', label: 'Viajes' },
    { value: 'ENTERTAIN', label: 'Entretenimiento' },
    { value: 'FINANCE', label: 'Finanzas' },
    { value: 'HOTEL', label: 'Hotel' },
    { value: 'EVENT_PLAN', label: 'Eventos' },
    { value: 'NONPROFIT', label: 'Sin ánimo de lucro' },
    { value: 'OTHER', label: 'Otro' },
];

interface WhatsAppProfilePanelProps {
    businessId?: number;
    enabled: boolean;
}

interface FormState {
    about: string;
    description: string;
    email: string;
    address: string;
    vertical: string;
    website: string;
}

const vacio: FormState = {
    about: '',
    description: '',
    email: '',
    address: '',
    vertical: '',
    website: '',
};

function aFormulario(perfil: WhatsAppBusinessProfile): FormState {
    return {
        about: perfil.about || '',
        description: perfil.description || '',
        email: perfil.email || '',
        address: perfil.address || '',
        vertical: perfil.vertical === 'UNDEFINED' ? '' : perfil.vertical || '',
        website: (perfil.websites && perfil.websites[0]) || '',
    };
}

export function WhatsAppProfilePanel({ businessId, enabled }: WhatsAppProfilePanelProps) {
    const [form, setForm] = useState<FormState>(vacio);
    const [original, setOriginal] = useState<FormState>(vacio);
    const [fotoURL, setFotoURL] = useState('');
    const [loading, setLoading] = useState(false);
    const [saving, setSaving] = useState(false);
    const [subiendoFoto, setSubiendoFoto] = useState(false);
    const [error, setError] = useState('');
    const [aviso, setAviso] = useState('');
    const inputFoto = useRef<HTMLInputElement>(null);

    const aplicar = useCallback((perfil: WhatsAppBusinessProfile) => {
        const datos = aFormulario(perfil);
        setForm(datos);
        setOriginal(datos);
        setFotoURL(perfil.profile_picture_url || '');
    }, []);

    const cargar = useCallback(async () => {
        if (!enabled) return;
        setLoading(true);
        setError('');
        const res = await getWhatsAppBusinessProfileAction(businessId);
        if (res.success && res.data) {
            aplicar(res.data);
        } else if (!res.success) {
            setError(res.message || 'No se pudo consultar el perfil');
        }
        setLoading(false);
    }, [businessId, enabled, aplicar]);

    useEffect(() => {
        cargar();
    }, [cargar]);

    if (!enabled) return null;

    const cambiado = JSON.stringify(form) !== JSON.stringify(original);

    const guardar = async () => {
        setSaving(true);
        setError('');
        setAviso('');

        const res = await updateWhatsAppBusinessProfileAction(
            {
                about: form.about,
                description: form.description,
                email: form.email,
                address: form.address,
                vertical: form.vertical,
                websites: form.website.trim() ? [form.website.trim()] : [],
            },
            businessId
        );

        if (res.success && res.data) {
            aplicar(res.data);
            setAviso('Perfil actualizado. Tus clientes lo verán al abrir el chat.');
        } else {
            setError(res.message || 'No se pudo guardar el perfil');
        }
        setSaving(false);
    };

    const cambiarFoto = async (file: File) => {
        if (file.size > MAX_FOTO_MB * 1024 * 1024) {
            setError(`La imagen no puede pesar más de ${MAX_FOTO_MB} MB`);
            return;
        }

        setSubiendoFoto(true);
        setError('');
        setAviso('');

        const res = await updateWhatsAppBusinessProfilePhotoAction(file, businessId);
        if (res.success && res.data) {
            aplicar(res.data);
            setAviso('Foto actualizada.');
        } else {
            setError(res.message || 'No se pudo actualizar la foto');
        }

        setSubiendoFoto(false);
        if (inputFoto.current) inputFoto.current.value = '';
    };

    return (
        <Card
            icon={<UserCircleIcon style={{ color: ACCENT, width: 16, height: 16 }} />}
            title={'Perfil de tu WhatsApp'}
            description={
                'Es lo que ve el cliente al abrir el chat o tocar tu nombre. El nombre del negocio no se edita aquí: lo aprueba Meta al conectar el número.'
            }
        >
            <div className="space-y-4">
                <div className="flex items-center gap-4">
                    <div
                        className="flex h-16 w-16 shrink-0 items-center justify-center overflow-hidden rounded-full bg-gray-100 dark:bg-gray-700"
                        style={{ border: `1px solid ${INPUT_BORDER}` }}
                    >
                        {fotoURL ? (
                            <img src={fotoURL} alt="Foto del perfil" className="h-full w-full object-cover" />
                        ) : (
                            <UserCircleIcon className="h-9 w-9 text-gray-300" />
                        )}
                    </div>
                    <div className="min-w-0">
                        <input
                            ref={inputFoto}
                            type="file"
                            accept="image/png,image/jpeg"
                            className="hidden"
                            onChange={(e) => {
                                const file = e.target.files?.[0];
                                if (file) cambiarFoto(file);
                            }}
                        />
                        <ActionButton
                            variant="ghost"
                            onClick={() => inputFoto.current?.click()}
                            loading={subiendoFoto}
                            disabled={loading}
                        >
                            {fotoURL ? 'Cambiar foto' : 'Subir foto'}
                        </ActionButton>
                        <p className={fieldHint}>{'Cuadrada, mínimo 192x192, hasta 5 MB. JPG o PNG.'}</p>
                    </div>
                </div>

                <div>
                    <label className={fieldLabel}>{'Información'}</label>
                    <input
                        type="text"
                        value={form.about}
                        maxLength={MAX_ABOUT}
                        onChange={(e) => setForm({ ...form, about: e.target.value })}
                        placeholder={'Envíos a todo el país'}
                        className={inputCls}
                        style={{ borderColor: INPUT_BORDER }}
                    />
                    <p className={fieldHint}>
                        {'Frase corta bajo tu nombre. '}
                        {form.about.length}/{MAX_ABOUT}
                    </p>
                </div>

                <div>
                    <label className={fieldLabel}>{'Descripción'}</label>
                    <textarea
                        value={form.description}
                        maxLength={MAX_DESCRIPTION}
                        rows={3}
                        onChange={(e) => setForm({ ...form, description: e.target.value })}
                        placeholder={'Qué vendes y cómo puedes ayudar a tu cliente'}
                        className={inputCls}
                        style={{ borderColor: INPUT_BORDER }}
                    />
                    <p className={fieldHint}>
                        {form.description.length}/{MAX_DESCRIPTION}
                    </p>
                </div>

                <div className="grid gap-3 sm:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>{'Correo'}</label>
                        <input
                            type="email"
                            value={form.email}
                            onChange={(e) => setForm({ ...form, email: e.target.value })}
                            placeholder="contacto@tunegocio.com"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                    </div>
                    <div>
                        <label className={fieldLabel}>{'Sitio web'}</label>
                        <input
                            type="url"
                            value={form.website}
                            onChange={(e) => setForm({ ...form, website: e.target.value })}
                            placeholder="https://tunegocio.com"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>{'Debe empezar por https://'}</p>
                    </div>
                </div>

                <div className="grid gap-3 sm:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>{'Dirección'}</label>
                        <input
                            type="text"
                            value={form.address}
                            onChange={(e) => setForm({ ...form, address: e.target.value })}
                            placeholder={'Calle 100 #15-20, Bogotá'}
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                    </div>
                    <div>
                        <label className={fieldLabel}>{'Categoría'}</label>
                        <select
                            value={form.vertical}
                            onChange={(e) => setForm({ ...form, vertical: e.target.value })}
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        >
                            {CATEGORIAS.map((c) => (
                                <option key={c.value} value={c.value}>
                                    {c.label}
                                </option>
                            ))}
                        </select>
                    </div>
                </div>

                {error && (
                    <p className="rounded-lg bg-red-50 dark:bg-red-900/20 px-3 py-2 text-[12px] text-red-700 dark:text-red-300">
                        {error}
                    </p>
                )}
                {aviso && !error && (
                    <p className="rounded-lg bg-emerald-50 dark:bg-emerald-900/20 px-3 py-2 text-[12px] text-emerald-700 dark:text-emerald-300">
                        {aviso}
                    </p>
                )}

                <div className="flex items-center gap-2">
                    <ActionButton onClick={guardar} disabled={!cambiado || loading} loading={saving}>
                        Guardar perfil
                    </ActionButton>
                    {cambiado && <span className={fieldHint}>Tienes cambios sin guardar.</span>}
                </div>
            </div>
        </Card>
    );
}
