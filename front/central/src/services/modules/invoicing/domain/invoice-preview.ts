export interface InvoicePreviewItem {
  code?: string;
  description?: string;
  quantity?: number | string;
  unitPrice?: number | string;
  discount?: number | string;
  taxRate?: number | string;
  taxValue?: number | string;
  total?: number | string;
}

export interface InvoicePreviewPayment {
  name?: string;
  value?: number | string;
}

export interface InvoicePreviewShipping {
  address?: string;
  city?: string;
  department?: string;
  phone?: string;
}

export interface InvoicePreviewCashReceipt {
  status?: string;
  paymentType?: string;
  amount?: number | string;
  message?: string;
  error?: string;
}

export interface InvoicePreview {
  provider?: string;
  documentNumber?: string;
  documentName?: string;
  documentDate?: string;
  customerName?: string;
  customerIdentification?: string;
  total?: number | string;
  tax?: number | string;
  discount?: number | string;
  withholding?: number | string;
  balance?: number | string;
  electronic?: boolean;
  stampStatus?: string;
  cufe?: string;
  publicUrl?: string;
  notes?: string;
  items: InvoicePreviewItem[];
  payments: InvoicePreviewPayment[];
  shipping?: InvoicePreviewShipping;
  cashReceipt?: InvoicePreviewCashReceipt;
  completedFields?: string[];
}

type Bruto = Record<string, any>;

const numero = (v: unknown): number | undefined => {
  if (v === null || v === undefined || v === '') return undefined;
  const n = Number(v);
  return Number.isFinite(n) ? n : undefined;
};

const texto = (v: unknown): string | undefined => {
  if (typeof v === 'string' && v.trim() !== '') return v;
  return undefined;
};

const desdeCanonico = (raw: Bruto): InvoicePreview => ({
  provider: texto(raw.provider),
  documentNumber: texto(raw.document_number),
  documentName: texto(raw.document_prefix),
  documentDate: texto(raw.document_date),
  customerName: texto(raw.customer_name),
  customerIdentification: texto(raw.customer_identification),
  total: numero(raw.total),
  tax: numero(raw.tax),
  discount: numero(raw.discount),
  withholding: numero(raw.withholding),
  balance: numero(raw.balance),
  electronic: raw.electronic === true,
  stampStatus: texto(raw.stamp_status),
  cufe: texto(raw.cufe),
  publicUrl: texto(raw.public_url),
  notes: texto(raw.notes),
  items: Array.isArray(raw.items)
    ? raw.items.map((it: Bruto) => ({
        code: texto(it.code),
        description: texto(it.description),
        quantity: numero(it.quantity),
        unitPrice: numero(it.unit_price),
        discount: numero(it.discount),
        taxRate: numero(it.tax_rate),
        taxValue: numero(it.tax_value),
        total: numero(it.total),
      }))
    : [],
  payments: Array.isArray(raw.payments)
    ? raw.payments.map((p: Bruto) => ({ name: texto(p.name), value: numero(p.value) }))
    : [],
  cashReceipt: raw.cash_receipt ? mapearRecibo(raw.cash_receipt) : undefined,
});

const desdeSoftpymes = (raw: Bruto): InvoicePreview => ({
  provider: 'softpymes',
  documentNumber: texto(raw.documentNumber),
  documentName: texto(raw.documentName),
  documentDate: texto(raw.documentDate),
  customerName: texto(raw.customerName),
  customerIdentification: texto(raw.customerIdentification),
  total: numero(raw.total),
  tax: numero(raw.totalIva),
  discount: numero(raw.totalDiscount),
  withholding: numero(raw.totalWithholdingTax),
  electronic: raw.electronicDocument === true,
  notes: texto(raw.comment),
  items: Array.isArray(raw.details)
    ? raw.details.map((d: Bruto) => ({
        code: texto(d.itemCode),
        description: texto(d.itemName),
        quantity: numero(d.quantity),
        discount: numero(d.discount),
        taxRate: numero(d.iva),
        total: numero(d.value),
      }))
    : [],
  payments: [],
  shipping: raw.shipInformation
    ? {
        address: texto(raw.shipInformation.shipAddress),
        city: texto(raw.shipInformation.shipCity),
        department: texto(raw.shipInformation.shipDepartment),
        phone: texto(raw.shipInformation.shipPhone),
      }
    : undefined,
  cashReceipt: raw.cash_receipt ? mapearRecibo(raw.cash_receipt) : undefined,
});

const mapearRecibo = (raw: Bruto): InvoicePreviewCashReceipt => ({
  status: texto(raw.status),
  paymentType: texto(raw.payment_type),
  amount: numero(raw.amount),
  message: texto(raw.message),
  error: texto(raw.error),
});

export interface InvoicePreviewFallback {
  customerName?: string;
  customerIdentification?: string;
  documentDate?: string;
  total?: number;
  tax?: number;
  discount?: number;
}

const vacio = (v: unknown) => v === undefined || v === null || v === '';

const completar = (preview: InvoicePreview, fallback?: InvoicePreviewFallback): InvoicePreview => {
  if (!fallback) return preview;

  const completados: string[] = [];
  const resultado = { ...preview };

  if (vacio(resultado.customerName) && fallback.customerName) {
    resultado.customerName = fallback.customerName;
    completados.push('cliente');
  }
  if (vacio(resultado.customerIdentification) && fallback.customerIdentification) {
    resultado.customerIdentification = fallback.customerIdentification;
    completados.push('identificacion');
  }
  if (vacio(resultado.documentDate) && fallback.documentDate) {
    resultado.documentDate = fallback.documentDate;
  }
  if (!Number(resultado.total) && fallback.total) {
    resultado.total = fallback.total;
  }
  if (!Number(resultado.tax) && fallback.tax) {
    resultado.tax = fallback.tax;
    completados.push('impuestos');
  }
  if (!Number(resultado.discount) && fallback.discount) {
    resultado.discount = fallback.discount;
    completados.push('descuento');
  }

  resultado.completedFields = completados;
  return resultado;
};

export const normalizeInvoicePreview = (
  raw?: Bruto | null,
  fallback?: InvoicePreviewFallback,
): InvoicePreview | null => {
  if (!raw || typeof raw !== 'object') return null;
  if (raw.preview_version) return completar(desdeCanonico(raw), fallback);
  if (raw.documentNumber || raw.details || raw.customerName) return completar(desdeSoftpymes(raw), fallback);
  if (raw.document_number || raw.items) return completar(desdeCanonico(raw), fallback);
  return null;
};
