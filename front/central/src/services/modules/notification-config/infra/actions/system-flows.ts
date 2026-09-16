"use server";

import { cookies } from "next/headers";
import { env } from "@/shared/config/env";
import { SystemFlow } from "../../domain/scheduled-types";

export async function listSystemFlowsAction(businessId?: number) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("session_token")?.value || "";
    const query = businessId ? `?business_id=${businessId}` : "";

    const response = await fetch(`${env.API_BASE_URL}/integrations/whatsapp/system-flows${query}`, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      cache: "no-store",
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "No se pudieron cargar los flujos del sistema", data: [] as SystemFlow[] };
    }
    return { success: true, data: (body.data || []) as SystemFlow[] };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message, data: [] as SystemFlow[] };
  }
}
