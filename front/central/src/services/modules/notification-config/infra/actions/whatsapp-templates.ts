"use server";

import { cookies } from "next/headers";
import { revalidatePath } from "next/cache";
import { env } from "@/shared/config/env";
import {
  CreateTemplateDTO,
  PaginatedResult,
  TemplateFlow,
  TemplateFlowInput,
  TemplateScope,
  UpdateTemplateDTO,
  WhatsappTemplate,
} from "../../domain/scheduled-types";

async function authHeaders() {
  const cookieStore = await cookies();
  const token = cookieStore.get("session_token")?.value || "";
  return {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
}

function withBusiness(path: string, businessId?: number) {
  if (!businessId) return path;
  const separator = path.includes("?") ? "&" : "?";
  return `${path}${separator}business_id=${businessId}`;
}

export async function listTemplatesAction(
  businessId?: number,
  scope: TemplateScope | "all" = "scheduled",
  status?: string,
  page = 1,
  pageSize = 20,
) {
  try {
    const query = new URLSearchParams({
      page: String(page),
      page_size: String(pageSize),
      scope,
    });
    if (status) query.set("status", status);

    const url = withBusiness(
      `${env.API_BASE_URL}/whatsapp-templates?${query.toString()}`,
      businessId,
    );

    const response = await fetch(url, {
      headers: await authHeaders(),
      cache: "no-store",
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error listando plantillas", data: [] };
    }

    const result = body as PaginatedResult<WhatsappTemplate>;
    return { success: true, data: result.data || [], total: result.total || 0 };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message, data: [] };
  }
}

export async function getTemplateVariablesAction(businessId?: number) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-templates/variables`, businessId);
    const response = await fetch(url, {
      headers: await authHeaders(),
      cache: "no-store",
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error", data: {} };
    }

    return { success: true, data: (body.data || {}) as Record<string, string> };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message, data: {} };
  }
}

export async function createTemplateAction(dto: CreateTemplateDTO, businessId?: number) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-templates`, businessId);
    const response = await fetch(url, {
      method: "POST",
      headers: await authHeaders(),
      body: JSON.stringify(dto),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error creando la plantilla" };
    }

    revalidatePath("/notification-config");
    return { success: true, data: body.data as WhatsappTemplate };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function updateTemplateAction(
  id: number,
  dto: UpdateTemplateDTO,
  businessId?: number,
) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-templates/${id}`, businessId);
    const response = await fetch(url, {
      method: "PUT",
      headers: await authHeaders(),
      body: JSON.stringify(dto),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error editando la plantilla" };
    }

    revalidatePath("/notification-config");
    return { success: true, data: body.data as WhatsappTemplate };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function uploadTemplateMediaAction(formData: FormData, businessId?: number) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("session_token")?.value || "";

    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-templates/media`, businessId);
    const response = await fetch(url, {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: formData,
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error subiendo la imagen" };
    }

    return { success: true, url: (body.data?.url || "") as string };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function listTemplateFlowsAction(templateId: number, businessId?: number) {
  try {
    const url = withBusiness(
      `${env.API_BASE_URL}/whatsapp-templates/${templateId}/flows`,
      businessId,
    );
    const response = await fetch(url, {
      headers: await authHeaders(),
      cache: "no-store",
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error listando el flujo", data: [] };
    }

    return { success: true, data: (body.data || []) as TemplateFlow[] };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message, data: [] };
  }
}

export async function listAllTemplateFlowsAction(businessId?: number) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-templates/flows`, businessId);
    const response = await fetch(url, {
      headers: await authHeaders(),
      cache: "no-store",
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error listando los flujos", data: [] };
    }

    return { success: true, data: (body.data || []) as TemplateFlow[] };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message, data: [] };
  }
}

export async function replaceTemplateFlowsAction(
  templateId: number,
  flows: TemplateFlowInput[],
  businessId?: number,
) {
  try {
    const url = withBusiness(
      `${env.API_BASE_URL}/whatsapp-templates/${templateId}/flows`,
      businessId,
    );
    const response = await fetch(url, {
      method: "PUT",
      headers: await authHeaders(),
      body: JSON.stringify({ flows }),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error guardando el flujo" };
    }

    revalidatePath("/notification-config");
    return { success: true, data: (body.data || []) as TemplateFlow[] };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function deleteTemplateAction(id: number, businessId?: number) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-templates/${id}`, businessId);
    const response = await fetch(url, {
      method: "DELETE",
      headers: await authHeaders(),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error eliminando la plantilla" };
    }

    revalidatePath("/notification-config");
    return { success: true };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function submitTemplateForReviewAction(id: number, businessId?: number) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-templates/${id}/submit`, businessId);
    const response = await fetch(url, {
      method: "POST",
      headers: await authHeaders(),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error enviando la plantilla a revisión" };
    }

    revalidatePath("/notification-config");
    return { success: true, data: body.data as WhatsappTemplate };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}
