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
  integration_type_id: number;
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
}

export interface GenerateOrdersDTO {
  count: number;
  integration_id?: number;
  random_products: boolean;
  max_items_per_order: number;
}

export interface CreatedOrder {
  id: string;
  order_number: string;
  total: number;
  customer_name: string;
}

export interface OrderError {
  index: number;
  message: string;
}

export interface APICallLog {
  index: number;
  success: boolean;
  timestamp: string;
  duration_ms: number;
  request: {
    method: string;
    url: string;
    body: Record<string, unknown>;
  };
  response: {
    status_code: number;
    body: string;
  };
}

export interface GenerateResult {
  total: number;
  created: number;
  failed: number;
  orders: CreatedOrder[] | null;
  errors: OrderError[] | null;
  api_logs: APICallLog[] | null;
}
