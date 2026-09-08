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
export type TicketArea = 'comercial' | 'soporte' | 'desarrollo';

export interface Ticket {
    id: number;
    code: string;
    business_id?: number | null;
    business_name?: string;
    created_by_id: number;
    created_by_name?: string;
    created_by_avatar_url?: string;
    assigned_to_id?: number | null;
    assigned_to_name?: string;
    assigned_to_avatar_url?: string;
    sprint_id?: number | null;
    sprint_name?: string;
    title: string;
    description: string;
    type: TicketType;
    category?: string;
    priority: TicketPriority;
    status: TicketStatus;
    source: TicketSource;
    severity?: TicketSeverity;
    area?: TicketArea;
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
    area?: TicketArea;
    assigned_to_id?: number | null;
    sprint_id?: number | null;
    due_date?: string | null;
}

export interface UpdateTicketDTO {
    title?: string;
    description?: string;
    type?: TicketType;
    category?: string;
    priority?: TicketPriority;
    severity?: TicketSeverity;
    area?: TicketArea;
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
    area?: string;
    source?: string;
    escalated?: boolean;
    sprint_id?: number | string;
    search?: string;
    only_mine?: boolean;
    assigned_to_id?: number;
    created_by_id?: number;
    sort_by?: string;
    sort_order?: 'asc' | 'desc';
}

export interface PaginatedTickets {
    data: Ticket[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export const STATUS_META: Record<TicketStatus, { label: string; color: string; bg: string; ring: string }> = {
    open:           { label: 'Abierto',         color: 'text-blue-700 dark:text-blue-300',     bg: 'bg-blue-100 dark:bg-blue-500/15',     ring: 'ring-blue-300 dark:ring-blue-500/30' },
    in_review:      { label: 'En revisi\u00f3n',     color: 'text-amber-700 dark:text-amber-300',   bg: 'bg-amber-100 dark:bg-amber-500/15',   ring: 'ring-amber-300 dark:ring-amber-500/30' },
    in_development: { label: 'En desarrollo',   color: 'text-purple-700 dark:text-purple-300', bg: 'bg-purple-100 dark:bg-purple-500/20', ring: 'ring-purple-300 dark:ring-purple-500/30' },
    testing:        { label: 'Pruebas',         color: 'text-cyan-700 dark:text-cyan-300',     bg: 'bg-cyan-100 dark:bg-cyan-500/15',     ring: 'ring-cyan-300 dark:ring-cyan-500/30' },
    blocked:        { label: 'Bloqueado',       color: 'text-red-700 dark:text-red-300',       bg: 'bg-red-100 dark:bg-red-500/15',       ring: 'ring-red-300 dark:ring-red-500/30' },
    resolved:       { label: 'Resuelto',        color: 'text-emerald-700 dark:text-emerald-300', bg: 'bg-emerald-100 dark:bg-emerald-500/15', ring: 'ring-emerald-300 dark:ring-emerald-500/30' },
    closed:         { label: 'Cerrado',         color: 'text-gray-700 dark:text-gray-300',     bg: 'bg-gray-200 dark:bg-white/10',        ring: 'ring-gray-300 dark:ring-white/15' },
    wont_fix:       { label: 'No se har\u00e1',      color: 'text-zinc-700 dark:text-zinc-300',     bg: 'bg-zinc-200 dark:bg-white/10',        ring: 'ring-zinc-300 dark:ring-white/15' },
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
    integration: { label: 'Integraci\u00f3n',        icon: 'INT' },
    support:     { label: 'Soporte',            icon: 'SUP' },
    complaint:   { label: 'Queja',              icon: 'QJA' },
    claim:       { label: 'Reclamo',            icon: 'RCL' },
    question:    { label: 'Pregunta',           icon: 'PRG' },
};

export const TICKET_STATUSES: TicketStatus[] = ['open', 'in_review', 'in_development', 'testing', 'blocked', 'resolved', 'closed', 'wont_fix'];
export const TICKET_PRIORITIES: TicketPriority[] = ['low', 'medium', 'high', 'critical'];
export const TICKET_TYPES: TicketType[] = ['bug', 'improvement', 'feature', 'data', 'integration', 'support', 'complaint', 'claim', 'question'];
export const TICKET_AREAS: TicketArea[] = ['comercial', 'soporte', 'desarrollo'];
export const TICKET_SEVERITIES: TicketSeverity[] = ['', 'low', 'medium', 'high'];
export const TICKET_MODULES: string[] = [
    '\u00d3rdenes',
    'Env\u00edos',
    'Facturaci\u00f3n',
    'Productos',
    'Inventario',
    'Clientes',
    'Pagos y billetera',
    'Integraciones',
    'WhatsApp',
    'Usuarios y permisos',
    'Dashboard',
    'P\u00e1gina web',
    'M\u00f3vil',
    'Infraestructura',
];

export const TICKET_SOURCES: TicketSource[] = ['internal', 'business'];

export const SEVERITY_META: Record<TicketSeverity, { label: string }> = {
    '':     { label: 'Sin severidad' },
    low:    { label: 'Baja' },
    medium: { label: 'Media' },
    high:   { label: 'Alta' },
};

export const SOURCE_META: Record<TicketSource, { label: string }> = {
    internal: { label: 'Interno' },
    business: { label: 'Negocio' },
};

export const AREA_META: Record<TicketArea, { label: string; color: string; bg: string }> = {
    comercial:  { label: 'Comercial',  color: 'text-pink-700 dark:text-pink-200',     bg: 'bg-pink-100 dark:bg-pink-900/40' },
    soporte:    { label: 'Soporte',    color: 'text-cyan-700 dark:text-cyan-200',     bg: 'bg-cyan-100 dark:bg-cyan-900/40' },
    desarrollo: { label: 'Desarrollo', color: 'text-purple-700 dark:text-purple-200', bg: 'bg-purple-100 dark:bg-purple-900/40' },
};

export const PRIORITY_ACCENT: Record<TicketPriority, string> = {
    low:      'border-l-gray-400 dark:border-l-gray-600',
    medium:   'border-l-blue-400 dark:border-l-blue-500',
    high:     'border-l-amber-400 dark:border-l-amber-500',
    critical: 'border-l-red-400 dark:border-l-red-500',
};

export const PRIORITY_DOT: Record<TicketPriority, string> = {
    low:      'text-gray-500 dark:text-gray-400',
    medium:   'text-blue-500 dark:text-blue-400',
    high:     'text-amber-500 dark:text-amber-400',
    critical: 'text-red-500 dark:text-red-400',
};

export const TYPE_CHIP: Record<TicketType, string> = {
    bug:         'text-red-700 bg-red-100 dark:text-red-400 dark:bg-red-500/15',
    improvement: 'text-indigo-700 bg-indigo-100 dark:text-indigo-300 dark:bg-indigo-500/20',
    feature:     'text-emerald-700 bg-emerald-100 dark:text-emerald-400 dark:bg-emerald-500/15',
    data:        'text-amber-700 bg-amber-100 dark:text-amber-300 dark:bg-amber-500/15',
    integration: 'text-cyan-700 bg-cyan-100 dark:text-cyan-300 dark:bg-cyan-500/15',
    support:     'text-sky-700 bg-sky-100 dark:text-sky-300 dark:bg-sky-500/15',
    complaint:   'text-orange-700 bg-orange-100 dark:text-orange-300 dark:bg-orange-500/15',
    claim:       'text-rose-700 bg-rose-100 dark:text-rose-300 dark:bg-rose-500/15',
    question:    'text-violet-700 bg-violet-100 dark:text-violet-300 dark:bg-violet-500/15',
};
