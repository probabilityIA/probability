export const GUIDE_FIELD_LIMITS = {
    address: 50,
    crossStreet: 50,
    reference: 25,
    suburb: 30,
} as const;

export const ADDRESS_PART_SEPARATOR = ' | ';

export interface GuideAddressSource {
    shipping_street?: string | null;
    shipping_neighborhood?: string | null;
}

export interface GuideDestinationParts {
    address: string;
    crossStreet: string;
    reference: string;
    suburb: string;
    dropped: string;
}

const collapse = (text: string): string => text.trim().split(/\s+/).filter(Boolean).join(' ');

export const clampGuideField = (text: string, limit: number): string => {
    const clean = collapse(text || '');
    if (clean.length <= limit) return clean;
    const cut = clean.slice(0, limit);
    const lastSpace = cut.lastIndexOf(' ');
    const trimmed = lastSpace >= Math.floor(limit / 2) ? cut.slice(0, lastSpace) : cut;
    return trimmed.replace(/[\s,.\-]+$/, '');
};

const fill = (chunks: string[], limit: number): { value: string; rest: string[] } => {
    let value = '';
    let i = 0;
    while (i < chunks.length) {
        const candidate = value ? `${value} ${chunks[i]}` : chunks[i];
        if (collapse(candidate).length > limit) break;
        value = candidate;
        i += 1;
    }
    if (!value && chunks.length > 0) {
        return { value: clampGuideField(chunks[0], limit), rest: chunks.slice(1) };
    }
    return { value: collapse(value), rest: chunks.slice(i) };
};

const COMPLEMENT_MARKERS = [
    'conjunto', 'conjuntos', 'urbanizacion', 'unidad', 'parcelacion', 'edificio', 'edif',
    'torre', 'bloque', 'manzana', 'mz', 'etapa', 'casa', 'apto', 'apartamento', 'apt',
    'interior', 'int', 'piso', 'local', 'oficina', 'barrio', 'sector', 'vereda', 'agrupacion',
];

const deaccent = (text: string): string => text.normalize('NFD').replace(/[\u0300-\u036f]/g, '').toLowerCase();

const findComplementCut = (text: string): number => {
    const words = text.split(' ');
    let offset = 0;
    for (let i = 0; i < words.length; i += 1) {
        const bare = deaccent(words[i]).replace(/[^a-z/]/g, '');
        if (i > 0 && (COMPLEMENT_MARKERS.includes(bare) || bare === 'b/')) return offset;
        offset += words[i].length + 1;
    }
    return -1;
};

const splitLongText = (text: string, limit: number): string[] => {
    const words = collapse(text).split(' ');
    const chunks: string[] = [];
    let current = '';
    for (const word of words) {
        const candidate = current ? `${current} ${word}` : word;
        if (candidate.length > limit && current) {
            chunks.push(current);
            current = word;
        } else {
            current = candidate;
        }
    }
    if (current) chunks.push(current);
    return chunks;
};

export const buildGuideDestination = (order: GuideAddressSource | null | undefined): GuideDestinationParts => {
    const street = collapse(order?.shipping_street || '');
    const neighborhood = collapse(order?.shipping_neighborhood || '');

    if (!street) {
        return {
            address: '',
            crossStreet: '',
            reference: '',
            suburb: clampGuideField(neighborhood, GUIDE_FIELD_LIMITS.suburb),
            dropped: '',
        };
    }

    const parts = street.split(ADDRESS_PART_SEPARATOR).map(collapse).filter(Boolean);

    if (parts.length > 1) {
        const address = clampGuideField(parts[0], GUIDE_FIELD_LIMITS.address);
        const suburbSource = parts[2] || neighborhood;
        const suburb = clampGuideField(suburbSource, GUIDE_FIELD_LIMITS.suburb);
        const reference = clampGuideField(parts[1], GUIDE_FIELD_LIMITS.reference);
        const complement = parts.slice(1).join(' ');
        const crossStreet = clampGuideField(complement, GUIDE_FIELD_LIMITS.crossStreet);
        return { address, crossStreet, reference, suburb, dropped: '' };
    }

    const cut = findComplementCut(street);
    if (cut > 0 && cut <= GUIDE_FIELD_LIMITS.address) {
        const via = collapse(street.slice(0, cut));
        const complement = collapse(street.slice(cut));
        const complementChunks = splitLongText(complement, GUIDE_FIELD_LIMITS.crossStreet);
        const crossFromCut = fill(complementChunks, GUIDE_FIELD_LIMITS.crossStreet);
        const referenceFromCut = fill(crossFromCut.rest, GUIDE_FIELD_LIMITS.reference);
        return {
            address: via,
            crossStreet: crossFromCut.value,
            reference: referenceFromCut.value,
            suburb: clampGuideField(neighborhood, GUIDE_FIELD_LIMITS.suburb),
            dropped: collapse(referenceFromCut.rest.join(' ')),
        };
    }

    const chunks = splitLongText(street, GUIDE_FIELD_LIMITS.address);
    const addressSlot = fill(chunks, GUIDE_FIELD_LIMITS.address);
    const crossSlot = fill(addressSlot.rest, GUIDE_FIELD_LIMITS.crossStreet);
    const referenceSlot = fill(crossSlot.rest, GUIDE_FIELD_LIMITS.reference);

    return {
        address: addressSlot.value,
        crossStreet: crossSlot.value,
        reference: referenceSlot.value,
        suburb: clampGuideField(neighborhood, GUIDE_FIELD_LIMITS.suburb),
        dropped: collapse(referenceSlot.rest.join(' ')),
    };
};
