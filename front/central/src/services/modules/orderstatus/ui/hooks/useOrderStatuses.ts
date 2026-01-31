"use client";

import { useState, useEffect } from "react";
import { getOrderStatusesSimpleAction } from "../../infra/actions";
import { OrderStatusSimple } from "../../domain/types";

/**
 * Hook para obtener la lista de estados de orden (OrderStatus)
 *
 * Este hook obtiene los estados de orden desde el backend y los almacena en estado local.
 * Se usa principalmente para:
 * - Mostrar checklist de estados en el formulario de NotificationConfig
 * - Selectores de estados en otros formularios
 * - Filtros por estado
 *
 * @param isActive - Si true, solo retorna estados activos. Si false/undefined, retorna todos.
 * @returns {orderStatuses, loading, error, refresh}
 *
 * @example
 * ```tsx
 * const { orderStatuses, loading } = useOrderStatuses(true);
 *
 * if (loading) return <Spinner />;
 *
 * return (
 *   <div>
 *     {orderStatuses.map(status => (
 *       <Checkbox key={status.id} label={status.name} />
 *     ))}
 *   </div>
 * );
 * ```
 */
export function useOrderStatuses(isActive: boolean = true) {
  const [orderStatuses, setOrderStatuses] = useState<OrderStatusSimple[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchOrderStatuses = async () => {
    setLoading(true);
    setError(null);

    try {
      const result = await getOrderStatusesSimpleAction(isActive);

      if (result.success) {
        setOrderStatuses(result.data);
      } else {
        setError(result.message || "Error al cargar estados de orden");
        setOrderStatuses([]);
      }
    } catch (err: any) {
      setError(err.message || "Error al cargar estados de orden");
      setOrderStatuses([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrderStatuses();
  }, [isActive]); // Re-fetch si cambia el filtro isActive

  return {
    orderStatuses,
    loading,
    error,
    refresh: fetchOrderStatuses, // Función para refrescar manualmente
  };
}
