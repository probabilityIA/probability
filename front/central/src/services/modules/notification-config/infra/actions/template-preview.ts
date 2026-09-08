"use server";

import { cookies } from "next/headers";
import { env } from "@/shared/config/env";
import { TemplatePreview } from "../../domain/scheduled-types";

export async function previewTemplatesByEventAction(
  eventCode: string,
  businessId?: number,
) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("session_token")?.value || "";

    const query = new URLSearchParams({ event_code: eventCode });
    if (businessId) query.set("business_id", String(businessId));

    const response = await fetch(
      `${env.API_BASE_URL}/integrations/whatsapp/templates/preview?${query.toString()}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        cache: "no-store",
      },
    );

    const body = await response.json();
    if (!response.ok) {
      return {
        success: false,
        error: body.error || "No se pudo consultar la plantilla en Meta",
        data: [],
      };
    }

    return { success: true, data: (body.data || []) as TemplatePreview[] };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message, data: [] };
  }
}
