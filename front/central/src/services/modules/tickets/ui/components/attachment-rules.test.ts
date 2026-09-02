import { describe, it, expect } from 'vitest';
import {
    MAX_UPLOAD_BYTES,
    ACCEPTED_TYPES,
    isImageMime,
    isPdfMime,
    isAllowedFile,
    formatSize,
    attachmentError,
} from './attachment-rules';

const makeFile = (name: string, type: string, size: number): File => {
    const file = new File(['x'], name, { type });
    Object.defineProperty(file, 'size', { value: size });
    return file;
};

describe('MAX_UPLOAD_BYTES', () => {
    it('vale exactamente 10 MB en bytes', () => {
        expect(MAX_UPLOAD_BYTES).toBe(10485760);
    });
});

describe('ACCEPTED_TYPES', () => {
    it('declara los cinco mime admitidos', () => {
        expect(ACCEPTED_TYPES.split(',')).toEqual([
            'image/png',
            'image/jpeg',
            'image/webp',
            'image/gif',
            'application/pdf',
        ]);
    });
});

describe('isImageMime', () => {
    it('acepta cualquier mime que empiece por image/', () => {
        expect(isImageMime('image/png')).toBe(true);
        expect(isImageMime('image/jpeg')).toBe(true);
        expect(isImageMime('image/svg+xml')).toBe(true);
    });

    it('rechaza mime vacio', () => {
        expect(isImageMime('')).toBe(false);
    });

    it('rechaza mime que no sea de imagen', () => {
        expect(isImageMime('application/pdf')).toBe(false);
        expect(isImageMime('text/image')).toBe(false);
    });
});

describe('isPdfMime', () => {
    it('solo acepta application/pdf exacto', () => {
        expect(isPdfMime('application/pdf')).toBe(true);
        expect(isPdfMime('application/x-pdf')).toBe(false);
        expect(isPdfMime('')).toBe(false);
    });
});

describe('isAllowedFile', () => {
    it('acepta imagenes', () => {
        expect(isAllowedFile(makeFile('foto.png', 'image/png', 100))).toBe(true);
    });

    it('acepta PDF', () => {
        expect(isAllowedFile(makeFile('doc.pdf', 'application/pdf', 100))).toBe(true);
    });

    it('rechaza otros tipos', () => {
        expect(isAllowedFile(makeFile('hoja.csv', 'text/csv', 100))).toBe(false);
        expect(isAllowedFile(makeFile('app.zip', 'application/zip', 100))).toBe(false);
    });

    it('rechaza archivo sin mime', () => {
        expect(isAllowedFile(makeFile('raro', '', 100))).toBe(false);
    });
});

describe('formatSize', () => {
    it('usa bytes por debajo de 1 KB', () => {
        expect(formatSize(0)).toBe('0 B');
        expect(formatSize(1)).toBe('1 B');
        expect(formatSize(1023)).toBe('1023 B');
    });

    it('usa KB desde 1024 bytes y hasta 1 MB exclusivo', () => {
        expect(formatSize(1024)).toBe('1.0 KB');
        expect(formatSize(1536)).toBe('1.5 KB');
        expect(formatSize(1048575)).toBe('1024.0 KB');
    });

    it('usa MB desde 1 MB', () => {
        expect(formatSize(1048576)).toBe('1.0 MB');
        expect(formatSize(MAX_UPLOAD_BYTES)).toBe('10.0 MB');
        expect(formatSize(15 * 1024 * 1024)).toBe('15.0 MB');
    });
});

describe('attachmentError', () => {
    it('no reporta error para imagen dentro del limite', () => {
        expect(attachmentError(makeFile('foto.png', 'image/png', 1024))).toBeNull();
    });

    it('no reporta error para PDF dentro del limite', () => {
        expect(attachmentError(makeFile('doc.pdf', 'application/pdf', 2048))).toBeNull();
    });

    it('acepta un archivo de exactamente 10 MB', () => {
        expect(attachmentError(makeFile('justo.png', 'image/png', MAX_UPLOAD_BYTES))).toBeNull();
    });

    it('rechaza un archivo de 10 MB mas un byte con el peso formateado', () => {
        const error = attachmentError(makeFile('grande.png', 'image/png', MAX_UPLOAD_BYTES + 1));

        expect(error).toBe(
            'No se puede adjuntar grande.png. Pesa 10.0 MB y el l\u00edmite es 10 MB.'
        );
    });

    it('reporta el peso real del archivo rechazado', () => {
        const error = attachmentError(makeFile('enorme.pdf', 'application/pdf', 25 * 1024 * 1024));

        expect(error).toContain('Pesa 25.0 MB');
        expect(error).toContain('enorme.pdf');
    });

    it('rechaza un tipo no permitido antes que el tamano', () => {
        const error = attachmentError(makeFile('datos.csv', 'text/csv', MAX_UPLOAD_BYTES * 3));

        expect(error).toBe(
            'No se puede adjuntar datos.csv. Solo se admiten im\u00e1genes (PNG, JPG, WEBP, GIF) y PDF.'
        );
    });

    it('rechaza un archivo sin mime', () => {
        const error = attachmentError(makeFile('sinmime', '', 10));

        expect(error).toContain('Solo se admiten im\u00e1genes');
    });
});
