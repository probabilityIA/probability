import { getCatalogAction } from '@/services/modules/storefront/infra/actions';
import { getStorefrontBusinessId } from '@/shared/utils/storefront-business';
import { CatalogGrid } from '@/services/modules/storefront/ui/components/CatalogGrid';
import { CatalogSearch } from './search';
import { StorefrontPagination } from './pagination';

interface PageProps {
    searchParams: Promise<{ search?: string; page?: string }>;
}

export default async function CatalogoPage({ searchParams }: PageProps) {
    const params = await searchParams;
    const businessId = await getStorefrontBusinessId();
    const page = params.page ? parseInt(params.page) : 1;
    const search = params.search || '';

    const data = await getCatalogAction({
        page,
        page_size: 12,
        search: search || undefined,
        business_id: businessId,
    });

    return (
        <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Catalogo</h1>
            <CatalogSearch initialSearch={search} />

            {data.data.length === 0 ? (
                <div className="text-center py-12">
                    <p className="text-gray-500 dark:text-gray-400 text-lg">No se encontraron productos</p>
                </div>
            ) : (
                <>
                    <CatalogGrid products={data.data} />
                    {data.total_pages > 1 && (
                        <StorefrontPagination
                            currentPage={data.page}
                            totalPages={data.total_pages}
                            total={data.total}
                            basePath="/storefront/catalogo"
                            label="productos"
                        />
                    )}
                </>
            )}
        </div>
    );
}
