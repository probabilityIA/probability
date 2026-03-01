'use client';

import dynamic from 'next/dynamic';
import { Warehouse } from '../../domain/types';
import { Button } from '@/shared/ui';
import { PencilIcon, ArrowLeftIcon } from '@heroicons/react/24/outline';

const MapComponent = dynamic(() => import('@/shared/ui/MapComponent'), {
    ssr: false,
    loading: () => (
        <div className="flex items-center justify-center bg-gray-100 rounded-lg" style={{ height: '350px' }}>
            <span className="text-sm text-gray-400">Cargando mapa...</span>
        </div>
    ),
});

interface WarehouseDetailProps {
    warehouse: Warehouse;
    onEdit: () => void;
    onBack: () => void;
}

function InfoRow({ label, value }: { label: string; value?: string | number | null }) {
    if (!value && value !== 0) return null;
    return (
        <div className="flex justify-between py-1.5 border-b border-gray-50 last:border-b-0">
            <span className="text-sm text-gray-500">{label}</span>
            <span className="text-sm text-gray-900 text-right max-w-[60%]">{value}</span>
        </div>
    );
}

function SectionTitle({ children }: { children: React.ReactNode }) {
    return <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2 mt-4 first:mt-0">{children}</h3>;
}

export default function WarehouseDetail({ warehouse, onEdit, onBack }: WarehouseDetailProps) {
    const hasMapData = !!(warehouse.latitude || warehouse.longitude || warehouse.address || warehouse.street || warehouse.city);
    const hasCarrierData = !!(warehouse.company || warehouse.first_name || warehouse.last_name || warehouse.email || warehouse.street || warehouse.suburb || warehouse.city_dane_code || warehouse.postal_code);

    const formatDate = (dateStr: string) => {
        if (!dateStr) return null;
        try {
            return new Date(dateStr).toLocaleDateString('es-CO', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
            });
        } catch {
            return dateStr;
        }
    };

    return (
        <div className="space-y-4">
            {/* Header with badges */}
            <div className="flex items-start justify-between">
                <div>
                    <h2 className="text-lg font-semibold text-gray-900">{warehouse.name}</h2>
                    <span className="text-sm font-mono text-gray-500">{warehouse.code}</span>
                    <div className="flex gap-2 mt-2">
                        {warehouse.is_default && (
                            <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                                Principal
                            </span>
                        )}
                        {warehouse.is_fulfillment && (
                            <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
                                Fulfillment
                            </span>
                        )}
                        <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${warehouse.is_active ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'}`}>
                            {warehouse.is_active ? 'Activa' : 'Inactiva'}
                        </span>
                    </div>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" onClick={onBack}>
                        <ArrowLeftIcon className="w-4 h-4 mr-1" />
                        Volver
                    </Button>
                    <Button variant="primary" onClick={onEdit}>
                        <PencilIcon className="w-4 h-4 mr-1" />
                        Editar
                    </Button>
                </div>
            </div>

            {/* Two-column layout */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Left column - Data */}
                <div className="space-y-1">
                    <SectionTitle>Dirección</SectionTitle>
                    <InfoRow label="Dirección" value={warehouse.address} />
                    <InfoRow label="Ciudad" value={warehouse.city} />
                    <InfoRow label="Departamento" value={warehouse.state} />
                    <InfoRow label="País" value={warehouse.country} />
                    <InfoRow label="Código postal" value={warehouse.zip_code} />

                    <SectionTitle>Contacto</SectionTitle>
                    <InfoRow label="Teléfono" value={warehouse.phone} />
                    <InfoRow label="Nombre contacto" value={warehouse.contact_name} />
                    <InfoRow label="Email contacto" value={warehouse.contact_email} />

                    {hasCarrierData && (
                        <>
                            <SectionTitle>Datos de transportadora</SectionTitle>
                            <InfoRow label="Empresa" value={warehouse.company} />
                            <InfoRow label="Nombre" value={[warehouse.first_name, warehouse.last_name].filter(Boolean).join(' ') || null} />
                            <InfoRow label="Email" value={warehouse.email} />
                            <InfoRow label="Calle" value={warehouse.street} />
                            <InfoRow label="Barrio" value={warehouse.suburb} />
                            <InfoRow label="Código DANE" value={warehouse.city_dane_code} />
                            <InfoRow label="Código postal" value={warehouse.postal_code} />
                        </>
                    )}

                    {(warehouse.latitude != null || warehouse.longitude != null) && (
                        <>
                            <SectionTitle>Coordenadas</SectionTitle>
                            <InfoRow label="Latitud" value={warehouse.latitude} />
                            <InfoRow label="Longitud" value={warehouse.longitude} />
                        </>
                    )}

                    <SectionTitle>Registro</SectionTitle>
                    <InfoRow label="Creada" value={formatDate(warehouse.created_at)} />
                    <InfoRow label="Actualizada" value={formatDate(warehouse.updated_at)} />
                </div>

                {/* Right column - Map */}
                <div>
                    {hasMapData ? (
                        <MapComponent
                            address={warehouse.street || warehouse.address || ''}
                            city={warehouse.city || ''}
                            height="350px"
                            latitude={warehouse.latitude}
                            longitude={warehouse.longitude}
                        />
                    ) : (
                        <div className="flex items-center justify-center bg-gray-50 rounded-lg text-sm text-gray-400" style={{ height: '350px' }}>
                            Sin datos de ubicación
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
