"use server";

import { cookies } from "next/headers";
import { env } from "@/shared/config/env";
import {
  NotificationEventType,
  CreateNotificationEventTypeDTO,
  UpdateNotificationEventTypeDTO,
} from "../../domain/types";

// ============================================
// NOTIFICATION EVENT TYPES - Server Actions
// ============================================

/**
 * Obtener eventos de notificación por tipo
 */
export async function getNotificationEventTypesAction(
  notificationTypeId?: number
) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("token")?.value || "";

    const url = notificationTypeId
      ? `${env.API_BASE_URL}/notification-event-types?notification_type_id=${notificationTypeId}`
      : `${env.API_BASE_URL}/notification-event-types`;

    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      cache: "no-store",
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data: NotificationEventType[] = await response.json();
    return { success: true, data };
  } catch (error: any) {
    console.error("Error getting notification event types:", error);
    return { success: false, error: error.message, data: [] };
  }
}

/**
 * Obtener un evento de notificación por ID
 */
export async function getNotificationEventTypeByIdAction(id: number) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("token")?.value || "";

    const response = await fetch(
      `${env.API_BASE_URL}/notification-event-types/${id}`,
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

    const data: NotificationEventType = await response.json();
    return { success: true, data };
  } catch (error: any) {
    console.error("Error getting notification event type:", error);
    return { success: false, error: error.message };
  }
}

/**
 * Crear un nuevo evento de notificación
 */
export async function createNotificationEventTypeAction(
  dto: CreateNotificationEventTypeDTO
) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("token")?.value || "";

    const response = await fetch(
      `${env.API_BASE_URL}/notification-event-types`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(dto),
      }
    );

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(
        errorData.error || `HTTP error! status: ${response.status}`
      );
    }

    const data: NotificationEventType = await response.json();
    return { success: true, data };
  } catch (error: any) {
    console.error("Error creating notification event type:", error);
    return { success: false, error: error.message };
  }
}

/**
 * Actualizar un evento de notificación
 */
export async function updateNotificationEventTypeAction(
  id: number,
  dto: UpdateNotificationEventTypeDTO
) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("token")?.value || "";

    const response = await fetch(
      `${env.API_BASE_URL}/notification-event-types/${id}`,
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
      throw new Error(
        errorData.error || `HTTP error! status: ${response.status}`
      );
    }

    const data: NotificationEventType = await response.json();
    return { success: true, data };
  } catch (error: any) {
    console.error("Error updating notification event type:", error);
    return { success: false, error: error.message };
  }
}

/**
 * Eliminar un evento de notificación
 */
export async function deleteNotificationEventTypeAction(id: number) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get("token")?.value || "";

    const response = await fetch(
      `${env.API_BASE_URL}/notification-event-types/${id}`,
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
      throw new Error(
        errorData.error || `HTTP error! status: ${response.status}`
      );
    }

    return { success: true };
  } catch (error: any) {
    console.error("Error deleting notification event type:", error);
    return { success: false, error: error.message };
  }
}
