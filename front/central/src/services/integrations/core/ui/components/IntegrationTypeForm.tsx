'use client';

import { useState, useEffect } from 'react';
import {
    EyeIcon,
    EyeSlashIcon,
    Cog6ToothIcon,
    PhotoIcon,
    CodeBracketIcon,
    GlobeAltIcon,
    DocumentTextIcon,
    KeyIcon,
    InformationCircleIcon,
} from '@heroicons/react/24/outline';
import { useIntegrationTypes } from '../hooks/useIntegrationTypes';
import { IntegrationType, CreateIntegrationTypeDTO, UpdateIntegrationTypeDTO, IntegrationCategory } from '../../domain/types';
import { Select, Button, Alert, FileInput } from '@/shared/ui';
import { getIntegrationCategoriesAction, getIntegrationTypePlatformCredentialsAction } from '../../infra/actions';
import { WhatsAppTypeCredentialsForm } from '@/services/integrations/messages/whatsapp/ui/components';
import type { WhatsAppPlatformCredentials } from '@/services/integrations/messages/whatsapp/ui/components';
import { BoldTypeCredentialsForm } from '@/services/integrations/pay/bold/ui/components';
import type { BoldPlatformCredentials } from '@/services/integrations/pay/bold/ui/components';
import { getActionError } from '@/shared/utils/action-result';

const WHATSAPP_TYPE_ID = 2;
const BOLD_TYPE_ID = 23;

const ACCENT = 'var(--color-primary)';
const ACCENT_DARK = 'color-mix(in srgb, var(--color-primary) 85%, black)';
const ACCENT_SOFT = 'color-mix(in srgb, var(--color-primary) 10%, white)';
const ACCENT_BORDER = 'color-mix(in srgb, var(--color-primary) 25%, white)';
const CARD_BG = '#fafafd';
const CARD_BORDER = '#eceaf3';
const INPUT_BORDER = '#e9e9f0';

const fieldLabel = 'block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1';
const fieldHint = 'text-[11px] text-gray-400 dark:text-gray-500 mt-1 flex items-start gap-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)]';
const jsonCls = 'w-full px-3 py-2 bg-gray-900 text-green-400 border border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/40 focus:border-[var(--color-primary)] font-mono text-xs';

interface IntegrationTypeFormProps {
    integrationType?: IntegrationType;
    onSuccess?: () => void;
    onCancel?: () => void;
}

function SectionCard({
    icon: Icon,
    title,
    bg = CARD_BG,
    children,
}: {
    icon: React.ComponentType<{ style?: React.CSSProperties }>;
    title: string;
    bg?: string;
    children: React.ReactNode;
}) {
    return (
        <div className="rounded-xl p-4 dark:bg-gray-800/60" style={{ backgroundColor: bg, border: `1px solid ${CARD_BORDER}` }}>
            <div className="flex items-center gap-2 mb-3">
                <span className="flex h-7 w-7 items-center justify-center rounded-md" style={{ backgroundColor: ACCENT_SOFT }}>
                    <Icon style={{ color: ACCENT, width: 16, height: 16 }} />
                </span>
                <h3 className="text-sm font-bold text-gray-900 dark:text-white">{title}</h3>
            </div>
            {children}
        </div>
    );
}

export default function IntegrationTypeForm({ integrationType, onSuccess, onCancel }: IntegrationTypeFormProps) {
    const { createIntegrationType, updateIntegrationType } = useIntegrationTypes();

    const [formData, setFormData] = useState({
        name: '',
        code: '',
        description: '',
        category_id: 0,
        is_active: true,
        config_schema: '{}',
        credentials_schema: '{}',
        setup_instructions: '',
        base_url: '',
        base_url_test: '',
        platform_credentials: '{}',
    });

    const [categories, setCategories] = useState<IntegrationCategory[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [imageFile, setImageFile] = useState<File | null>(null);
    const [removeImage, setRemoveImage] = useState(false);
    const [imagePreview, setImagePreview] = useState<string | null>(null);
    const [showPlatformCredentials, setShowPlatformCredentials] = useState(false);
    const [whatsappCredentials, setWhatsappCredentials] = useState<WhatsAppPlatformCredentials>({
        whatsapp_url: '',
        webhook_callback_url: '',
        phone_number_id: '',
        access_token: '',
        verify_token: '',
        webhook_secret: '',
        test_phone_number: '',
        ai_sales_enabled: false,
        ai_sales_model_id: 'amazon.nova-micro-v1:0',
        ai_sales_session_ttl_minutes: '20',
        ai_sales_max_tool_iterations: '5',
        ai_sales_demo_business_id: '1',
    });
    const [boldCredentials, setBoldCredentials] = useState<BoldPlatformCredentials>({
        api_key: '',
        secret_key: '',
        test_api_key: '',
        test_secret_key: '',
        link_api_key: '',
        link_secret_key: '',
        test_link_api_key: '',
        test_link_secret_key: '',
    });
    const [boldWebhookUrls, setBoldWebhookUrls] = useState<{ production?: string; sandbox?: string }>({});

    useEffect(() => {
        getIntegrationCategoriesAction()
            .then((res) => {
                if (res.success && res.data.length > 0) {
                    setCategories(res.data);
                    if (!integrationType) {
                        setFormData((prev) => ({ ...prev, category_id: res.data[0].id }));
                    }
                }
            })
            .catch(() => {});
    }, []);

    useEffect(() => {
        if (integrationType) {
            setFormData({
                name: integrationType.name,
                code: integrationType.code,
                description: integrationType.description || '',
                category_id: integrationType.category_id || 0,
                is_active: integrationType.is_active,
                config_schema: JSON.stringify(integrationType.config_schema || {}, null, 2),
                credentials_schema: JSON.stringify(integrationType.credentials_schema || {}, null, 2),
                setup_instructions: integrationType.setup_instructions || '',
                base_url: integrationType.base_url || '',
                base_url_test: integrationType.base_url_test || '',
                platform_credentials: '{}',
            });
            if (integrationType.image_url) {
                setImagePreview(integrationType.image_url);
            }
            if (integrationType.has_platform_credentials) {
                getIntegrationTypePlatformCredentialsAction(integrationType.id)
                    .then((res) => {
                        if (res.success && res.data && Object.keys(res.data).length > 0) {
                            if (integrationType.id === WHATSAPP_TYPE_ID) {
                                const d = res.data as Record<string, unknown>;
                                setWhatsappCredentials({
                                    whatsapp_url: String(d.whatsapp_url || ''),
                                    webhook_callback_url: String(d.webhook_callback_url || ''),
                                    phone_number_id: String(d.phone_number_id || ''),
                                    access_token: String(d.access_token || ''),
                                    verify_token: String(d.verify_token || ''),
                                    webhook_secret: String(d.webhook_secret || ''),
                                    test_phone_number: String(d.test_phone_number || ''),
                                    ai_sales_enabled: d.ai_sales_enabled === true,
                                    ai_sales_model_id: String(d.ai_sales_model_id || 'amazon.nova-micro-v1:0'),
                                    ai_sales_session_ttl_minutes: String(d.ai_sales_session_ttl_minutes || '20'),
                                    ai_sales_max_tool_iterations: String(d.ai_sales_max_tool_iterations || '5'),
                                    ai_sales_demo_business_id: String(d.ai_sales_demo_business_id || '1'),
                                });
                            } else if (integrationType.id === BOLD_TYPE_ID) {
                                const d = res.data as Record<string, unknown>;
                                setBoldCredentials({
                                    api_key: String(d.api_key || ''),
                                    secret_key: String(d.secret_key || ''),
                                    test_api_key: String(d.test_api_key || ''),
                                    test_secret_key: String(d.test_secret_key || ''),
                                    link_api_key: String(d.link_api_key || ''),
                                    link_secret_key: String(d.link_secret_key || ''),
                                    test_link_api_key: String(d.test_link_api_key || ''),
                                    test_link_secret_key: String(d.test_link_secret_key || ''),
                                });
                                if (res.webhook_urls) {
                                    setBoldWebhookUrls({
                                        production: res.webhook_urls.production,
                                        sandbox: res.webhook_urls.sandbox,
                                    });
                                }
                            } else {
                                setFormData((prev) => ({
                                    ...prev,
                                    platform_credentials: JSON.stringify(res.data, null, 2),
                                }));
                            }
                        }
                    })
                    .catch(() => {});
            }
        }
    }, [integrationType]);

    const handleImageChange = (file: File | null) => {
        setImageFile(file);
        setRemoveImage(false);
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => {
                setImagePreview(reader.result as string);
            };
            reader.readAsDataURL(file);
        } else {
            setImagePreview(integrationType?.image_url || null);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            let success = false;
            let platformCredentials: Record<string, unknown> | undefined;
            if (integrationType?.id === WHATSAPP_TYPE_ID) {
                const wa: Record<string, unknown> = {};
                if (whatsappCredentials.whatsapp_url.trim()) wa.whatsapp_url = whatsappCredentials.whatsapp_url.trim();
                if (whatsappCredentials.webhook_callback_url.trim()) wa.webhook_callback_url = whatsappCredentials.webhook_callback_url.trim();
                if (whatsappCredentials.phone_number_id.trim()) wa.phone_number_id = whatsappCredentials.phone_number_id.trim();
                if (whatsappCredentials.access_token.trim()) wa.access_token = whatsappCredentials.access_token.trim();
                if (whatsappCredentials.verify_token.trim()) wa.verify_token = whatsappCredentials.verify_token.trim();
                if (whatsappCredentials.webhook_secret.trim()) wa.webhook_secret = whatsappCredentials.webhook_secret.trim();
                if (whatsappCredentials.test_phone_number.trim()) wa.test_phone_number = whatsappCredentials.test_phone_number.trim();
                wa.ai_sales_enabled = whatsappCredentials.ai_sales_enabled;
                if (whatsappCredentials.ai_sales_model_id.trim()) wa.ai_sales_model_id = whatsappCredentials.ai_sales_model_id.trim();
                if (whatsappCredentials.ai_sales_session_ttl_minutes.trim()) wa.ai_sales_session_ttl_minutes = Number(whatsappCredentials.ai_sales_session_ttl_minutes);
                if (whatsappCredentials.ai_sales_max_tool_iterations.trim()) wa.ai_sales_max_tool_iterations = Number(whatsappCredentials.ai_sales_max_tool_iterations);
                if (whatsappCredentials.ai_sales_demo_business_id.trim()) wa.ai_sales_demo_business_id = Number(whatsappCredentials.ai_sales_demo_business_id);
                if (Object.keys(wa).length > 0) platformCredentials = wa;
            } else if (integrationType?.id === BOLD_TYPE_ID) {
                const bold: Record<string, unknown> = {};
                if (boldCredentials.api_key.trim()) bold.api_key = boldCredentials.api_key.trim();
                if (boldCredentials.secret_key.trim()) bold.secret_key = boldCredentials.secret_key.trim();
                if (boldCredentials.test_api_key.trim()) bold.test_api_key = boldCredentials.test_api_key.trim();
                if (boldCredentials.test_secret_key.trim()) bold.test_secret_key = boldCredentials.test_secret_key.trim();
                if (boldCredentials.link_api_key.trim()) bold.link_api_key = boldCredentials.link_api_key.trim();
                if (boldCredentials.link_secret_key.trim()) bold.link_secret_key = boldCredentials.link_secret_key.trim();
                if (boldCredentials.test_link_api_key.trim()) bold.test_link_api_key = boldCredentials.test_link_api_key.trim();
                if (boldCredentials.test_link_secret_key.trim()) bold.test_link_secret_key = boldCredentials.test_link_secret_key.trim();
                if (Object.keys(bold).length > 0) platformCredentials = bold;
            } else {
                try {
                    const parsed = formData.platform_credentials ? JSON.parse(formData.platform_credentials) : {};
                    if (parsed && typeof parsed === 'object' && Object.keys(parsed).length > 0) {
                        platformCredentials = parsed;
                    }
                } catch {
                    throw new Error('Las credenciales de plataforma no son un JSON valido');
                }
            }

            if (integrationType) {
                const updateData: UpdateIntegrationTypeDTO = {
                    name: formData.name,
                    code: formData.code,
                    description: formData.description,
                    category_id: formData.category_id,
                    is_active: formData.is_active,
                    config_schema: formData.config_schema ? JSON.parse(formData.config_schema) : undefined,
                    credentials_schema: formData.credentials_schema ? JSON.parse(formData.credentials_schema) : undefined,
                    setup_instructions: formData.setup_instructions,
                    image_file: imageFile || undefined,
                    remove_image: removeImage || undefined,
                    base_url: formData.base_url || undefined,
                    base_url_test: formData.base_url_test || undefined,
                    platform_credentials: platformCredentials,
                };
                success = await updateIntegrationType(integrationType.id, updateData);
            } else {
                const createData: CreateIntegrationTypeDTO = {
                    name: formData.name,
                    code: formData.code,
                    description: formData.description,
                    category_id: formData.category_id,
                    is_active: formData.is_active,
                    config_schema: formData.config_schema ? JSON.parse(formData.config_schema) : undefined,
                    credentials_schema: formData.credentials_schema ? JSON.parse(formData.credentials_schema) : undefined,
                    setup_instructions: formData.setup_instructions,
                    image_file: imageFile || undefined,
                    base_url: formData.base_url || undefined,
                    base_url_test: formData.base_url_test || undefined,
                    platform_credentials: platformCredentials,
                };
                success = await createIntegrationType(createData);
            }

            if (success) {
                onSuccess?.();
            }
        } catch (err: any) {
            console.error('Error saving integration type:', err);
            setError(getActionError(err, 'Error al guardar el tipo de integracion'));
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <SectionCard icon={Cog6ToothIcon} title="Informacion basica">
                <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>
                            Nombre <span style={{ color: ACCENT }}>*</span>
                        </label>
                        <input
                            type="text"
                            required
                            value={formData.name}
                            onChange={(e) => {
                                const name = e.target.value;
                                setFormData({
                                    ...formData,
                                    name,
                                    code: integrationType ? formData.code : name.toLowerCase().replace(/\s+/g, '_').replace(/[^a-z0-9_]/g, ''),
                                });
                            }}
                            placeholder="Ej: WooCommerce"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                    </div>

                    <div>
                        <label className={fieldLabel}>
                            Categoria <span style={{ color: ACCENT }}>*</span>
                        </label>
                        <Select
                            required
                            value={String(formData.category_id)}
                            onChange={(e) => setFormData({ ...formData, category_id: Number(e.target.value) })}
                            options={categories.map((cat) => ({ value: String(cat.id), label: cat.name }))}
                            className="bg-white dark:bg-gray-800"
                        />
                    </div>

                    <div className="md:col-span-2">
                        <label className={fieldLabel}>Logo del Tipo de Integracion</label>
                        <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
                            <img
                                src={imagePreview || ''}
                                alt="Logo"
                                className={`h-20 w-20 flex-shrink-0 rounded-xl object-contain p-2 ${imagePreview ? '' : 'hidden'}`}
                                style={{ border: `1px solid ${INPUT_BORDER}`, backgroundColor: '#ffffff' }}
                            />
                            <div className="flex-1">
                                <FileInput
                                    accept="image/*"
                                    onChange={handleImageChange}
                                    buttonText="Seleccionar imagen"
                                    helperText="Formatos soportados: JPG, PNG, GIF, WEBP. Tamano maximo: 10MB"
                                />
                                {imagePreview && (
                                    <p className={fieldHint}>
                                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                        <span>{imageFile ? 'Nueva imagen seleccionada' : 'Imagen actual'}</span>
                                    </p>
                                )}
                            </div>
                        </div>
                        {integrationType && integrationType.image_url && (
                            <label className="mt-3 flex items-center gap-2">
                                <input
                                    type="checkbox"
                                    checked={removeImage}
                                    onChange={(e) => {
                                        setRemoveImage(e.target.checked);
                                        if (e.target.checked) {
                                            setImageFile(null);
                                            setImagePreview(null);
                                        } else {
                                            setImagePreview(integrationType.image_url || null);
                                        }
                                    }}
                                    className="h-4 w-4 rounded border-gray-300 dark:border-gray-600 accent-[var(--color-primary)]"
                                />
                                <span className="text-[13px] text-gray-700 dark:text-gray-200">Eliminar imagen actual</span>
                            </label>
                        )}
                    </div>

                    <div className="md:col-span-2">
                        <label className={fieldLabel}>Descripcion</label>
                        <textarea
                            value={formData.description}
                            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                            rows={2}
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                    </div>
                </div>
            </SectionCard>

            <SectionCard icon={CodeBracketIcon} title="Esquemas de datos">
                <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>Config Schema (JSON)</label>
                        <textarea
                            value={formData.config_schema}
                            onChange={(e) => setFormData({ ...formData, config_schema: e.target.value })}
                            rows={12}
                            className={jsonCls}
                            placeholder='{"type": "object", "properties": {...}}'
                            spellCheck={false}
                        />
                        <p className={fieldHint}>
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>Campos de configuracion (no sensibles)</span>
                        </p>
                    </div>

                    <div>
                        <label className={fieldLabel}>Credentials Schema (JSON)</label>
                        <textarea
                            value={formData.credentials_schema}
                            onChange={(e) => setFormData({ ...formData, credentials_schema: e.target.value })}
                            rows={12}
                            className={jsonCls}
                            placeholder='{"type": "object", "properties": {...}}'
                            spellCheck={false}
                        />
                        <p className={fieldHint}>
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>Campos de credenciales (tokens, keys, etc.)</span>
                        </p>
                    </div>
                </div>
            </SectionCard>

            <SectionCard icon={GlobeAltIcon} title="URLs del API">
                <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>URL de Produccion</label>
                        <input
                            type="url"
                            value={formData.base_url}
                            onChange={(e) => setFormData({ ...formData, base_url: e.target.value })}
                            placeholder="https://api.ejemplo.com/v1"
                            className={`${inputCls} font-mono`}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>URL base del API en produccion</span>
                        </p>
                    </div>
                    <div>
                        <label className={fieldLabel}>URL de Pruebas (Sandbox)</label>
                        <input
                            type="url"
                            value={formData.base_url_test}
                            onChange={(e) => setFormData({ ...formData, base_url_test: e.target.value })}
                            placeholder="https://sandbox.ejemplo.com/v1"
                            className={`${inputCls} font-mono`}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>URL del entorno sandbox para modo de pruebas</span>
                        </p>
                    </div>
                </div>
            </SectionCard>

            <SectionCard icon={DocumentTextIcon} title="Instrucciones de Configuracion" bg="#ffffff">
                <textarea
                    value={formData.setup_instructions}
                    onChange={(e) => setFormData({ ...formData, setup_instructions: e.target.value })}
                    rows={6}
                    className={inputCls}
                    style={{ borderColor: INPUT_BORDER }}
                    placeholder="Pasos para configurar esta integracion:&#10;&#10;1. Ve a...&#10;2. Configura...&#10;3. Copia..."
                />
                <p className={fieldHint}>
                    <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                    <span>Instrucciones paso a paso para el usuario</span>
                </p>
            </SectionCard>

            {integrationType?.id === WHATSAPP_TYPE_ID ? (
                <WhatsAppTypeCredentialsForm
                    credentials={whatsappCredentials}
                    onChange={setWhatsappCredentials}
                    isEditing={!!integrationType}
                />
            ) : integrationType?.id === BOLD_TYPE_ID ? (
                <BoldTypeCredentialsForm
                    credentials={boldCredentials}
                    onChange={setBoldCredentials}
                    isEditing={!!integrationType}
                    webhookUrlProd={boldWebhookUrls.production}
                    webhookUrlTest={boldWebhookUrls.sandbox}
                />
            ) : (
                <SectionCard icon={KeyIcon} title="Credenciales de Plataforma">
                    <div className="flex items-center justify-end mb-2">
                        <button
                            type="button"
                            onClick={() => setShowPlatformCredentials((v) => !v)}
                            className="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200"
                        >
                            {showPlatformCredentials ? (
                                <>
                                    <EyeSlashIcon className="w-4 h-4" />
                                    Ocultar
                                </>
                            ) : (
                                <>
                                    <EyeIcon className="w-4 h-4" />
                                    Mostrar
                                </>
                            )}
                        </button>
                    </div>
                    <textarea
                        value={showPlatformCredentials
                            ? formData.platform_credentials
                            : formData.platform_credentials.replace(/:\s*"([^"]*)"/g, ': "********"')
                        }
                        onChange={(e) => {
                            if (showPlatformCredentials) {
                                setFormData({ ...formData, platform_credentials: e.target.value });
                            }
                        }}
                        readOnly={!showPlatformCredentials}
                        rows={6}
                        className={`${jsonCls} ${showPlatformCredentials ? '' : 'text-gray-500 cursor-default'}`}
                        placeholder={showPlatformCredentials ? '{\n  "api_key": "tu-api-key-aqui"\n}' : ''}
                        spellCheck={false}
                    />
                    <p className={fieldHint}>
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>Credenciales globales del proveedor (se encriptaran). Deja <code>{'{}'}</code> para no cambiarlas.</span>
                    </p>
                </SectionCard>
            )}

            <div
                className="flex items-start gap-4 rounded-xl p-4 dark:bg-gray-800/60"
                style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
            >
                <button
                    type="button"
                    role="switch"
                    aria-checked={formData.is_active}
                    onClick={() => setFormData({ ...formData, is_active: !formData.is_active })}
                    className="relative inline-flex h-7 w-12 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 mt-0.5"
                    style={{ backgroundColor: formData.is_active ? ACCENT : '#d1d5db' }}
                >
                    <span
                        className={`pointer-events-none inline-block h-6 w-6 transform rounded-full bg-white shadow ring-0 transition duration-200 ${
                            formData.is_active ? 'translate-x-5' : 'translate-x-0'
                        }`}
                    />
                </button>
                <div className="flex-1">
                    <span className="block text-base font-semibold" style={{ color: formData.is_active ? ACCENT_DARK : '#6b7280' }}>
                        {formData.is_active ? 'Activo' : 'Desactivado'}
                    </span>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mt-0.5">
                        {formData.is_active ? 'Este tipo de integracion esta disponible para los negocios.' : 'Este tipo de integracion esta oculto para los negocios.'}
                    </p>
                </div>
            </div>

            <div className="flex justify-end space-x-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                {onCancel && (
                    <Button type="button" onClick={onCancel} variant="outline">
                        Cancelar
                    </Button>
                )}
                <Button type="submit" disabled={loading} loading={loading} variant="primary">
                    {integrationType ? 'Actualizar' : 'Crear'}
                </Button>
            </div>
        </form>
    );
}
