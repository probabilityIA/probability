'use client';

import type { InvoicingFilters } from '@/services/modules/invoicing/domain/types';
import { MultiSelect } from './MultiSelect';

interface InvoicingFilterBuilderProps {
  filters: InvoicingFilters;
  onChange: (filters: InvoicingFilters) => void;
  disabled?: boolean;
}

/**
 * Constructor visual de filtros de facturación
 */
export function InvoicingFilterBuilder({
  filters,
  onChange,
  disabled = false,
}: InvoicingFilterBuilderProps) {
  const updateFilter = <K extends keyof InvoicingFilters>(
    key: K,
    value: InvoicingFilters[K]
  ) => {
    onChange({ ...filters, [key]: value });
  };

  return (
    <div className="space-y-6 border border-gray-200 rounded-lg p-6 bg-gray-50">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">
        Filtros de Facturación Automática
      </h3>

      {/* ============================================ */}
      {/* SECCIÓN: MONTO */}
      {/* ============================================ */}
      <div className="bg-white p-4 rounded-lg border border-gray-200">
        <h4 className="font-medium text-gray-900 mb-3">Filtros de Monto</h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Monto mínimo */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Monto mínimo (COP)
            </label>
            <input
              type="number"
              value={filters.min_amount || ''}
              onChange={(e) =>
                updateFilter('min_amount', Number(e.target.value) || undefined)
              }
              placeholder="100000"
              disabled={disabled}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
            />
            <p className="text-xs text-gray-500 mt-1">
              Solo facturar órdenes con este monto mínimo
            </p>
          </div>

          {/* Monto máximo */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Monto máximo (COP)
            </label>
            <input
              type="number"
              value={filters.max_amount || ''}
              onChange={(e) =>
                updateFilter('max_amount', Number(e.target.value) || undefined)
              }
              placeholder="10000000"
              disabled={disabled}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
            />
            <p className="text-xs text-gray-500 mt-1">
              Solo facturar órdenes hasta este monto
            </p>
          </div>
        </div>
      </div>

      {/* ============================================ */}
      {/* SECCIÓN: PAGO */}
      {/* ============================================ */}
      <div className="bg-white p-4 rounded-lg border border-gray-200">
        <h4 className="font-medium text-gray-900 mb-3">Filtros de Pago</h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Estado de pago */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Estado de pago
            </label>
            <select
              value={filters.payment_status || ''}
              onChange={(e) =>
                updateFilter(
                  'payment_status',
                  e.target.value as any || undefined
                )
              }
              disabled={disabled}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
            >
              <option value="">Sin filtro</option>
              <option value="paid">Solo pagadas</option>
              <option value="unpaid">Solo sin pagar</option>
              <option value="partial">Pago parcial</option>
            </select>
          </div>
        </div>
      </div>

      {/* ============================================ */}
      {/* SECCIÓN: ORDEN */}
      {/* ============================================ */}
      <div className="bg-white p-4 rounded-lg border border-gray-200">
        <h4 className="font-medium text-gray-900 mb-3">Filtros de Orden</h4>
        <div className="space-y-4">
          {/* Tipos de orden */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Tipos de orden permitidos
            </label>
            <MultiSelect
              options={[
                { value: 'delivery', label: 'Domicilio' },
                { value: 'pickup', label: 'Recoger en tienda' },
                { value: 'dine_in', label: 'Comer en restaurante' },
                { value: 'online', label: 'Online' },
              ]}
              value={filters.order_types || []}
              onChange={(types) => updateFilter('order_types', types as string[])}
              placeholder="Todos los tipos"
              disabled={disabled}
            />
            <p className="text-xs text-gray-500 mt-1">
              Dejar vacío para permitir todos los tipos
            </p>
          </div>

          {/* Estados excluidos */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Excluir estados
            </label>
            <MultiSelect
              options={[
                { value: 'cancelled', label: 'Cancelado' },
                { value: 'returned', label: 'Devuelto' },
                { value: 'refunded', label: 'Reembolsado' },
                { value: 'pending', label: 'Pendiente' },
              ]}
              value={filters.exclude_statuses || []}
              onChange={(statuses) =>
                updateFilter('exclude_statuses', statuses as string[])
              }
              placeholder="No excluir ninguno"
              disabled={disabled}
            />
            <p className="text-xs text-gray-500 mt-1">
              Órdenes con estos estados NO se facturarán
            </p>
          </div>
        </div>
      </div>

      {/* ============================================ */}
      {/* SECCIÓN: PRODUCTOS */}
      {/* ============================================ */}
      <div className="bg-white p-4 rounded-lg border border-gray-200">
        <h4 className="font-medium text-gray-900 mb-3">Filtros de Productos</h4>
        <div className="space-y-4">
          {/* Productos excluidos */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              SKUs excluidos
            </label>
            <input
              type="text"
              value={filters.exclude_products?.join(', ') || ''}
              onChange={(e) =>
                updateFilter(
                  'exclude_products',
                  e.target.value
                    .split(',')
                    .map((s) => s.trim())
                    .filter(Boolean)
                )
              }
              placeholder="GIFT-CARD-001, SKU-123"
              disabled={disabled}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
            />
            <p className="text-xs text-gray-500 mt-1">
              Separar múltiples SKUs con comas. Órdenes que contengan estos
              productos NO se facturarán
            </p>
          </div>

          {/* Productos incluidos únicamente */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              SKUs permitidos (exclusivo)
            </label>
            <input
              type="text"
              value={filters.include_products_only?.join(', ') || ''}
              onChange={(e) =>
                updateFilter(
                  'include_products_only',
                  e.target.value
                    .split(',')
                    .map((s) => s.trim())
                    .filter(Boolean)
                )
              }
              placeholder="PROD-001, PROD-002"
              disabled={disabled}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
            />
            <p className="text-xs text-gray-500 mt-1">
              Solo facturar órdenes que contengan ÚNICAMENTE estos productos
            </p>
          </div>

          {/* Cantidad de ítems */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Mínimo ítems
              </label>
              <input
                type="number"
                value={filters.min_items_count || ''}
                onChange={(e) =>
                  updateFilter(
                    'min_items_count',
                    Number(e.target.value) || undefined
                  )
                }
                placeholder="1"
                min="1"
                disabled={disabled}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Máximo ítems
              </label>
              <input
                type="number"
                value={filters.max_items_count || ''}
                onChange={(e) =>
                  updateFilter(
                    'max_items_count',
                    Number(e.target.value) || undefined
                  )
                }
                placeholder="100"
                min="1"
                disabled={disabled}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
              />
            </div>
          </div>
        </div>
      </div>

      {/* ============================================ */}
      {/* SECCIÓN: UBICACIÓN */}
      {/* ============================================ */}
      <div className="bg-white p-4 rounded-lg border border-gray-200">
        <h4 className="font-medium text-gray-900 mb-3">Filtros de Ubicación</h4>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Regiones de envío permitidas
          </label>
          <MultiSelect
            options={[
              { value: 'Bogotá', label: 'Bogotá' },
              { value: 'Medellín', label: 'Medellín' },
              { value: 'Cali', label: 'Cali' },
              { value: 'Barranquilla', label: 'Barranquilla' },
              { value: 'Cartagena', label: 'Cartagena' },
              { value: 'Bucaramanga', label: 'Bucaramanga' },
              { value: 'Pereira', label: 'Pereira' },
              { value: 'Santa Marta', label: 'Santa Marta' },
            ]}
            value={filters.shipping_regions || []}
            onChange={(regions) =>
              updateFilter('shipping_regions', regions as string[])
            }
            placeholder="Todas las regiones"
            disabled={disabled}
          />
          <p className="text-xs text-gray-500 mt-1">
            Solo facturar envíos a estas regiones
          </p>
        </div>
      </div>

      {/* Resumen de filtros activos */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <h4 className="font-medium text-blue-900 mb-2">Filtros Activos</h4>
        <div className="text-sm text-blue-800">
          {Object.entries(filters).filter(
            ([_, value]) =>
              value !== undefined &&
              value !== null &&
              (Array.isArray(value) ? value.length > 0 : true)
          ).length === 0 ? (
            <p>No hay filtros configurados. Se facturarán todas las órdenes.</p>
          ) : (
            <ul className="list-disc list-inside space-y-1">
              {filters.min_amount && (
                <li>Monto mínimo: ${filters.min_amount.toLocaleString()}</li>
              )}
              {filters.max_amount && (
                <li>Monto máximo: ${filters.max_amount.toLocaleString()}</li>
              )}
              {filters.payment_status && (
                <li>Estado de pago: {filters.payment_status}</li>
              )}
              {filters.order_types && filters.order_types.length > 0 && (
                <li>Tipos de orden: {filters.order_types.join(', ')}</li>
              )}
              {filters.exclude_statuses &&
                filters.exclude_statuses.length > 0 && (
                  <li>
                    Excluir estados: {filters.exclude_statuses.join(', ')}
                  </li>
                )}
              {filters.exclude_products &&
                filters.exclude_products.length > 0 && (
                  <li>
                    Productos excluidos: {filters.exclude_products.join(', ')}
                  </li>
                )}
              {filters.shipping_regions &&
                filters.shipping_regions.length > 0 && (
                  <li>
                    Regiones: {filters.shipping_regions.join(', ')}
                  </li>
                )}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
}
