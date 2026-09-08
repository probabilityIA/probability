export type CampaignStatus =
  | "draft"
  | "scheduled"
  | "running"
  | "paused"
  | "completed"
  | "cancelled";

export type CampaignAudienceType = "all_clients" | "filtered_clients";

export type CampaignSendStatus =
  | "pending"
  | "queued"
  | "sent"
  | "delivered"
  | "read"
  | "replied"
  | "failed"
  | "skipped";

export interface Campaign {
  ID: number;
  BusinessID: number;
  IntegrationID: number | null;
  WhatsappTemplateID: number | null;
  Name: string;
  Description: string;
  SenderName: string;
  AudienceType: CampaignAudienceType;
  AudienceParams: {
    City: string;
    CreatedFromDays: number;
    OnlyWithoutOrder: boolean;
    ClientIDs: number[] | null;
  };
  VariableValues: Record<string, string> | null;
  Timezone: string;
  SendWindowStart: string;
  SendWindowEnd: string;
  ScheduledAt: string | null;
  DailySendCap: number;
  BatchSize: number;
  Status: CampaignStatus;
  AudienceCount: number;
  QueuedCount: number;
  SentCount: number;
  FailedCount: number;
  SkippedCount: number;
  RepliedCount: number;
  StartedAt: string | null;
  FinishedAt: string | null;
  LastBatchAt: string | null;
  ErrorMessage: string;
  CreatedAt: string;
}

export interface CampaignSend {
  ID: number;
  CampaignID: number;
  ClientID: number;
  Phone: string;
  ClientName: string;
  Status: CampaignSendStatus;
  MessageID: string;
  ErrorMessage: string;
  QueuedAt: string | null;
  SentAt: string | null;
  DeliveredAt: string | null;
  ReadAt: string | null;
  RepliedAt: string | null;
}

export interface CampaignAudiencePreview {
  total: number;
  opted_out: number;
  no_phone: number;
  reachable: number;
  sample_names: string[] | null;
}

export interface CreateCampaignDTO {
  whatsapp_template_id: number;
  name: string;
  description?: string;
  sender_name?: string;
  audience_type: CampaignAudienceType;
  city?: string;
  created_from_days?: number;
  only_without_order?: boolean;
  client_ids?: number[];
  variable_values?: Record<string, string>;
  timezone?: string;
  send_window_start?: string;
  send_window_end?: string;
  scheduled_at?: string | null;
  daily_send_cap?: number;
  batch_size?: number;
}
