'use client';

import { useState, useEffect } from 'react';
import { ConfigsClient } from './ConfigsClient';
import { getConfigsAction } from '@/services/modules/invoicing/infra/actions';
import { useInvoicingBusiness } from '@/shared/contexts/invoicing-business-context';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import type { InvoicingConfig } from '@/services/modules/invoicing/domain/types';

export default function InvoicingConfigsPage() {
  const { selectedBusinessId } = useInvoicingBusiness();
  const { isSuperAdmin } = usePermissions();
  const { businesses } = useBusinessesSimple();
  const [configs, setConfigs] = useState<InvoicingConfig[]>([]);
  const [loading, setLoading] = useState(false);

  const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

  useEffect(() => {
    if (requiresBusinessSelection) {
      setConfigs([]);
      return;
    }

    const fetchConfigs = async () => {
      setLoading(true);
      try {
        const filters = selectedBusinessId ? { business_id: selectedBusinessId } : {};
        const response = await getConfigsAction(filters);
        setConfigs(response.data || []);
      } catch (error) {
        console.error('Error loading invoicing configs:', error);
        setConfigs([]);
      } finally {
        setLoading(false);
      }
    };

    fetchConfigs();
  }, [selectedBusinessId, requiresBusinessSelection]);

  if (requiresBusinessSelection) {
    return (
      <div className="p-8 min-h-screen bg-gray-50 dark:bg-gray-900">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm dark:shadow-lg border border-gray-200 dark:border-gray-700 p-16 text-center text-gray-500 dark:text-gray-400">
          Selecciona un negocio para ver las configuraciones de facturación
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="p-8 min-h-screen bg-gray-50 dark:bg-gray-900">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm dark:shadow-lg border border-gray-200 dark:border-gray-700 p-16 text-center text-gray-500 dark:text-gray-400">
          Cargando configuraciones...
        </div>
      </div>
    );
  }

  return (
    <ConfigsClient
      initialConfigs={configs}
      businesses={businesses}
      isSuperAdmin={isSuperAdmin}
      selectedBusinessId={selectedBusinessId}
    />
  );
}
