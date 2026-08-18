const MELI_LABEL_MARKERS = ['api.mercadolibre.com', 'shipment_labels'];

export function isMeliLabelUrl(guideUrl?: string | null): boolean {
    if (!guideUrl) return false;
    const value = guideUrl.toLowerCase();
    return MELI_LABEL_MARKERS.some((marker) => value.includes(marker));
}

// La etiqueta de MercadoLibre vive en su API y exige el token del vendedor:
// abrirla directo devuelve 401. Se redirige al proxy interno, que la pide
// autenticada y devuelve el PDF.
export function guideHref(
    shipmentId: number | string | null | undefined,
    guideUrl?: string | null,
    businessId?: number | null,
): string {
    if (!guideUrl) return '';
    if (!isMeliLabelUrl(guideUrl)) return guideUrl;
    if (!shipmentId) return guideUrl;
    const qs = businessId ? `?business_id=${businessId}` : '';
    return `/internal/meli-label/${shipmentId}${qs}`;
}
