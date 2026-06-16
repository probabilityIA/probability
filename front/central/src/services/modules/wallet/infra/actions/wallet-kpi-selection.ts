'use server';

import { cookies } from 'next/headers';
import { env } from '@/shared/config/env';

async function getAuthHeader() {
	const cookieStore = await cookies();
	const token = cookieStore.get('session_token')?.value;
	return {
		'Authorization': `Bearer ${token}`,
		'Content-Type': 'application/json',
	};
}

export async function getWalletKPISelectionAction() {
	try {
		const headers = await getAuthHeader();
		const res = await fetch(`${env.API_BASE_URL}/pay/wallet/admin/kpi-selection`, {
			headers,
		});

		if (!res.ok) throw new Error(`Failed to fetch: ${res.status}`);
		return await res.json();
	} catch (error: any) {
		return { success: false, error: error.message };
	}
}

export async function updateWalletKPISelectionAction(selectedBusinessIDs: number[]) {
	try {
		const headers = await getAuthHeader();
		const res = await fetch(`${env.API_BASE_URL}/pay/wallet/admin/kpi-selection`, {
			method: 'POST',
			headers,
			body: JSON.stringify({ selected_business_ids: selectedBusinessIDs }),
		});

		if (!res.ok) throw new Error(`Failed to update: ${res.status}`);
		return await res.json();
	} catch (error: any) {
		return { success: false, error: error.message };
	}
}
