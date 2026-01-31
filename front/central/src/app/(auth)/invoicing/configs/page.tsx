/**
 * Página de Configuración de Facturación
 * Configuración por integración (qué integraciones facturan automáticamente)
 */

'use client';

import { useState, useEffect } from 'react';
import { Button } from '@/shared/ui/button';
import { Table } from '@/shared/ui/table';
import { Badge } from '@/shared/ui/badge';
import { Spinner } from '@/shared/ui/spinner';
import { useToast } from '@/shared/providers/toast-provider';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import {
  getConfigsAction,
  deleteConfigAction,
  updateConfigAction,
} from '@/services/modules/invoicing/infra/actions';
import type { InvoicingConfig } from '@/services/modules/invoicing/domain/types';

export default function InvoicingConfigsPage() {
  const { showToast } = useToast();
  const { permissions } = usePermissions();
  const [configs, setConfigs] = useState<InvoicingConfig[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedConfig, setSelectedConfig] = useState<InvoicingConfig | null>(null);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

  const businessId = permissions?.business_id || 0;

  const loadConfigs = async () => {
    if (!businessId || businessId === 0) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      const response = await getConfigsAction({ business_id: businessId });
      setConfigs(response.data);
    } catch (error: any) {
      showToast('Error al cargar configuraciones: ' + error.message, 'error');
      setConfigs([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (businessId) {
      loadConfigs();
    }
  }, [businessId]);

  const handleToggleEnabled = async (config: InvoicingConfig) => {
    try {
      setActionLoading(true);
      await updateConfigAction(config.id, { enabled: !config.enabled });
      showToast(
        config.enabled ? 'Facturación desactivada' : 'Facturación activada',
        'success'
      );
      loadConfigs();
    } catch (error: any) {
      showToast('Error al actualizar configuración: ' + error.message, 'error');
    } finally {
      setActionLoading(false);
    }
  };

  const handleToggleAutoInvoice = async (config: InvoicingConfig) => {
    try {
      setActionLoading(true);
      await updateConfigAction(config.id, { auto_invoice: !config.auto_invoice });
      showToast(
        config.auto_invoice
          ? 'Facturación automática desactivada'
          : 'Facturación automática activada',
        'success'
      );
      loadConfigs();
    } catch (error: any) {
      showToast('Error al actualizar configuración: ' + error.message, 'error');
    } finally {
      setActionLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!selectedConfig) return;

    try {
      setActionLoading(true);
      await deleteConfigAction(selectedConfig.id);
      showToast('Configuración eliminada exitosamente', 'success');
      setShowDeleteModal(false);
      loadConfigs();
    } catch (error: any) {
      showToast('Error al eliminar configuración: ' + error.message, 'error');
    } finally {
      setActionLoading(false);
    }
  };

  const columns = [
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
      render: (_: unknown, config: InvoicingConfig) => (
        <div className="text-sm text-gray-700">
          {config.provider_name || `ID: ${config.invoicing_provider_id}`}
        </div>
      ),
    },
    {
      key: 'enabled',
      label: 'Habilitado',
      render: (_: unknown, config: InvoicingConfig) => (
        <button
          onClick={() => handleToggleEnabled(config)}
          disabled={actionLoading}
          className="cursor-pointer"
        >
          <Badge color={config.enabled ? 'green' : 'gray'}>
            {config.enabled ? 'Sí' : 'No'}
          </Badge>
        </button>
      ),
    },
    {
      key: 'auto_invoice',
      label: 'Auto-facturar',
      render: (_: unknown, config: InvoicingConfig) => (
        <button
          onClick={() => handleToggleAutoInvoice(config)}
          disabled={actionLoading}
          className="cursor-pointer"
        >
          <Badge color={config.auto_invoice ? 'blue' : 'gray'}>
            {config.auto_invoice ? 'Automático' : 'Manual'}
          </Badge>
        </button>
      ),
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
          <h1 className="text-3xl font-bold text-gray-900">Configuración de Facturación</h1>
          <p className="text-gray-600 mt-2">
            Define qué integraciones deben generar facturas automáticamente
          </p>
        </div>
      </div>

      {configs.length === 0 ? (
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
        <div className="bg-white rounded-lg shadow p-6">
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
