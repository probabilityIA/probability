import {
    Ticket,
    TicketComment,
    TicketAttachment,
    TicketHistoryEntry,
    PaginatedTickets,
    ListTicketsParams,
    CreateTicketDTO,
    UpdateTicketDTO,
} from './types';

export interface ITicketRepository {
    list(params?: ListTicketsParams): Promise<PaginatedTickets>;
    get(id: number, businessId?: number): Promise<Ticket>;
    create(data: CreateTicketDTO): Promise<Ticket>;
    update(id: number, data: UpdateTicketDTO): Promise<Ticket>;
    remove(id: number): Promise<void>;
    changeStatus(id: number, status: string, note?: string): Promise<Ticket>;
    assign(id: number, assignedToId: number | null): Promise<Ticket>;
    escalate(id: number, note?: string): Promise<Ticket>;

    listComments(id: number, businessId?: number): Promise<TicketComment[]>;
    addComment(id: number, body: string, isInternal: boolean): Promise<TicketComment>;

    listAttachments(id: number, businessId?: number): Promise<TicketAttachment[]>;
    uploadAttachment(id: number, file: File, commentId?: number): Promise<TicketAttachment>;
    deleteAttachment(attachmentId: number): Promise<void>;

    listHistory(id: number, businessId?: number): Promise<TicketHistoryEntry[]>;
}
