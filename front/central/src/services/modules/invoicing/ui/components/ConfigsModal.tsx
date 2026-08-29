'use client';

import { useState, useEffect } from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline';
import { ConfigsClient } from '@/app/(auth)/invoicing/configs/ConfigsClient';
import { InvoicingConfigForm } from './InvoicingConfigForm';
import { getConfigsAction } from '../../infra/actions';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';
import type { InvoicingConfig } from '../../domain/types';

interface ConfigsModalProps {
  isOpen: boolean;
  onClose: () => void;
  selectedBusinessId?: number | null;
}

export function ConfigsModal({ isOpen, onClose, selectedBusinessId }: ConfigsModalProps) {
  const { isSuperAdmin } = usePermissions();
  const { businesses } = useBusinessesSimple();
  const [configs, setConfigs] = useState<InvoicingConfig[]>([]);
  const [loading, setLoading] = useState(false);

  const loadConfigs = async () => {
    setLoading(true);
    try {
      const filters = selectedBusinessId ? { business_id: selectedBusinessId } : {};
      const response = await getConfigsAction(filters);
      setConfigs(response.data || []);
    } catch {
      setConfigs([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isOpen) {
      loadConfigs();
    }
  }, [isOpen, selectedBusinessId]);

  if (!isOpen) return null;

  const configUnica = configs.length === 1 ? configs[0] : null;
  const titulo = configUnica ? 'Configuración de Facturación' : 'Configuraciones de Facturación';
  const facturadorNombre = configUnica?.provider_name ?? '';
  const facturadorLogo = configUnica?.provider_image_url ?? '';

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div
        className="bg-white rounded-2xl shadow-2xl flex flex-col overflow-hidden"
        style={{ width: configUnica ? 'min(94vw, 58rem)' : '70vw', maxHeight: '90vh' }}
      >
        {/* Header */}
        <div className="flex items-center justify-between gap-4 px-6 py-4 border-b border-gray-200 flex-shrink-0">
          <div className="flex min-w-0 items-center gap-3">
            {facturadorLogo ? (
              <img src={facturadorLogo} alt={facturadorNombre} className="h-9 w-9 shrink-0 object-contain" />
            ) : (
              <svg className="h-6 w-6 shrink-0" style={{ color: 'var(--color-primary)' }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            )}
            <div className="min-w-0">
              <h2 className="truncate text-xl font-bold" style={{ color: 'var(--color-primary)' }}>{titulo}</h2>
              {facturadorNombre && (
                <p className="truncate text-xs text-gray-500">Facturador electronico: {facturadorNombre}</p>
              )}
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <XMarkIcon className="w-5 h-5" />
          </button>
        </div>

        {/* Body */}
        <div className="flex min-h-0 flex-1 flex-col">
          {loading ? (
            <div className="flex items-center justify-center h-48 text-gray-500">
              Cargando configuraciones...
            </div>
          ) : configUnica ? (
            <div className="flex min-h-0 flex-1 flex-col">
              <InvoicingConfigForm
                alturaFija
                integrationIds={configUnica.integration_ids ?? []}
                invoicingIntegrationId={configUnica.invoicing_integration_id ?? configUnica.invoicing_provider_id ?? 0}
                businessId={configUnica.business_id}
                initialData={configUnica}
                onSuccess={() => { loadConfigs(); onClose(); }}
                onCancel={onClose}
              />
            </div>
          ) : (
            <ConfigsClient
              initialConfigs={configs}
              businesses={businesses}
              isSuperAdmin={isSuperAdmin}
              selectedBusinessId={selectedBusinessId}
              hideNavbarButton
              onConfigsChange={loadConfigs}
            />
          )}
        </div>
      </div>
    </div>
  );
}
