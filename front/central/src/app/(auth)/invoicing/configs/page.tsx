/**
 * Server Component - Página de Configuración de Facturación
 * Fetch inicial en servidor (logs visibles en terminal backend)
 */

import { ConfigsClient } from './ConfigsClient';
import { getConfigsAction } from '@/services/modules/invoicing/infra/actions';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';

interface PageProps {
  searchParams: Promise<{ business_id?: string }>;
}

export default async function InvoicingConfigsPage({ searchParams }: PageProps) {
  const params = await searchParams;

  // Determinar filtros desde query params
  const selectedBusinessId = params.business_id ? parseInt(params.business_id) : null;
  const filters = selectedBusinessId ? { business_id: selectedBusinessId } : {};

  // Fetch de configuraciones (SE EJECUTA EN EL SERVIDOR)
  // El backend filtra según los permisos del token
  console.log('🔍 [SERVER] Fetching invoicing configs with filters:', filters);
  let configs = [];
  try {
    const response = await getConfigsAction(filters);
    configs = response.data || [];
    console.log('✅ [SERVER] Invoicing configs loaded:', configs.length, 'items');
  } catch (error: any) {
    console.error('❌ [SERVER] Error loading invoicing configs:', error.message);
  }

  // Cargar businesses para el dropdown (el backend retorna según permisos)
  let businesses = [];
  try {
    const businessesResponse = await getBusinessesAction({});
    businesses = businessesResponse.data || [];
  } catch (error) {
    console.error('Error loading businesses:', error);
  }

  // Si hay más de 1 business disponible, es super admin
  const isSuperAdmin = businesses.length > 1;

  return (
    <ConfigsClient
      initialConfigs={configs}
      businesses={businesses}
      isSuperAdmin={isSuperAdmin}
    />
  );
}
