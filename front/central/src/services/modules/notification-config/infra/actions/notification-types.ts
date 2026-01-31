"use server";

import { cookies } from "next/headers";
import { env } from "@/shared/config/env";
import {
  NotificationType,
  CreateNotificationTypeDTO,
  UpdateNotificationTypeDTO,
} from "../../domain/types";

// ============================================
// NOTIFICATION TYPES - Server Actions
// ============================================

/**
 * Obtener todos los tipos de notificación
 */
export async function getNotificationTypesAction() {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("token")?.value || "";

    const response = await fetch(`${env.API_BASE_URL}/notification-types`, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      cache: "no-store",
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: NotificationType[] = await response.json();
    return { success: true, data };
  } catch (error: any) {
    console.error("Error getting notification types:", error);
    return { success: false, error: error.message, data: [] };
  }
}

/**
 * Obtener un tipo de notificación por ID
 */
export async function getNotificationTypeByIdAction(id: number) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("token")?.value || "";

    const response = await fetch(
      `${env.API_BASE_URL}/notification-types/${id}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        cache: "no-store",
      }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: NotificationType = await response.json();
    return { success: true, data };
  } catch (error: any) {
    console.error("Error getting notification type:", error);
    return { success: false, error: error.message };
  }
}

/**
 * Crear un nuevo tipo de notificación
 */
export async function createNotificationTypeAction(
  dto: CreateNotificationTypeDTO
) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("token")?.value || "";

    const response = await fetch(`${env.API_BASE_URL}/notification-types`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(dto),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }

    const data: NotificationType = await response.json();
    return { success: true, data };
  } catch (error: any) {
    console.error("Error creating notification type:", error);
    return { success: false, error: error.message };
  }
}

/**
 * Actualizar un tipo de notificación
 */
export async function updateNotificationTypeAction(
  id: number,
  dto: UpdateNotificationTypeDTO
) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("token")?.value || "";

    const response = await fetch(
      `${env.API_BASE_URL}/notification-types/${id}`,
      {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(dto),
      }
    );

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }

    const data: NotificationType = await response.json();
    return { success: true, data };
  } catch (error: any) {
    console.error("Error updating notification type:", error);
    return { success: false, error: error.message };
  }
}

/**
 * Eliminar un tipo de notificación
 */
export async function deleteNotificationTypeAction(id: number) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("token")?.value || "";

    const response = await fetch(
      `${env.API_BASE_URL}/notification-types/${id}`,
      {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      }
    );

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }

    return { success: true };
  } catch (error: any) {
    console.error("Error deleting notification type:", error);
    return { success: false, error: error.message };
  }
}
