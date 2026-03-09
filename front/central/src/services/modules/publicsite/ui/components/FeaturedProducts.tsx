import { PublicProduct } from '../../domain/types';
import { PublicProductCard } from './PublicProductCard';

interface FeaturedProductsProps {
    products: PublicProduct[];
    slug: string;
}

export function FeaturedProducts({ products, slug }: FeaturedProductsProps) {
    return (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {products.map((product) => (
                <PublicProductCard key={product.id} product={product} slug={slug} />
            ))}
        </div>
    );
}
