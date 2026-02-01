/**
 * Página de Proveedores de Facturación
 * Gestión de proveedores de facturación electrónica (Softpymes, Siigo, etc.)
 * Ahora usa integrations/core en lugar del módulo deprecado
 */

'use client';

import { useState, useEffect } from 'react';
import { Button } from '@/shared/ui/button';
import { Table } from '@/shared/ui/table';
import { Badge } from '@/shared/ui/badge';
import { Spinner } from '@/shared/ui/spinner';
import { useToast } from '@/shared/providers/toast-provider';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { getIntegrationsAction } from '@/services/integrations/core/infra/actions';
import type { Integration } from '@/services/integrations/core/domain/types';

export default function InvoicingProvidersPage() {
  const { showToast } = useToast();
  const { permissions } = usePermissions();
  const [providers, setProviders] = useState<Integration[]>([]);
  const [loading, setLoading] = useState(true);

  const businessId = permissions?.business_id || 0;
  const isSuperAdmin = permissions?.is_super || false;

  const loadProviders = async () => {
    try {
      setLoading(true);
      // Filtrar por categoría "invoicing" (facturación electrónica)
      const filters: any = {
        category: 'invoicing',
        page: 1,
        page_size: 100
      };

      // Si no es super admin, filtrar también por business_id
      if (!isSuperAdmin && businessId) {
        filters.business_id = businessId;
      }

      const response = await getIntegrationsAction(filters);
      setProviders(response.data || []);
    } catch (error: any) {
      showToast('Error al cargar proveedores: ' + error.message, 'error');
      setProviders([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProviders();
  }, [businessId]);

  const columns = [
    {
      key: 'name',
      label: 'Nombre',
      render: (_: unknown, provider: Integration) => (
        <div>
          <div className="font-medium">{provider.name}</div>
          {provider.description && (
            <div className="text-xs text-gray-500">{provider.description}</div>
          )}
        </div>
      ),
    },
    {
      key: 'integration_type',
      label: 'Tipo',
      render: (_: unknown, provider: Integration) => (
        <Badge color="blue">
          {provider.integration_type?.name || provider.integration_type?.code || 'N/A'}
        </Badge>
      ),
    },
    {
      key: 'is_active',
      label: 'Estado',
      render: (_: unknown, provider: Integration) => (
        <Badge color={provider.is_active ? 'green' : 'gray'}>
          {provider.is_active ? 'Activo' : 'Inactivo'}
        </Badge>
      ),
    },
    {
      key: 'is_default',
      label: 'Por Defecto',
      render: (_: unknown, provider: Integration) =>
        provider.is_default ? (
          <Badge color="purple">Predeterminado</Badge>
        ) : (
          <span className="text-gray-400">-</span>
        ),
    },
    {
      key: 'created_at',
      label: 'Creado',
      render: (_: unknown, provider: Integration) => (
        <div className="text-sm text-gray-600">
          {new Date(provider.created_at).toLocaleDateString('es-CO')}
        </div>
      ),
    },
  ];

  if (loading) {
    return (
      <div className="flex justify-center items-center h-screen">
        <Spinner />
      </div>
    );
  }

  return (
    <div className="p-8">
      <div className="mb-6 flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Proveedores de Facturación</h1>
          <p className="text-gray-600 mt-2">
            Configura los proveedores de facturación electrónica para tu negocio
          </p>
        </div>
        <Button
          variant="primary"
          onClick={() => window.location.href = '/integrations'}
        >
          + Nuevo Proveedor
        </Button>
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <Table
          data={providers}
          columns={columns}
          emptyMessage="No hay proveedores configurados. Crea uno desde Integraciones → Facturación Electrónica."
        />
      </div>
    </div>
  );
}
