export type TicketStatus =
    | 'open'
    | 'in_review'
    | 'in_development'
    | 'testing'
    | 'blocked'
    | 'resolved'
    | 'closed'
    | 'wont_fix';

export type TicketPriority = 'low' | 'medium' | 'high' | 'critical';

export type TicketType =
    | 'bug'
    | 'improvement'
    | 'feature'
    | 'data'
    | 'integration'
    | 'support'
    | 'complaint'
    | 'claim'
    | 'question';

export type TicketSeverity = 'low' | 'medium' | 'high' | '';
export type TicketSource = 'internal' | 'business';

export interface Ticket {
    id: number;
    code: string;
    business_id?: number | null;
    business_name?: string;
    created_by_id: number;
    created_by_name?: string;
    assigned_to_id?: number | null;
    assigned_to_name?: string;
    title: string;
    description: string;
    type: TicketType;
    category?: string;
    priority: TicketPriority;
    status: TicketStatus;
    source: TicketSource;
    severity?: TicketSeverity;
    escalated_to_dev: boolean;
    escalated_at?: string | null;
    due_date?: string | null;
    resolved_at?: string | null;
    closed_at?: string | null;
    created_at: string;
    updated_at: string;
    comments_count: number;
    attachments_count: number;
}

export interface TicketComment {
    id: number;
    ticket_id: number;
    user_id: number;
    user_name: string;
    body: string;
    is_internal: boolean;
    created_at: string;
    attachments?: TicketAttachment[];
}

export interface TicketAttachment {
    id: number;
    ticket_id: number;
    comment_id?: number | null;
    uploaded_by_id: number;
    uploaded_by_name: string;
    file_url: string;
    file_name: string;
    mime_type: string;
    size: number;
    created_at: string;
}

export interface TicketHistoryEntry {
    id: number;
    ticket_id: number;
    from_status: string;
    to_status: string;
    changed_by_id: number;
    changed_by_name: string;
    note: string;
    created_at: string;
}

export interface CreateTicketDTO {
    business_id?: number | null;
    title: string;
    description: string;
    type?: TicketType;
    category?: string;
    priority?: TicketPriority;
    severity?: TicketSeverity;
    source?: TicketSource;
    assigned_to_id?: number | null;
    due_date?: string | null;
}

export interface UpdateTicketDTO {
    title?: string;
    description?: string;
    type?: TicketType;
    category?: string;
    priority?: TicketPriority;
    severity?: TicketSeverity;
    assigned_to_id?: number | null;
    due_date?: string | null;
    clear_due_date?: boolean;
}

export interface ListTicketsParams {
    page?: number;
    page_size?: number;
    business_id?: number;
    status?: string;
    priority?: string;
    type?: string;
    source?: string;
    escalated?: boolean;
    search?: string;
    only_mine?: boolean;
    assigned_to_id?: number;
    created_by_id?: number;
}

export interface PaginatedTickets {
    data: Ticket[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export const STATUS_META: Record<TicketStatus, { label: string; color: string; bg: string; ring: string }> = {
    open:           { label: 'Abierto',         color: 'text-blue-700 dark:text-blue-200',     bg: 'bg-blue-100 dark:bg-blue-900/40',     ring: 'ring-blue-300' },
    in_review:      { label: 'En revision',     color: 'text-amber-700 dark:text-amber-200',   bg: 'bg-amber-100 dark:bg-amber-900/40',   ring: 'ring-amber-300' },
    in_development: { label: 'En desarrollo',   color: 'text-purple-700 dark:text-purple-200', bg: 'bg-purple-100 dark:bg-purple-900/40', ring: 'ring-purple-300' },
    testing:        { label: 'Pruebas',         color: 'text-cyan-700 dark:text-cyan-200',     bg: 'bg-cyan-100 dark:bg-cyan-900/40',     ring: 'ring-cyan-300' },
    blocked:        { label: 'Bloqueado',       color: 'text-red-700 dark:text-red-200',       bg: 'bg-red-100 dark:bg-red-900/40',       ring: 'ring-red-300' },
    resolved:       { label: 'Resuelto',        color: 'text-emerald-700 dark:text-emerald-200', bg: 'bg-emerald-100 dark:bg-emerald-900/40', ring: 'ring-emerald-300' },
    closed:         { label: 'Cerrado',         color: 'text-gray-700 dark:text-gray-200',     bg: 'bg-gray-200 dark:bg-gray-700',        ring: 'ring-gray-300' },
    wont_fix:       { label: 'No se hara',      color: 'text-zinc-700 dark:text-zinc-200',     bg: 'bg-zinc-200 dark:bg-zinc-700',        ring: 'ring-zinc-300' },
};

export const PRIORITY_META: Record<TicketPriority, { label: string; color: string; bg: string }> = {
    low:      { label: 'Baja',     color: 'text-gray-700 dark:text-gray-200',     bg: 'bg-gray-100 dark:bg-gray-700' },
    medium:   { label: 'Media',    color: 'text-blue-700 dark:text-blue-200',     bg: 'bg-blue-100 dark:bg-blue-900/40' },
    high:     { label: 'Alta',     color: 'text-orange-700 dark:text-orange-200', bg: 'bg-orange-100 dark:bg-orange-900/40' },
    critical: { label: 'Critica',  color: 'text-red-700 dark:text-red-200',       bg: 'bg-red-100 dark:bg-red-900/40' },
};

export const TYPE_META: Record<TicketType, { label: string; icon: string }> = {
    bug:         { label: 'Bug',                icon: 'BUG' },
    improvement: { label: 'Mejora',             icon: 'IMP' },
    feature:     { label: 'Nueva funcionalidad', icon: 'NEW' },
    data:        { label: 'Datos',              icon: 'DAT' },
    integration: { label: 'Integracion',        icon: 'INT' },
    support:     { label: 'Soporte',            icon: 'SUP' },
    complaint:   { label: 'Queja',              icon: 'QJA' },
    claim:       { label: 'Reclamo',            icon: 'RCL' },
    question:    { label: 'Pregunta',           icon: 'PRG' },
};

export const TICKET_STATUSES: TicketStatus[] = ['open', 'in_review', 'in_development', 'testing', 'blocked', 'resolved', 'closed', 'wont_fix'];
export const TICKET_PRIORITIES: TicketPriority[] = ['low', 'medium', 'high', 'critical'];
export const TICKET_TYPES: TicketType[] = ['bug', 'improvement', 'feature', 'data', 'integration', 'support', 'complaint', 'claim', 'question'];
