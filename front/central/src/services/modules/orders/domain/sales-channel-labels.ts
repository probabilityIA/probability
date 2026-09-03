const CHANNEL_LABELS: Record<string, string> = {
    shopify: 'Shopify',
    woocommerce: 'WooCommerce',
    mercadolibre: 'MercadoLibre',
    amazon: 'Amazon',
    jumpseller: 'Jumpseller',
    prestashop: 'PrestaShop',
    magento: 'Magento',
    vtex: 'VTEX',
    tiendanube: 'Tiendanube',
    shein: 'SHEIN',
    tiktok: 'TikTok Shop',
    whatsapp: 'WhatsApp',
    siigo: 'Siigo',
    softpymes: 'Softpymes',
    manual: 'Manual',
    api: 'API',
};

const CANALES_PROPIOS = ['manual', 'plataforma', 'platform', 'api', ''];

export function esCanalDeVenta(platform?: string | null): boolean {
    if (!platform) return false;
    const key = platform.trim().toLowerCase().replace(/[\s_-]+/g, '');
    return !CANALES_PROPIOS.includes(key);
}

export function salesChannelLabel(platform?: string | null): string {
    if (!platform) return '-';
    const key = platform.trim().toLowerCase().replace(/[\s_-]+/g, '');
    if (CHANNEL_LABELS[key]) return CHANNEL_LABELS[key];
    return platform.charAt(0).toUpperCase() + platform.slice(1);
}
