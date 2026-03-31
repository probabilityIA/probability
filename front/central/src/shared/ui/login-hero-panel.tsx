'use client';

import React from 'react';

export const LoginHeroPanel = () => {
  return (
    <div className="hidden lg:flex lg:w-[60%] h-full relative overflow-hidden bg-[#F0EEFF]">
      {/* Blobs orgánicos grandes de fondo */}
      {/* Blob arriba izquierda - GRANDE */}
      <div
        className="absolute opacity-60"
        style={{
          background: '#E2DCFF',
          borderRadius: '40% 60% 50% 50% / 50% 50% 50% 50%',
          width: '550px',
          height: '550px',
          top: '-150px',
          left: '-100px',
        }}
      />

      {/* Blob centro-derecha - MEDIANO GRANDE */}
      <div
        className="absolute opacity-60"
        style={{
          background: '#E2DCFF',
          borderRadius: '60% 40% 50% 50% / 50% 50% 50% 50%',
          width: '400px',
          height: '400px',
          top: '100px',
          right: '-50px',
        }}
      />

      {/* Blob abajo - PEQUEÑO MEDIANO */}
      <div
        className="absolute opacity-60"
        style={{
          background: '#E2DCFF',
          borderRadius: '50% 50% 40% 60% / 50% 50% 50% 50%',
          width: '300px',
          height: '300px',
          bottom: '-80px',
          left: '50px',
        }}
      />

      {/* Puntos decorativos dispersos */}
      <div className="absolute w-2 h-2 bg-[#6B21E8] rounded-full opacity-80" style={{ top: '80px', left: '370px' }} />
      <div className="absolute w-2 h-2 bg-[#22C55E] rounded-full opacity-80" style={{ top: '150px', left: '385px' }} />
      <div className="absolute w-2 h-2 bg-[#6B21E8] rounded-full opacity-80" style={{ top: '320px', left: '240px' }} />
      <div className="absolute w-2 h-2 bg-[#6B21E8] rounded-full opacity-80" style={{ top: '420px', left: '580px' }} />
      <div className="absolute w-2 h-2 bg-[#22C55E] rounded-full opacity-80" style={{ top: '70px', right: '200px' }} />
      <div className="absolute w-2 h-2 bg-[#22C55E] rounded-full opacity-80" style={{ top: '380px', right: '120px' }} />
      <div className="absolute w-2 h-2 bg-[#6B21E8] rounded-full opacity-80" style={{ top: '640px', left: '350px' }} />
      <div className="absolute w-2 h-2 bg-[#6B21E8] rounded-full opacity-80" style={{ top: '300px', right: '80px' }} />

      {/* IPHONE MOCKUP - Centro superior */}
      <div
        className="absolute rounded-[48px] bg-[#D1D5DB] shadow-2xl"
        style={{
          width: '260px',
          height: '520px',
          top: '15%',
          left: '50%',
          transform: 'translateX(-50%)',
          zIndex: 10,
          padding: '14px',
          boxShadow: '0 25px 70px rgba(0, 0, 0, 0.2)',
          border: '8px solid #E5E7EB',
        }}
      >
        {/* Notch del iPhone */}
        <div
          className="absolute top-0 left-1/2 transform -translate-x-1/2 w-32 h-6 bg-[#D1D5DB] rounded-b-3xl"
          style={{ zIndex: 15 }}
        />

        {/* Pantalla del iPhone */}
        <div className="w-full h-full rounded-[40px] bg-[#1a1a2e]" />
      </div>

      {/* CARD "Órdenes del mes" - Arriba derecha del iPhone */}
      <div
        className="absolute bg-white rounded-[20px] shadow-2xl p-6"
        style={{
          width: '300px',
          top: '50px',
          right: '40px',
          zIndex: 20,
          boxShadow: '0 12px 40px rgba(0, 0, 0, 0.12)',
        }}
      >
        <h3 className="text-gray-900 font-semibold text-base mb-6">Órdenes del mes</h3>

        {/* Gráfico de barras con espacios */}
        <div className="flex items-end gap-4 h-40 justify-center">
          {/* Barra 1 - 40% - Púrpura claro #A78BFA */}
          <div className="flex-1 bg-[#A78BFA] rounded-t-lg" style={{ height: '40%', minWidth: '35px' }} />
          {/* Barra 2 - 70% - Púrpura oscuro #7C3AED */}
          <div className="flex-1 bg-[#7C3AED] rounded-t-lg" style={{ height: '70%', minWidth: '35px' }} />
          {/* Barra 3 - 50% - Verde claro #6EE7B7 */}
          <div className="flex-1 bg-[#6EE7B7] rounded-t-lg" style={{ height: '50%', minWidth: '35px' }} />
          {/* Barra 4 - 85% - Verde oscuro #22C55E */}
          <div className="flex-1 bg-[#22C55E] rounded-t-lg" style={{ height: '85%', minWidth: '35px' }} />
        </div>
      </div>

      {/* BADGE "Envío entregado" - Arriba izquierda */}
      <div
        className="absolute bg-white rounded-full border-2 border-[#7C3AED] px-7 py-4 shadow-2xl flex items-center gap-3"
        style={{
          left: '60px',
          top: '90px',
          zIndex: 20,
          boxShadow: '0 12px 40px rgba(124, 58, 237, 0.2)',
          whiteSpace: 'nowrap',
        }}
      >
        <div className="w-3 h-3 bg-[#22C55E] rounded-full flex-shrink-0 animate-pulse" />
        <span className="text-gray-900 font-semibold text-sm">
          Envío entregado • Pedido #1284
        </span>
      </div>

      {/* CARD TRACKING STEPS - Abajo derecha del iPhone */}
      <div
        className="absolute bg-white rounded-[16px] shadow-2xl p-4"
        style={{
          width: '320px',
          bottom: '140px',
          right: '50px',
          zIndex: 20,
          boxShadow: '0 12px 40px rgba(0, 0, 0, 0.12)',
        }}
      >
        {/* Stepper horizontal */}
        <div className="flex items-center justify-between gap-1">
          {/* Step 1 - Recogido */}
          <div className="flex flex-col items-center flex-1">
            <div className="w-7 h-7 bg-[#7C3AED] rounded-full flex items-center justify-center mb-2 flex-shrink-0">
              <svg className="w-3.5 h-3.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <span className="text-xs text-gray-700 text-center font-medium">Recogido</span>
          </div>

          {/* Línea conectora 1 */}
          <div className="flex-1 h-1 bg-[#7C3AED] mb-8 -mx-0.5" />

          {/* Step 2 - En Tránsito */}
          <div className="flex flex-col items-center flex-1">
            <div className="w-7 h-7 bg-[#7C3AED] rounded-full flex items-center justify-center mb-2 flex-shrink-0">
              <svg className="w-3.5 h-3.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <span className="text-xs text-gray-700 text-center font-medium">En Tránsito</span>
          </div>

          {/* Línea conectora 2 */}
          <div className="flex-1 h-1 bg-[#7C3AED] mb-8 -mx-0.5" />

          {/* Step 3 - En Reparto */}
          <div className="flex flex-col items-center flex-1">
            <div className="w-7 h-7 bg-[#7C3AED] rounded-full flex items-center justify-center mb-2 flex-shrink-0">
              <svg className="w-3.5 h-3.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <span className="text-xs text-gray-700 text-center font-medium">En Reparto</span>
          </div>

          {/* Línea conectora 3 */}
          <div className="flex-1 h-1 bg-gray-300 mb-8 -mx-0.5" />

          {/* Step 4 - Entregado */}
          <div className="flex flex-col items-center flex-1">
            <div className="w-7 h-7 bg-gray-300 rounded-full mb-2 flex-shrink-0" />
            <span className="text-xs text-gray-900 text-center font-bold">Entregado</span>
          </div>
        </div>
      </div>

      {/* TEXTO - Abajo Izquierda */}
      <div
        className="absolute"
        style={{
          bottom: '80px',
          left: '50px',
          zIndex: 15,
          maxWidth: '350px',
        }}
      >
        {/* Título - Predice y elimina devoluciones antes de que sucedan (Púrpura) */}
        <h1 className="text-[2.4rem] font-black text-[#7C3AED] leading-tight mb-6">
          Predice y elimina<br />devoluciones antes de que sucedan
        </h1>

        {/* Subtítulo */}
        <p className="text-base text-[#6B7280] mb-8 leading-relaxed font-medium">
          Shopify, Softpymes y muchas<br />transportadoras
        </p>

        {/* Pills */}
        <div className="flex gap-2.5 flex-wrap">
          <div className="bg-[#E9D5FF] text-[#6B21E8] px-5 py-2.5 rounded-full text-sm font-semibold whitespace-nowrap">
            🛍️ Shopify
          </div>
          <div className="bg-[#E9D5FF] text-[#6B21E8] px-5 py-2.5 rounded-full text-sm font-semibold whitespace-nowrap">
            💬 WhatsApp
          </div>
          <div className="bg-[#E9D5FF] text-[#6B21E8] px-5 py-2.5 rounded-full text-sm font-semibold whitespace-nowrap">
            🛒 Woocommerce
          </div>
          <div className="bg-[#E9D5FF] text-[#6B21E8] px-5 py-2.5 rounded-full text-sm font-semibold whitespace-nowrap">
            📦 Mercado Libre
          </div>
          <div className="bg-[#E9D5FF] text-[#6B21E8] px-5 py-2.5 rounded-full text-sm font-semibold whitespace-nowrap">
            🏬 Falabella
          </div>
        </div>
      </div>
    </div>
  );
};
