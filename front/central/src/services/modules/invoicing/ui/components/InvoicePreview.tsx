'use client';

import type { InvoicePreview as InvoicePreviewData } from '../../domain/invoice-preview';

interface InvoicePreviewProps {
  data: InvoicePreviewData;
  raw?: Record<string, any> | null;
  copySlot?: React.ReactNode;
}

const moneda = (v: number | string | undefined) => {
  if (v === undefined || v === null || v === '') return null;
  const n = Number(v);
  if (!Number.isFinite(n)) return String(v);
  return new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 2 }).format(n);
};

const cantidad = (v: number | string | undefined) => {
  if (v === undefined || v === null || v === '') return null;
  const n = Number(v);
  if (!Number.isFinite(n)) return String(v);
  return new Intl.NumberFormat('es-CO', { maximumFractionDigits: 2 }).format(n);
};

const Campo = ({ etiqueta, valor, mono }: { etiqueta: string; valor?: string | null; mono?: boolean }) => {
  if (!valor) return null;
  return (
    <div>
      <p className="text-xs text-gray-500 dark:text-gray-400 mb-0.5">{etiqueta}</p>
      <p className={`text-xs font-medium text-gray-700 dark:text-gray-200 ${mono ? 'font-mono break-all' : ''}`}>{valor}</p>
    </div>
  );
};

export function InvoicePreview({ data, raw, copySlot }: InvoicePreviewProps) {
  const timbradoPendiente = !!data.stampStatus && !data.cufe;

  return (
    <div className="mb-6 p-4 bg-blue-50 dark:bg-blue-900/10 border border-blue-200 dark:border-blue-800 rounded-lg">
      <p className="text-xs text-blue-600 uppercase tracking-wide font-semibold mb-3">Vista Previa de Factura</p>
      <details className="group" open>
        <summary className="text-xs text-gray-700 dark:text-gray-200 cursor-pointer hover:text-blue-600 font-medium flex items-center gap-2">
          <span>Detalles del documento</span>
          <span className="text-gray-400 group-open:rotate-180 transition-transform">{'\u25bc'}</span>
        </summary>

        <div className="mt-3 space-y-3">
          <div className="grid grid-cols-2 gap-3">
            <Campo etiqueta="Número de Documento" valor={data.documentNumber} mono />
            <Campo etiqueta="Fecha" valor={data.documentDate} />
            <Campo etiqueta="Cliente" valor={data.customerName} />
            <Campo etiqueta="Identificación" valor={data.customerIdentification} mono />
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-3 p-3 bg-white dark:bg-gray-800/60 rounded">
            <Campo etiqueta="Total" valor={moneda(data.total)} />
            <Campo etiqueta="IVA" valor={moneda(data.tax)} />
            <Campo etiqueta="Descuento" valor={moneda(data.discount)} />
            <Campo etiqueta="Retención" valor={moneda(data.withholding)} />
          </div>

          {(data.electronic || data.stampStatus || data.cufe) && (
            <div className="p-3 bg-white dark:bg-gray-800/60 rounded space-y-2">
              <p className="text-xs text-gray-500 dark:text-gray-400 font-medium">Facturación electrónica</p>
              <div className="grid grid-cols-2 gap-3">
                <Campo etiqueta="Estado DIAN" valor={data.stampStatus} />
                <Campo etiqueta="CUFE" valor={data.cufe} mono />
              </div>
              {timbradoPendiente && (
                <p className="text-xs text-amber-700 dark:text-amber-400">
                  El proveedor aun esta enviando el documento a la DIAN: el CUFE aparece cuando termina el timbrado.
                </p>
              )}
              {data.publicUrl && (
                <a
                  href={data.publicUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex text-xs font-medium text-blue-600 hover:text-blue-700 underline"
                >
                  Ver factura en el proveedor
                </a>
              )}
            </div>
          )}

          {data.items.length > 0 && (
            <div>
              <p className="text-xs text-gray-500 dark:text-gray-400 mb-2">Items ({data.items.length})</p>
              <div className="space-y-2">
                {data.items.map((item, idx) => (
                  <div
                    key={`${item.code ?? 'item'}-${idx}`}
                    className="p-2 bg-white dark:bg-gray-800/80 rounded text-xs border border-gray-200 dark:border-gray-700"
                  >
                    <div className="flex justify-between items-start mb-1 gap-2">
                      <span className="font-medium text-gray-900 dark:text-white">{item.description || item.code}</span>
                      <span className="font-semibold text-gray-900 dark:text-white whitespace-nowrap">{moneda(item.total)}</span>
                    </div>
                    <div className="flex flex-wrap gap-3 text-gray-600 dark:text-gray-300">
                      {item.code && item.description && <span className="font-mono">{item.code}</span>}
                      {cantidad(item.quantity) && <span>Cant: {cantidad(item.quantity)}</span>}
                      {moneda(item.unitPrice) && <span>Unit: {moneda(item.unitPrice)}</span>}
                      {cantidad(item.taxRate) && <span>IVA: {cantidad(item.taxRate)}%</span>}
                      {Number(item.discount) > 0 && <span>Desc: {moneda(item.discount)}</span>}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {data.payments.length > 0 && (
            <div className="p-3 bg-white dark:bg-gray-800/60 rounded">
              <p className="text-xs text-gray-500 dark:text-gray-400 mb-2 font-medium">Medios de pago</p>
              <div className="space-y-1 text-xs text-gray-700 dark:text-gray-200">
                {data.payments.map((pago, idx) => (
                  <div key={`${pago.name ?? 'pago'}-${idx}`} className="flex justify-between gap-3">
                    <span>{pago.name}</span>
                    <span className="font-medium">{moneda(pago.value)}</span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {data.shipping && (
            <div className="p-3 bg-white dark:bg-gray-800/60 rounded">
              <p className="text-xs text-gray-500 dark:text-gray-400 mb-2 font-medium">Información de Envío</p>
              <div className="space-y-1 text-xs text-gray-700 dark:text-gray-200">
                {data.shipping.address && <p>{data.shipping.address}</p>}
                {data.shipping.city && (
                  <p>
                    {data.shipping.city}
                    {data.shipping.department ? `, ${data.shipping.department}` : ''}
                  </p>
                )}
                {data.shipping.phone && <p>{data.shipping.phone}</p>}
              </div>
            </div>
          )}

          {data.cashReceipt && (
            <div
              className={`p-3 rounded ${
                data.cashReceipt.status === 'success'
                  ? 'bg-green-50 dark:bg-green-900/10 border border-green-200 dark:border-green-800'
                  : 'bg-red-50 dark:bg-red-900/10 border border-red-200 dark:border-red-800'
              }`}
            >
              <p
                className={`text-xs uppercase tracking-wide font-medium mb-2 ${
                  data.cashReceipt.status === 'success' ? 'text-green-600' : 'text-red-600'
                }`}
              >
                Recibo de Caja {data.cashReceipt.status === 'success' ? '- Registrado' : '- Error'}
              </p>
              {data.cashReceipt.status === 'success' ? (
                <div className="grid grid-cols-2 gap-2 text-xs">
                  <Campo etiqueta="Medio de pago" valor={data.cashReceipt.paymentType} />
                  <Campo etiqueta="Monto" valor={moneda(data.cashReceipt.amount)} />
                  {data.cashReceipt.message && (
                    <div className="col-span-2">
                      <Campo etiqueta="Respuesta" valor={data.cashReceipt.message} />
                    </div>
                  )}
                </div>
              ) : (
                <p className="text-xs text-red-700 dark:text-red-400">{data.cashReceipt.error}</p>
              )}
            </div>
          )}

          {data.notes && <Campo etiqueta="Observaciones" valor={data.notes} />}

          {data.completedFields && data.completedFields.length > 0 && (
            <p className="text-xs text-gray-500 dark:text-gray-400">
              {'El proveedor no devolvió '}
              {data.completedFields.join(', ')}
              {': esos campos se muestran con los datos de la orden en Probability.'}
            </p>
          )}

          {raw && (
            <details className="mt-3">
              <summary className="text-xs text-gray-500 dark:text-gray-400 cursor-pointer hover:text-gray-700 flex items-center gap-1">
                <span>Ver JSON completo</span>
                {copySlot}
              </summary>
              <pre className="mt-2 text-xs bg-white dark:bg-gray-800/80 rounded p-3 overflow-x-auto max-h-64 border border-gray-200 dark:border-gray-700 font-mono text-gray-700 dark:text-gray-200">
                {JSON.stringify(raw, null, 2)}
              </pre>
            </details>
          )}
        </div>
      </details>
    </div>
  );
}
