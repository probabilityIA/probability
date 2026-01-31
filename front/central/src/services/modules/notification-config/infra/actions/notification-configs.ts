"use server";

import { cookies } from "next/headers";
import { env } from "@/shared/config/env";
import { NotificationConfigApiRepository } from "../repository/api-repository";
import { CreateConfigDTO, UpdateConfigDTO, ConfigFilter } from "../../domain/types";

const getRepository = async () => {
  const cookieStore = await cookies();
  const token = cookieStore.get("token")?.value || "";
  return new NotificationConfigApiRepository(env.API_BASE_URL, token);
};

export async function createConfigAction(dto: CreateConfigDTO) {
  try {
    const repo = await getRepository();
    const config = await repo.create(dto);
    return { success: true, data: config };
  } catch (error: any) {
    return { success: false, error: error.message };
  }
}

export async function updateConfigAction(id: number, dto: UpdateConfigDTO) {
  try {
    const repo = await getRepository();
    const config = await repo.update(id, dto);
    return { success: true, data: config };
  } catch (error: any) {
    return { success: false, error: error.message };
  }
}

export async function deleteConfigAction(id: number) {
  try {
    const repo = await getRepository();
    await repo.delete(id);
    return { success: true };
  } catch (error: any) {
    return { success: false, error: error.message };
  }
}

export async function listConfigsAction(filter?: ConfigFilter) {
  try {
    const repo = await getRepository();
    const configs = await repo.list(filter);
    return { success: true, data: configs };
  } catch (error: any) {
    return { success: false, error: error.message };
  }
}

// Action temporal con paginación (compatible con tabla global)
export async function getConfigsAction(params?: any) {
  try {
    const repo = await getRepository();
    const configs = await repo.list(params);

    // Simular paginación hasta que el backend la implemente
    const page = params?.page || 1;
    const pageSize = params?.page_size || 20;
    const total = configs.length;
    const totalPages = Math.ceil(total / pageSize);
    const start = (page - 1) * pageSize;
    const end = start + pageSize;
    const paginatedData = configs.slice(start, end);

    return {
      success: true,
      data: paginatedData,
      total,
      page,
      page_size: pageSize,
      total_pages: totalPages,
    };
  } catch (error: any) {
    return {
      success: false,
      error: error.message,
      data: [],
      total: 0,
      total_pages: 0,
    };
  }
}

export async function getConfigAction(id: number) {
  try {
    const repo = await getRepository();
    const config = await repo.getById(id);
    return { success: true, data: config };
  } catch (error: any) {
    return { success: false, error: error.message };
  }
}
