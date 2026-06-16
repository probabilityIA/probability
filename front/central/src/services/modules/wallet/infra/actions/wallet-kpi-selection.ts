'use server';

import { API_BASE_URL } from '@/shared/utils/constants';
import { getToken } from '@/services/auth/infra/actions/auth-token';

export async function getWalletKPISelectionAction() {
	try {
		const token = await getToken();
		const res = await fetch(`${API_BASE_URL}/pay/wallet/admin/kpi-selection`, {
			headers: {
				Authorization: `Bearer ${token}`,
			},
		});

		if (!res.ok) throw new Error(`Failed to fetch: ${res.status}`);
		return await res.json();
	} catch (error: any) {
		return { success: false, error: error.message };
	}
}

export async function updateWalletKPISelectionAction(selectedBusinessIDs: number[]) {
	try {
		const token = await getToken();
		const res = await fetch(`${API_BASE_URL}/pay/wallet/admin/kpi-selection`, {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${token}`,
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ selected_business_ids: selectedBusinessIDs }),
		});

		if (!res.ok) throw new Error(`Failed to update: ${res.status}`);
		return await res.json();
	} catch (error: any) {
		return { success: false, error: error.message };
	}
}
