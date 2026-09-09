import { CatalogBanner } from '../../domain/types';

interface CatalogBannerDisplayProps {
    banner: CatalogBanner;
}

export function CatalogBannerDisplay({ banner }: CatalogBannerDisplayProps) {
    if (!banner.enabled || !banner.image_url) {
        return null;
    }

    return (
        <div className="w-full aspect-[16/5] sm:aspect-[21/5] rounded-2xl overflow-hidden border border-gray-200 dark:border-gray-700 mb-6 bg-gray-50 dark:bg-gray-900">
            <img src={banner.image_url} alt="Banner del catalogo" className="w-full h-full object-cover" />
        </div>
    );
}
