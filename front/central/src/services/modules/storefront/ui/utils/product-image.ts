const ALLOWED_IMAGE_HOST = /\.s3([.-][a-z0-9-]+)?\.amazonaws\.com$/i;

export function isOptimizableProductImage(url: string): boolean {
    if (!url) return false;
    try {
        const { protocol, hostname } = new URL(url);
        return protocol === 'https:' && ALLOWED_IMAGE_HOST.test(hostname);
    } catch {
        return false;
    }
}
