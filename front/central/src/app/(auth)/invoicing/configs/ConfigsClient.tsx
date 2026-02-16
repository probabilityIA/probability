/**
 * Client Component para la tabla de configuraciones de facturación
 * Maneja la interactividad (toggle, delete, etc.)
 */

'use client';

import { useState } from 'react';
import { Button } from '@/shared/ui/button';
import { Table } from '@/shared/ui/table';
import { Badge } from '@/shared/ui/badge';
import { useToast } from '@/shared/providers/toast-provider';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { InvoicingHeader } from '@/services/modules/invoicing/ui/components/InvoicingHeader';
import {
  deleteConfigAction,
  enableConfigAction,
  disableConfigAction,
  enableAutoInvoiceAction,
  disableAutoInvoiceAction,
} from '@/services/modules/invoicing/infra/actions';
import type { InvoicingConfig } from '@/services/modules/invoicing/domain/types';
import type { Business } from '@/services/auth/business/domain/types';
import { useRouter } from 'next/navigation';

interface ConfigsClientProps {
  initialConfigs: InvoicingConfig[];
  businesses?: Business[];
  isSuperAdmin: boolean;
  selectedBusinessId?: number | null;
}

export function ConfigsClient({ initialConfigs, businesses, isSuperAdmin, selectedBusinessId }: ConfigsClientProps) {
  const { showToast } = useToast();
  const router = useRouter();
  const [configs, setConfigs] = useState<InvoicingConfig[]>(initialConfigs);
  const [selectedConfig, setSelectedConfig] = useState<InvoicingConfig | null>(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

  const handleToggleEnabled = async (config: InvoicingConfig) => {
    try {
      setActionLoading(true);

      // Actualizar estado local INMEDIATAMENTE (optimista)
      setConfigs(prevConfigs =>
        prevConfigs.map(c =>
          c.id === config.id ? { ...c, enabled: !c.enabled } : c
        )
      );

      // Llamar al servidor en background
      if (config.enabled) {
        await disableConfigAction(config.id);
        showToast('Configuración desactivada exitosamente', 'success');
      } else {
        await enableConfigAction(config.id);
        showToast('Configuración activada exitosamente', 'success');
      }

      // Sincronizar con servidor (sin recargar página)
      router.refresh();
    } catch (error: any) {
      // Si falla, revertir el cambio optimista
      setConfigs(prevConfigs =>
        prevConfigs.map(c =>
          c.id === config.id ? { ...c, enabled: config.enabled } : c
        )
      );
      showToast('Error al actualizar configuración: ' + error.message, 'error');
    } finally {
      setActionLoading(false);
    }
  };

  const handleToggleAutoInvoice = async (config: InvoicingConfig) => {
    try {
      setActionLoading(true);

      // Actualizar estado local INMEDIATAMENTE (optimista)
      setConfigs(prevConfigs =>
        prevConfigs.map(c =>
          c.id === config.id ? { ...c, auto_invoice: !c.auto_invoice } : c
        )
      );

      // Llamar al servidor en background
      if (config.auto_invoice) {
        await disableAutoInvoiceAction(config.id);
        showToast('Facturación automática desactivada', 'success');
      } else {
        await enableAutoInvoiceAction(config.id);
        showToast('Facturación automática activada', 'success');
      }

      // Sincronizar con servidor (sin recargar página)
      router.refresh();
    } catch (error: any) {
      // Si falla, revertir el cambio optimista
      setConfigs(prevConfigs =>
        prevConfigs.map(c =>
          c.id === config.id ? { ...c, auto_invoice: config.auto_invoice } : c
        )
      );
      showToast('Error al actualizar facturación automática: ' + error.message, 'error');
    } finally {
      setActionLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!selectedConfig) return;

    try {
      setActionLoading(true);
      await deleteConfigAction(selectedConfig.id);

      // Actualizar estado local INMEDIATAMENTE
      setConfigs(prevConfigs => prevConfigs.filter(c => c.id !== selectedConfig.id));

      showToast('Configuración eliminada exitosamente', 'success');
      setShowDeleteModal(false);
      router.refresh();
    } catch (error: any) {
      showToast('Error al eliminar configuración: ' + error.message, 'error');
    } finally {
      setActionLoading(false);
    }
  };

  const handleBusinessChange = (businessId: string) => {
    const params = new URLSearchParams();
    if (businessId) {
      params.set('business_id', businessId);
    }
    router.push(`/invoicing/configs?${params.toString()}`);
  };

  const columns = [
    // Columna de Negocio
    {
      key: 'business',
      label: 'Negocio',
      render: (_: unknown, config: InvoicingConfig) => {
        if (!isSuperAdmin) {
          return <span className="text-sm font-medium text-gray-700">Mi Negocio</span>;
        }

        const business = businesses?.find(b => b.id === config.business_id);
        return (
          <div className="text-sm text-gray-700">
            {business ? (
              <span className="font-medium">{business.name}</span>
            ) : (
              <Badge type="warning">ID: {config.business_id}</Badge>
            )}
          </div>
        );
      },
    },
    {
      key: 'integration',
      label: 'Integración',
      render: (_: unknown, config: InvoicingConfig) => (
        <div>
          <div className="font-medium">{config.integration_name || `ID: ${config.integration_id}`}</div>
          {config.description && (
            <div className="text-xs text-gray-500">{config.description}</div>
          )}
        </div>
      ),
    },
    {
      key: 'provider',
      label: 'Proveedor',
      render: (_: unknown, config: InvoicingConfig) => {
        if (config.provider_image_url) {
          // Mostrar solo el logo si está disponible
          return (
            <div className="flex items-center justify-center">
              <img
                src={config.provider_image_url}
                alt={config.provider_name || 'Proveedor'}
                className="h-8 w-auto object-contain"
                title={config.provider_name || 'Proveedor'}
              />
            </div>
          );
        }

        // Fallback: mostrar inicial con color si no hay logo
        const providerName = config.provider_name || `ID: ${config.invoicing_provider_id}`;
        const firstLetter = providerName.charAt(0).toUpperCase();

        const providerColors: Record<string, string> = {
          'softpymes': 'bg-blue-500',
          'alegra': 'bg-green-500',
          'siigo': 'bg-purple-500',
          'default': 'bg-gray-500'
        };

        const providerKey = providerName.toLowerCase();
        const bgColor = providerColors[providerKey] || providerColors['default'];

        return (
          <div className="flex items-center justify-center">
            <div
              className={`w-10 h-10 ${bgColor} rounded-full flex items-center justify-center text-white font-bold text-sm`}
              title={providerName}
            >
              {firstLetter}
            </div>
          </div>
        );
      },
    },
    {
      key: 'created_at',
      label: 'Creado',
      render: (_: unknown, config: InvoicingConfig) => (
        <div className="text-sm text-gray-600">
          {new Date(config.created_at).toLocaleDateString('es-CO')}
        </div>
      ),
    },
    {
      key: 'status',
      label: 'Estado',
      render: (_: unknown, config: InvoicingConfig) => (
        <Button
          variant={config.enabled ? 'success' : 'danger'}
          size="sm"
          onClick={() => handleToggleEnabled(config)}
          disabled={actionLoading}
        >
          {config.enabled ? 'Activo' : 'Inactivo'}
        </Button>
      ),
    },
    {
      key: 'auto_invoice',
      label: 'Facturación Automática',
      render: (_: unknown, config: InvoicingConfig) => (
        <Button
          variant={config.auto_invoice ? 'success' : 'danger'}
          size="sm"
          onClick={() => handleToggleAutoInvoice(config)}
          disabled={actionLoading}
        >
          {config.auto_invoice ? 'Sí' : 'No'}
        </Button>
      ),
    },
    {
      key: 'actions',
      label: 'Acciones',
      render: (_: unknown, config: InvoicingConfig) => (
        <Button
          variant="danger"
          size="sm"
          onClick={() => {
            setSelectedConfig(config);
            setShowDeleteModal(true);
          }}
          disabled={actionLoading}
        >
          Eliminar
        </Button>
      ),
    },
  ];

  return (
    <div className="p-8">
      <InvoicingHeader
        title="Configuración de Facturación"
        description="Define qué integraciones deben generar facturas automáticamente"
      >
        <button
          onClick={() => router.push('/invoicing/configs/new')}
          className="px-6 py-2.5 bg-gradient-to-r from-[#7c3aed] to-[#6d28d9] text-white font-bold rounded-full shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-300 flex items-center gap-2"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Nueva Configuración
        </button>
      </InvoicingHeader>

      {/* Dropdown de Business para Super Admins */}
      {isSuperAdmin && businesses && businesses.length > 0 && (
        <div className="mb-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Seleccionar Negocio (Super Admin)
          </label>
          <select
            defaultValue={selectedBusinessId?.toString() || ''}
            onChange={(e) => handleBusinessChange(e.target.value)}
            className="w-full max-w-md px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">Todos los negocios</option>
            {businesses.map((business) => (
              <option key={business.id} value={business.id}>
                {business.name} (ID: {business.id})
              </option>
            ))}
          </select>
        </div>
      )}

      {!configs || configs.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-12 text-center">
          <svg
            className="mx-auto h-12 w-12 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
            />
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          <h3 className="mt-2 text-sm font-medium text-gray-900">
            No hay configuraciones
          </h3>
          <p className="mt-1 text-sm text-gray-500">
            Las configuraciones se crean automáticamente cuando se conecta una integración.
            <br />
            Ve a <strong>Integraciones</strong> para conectar Shopify, MercadoLibre, etc.
          </p>
        </div>
      ) : (
        <div className="configsTable">
          <div className="mb-4">
            <p className="text-sm text-gray-600">
              <strong>Tip:</strong> Haz clic en los badges para activar/desactivar rápidamente
            </p>
          </div>
          <Table
            data={configs}
            columns={columns}
            emptyMessage="No hay configuraciones para mostrar"
          />

          <style jsx>{`
            /* Tabla mejorada similar a Facturas */
            .configsTable :global(.table) {
              border-collapse: separate;
              border-spacing: 0 10px;
              background: transparent;
            }

            /* Quitar el borde del contenedor global de Table */
            .configsTable :global(div.overflow-hidden.w-full.rounded-lg.border.border-gray-200.bg-white) {
              border: none !important;
              background: transparent !important;
            }

            .configsTable :global(.table th) {
              background: linear-gradient(135deg, #7c3aed 0%, #6d28d9 100%);
              color: #fff;
              position: sticky;
              top: 0;
              z-index: 1;
            }

            /* Header más llamativo + bordes redondeados */
            .configsTable :global(.table thead th) {
              padding-top: 18px;
              padding-bottom: 18px;
              font-size: 0.875rem;
              font-weight: 800;
              letter-spacing: 0.06em;
              text-transform: uppercase;
              box-shadow: 0 10px 25px rgba(124, 58, 237, 0.18);
            }

            .configsTable :global(.table thead th:first-child) {
              border-top-left-radius: 14px;
              border-bottom-left-radius: 14px;
            }

            .configsTable :global(.table thead th:last-child) {
              border-top-right-radius: 14px;
              border-bottom-right-radius: 14px;
            }

            .configsTable :global(.table tbody tr) {
              background: rgba(255, 255, 255, 0.95);
              box-shadow: 0 1px 0 rgba(17, 24, 39, 0.04);
              transition: transform 180ms ease, box-shadow 180ms ease, background 180ms ease;
            }

            /* Zebra suave en morado */
            .configsTable :global(.table tbody tr:nth-child(even)) {
              background: rgba(124, 58, 237, 0.03);
            }

            .configsTable :global(.table tbody tr:hover) {
              background: rgba(124, 58, 237, 0.06);
              box-shadow: 0 10px 25px rgba(17, 24, 39, 0.08);
              transform: translateY(-1px);
            }

            .configsTable :global(.table td) {
              border-top: none;
            }

            /* Redondeo de cada fila */
            .configsTable :global(.table tbody td:first-child) {
              border-top-left-radius: 12px;
              border-bottom-left-radius: 12px;
            }
            .configsTable :global(.table tbody td:last-child) {
              border-top-right-radius: 12px;
              border-bottom-right-radius: 12px;
            }

            /* Acciones: focus consistente */
            .configsTable :global(a),
            .configsTable :global(button) {
              outline-color: rgba(124, 58, 237, 0.35);
            }
          `}</style>
        </div>
      )}

      <ConfirmModal
        isOpen={showDeleteModal}
        onClose={() => setShowDeleteModal(false)}
        onConfirm={handleDelete}
        title="Eliminar Configuración"
        message="¿Estás seguro de que deseas eliminar esta configuración de facturación?"
        confirmText="Sí, eliminar"
        cancelText="Cancelar"
        type="danger"
      />
    </div>
  );
}
