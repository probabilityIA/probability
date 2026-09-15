'use server';

import { cookies } from 'next/headers';
import { AccessApiRepository } from '../repository/api-repository';
import type { AccessResult } from '../../domain/types';

export async function getMyAccessAction(businessId?: number): Promise<AccessResult> {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value;
        if (!token) {
            return { success: false, error: 'sin sesion' };
        }
        const data = await new AccessApiRepository(token).getMyAccess(businessId);
        return { success: true, data };
    } catch (error: any) {
        return { success: false, error: error?.message || 'No se pudo obtener el acceso' };
    }
}
