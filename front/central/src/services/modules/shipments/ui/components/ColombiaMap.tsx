'use client';

import React, { useMemo } from 'react';

interface ColombiaMapProps {
  highlightedDepartment?: string;
  status?: string;
}

// Normalizar texto: lowercase y remover tildes
function normalizeName(name: string): string {
  return name
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .trim();
}

// Mapeo de departamentos con alias normalizados
const DEPARTMENTS_MAP: Record<string, string[]> = {
  'Amazonas': ['amazonas'],
  'Antioquia': ['antioquia'],
  'Arauca': ['arauca'],
  'Atlántico': ['atlantico', 'atlanticode'],
  'Bolívar': ['bolivar', 'boliva'],
  'Boyacá': ['boyaca'],
  'Caldas': ['caldas'],
  'Caquetá': ['caqueta', 'caquete'],
  'Casanare': ['casanare'],
  'Cauca': ['cauca'],
  'Cesar': ['cesar', 'cesare'],
  'Chocó': ['choco'],
  'Córdoba': ['cordoba'],
  'Cundinamarca': ['cundinamarca'],
  'Guainía': ['guainia'],
  'Guaviare': ['guaviare'],
  'Huila': ['huila'],
  'La Guajira': ['guajira', 'laguajira'],
  'Magdalena': ['magdalena'],
  'Meta': ['meta'],
  'Nariño': ['narino', 'narin'],
  'Norte de Santander': ['norte santander', 'nortesantander', 'santander'],
  'Putumayo': ['putumayo'],
  'Quindío': ['quindio'],
  'Risaralda': ['risaralda'],
  'San Andrés': ['san andres', 'sandandres'],
  'Santander': ['santander'],
  'Sucre': ['sucre'],
  'Tolima': ['tolima'],
  'Valle del Cauca': ['valle cauca', 'valledelcauca'],
  'Vaupés': ['vaupes'],
  'Vichada': ['vichada'],
};

const STATUS_COLORS: Record<string, string> = {
  'pending': '#fbbf24',      // amber
  'in_transit': '#3b82f6',   // blue
  'delivered': '#10b981',    // emerald
  'failed': '#ef4444',       // red
};

// Función para encontrar el departamento por nombre
function findDepartment(search?: string): string | null {
  if (!search) return null;
  const normalized = normalizeName(search);

  for (const [dept, aliases] of Object.entries(DEPARTMENTS_MAP)) {
    for (const alias of aliases) {
      if (normalized.includes(alias) || alias.includes(normalized)) {
        return dept;
      }
    }
  }
  return null;
}

export function ColombiaMap({ highlightedDepartment, status = 'in_transit' }: ColombiaMapProps) {
  const foundDepartment = useMemo(() => findDepartment(highlightedDepartment), [highlightedDepartment]);
  const highlightColor = useMemo(() => STATUS_COLORS[status] || '#3b82f6', [status]);

  return (
    <div className="w-full flex justify-center items-center">
      <svg
        viewBox="0 0 600 700"
        className="w-full max-w-sm h-auto"
        xmlns="http://www.w3.org/2000/svg"
      >
        {/* Fondo */}
        <rect width="600" height="700" fill="#ffffff" />

        {/* Departamentos */}
        <g className="departments">
          {/* Amazonas */}
          <path
            d="M 380 620 L 420 630 L 440 600 L 400 590 Z"
            fill={foundDepartment === 'Amazonas' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Amazonas"
          />

          {/* Antioquia */}
          <path
            d="M 120 200 L 180 180 L 200 260 L 140 280 Z"
            fill={foundDepartment === 'Antioquia' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Antioquia"
          />

          {/* Arauca */}
          <path
            d="M 300 250 L 360 240 L 370 280 L 310 290 Z"
            fill={foundDepartment === 'Arauca' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Arauca"
          />

          {/* Atlántico */}
          <path
            d="M 90 150 L 130 140 L 135 170 L 95 175 Z"
            fill={foundDepartment === 'Atlántico' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Atlántico"
          />

          {/* Bolívar */}
          <path
            d="M 120 240 L 160 230 L 170 290 L 130 300 Z"
            fill={foundDepartment === 'Bolívar' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Bolívar"
          />

          {/* Boyacá */}
          <path
            d="M 220 300 L 280 290 L 290 360 L 220 370 Z"
            fill={foundDepartment === 'Boyacá' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Boyacá"
          />

          {/* Caldas */}
          <path
            d="M 160 320 L 200 310 L 210 360 L 170 370 Z"
            fill={foundDepartment === 'Caldas' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Caldas"
          />

          {/* Caquetá */}
          <path
            d="M 340 480 L 400 470 L 410 540 L 350 550 Z"
            fill={foundDepartment === 'Caquetá' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Caquetá"
          />

          {/* Casanare */}
          <path
            d="M 320 360 L 380 350 L 390 420 L 330 430 Z"
            fill={foundDepartment === 'Casanare' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Casanare"
          />

          {/* Cauca */}
          <path
            d="M 160 420 L 220 410 L 230 480 L 170 490 Z"
            fill={foundDepartment === 'Cauca' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Cauca"
          />

          {/* Cesar */}
          <path
            d="M 150 160 L 210 150 L 220 210 L 160 220 Z"
            fill={foundDepartment === 'Cesar' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Cesar"
          />

          {/* Chocó */}
          <path
            d="M 80 320 L 140 310 L 150 380 L 90 390 Z"
            fill={foundDepartment === 'Chocó' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Chocó"
          />

          {/* Córdoba */}
          <path
            d="M 100 220 L 160 210 L 170 270 L 110 280 Z"
            fill={foundDepartment === 'Córdoba' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Córdoba"
          />

          {/* Cundinamarca */}
          <path
            d="M 240 330 L 310 320 L 320 390 L 250 400 Z"
            fill={foundDepartment === 'Cundinamarca' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Cundinamarca"
          />

          {/* Guainía */}
          <path
            d="M 420 580 L 480 570 L 490 640 L 430 650 Z"
            fill={foundDepartment === 'Guainía' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Guainía"
          />

          {/* Guaviare */}
          <path
            d="M 400 500 L 460 490 L 470 560 L 410 570 Z"
            fill={foundDepartment === 'Guaviare' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Guaviare"
          />

          {/* Huila */}
          <path
            d="M 240 420 L 300 410 L 310 480 L 250 490 Z"
            fill={foundDepartment === 'Huila' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Huila"
          />

          {/* La Guajira */}
          <path
            d="M 140 100 L 190 90 L 200 140 L 150 150 Z"
            fill={foundDepartment === 'La Guajira' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="La Guajira"
          />

          {/* Magdalena */}
          <path
            d="M 110 170 L 160 160 L 170 220 L 120 230 Z"
            fill={foundDepartment === 'Magdalena' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Magdalena"
          />

          {/* Meta */}
          <path
            d="M 320 410 L 380 400 L 390 470 L 330 480 Z"
            fill={foundDepartment === 'Meta' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Meta"
          />

          {/* Nariño */}
          <path
            d="M 180 540 L 240 530 L 250 600 L 190 610 Z"
            fill={foundDepartment === 'Nariño' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Nariño"
          />

          {/* Norte de Santander */}
          <path
            d="M 260 280 L 320 270 L 330 330 L 270 340 Z"
            fill={foundDepartment === 'Norte de Santander' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Norte de Santander"
          />

          {/* Putumayo */}
          <path
            d="M 300 540 L 360 530 L 370 600 L 310 610 Z"
            fill={foundDepartment === 'Putumayo' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Putumayo"
          />

          {/* Quindío */}
          <path
            d="M 180 380 L 220 370 L 225 410 L 185 420 Z"
            fill={foundDepartment === 'Quindío' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Quindío"
          />

          {/* Risaralda */}
          <path
            d="M 170 350 L 210 340 L 220 380 L 180 390 Z"
            fill={foundDepartment === 'Risaralda' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Risaralda"
          />

          {/* San Andrés */}
          <path
            d="M 40 100 L 60 95 L 65 115 L 45 120 Z"
            fill={foundDepartment === 'San Andrés' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="San Andrés"
          />

          {/* Santander */}
          <path
            d="M 240 300 L 300 290 L 310 350 L 250 360 Z"
            fill={foundDepartment === 'Santander' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Santander"
          />

          {/* Sucre */}
          <path
            d="M 130 280 L 180 270 L 190 320 L 140 330 Z"
            fill={foundDepartment === 'Sucre' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Sucre"
          />

          {/* Tolima */}
          <path
            d="M 200 380 L 260 370 L 270 440 L 210 450 Z"
            fill={foundDepartment === 'Tolima' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Tolima"
          />

          {/* Valle del Cauca */}
          <path
            d="M 140 440 L 200 430 L 210 500 L 150 510 Z"
            fill={foundDepartment === 'Valle del Cauca' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Valle del Cauca"
          />

          {/* Vaupés */}
          <path
            d="M 460 550 L 520 540 L 530 610 L 470 620 Z"
            fill={foundDepartment === 'Vaupés' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Vaupés"
          />

          {/* Vichada */}
          <path
            d="M 400 450 L 460 440 L 470 510 L 410 520 Z"
            fill={foundDepartment === 'Vichada' ? highlightColor : '#e5e7eb'}
            className="transition-all duration-300 hover:opacity-80 cursor-pointer"
            title="Vichada"
          />
        </g>

        {/* Animación de pulso para el departamento resaltado */}
        {foundDepartment && (
          <style>{`
            .departments path[title="${foundDepartment}"] {
              animation: departmentPulse 2s ease-in-out infinite;
            }
            @keyframes departmentPulse {
              0%, 100% { opacity: 1; }
              50% { opacity: 0.7; }
            }
          `}</style>
        )}
      </svg>
    </div>
  );
}
