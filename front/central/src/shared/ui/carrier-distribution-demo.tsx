'use client';

import { CarrierDistributionCard } from './carrier-distribution-card';

// Demo data for testing
const demoData = [
  { carrier: 'INTERRAPIDISIMO', income: 4500000 },
  { carrier: 'ENVIA', income: 3200000 },
  { carrier: 'COORDINADORA', income: 2800000 },
  { carrier: 'SERVIENTREGA', income: 1900000 },
  { carrier: 'TODOCARGO', income: 950000 },
  { carrier: 'DHLEXPRESS', income: 650000 },
];

export function CarrierDistributionDemo() {
  return (
    <div className="p-6">
      <CarrierDistributionCard
        data={demoData}
        currency="COP"
        title="Ingresos por Familia"
        subtitle="Distribución de ingresos por transportadora"
      />
    </div>
  );
}
