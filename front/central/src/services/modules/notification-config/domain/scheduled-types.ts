export type TemplateStatus =
  | "draft"
  | "pending"
  | "approved"
  | "rejected"
  | "paused"
  | "disabled"
  | "failed";

export type TemplateCategory = "MARKETING" | "UTILITY";

export type TemplateHeaderType = "TEXT" | "IMAGE";

export type TemplateScope = "order_event" | "scheduled" | "campaign";

export interface TemplateVariable {
  Position: number;
  Source: string;
  Label: string;
  Fallback: string;
}

export interface TemplateButton {
  Type: string;
  Text: string;
  URL: string;
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
  HeaderType: TemplateHeaderType;
  HeaderMediaURL: string;
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
  header_type?: TemplateHeaderType;
  header_media_url?: string;
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

export interface UpdateTemplateDTO {
  category: TemplateCategory;
  header_text?: string;
  header_type?: TemplateHeaderType;
  header_media_url?: string;
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

export interface Flow {
  ID: number;
  BusinessID: number;
  Name: string;
  Description: string;
  RootTemplateID: number | null;
  Enabled: boolean;
  RootTemplateName: string;
  RootTemplateStatus: TemplateStatus | "";
  StepCount: number;
  PendingCount: number;
}

export interface FlowInput {
  name: string;
  description?: string;
  root_template_id?: number | null;
  enabled?: boolean;
}

export interface TemplateFlow {
  ID: number;
  BusinessID: number;
  FlowID: number | null;
  SourceTemplateID: number;
  ButtonText: string;
  TargetTemplateID: number;
  Enabled: boolean;
  SourceName: string;
  TargetName: string;
  TargetLanguage: string;
  TargetStatus: TemplateStatus;
  TargetHeaderMediaURL: string;
}

export interface TemplateFlowInput {
  button_text: string;
  target_template_id: number;
  enabled?: boolean;
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
