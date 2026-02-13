/**
 * Header reutilizable para el módulo de Facturación
 * Muestra los 3 menús principales: Facturas, Proveedores, Configuraciones
 */

'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

interface InvoicingHeaderProps {
  title: string;
  description?: string;
  children?: React.ReactNode;
}

export function InvoicingHeader({ title, description, children }: InvoicingHeaderProps) {
  const pathname = usePathname();

  const isActive = (path: string) => pathname === path;

  const menuItems = [
    { label: 'Facturas', href: '/invoicing/invoices' },
    { label: 'Proveedores', href: '/invoicing/providers' },
    { label: 'Configuraciones', href: '/invoicing/configs' },
  ];

  return (
    <div className="mb-8">
      {/* Título con menús integrados */}
      <div className="flex items-center justify-between mb-8 gap-6">
        <div>
          <h1 className="text-4xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
            {title}
          </h1>
          {description && (
            <p className="text-gray-600 mt-3 text-base">{description}</p>
          )}
        </div>

        <div className="flex items-center gap-3">
          {/* Menús integrados a la derecha */}
          <div className="flex items-center gap-3 bg-gradient-to-r from-[#7c3aed]/10 to-[#6d28d9]/10 rounded-full p-2 shadow-lg border-2 border-[#7c3aed]/30 backdrop-blur-sm h-fit">
            {menuItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`
                  px-6 py-3 rounded-full text-sm font-bold transition-all duration-300 flex items-center gap-2 whitespace-nowrap
                  ${isActive(item.href)
                    ? 'bg-gradient-to-r from-[#7c3aed] to-[#6d28d9] text-white shadow-xl scale-110 transform'
                    : 'text-[#7c3aed] hover:bg-[#7c3aed]/20 hover:text-[#6d28d9]'
                  }
                `}
              >
                {item.label}
              </Link>
            ))}
          </div>

          {/* Botón más a la derecha */}
          {children && <div className="flex-shrink-0">{children}</div>}
        </div>
      </div>
    </div>
  );
}
