/**
 * EJEMPLO: Cómo refactorizar ShipmentGuideModal para usar colores dinámicos
 *
 * Este archivo muestra el ANTES y DESPUÉS de aplicar variables CSS dinámicas
 */

// ============================================
// PASO 1: ORIGEN Y DESTINO - ANTES (COLORS HARDCODEADOS)
// ============================================

export function ShipmentGuideModalBEFORE() {
  return (
    <div className="p-3 flex flex-col flex-1 overflow-hidden min-h-0">
      <div className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          {/* SECCIÓN ORIGEN - CON PURPLE HARDCODEADO */}
          <div className="bg-purple-50/50 dark:bg-purple-900/10 border border-purple-100 dark:border-purple-800/30 rounded-xl p-4 space-y-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                {/* ICON HARDCODEADO CON PURPLE */}
                <div className="w-8 h-8 rounded-lg bg-purple-100 dark:bg-purple-800/40 flex items-center justify-center text-purple-600 dark:text-purple-400 text-sm font-bold">
                  A
                </div>
                {/* TÍTULO HARDCODEADO CON PURPLE */}
                <h3 className="font-semibold text-base text-purple-700 dark:text-purple-400">
                  Origen
                </h3>
              </div>
            </div>

            {/* INPUT CON FOCUS PURPLE */}
            <input
              type="text"
              className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border-gray-300 dark:border-gray-600"
              placeholder="Buscar ciudad..."
            />
          </div>

          {/* SECCIÓN DESTINO - CON BLUE HARDCODEADO */}
          <div className="bg-blue-50/50 dark:bg-blue-900/10 border border-blue-100 dark:border-blue-800/30 rounded-xl p-4 space-y-2">
            <div className="flex items-center gap-2">
              {/* ICON HARDCODEADO CON BLUE */}
              <div className="w-8 h-8 rounded-lg bg-blue-100 dark:bg-blue-800/40 flex items-center justify-center text-blue-600 dark:text-blue-400 text-sm font-bold">
                B
              </div>
              {/* TÍTULO HARDCODEADO CON BLUE */}
              <h3 className="font-semibold text-base text-blue-700 dark:text-blue-400">
                Destino
              </h3>
            </div>

            {/* INPUT CON FOCUS BLUE */}
            <input
              type="text"
              className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white border-gray-300 dark:border-gray-600"
              placeholder="Buscar ciudad..."
            />
          </div>
        </div>
      </div>
    </div>
  );
}

// ============================================
// PASO 1: ORIGEN Y DESTINO - DESPUÉS (VARIABLES DINÁMICAS)
// ============================================

export function ShipmentGuideModalAFTER() {
  return (
    <div className="p-3 flex flex-col flex-1 overflow-hidden min-h-0">
      <div className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          {/* SECCIÓN ORIGEN - USA VARIABLES DINÁMICAS */}
          <div className="shipment-section-origin">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                {/* ICON CON VARIABLE DINÁMICA */}
                <div className="shipment-section-origin-icon w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold">
                  A
                </div>
                {/* TÍTULO CON VARIABLE DINÁMICA */}
                <h3 className="shipment-section-origin-label">Origen</h3>
              </div>
            </div>

            {/* INPUT CON VARIABLE DINÁMICA */}
            <input
              type="text"
              className="shipment-input"
              placeholder="Buscar ciudad..."
            />
          </div>

          {/* SECCIÓN DESTINO - USA VARIABLES DINÁMICAS */}
          <div className="shipment-section-destination">
            <div className="flex items-center gap-2">
              {/* ICON CON VARIABLE DINÁMICA */}
              <div className="shipment-section-destination-icon w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold">
                B
              </div>
              {/* TÍTULO CON VARIABLE DINÁMICA */}
              <h3 className="shipment-section-destination-label">Destino</h3>
            </div>

            {/* INPUT CON VARIABLE DINÁMICA */}
            <input
              type="text"
              className="shipment-input"
              placeholder="Buscar ciudad..."
            />
          </div>
        </div>
      </div>
    </div>
  );
}

// ============================================
// PASO 2: TARJETAS DE TRANSPORTISTA - ANTES (COLORS HARDCODEADOS)
// ============================================

export function CarrierCardBEFORE() {
  return (
    <div className="grid grid-cols-4 gap-3 auto-rows-max">
      <div
        onClick={() => {}}
        className="border border-gray-200 dark:border-gray-600 rounded-lg p-3 hover:border-purple-500 hover:shadow-md cursor-pointer transition-all bg-white dark:bg-gray-800"
      >
        <div className="grid grid-cols-3 gap-3 h-full">
          <div className="col-span-1 flex flex-col items-center justify-center">
            {/* BG PURPLE HARDCODEADO */}
            <div className="w-20 h-20 bg-purple-50 rounded-lg flex items-center justify-center overflow-hidden">
              <img
                src="https://logo.png"
                alt="Carrier"
                className="w-18 h-18 object-contain"
              />
            </div>
            <div className="font-semibold text-xs text-center mt-2">
              COORDINADORA
            </div>
          </div>

          <div className="col-span-2 flex flex-col text-[11px] text-gray-700 dark:text-gray-200">
            <div className="flex justify-between">
              <span>Guía</span>
              <span>$45.000</span>
            </div>
            <div className="border-t border-gray-300 dark:border-gray-600 mt-1 pt-1 flex justify-between items-baseline">
              {/* PURPLE HARDCODEADO */}
              <span className="font-semibold">Total</span>
              <span className="text-base font-bold text-purple-600">
                $50.000 <span className="text-[9px] font-normal text-gray-500">COP</span>
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// ============================================
// PASO 2: TARJETAS DE TRANSPORTISTA - DESPUÉS (VARIABLES DINÁMICAS)
// ============================================

export function CarrierCardAFTER() {
  const selectedRate = null; // Simulado

  return (
    <div className="grid grid-cols-4 gap-3 auto-rows-max">
      <div
        onClick={() => {}}
        className={`shipment-carrier-card ${
          selectedRate ? 'shipment-carrier-card-selected' : ''
        }`}
      >
        <div className="grid grid-cols-3 gap-3 h-full">
          <div className="col-span-1 flex flex-col items-center justify-center">
            {/* BG DINÁMICO */}
            <div className="w-20 h-20 shipment-carrier-logo-container rounded-lg flex items-center justify-center overflow-hidden">
              <img
                src="https://logo.png"
                alt="Carrier"
                className="w-18 h-18 object-contain"
              />
            </div>
            <div className="font-semibold text-xs text-center mt-2">
              COORDINADORA
            </div>
          </div>

          <div className="col-span-2 flex flex-col text-[11px] text-gray-700 dark:text-gray-200">
            <div className="flex justify-between">
              <span>Guía</span>
              <span>$45.000</span>
            </div>
            <div className="border-t border-gray-300 dark:border-gray-600 mt-1 pt-1 flex justify-between items-baseline">
              {/* COLOR DINÁMICO */}
              <span className="font-semibold">Total</span>
              <span className="shipment-cost-amount">
                $50.000 <span className="text-[9px] font-normal text-gray-500">COP</span>
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// ============================================
// BOTONES - ANTES (COLORS HARDCODEADOS)
// ============================================

export function ButtonsBEFORE() {
  return (
    <div className="flex justify-end gap-3">
      <button className="bg-gray-200 hover:bg-gray-300 text-gray-800 px-6 py-3 rounded-lg">
        Cancelar
      </button>
      {/* PURPLE HARDCODEADO */}
      <button style={{ background: '#7c3aed' }} className="text-white px-6 py-3 rounded-lg">
        Siguiente
      </button>
      <button className="bg-green-600 hover:bg-green-700 text-white px-6 py-3 rounded-lg">
        Pagar
      </button>
    </div>
  );
}

// ============================================
// BOTONES - DESPUÉS (VARIABLES DINÁMICAS)
// ============================================

export function ButtonsAFTER() {
  return (
    <div className="flex justify-end gap-3">
      <button className="shipment-btn-outline">Cancelar</button>
      <button className="shipment-btn-primary">Siguiente</button>
      <button className="shipment-btn-secondary">Pagar</button>
    </div>
  );
}

// ============================================
// BADGES - ANTES (COLORS HARDCODEADOS)
// ============================================

export function BadgesBEFORE() {
  return (
    <div className="flex gap-2">
      {/* AMBER HARDCODEADO */}
      <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-100 text-amber-700 border border-amber-300 dark:bg-amber-900/30 dark:text-amber-300 dark:border-amber-600">
        Contra Entrega - Solo opciones contra entrega
      </span>

      {/* GREEN HARDCODEADO */}
      <span className="inline-block px-2 py-1 bg-green-100 text-green-700 rounded text-xs">
        ✓ Cotizada
      </span>

      {/* PURPLE HARDCODEADO */}
      <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-purple-100 text-purple-700 border border-purple-300 dark:bg-purple-900/30 dark:text-purple-300 dark:border-purple-600">
        Filtrado: COORDINADORA
      </span>
    </div>
  );
}

// ============================================
// BADGES - DESPUÉS (VARIABLES DINÁMICAS)
// ============================================

export function BadgesAFTER() {
  return (
    <div className="flex gap-2">
      {/* BADGE WARNING (Amarillo) */}
      <span className="shipment-badge-warning">
        Contra Entrega - Solo opciones contra entrega
      </span>

      {/* BADGE SUCCESS (Verde) */}
      <span className="shipment-badge-success">✓ Cotizada</span>

      {/* BADGE PRIMARY (Color del negocio) */}
      <span className="shipment-badge-primary">Filtrado: COORDINADORA</span>
    </div>
  );
}

// ============================================
// ALERTAS - ANTES (COLORS HARDCODEADOS)
// ============================================

export function AlertsBEFORE() {
  const error = "Por favor completa los siguientes campos";
  const success = "✅ Guía generada exitosamente";

  return (
    <>
      {/* ERROR - RED HARDCODEADO */}
      {error && (
        <div className="mb-3 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-700 rounded-lg text-red-700 dark:text-red-400 text-sm">
          {error}
        </div>
      )}

      {/* SUCCESS - GREEN HARDCODEADO */}
      {success && (
        <div className="mb-2 p-2 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-700 rounded-lg text-green-700 dark:text-green-400">
          {success}
        </div>
      )}
    </>
  );
}

// ============================================
// ALERTAS - DESPUÉS (VARIABLES DINÁMICAS)
// ============================================

export function AlertsAFTER() {
  const error = "Por favor completa los siguientes campos";
  const success = "✅ Guía generada exitosamente";

  return (
    <>
      {/* ERROR - Dinámico */}
      {error && <div className="shipment-alert shipment-alert-error text-sm">{error}</div>}

      {/* SUCCESS - Verde fijo (adecuado para éxito) */}
      {success && <div className="shipment-alert shipment-alert-success">{success}</div>}

      {/* Otras variantes disponibles */}
      {/* <div className="shipment-alert shipment-alert-primary">Info primaria</div> */}
      {/* <div className="shipment-alert shipment-alert-warning">Advertencia</div> */}
    </>
  );
}

// ============================================
// CASO DE USO COMPLETO: Importar en el Modal
// ============================================

/**
 * En ShipmentGuideModal, al inicio del archivo:
 */

import '@/shared/ui/styles/shipment-modals.css'; // ← AGREGAR ESTO

export default function ShipmentGuideModal() {
  // ... resto del código del modal
  // Ahora usa las clases shipment-* en lugar de colores hardcodeados
  return (
    <div className="fixed inset-0 bg-black/20 backdrop-blur-sm flex items-center justify-center z-50 p-2">
      <div className="shipment-modal-content rounded-2xl shadow-xl flex flex-col overflow-hidden">
        {/* ... */}
      </div>
    </div>
  );
}
