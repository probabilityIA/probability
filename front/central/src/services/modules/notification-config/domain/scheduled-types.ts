export type TemplateStatus =
  | "draft"
  | "pending"
  | "approved"
  | "rejected"
  | "paused"
  | "disabled"
  | "failed";

export type TemplateCategory = "MARKETING" | "UTILITY";

export type TemplateScope = "order_event" | "scheduled";

export interface TemplateVariable {
  position: number;
  source: string;
  label: string;
  fallback: string;
}

export interface TemplateButton {
  type: string;
  text: string;
  url: string;
}

export interface WhatsappTemplate {
  ID: number;
  BusinessID: number;
  Name: string;
  Language: string;
  Category: TemplateCategory;
  BodyText: string;
  HeaderText: string;
  FooterText: string;
  Origin: string;
  Scope: string;
  Variables: TemplateVariable[] | null;
  Buttons: TemplateButton[] | null;
  MetaTemplateID: string;
  Status: TemplateStatus;
  RejectedReason: string;
  SubmittedAt: string | null;
  ReviewedAt: string | null;
  CreatedAt: string;
}

export interface CreateTemplateDTO {
  scope: TemplateScope;
  name: string;
  language: string;
  category: TemplateCategory;
  header_text?: string;
  body_text: string;
  footer_text?: string;
  variables: Array<{
    position: number;
    source: string;
    label?: string;
    fallback?: string;
  }>;
  buttons?: Array<{ type: string; text: string; url?: string }>;
}

export interface ScheduledRule {
  ID: number;
  BusinessID: number;
  NotificationTypeID: number;
  WhatsappTemplateID: number | null;
  Name: string;
  Description: string;
  SegmentType: string;
  SegmentParams: {
    DaysWithoutPurchase: number;
    MinOrders: number;
    MaxOrders: number;
  };
  Timezone: string;
  FrequencyMinutes: number;
  SendWindowStart: string;
  SendWindowEnd: string;
  CooldownDays: number;
  DailySendCap: number;
  BatchSizeCap: number;
  RequiresOptIn: boolean;
  Enabled: boolean;
  LastRunAt: string | null;
  NextRunAt: string | null;
}

export interface CreateScheduledRuleDTO {
  notification_type_id: number;
  whatsapp_template_id: number;
  name: string;
  description?: string;
  segment_type: string;
  days_without_purchase: number;
  min_orders?: number;
  max_orders?: number;
  timezone?: string;
  frequency_minutes?: number;
  send_window_start?: string;
  send_window_end?: string;
  cooldown_days?: number;
  daily_send_cap?: number;
  batch_size_cap?: number;
  requires_opt_in?: boolean;
  enabled?: boolean;
}

export interface ScheduledRun {
  ID: number;
  RuleID: number;
  StartedAt: string;
  FinishedAt: string | null;
  Status: string;
  MatchedCount: number;
  QueuedCount: number;
  SkippedCount: number;
  FailedCount: number;
  CapReachedFlag: boolean;
  ErrorMessage: string;
}

export interface PaginatedResult<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export const TEMPLATE_STATUS_LABEL: Record<TemplateStatus, string> = {
  draft: "Borrador",
  pending: "En revisión de Meta",
  approved: "Aprobada",
  rejected: "Rechazada",
  paused: "Pausada por Meta",
  disabled: "Deshabilitada",
  failed: "Falló el envío a Meta",
};

export const SEGMENT_LABEL: Record<string, string> = {
  customers_inactive: "Clientes sin comprar hace N días",
};

export interface TemplatePreview {
  name: string;
  language: string;
  status: string;
  category: string;
  meta_id: string;
  header: string;
  body: string;
  footer: string;
  buttons: string[] | null;
  variables: number;
  condition: string;
  found_in_whatsapp: boolean;
  rejected_reason: string;
}
