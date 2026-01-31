/**
 * Página de Proveedores de Facturación
 * Gestión de proveedores de facturación electrónica (Softpymes, Siigo, etc.)
 */

'use client';

import { useState, useEffect } from 'react';
import { Button } from '@/shared/ui/button';
import { Table } from '@/shared/ui/table';
import { Badge } from '@/shared/ui/badge';
import { Spinner } from '@/shared/ui/spinner';
import { useToast } from '@/shared/providers/toast-provider';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { ProviderForm } from '@/services/modules/invoicing/ui/components/ProviderForm';
import { getProvidersAction } from '@/services/modules/invoicing/infra/actions';
import type { InvoicingProvider } from '@/services/modules/invoicing/domain/types';

export default function InvoicingProvidersPage() {
  const { showToast } = useToast();
  const { permissions } = usePermissions();
  const [providers, setProviders] = useState<InvoicingProvider[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [selectedProvider, setSelectedProvider] = useState<InvoicingProvider | undefined>();

  const businessId = permissions?.business_id || 0;

  const loadProviders = async () => {
    if (!businessId || businessId === 0) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      const response = await getProvidersAction({ business_id: businessId });
      setProviders(response.data);
    } catch (error: any) {
      showToast('Error al cargar proveedores: ' + error.message, 'error');
      setProviders([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (businessId) {
      loadProviders();
    }
  }, [businessId]);

  const handleEdit = (provider: InvoicingProvider) => {
    setSelectedProvider(provider);
    setShowForm(true);
  };

  const handleNew = () => {
    setSelectedProvider(undefined);
    setShowForm(true);
  };

  const handleClose = () => {
    setShowForm(false);
    setSelectedProvider(undefined);
  };

  const columns = [
    {
      key: 'name',
      label: 'Nombre',
      render: (_: unknown, provider: InvoicingProvider) => (
        <div>
          <div className="font-medium">{provider.name}</div>
          {provider.description && (
            <div className="text-xs text-gray-500">{provider.description}</div>
          )}
        </div>
      ),
    },
    {
      key: 'provider_type_code',
      label: 'Tipo',
      render: (_: unknown, provider: InvoicingProvider) => (
        <Badge color="blue">{provider.provider_type_code.toUpperCase()}</Badge>
      ),
    },
    {
      key: 'is_active',
      label: 'Estado',
      render: (_: unknown, provider: InvoicingProvider) => (
        <Badge color={provider.is_active ? 'green' : 'gray'}>
          {provider.is_active ? 'Activo' : 'Inactivo'}
        </Badge>
      ),
    },
    {
      key: 'is_default',
      label: 'Por Defecto',
      render: (_: unknown, provider: InvoicingProvider) =>
        provider.is_default ? (
          <Badge color="purple">Predeterminado</Badge>
        ) : (
          <span className="text-gray-400">-</span>
        ),
    },
    {
      key: 'created_at',
      label: 'Creado',
      render: (_: unknown, provider: InvoicingProvider) => (
        <div className="text-sm text-gray-600">
          {new Date(provider.created_at).toLocaleDateString('es-CO')}
        </div>
      ),
    },
    {
      key: 'actions',
      label: 'Acciones',
      render: (_: unknown, provider: InvoicingProvider) => (
        <Button variant="secondary" size="sm" onClick={() => handleEdit(provider)}>
          Editar
        </Button>
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
        <Button variant="primary" onClick={handleNew}>
          + Nuevo Proveedor
        </Button>
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <Table
          data={providers}
          columns={columns}
          emptyMessage="No hay proveedores configurados. Crea uno para empezar a facturar."
        />
      </div>

      <ProviderForm
        isOpen={showForm}
        onClose={handleClose}
        onSuccess={() => {
          loadProviders();
          handleClose();
        }}
        provider={selectedProvider}
        businessId={businessId}
      />
    </div>
  );
}
