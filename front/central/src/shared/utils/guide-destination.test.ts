import { describe, it, expect } from 'vitest';
import { buildGuideDestination, buildComplement, clampGuideField, GUIDE_FIELD_LIMITS } from './guide-destination';

const within = (parts: ReturnType<typeof buildGuideDestination>) => {
    expect(parts.address.length).toBeLessThanOrEqual(GUIDE_FIELD_LIMITS.address);
    expect(parts.crossStreet.length).toBeLessThanOrEqual(GUIDE_FIELD_LIMITS.crossStreet);
    expect(parts.reference.length).toBeLessThanOrEqual(GUIDE_FIELD_LIMITS.reference);
    expect(parts.suburb.length).toBeLessThanOrEqual(GUIDE_FIELD_LIMITS.suburb);
};

describe('buildGuideDestination con el formato de tres campos', () => {
    it('MYS-0922 mantiene el reparto que ya funcionaba', () => {
        const parts = buildGuideDestination({
            shipping_street: 'Urbanización el Lago   | Casa 58,  La troncal. | Caucasia - Antioquia',
        });
        expect(parts.address).toBe('Urbanización el Lago');
        expect(parts.reference).toBe('Casa 58, La troncal.');
        expect(parts.suburb).toBe('Caucasia - Antioquia');
        within(parts);
    });

    it('el complemento completo viaja en crossStreet en vez de duplicar la direccion', () => {
        const parts = buildGuideDestination({
            shipping_street: 'Calle 146 #7b-80  | edificio siena | apto 704 Belmira',
        });
        expect(parts.address).toBe('Calle 146 #7b-80');
        expect(parts.crossStreet).toBe('edificio siena apto 704 Belmira');
        expect(parts.crossStreet).not.toBe(parts.address);
        within(parts);
    });
});

describe('buildGuideDestination con una sola cadena', () => {
    it('VIG-0150 ya no pierde la torre ni el apartamento', () => {
        const parts = buildGuideDestination({
            shipping_street: 'Calle 6 # 3-184 Conjunto residencial Alameda de Albornoz Torre B Apartamento 503',
        });
        const enviado = [parts.address, parts.crossStreet, parts.reference].join(' ');
        expect(enviado).toContain('Torre B');
        expect(enviado).toContain('Apartamento 503');
        expect(parts.dropped).toBe('');
        within(parts);
    });

    it('VIG-0149 ya no pierde el bloque ni el apartamento', () => {
        const parts = buildGuideDestination({
            shipping_street: 'Calle 64B#31-15 Conjuntos conucos entrada 4 ,bloque 25 apto 4-02',
        });
        const enviado = [parts.address, parts.crossStreet, parts.reference].join(' ');
        expect(enviado).toContain('bloque 25');
        expect(enviado).toContain('apto 4-02');
        expect(parts.dropped).toBe('');
        within(parts);
    });

    it('VIG-0154 reparte sin perder el sector', () => {
        const parts = buildGuideDestination({
            shipping_street: 'CRA 10 # 112-42 barrio la paz ( sector las mangueitas)',
        });
        const enviado = [parts.address, parts.crossStreet, parts.reference].join(' ');
        expect(enviado).toContain('barrio la paz');
        expect(enviado).toContain('mangueitas');
        within(parts);
    });

    it('corta en el complemento, no a mitad de frase', () => {
        const parts = buildGuideDestination({
            shipping_street: 'Calle 6 # 3-184 Conjunto residencial Alameda de Albornoz Torre B Apartamento 503',
        });
        expect(parts.address).toBe('Calle 6 # 3-184');
        expect(parts.crossStreet).toBe('Conjunto residencial Alameda de Albornoz Torre B');
        expect(parts.reference).toBe('Apartamento 503');
    });

    it('reconoce el marcador barrio', () => {
        const parts = buildGuideDestination({ shipping_street: 'CRA 1D #25-25 barrio bello horizonte' });
        expect(parts.address).toBe('CRA 1D #25-25');
        expect(parts.crossStreet).toBe('barrio bello horizonte');
    });

    it('sin marcador de complemento reparte por palabras sin perder nada', () => {
        const parts = buildGuideDestination({
            shipping_street: 'Avenida Siempre Viva numero 742 esquina con la quinta transversal larga',
        });
        const enviado = [parts.address, parts.crossStreet, parts.reference].join(' ');
        expect(enviado).toContain('transversal larga');
        expect(parts.dropped).toBe('');
        within(parts);
    });

    it('una direccion corta cabe entera en address y no inventa complemento', () => {
        const parts = buildGuideDestination({ shipping_street: 'cra 109  #17 - 24 Fontibon' });
        expect(parts.address).toBe('cra 109 #17 - 24 Fontibon');
        expect(parts.crossStreet).toBe('');
        expect(parts.reference).toBe('');
        within(parts);
    });

    it('el barrio sale de shipping_neighborhood, nunca del departamento', () => {
        const parts = buildGuideDestination({
            shipping_street: 'Calle 8 #36-16',
            shipping_neighborhood: 'La Troncal',
        });
        expect(parts.suburb).toBe('La Troncal');
    });

    it('sin barrio conocido el campo queda vacio', () => {
        const parts = buildGuideDestination({ shipping_street: 'Calle 8 #36-16' });
        expect(parts.suburb).toBe('');
    });
});

describe('con los campos estructurados del formulario', () => {
    it('VIG-0150 capturada con campos separados no necesita heuristica', () => {
        const parts = buildGuideDestination({
            shipping_street: 'Calle 6 # 3-184',
            shipping_complement_type: 'apartamento',
            shipping_complement_number: '503',
            shipping_tower: 'Torre B',
            shipping_building: 'Alameda de Albornoz',
            shipping_neighborhood: 'Albornoz',
        });
        expect(parts.address).toBe('Calle 6 # 3-184');
        expect(parts.reference).toBe('Apto 503');
        expect(parts.crossStreet).toBe('Torre B Alameda de Albornoz');
        expect(parts.suburb).toBe('Albornoz');
        expect(parts.dropped).toBe('');
        within(parts);
    });

    it('los campos estructurados le ganan al parseo del texto', () => {
        const parts = buildGuideDestination({
            shipping_street: 'Calle 6 # 3-184 | Apto 503 Torre B | Albornoz',
            shipping_complement_type: 'casa',
            shipping_complement_number: '17',
        });
        expect(parts.reference).toBe('Casa 17');
        expect(parts.suburb).toBe('Albornoz');
    });

    it('buildComplement arma el texto que ve el transportador', () => {
        expect(buildComplement('apartamento', '503', 'Torre B', 'Alameda')).toBe('Apto 503 Torre B Alameda');
        expect(buildComplement('casa', '17', '', '')).toBe('Casa 17');
        expect(buildComplement('', '', '', '')).toBe('');
        expect(buildComplement('', '  ', 'Torre 2', '')).toBe('Torre 2');
    });
});

describe('casos limite', () => {
    it('una orden sin direccion no revienta', () => {
        const parts = buildGuideDestination({ shipping_street: '' });
        expect(parts.address).toBe('');
        expect(parts.dropped).toBe('');
    });

    it('null tampoco revienta', () => {
        expect(buildGuideDestination(null).address).toBe('');
    });

    it('reporta en dropped lo que no cabe en ningun campo', () => {
        const larga = Array.from({ length: 40 }, (_, i) => `palabra${i}`).join(' ');
        const parts = buildGuideDestination({ shipping_street: larga });
        expect(parts.dropped.length).toBeGreaterThan(0);
        within(parts);
    });

    it('clampGuideField corta por palabra y limpia la puntuacion del final', () => {
        expect(clampGuideField('Calle 6 # 3-184 Conjunto residencial,', 20)).toBe('Calle 6 # 3-184');
        expect(clampGuideField('  doble   espacio  ', 50)).toBe('doble espacio');
    });
});
