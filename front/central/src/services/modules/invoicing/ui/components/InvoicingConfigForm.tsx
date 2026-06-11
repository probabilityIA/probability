'use client';

import { useState, useEffect, FormEvent } from 'react';
import type {
  InvoicingConfig,
  CreateConfigDTO,
  BankAccountResult,
} from '@/services/modules/invoicing/domain/types';
import type { Integration } from '@/services/integrations/core/domain/types';
import { getIntegrationsAction } from '@/services/integrations/core/infra/actions';
import { useInvoicingConfig } from '@/services/modules/invoicing/ui/hooks/useInvoicingConfig';
import {
  requestListBankAccountsAction,
  getListBankAccountsResultAction,
} from '@/services/modules/invoicing/infra/actions';

interface InvoicingConfigFormProps {
  integrationIds: number[];
  invoicingIntegrationId: number;
  businessId: number;
  onSuccess?: () => void;
  onCancel?: () => void;
  initialData?: InvoicingConfig;
}

export function InvoicingConfigForm({
  integrationIds,
  invoicingIntegrationId,
  businessId,
  onSuccess,
  onCancel,
  initialData,
}: InvoicingConfigFormProps) {
  const { createConfig, updateConfig, loading } = useInvoicingConfig(businessId);

  const [formData, setFormData] = useState({
    enabled: initialData?.enabled ?? true,
    auto_invoice: initialData?.auto_invoice ?? false,
    payment_status: (initialData?.filters?.payment_status as string) ?? '',
    invoice_cod: (initialData?.filters?.invoice_cod as boolean) ?? false,
    force_default_customer: initialData?.config?.force_default_customer ?? false,
    send_cash_receipt: initialData?.config?.send_cash_receipt ?? false,
    payment_type: (initialData?.config?.payment_type as string) ?? 'EF',
    payment_bank_account_id: initialData?.config?.payment_bank_account_id ?? '' as string | number,
    cod_use_alternate_bank: initialData?.config?.cod_use_alternate_bank ?? false,
    cod_payment_bank_account_id: initialData?.config?.cod_payment_bank_account_id ?? '' as string | number,
    payment_financial_entity_id: initialData?.config?.payment_financial_entity_id ?? '' as string | number,
    payment_bonus_code: (initialData?.config?.payment_bonus_code as string) ?? '',
    payment_bank_name: (initialData?.config?.payment_bank_name as string) ?? '',
    payment_account_number: (initialData?.config?.payment_account_number as string) ?? '',
    // Item mappings
    item_mappings_shipping: (initialData?.config?.item_mappings?.shipping as string) ?? '',
    item_mappings_membership: (initialData?.config?.item_mappings?.membership as string) ?? '',
    item_mappings_tip: (initialData?.config?.item_mappings?.tip as string) ?? '',
    // Siigo-specific ids
    siigo_document_id: (initialData?.config?.document_id as number | string) ?? '',
    siigo_payment_method_id: (initialData?.config?.payment_method_id as number | string) ?? '',
    siigo_tax_id: (initialData?.config?.tax_id as number | string) ?? '',
    siigo_seller_id: (initialData?.config?.seller_id as number | string) ?? '',
    siigo_cash_receipt_document_id: (initialData?.config?.cash_receipt_document_id as number | string) ?? '',
    siigo_cash_receipt_payment_id: (initialData?.config?.cash_receipt_payment_id as number | string) ?? '',
    siigo_credit_note_document_id: (initialData?.config?.credit_note_document_id as number | string) ?? '',
  });

  const providerName = initialData?.provider_name ?? '';
  const providerImageUrl = initialData?.provider_image_url;
  const isSiigo = providerName.toLowerCase().includes('siigo');
  const cashReceiptDesc = isSiigo
    ? 'Registra un recibo de caja en Siigo al crear la factura'
    : 'Registra el pago en Softpymes al crear la factura (mueve cuentas por cobrar al medio de pago)';

  const [showItemMappings, setShowItemMappings] = useState(
    !!(initialData?.config?.item_mappings?.shipping ||
       initialData?.config?.item_mappings?.membership ||
       initialData?.config?.item_mappings?.tip)
  );

  // Selección de integraciones de origen
  const initialSelected = initialData?.integration_ids?.length
    ? initialData.integration_ids
    : integrationIds;
  const [selectedIntegrationIds, setSelectedIntegrationIds] = useState<number[]>(initialSelected);
  const [availableIntegrations, setAvailableIntegrations] = useState<Integration[]>([]);
  const [loadingIntegrations, setLoadingIntegrations] = useState(false);

  const [error, setError] = useState<string | null>(null);

  // Bank accounts state
  const [loadingBankAccounts, setLoadingBankAccounts] = useState(false);
  const [bankAccounts, setBankAccounts] = useState<BankAccountResult[] | null>(null);

  // Cargar integraciones de origen disponibles del negocio (ecommerce + platform)
  useEffect(() => {
    if (!businessId) return;
    setLoadingIntegrations(true);
    Promise.all([
      getIntegrationsAction({ business_id: businessId, category: 'ecommerce', page_size: 100 }),
      getIntegrationsAction({ business_id: businessId, category: 'platform', page_size: 100 }),
    ])
      .then(([ecommerce, platform]) => {
        setAvailableIntegrations([...(ecommerce.data ?? []), ...(platform.data ?? [])]);
      })
      .catch(() => setAvailableIntegrations([]))
      .finally(() => setLoadingIntegrations(false));
  }, [businessId]);

  const toggleIntegration = (id: number) => {
    setSelectedIntegrationIds((prev) =>
      prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]
    );
  };

  const handleFetchBankAccounts = async () => {
    setLoadingBankAccounts(true);
    setBankAccounts(null);
    try {
      const result = await requestListBankAccountsAction(businessId);
      const correlationId = result.correlation_id;

      // Poll every 2 seconds for up to 30 seconds
      let attempts = 0;
      const maxAttempts = 15;
      const poll = setInterval(async () => {
        attempts++;
        try {
          const data = await getListBankAccountsResultAction(correlationId, businessId);
          if (data !== null) {
            setBankAccounts(data.results);
            setLoadingBankAccounts(false);
            clearInterval(poll);
          }
        } catch {
          // Ignore polling errors
        }
        if (attempts >= maxAttempts) {
          setLoadingBankAccounts(false);
          clearInterval(poll);
        }
      }, 2000);
    } catch {
      setLoadingBankAccounts(false);
    }
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);

    if (selectedIntegrationIds.length === 0) {
      setError('Debes seleccionar al menos una integración de origen de órdenes');
      return;
    }

    const filters: Record<string, any> = {};
    if (formData.payment_status) {
      filters.payment_status = formData.payment_status as 'paid' | 'unpaid' | 'partial';
    }
    if (formData.invoice_cod) {
      filters.invoice_cod = true;
    }

    // Build invoice_config with cash receipt settings
    const invoiceConfig: Record<string, any> = {};
    if (formData.force_default_customer) {
      invoiceConfig.force_default_customer = true;
    }
    if (isSiigo) {
      if (formData.siigo_document_id) invoiceConfig.document_id = Number(formData.siigo_document_id);
      if (formData.siigo_payment_method_id) invoiceConfig.payment_method_id = Number(formData.siigo_payment_method_id);
      if (formData.siigo_tax_id) invoiceConfig.tax_id = Number(formData.siigo_tax_id);
      if (formData.siigo_seller_id) invoiceConfig.seller_id = Number(formData.siigo_seller_id);
      if (formData.siigo_credit_note_document_id) invoiceConfig.credit_note_document_id = Number(formData.siigo_credit_note_document_id);
    }

    if (formData.send_cash_receipt) {
      invoiceConfig.send_cash_receipt = true;
      if (isSiigo) {
        if (formData.siigo_cash_receipt_document_id) invoiceConfig.cash_receipt_document_id = Number(formData.siigo_cash_receipt_document_id);
        if (formData.siigo_cash_receipt_payment_id) invoiceConfig.cash_receipt_payment_id = Number(formData.siigo_cash_receipt_payment_id);
      } else {
        invoiceConfig.payment_type = formData.payment_type || 'EF';
        if (formData.payment_type === 'TR' && formData.payment_bank_account_id)
          invoiceConfig.payment_bank_account_id = String(formData.payment_bank_account_id);
        if (formData.payment_type === 'CH') {
          if (formData.payment_account_number) invoiceConfig.payment_account_number = formData.payment_account_number;
          if (formData.payment_bank_name) invoiceConfig.payment_bank_name = formData.payment_bank_name;
        }
        if ((formData.payment_type === 'TC' || formData.payment_type === 'TD') && formData.payment_financial_entity_id)
          invoiceConfig.payment_financial_entity_id = Number(formData.payment_financial_entity_id);
        if (formData.payment_type === 'BN' && formData.payment_bonus_code)
          invoiceConfig.payment_bonus_code = formData.payment_bonus_code;
        if (formData.cod_use_alternate_bank) {
          invoiceConfig.cod_use_alternate_bank = true;
          if (formData.cod_payment_bank_account_id)
            invoiceConfig.cod_payment_bank_account_id = String(formData.cod_payment_bank_account_id);
        }
      }
    }

    // Build item_mappings
    const itemMappings: Record<string, string> = {};
    if (formData.item_mappings_shipping) itemMappings.shipping = formData.item_mappings_shipping;
    if (formData.item_mappings_membership) itemMappings.membership = formData.item_mappings_membership;
    if (formData.item_mappings_tip) itemMappings.tip = formData.item_mappings_tip;
    if (Object.keys(itemMappings).length > 0) invoiceConfig.item_mappings = itemMappings;

    try {
      if (initialData?.id) {
        const result = await updateConfig(initialData.id, {
          enabled: formData.enabled,
          auto_invoice: formData.auto_invoice,
          filters,
          integration_ids: selectedIntegrationIds,
          config: invoiceConfig,
        });

        if (result.success) {
          onSuccess?.();
        } else {
          setError(result.error || 'Error al actualizar configuración');
        }
      } else {
        const createData: CreateConfigDTO = {
          business_id: businessId,
          integration_ids: selectedIntegrationIds,
          invoicing_integration_id: invoicingIntegrationId,
          enabled: formData.enabled,
          auto_invoice: formData.auto_invoice,
          filters,
          config: invoiceConfig,
        };

        const result = await createConfig(createData);

        if (result.success) {
          onSuccess?.();
        } else {
          setError(result.error || 'Error al crear la configuración');
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error desconocido');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-sm text-red-800">{error}</p>
        </div>
      )}

      {providerName && (
        <div className="flex items-center gap-3 bg-gradient-to-r from-blue-50 to-cyan-50 dark:from-gray-800 dark:to-gray-800 p-3 rounded-lg border border-blue-100 dark:border-gray-700">
          {providerImageUrl ? (
            <img src={providerImageUrl} alt={providerName} className="w-8 h-8 object-contain flex-shrink-0" />
          ) : (
            <div className="w-8 h-8 rounded-full bg-blue-500 text-white flex items-center justify-center text-sm font-bold flex-shrink-0">
              {providerName.charAt(0)}
            </div>
          )}
          <div className="min-w-0">
            <p className="text-xs text-gray-500 dark:text-gray-400">Facturador electrónico</p>
            <p className="text-sm font-semibold text-gray-900 dark:text-white truncate">{providerName}</p>
          </div>
        </div>
      )}

      {isSiigo && (
        <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
          <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">Identificadores de Siigo</h4>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Tipo de documento (FV)</label>
              <input type="number" value={formData.siigo_document_id} onChange={(e) => setFormData({ ...formData, siigo_document_id: e.target.value })} placeholder="document_id" disabled={loading} className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50" />
            </div>
            <div>
              <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Medio de pago (factura)</label>
              <input type="number" value={formData.siigo_payment_method_id} onChange={(e) => setFormData({ ...formData, siigo_payment_method_id: e.target.value })} placeholder="payment_method_id" disabled={loading} className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50" />
            </div>
            <div>
              <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Impuesto (IVA) — opcional</label>
              <input type="number" value={formData.siigo_tax_id} onChange={(e) => setFormData({ ...formData, siigo_tax_id: e.target.value })} placeholder="tax_id" disabled={loading} className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50" />
            </div>
            <div>
              <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Vendedor — opcional</label>
              <input type="number" value={formData.siigo_seller_id} onChange={(e) => setFormData({ ...formData, siigo_seller_id: e.target.value })} placeholder="seller_id" disabled={loading} className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50" />
            </div>
            <div>
              <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Tipo doc. Nota Crédito</label>
              <input type="number" value={formData.siigo_credit_note_document_id} onChange={(e) => setFormData({ ...formData, siigo_credit_note_document_id: e.target.value })} placeholder="credit_note_document_id" disabled={loading} className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50" />
            </div>
          </div>
          <p className="text-xs text-gray-400 mt-2">IDs del catálogo de Siigo (consúltalos en tu cuenta Siigo).</p>
        </div>
      )}

      {/* Facturación automática */}
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
        <label className="flex items-center gap-3 cursor-pointer">
          <input
            type="checkbox"
            checked={formData.auto_invoice}
            onChange={(e) => setFormData({ ...formData, auto_invoice: e.target.checked })}
            disabled={loading}
            className="w-5 h-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
          />
          <div>
            <span className="text-sm font-medium text-gray-900 dark:text-white">Facturación automática</span>
            <p className="text-xs text-gray-500 dark:text-gray-400">Las órdenes que cumplan los filtros se facturarán automáticamente</p>
          </div>
        </label>
      </div>

      {/* Filtros de Pago */}
      {formData.auto_invoice && (
        <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
          <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">Filtros de Pago</h4>
          <div>
            <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Estado de pago</label>
            <select
              value={formData.payment_status}
              onChange={(e) => setFormData({ ...formData, payment_status: e.target.value })}
              disabled={loading}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
            >
              <option value="">Sin filtro (todas las órdenes)</option>
              <option value="paid">Solo pagadas</option>
              <option value="unpaid">Solo sin pagar</option>
              <option value="partial">Pago parcial</option>
            </select>
          </div>
        </div>
      )}

      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
        <label className="flex items-center gap-3 cursor-pointer">
          <input
            type="checkbox"
            checked={formData.invoice_cod}
            onChange={(e) => setFormData({ ...formData, invoice_cod: e.target.checked })}
            disabled={loading}
            className="w-5 h-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
          />
          <div>
            <span className="text-sm font-medium text-gray-900 dark:text-white">Facturar contra entrega</span>
            <p className="text-xs text-gray-500 dark:text-gray-400">Permite facturar ordenes de pago contra entrega aunque no esten pagadas. Si esta desactivado, las contra entrega se bloquean.</p>
          </div>
        </label>
      </div>

      {/* Facturar como Consumidor Final */}
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
        <label className="flex items-center gap-3 cursor-pointer">
          <input
            type="checkbox"
            checked={formData.force_default_customer}
            onChange={(e) => setFormData({ ...formData, force_default_customer: e.target.checked })}
            disabled={loading}
            className="w-5 h-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
          />
          <div>
            <span className="text-sm font-medium text-gray-900 dark:text-white">Facturar como Consumidor Final</span>
            <p className="text-xs text-gray-500 dark:text-gray-400">Todas las facturas se generan a nombre de CONSUMIDOR FINAL (222222222222), sin importar los datos del cliente</p>
          </div>
        </label>
      </div>

      {/* Recibo de Caja */}
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
          <label className="flex items-center gap-3 cursor-pointer mb-3">
            <input
              type="checkbox"
              checked={formData.send_cash_receipt}
              onChange={(e) => setFormData({ ...formData, send_cash_receipt: e.target.checked })}
              disabled={loading}
              className="w-5 h-5 rounded border-gray-300 text-green-600 focus:ring-green-500 disabled:opacity-50"
            />
            <div>
              <span className="text-sm font-medium text-gray-900 dark:text-white">Enviar recibo de caja</span>
              <p className="text-xs text-gray-500 dark:text-gray-400">{cashReceiptDesc}</p>
            </div>
          </label>

          {formData.send_cash_receipt && isSiigo && (
            <div className="space-y-3 pl-8">
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Tipo doc. Recibo (RC)</label>
                  <input type="number" value={formData.siigo_cash_receipt_document_id} onChange={(e) => setFormData({ ...formData, siigo_cash_receipt_document_id: e.target.value })} placeholder="cash_receipt_document_id" disabled={loading} className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50" />
                </div>
                <div>
                  <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Medio de pago (Siigo)</label>
                  <input type="number" value={formData.siigo_cash_receipt_payment_id} onChange={(e) => setFormData({ ...formData, siigo_cash_receipt_payment_id: e.target.value })} placeholder="cash_receipt_payment_id" disabled={loading} className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50" />
                </div>
              </div>
              <button type="button" onClick={handleFetchBankAccounts} disabled={loadingBankAccounts} className="px-3 py-1.5 text-xs font-medium text-green-700 bg-green-50 border border-green-200 rounded-md hover:bg-green-100 disabled:opacity-50">
                {loadingBankAccounts ? 'Consultando...' : 'Consultar medios de pago de Siigo'}
              </button>
              {bankAccounts && bankAccounts.length > 0 && (
                <div className="space-y-1">
                  {bankAccounts.map((account, idx) => (
                    <button key={idx} type="button" onClick={() => setFormData({ ...formData, siigo_cash_receipt_payment_id: account.account_number })} className={`w-full text-left p-2 rounded text-xs border ${String(formData.siigo_cash_receipt_payment_id) === String(account.account_number) ? 'border-green-500 bg-green-50' : 'border-gray-200 dark:border-gray-700 hover:bg-gray-50'}`}>
                      <span className="font-medium">{account.account_number}</span>
                      <span className="text-gray-500 dark:text-gray-400 ml-2">{account.name}</span>
                    </button>
                  ))}
                </div>
              )}
            </div>
          )}

          {formData.send_cash_receipt && !isSiigo && (
            <div className="space-y-3 pl-8">
              <div>
                <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Medio de pago</label>
                <select
                  value={formData.payment_type}
                  onChange={(e) => setFormData({ ...formData, payment_type: e.target.value })}
                  disabled={loading}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
                >
                  <option value="EF">EF — Efectivo</option>
                  <option value="TR">TR — Transferencia bancaria</option>
                  <option value="TC">TC — Tarjeta de credito</option>
                  <option value="TD">TD — Tarjeta de debito</option>
                  <option value="CH">CH — Cheque</option>
                  <option value="BN">BN — Bonos</option>
                </select>
              </div>

              {/* TR: bankAccountId */}
              {formData.payment_type === 'TR' && (
                <div>
                  <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Numero de cuenta bancaria</label>
                  <input
                    type="text"
                    value={formData.payment_bank_account_id}
                    onChange={(e) => setFormData({ ...formData, payment_bank_account_id: e.target.value })}
                    placeholder="Ej: 1"
                    disabled={loading}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
                  />
                  <p className="text-xs text-gray-400 mt-1">
                    Numero de cuenta registrada en Softpymes (consultar en Utilidades → Buscar cuentas bancarias)
                  </p>
                  <button
                    type="button"
                    onClick={handleFetchBankAccounts}
                    disabled={loadingBankAccounts}
                    className="mt-2 px-3 py-1.5 text-xs font-medium text-green-700 bg-green-50 border border-green-200 rounded-md hover:bg-green-100 disabled:opacity-50"
                  >
                    {loadingBankAccounts ? 'Consultando...' : 'Consultar cuentas en Softpymes'}
                  </button>

                  {bankAccounts && bankAccounts.length > 0 && (
                    <div className="mt-2 space-y-1">
                      {bankAccounts.map((account, idx) => (
                        <button
                          key={idx}
                          type="button"
                          onClick={() => setFormData({ ...formData, payment_bank_account_id: account.account_number })}
                          className={`w-full text-left p-2 rounded text-xs border ${
                            formData.payment_bank_account_id === account.account_number
                              ? 'border-green-500 bg-green-50'
                              : 'border-gray-200 dark:border-gray-700 hover:bg-gray-50'
                          }`}
                        >
                          <span className="font-medium">{account.account_number}</span>
                          <span className="text-gray-500 dark:text-gray-400 ml-2">{account.name}</span>
                          <span className="text-gray-400 ml-1">({account.name_type})</span>
                        </button>
                      ))}
                    </div>
                  )}
                </div>
              )}

              {/* CH: accountNumber + bankName */}
              {formData.payment_type === 'CH' && (
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Numero de cuenta</label>
                    <input
                      type="text"
                      value={formData.payment_account_number}
                      onChange={(e) => setFormData({ ...formData, payment_account_number: e.target.value })}
                      placeholder="Numero de cuenta"
                      disabled={loading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
                    />
                  </div>
                  <div>
                    <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Nombre del banco</label>
                    <input
                      type="text"
                      value={formData.payment_bank_name}
                      onChange={(e) => setFormData({ ...formData, payment_bank_name: e.target.value })}
                      placeholder="Ej: Bancolombia"
                      disabled={loading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
                    />
                  </div>
                </div>
              )}

              {/* TC/TD: finantialEntityId */}
              {(formData.payment_type === 'TC' || formData.payment_type === 'TD') && (
                <div>
                  <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">ID Entidad Financiera (Softpymes)</label>
                  <input
                    type="number"
                    value={formData.payment_financial_entity_id}
                    onChange={(e) => setFormData({ ...formData, payment_financial_entity_id: e.target.value })}
                    placeholder="ID numerico de la entidad"
                    disabled={loading}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
                  />
                </div>
              )}

              <div className="border-t pt-3 mt-2">
                <label className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={formData.cod_use_alternate_bank}
                    onChange={(e) => setFormData({ ...formData, cod_use_alternate_bank: e.target.checked })}
                    disabled={loading}
                    className="w-5 h-5 rounded border-gray-300 text-green-600 focus:ring-green-500 disabled:opacity-50"
                  />
                  <div>
                    <span className="text-sm font-medium text-gray-900 dark:text-white">Usar cuenta alterna para contra entrega</span>
                    <p className="text-xs text-gray-500 dark:text-gray-400">Si esta activo, el recibo de caja de las ordenes contra entrega se registra en otra cuenta bancaria.</p>
                  </div>
                </label>

                {formData.cod_use_alternate_bank && (
                  <div className="mt-3">
                    <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Numero de cuenta bancaria contra entrega</label>
                    <input
                      type="text"
                      value={formData.cod_payment_bank_account_id}
                      onChange={(e) => setFormData({ ...formData, cod_payment_bank_account_id: e.target.value })}
                      placeholder="Ej: 2"
                      disabled={loading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
                    />
                    <p className="text-xs text-gray-400 mt-1">
                      Numero de cuenta registrada en Softpymes que se usara solo para ordenes contra entrega.
                    </p>
                  </div>
                )}
              </div>

              {/* BN: code */}
              {formData.payment_type === 'BN' && (
                <div>
                  <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Codigo del bono</label>
                  <input
                    type="text"
                    value={formData.payment_bonus_code}
                    onChange={(e) => setFormData({ ...formData, payment_bonus_code: e.target.value })}
                    placeholder="Codigo identificador"
                    disabled={loading}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50"
                  />
                </div>
              )}
            </div>
          )}
        </div>

      {/* Mapeo de Servicios */}
        <div className="bg-purple-50 p-4 rounded-lg border border-purple-100">
          <div className="flex items-center justify-between">
            <div>
              <span className="text-sm font-medium text-gray-900 dark:text-white">Mapeo de servicios</span>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
                Probability traduce internamente conceptos como envío, membresía y propina. Aquí defines con qué código se facturan.
              </p>
            </div>
            <button
              type="button"
              onClick={() => setShowItemMappings(!showItemMappings)}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none flex-shrink-0 ml-4 ${showItemMappings ? 'bg-purple-500' : 'bg-gray-200'}`}
            >
              <span className={`inline-block h-4 w-4 transform rounded-full bg-white dark:bg-gray-800 transition-transform ${showItemMappings ? 'translate-x-6' : 'translate-x-1'}`} />
            </button>
          </div>

          {showItemMappings && (
            <div className="space-y-3 mt-4">
              <div>
                <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Servicio de Envío</label>
                <input
                  type="text"
                  value={formData.item_mappings_shipping}
                  onChange={(e) => setFormData({ ...formData, item_mappings_shipping: e.target.value })}
                  placeholder="Ej: SE02001 (vacío = SHIPPING)"
                  disabled={loading}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-purple-500 disabled:opacity-50"
                />
                <p className="text-xs text-gray-400 mt-1">Código con el que se factura el costo de envío</p>
              </div>
              <div>
                <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Servicio de Membresía</label>
                <input
                  type="text"
                  value={formData.item_mappings_membership}
                  onChange={(e) => setFormData({ ...formData, item_mappings_membership: e.target.value })}
                  placeholder="Ej: SE01001"
                  disabled={loading}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-purple-500 disabled:opacity-50"
                />
                <p className="text-xs text-gray-400 mt-1">Código con el que se facturan las membresías</p>
              </div>
              <div>
                <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Servicio de Propina</label>
                <input
                  type="text"
                  value={formData.item_mappings_tip}
                  onChange={(e) => setFormData({ ...formData, item_mappings_tip: e.target.value })}
                  placeholder="Ej: SE03001"
                  disabled={loading}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-purple-500 disabled:opacity-50"
                />
                <p className="text-xs text-gray-400 mt-1">Código con el que se facturan las propinas</p>
              </div>
            </div>
          )}
        </div>

      {/* Integraciones de origen de órdenes */}
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg border border-gray-200 dark:border-gray-700">
        <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-1">
          Fuentes de órdenes
        </h4>
        <p className="text-xs text-gray-500 dark:text-gray-400 mb-3">
          Selecciona las integraciones desde las cuales se facturarán las órdenes
        </p>

        {loadingIntegrations ? (
          <p className="text-xs text-gray-400">Cargando integraciones...</p>
        ) : availableIntegrations.length === 0 ? (
          <p className="text-xs text-gray-400">No hay integraciones e-commerce disponibles</p>
        ) : (
          <div className="space-y-2">
            {availableIntegrations.map((integration) => (
              <label
                key={integration.id}
                className="flex items-center gap-3 cursor-pointer p-2 rounded-md hover:bg-gray-50"
              >
                <input
                  type="checkbox"
                  checked={selectedIntegrationIds.includes(integration.id)}
                  onChange={() => toggleIntegration(integration.id)}
                  disabled={loading}
                  className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
                />
                <div className="flex items-center gap-2 min-w-0">
                  {integration.integration_type?.image_url && (
                    <img
                      src={integration.integration_type.image_url}
                      alt={integration.name}
                      className="w-5 h-5 object-contain flex-shrink-0"
                    />
                  )}
                  <span className="text-sm text-gray-800 dark:text-gray-100 truncate">{integration.name}</span>
                  {integration.type && (
                    <span className="text-xs text-gray-400 flex-shrink-0">({integration.type})</span>
                  )}
                </div>
              </label>
            ))}
          </div>
        )}
      </div>

      {/* Acciones */}
      <div className="flex items-center gap-3 pt-2 border-t">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            disabled={loading}
            className="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-800 border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
          >
            Cancelar
          </button>
        )}
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? 'Guardando...' : initialData?.id ? 'Actualizar' : 'Crear Configuración'}
        </button>
      </div>
    </form>
  );
}
