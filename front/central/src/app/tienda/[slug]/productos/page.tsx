import { notFound } from 'next/navigation';
import { getPublicBusinessAction, getPublicCatalogAction } from '@/services/modules/publicsite/infra/actions';
import { getTemplate } from '@/services/modules/publicsite/ui/templates/registry';
import { CatalogSearch } from '@/services/modules/publicsite/ui/components/CatalogSearch';
import { CatalogPagination } from '@/services/modules/publicsite/ui/components/CatalogPagination';

interface PageProps {
    params: Promise<{ slug: string }>;
    searchParams: Promise<{ page?: string; search?: string; category?: string }>;
}

export default async function ProductosPage({ params, searchParams }: PageProps) {
    const { slug } = await params;
    const sp = await searchParams;

    const business = await getPublicBusinessAction(slug);
    if (!business) return notFound();

    const template = getTemplate(business.website_config?.template || 'default');
    const ProductCard = template.ProductCard;

    const page = sp.page ? parseInt(sp.page) : 1;
    const search = sp.search || '';
    const category = sp.category || '';

    const catalog = await getPublicCatalogAction(slug, {
        page,
        page_size: 12,
        search: search || undefined,
        category: category || undefined,
    });

    return (
        <div className="py-8 px-4 max-w-7xl mx-auto">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-8">Catálogo de Productos</h1>

            <CatalogSearch basePath={`/tienda/${slug}/productos`} initialSearch={search} />

            {catalog.data.length === 0 ? (
                <div className="text-center py-16 text-gray-500 dark:text-gray-400">
                    <p className="text-lg">No se encontraron productos</p>
                    {search && <p className="mt-2">Intenta con otro término de búsqueda</p>}
                </div>
            ) : (
                <>
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 mt-6">
                        {catalog.data.map((product) => (
                            <ProductCard key={product.id} product={product} slug={slug} />
                        ))}
                    </div>

                    {catalog.total_pages > 1 && (
                        <CatalogPagination
                            currentPage={catalog.page}
                            totalPages={catalog.total_pages}
                            total={catalog.total}
                            basePath={`/tienda/${slug}/productos`}
                        />
                    )}
                </>
            )}
        </div>
    );
}
