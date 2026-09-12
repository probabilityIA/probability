'use client';

import { useCallback, useEffect, useState } from 'react';
import { Alert } from '@/shared/ui';
import { DevicePhoneMobileIcon } from '@heroicons/react/24/outline';
import { ACCENT, ActionButton, Card, Pill, Steps, fieldHint, fieldLabel, inputCls } from './ui-kit';
import {
    addWhatsAppNumberAction,
    getWhatsAppNumberStateAction,
    registerWhatsAppNumberAction,
    removeWhatsAppNumberAction,
    requestWhatsAppNumberCodeAction,
    verifyWhatsAppNumberCodeAction,
} from '../../infra/actions';
import { WhatsAppNumberState } from '../../domain/types';

interface WhatsAppNumberWizardProps {
    businessId?: number;
    onChanged?: () => void;
}

const PASOS = [
    { key: 'sin_numero', label: 'Tu número' },
    { key: 'esperando_codigo', label: 'Código' },
    { key: 'verificado', label: 'Activación' },
    { key: 'registrado', label: 'Listo' },
];

function indiceDePaso(status: string): number {
    const idx = PASOS.findIndex((paso) => paso.key === status);
    if (idx >= 0) return idx;
    return status === 'nombre_en_revision' ? 3 : 0;
}

export default function WhatsAppNumberWizard({ businessId, onChanged }: WhatsAppNumberWizardProps) {
    const [state, setState] = useState<WhatsAppNumberState | null>(null);
    const [countryCode, setCountryCode] = useState('57');
    const [phone, setPhone] = useState('');
    const [name, setName] = useState('');
    const [code, setCode] = useState('');
    const [pin, setPin] = useState('');
    const [busy, setBusy] = useState(false);
    const [confirmando, setConfirmando] = useState(false);
    const [confirmandoAlta, setConfirmandoAlta] = useState(false);
    const [loading, setLoading] = useState(true);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        const result = await getWhatsAppNumberStateAction(businessId);
        if (result.success && result.data) {
            setState(result.data);
        }
        setLoading(false);
    }, [businessId]);

    useEffect(() => {
        load();
    }, [load]);

    const correr = async (accion: () => Promise<any>, exito: string) => {
        setBusy(true);
        setMessage(null);
        try {
            const result = await accion();
            if (result.success && result.data) {
                setState(result.data);
                if (result.data.pin) setPin(result.data.pin);
                setMessage({ type: 'success', text: exito });
                onChanged?.();
            } else {
                setMessage({ type: 'error', text: result.message || 'No se pudo completar el paso' });
            }
        } catch (err: any) {
            setMessage({ type: 'error', text: err?.message || 'No se pudo completar el paso' });
        } finally {
            setBusy(false);
        }
    };

    const quitar = async () => {
        setBusy(true);
        setMessage(null);
        try {
            const result = await removeWhatsAppNumberAction(businessId);
            if (result.success) {
                setState(result.data || null);
                setConfirmando(false);
                setPhone('');
                setName('');
                setCode('');
                setPin('');
                setMessage({ type: 'success', text: 'N\u00famero quitado. Ya puedes agregar el correcto.' });
                onChanged?.();
            } else {
                setMessage({ type: 'error', text: result.message || 'No se pudo quitar el n\u00famero' });
            }
        } catch (err: any) {
            setMessage({ type: 'error', text: err?.message || 'No se pudo quitar el n\u00famero' });
        } finally {
            setBusy(false);
        }
    };

    const soloDigitos = (valor: string) => valor.replace(/\D/g, '');
    const cc = soloDigitos(countryCode);
    const numero = soloDigitos(phone);

    const problemaDelNumero = (): string | null => {
        if (!cc || !numero) return 'Escribe el indicativo y el n\u00famero.';
        if (cc.length > 3) return 'El indicativo no puede pasar de 3 d\u00edgitos.';
        if (cc.length + numero.length > 15) {
            return `Indicativo y n\u00famero suman ${cc.length + numero.length} d\u00edgitos: un n\u00famero internacional no pasa de 15.`;
        }
        if (cc === '57' && (numero.length !== 10 || !numero.startsWith('3'))) {
            return `Un celular colombiano tiene 10 d\u00edgitos y empieza por 3; escribiste ${numero.length}.`;
        }
        if (numero.length < 6) return 'El n\u00famero es demasiado corto.';
        return null;
    };

    const problema = problemaDelNumero();

    const status = state?.status || 'sin_numero';
    const pasoActual = indiceDePaso(status);
    const conectado = status === 'registrado' || status === 'nombre_en_revision';

    return (
        <Card
            icon={<DevicePhoneMobileIcon style={{ color: ACCENT, width: 16, height: 16 }} />}
            title={'Conectar tu propio n\u00famero'}
            description={'Tu n\u00famero queda dentro de la cuenta de WhatsApp de Probability: usa las plantillas ya aprobadas y nosotros gestionamos el consumo. El n\u00famero no puede estar activo en WhatsApp normal.'}
            action={
                status === 'registrado' ? (
                    <Pill tone="ok">Conectado</Pill>
                ) : status === 'nombre_en_revision' ? (
                    <Pill tone="warn">{'Nombre en revisi\u00f3n'}</Pill>
                ) : status === 'sin_numero' ? null : (
                    <Pill tone="warn">En proceso</Pill>
                )
            }
        >
            <Steps pasos={PASOS} actual={pasoActual} />

            {loading ? (
                <p className="text-[13px] text-gray-500 dark:text-gray-400">{'Consultando el estado...'}</p>
            ) : (
                <div className="space-y-3">
                    {status === 'sin_numero' && (
                        <>
                            <div className="grid grid-cols-1 gap-x-4 gap-y-3 sm:grid-cols-3">
                                <div>
                                    <label className={fieldLabel}>
                                        Indicativo <span style={{ color: ACCENT }}>*</span>
                                    </label>
                                    <input
                                        className={inputCls}
                                        style={{ borderColor: '#e9e9f0' }}
                                        value={countryCode}
                                        onChange={(e) => setCountryCode(e.target.value)}
                                        placeholder="57"
                                    />
                                </div>
                                <div className="sm:col-span-2">
                                    <label className={fieldLabel}>
                                        {'N\u00famero de tel\u00e9fono'} <span style={{ color: ACCENT }}>*</span>
                                    </label>
                                    <input
                                        className={`${inputCls} font-mono`}
                                        style={{ borderColor: '#e9e9f0' }}
                                        value={phone}
                                        onChange={(e) => setPhone(e.target.value)}
                                        placeholder="3001234567"
                                    />
                                </div>
                            </div>

                            <div>
                                <label className={fieldLabel}>
                                    {'Nombre que ver\u00e1n tus clientes'} <span style={{ color: ACCENT }}>*</span>
                                </label>
                                <input
                                    className={inputCls}
                                    style={{ borderColor: '#e9e9f0' }}
                                    value={name}
                                    onChange={(e) => setName(e.target.value)}
                                    placeholder="Mi Tienda"
                                />
                                <p className={fieldHint}>
                                    {'Meta revisa este nombre. Debe corresponder a tu negocio real.'}
                                </p>
                            </div>

                            {phone.trim() !== '' && problema && (
                                <p className="text-[12px] font-medium" style={{ color: '#dc2626' }}>
                                    {problema}
                                </p>
                            )}

                            <ActionButton
                                disabled={busy || problema !== null || name.trim() === ''}
                                onClick={() => setConfirmandoAlta(true)}
                            >
                                {'Agregar n\u00famero'}
                            </ActionButton>

                            {confirmandoAlta && (
                                <Alert type="warning">
                                    <p>{'Revisa el n\u00famero antes de crearlo en Meta:'}</p>
                                    <p className="my-2 font-mono text-[18px] font-bold tracking-wider">
                                        {`+${cc} ${numero}`}
                                    </p>
                                    <p>
                                        {`${numero.length} d\u00edgitos, a nombre de "${name.trim()}".`}
                                    </p>
                                    <p className="mt-2">
                                        {'Un n\u00famero de WhatsApp no se puede editar despu\u00e9s: si queda mal escrito hay que quitarlo y volver a crearlo, y el c\u00f3digo de verificaci\u00f3n nunca llega.'}
                                    </p>
                                    <div className="mt-2 flex flex-wrap gap-2">
                                        <ActionButton
                                            disabled={busy}
                                            loading={busy}
                                            onClick={async () => {
                                                await correr(
                                                    () =>
                                                        addWhatsAppNumberAction(
                                                            { country_code: cc, phone_number: numero, verified_name: name },
                                                            businessId
                                                        ),
                                                    'N\u00famero agregado. Pide el c\u00f3digo para verificarlo.'
                                                );
                                                setConfirmandoAlta(false);
                                            }}
                                        >
                                            {'S\u00ed, es correcto'}
                                        </ActionButton>
                                        <ActionButton
                                            variant="ghost"
                                            disabled={busy}
                                            onClick={() => setConfirmandoAlta(false)}
                                        >
                                            {'Corregirlo'}
                                        </ActionButton>
                                    </div>
                                </Alert>
                            )}
                        </>
                    )}

                    {(status === 'esperando_codigo' || status === 'verificado') && (
                        <div
                            className="flex flex-wrap items-center justify-between gap-2 rounded-lg bg-white dark:bg-gray-800 px-3 py-2.5"
                            style={{ border: '1px solid #e9e9f0' }}
                        >
                            <div>
                                <p className="text-[13px] text-gray-800 dark:text-gray-100">
                                    {'Verificando '}
                                    <span className="font-mono font-semibold">
                                        {state?.display_phone_number || state?.phone_number_id}
                                    </span>
                                </p>
                                <p className={fieldHint}>
                                    {'\u00bfEst\u00e1 mal escrito? Qu\u00edtalo y agrega el correcto.'}
                                </p>
                            </div>
                            <ActionButton variant="ghost" disabled={busy} onClick={() => setConfirmando(true)}>
                                {'Cambiar n\u00famero'}
                            </ActionButton>
                        </div>
                    )}

                    {status === 'esperando_codigo' && (
                        <>
                            <p className="text-[13px] text-gray-600 dark:text-gray-300">
                                {'Te enviaremos un c\u00f3digo de 6 d\u00edgitos al tel\u00e9fono para confirmar que es tuyo.'}
                            </p>
                            <div className="flex flex-wrap gap-2">
                                <ActionButton
                                    variant="ghost"
                                    disabled={busy}
                                    onClick={() =>
                                        correr(() => requestWhatsAppNumberCodeAction('SMS', businessId), 'C\u00f3digo enviado por SMS')
                                    }
                                >
                                    {'Enviar por SMS'}
                                </ActionButton>
                                <ActionButton
                                    variant="ghost"
                                    disabled={busy}
                                    onClick={() =>
                                        correr(() => requestWhatsAppNumberCodeAction('VOICE', businessId), 'Te llamaremos con el c\u00f3digo')
                                    }
                                >
                                    {'Recibir llamada'}
                                </ActionButton>
                            </div>
                            <div className="flex flex-col gap-2 sm:flex-row sm:items-end">
                                <div className="flex-1">
                                    <label className={fieldLabel}>{'C\u00f3digo recibido'}</label>
                                    <input
                                        className={`${inputCls} font-mono`}
                                        style={{ borderColor: '#e9e9f0' }}
                                        value={code}
                                        onChange={(e) => setCode(e.target.value)}
                                        placeholder="123456"
                                    />
                                </div>
                                <ActionButton
                                    disabled={busy || !code.trim()}
                                    loading={busy}
                                    onClick={() =>
                                        correr(() => verifyWhatsAppNumberCodeAction(code, businessId), 'N\u00famero verificado')
                                    }
                                >
                                    Verificar
                                </ActionButton>
                            </div>
                        </>
                    )}

                    {status === 'verificado' && (
                        <>
                            <p className="text-[13px] text-gray-600 dark:text-gray-300">
                                {'El n\u00famero es tuyo. Falta activarlo para enviar y recibir mensajes.'}
                            </p>
                            <ActionButton
                                disabled={busy}
                                loading={busy}
                                onClick={() => correr(() => registerWhatsAppNumberAction(businessId), 'N\u00famero activado')}
                            >
                                {'Activar n\u00famero'}
                            </ActionButton>
                        </>
                    )}

                    {(status === 'registrado' || status === 'nombre_en_revision') && (
                        <div className="rounded-lg bg-white dark:bg-gray-800 px-3 py-2.5" style={{ border: '1px solid #e9e9f0' }}>
                            <p className="text-[13px] text-gray-800 dark:text-gray-100">
                                {'Tus mensajes salen desde '}
                                <span className="font-mono font-semibold">
                                    {state?.display_phone_number || state?.phone_number_id}
                                </span>
                                {state?.quality_rating ? ` \u00b7 calidad ${state.quality_rating}` : ''}
                            </p>
                            {status === 'nombre_en_revision' && (
                                <p className={fieldHint}>
                                    {'Meta todav\u00eda est\u00e1 revisando el nombre que ver\u00e1n tus clientes.'}
                                </p>
                            )}
                            <div className="mt-2">
                                <ActionButton variant="ghost" disabled={busy} onClick={() => setConfirmando(true)}>
                                    {'Desconectar n\u00famero'}
                                </ActionButton>
                            </div>
                        </div>
                    )}

                    {confirmando && (
                        <Alert type="error">
                            <p className="font-semibold">
                                {'Esta acci\u00f3n es destructiva y no se puede deshacer.'}
                            </p>
                            <p className="mt-1">
                                {'El n\u00famero '}
                                <span className="font-mono font-semibold">
                                    {state?.display_phone_number || state?.phone_number_id}
                                </span>
                                {' se elimina por completo de la cuenta de WhatsApp de Probability, no solo de esta pantalla. Para volver a usarlo hay que crearlo y verificarlo desde cero.'}
                            </p>
                            {conectado && (
                                <p className="mt-1">
                                    {'Dejas de enviar y recibir por \u00e9l de inmediato; tus mensajes vuelven a salir desde el n\u00famero de Probability.'}
                                </p>
                            )}
                            <p className="mt-1">
                                {'Tu historial de conversaciones y mensajes NO se borra: queda intacto en Probability.'}
                            </p>
                            <div className="mt-2 flex flex-wrap gap-2">
                                <ActionButton disabled={busy} loading={busy} onClick={quitar}>
                                    {conectado ? 'S\u00ed, desconectar' : 'S\u00ed, quitarlo'}
                                </ActionButton>
                                <ActionButton variant="ghost" disabled={busy} onClick={() => setConfirmando(false)}>
                                    Cancelar
                                </ActionButton>
                            </div>
                        </Alert>
                    )}

                    {pin && (
                        <Alert type="success">
                            {'Guarda este PIN de dos pasos, no se vuelve a mostrar: '}
                            <span className="font-mono font-semibold">{pin}</span>
                        </Alert>
                    )}

                    {message && (
                        <Alert type={message.type} onClose={() => setMessage(null)}>
                            {message.text}
                        </Alert>
                    )}
                </div>
            )}
        </Card>
    );
}
