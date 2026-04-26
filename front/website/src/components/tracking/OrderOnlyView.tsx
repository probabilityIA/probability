/** @jsxImportSource react */
import type { OrderPublicTracking } from '../../types/tracking';

interface Props {
  order: OrderPublicTracking;
}

const STATUS_CONFIG: Record<string, { label: string; color: string; icon: string; description: string }> = {
  pending: { label: 'Pendiente', color: '#F59E0B', icon: '⏳', description: 'Tu pedido está siendo preparado para despacho' },
  open: { label: 'Abierto', color: '#3B82F6', icon: '📋', description: 'Tu pedido fue recibido y está en revisión' },
  confirmed: { label: 'Confirmado', color: '#10B981', icon: '✅', description: 'Tu pedido fue confirmado y será despachado pronto' },
  processing: { label: 'En proceso', color: '#8B5CF6', icon: '⚙️', description: 'Tu pedido está siendo preparado' },
  in_transit: { label: 'En tránsito', color: '#3B82F6', icon: '🚚', description: 'Tu pedido está en camino' },
  delivered: { label: 'Entregado', color: '#22C55E', icon: '🎉', description: 'Tu pedido fue entregado' },
  cancelled: { label: 'Cancelado', color: '#EF4444', icon: '❌', description: 'Este pedido fue cancelado' },
};

function formatMoney(amount: number, currency: string = 'COP') {
  try {
    return new Intl.NumberFormat('es-CO', { style: 'currency', currency, maximumFractionDigits: 0 }).format(amount);
  } catch {
    return `$${amount}`;
  }
}

function formatDate(s?: string | null) {
  if (!s) return '';
  try {
    return new Date(s).toLocaleString('es-CO', {
      day: '2-digit', month: 'short', year: 'numeric',
      hour: '2-digit', minute: '2-digit', hour12: false, timeZone: 'America/Bogota',
    });
  } catch {
    return s;
  }
}

export default function OrderOnlyView({ order }: Props) {
  const cfg = STATUS_CONFIG[order.Status?.toLowerCase()] || {
    label: order.Status || 'Recibido',
    color: '#6B7280',
    icon: '📦',
    description: 'Tu pedido fue registrado en nuestro sistema',
  };

  const fullAddress = [order.ShippingStreet, order.ShippingCity, order.ShippingState]
    .filter(Boolean)
    .join(', ');

  return (
    <div class="space-y-6">
      <div class="bg-white rounded-2xl shadow-lg p-8">
        <div class="flex items-start justify-between gap-4 mb-6">
          <div>
            <h2 class="text-2xl font-bold text-gray-900">Estado del Pedido</h2>
            <p class="text-sm text-gray-500 mt-1">Pedido #{order.OrderNumber}</p>
          </div>
          <div
            class="px-4 py-2 rounded-full text-white text-sm font-bold shadow-sm"
            style={{ backgroundColor: cfg.color }}
          >
            {cfg.label}
          </div>
        </div>

        <div class="bg-gradient-to-br from-amber-50 to-orange-50 border-2 border-amber-200 rounded-xl p-6 mb-6 flex items-start gap-4">
          <div class="text-4xl flex-shrink-0">{cfg.icon}</div>
          <div class="flex-1">
            <h3 class="font-bold text-amber-900 text-lg mb-1">Aún sin guía de envío</h3>
            <p class="text-amber-800 text-sm">
              {cfg.description}. Cuando se genere la guía con la transportadora, recibirás
              una notificación por WhatsApp con el número de rastreo.
            </p>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="bg-slate-50 rounded-xl p-5 border border-slate-200">
            <div class="flex items-center gap-2 mb-2">
              <span class="text-xl">👤</span>
              <p class="text-xs font-bold text-gray-500 uppercase tracking-widest">Cliente</p>
            </div>
            <p class="text-base font-semibold text-gray-900">{order.CustomerName || '—'}</p>
            {order.CustomerPhone && (
              <p class="text-xs text-gray-500 mt-1">📞 {order.CustomerPhone}</p>
            )}
          </div>

          <div class="bg-slate-50 rounded-xl p-5 border border-slate-200">
            <div class="flex items-center gap-2 mb-2">
              <span class="text-xl">🏪</span>
              <p class="text-xs font-bold text-gray-500 uppercase tracking-widest">Tienda</p>
            </div>
            <p class="text-base font-semibold text-gray-900">{order.BusinessName || '—'}</p>
          </div>

          {fullAddress && (
            <div class="bg-slate-50 rounded-xl p-5 border border-slate-200 md:col-span-2">
              <div class="flex items-center gap-2 mb-2">
                <span class="text-xl">📍</span>
                <p class="text-xs font-bold text-gray-500 uppercase tracking-widest">Dirección de Envío</p>
              </div>
              <p class="text-sm text-gray-900">{fullAddress}</p>
              {order.ShippingPostalCode && (
                <p class="text-xs text-gray-500 mt-1">CP: {order.ShippingPostalCode}</p>
              )}
            </div>
          )}

          <div class="bg-emerald-50 rounded-xl p-5 border border-emerald-200">
            <div class="flex items-center gap-2 mb-2">
              <span class="text-xl">💰</span>
              <p class="text-xs font-bold text-emerald-700 uppercase tracking-widest">Total</p>
            </div>
            <p class="text-lg font-bold text-emerald-900">
              {formatMoney(order.TotalAmount, order.Currency || 'COP')}
            </p>
            {order.IsPaid ? (
              <p class="text-xs text-emerald-600 mt-1 font-semibold">✅ Pagado</p>
            ) : (
              <p class="text-xs text-amber-600 mt-1 font-semibold">⏳ Pendiente de pago</p>
            )}
          </div>

          {order.CodTotal != null && order.CodTotal > 0 && (
            <div class="bg-orange-50 rounded-xl p-5 border border-orange-200">
              <div class="flex items-center gap-2 mb-2">
                <span class="text-xl">💵</span>
                <p class="text-xs font-bold text-orange-700 uppercase tracking-widest">Contra Entrega</p>
              </div>
              <p class="text-lg font-bold text-orange-900">
                {formatMoney(order.CodTotal, order.Currency || 'COP')}
              </p>
              <p class="text-xs text-orange-600 mt-1">A recaudar al entregar</p>
            </div>
          )}

          {order.CreatedAt && (
            <div class="bg-slate-50 rounded-xl p-5 border border-slate-200 md:col-span-2">
              <div class="flex items-center gap-2 mb-2">
                <span class="text-xl">📅</span>
                <p class="text-xs font-bold text-gray-500 uppercase tracking-widest">Fecha del Pedido</p>
              </div>
              <p class="text-sm text-gray-900">{formatDate(order.OccurredAt || order.CreatedAt)}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
