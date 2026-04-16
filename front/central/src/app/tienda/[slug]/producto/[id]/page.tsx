import { notFound } from 'next/navigation';
import { getPublicBusinessAction, getPublicProductAction } from '@/services/modules/publicsite/infra/actions';
import Link from 'next/link';

interface PageProps {
    params: Promise<{ slug: string; id: string }>;
}

export default async function ProductoDetailPage({ params }: PageProps) {
    const { slug, id } = await params;

    const [business, product] = await Promise.all([
        getPublicBusinessAction(slug),
        getPublicProductAction(slug, id),
    ]);

    if (!business || !product) return notFound();

    return (
        <div className="py-8 px-4 max-w-7xl mx-auto">
            <Link
                href={`/tienda/${slug}/productos`}
                className="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:text-gray-200 mb-6 inline-flex items-center gap-1"
            >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                </svg>
                Volver al catálogo
            </Link>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mt-4">
                {/* Product Image */}
                <div className="aspect-square bg-gray-100 rounded-xl overflow-hidden">
                    {product.image_url ? (
                        <img
                            src={product.image_url}
                            alt={product.name}
                            className="w-full h-full object-cover"
                        />
                    ) : (
                        <div className="w-full h-full flex items-center justify-center text-gray-400">
                            <svg className="w-24 h-24" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                            </svg>
                        </div>
                    )}
                </div>

                {/* Product Info */}
                <div>
                    {product.brand && (
                        <p className="text-sm text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">{product.brand}</p>
                    )}
                    <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">{product.name}</h1>

                    <div className="flex items-baseline gap-3 mb-6">
                        <span className="text-3xl font-bold" style={{ color: 'var(--brand-secondary)' }}>
                            ${product.price.toLocaleString('es-CO')} {product.currency}
                        </span>
                        {product.compare_at_price && product.compare_at_price > product.price && (
                            <span className="text-lg text-gray-400 line-through">
                                ${product.compare_at_price.toLocaleString('es-CO')}
                            </span>
                        )}
                    </div>

                    {product.category && (
                        <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">
                            Categoría: <span className="font-medium text-gray-700 dark:text-gray-200">{product.category}</span>
                        </p>
                    )}

                    {product.description && (
                        <div className="prose prose-gray max-w-none mb-8">
                            <p className="text-gray-600 dark:text-gray-300 whitespace-pre-line">{product.description}</p>
                        </div>
                    )}

                    <div className="flex items-center gap-2 mb-8">
                        {product.stock_quantity > 0 ? (
                            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium bg-green-50 text-green-700">
                                <span className="w-2 h-2 rounded-full bg-green-500"></span>
                                En stock
                            </span>
                        ) : (
                            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium bg-red-50 text-red-700">
                                <span className="w-2 h-2 rounded-full bg-red-500"></span>
                                Agotado
                            </span>
                        )}
                    </div>

                    <Link
                        href={`/login?redirect=/storefront/catalogo&business_code=${slug}`}
                        className="inline-block w-full text-center px-8 py-4 rounded-lg text-white font-bold text-lg transition-transform hover:scale-[1.02]"
                        style={{ backgroundColor: 'var(--brand-secondary)' }}
                    >
                        Hacer Pedido
                    </Link>
                </div>
            </div>
        </div>
    );
}
