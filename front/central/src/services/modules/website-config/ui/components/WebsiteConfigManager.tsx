'use client';

import { useState, useEffect, useCallback } from 'react';
import { getWebsiteConfigAction, updateWebsiteConfigAction } from '../../infra/actions';
import { WebsiteConfigData, UpdateWebsiteConfigDTO } from '../../domain/types';
import { SectionToggle } from './SectionToggle';
import { PreviewLink } from './PreviewLink';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { TokenStorage } from '@/shared/utils/token-storage';
import { getAvailableTemplates } from '@/services/modules/publicsite/ui/templates/registry';

export function WebsiteConfigManager() {
    const { isSuperAdmin } = usePermissions();
    const { businesses, loading: loadingBusinesses } = useBusinessesSimple();
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [businessCode, setBusinessCode] = useState<string>('');

    const [config, setConfig] = useState<WebsiteConfigData | null>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [message, setMessage] = useState('');

    const effectiveBusinessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    // Resolve business code for the preview link
    useEffect(() => {
        if (isSuperAdmin) {
            if (selectedBusinessId && businesses.length > 0) {
                const selected = businesses.find(b => b.id === selectedBusinessId);
                setBusinessCode(selected?.code || '');
            } else {
                setBusinessCode('');
            }
        } else {
            const businessesData = TokenStorage.getBusinessesData();
            if (businessesData && businessesData.length > 0) {
                setBusinessCode(businessesData[0].code || '');
            }
        }
    }, [isSuperAdmin, selectedBusinessId, businesses]);

    const loadConfig = useCallback(async () => {
        setLoading(true);
        const result = await getWebsiteConfigAction(effectiveBusinessId);
        if (result) {
            setConfig(result);
        }
        setLoading(false);
    }, [effectiveBusinessId]);

    useEffect(() => {
        if (isSuperAdmin && !selectedBusinessId) {
            setLoading(false);
            return;
        }
        loadConfig();
    }, [loadConfig, isSuperAdmin, selectedBusinessId]);

    const handleToggle = (key: string, value: boolean) => {
        if (!config) return;
        setConfig({ ...config, [key]: value });
    };

    const handleSave = async () => {
        if (!config) return;
        setSaving(true);
        setMessage('');

        const dto: UpdateWebsiteConfigDTO = {
            template: config.template,
            show_hero: config.show_hero,
            show_about: config.show_about,
            show_featured_products: config.show_featured_products,
            show_full_catalog: config.show_full_catalog,
            show_testimonials: config.show_testimonials,
            show_location: config.show_location,
            show_contact: config.show_contact,
            show_social_media: config.show_social_media,
            show_whatsapp: config.show_whatsapp,
            hero_content: config.hero_content || undefined,
            about_content: config.about_content || undefined,
            testimonials_content: config.testimonials_content || undefined,
            location_content: config.location_content || undefined,
            contact_content: config.contact_content || undefined,
            social_media_content: config.social_media_content || undefined,
            whatsapp_content: config.whatsapp_content || undefined,
        };

        const result = await updateWebsiteConfigAction(dto, effectiveBusinessId);
        setSaving(false);

        if (result && 'id' in result) {
            setConfig(result as WebsiteConfigData);
            setMessage('Configuracion guardada correctamente');
            setTimeout(() => setMessage(''), 3000);
        } else {
            setMessage('Error al guardar la configuracion');
        }
    };

    const requiresBusinessSelection = isSuperAdmin && !selectedBusinessId;

    return (
        <div className="space-y-6">
            {/* Super admin business selector */}
            {isSuperAdmin && (
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <label className="block text-sm font-medium text-blue-800 mb-2">Seleccionar Negocio</label>
                    <select
                        value={selectedBusinessId?.toString() ?? ''}
                        onChange={(e) => {
                            const val = e.target.value ? Number(e.target.value) : null;
                            setSelectedBusinessId(val);
                            setConfig(null);
                        }}
                        className="w-full px-3 py-2 border border-blue-300 rounded-lg bg-white dark:bg-gray-800 text-gray-900"
                        disabled={loadingBusinesses}
                    >
                        <option value="">-- Selecciona un negocio --</option>
                        {businesses?.map(b => (
                            <option key={b.id} value={b.id}>{b.name} (ID: {b.id})</option>
                        ))}
                    </select>
                </div>
            )}

            {requiresBusinessSelection ? (
                <div className="text-center py-16 text-gray-500">
                    Selecciona un negocio para configurar su sitio web
                </div>
            ) : loading ? (
                <div className="text-center py-16 text-gray-500">Cargando configuracion...</div>
            ) : config ? (
                <>
                    <PreviewLink businessCode={businessCode} />

                    {/* Template selector */}
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-6">
                        <h2 className="text-lg font-semibold text-gray-900 mb-2">Plantilla del sitio</h2>
                        <p className="text-sm text-gray-500 mb-4">
                            Selecciona el diseno visual de tu sitio web
                        </p>
                        <select
                            value={config.template || 'default'}
                            onChange={(e) => setConfig({ ...config, template: e.target.value })}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-white dark:bg-gray-800 text-gray-900"
                        >
                            {getAvailableTemplates().map(t => (
                                <option key={t.id} value={t.id}>{t.name} — {t.description}</option>
                            ))}
                        </select>
                    </div>

                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-6 space-y-4">
                        <h2 className="text-lg font-semibold text-gray-900 mb-4">Secciones de la pagina</h2>

                        <SectionToggle
                            label="Hero / Banner Principal"
                            description="Banner grande con titulo, subtitulo y boton de accion"
                            checked={config.show_hero}
                            onChange={(v) => handleToggle('show_hero', v)}
                        />
                        <SectionToggle
                            label="Sobre Nosotros"
                            description="Seccion con informacion de la empresa, mision y vision"
                            checked={config.show_about}
                            onChange={(v) => handleToggle('show_about', v)}
                        />
                        <SectionToggle
                            label="Productos Destacados"
                            description="Muestra los productos marcados como destacados"
                            checked={config.show_featured_products}
                            onChange={(v) => handleToggle('show_featured_products', v)}
                        />
                        <SectionToggle
                            label="Catalogo Completo"
                            description="Link al catalogo completo de productos y CTA de pedido"
                            checked={config.show_full_catalog}
                            onChange={(v) => handleToggle('show_full_catalog', v)}
                        />
                        <SectionToggle
                            label="Testimonios"
                            description="Opiniones y resenas de clientes"
                            checked={config.show_testimonials}
                            onChange={(v) => handleToggle('show_testimonials', v)}
                        />
                        <SectionToggle
                            label="Ubicacion"
                            description="Mapa con la direccion y horarios"
                            checked={config.show_location}
                            onChange={(v) => handleToggle('show_location', v)}
                        />
                        <SectionToggle
                            label="Formulario de Contacto"
                            description="Formulario para que los visitantes envien mensajes"
                            checked={config.show_contact}
                            onChange={(v) => handleToggle('show_contact', v)}
                        />
                        <SectionToggle
                            label="Redes Sociales"
                            description="Links a Facebook, Instagram, Twitter, TikTok"
                            checked={config.show_social_media}
                            onChange={(v) => handleToggle('show_social_media', v)}
                        />
                        <SectionToggle
                            label="Boton WhatsApp"
                            description="Boton flotante de WhatsApp para contacto rapido"
                            checked={config.show_whatsapp}
                            onChange={(v) => handleToggle('show_whatsapp', v)}
                        />
                    </div>

                    {/* Save button */}
                    <div className="flex items-center gap-4">
                        <button
                            onClick={handleSave}
                            disabled={saving}
                            className="px-6 py-3 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors disabled:opacity-50"
                        >
                            {saving ? 'Guardando...' : 'Guardar Configuracion'}
                        </button>
                        {message && (
                            <span className={`text-sm ${message.includes('Error') ? 'text-red-600' : 'text-green-600'}`}>
                                {message}
                            </span>
                        )}
                    </div>
                </>
            ) : (
                <div className="text-center py-16 text-gray-500">No se pudo cargar la configuracion</div>
            )}
        </div>
    );
}
