import { CartProvider } from '@/services/modules/storefront/ui/cart/cart-context';
import { CartButton } from '@/services/modules/storefront/ui/cart/CartButton';
import { CartDrawer } from '@/services/modules/storefront/ui/cart/CartDrawer';

export default function StorefrontLayout({ children }: { children: React.ReactNode }) {
    return (
        <CartProvider>
            {children}
            <CartButton />
            <CartDrawer />
        </CartProvider>
    );
}
