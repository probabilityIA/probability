import Link from 'next/link';
import { HeroContent, PublicBusiness } from '../../domain/types';

interface HeroSectionProps {
    content: HeroContent | null;
    business: PublicBusiness;
    slug: string;
}

export function HeroSection({ content, business, slug }: HeroSectionProps) {
    const title = content?.title || `Bienvenido a ${business.name}`;
    const subtitle = content?.subtitle || business.description || 'Descubre nuestros productos';
    const ctaText = content?.cta_text || 'Ver Productos';
    const bgImage = content?.background_image;

    return (
        <section
            className="relative py-24 px-4 flex items-center justify-center text-center min-h-[400px]"
            style={{
                backgroundColor: bgImage ? undefined : 'var(--brand-primary)',
                backgroundImage: bgImage ? `url(${bgImage})` : undefined,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
            }}
        >
            {bgImage && <div className="absolute inset-0 bg-black/50" />}
            <div className="relative z-10 max-w-3xl mx-auto">
                <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">{title}</h1>
                <p className="text-xl text-white/80 mb-8">{subtitle}</p>
                <Link
                    href={`/tienda/${slug}/productos`}
                    className="inline-block px-8 py-4 bg-white rounded-lg font-bold text-lg transition-transform hover:scale-105"
                    style={{ color: 'var(--brand-primary)' }}
                >
                    {ctaText}
                </Link>
            </div>
        </section>
    );
}
