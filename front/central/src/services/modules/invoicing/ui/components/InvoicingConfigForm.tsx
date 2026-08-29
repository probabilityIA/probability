'use client';

import { useState, useEffect, FormEvent } from 'react';
import type {
  InvoicingConfig,
  CreateConfigDTO,
  BankAccountResult,
} from '@/services/modules/invoicing/domain/types';
import type { Integration } from '@/services/integrations/core/domain/types';
import { getIntegrationsAction } from '@/services/integrations/core/infra/actions';
import { getSiigoCatalogsAction, type SiigoCatalogs } from '@/services/integrations/invoicing/siigo/infra/actions';
import { CampoCatalogo } from './CampoCatalogo';
import { CampoServicio } from './CampoServicio';
import { FilaSwitch } from './FilaSwitch';
import { useInvoicingConfig } from '@/services/modules/invoicing/ui/hooks/useInvoicingConfig';
import {
  requestListBankAccountsAction,
  getListBankAccountsResultAction,
} from '@/services/modules/invoicing/infra/actions';

interface ProveedorFacturador {
  nombre: string;
  imagenUrl?: string;
}

interface InvoicingConfigFormProps {
  integrationIds: number[];
  invoicingIntegrationId: number;
  businessId: number;
  proveedor?: ProveedorFacturador;
  alturaFija?: boolean;
  onSuccess?: () => void;
  onCancel?: () => void;
  initialData?: InvoicingConfig;
}

export function InvoicingConfigForm({
  integrationIds,
  invoicingIntegrationId,
  businessId,
  proveedor,
  alturaFija,
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
    item_mappings_shipping: (initialData?.config?.item_mappings?.shipping as string) ?? '',
    item_mappings_membership: (initialData?.config?.item_mappings?.membership as string) ?? '',
    item_mappings_tip: (initialData?.config?.item_mappings?.tip as string) ?? '',
    siigo_document_id: (initialData?.config?.document_id as number | string) ?? '',
    siigo_payment_method_id: (initialData?.config?.payment_method_id as number | string) ?? '',
    siigo_tax_id: (initialData?.config?.tax_id as number | string) ?? '',
    siigo_seller_id: (initialData?.config?.seller_id as number | string) ?? '',
    siigo_cash_receipt_document_id: (initialData?.config?.cash_receipt_document_id as number | string) ?? '',
    siigo_cash_receipt_payment_id: (initialData?.config?.cash_receipt_payment_id as number | string) ?? '',
    siigo_credit_note_document_id: (initialData?.config?.credit_note_document_id as number | string) ?? '',
    final_customer_when_no_id: initialData?.config?.final_customer_when_no_id ?? false,
    inventory_exit_only: initialData?.config?.inventory_exit_only ?? false,
    inventory_exit_document_id: (initialData?.config?.inventory_exit_document_id as number | string) ?? '',
    inventory_exit_account_code: (initialData?.config?.inventory_exit_account_code as string) ?? '',
    inventory_exit_offset_account_code: (initialData?.config?.inventory_exit_offset_account_code as string) ?? '',
    inventory_exit_warehouse_id: (initialData?.config?.inventory_exit_warehouse_id as number | string) ?? '',
  });

  const [catalogs, setCatalogs] = useState<SiigoCatalogs | null>(null);
  const [loadingCatalogs, setLoadingCatalogs] = useState(false);
  const [catalogsError, setCatalogsError] = useState<string | null>(null);

  const providerName = initialData?.provider_name ?? proveedor?.nombre ?? '';
  const providerImageUrl = initialData?.provider_image_url ?? proveedor?.imagenUrl;
  const isSiigo = providerName.toLowerCase().includes('siigo') || (providerImageUrl ?? '').toLowerCase().includes('siigo');
  const cashReceiptDesc = isSiigo
    ? 'Registra un recibo de caja en Siigo al crear la factura'
    : 'Registra el pago en Softpymes al crear la factura (mueve cuentas por cobrar al medio de pago)';

  useEffect(() => {
    if (!isSiigo) return;
    if (!invoicingIntegrationId) {
      setCatalogsError('No se pudo identificar la integraci\u00f3n de Siigo de esta configuraci\u00f3n');
      return;
    }
    let cancelado = false;
    setLoadingCatalogs(true);
    setCatalogsError(null);
    getSiigoCatalogsAction(invoicingIntegrationId)
      .then(res => {
        if (cancelado) return;
        if (res.success && res.data) {
          setCatalogs(res.data);
          if (res.data.errors?.length) setCatalogsError(res.data.errors.join(' | '));
        } else {
          setCatalogsError(res.message ?? 'No se pudieron cargar los catalogos de Siigo');
        }
      })
      .catch(() => { if (!cancelado) setCatalogsError('No se pudieron cargar los catalogos de Siigo'); })
      .finally(() => { if (!cancelado) setLoadingCatalogs(false); });
    return () => { cancelado = true; };
  }, [isSiigo, invoicingIntegrationId]);

  const [showItemMappings, setShowItemMappings] = useState(
    !!(initialData?.config?.item_mappings?.shipping ||
       initialData?.config?.item_mappings?.membership ||
       initialData?.config?.item_mappings?.tip)
  );


  const initialSelected = initialData?.integration_ids?.length
    ? initialData.integration_ids
    : integrationIds;
  const [selectedIntegrationIds, setSelectedIntegrationIds] = useState<number[]>(initialSelected);
  const [availableIntegrations, setAvailableIntegrations] = useState<Integration[]>([]);
  const [loadingIntegrations, setLoadingIntegrations] = useState(false);

  const [error, setError] = useState<string | null>(null);
  const [loadingBankAccounts, setLoadingBankAccounts] = useState(false);
  const [bankAccounts, setBankAccounts] = useState<BankAccountResult[] | null>(null);
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
      setError('Debes seleccionar al menos una integraci\u00f3n de origen de \u00f3rdenes');
      return;
    }

    const filters: Record<string, any> = {};
    if (formData.payment_status) {
      filters.payment_status = formData.payment_status as 'paid' | 'unpaid' | 'partial';
    }
    if (formData.invoice_cod) {
      filters.invoice_cod = true;
    }
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
      if (formData.final_customer_when_no_id) {
        invoiceConfig.final_customer_when_no_id = true;
      }
      if (formData.inventory_exit_only) {
        invoiceConfig.inventory_exit_only = true;
        if (formData.inventory_exit_document_id) invoiceConfig.inventory_exit_document_id = Number(formData.inventory_exit_document_id);
        if (formData.inventory_exit_account_code) invoiceConfig.inventory_exit_account_code = formData.inventory_exit_account_code;
        if (formData.inventory_exit_offset_account_code) invoiceConfig.inventory_exit_offset_account_code = formData.inventory_exit_offset_account_code;
        if (formData.inventory_exit_warehouse_id) invoiceConfig.inventory_exit_warehouse_id = Number(formData.inventory_exit_warehouse_id);
      }
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
          setError(result.error || 'Error al actualizar configuraci\u00f3n');
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
          setError(result.error || 'Error al crear la configuraci\u00f3n');
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error desconocido');
    }
  };

  const claseForm = alturaFija ? 'flex h-full min-h-0 flex-col' : 'space-y-4';
  const claseEncabezado = alturaFija
    ? 'shrink-0 space-y-3 border-b border-gray-200 px-6 py-4 dark:border-gray-700'
    : 'space-y-4';
  const claseCuerpo = alturaFija ? 'min-h-0 flex-1 overflow-y-auto px-6 py-4' : '';
  const mostrarEncabezado = !!error || (!!providerName && !alturaFija);
  const clasePie = alturaFija
    ? 'shrink-0 flex items-center justify-end gap-3 border-t border-gray-200 bg-white px-6 py-3 dark:border-gray-700 dark:bg-gray-800'
    : 'flex items-center justify-end gap-3 pt-2 border-t';

  return (
    <form onSubmit={handleSubmit} className={claseForm}>
      {mostrarEncabezado && (
      <div className={claseEncabezado}>
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-sm text-red-800">{error}</p>
        </div>
      )}

      {providerName && !alturaFija && (
        <div
          className="flex items-center gap-3 p-3 rounded-lg border dark:bg-gray-800 dark:border-gray-700"
          style={{
            backgroundColor: 'color-mix(in srgb, var(--color-primary) 8%, white)',
            borderColor: 'color-mix(in srgb, var(--color-primary) 20%, white)',
          }}>
          {providerImageUrl ? (
            <img src={providerImageUrl} alt={providerName} className="w-8 h-8 object-contain flex-shrink-0" />
          ) : (
            <div
              className="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold flex-shrink-0"
              style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, #fff)' }}>
              {providerName.charAt(0)}
            </div>
          )}
          <div className="min-w-0">
            <p className="text-xs text-gray-500 dark:text-gray-400">Facturador {'electr\u00f3nico'}</p>
            <p className="text-sm font-semibold text-gray-900 dark:text-white truncate">{providerName}</p>
          </div>
        </div>
      )}

      </div>
      )}

      <div className={claseCuerpo}>
      <div className="grid grid-cols-1 gap-4 items-start lg:grid-cols-2">

      {isSiigo && (
        <div className="bg-[#fafafd] dark:bg-gray-800/60 p-4 rounded-xl border border-[#eceaf3] dark:border-gray-700 lg:col-span-2">
          <div className="mb-3 flex items-center justify-between">
            <h4 className="text-sm font-medium text-gray-900 dark:text-white">Identificadores de Siigo</h4>
            {loadingCatalogs && <span className="text-xs text-gray-400">Cargando desde Siigo...</span>}
          </div>

          {catalogsError && (
            <p className="mb-3 rounded bg-amber-50 px-2 py-1.5 text-[11px] leading-snug text-amber-700 dark:bg-amber-900/20 dark:text-amber-300">
              {catalogsError}. Puedes escribir los IDs a mano si conoces sus valores.
            </p>
          )}

          <div className="grid grid-cols-2 gap-3">
            <CampoCatalogo
              etiqueta="Tipo de documento (FV)"
              requerido
              ayuda="La resoluci\u00f3n de facturaci\u00f3n de venta con la que la DIAN te autorizo numerar. En Siigo esta en Configuraci\u00f3n > Tipos de comprobante > Factura de venta."
              opciones={catalogs?.document_types_fv}
              valor={formData.siigo_document_id}
              onChange={(v) => setFormData({ ...formData, siigo_document_id: v })}
              disabled={loading}
            />
            <CampoCatalogo
              etiqueta="Medio de pago (factura)"
              requerido
              ayuda="Con que forma de pago queda registrada la factura en Siigo. Es el valor por defecto: se aplica a todas las facturas de esta configuraci\u00f3n."
              opciones={catalogs?.payment_types_fv}
              valor={formData.siigo_payment_method_id}
              onChange={(v) => setFormData({ ...formData, siigo_payment_method_id: v })}
              disabled={loading}
            />
            <CampoCatalogo
              etiqueta="Vendedor"
              requerido
              ayuda="El vendedor al que Siigo le atribuye la venta. Si manejas comisiones, elige el que corresponda a las ventas por canal digital."
              opciones={catalogs?.sellers}
              valor={formData.siigo_seller_id}
              onChange={(v) => setFormData({ ...formData, siigo_seller_id: v })}
              disabled={loading}
            />
            <CampoCatalogo
              etiqueta="Impuesto (IVA)"
              ayuda="Se aplica solo a los productos que en la orden ya traen impuesto. Si lo dejas vacio, las facturas salen sin IVA."
              opciones={catalogs?.taxes}
              valor={formData.siigo_tax_id}
              onChange={(v) => setFormData({ ...formData, siigo_tax_id: v })}
              disabled={loading}
            />
            <CampoCatalogo
              etiqueta="Tipo doc. Nota Cr\u00e9dito"
              ayuda="Necesario para anular. Sin el, una factura ya emitida no se puede reversar desde Probability."
              opciones={catalogs?.document_types_nc}
              valor={formData.siigo_credit_note_document_id}
              onChange={(v) => setFormData({ ...formData, siigo_credit_note_document_id: v })}
              disabled={loading}
            />
          </div>
          <p className="text-xs text-gray-400 mt-2">Se leen de la cuenta Siigo del negocio. Los tres primeros son obligatorios: sin ellos Siigo rechaza la factura.</p>
        </div>
      )}
      <div className="bg-[#fafafd] dark:bg-gray-800/60 p-4 rounded-xl border border-[#eceaf3] dark:border-gray-700 lg:col-span-2">
        <h4 className="text-sm font-medium text-gray-900 dark:text-white">Opciones de {'facturaci\u00f3n'}</h4>
        <p className="mt-0.5 text-xs text-gray-500 dark:text-gray-400">Que se factura, a nombre de quien y que se registra {'adem\u00e1s'} de la factura.</p>

        <div className="mt-2 divide-y divide-gray-100 dark:divide-gray-700">
          <FilaSwitch
            titulo="Facturaci\u00f3n autom\u00e1tica"
            descripcion={formData.inventory_exit_only
              ? 'Desactivada porque la salida de inventario sin facturar est\u00e1 activa'
              : 'Las \u00f3rdenes que cumplan los filtros se facturaran autom\u00e1ticamente'}
            checked={formData.auto_invoice}
            onToggle={(v) => setFormData({ ...formData, auto_invoice: v, inventory_exit_only: v ? false : formData.inventory_exit_only })}
            disabled={loading || formData.inventory_exit_only}
          >
            {formData.auto_invoice && (
              <div className="mt-3 border-t border-gray-100 pt-3 dark:border-gray-700 sm:max-w-md">
            <div>
              <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">Estado de pago</label>
              <p className="mb-1 text-[11px] leading-snug text-gray-400">Que {'\u00f3rdenes'} entran a la {'facturaci\u00f3n'} {'autom\u00e1tica'} {'seg\u00fan'} como esten pagadas. Las contra entrega se controlan aparte, con el interruptor de {'m\u00e1s'} abajo.</p>
              <select
                value={formData.payment_status}
                onChange={(e) => setFormData({ ...formData, payment_status: e.target.value })}
                disabled={loading}
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50"
              >
                <option value="">Sin filtro (todas las {'\u00f3rdenes'})</option>
                <option value="paid">Solo pagadas</option>
                <option value="unpaid">Solo sin pagar</option>
                <option value="partial">Pago parcial</option>
              </select>
            </div>
              </div>
            )}
          </FilaSwitch>

          {isSiigo && (
            <FilaSwitch
              titulo="Salida de inventario sin facturar"
              descripcion={formData.auto_invoice
                ? 'Desactivada porque la facturaci\u00f3n autom\u00e1tica est\u00e1 activa'
                : 'Descarga el inventario en Siigo con un comprobante contable, sin emitir factura. Excluyente con la facturaci\u00f3n autom\u00e1tica para evitar el doble descargue.'}
              checked={formData.inventory_exit_only}
              onToggle={(v) => setFormData({ ...formData, inventory_exit_only: v, auto_invoice: v ? false : formData.auto_invoice })}
              disabled={loading || formData.auto_invoice}
            >
            {formData.inventory_exit_only && (
              <div className="mt-3 grid grid-cols-2 gap-3 border-t border-gray-100 pt-3 dark:border-gray-700">
                <div>
                  <label className="block text-xs text-gray-500 dark:text-gray-400 mb-1">Tipo doc. Comprobante (CC)</label>
                  <input type="number" value={formData.inventory_exit_document_id} onChange={(e) => setFormData({ ...formData, inventory_exit_document_id: e.target.value })} placeholder="document_id" disabled={loading} className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50" />
                  <p className="mt-1 text-[11px] leading-snug text-gray-400">El comprobante contable donde queda la salida. En Siigo: {'Configuraci\u00f3n'} &gt; Tipos de comprobante &gt; Comprobante contable.</p>
                </div>
                <div>
                  <label className="block text-xs text-gray-500 dark:text-gray-400 mb-1">Bodega - opcional</label>
                  <input type="number" value={formData.inventory_exit_warehouse_id} onChange={(e) => setFormData({ ...formData, inventory_exit_warehouse_id: e.target.value })} placeholder="warehouse_id" disabled={loading} className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50" />
                  <p className="mt-1 text-[11px] leading-snug text-gray-400">De cual bodega de Siigo se descuenta. Vacio = la bodega por defecto de la cuenta.</p>
                </div>
                <div>
                  <label className="block text-xs text-gray-500 dark:text-gray-400 mb-1">Cuenta de inventario ({'cr\u00e9dito'})</label>
                  <input type="text" value={formData.inventory_exit_account_code} onChange={(e) => setFormData({ ...formData, inventory_exit_account_code: e.target.value })} placeholder="14350501" disabled={loading} className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50" />
                  <p className="mt-1 text-[11px] leading-snug text-gray-400">La cuenta del PUC de donde sale la {'mercanc\u00eda'}. Suele empezar en 1435.</p>
                </div>
                <div>
                  <label className="block text-xs text-gray-500 dark:text-gray-400 mb-1">Contrapartida ({'d\u00e9bito'})</label>
                  <input type="text" value={formData.inventory_exit_offset_account_code} onChange={(e) => setFormData({ ...formData, inventory_exit_offset_account_code: e.target.value })} placeholder="61350501" disabled={loading} className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50" />
                  <p className="mt-1 text-[11px] leading-snug text-gray-400">Contra que cuenta se registra la salida, normalmente el costo de venta (61xx). Confirmalo con tu contador.</p>
                </div>
                <p className="col-span-2 text-xs text-gray-400">Los {'c\u00f3digos'} de cuenta los define el contador y deben existir y estar activos en Siigo.</p>
              </div>
            )}
            </FilaSwitch>
          )}

          <FilaSwitch
            titulo="Facturar contra entrega"
            descripcion="Permite facturar \u00f3rdenes de pago contra entrega aunque no esten pagadas. Si esta desactivado, las contra entrega se bloquean."
            checked={formData.invoice_cod}
            onToggle={(v) => setFormData({ ...formData, invoice_cod: v })}
            disabled={loading}
          />

          {isSiigo && (
            <FilaSwitch
              titulo="Facturar a consumidor final sin cedula"
              descripcion="Si la orden no trae cedula, NIT ni documento del cliente, la factura se emite al tercero generico CONSUMIDOR FINAL. Siigo exige identificaci\u00f3n siempre."
              checked={formData.final_customer_when_no_id}
              onToggle={(v) => setFormData({ ...formData, final_customer_when_no_id: v })}
              disabled={loading}
            >
            {!formData.final_customer_when_no_id && (
              <p className="mt-2 rounded bg-amber-50 px-2 py-1.5 text-[11px] leading-snug text-amber-700 dark:bg-amber-900/20 dark:text-amber-300">
                Desactivado: las ordenes sin documento del cliente quedaran en estado Fallida, con el motivo en el historial de sincronizacion. Hay que completar el documento en la orden y reintentar.
              </p>
            )}
            {formData.final_customer_when_no_id && (
              <p className="mt-2 rounded px-2 py-1.5 text-[11px] leading-snug" style={{ backgroundColor: 'color-mix(in srgb, var(--color-primary) 10%, white)', color: 'color-mix(in srgb, var(--color-primary) 80%, black)' }}>
                Activado: esas facturas no identifican al comprador y salen a nombre de CONSUMIDOR FINAL, no del cliente real.
              </p>
            )}
            </FilaSwitch>
          )}

          <FilaSwitch
            titulo="Facturar como Consumidor Final"
            descripcion="Todas las facturas se generan a nombre de CONSUMIDOR FINAL (222222222222), sin importar los datos del cliente"
            checked={formData.force_default_customer}
            onToggle={(v) => setFormData({ ...formData, force_default_customer: v })}
            disabled={loading}
          />

          <FilaSwitch
            titulo="Enviar recibo de caja"
            descripcion={cashReceiptDesc}
            checked={formData.send_cash_receipt}
            onToggle={(v) => setFormData({ ...formData, send_cash_receipt: v })}
            disabled={loading}
          >
            {formData.send_cash_receipt && isSiigo && (
              <div className="space-y-3 pl-8">
                <div className="grid grid-cols-2 gap-3">
                  <CampoCatalogo
                    etiqueta="Tipo doc. Recibo (RC)"
                    requerido
                    ayuda="El comprobante de recibo de caja en Siigo. Es distinto del de la factura: Configuraci\u00f3n > Tipos de comprobante > Recibo de caja."
                    opciones={catalogs?.document_types_rc}
                    valor={formData.siigo_cash_receipt_document_id}
                    onChange={(v) => setFormData({ ...formData, siigo_cash_receipt_document_id: v })}
                    disabled={loading}
                  />
                  <CampoCatalogo
                    etiqueta="Medio de pago (Siigo)"
                    requerido
                    ayuda="La cuenta o caja donde entra la plata. Usa el bot\u00f3n de abajo para traer las cuentas reales de Siigo y elegir una."
                    opciones={catalogs?.payment_types_rc}
                    valor={formData.siigo_cash_receipt_payment_id}
                    onChange={(v) => setFormData({ ...formData, siigo_cash_receipt_payment_id: v })}
                    disabled={loading}
                  />
                </div>
                <button type="button" onClick={handleFetchBankAccounts} disabled={loadingBankAccounts} className="px-3 py-1.5 text-xs font-medium rounded-md border disabled:opacity-50"
                style={{ color: 'var(--color-primary)', borderColor: 'color-mix(in srgb, var(--color-primary) 30%, white)', backgroundColor: 'color-mix(in srgb, var(--color-primary) 8%, white)' }}>
                  {loadingBankAccounts ? 'Consultando...' : 'Consultar medios de pago de Siigo'}
                </button>
                {bankAccounts && bankAccounts.length > 0 && (
                  <div className="space-y-1">
                    {bankAccounts.map((account, idx) => (
                      <button key={idx} type="button" onClick={() => setFormData({ ...formData, siigo_cash_receipt_payment_id: account.account_number })} className={`w-full text-left p-2 rounded text-xs border ${String(formData.siigo_cash_receipt_payment_id) === String(account.account_number) ? 'border-[var(--color-primary)] bg-[color-mix(in_srgb,var(--color-primary)_10%,white)]' : 'border-gray-200 dark:border-gray-700 hover:bg-gray-50'}`}>
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
                    className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50"
                  >
                    <option value="EF">EF - Efectivo</option>
                    <option value="TR">TR - Transferencia bancaria</option>
                    <option value="TC">TC - Tarjeta de {'cr\u00e9dito'}</option>
                    <option value="TD">TD - Tarjeta de {'d\u00e9bito'}</option>
                    <option value="CH">CH - Cheque</option>
                    <option value="BN">BN - Bonos</option>
                  </select>
                </div>
                {formData.payment_type === 'TR' && (
                  <div>
                    <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">{'N\u00famero'} de cuenta bancaria</label>
                    <input
                      type="text"
                      value={formData.payment_bank_account_id}
                      onChange={(e) => setFormData({ ...formData, payment_bank_account_id: e.target.value })}
                      placeholder="Ej: 1"
                      disabled={loading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50"
                    />
                    <p className="text-xs text-gray-400 mt-1">
                      Numero de cuenta registrada en Softpymes (consultar en Utilidades &gt; Buscar cuentas bancarias)
                    </p>
                    <button
                      type="button"
                      onClick={handleFetchBankAccounts}
                      disabled={loadingBankAccounts}
                      className="mt-2 px-3 py-1.5 text-xs font-medium rounded-md border disabled:opacity-50"
                    style={{ color: 'var(--color-primary)', borderColor: 'color-mix(in srgb, var(--color-primary) 30%, white)', backgroundColor: 'color-mix(in srgb, var(--color-primary) 8%, white)' }}
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
                                ? 'border-[var(--color-primary)] bg-[color-mix(in_srgb,var(--color-primary)_10%,white)]'
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
                {formData.payment_type === 'CH' && (
                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">{'N\u00famero'} de cuenta</label>
                      <input
                        type="text"
                        value={formData.payment_account_number}
                        onChange={(e) => setFormData({ ...formData, payment_account_number: e.target.value })}
                        placeholder="N\u00famero de cuenta"
                        disabled={loading}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50"
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
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50"
                      />
                    </div>
                  </div>
                )}
                {(formData.payment_type === 'TC' || formData.payment_type === 'TD') && (
                  <div>
                    <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">ID Entidad Financiera (Softpymes)</label>
                    <input
                      type="number"
                      value={formData.payment_financial_entity_id}
                      onChange={(e) => setFormData({ ...formData, payment_financial_entity_id: e.target.value })}
                      placeholder="ID num\u00e9rico de la entidad"
                      disabled={loading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50"
                    />
                  </div>
                )}

                <div className="border-t pt-1 mt-2">
                  <FilaSwitch
                    titulo="Usar cuenta alterna para contra entrega"
                    descripcion="Si est\u00e1 activo, el recibo de caja de las \u00f3rdenes contra entrega se registra en otra cuenta bancaria."
                    checked={formData.cod_use_alternate_bank}
                    onToggle={(v) => setFormData({ ...formData, cod_use_alternate_bank: v })}
                    disabled={loading}
                  />

                  {formData.cod_use_alternate_bank && (
                    <div className="mt-3">
                      <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">{'N\u00famero'} de cuenta bancaria contra entrega</label>
                      <input
                        type="text"
                        value={formData.cod_payment_bank_account_id}
                        onChange={(e) => setFormData({ ...formData, cod_payment_bank_account_id: e.target.value })}
                        placeholder="Ej: 2"
                        disabled={loading}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50"
                      />
                      <p className="text-xs text-gray-400 mt-1">
                        Numero de cuenta registrada en Softpymes que se usara solo para ordenes contra entrega.
                      </p>
                    </div>
                  )}
                </div>
                {formData.payment_type === 'BN' && (
                  <div>
                    <label className="block text-sm text-gray-700 dark:text-gray-200 mb-1">{'C\u00f3digo'} del bono</label>
                    <input
                      type="text"
                      value={formData.payment_bonus_code}
                      onChange={(e) => setFormData({ ...formData, payment_bonus_code: e.target.value })}
                      placeholder="C\u00f3digo identificador"
                      disabled={loading}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] disabled:opacity-50"
                    />
                  </div>
                )}
              </div>
            )}
          </FilaSwitch>
        </div>
      </div>

      
        <div className="bg-[#fafafd] dark:bg-gray-800/60 p-4 rounded-xl border border-[#eceaf3] dark:border-gray-700 lg:col-span-2">
          <div className="flex items-center justify-between">
            <div>
              <span className="text-sm font-medium text-gray-900 dark:text-white">Mapeo de servicios</span>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
                Probability traduce internamente conceptos como envio, membresia y propina. Aqui defines con que codigo se facturan.
              </p>
            </div>
            <button
              type="button"
              onClick={() => setShowItemMappings(!showItemMappings)}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none flex-shrink-0 ml-4 ${showItemMappings ? '' : 'bg-gray-200'}`}
              style={{ backgroundColor: showItemMappings ? 'var(--color-primary)' : undefined }}
            >
              <span className={`inline-block h-4 w-4 transform rounded-full bg-white dark:bg-gray-800 transition-transform ${showItemMappings ? 'translate-x-6' : 'translate-x-1'}`} />
            </button>
          </div>

          {showItemMappings && (
            <div className="mt-4 space-y-3">
              <p className="rounded bg-amber-50 px-2 py-1.5 text-[11px] leading-snug text-amber-700 dark:bg-amber-900/20 dark:text-amber-300">
                {isSiigo
                  ? 'Busca por nombre o c\u00f3digo la cuenta que ya tienes creada en tu facturador. Lo recomendado es que sea un servicio dedicado al flete. Si el c\u00f3digo no existe alla, el facturador rechaza la factura completa, no solo esa l\u00ednea.'
                  : 'Estos c\u00f3digos tienen que existir antes en tu facturador, creados como producto de tipo servicio. No los creamos nosotros. Si el c\u00f3digo no existe, el facturador rechaza la factura completa, no solo esa l\u00ednea.'}
              </p>
              <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              <CampoServicio
                etiqueta="Servicio de Env\u00edo"
                valor={formData.item_mappings_shipping}
                onChange={(v) => setFormData({ ...formData, item_mappings_shipping: v })}
                disabled={loading}
                placeholder="Ej: SE02001"
                ayuda="El costo del env\u00edo se factura como una l\u00ednea aparte, con este c\u00f3digo y cantidad 1. Se agrega sola cuando la orden trae flete."
                advertencia="Sin elegir: mandaremos el c\u00f3digo literal SHIPPING. Si no existe un servicio con ese c\u00f3digo exacto, toda orden con env\u00edo va a fallar al facturar."
                integrationId={invoicingIntegrationId}
                buscable={isSiigo}
              />
              <CampoServicio
                etiqueta="Servicio de Membresia"
                valor={formData.item_mappings_membership}
                onChange={(v) => setFormData({ ...formData, item_mappings_membership: v })}
                disabled={loading}
                placeholder="Ej: SE01001"
                ayuda="Solo aplica si vendes membresias o suscripciones. Si no las manejas, dejalo vacio."
                integrationId={invoicingIntegrationId}
                buscable={isSiigo}
              />
              <CampoServicio
                etiqueta="Servicio de Propina"
                valor={formData.item_mappings_tip}
                onChange={(v) => setFormData({ ...formData, item_mappings_tip: v })}
                disabled={loading}
                placeholder="Ej: SE03001"
                ayuda="Solo aplica si tus \u00f3rdenes traen propina. Si no las manejas, dejalo vacio."
                integrationId={invoicingIntegrationId}
                buscable={isSiigo}
              />
              </div>
            </div>
          )}
        </div>
      <div className="bg-[#fafafd] dark:bg-gray-800/60 p-4 rounded-xl border border-[#eceaf3] dark:border-gray-700 lg:col-span-2">
        <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-1">
          Fuentes de ordenes
        </h4>
        <p className="text-xs text-gray-500 dark:text-gray-400 mb-3">
          Selecciona las integraciones desde las cuales se facturaran las ordenes
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
                  className="w-4 h-4 rounded border-gray-300 text-[var(--color-primary)] focus:ring-[var(--color-primary)] disabled:opacity-50"
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
      </div>

      </div>

      <div className={clasePie}>
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
          className="px-4 py-2 text-sm font-medium rounded-md disabled:opacity-50 disabled:cursor-not-allowed"
          style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, #fff)' }}
        >
          {loading ? 'Guardando...' : initialData?.id ? 'Actualizar' : 'Crear Configuraci\u00f3n'}
        </button>
      </div>
    </form>
  );
}
