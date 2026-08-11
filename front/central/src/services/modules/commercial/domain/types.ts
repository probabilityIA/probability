export interface Prospect {
  id: number;
  domain: string;
  profile_url: string;
  meli_seller_id?: string;
  business_name?: string;
  platform: string;
  city?: string;
  phone?: string;
  whatsapp_number?: string;
  whatsapp_url?: string;
  instagram_handle?: string;
  instagram_url?: string;
  instagram_followers: number;
  product_count: number;
  has_cod: boolean;
  num_couriers_mentioned: number;
  uses_probability_channel: boolean;
  discovery_query?: string;
  score: number;
  score_breakdown?: Record<string, number>;
  status: string;
  seen: boolean;
  source_discovered_at?: string;
  source_last_scanned_at?: string;
}

export interface GetProspectsParams {
  page?: number;
  page_size?: number;
  platform?: string;
  status?: string;
  city?: string;
  search?: string;
  min_score?: number;
  max_score?: number;
  has_cod?: boolean;
}

export interface UpdateSeenResponse {
  success: boolean;
  error?: string;
}

export interface PaginatedProspectsResponse {
  success: boolean;
  error?: string;
  data: Prospect[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface CountByKey {
  key: string;
  count: number;
}

export interface ScoreBucket {
  range_label: string;
  count: number;
}

export interface ProspectStats {
  total: number;
  by_platform: CountByKey[];
  by_status: CountByKey[];
  by_city: CountByKey[];
  score_buckets: ScoreBucket[];
  has_cod_count: number;
  multi_channel_count: number;
  average_score: number;
  top_score_domain: string;
  top_score: number;
}

export interface ProspectStatsResponse {
  success: boolean;
  error?: string;
  data: ProspectStats;
}
