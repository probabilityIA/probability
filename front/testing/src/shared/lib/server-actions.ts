"use server";

const TESTING_URL =
  process.env.TESTING_BACKEND_URL ||
  process.env.TESTING_API_URL ||
  "http://localhost:9092";

async function fetchServer(
  url: string,
  token: string,
  options: RequestInit = {}
) {
  const res = await fetch(url, {
    ...options,
    cache: "no-store",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...((options.headers as Record<string, string>) || {}),
    },
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(
      data.error || data.message || `Request failed: ${res.status}`
    );
  }

  return data;
}

export async function fetchBusinessesAction(token: string) {
  const data = await fetchServer(`${TESTING_URL}/api/v1/businesses`, token);
  return data.data as { id: number; name: string }[];
}

export async function fetchReferenceDataAction(
  token: string,
  businessId: number
) {
  const data = await fetchServer(
    `${TESTING_URL}/api/v1/orders/reference-data?business_id=${businessId}`,
    token
  );
  return data.data;
}

export async function generateOrdersAction(
  token: string,
  businessId: number,
  dto: {
    count: number;
    integration_id?: number;
    random_products: boolean;
    max_items_per_order: number;
  }
) {
  const data = await fetchServer(
    `${TESTING_URL}/api/v1/orders/generate?business_id=${businessId}`,
    token,
    {
      method: "POST",
      body: JSON.stringify(dto),
    }
  );
  return data.data;
}
