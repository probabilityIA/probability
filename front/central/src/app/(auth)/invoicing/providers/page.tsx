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
import { InvoicingHeader } from '@/services/modules/invoicing/ui/components/InvoicingHeader';
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
        <Badge type="warning">
          {provider.integration_type?.name || provider.integration_type?.code || 'N/A'}
        </Badge>
      ),
    },
    {
      key: 'is_active',
      label: 'Estado',
      render: (_: unknown, provider: Integration) => (
        <Badge type="success">
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
      <InvoicingHeader
        title="Proveedores de Facturación"
        description="Configura los proveedores de facturación electrónica para tu negocio"
      >
        <button
          onClick={() => window.location.href = '/integrations'}
          className="px-6 py-2.5 bg-gradient-to-r from-[#7c3aed] to-[#6d28d9] text-white font-bold rounded-full shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-300 flex items-center gap-2"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Nuevo Proveedor
        </button>
      </InvoicingHeader>

      <div className="providersTable">
        <Table
          data={providers}
          columns={columns}
          emptyMessage="No hay proveedores configurados. Crea uno desde Integraciones → Facturación Electrónica."
        />

        <style jsx>{`
          /* Tabla mejorada similar a Facturas */
          .providersTable :global(.table) {
            border-collapse: separate;
            border-spacing: 0 10px;
            background: transparent;
          }

          /* Quitar el borde del contenedor global de Table */
          .providersTable :global(div.overflow-hidden.w-full.rounded-lg.border.border-gray-200.bg-white) {
            border: none !important;
            background: transparent !important;
          }

          .providersTable :global(.table th) {
            background: linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%);
            color: #fff;
            position: sticky;
            top: 0;
            z-index: 1;
          }

          /* Header más llamativo + bordes redondeados */
          .providersTable :global(.table thead th) {
            padding-top: 18px;
            padding-bottom: 18px;
            font-size: 0.875rem;
            font-weight: 800;
            letter-spacing: 0.06em;
            text-transform: uppercase;
            box-shadow: 0 10px 25px rgba(124, 58, 237, 0.18);
          }

          .providersTable :global(.table thead th:first-child) {
            border-top-left-radius: 14px;
            border-bottom-left-radius: 14px;
          }

          .providersTable :global(.table thead th:last-child) {
            border-top-right-radius: 14px;
            border-bottom-right-radius: 14px;
          }

          .providersTable :global(.table tbody tr) {
            background: rgba(255, 255, 255, 0.95);
            box-shadow: 0 1px 0 rgba(17, 24, 39, 0.04);
            transition: transform 180ms ease, box-shadow 180ms ease, background 180ms ease;
          }

          /* Zebra suave en morado */
          .providersTable :global(.table tbody tr:nth-child(even)) {
            background: rgba(124, 58, 237, 0.03);
          }

          .providersTable :global(.table tbody tr:hover) {
            background: rgba(124, 58, 237, 0.06);
            box-shadow: 0 10px 25px rgba(17, 24, 39, 0.08);
            transform: translateY(-1px);
          }

          .providersTable :global(.table td) {
            border-top: none;
          }

          /* Redondeo de cada fila */
          .providersTable :global(.table tbody td:first-child) {
            border-top-left-radius: 12px;
            border-bottom-left-radius: 12px;
          }
          .providersTable :global(.table tbody td:last-child) {
            border-top-right-radius: 12px;
            border-bottom-right-radius: 12px;
          }

          /* Acciones: focus consistente */
          .providersTable :global(a),
          .providersTable :global(button) {
            outline-color: rgba(124, 58, 237, 0.35);
          }
        `}</style>
      </div>
    </div>
  );
}
