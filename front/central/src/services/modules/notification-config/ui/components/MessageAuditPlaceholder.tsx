'use client';

export function MessageAuditPlaceholder() {
  return (
    <div className="relative bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      {/* Header */}
      <div className="p-4 border-b flex items-center gap-2">
        <h3 className="text-sm font-medium text-gray-900">Auditoría de Mensajes</h3>
        <span className="px-1.5 py-0.5 text-[10px] font-semibold bg-purple-100 text-purple-700 rounded">
          Beta
        </span>
      </div>

      {/* Content */}
      <div className="p-4 space-y-4">
        {/* Stat cards */}
        <div className="grid grid-cols-3 gap-3">
          <div className="rounded-lg border border-gray-200 p-3 text-center">
            <p className="text-xs text-gray-500">Enviados</p>
            <p className="text-xl font-bold text-gray-300 mt-1">--</p>
          </div>
          <div className="rounded-lg border border-gray-200 p-3 text-center">
            <p className="text-xs text-gray-500">Fallidos</p>
            <p className="text-xl font-bold text-gray-300 mt-1">--</p>
          </div>
          <div className="rounded-lg border border-gray-200 p-3 text-center">
            <p className="text-xs text-gray-500">Tasa de Éxito</p>
            <p className="text-xl font-bold text-gray-300 mt-1">--</p>
          </div>
        </div>

        {/* Chart placeholder */}
        <div className="h-32 rounded-lg border-2 border-dashed border-gray-200 flex items-center justify-center">
          <div className="text-center text-gray-300">
            <svg className="w-8 h-8 mx-auto mb-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            <p className="text-xs">Gráfico semanal</p>
          </div>
        </div>

        {/* Skeleton rows */}
        <div className="space-y-2">
          <p className="text-xs font-medium text-gray-400">Últimos mensajes</p>
          {[1, 2, 3].map((i) => (
            <div key={i} className="flex items-center gap-3 p-2 rounded-lg bg-gray-50 animate-pulse">
              <div className="w-8 h-8 rounded-full bg-gray-200 shrink-0" />
              <div className="flex-1 space-y-1.5">
                <div className="h-3 bg-gray-200 rounded w-3/4" />
                <div className="h-2.5 bg-gray-200 rounded w-1/2" />
              </div>
              <div className="h-5 w-14 bg-gray-200 rounded-full shrink-0" />
            </div>
          ))}
        </div>
      </div>

      {/* Overlay */}
      <div className="absolute inset-0 bg-white/60 backdrop-blur-[1px] flex items-center justify-center rounded-lg">
        <div className="text-center">
          <svg className="w-10 h-10 mx-auto text-gray-400 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p className="text-sm font-semibold text-gray-600">Próximamente</p>
          <p className="text-xs text-gray-400 mt-0.5">Auditoría de mensajes enviados</p>
        </div>
      </div>
    </div>
  );
}
