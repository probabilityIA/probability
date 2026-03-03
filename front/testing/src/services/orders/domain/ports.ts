import { ReferenceData, GenerateOrdersDTO, GenerateResult } from "./types";

export interface IOrdersRepository {
  getReferenceData(businessId: number): Promise<ReferenceData>;
  generateOrders(businessId: number, dto: GenerateOrdersDTO): Promise<GenerateResult>;
}
