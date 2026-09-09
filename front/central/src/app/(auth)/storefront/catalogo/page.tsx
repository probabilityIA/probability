import Link from 'next/link';
import { UsersIcon, ClipboardDocumentListIcon } from '@heroicons/react/24/outline';
import { getCatalogAction, getCatalogFiltersAction } from '@/services/modules/storefront/infra/actions';
import { getStorefrontBusinessId } from '@/shared/utils/storefront-business';
import { CatalogGrid } from '@/services/modules/storefront/ui/components/CatalogGrid';
import { CatalogBannerDisplay } from '@/services/modules/storefront/ui/components/CatalogBannerDisplay';
import { CatalogSearch } from './search';
import { CatalogFilters } from './filters';
import { CatalogLayoutSettings } from './layout-settings';
import { BannerSettings } from './banner-settings';
import { StorefrontPagination } from './pagination';

interface PageProps {
    searchParams: Promise<{ search?: string; category?: string; family_id?: string; page?: string }>;
}

export default async function CatalogoPage({ searchParams }: PageProps) {
    const params = await searchParams;
    const businessId = await getStorefrontBusinessId();
    const page = params.page ? parseInt(params.page) : 1;
    const search = params.search || '';
    const category = params.category || '';
    const familyId = params.family_id ? parseInt(params.family_id) : undefined;

    const filters = await getCatalogFiltersAction(businessId);
    const pageSize = filters.layout.columns * filters.layout.rows;

    const data = await getCatalogAction({
        page,
        page_size: pageSize,
        search: search || undefined,
        category: category || undefined,
        family_id: familyId,
        business_id: businessId,
    });

    return (
        <div className="p-6 sm:p-8">
            <div className="flex flex-wrap items-center justify-between gap-3 mb-6">
                <div className="flex flex-wrap items-center gap-2">
                    <CatalogSearch initialSearch={search} />
                    <CatalogFilters
                        categories={filters.categories}
                        families={filters.families}
                        selectedCategory={category}
                        selectedFamilyId={familyId ?? null}
                    />
                </div>
                <div className="flex gap-2">
                    <BannerSettings banner={filters.banner} businessId={businessId} />
                    <CatalogLayoutSettings layout={filters.layout} businessId={businessId} />
                    <Link
                        href="/storefront/clientes"
                        className="flex items-center gap-2 px-4 py-2 h-[42px] text-sm font-medium border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                    >
                        <UsersIcon className="w-4 h-4" />
                        Clientes
                    </Link>
                    <Link
                        href="/storefront/pedidos"
                        className="flex items-center gap-2 px-4 py-2 h-[42px] text-sm font-medium border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                    >
                        <ClipboardDocumentListIcon className="w-4 h-4" />
                        Pedidos
                    </Link>
                </div>
            </div>

            <CatalogBannerDisplay banner={filters.banner} />

            <CatalogGrid products={data.data} columns={filters.layout.columns} />

            {data.total_pages > 1 && (
                <StorefrontPagination
                    currentPage={data.page}
                    totalPages={data.total_pages}
                    total={data.total}
                    basePath="/storefront/catalogo"
                    label="productos"
                />
            )}
        </div>
    );
}
