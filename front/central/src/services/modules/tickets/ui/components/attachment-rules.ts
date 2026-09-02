export const MAX_UPLOAD_BYTES = 10 * 1024 * 1024;
export const ACCEPTED_TYPES = 'image/png,image/jpeg,image/webp,image/gif,application/pdf';

export const isImageMime = (mime: string) => !!mime && mime.startsWith('image/');
export const isPdfMime = (mime: string) => mime === 'application/pdf';
export const isAllowedFile = (file: File) => isImageMime(file.type) || isPdfMime(file.type);

export const formatSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
};

export const attachmentError = (file: File): string | null => {
    if (!isAllowedFile(file)) {
        return 'No se puede adjuntar ' + file.name + '. Solo se admiten im\u00e1genes (PNG, JPG, WEBP, GIF) y PDF.';
    }
    if (file.size > MAX_UPLOAD_BYTES) {
        return 'No se puede adjuntar ' + file.name + '. Pesa ' + formatSize(file.size) + ' y el l\u00edmite es 10 MB.';
    }
    return null;
};
