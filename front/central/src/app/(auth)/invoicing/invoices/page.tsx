'use client';

import { useRef } from 'react';
import { InvoiceList } from '@/services/modules/invoicing/ui/components/InvoiceList';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInvoicingBusiness } from '@/shared/contexts/invoicing-business-context';

export default function InvoicesPage() {
  const { permissions, isSuperAdmin } = usePermissions();
  const { selectedBusinessId } = useInvoicingBusiness();
  const businessId = permissions?.business_id || 0;
  const invoiceListRef = useRef<any>(null);

  if (isSuperAdmin && !selectedBusinessId) {
    return (
      <div className="p-8">
        <div className="flex flex-col items-center justify-center py-24 text-gray-500 dark:text-gray-400">
          <svg className="w-16 h-16 mb-4 text-purple-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
          </svg>
          <p className="text-lg font-medium">Selecciona un negocio para continuar</p>
          <p className="text-sm mt-1">Usa el selector en la barra superior</p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8 min-h-screen bg-gray-50 dark:bg-gray-900">
      <div ref={invoiceListRef} />
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm dark:shadow-lg border border-gray-200 dark:border-gray-700 p-6">
        <InvoiceList ref={invoiceListRef} businessId={businessId} selectedBusinessId={selectedBusinessId} />
      </div>
    </div>
  );
}
