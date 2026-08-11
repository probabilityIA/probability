'use server';

import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { CommercialApiRepository } from '../repository/api-repository';
import { CommercialUseCases } from '../../app/use-cases';
import { GetProspectsParams } from '../../domain/types';

async function getUseCases() {
  const cookieStore = await cookies();
  const token = cookieStore.get('session_token')?.value || null;
  const repository = new CommercialApiRepository(token);
  return new CommercialUseCases(repository);
}

export const getProspectsAction = async (params?: GetProspectsParams) => {
  try {
    return await (await getUseCases()).getProspects(params);
  } catch (error: any) {
    console.error('Get Prospects Action Error:', error.message);
    throw new Error(error.message);
  }
};

export const getProspectsStatsAction = async () => {
  try {
    return await (await getUseCases()).getProspectsStats();
  } catch (error: any) {
    console.error('Get Prospects Stats Action Error:', error.message);
    throw new Error(error.message);
  }
};

export const updateProspectSeenAction = async (id: number, seen: boolean) => {
  try {
    const result = await (await getUseCases()).updateProspectSeen(id, seen);
    revalidatePath('/commercial');
    return result;
  } catch (error: any) {
    console.error('Update Prospect Seen Action Error:', error.message);
    throw new Error(error.message);
  }
};
