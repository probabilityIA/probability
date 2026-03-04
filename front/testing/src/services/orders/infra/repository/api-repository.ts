import { testingAPI } from "@/shared/lib/api";
import type { IOrdersRepository } from "../../domain/ports";
import type { ReferenceData, GenerateOrdersDTO, GenerateResult, DeleteResult } from "../../domain/types";

export class OrdersApiRepository implements IOrdersRepository {
  async getReferenceData(businessId: number): Promise<ReferenceData> {
    const res = await testingAPI<{ data: ReferenceData }>(
      `/orders/reference-data?business_id=${businessId}`
    );
    return res.data;
  }

  async generateOrders(businessId: number, dto: GenerateOrdersDTO): Promise<GenerateResult> {
    const res = await testingAPI<{ data: GenerateResult }>(
      `/orders/generate?business_id=${businessId}`,
      {
        method: "POST",
        body: JSON.stringify(dto),
      }
    );
    return res.data;
  }

  async deleteAllOrders(businessId: number): Promise<DeleteResult> {
    return testingAPI<DeleteResult>(
      `/orders?business_id=${businessId}`,
      { method: "DELETE" }
    );
  }
}
