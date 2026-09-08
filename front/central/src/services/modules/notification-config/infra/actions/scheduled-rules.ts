"use server";

import { cookies } from "next/headers";
import { revalidatePath } from "next/cache";
import { env } from "@/shared/config/env";
import {
  CreateScheduledRuleDTO,
  PaginatedResult,
  ScheduledRule,
  ScheduledRun,
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

const BASE = "scheduled-notification-rules";

export async function listScheduledRulesAction(businessId?: number, page = 1, pageSize = 20) {
  try {
    const query = new URLSearchParams({
      page: String(page),
      page_size: String(pageSize),
    });

    const url = withBusiness(`${env.API_BASE_URL}/${BASE}?${query.toString()}`, businessId);
    const response = await fetch(url, {
      headers: await authHeaders(),
      cache: "no-store",
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error listando reglas", data: [] };
    }

    const result = body as PaginatedResult<ScheduledRule>;
    return { success: true, data: result.data || [], total: result.total || 0 };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message, data: [] };
  }
}

export async function createScheduledRuleAction(
  dto: CreateScheduledRuleDTO,
  businessId?: number,
) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/${BASE}`, businessId);
    const response = await fetch(url, {
      method: "POST",
      headers: await authHeaders(),
      body: JSON.stringify(dto),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error creando la regla" };
    }

    revalidatePath("/notification-config");
    return { success: true, data: body.data as ScheduledRule };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function updateScheduledRuleAction(
  id: number,
  dto: CreateScheduledRuleDTO,
  businessId?: number,
) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/${BASE}/${id}`, businessId);
    const response = await fetch(url, {
      method: "PUT",
      headers: await authHeaders(),
      body: JSON.stringify(dto),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error actualizando la regla" };
    }

    revalidatePath("/notification-config");
    return { success: true, data: body.data as ScheduledRule };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function deleteScheduledRuleAction(id: number, businessId?: number) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/${BASE}/${id}`, businessId);
    const response = await fetch(url, {
      method: "DELETE",
      headers: await authHeaders(),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error eliminando la regla" };
    }

    revalidatePath("/notification-config");
    return { success: true };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function runScheduledRuleNowAction(id: number, businessId?: number) {
  try {
    const url = withBusiness(`${env.API_BASE_URL}/${BASE}/${id}/run`, businessId);
    const response = await fetch(url, {
      method: "POST",
      headers: await authHeaders(),
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error ejecutando la regla" };
    }

    revalidatePath("/notification-config");
    return { success: true, data: body.data as ScheduledRun };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message };
  }
}

export async function listScheduledRunsAction(
  ruleId: number,
  businessId?: number,
  page = 1,
  pageSize = 10,
) {
  try {
    const query = new URLSearchParams({
      page: String(page),
      page_size: String(pageSize),
    });

    const url = withBusiness(
      `${env.API_BASE_URL}/${BASE}/${ruleId}/runs?${query.toString()}`,
      businessId,
    );

    const response = await fetch(url, {
      headers: await authHeaders(),
      cache: "no-store",
    });

    const body = await response.json();
    if (!response.ok) {
      return { success: false, error: body.error || "Error listando corridas", data: [] };
    }

    return { success: true, data: (body.data || []) as ScheduledRun[], total: body.total || 0 };
  } catch (error) {
    const message = error instanceof Error ? error.message : "Error desconocido";
    return { success: false, error: message, data: [] };
  }
}
