import { ReferenceData, GenerateOrdersDTO, GenerateResult, DeleteResult } from "./types";

export interface IOrdersRepository {
  getReferenceData(businessId: number): Promise<ReferenceData>;
  generateOrders(businessId: number, dto: GenerateOrdersDTO): Promise<GenerateResult>;
  deleteAllOrders(businessId: number): Promise<DeleteResult>;
}
