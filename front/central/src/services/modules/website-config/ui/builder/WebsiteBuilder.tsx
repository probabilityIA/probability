'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { getWebsiteConfigAction, updateWebsiteConfigAction, deleteWebsiteImageAction, getWebsiteCategoriesAction } from '../../infra/actions';
import { getProductsAction } from '@/services/modules/products/infra/actions';
import { WebsiteConfigData, UpdateWebsiteConfigDTO, WebsiteCategory } from '../../domain/types';
import { PublicBusiness, PublicProduct, WebsiteConfig, SectionKey } from '@/services/modules/publicsite/domain/types';
import { getTemplate, getAvailableTemplates } from '@/services/modules/publicsite/ui/templates/registry';
import { SectionRenderer } from '@/services/modules/publicsite/ui/components/SectionRenderer';
import { CartProvider } from '@/services/modules/publicsite/ui/cart/cart-context';
import { buildThemeStyle } from '@/services/modules/publicsite/ui/theme';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useSelectedBusiness } from '@/shared/contexts/selected-business-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import { TokenStorage } from '@/shared/utils/token-storage';
import { SECTION_META, normalizeOrder } from './section-meta';
import { SectionOverlay } from './SectionOverlay';
import { ThemeEditor } from './ThemeEditor';
import { HeroEditor, AboutEditor, TestimonialsEditor, LocationEditor, ContactEditor, SocialMediaEditor, WhatsAppEditor } from './editors';
import { ArrowTopRightOnSquareIcon } from '@heroicons/react/24/outline';

function collectWebsiteImages(c: WebsiteConfigData | null): string[] {
    if (!c) return [];
    const urls: string[] = [];
    const hero = c.hero_content as Record<string, any> | null;
    if (hero?.background_image) urls.push(hero.background_image);
    const about = c.about_content as Record<string, any> | null;
    if (about?.image) urls.push(about.image);
    (c.testimonials_content || []).forEach((t: Record<string, any>) => {
        if (t?.avatar) urls.push(t.avatar);
    });
    const themeContent = c.theme_content as Record<string, any> | null;
    if (themeContent?.background_image) urls.push(themeContent.background_image);
    return urls.filter(u => typeof u === 'string' && u.includes('/website/'));
}

interface BusinessInfo {
    id: number;
    name: string;
    code: string;
    logo_url?: string;
    primary_color?: string;
    secondary_color?: string;
    tertiary_color?: string;
    quaternary_color?: string;
}

export function WebsiteBuilder() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useSelectedBusiness();
    const { businesses } = useBusinessesSimple();

    const [config, setConfig] = useState<WebsiteConfigData | null>(null);
    const [order, setOrder] = useState<SectionKey[]>(normalizeOrder(null));
    const [products, setProducts] = useState<PublicProduct[]>([]);
    const [categories, setCategories] = useState<WebsiteCategory[]>([]);
    const [selectedSection, setSelectedSection] = useState<SectionKey | null>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [dirty, setDirty] = useState(false);
    const [message, setMessage] = useState('');
    const dragIndex = useRef<number | null>(null);
    const savedImagesRef = useRef<string[]>([]);
    const autoSaveRef = useRef(false);

    const effectiveBusinessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    const businessInfo: BusinessInfo | null = (() => {
        if (isSuperAdmin) {
            const b = businesses.find(x => x.id === selectedBusinessId);
            return b ? { id: b.id, name: b.name, code: b.code || '', logo_url: b.logo_url, primary_color: b.primary_color, secondary_color: b.secondary_color, tertiary_color: b.tertiary_color, quaternary_color: b.quaternary_color } : null;
        }
        const data = TokenStorage.getBusinessesData();
        if (data && data.length > 0) {
            const b = data[0];
            return { id: b.id, name: b.name, code: b.code, logo_url: b.logo_url, primary_color: b.primary_color, secondary_color: b.secondary_color, tertiary_color: b.tertiary_color, quaternary_color: b.quaternary_color };
        }
        return null;
    })();

    const loadAll = useCallback(async () => {
        setLoading(true);
        const [cfg, prods, cats] = await Promise.all([
            getWebsiteConfigAction(effectiveBusinessId),
            getProductsAction({ page: 1, page_size: 8, business_id: effectiveBusinessId }),
            getWebsiteCategoriesAction(effectiveBusinessId),
        ]);
        setCategories(cats || []);
        if (cfg) {
            setConfig(cfg);
            setOrder(normalizeOrder(cfg.sections_order));
            savedImagesRef.current = collectWebsiteImages(cfg);
        }
        const list = (prods?.data || []) as any[];
        const featured = list.filter(p => p.is_featured);
        const source = featured.length > 0 ? featured : list.slice(0, 4);
        setProducts(source.map(p => ({
            id: String(p.id),
            name: p.name,
            description: p.description || '',
            short_description: p.short_description || '',
            price: p.price || 0,
            currency: p.currency || 'COP',
            image_url: p.image_url || '',
            sku: p.sku || '',
            stock_quantity: p.stock_quantity ?? 0,
            category: p.category || '',
            brand: p.brand || '',
            is_featured: !!p.is_featured,
            created_at: p.created_at || '',
        })));
        setDirty(false);
        setSelectedSection(null);
        setLoading(false);
    }, [effectiveBusinessId]);

    useEffect(() => {
        if (isSuperAdmin && !selectedBusinessId) {
            setLoading(false);
            return;
        }
        loadAll();
    }, [loadAll, isSuperAdmin, selectedBusinessId]);

    const patchConfig = (patch: Partial<WebsiteConfigData>) => {
        setConfig(prev => prev ? { ...prev, ...patch } : prev);
        setDirty(true);
    };

    const moveSection = (from: number, to: number) => {
        if (to < 0 || to >= order.length) return;
        const next = [...order];
        const [item] = next.splice(from, 1);
        next.splice(to, 0, item);
        setOrder(next);
        setDirty(true);
    };

    const handleSave = async () => {
        if (!config) return;
        setSaving(true);
        setMessage('');
        const dto: UpdateWebsiteConfigDTO = {
            template: config.template,
            sections_order: order,
            theme_content: config.theme_content || undefined,
            hidden_categories: config.hidden_categories || [],
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
            const saved = result as WebsiteConfigData;
            setConfig(saved);
            setOrder(normalizeOrder(saved.sections_order));
            setDirty(false);
            setMessage('Cambios publicados');
            setTimeout(() => setMessage(''), 3000);

            const currentImages = collectWebsiteImages(saved);
            const orphans = savedImagesRef.current.filter(u => !currentImages.includes(u));
            savedImagesRef.current = currentImages;
            for (const url of orphans) {
                deleteWebsiteImageAction(url, effectiveBusinessId);
            }
        } else {
            setMessage('Error al guardar');
        }
    };

    useEffect(() => {
        if (autoSaveRef.current) {
            autoSaveRef.current = false;
            handleSave();
        }
    });

    const handleImageDeleted = () => {
        autoSaveRef.current = true;
    };

    if (isSuperAdmin && !selectedBusinessId) {
        return (
            <div className="text-center py-24 text-gray-500 dark:text-gray-400">
                Selecciona un negocio en la barra superior para disenar su tienda
            </div>
        );
    }

    if (loading) {
        return <div className="text-center py-24 text-gray-500 dark:text-gray-400">Cargando el editor de tu tienda...</div>;
    }

    if (!config || !businessInfo) {
        return <div className="text-center py-24 text-gray-500 dark:text-gray-400">No se pudo cargar la configuracion del sitio</div>;
    }

    const template = getTemplate(config.template || 'default');

    const theme = (config.theme_content || {}) as Record<string, string>;
    const fallbackColors = {
        primary: businessInfo.primary_color || '#1f2937',
        secondary: businessInfo.secondary_color || '#3b82f6',
        tertiary: businessInfo.tertiary_color || '#10b981',
        quaternary: businessInfo.quaternary_color || '#fbbf24',
        background: '#ffffff',
    };

    const canvasConfig: WebsiteConfig = {
        template: config.template,
        sections_order: order,
        theme_content: theme,
        show_hero: config.show_hero,
        show_about: config.show_about,
        show_featured_products: config.show_featured_products,
        show_full_catalog: config.show_full_catalog,
        show_testimonials: config.show_testimonials,
        show_location: config.show_location,
        show_contact: config.show_contact,
        show_social_media: config.show_social_media,
        show_whatsapp: config.show_whatsapp,
        hero_content: config.hero_content as any,
        about_content: config.about_content as any,
        testimonials_content: config.testimonials_content as any,
        location_content: config.location_content as any,
        contact_content: config.contact_content as any,
        social_media_content: config.social_media_content as any,
        whatsapp_content: config.whatsapp_content as any,
    };

    const canvasBusiness: PublicBusiness = {
        id: businessInfo.id,
        name: businessInfo.name,
        code: businessInfo.code,
        description: '',
        logo_url: businessInfo.logo_url || '',
        primary_color: theme.primary || fallbackColors.primary,
        secondary_color: theme.secondary || fallbackColors.secondary,
        tertiary_color: theme.tertiary || fallbackColors.tertiary,
        quaternary_color: theme.quaternary || fallbackColors.quaternary,
        navbar_image_url: '',
        website_config: canvasConfig,
        featured_products: products,
    };

    const meta = selectedSection ? SECTION_META[selectedSection] : null;

    return (
        <div className="flex h-[calc(100vh-4rem)] overflow-hidden">
            <div
                className="flex-1 overflow-y-auto bg-gray-100 dark:bg-gray-900 p-4"
                onClickCapture={(e) => {
                    const anchor = (e.target as HTMLElement).closest('a');
                    if (anchor) e.preventDefault();
                }}
                onClick={() => setSelectedSection(null)}
            >
                <div
                    className="mx-auto max-w-5xl shadow-xl rounded-lg overflow-hidden"
                    style={buildThemeStyle(theme, fallbackColors)}
                >
                    <CartProvider slug={businessInfo.code}>
                        <template.Nav business={canvasBusiness} />
                        {order.map((key, index) => {
                            const m = SECTION_META[key];
                            const visible = (config as unknown as Record<string, boolean>)[m.toggleField];
                            const rendered = (
                                <SectionRenderer
                                    sectionKey={key}
                                    business={canvasBusiness}
                                    slug={businessInfo.code}
                                    config={canvasConfig}
                                    template={template}
                                />
                            );
                            const isEmpty = m.contentField && !(config as unknown as Record<string, unknown>)[m.contentField]
                                && key !== 'hero' && key !== 'contact';
                            return (
                                <SectionOverlay
                                    key={key}
                                    label={m.label}
                                    visible={visible}
                                    selected={selectedSection === key}
                                    isFirst={index === 0}
                                    isLast={index === order.length - 1}
                                    hasEditor={!!m.contentField}
                                    onSelect={() => setSelectedSection(key)}
                                    onToggleVisible={() => patchConfig({ [m.toggleField]: !visible } as Partial<WebsiteConfigData>)}
                                    onMoveUp={() => moveSection(index, index - 1)}
                                    onMoveDown={() => moveSection(index, index + 1)}
                                    onDragStart={() => { dragIndex.current = index; }}
                                    onDragOver={(e) => e.preventDefault()}
                                    onDrop={(e) => {
                                        e.preventDefault();
                                        if (dragIndex.current !== null && dragIndex.current !== index) {
                                            moveSection(dragIndex.current, index);
                                        }
                                        dragIndex.current = null;
                                    }}
                                >
                                    {isEmpty ? (
                                        <div className="py-16 text-center text-gray-400 border-y border-dashed border-gray-300 dark:border-gray-700">
                                            <p className="font-medium">{m.label}</p>
                                            <p className="text-sm">Sin contenido: haz clic para editar esta seccion</p>
                                        </div>
                                    ) : rendered}
                                </SectionOverlay>
                            );
                        })}
                        <template.Footer business={canvasBusiness} />
                    </CartProvider>
                </div>
            </div>

            <aside className="w-96 shrink-0 border-l border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 overflow-y-auto">
                <div className="p-4 border-b border-gray-200 dark:border-gray-700 sticky top-0 bg-white dark:bg-gray-800 z-10 space-y-3">
                    <div className="flex items-center justify-between">
                        <h2 className="font-semibold text-gray-900 dark:text-white">Disena tu tienda</h2>
                        <a
                            href={`/tienda/${businessInfo.code}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center gap-1 text-sm text-blue-600 hover:underline"
                        >
                            Ver tienda <ArrowTopRightOnSquareIcon className="w-4 h-4" />
                        </a>
                    </div>
                    <button
                        onClick={handleSave}
                        disabled={saving || !dirty}
                        className="w-full py-2.5 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors disabled:opacity-40"
                    >
                        {saving ? 'Publicando...' : dirty ? 'Publicar cambios' : 'Sin cambios pendientes'}
                    </button>
                    {message && (
                        <p className={`text-sm text-center ${message.includes('Error') ? 'text-red-600' : 'text-green-600'}`}>{message}</p>
                    )}
                </div>

                <div className="p-4 space-y-6">
                    {!selectedSection && (
                        <>
                            <div>
                                <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">Plantilla</h3>
                                <select
                                    value={config.template || 'default'}
                                    onChange={(e) => patchConfig({ template: e.target.value })}
                                    className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                                >
                                    {getAvailableTemplates().map(t => (
                                        <option key={t.id} value={t.id}>{t.name}</option>
                                    ))}
                                </select>
                            </div>

                            <div>
                                <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">Tema de colores</h3>
                                <ThemeEditor
                                    theme={config.theme_content || null}
                                    fallback={fallbackColors}
                                    onChange={(t) => patchConfig({ theme_content: t })}
                                    businessId={effectiveBusinessId}
                                    onImageDeleted={handleImageDeleted}
                                />
                            </div>

                            <div>
                                <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">Boton de WhatsApp</h3>
                                <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200 mb-3">
                                    <input
                                        type="checkbox"
                                        checked={config.show_whatsapp}
                                        onChange={(e) => patchConfig({ show_whatsapp: e.target.checked })}
                                        className="rounded border-gray-300"
                                    />
                                    Activar WhatsApp en la tienda
                                </label>
                                {config.show_whatsapp && (
                                    <WhatsAppEditor
                                        content={config.whatsapp_content}
                                        onChange={(c) => patchConfig({ whatsapp_content: c })}
                                    />
                                )}
                            </div>

                            {categories.length > 0 && (
                                <div>
                                    <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-2">Categorias en la tienda</h3>
                                    <p className="text-xs text-gray-400 mb-2">Desactiva las categorias que no quieres mostrar; sus productos se ocultan del catalogo.</p>
                                    <div className="space-y-1.5">
                                        {categories.map(cat => {
                                            const hidden = (config.hidden_categories || []).includes(cat.name);
                                            return (
                                                <label key={cat.name} className="flex items-center justify-between gap-2 text-sm text-gray-700 dark:text-gray-200">
                                                    <span className="flex items-center gap-2">
                                                        <input
                                                            type="checkbox"
                                                            checked={!hidden}
                                                            onChange={(e) => {
                                                                const current = config.hidden_categories || [];
                                                                const next = e.target.checked
                                                                    ? current.filter(c => c !== cat.name)
                                                                    : [...current, cat.name];
                                                                patchConfig({ hidden_categories: next });
                                                            }}
                                                            className="rounded border-gray-300"
                                                        />
                                                        {cat.name}
                                                    </span>
                                                    <span className="text-xs text-gray-400">{cat.product_count}</span>
                                                </label>
                                            );
                                        })}
                                    </div>
                                </div>
                            )}

                            <div className="text-sm text-gray-500 dark:text-gray-400 border-t border-gray-200 dark:border-gray-700 pt-4 space-y-2">
                                <p>Haz clic en una seccion del lienzo para editar su contenido.</p>
                                <p>Usa las flechas o arrastra para reordenar, y el ojo para mostrar u ocultar.</p>
                            </div>
                        </>
                    )}

                    {selectedSection && meta && (
                        <div className="space-y-4">
                            <div className="flex items-center justify-between">
                                <h3 className="text-sm font-semibold text-gray-900 dark:text-white">{meta.label}</h3>
                                <button
                                    onClick={() => setSelectedSection(null)}
                                    className="text-xs text-gray-400 hover:text-gray-600"
                                >
                                    Cerrar
                                </button>
                            </div>

                            <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200">
                                <input
                                    type="checkbox"
                                    checked={(config as unknown as Record<string, boolean>)[meta.toggleField]}
                                    onChange={(e) => patchConfig({ [meta.toggleField]: e.target.checked } as Partial<WebsiteConfigData>)}
                                    className="rounded border-gray-300"
                                />
                                Seccion visible
                            </label>

                            {selectedSection === 'hero' && (
                                <HeroEditor content={config.hero_content} onChange={(c) => patchConfig({ hero_content: c })} businessId={effectiveBusinessId} onImageDeleted={handleImageDeleted} />
                            )}
                            {selectedSection === 'about' && (
                                <AboutEditor content={config.about_content} onChange={(c) => patchConfig({ about_content: c })} businessId={effectiveBusinessId} onImageDeleted={handleImageDeleted} />
                            )}
                            {selectedSection === 'testimonials' && (
                                <TestimonialsEditor content={config.testimonials_content} onChange={(c) => patchConfig({ testimonials_content: c })} />
                            )}
                            {selectedSection === 'location' && (
                                <LocationEditor content={config.location_content} onChange={(c) => patchConfig({ location_content: c })} businessId={effectiveBusinessId} />
                            )}
                            {selectedSection === 'contact' && (
                                <ContactEditor content={config.contact_content} onChange={(c) => patchConfig({ contact_content: c })} />
                            )}
                            {selectedSection === 'social_media' && (
                                <SocialMediaEditor content={config.social_media_content} onChange={(c) => patchConfig({ social_media_content: c })} />
                            )}
                            {selectedSection === 'featured_products' && (
                                <p className="text-sm text-gray-500 dark:text-gray-400">
                                    Esta seccion muestra los productos marcados como destacados. Gestionalo desde el modulo de Productos.
                                </p>
                            )}
                            {selectedSection === 'full_catalog' && (
                                <p className="text-sm text-gray-500 dark:text-gray-400">
                                    Bloque de llamado a la accion que invita a explorar el catalogo completo.
                                </p>
                            )}
                        </div>
                    )}
                </div>
            </aside>
        </div>
    );
}
