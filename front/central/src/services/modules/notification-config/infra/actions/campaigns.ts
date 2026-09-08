"use server";

import { cookies } from "next/headers";
import { revalidatePath } from "next/cache";
import { env } from "@/shared/config/env";
import {
  Campaign,
  CampaignAudiencePreview,
  CampaignSend,
  CreateCampaignDTO,
} from "../../domain/campaign-types";
import { PaginatedResult } from "../../domain/scheduled-types";

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

export async function listCampaignsAction(
  businessId?: number,
  status?: string,
  page = 1,
  pageSize = 20,
) {
  try {
    const query = new URLSearchParams({
      page: String(page),
      page_size: String(pageSize),
    });
    if (status) query.set("status", status);

    const url = withBusiness(
      `${env.API_BASE_URL}/whatsapp-campaigns?${query.toString()}`,
      businessId,
    );

    const response = await fetch(url, {
      headers: await authHeaders(),
      cache: "no-store",
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error listando campanas", data: [] };
    }

    const result = body as PaginatedResult<Campaign>;
    return { success: true, data: result.data || [], total: result.total || 0 };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message, data: [] };
  }
}

export async function getCampaignAction(id: number, businessId?: number) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-campaigns/${id}`, businessId);
    const response = await fetch(url, {
      headers: await authHeaders(),
      cache: "no-store",
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error obteniendo la campana" };
    }

    return { success: true, data: body.data as Campaign };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function createCampaignAction(dto: CreateCampaignDTO, businessId?: number) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-campaigns`, businessId);
    const response = await fetch(url, {
      method: "POST",
      headers: await authHeaders(),
      body: JSON.stringify(dto),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error creando la campana" };
    }

    revalidatePath("/notification-config");
    return { success: true, data: body.data as Campaign };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function updateCampaignAction(
  id: number,
  dto: CreateCampaignDTO,
  businessId?: number,
) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-campaigns/${id}`, businessId);
    const response = await fetch(url, {
      method: "PUT",
      headers: await authHeaders(),
      body: JSON.stringify(dto),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error actualizando la campana" };
    }

    revalidatePath("/notification-config");
    return { success: true, data: body.data as Campaign };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function deleteCampaignAction(id: number, businessId?: number) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-campaigns/${id}`, businessId);
    const response = await fetch(url, {
      method: "DELETE",
      headers: await authHeaders(),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error eliminando la campana" };
    }

    revalidatePath("/notification-config");
    return { success: true };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function previewCampaignAudienceAction(
  dto: CreateCampaignDTO,
  businessId?: number,
) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/whatsapp-campaigns/audience-preview`, businessId);
    const response = await fetch(url, {
      method: "POST",
      headers: await authHeaders(),
      body: JSON.stringify(dto),
      cache: "no-store",
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error calculando la audiencia" };
    }

    return { success: true, data: body.data as CampaignAudiencePreview };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

async function lifecycleAction(id: number, action: string, businessId?: number) {
  try {
    const url = withBusiness(
      `${env.API_BASE_URL}/whatsapp-campaigns/${id}/${action}`,
      businessId,
    );
    const response = await fetch(url, {
      method: "POST",
      headers: await authHeaders(),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "No se pudo cambiar el estado" };
    }

    revalidatePath("/notification-config");
    return { success: true, data: body.data as Campaign };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function launchCampaignAction(id: number, businessId?: number) {
  return lifecycleAction(id, "launch", businessId);
}

export async function pauseCampaignAction(id: number, businessId?: number) {
  return lifecycleAction(id, "pause", businessId);
}

export async function resumeCampaignAction(id: number, businessId?: number) {
  return lifecycleAction(id, "resume", businessId);
}

export async function cancelCampaignAction(id: number, businessId?: number) {
  return lifecycleAction(id, "cancel", businessId);
}

export async function listCampaignSendsAction(
  campaignId: number,
  businessId?: number,
  status?: string,
  page = 1,
  pageSize = 20,
) {
  try {
    const query = new URLSearchParams({
      page: String(page),
      page_size: String(pageSize),
    });
    if (status) query.set("status", status);

    const url = withBusiness(
      `${env.API_BASE_URL}/whatsapp-campaigns/${campaignId}/sends?${query.toString()}`,
      businessId,
    );

    const response = await fetch(url, {
      headers: await authHeaders(),
      cache: "no-store",
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error listando envios", data: [] };
    }

    const result = body as PaginatedResult<CampaignSend>;
    return { success: true, data: result.data || [], total: result.total || 0 };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message, data: [] };
  }
}
