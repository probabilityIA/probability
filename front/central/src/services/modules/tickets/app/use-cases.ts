import { ITicketRepository } from '../domain/ports';
import { CreateTicketDTO, UpdateTicketDTO, ListTicketsParams } from '../domain/types';

export class TicketUseCases {
    constructor(private repo: ITicketRepository) {}

    list(params?: ListTicketsParams) { return this.repo.list(params); }
    get(id: number, businessId?: number) { return this.repo.get(id, businessId); }
    create(data: CreateTicketDTO) { return this.repo.create(data); }
    update(id: number, data: UpdateTicketDTO) { return this.repo.update(id, data); }
    remove(id: number) { return this.repo.remove(id); }
    changeStatus(id: number, status: string, note?: string) { return this.repo.changeStatus(id, status, note); }
    assign(id: number, assignedToId: number | null) { return this.repo.assign(id, assignedToId); }
    escalate(id: number, note?: string) { return this.repo.escalate(id, note); }

    listComments(id: number, businessId?: number) { return this.repo.listComments(id, businessId); }
    addComment(id: number, body: string, isInternal: boolean) { return this.repo.addComment(id, body, isInternal); }

    listAttachments(id: number, businessId?: number) { return this.repo.listAttachments(id, businessId); }
    uploadAttachment(id: number, file: File, commentId?: number) { return this.repo.uploadAttachment(id, file, commentId); }
    deleteAttachment(attachmentId: number) { return this.repo.deleteAttachment(attachmentId); }

    listHistory(id: number, businessId?: number) { return this.repo.listHistory(id, businessId); }
}
