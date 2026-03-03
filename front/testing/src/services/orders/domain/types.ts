export interface Product {
  id: string;
  name: string;
  sku: string;
  price: number;
  currency: string;
}

export interface Integration {
  id: number;
  name: string;
  code: string;
  category: string;
  category_id: number;
  integration_type_id: number;
  integration_type_code: string;
}

export interface PaymentMethod {
  id: number;
  code: string;
  name: string;
}

export interface OrderStatus {
  id: number;
  code: string;
  name: string;
}

export interface ReferenceData {
  products: Product[];
  integrations: Integration[];
  payment_methods: PaymentMethod[];
  order_statuses: OrderStatus[];
  webhook_topics: Record<string, string[]>;
}

export interface WebhookPayload {
  url: string;
  method: string;
  headers: Record<string, string>;
  body: Record<string, unknown>;
  raw_body?: string; // exact bytes for HMAC webhooks — send this instead of re-serializing body
  hmac_secret?: string; // debug: secret used to compute HMAC
}

export interface GenerateOrdersDTO {
  count: number;
  integration_id?: number;
  random_products: boolean;
  max_items_per_order: number;
  topic: string;
}

export interface OrderError {
  index: number;
  message: string;
}

export interface GenerateResult {
  total: number;
  payloads: WebhookPayload[] | null;
  errors: OrderError[] | null;
}

export interface APICallLog {
  index: number;
  success: boolean;
  timestamp: string;
  duration_ms: number;
  request: {
    method: string;
    url: string;
    headers?: Record<string, string>;
    body: Record<string, unknown>;
    hmac_secret?: string;
  };
  response: {
    status_code: number;
    body: string;
  };
}
