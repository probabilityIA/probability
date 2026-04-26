export type TrackingStatus = 'pending' | 'picked_up' | 'in_transit' | 'out_for_delivery' | 'delivered' | 'failed';

export interface TrackingSearchResult {
  id: number;
  order_id?: string;
  tracking_number: string;
  carrier: string;
  carrier_code?: string;
  status: TrackingStatus;
  client_name?: string;
  destination_address?: string;
  estimated_delivery?: string;
  shipped_at?: string;
  delivered_at?: string;
  shipping_cost?: number;
  total_cost?: number;
  tracking_url?: string;
  guide_url?: string;
  is_test: boolean;
}

export interface TrackingHistory {
  date: string;
  status: string;
  description: string;
  location: string;
}

export interface OrderPublicTracking {
  ID: string;
  OrderNumber: string;
  BusinessName: string;
  Status: string;
  IsPaid: boolean;
  TotalAmount: number;
  CodTotal?: number | null;
  Currency: string;
  CustomerName: string;
  CustomerPhone: string;
  ShippingStreet: string;
  ShippingCity: string;
  ShippingState: string;
  ShippingPostalCode: string;
  CreatedAt: string;
  OccurredAt?: string | null;
}
