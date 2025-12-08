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

export async function getConfigAction(id: number) {
    try {
        const repo = await getRepository();
        const config = await repo.getById(id);
        return { success: true, data: config };
    } catch (error: any) {
        return { success: false, error: error.message };
    }
}
