'use client';

import { useCallback, useEffect, useState } from 'react';
import { HexColorPicker, HexColorInput } from 'react-colorful';
import { Modal } from '@/shared/ui/modal';
import { Input } from '@/shared/ui/input';
import { Alert } from '@/shared/ui/alert';
import { FileInput } from '@/shared/ui/file-input';
import { getBusinessByIdAction, updateBusinessAction } from '../../infra/actions';
import { COLOR_PALETTES } from '../../domain/color-palettes';
import type { Business } from '../../domain/types';
import { applyBusinessTheme } from '@/shared/utils/apply-business-theme';
import { getActionError } from '@/shared/utils/action-result';
import { BusinessThemePreview } from './BusinessThemePreview';

interface BusinessSettingsModalProps {
    isOpen: boolean;
    onClose: () => void;
    businessId: number | null;
}

interface Colors {
    primary: string;
    secondary: string;
    tertiary: string;
    quaternary: string;
}

type ColorKey = keyof Colors;

interface BarTokens {
    sidebar: ColorKey;
    topbar: ColorKey;
}

const DEFAULT_BAR_TOKENS: BarTokens = {
    sidebar: 'primary',
    topbar: 'primary',
};

const DEFAULT_COLORS: Colors = {
    primary: '#0f172a',
    secondary: '#be185d',
    tertiary: '#06b6d4',
    quaternary: '#f59e0b',
};

const COLOR_LABELS: Record<keyof Colors, string> = {
    primary: 'Primario',
    secondary: 'Secundario',
    tertiary: 'Terciario',
    quaternary: 'Cuaternario',
};

function matchToken(palette: Colors, color?: string): ColorKey {
    if (!color) return 'primary';
    const found = (Object.keys(palette) as ColorKey[]).find(
        key => palette[key].toLowerCase() === color.toLowerCase(),
    );
    return found || 'primary';
}

interface SwatchProps {
    label: string;
    value: string;
    onChange: (value: string) => void;
}

function ColorSwatch({ label, value, onChange }: SwatchProps) {
    const [open, setOpen] = useState(false);

    return (
        <div className="relative">
            <button
                type="button"
                onClick={() => setOpen(prev => !prev)}
                className="flex w-full items-center gap-2 rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 px-2 py-1.5 text-left transition-colors hover:border-gray-400"
            >
                <span
                    className="h-6 w-6 flex-shrink-0 rounded-md border border-black/10 shadow-sm"
                    style={{ backgroundColor: value }}
                />
                <span className="min-w-0 flex-1">
                    <span className="block text-[10px] font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500">
                        {label}
                    </span>
                    <span className="block truncate text-xs font-mono text-gray-700 dark:text-gray-200">{value}</span>
                </span>
            </button>

            {open && (
                <>
                    <span className="fixed inset-0 z-10" onClick={() => setOpen(false)} />
                    <div className="absolute left-0 top-full z-20 mt-1 w-52 rounded-xl border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 p-3 shadow-xl">
                        <div className="colorPickerBox">
                            <HexColorPicker color={value} onChange={onChange} />
                        </div>
                        <HexColorInput
                            color={value}
                            onChange={onChange}
                            prefixed
                            className="mt-2 w-full rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 px-2 py-1 text-center text-xs font-mono text-gray-700 dark:text-gray-200"
                        />
                        <style jsx>{`
                            .colorPickerBox :global(.react-colorful) {
                                width: 100%;
                                height: 120px;
                            }
                            .colorPickerBox :global(.react-colorful__saturation) {
                                border-radius: 8px 8px 0 0;
                            }
                            .colorPickerBox :global(.react-colorful__hue) {
                                height: 14px;
                                border-radius: 0 0 8px 8px;
                            }
                        `}</style>
                    </div>
                </>
            )}
        </div>
    );
}

export function BusinessSettingsModal({ isOpen, onClose, businessId }: BusinessSettingsModalProps) {
    const [business, setBusiness] = useState<Business | null>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [saved, setSaved] = useState(false);

    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [address, setAddress] = useState('');
    const [colors, setColors] = useState<Colors>(DEFAULT_COLORS);
    const [barTokens, setBarTokens] = useState<BarTokens>(DEFAULT_BAR_TOKENS);
    const [tintBars, setTintBars] = useState(false);
    const [logoFile, setLogoFile] = useState<File | null>(null);
    const [logoPreview, setLogoPreview] = useState<string | null>(null);

    const load = useCallback(async () => {
        if (!businessId) return;

        setLoading(true);
        setError(null);
        setSaved(false);
        setLogoFile(null);

        try {
            const res = await getBusinessByIdAction(businessId);
            const data = (res as { success?: boolean; data?: Business })?.data;
            if (!data) {
                setError('No se pudo cargar la informacion del negocio');
                return;
            }

            setBusiness(data);
            setName(data.name || '');
            setDescription(data.description || '');
            setAddress(data.address || '');
            setColors({
                primary: data.primary_color || DEFAULT_COLORS.primary,
                secondary: data.secondary_color || DEFAULT_COLORS.secondary,
                tertiary: data.tertiary_color || DEFAULT_COLORS.tertiary,
                quaternary: data.quaternary_color || DEFAULT_COLORS.quaternary,
            });
            const paletteColors: Colors = {
                primary: data.primary_color || DEFAULT_COLORS.primary,
                secondary: data.secondary_color || DEFAULT_COLORS.secondary,
                tertiary: data.tertiary_color || DEFAULT_COLORS.tertiary,
                quaternary: data.quaternary_color || DEFAULT_COLORS.quaternary,
            };
            setBarTokens({
                sidebar: matchToken(paletteColors, data.sidebar_color),
                topbar: matchToken(paletteColors, data.topbar_color),
            });
            setTintBars(Boolean(data.sidebar_color || data.topbar_color));
            setLogoPreview(data.logo_url || null);
        } catch (err) {
            setError(getActionError(err, 'Error cargando el negocio'));
        } finally {
            setLoading(false);
        }
    }, [businessId]);

    useEffect(() => {
        if (isOpen) load();
    }, [isOpen, load]);

    useEffect(() => {
        if (!logoFile) return;
        const url = URL.createObjectURL(logoFile);
        setLogoPreview(url);
        return () => URL.revokeObjectURL(url);
    }, [logoFile]);

    const setColor = (key: keyof Colors) => (value: string) => {
        setColors(prev => ({ ...prev, [key]: value }));
        setSaved(false);
    };

    const applyPalette = (palette: Colors) => {
        setColors(palette);
        setSaved(false);
    };

    const setBarToken = (bar: keyof BarTokens, token: ColorKey) => {
        setBarTokens(prev => ({ ...prev, [bar]: token }));
        setSaved(false);
    };

    const bars = {
        sidebar: colors[barTokens.sidebar],
        topbar: colors[barTokens.topbar],
    };

    const handleSave = async () => {
        if (!business) return;

        setSaving(true);
        setError(null);

        try {
            const payload: Record<string, unknown> = {
                name,
                description,
                address,
                primary_color: colors.primary,
                secondary_color: colors.secondary,
                tertiary_color: colors.tertiary,
                quaternary_color: colors.quaternary,
                sidebar_color: tintBars ? bars.sidebar : '',
                topbar_color: tintBars ? bars.topbar : '',
            };
            if (logoFile) payload.logo_file = logoFile;

            const res = await updateBusinessAction(business.id, payload);
            const updated = (res as { success?: boolean; data?: Business })?.data;

            if ((res as { success?: boolean })?.success === false) {
                setError(getActionError(res, 'No se pudo guardar el negocio'));
                return;
            }

            applyBusinessTheme({
                name,
                logo_url: updated?.logo_url || business.logo_url,
                primary_color: colors.primary,
                secondary_color: colors.secondary,
                tertiary_color: colors.tertiary,
                quaternary_color: colors.quaternary,
                sidebar_color: tintBars ? bars.sidebar : undefined,
                topbar_color: tintBars ? bars.topbar : undefined,
            });

            setSaved(true);
            if (updated) setBusiness(updated);
            setLogoFile(null);
        } catch (err) {
            setError(getActionError(err, 'Error guardando el negocio'));
        } finally {
            setSaving(false);
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Informacion del Negocio" size="5xl">
            {loading ? (
                <div className="flex items-center justify-center py-16">
                    <div className="h-8 w-8 animate-spin rounded-full border-b-2" style={{ borderColor: 'var(--color-primary)' }} />
                </div>
            ) : (
                <div className="space-y-4">
                    {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}
                    {saved && <Alert type="success" onClose={() => setSaved(false)}>Cambios guardados</Alert>}

                    <div className="grid grid-cols-1 gap-5 lg:grid-cols-2">
                        <div className="space-y-4">
                            <div className="flex items-center gap-3">
                                <div className="flex h-16 w-16 flex-shrink-0 items-center justify-center overflow-hidden rounded-xl border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-700">
                                    {logoPreview ? (
                                        <img src={logoPreview} alt="Logo" className="max-h-full max-w-full object-contain" />
                                    ) : (
                                        <span className="text-xl font-bold text-gray-300">{name.charAt(0) || '?'}</span>
                                    )}
                                </div>
                                <div className="min-w-0 flex-1">
                                    <FileInput
                                        label="Logo"
                                        accept="image/*"
                                        buttonText="Cambiar logo"
                                        onChange={(file: File | null) => setLogoFile(file)}
                                        helperText="JPG, PNG, WEBP"
                                    />
                                </div>
                            </div>

                            <Input label="Nombre" value={name} onChange={e => setName(e.target.value)} />
                            <Input label="Descripcion" value={description} onChange={e => setDescription(e.target.value)} />
                            <Input label="Direccion" value={address} onChange={e => setAddress(e.target.value)} />

                            <dl className="grid grid-cols-3 gap-2 rounded-xl bg-gray-50 dark:bg-gray-700/40 p-3 text-center">
                                <div>
                                    <dt className="text-[10px] font-semibold uppercase tracking-wide text-gray-400">Codigo</dt>
                                    <dd className="truncate text-xs font-medium text-gray-700 dark:text-gray-200">{business?.code || '-'}</dd>
                                </div>
                                <div>
                                    <dt className="text-[10px] font-semibold uppercase tracking-wide text-gray-400">Estado</dt>
                                    <dd className="text-xs font-medium text-gray-700 dark:text-gray-200">
                                        {business?.is_active ? 'Activo' : 'Inactivo'}
                                    </dd>
                                </div>
                                <div>
                                    <dt className="text-[10px] font-semibold uppercase tracking-wide text-gray-400">Zona horaria</dt>
                                    <dd className="truncate text-xs font-medium text-gray-700 dark:text-gray-200">{business?.timezone || '-'}</dd>
                                </div>
                            </dl>
                        </div>

                        <div className="space-y-3">
                            <div>
                                <span className="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                                    Colores de la plataforma
                                </span>
                                <div className="mt-2 grid grid-cols-2 gap-2">
                                    {(Object.keys(COLOR_LABELS) as (keyof Colors)[]).map(key => (
                                        <ColorSwatch
                                            key={key}
                                            label={COLOR_LABELS[key]}
                                            value={colors[key]}
                                            onChange={setColor(key)}
                                        />
                                    ))}
                                </div>
                            </div>

                            <div>
                                <span className="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                                    Combinaciones sugeridas
                                </span>
                                <div className="mt-2 grid grid-cols-5 gap-1.5">
                                    {COLOR_PALETTES.map(palette => (
                                        <button
                                            key={palette.name}
                                            type="button"
                                            title={palette.name}
                                            onClick={() => applyPalette(palette.colors)}
                                            className="overflow-hidden rounded-md border border-gray-200 dark:border-gray-600 transition-transform hover:scale-105"
                                        >
                                            <span className="flex h-5 w-full">
                                                <span className="flex-1" style={{ backgroundColor: palette.colors.primary }} />
                                                <span className="flex-1" style={{ backgroundColor: palette.colors.secondary }} />
                                                <span className="flex-1" style={{ backgroundColor: palette.colors.tertiary }} />
                                                <span className="flex-1" style={{ backgroundColor: palette.colors.quaternary }} />
                                            </span>
                                        </button>
                                    ))}
                                </div>
                            </div>

                            <div className="rounded-xl border border-gray-200 dark:border-gray-600 p-2.5">
                                <label className="flex cursor-pointer items-center justify-between gap-2">
                                    <span className="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                                        Colorear barras de navegacion
                                    </span>
                                    <input
                                        type="checkbox"
                                        checked={tintBars}
                                        onChange={e => {
                                            setTintBars(e.target.checked);
                                            setSaved(false);
                                        }}
                                        className="h-4 w-4 cursor-pointer"
                                    />
                                </label>

                                {tintBars && (
                                    <div className="mt-2 space-y-2">
                                        {(['sidebar', 'topbar'] as (keyof BarTokens)[]).map(bar => (
                                            <div key={bar} className="flex items-center justify-between gap-2">
                                                <span className="text-[10px] font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500">
                                                    {bar === 'sidebar' ? 'Barra izquierda' : 'Barra superior'}
                                                </span>
                                                <div className="flex items-center gap-1.5">
                                                    {(Object.keys(COLOR_LABELS) as ColorKey[]).map(token => (
                                                        <button
                                                            key={token}
                                                            type="button"
                                                            title={COLOR_LABELS[token]}
                                                            onClick={() => setBarToken(bar, token)}
                                                            className={`h-6 w-6 rounded-md border transition-transform hover:scale-110 ${
                                                                barTokens[bar] === token
                                                                    ? 'border-gray-900 dark:border-white ring-2 ring-offset-1 ring-gray-300 dark:ring-gray-500'
                                                                    : 'border-black/10'
                                                            }`}
                                                            style={{ backgroundColor: colors[token] }}
                                                        />
                                                    ))}
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                )}
                            </div>

                            <BusinessThemePreview
                                colors={colors}
                                bars={tintBars ? bars : null}
                                businessName={name}
                                logoUrl={logoPreview}
                            />
                        </div>
                    </div>

                    <div className="flex justify-end gap-2 border-t border-gray-100 dark:border-gray-700 pt-4">
                        <button
                            type="button"
                            onClick={onClose}
                            className="rounded-lg border border-gray-200 dark:border-gray-600 px-4 py-2 text-sm font-medium text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700"
                        >
                            Cerrar
                        </button>
                        <button
                            type="button"
                            onClick={handleSave}
                            disabled={saving}
                            className="btn-primary rounded-lg px-5 py-2 text-sm font-semibold text-white disabled:opacity-50"
                        >
                            {saving ? 'Guardando...' : 'Guardar cambios'}
                        </button>
                    </div>
                </div>
            )}
        </Modal>
    );
}
