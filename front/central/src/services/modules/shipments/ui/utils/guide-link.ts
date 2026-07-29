const MELI_LABEL_MARKERS = ['api.mercadolibre.com', 'shipment_labels'];

export function isMeliLabelUrl(guideUrl?: string | null): boolean {
    if (!guideUrl) return false;
    const value = guideUrl.toLowerCase();
    return MELI_LABEL_MARKERS.some((marker) => value.includes(marker));
}

export function guideHref(
    shipmentId: number | string,
    guideUrl?: string | null,
    businessId?: number | null,
): string {
    if (!guideUrl) return '';
    if (!isMeliLabelUrl(guideUrl)) return guideUrl;
    const qs = businessId ? `?business_id=${businessId}` : '';
    return `/internal/meli-label/${shipmentId}${qs}`;
}
