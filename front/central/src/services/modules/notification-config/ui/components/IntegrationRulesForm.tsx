"use client";

import { useState, useEffect } from "react";
import { useToast } from "@/shared/providers/toast-provider";
import { useOrderStatuses } from "@/services/modules/orderstatus/ui";
import { RuleCard, LocalRule } from "./RuleCard";
import { SyncConfigsDTO } from "../../domain/types";
import { getConfigsAction, syncConfigsAction } from "../../infra/actions";
import type { IntegrationSimple } from "@/services/integrations/core/domain/types";

interface IntegrationRulesFormProps {
  integration: IntegrationSimple;
  businessId: number;
  onSuccess: () => void;
  onCancel: () => void;
}

function generateTempId(): string {
  return Math.random().toString(36).substring(2, 11);
}

export function IntegrationRulesForm({
  integration,
  businessId,
  onSuccess,
  onCancel,
}: IntegrationRulesFormProps) {
  const [rules, setRules] = useState<LocalRule[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingExisting, setLoadingExisting] = useState(true);
  const { showToast } = useToast();
  const { orderStatuses, loading: loadingOrderStatuses } = useOrderStatuses(true);

  // Load existing configs for this integration
  useEffect(() => {
    const loadExisting = async () => {
      setLoadingExisting(true);
      try {
        const result = await getConfigsAction({
          integration_id: integration.id,
          business_id: businessId,
        });
        if (result.success && result.data) {
          const existingRules: LocalRule[] = result.data.map((config: any) => ({
            _tempId: generateTempId(),
            id: config.id,
            notification_type_id: config.notification_type_id,
            notification_event_type_id: config.notification_event_type_id,
            enabled: config.enabled,
            description: config.description || "",
            order_status_ids: config.order_status_ids || [],
            _deleted: false,
          }));
          setRules(existingRules);
        }
      } catch (error) {
        showToast("Error al cargar reglas existentes", "error");
      } finally {
        setLoadingExisting(false);
      }
    };

    loadExisting();
  }, [integration.id, businessId]);

  const handleAddRule = () => {
    setRules((prev) => [
      ...prev,
      {
        _tempId: generateTempId(),
        notification_type_id: 0,
        notification_event_type_id: 0,
        enabled: true,
        description: "",
        order_status_ids: [],
        _deleted: false,
      },
    ]);
  };

  const handleRuleChange = (index: number, updated: LocalRule) => {
    setRules((prev) => prev.map((r, i) => (i === index ? updated : r)));
  };

  const handleRuleDelete = (index: number) => {
    setRules((prev) =>
      prev.map((r, i) => {
        if (i !== index) return r;
        // If it has an id (exists in DB), mark as deleted
        if (r.id) return { ...r, _deleted: true };
        // If new (no id), remove from array
        return r;
      }).filter((r) => !(!r.id && r._deleted))
    );
    // Remove new rules without id that were marked deleted
    setRules((prev) => prev.filter((r) => r.id || !r._deleted));
  };

  const handleSave = async () => {
    // Filter out deleted rules - these will be deleted by the backend (not in the sync rules)
    const activeRules = rules.filter((r) => !r._deleted);

    // Validate rules
    for (let i = 0; i < activeRules.length; i++) {
      const rule = activeRules[i];
      if (!rule.notification_type_id) {
        showToast(`Regla ${i + 1}: selecciona un canal`, "error");
        return;
      }
      if (!rule.notification_event_type_id) {
        showToast(`Regla ${i + 1}: selecciona un evento`, "error");
        return;
      }
    }

    // Check for duplicates (same notification_type_id + notification_event_type_id)
    const seen = new Set<string>();
    for (const rule of activeRules) {
      const key = `${rule.notification_type_id}-${rule.notification_event_type_id}`;
      if (seen.has(key)) {
        showToast("Hay reglas duplicadas (mismo canal + evento)", "error");
        return;
      }
      seen.add(key);
    }

    setLoading(true);
    try {
      const dto: SyncConfigsDTO = {
        integration_id: integration.id,
        rules: activeRules.map((r) => ({
          id: r.id,
          notification_type_id: r.notification_type_id,
          notification_event_type_id: r.notification_event_type_id,
          enabled: r.enabled,
          description: r.description,
          order_status_ids: r.order_status_ids,
        })),
      };

      const result = await syncConfigsAction(dto, businessId);

      if (result.success) {
        const data = result.data;
        showToast(
          `Sincronizado: ${data?.created || 0} creadas, ${data?.updated || 0} actualizadas, ${data?.deleted || 0} eliminadas`,
          "success"
        );
        onSuccess();
      } else {
        showToast(result.error || "Error al sincronizar reglas", "error");
      }
    } catch (error: any) {
      showToast(error.message || "Error inesperado", "error");
    } finally {
      setLoading(false);
    }
  };

  const activeRulesCount = rules.filter((r) => !r._deleted).length;

  return (
    <div className="space-y-4">
      {/* Integration header */}
      <div className="flex items-center gap-3 p-3 rounded-lg bg-gray-50 border border-gray-200">
        {integration.image_url ? (
          <img
            src={integration.image_url}
            alt={integration.name}
            className="w-8 h-8 object-contain rounded"
          />
        ) : (
          <div className="w-8 h-8 rounded bg-gray-200 flex items-center justify-center">
            <span className="text-xs font-bold text-gray-500">
              {integration.type?.charAt(0).toUpperCase() || "?"}
            </span>
          </div>
        )}
        <div>
          <h3 className="font-medium text-gray-900">{integration.name}</h3>
          <p className="text-xs text-gray-500">{integration.category_name || integration.type}</p>
        </div>
      </div>

      {/* Loading state */}
      {loadingExisting ? (
        <div className="text-center py-8 text-gray-500">Cargando reglas...</div>
      ) : (
        <>
          {/* Rules table */}
          <div className="border border-gray-200 rounded-lg overflow-hidden">
            {rules.filter((r) => !r._deleted).length === 0 ? (
              <div className="text-center py-8 text-gray-400">
                No hay reglas configuradas. Agrega una para empezar.
              </div>
            ) : (
              <table className="w-full">
                <thead>
                  <tr className="bg-gray-50 border-b border-gray-200">
                    <th className="py-2 px-3 text-[10px] font-semibold text-gray-500 uppercase text-center w-10">#</th>
                    <th className="py-2 px-3 text-[10px] font-semibold text-gray-500 uppercase text-left w-[120px]">Canal</th>
                    <th className="py-2 px-3 text-[10px] font-semibold text-gray-500 uppercase text-left w-[160px]">Evento</th>
                    <th className="py-2 px-3 text-[10px] font-semibold text-gray-500 uppercase text-left">Estados</th>
                    <th className="py-2 px-3 text-[10px] font-semibold text-gray-500 uppercase text-center w-[70px]">Activo</th>
                    <th className="py-2 px-3 text-[10px] font-semibold text-gray-500 uppercase text-center w-[50px]"></th>
                  </tr>
                </thead>
                <tbody>
                  {rules.map((rule, index) =>
                    !rule._deleted ? (
                      <RuleCard
                        key={rule._tempId}
                        rule={rule}
                        index={rules.filter((r, i) => i <= index && !r._deleted).length - 1}
                        orderStatuses={orderStatuses}
                        onChange={(updated) => handleRuleChange(index, updated)}
                        onDelete={() => handleRuleDelete(index)}
                      />
                    ) : null
                  )}
                </tbody>
              </table>
            )}
          </div>

          {/* Add rule button */}
          <div className="flex justify-center">
            <button
              type="button"
              onClick={handleAddRule}
              className="p-2 rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 transition-colors"
              title="Agregar regla"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </button>
          </div>
        </>
      )}

      {/* Footer actions */}
      <div className="flex items-center justify-between pt-4 border-t">
        <span className="text-sm text-gray-500">
          {activeRulesCount} regla(s) activa(s)
        </span>
        <div className="flex gap-2">
          <button
            type="button"
            onClick={onCancel}
            disabled={loading}
            className="p-2 rounded-lg bg-gray-100 text-gray-500 hover:bg-gray-200 transition-colors disabled:opacity-40"
            title="Cancelar"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          <button
            type="button"
            onClick={handleSave}
            disabled={loading || loadingExisting}
            className="p-2 rounded-lg bg-green-50 text-green-600 hover:bg-green-100 transition-colors disabled:opacity-40"
            title={loading ? "Guardando..." : "Guardar reglas"}
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}
