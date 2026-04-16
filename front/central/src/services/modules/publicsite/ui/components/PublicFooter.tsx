import Link from 'next/link';
import { PublicBusiness } from '../../domain/types';

interface PublicFooterProps {
    business: PublicBusiness;
}

export function PublicFooter({ business }: PublicFooterProps) {
    const slug = business.code;

    return (
        <footer className="bg-gray-900 text-white py-12">
            <div className="max-w-7xl mx-auto px-4">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    <div>
                        {business.logo_url ? (
                            <img src={business.logo_url} alt={business.name} className="h-10 object-contain mb-4 brightness-0 invert" />
                        ) : (
                            <h3 className="text-xl font-bold mb-4">{business.name}</h3>
                        )}
                        {business.description && (
                            <p className="text-gray-400 text-sm">{business.description}</p>
                        )}
                    </div>
                    <div>
                        <h4 className="font-semibold mb-4">Enlaces</h4>
                        <ul className="space-y-2 text-sm text-gray-400">
                            <li><Link href={`/tienda/${slug}`} className="hover:text-white transition-colors">Inicio</Link></li>
                            <li><Link href={`/tienda/${slug}/productos`} className="hover:text-white transition-colors">Productos</Link></li>
                            <li><Link href={`/tienda/${slug}/contacto`} className="hover:text-white transition-colors">Contacto</Link></li>
                        </ul>
                    </div>
                    <div>
                        <h4 className="font-semibold mb-4">Mi Cuenta</h4>
                        <ul className="space-y-2 text-sm text-gray-400">
                            <li><Link href={`/login?redirect=/storefront/catalogo&business_code=${slug}`} className="hover:text-white transition-colors">Iniciar Sesión</Link></li>
                            <li><Link href={`/login?redirect=/storefront/catalogo&business_code=${slug}`} className="hover:text-white transition-colors">Hacer Pedido</Link></li>
                        </ul>
                    </div>
                </div>
                <div className="border-t border-gray-800 mt-8 pt-8 text-center text-sm text-gray-500 dark:text-gray-400">
                    <p>&copy; {new Date().getFullYear()} {business.name}. Powered by Probability.</p>
                </div>
            </div>
        </footer>
    );
}
