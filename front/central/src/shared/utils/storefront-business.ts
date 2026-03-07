import { cookies } from 'next/headers';

export async function getStorefrontBusinessId(): Promise<number | undefined> {
    const cookieStore = await cookies();
    const val = cookieStore.get('storefront_business_id')?.value;
    return val ? parseInt(val, 10) : undefined;
}
