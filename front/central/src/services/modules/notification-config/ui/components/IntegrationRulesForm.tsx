"use client";

import { useState, useEffect, useMemo } from "react";
import { useToast } from "@/shared/providers/toast-provider";
import { useOrderStatuses } from "@/services/modules/orderstatus/ui";
import { useIntegrationsSimple } from "@/services/integrations/core/ui/hooks/useIntegrationsSimple";
import { RuleCard, LocalRule } from "./RuleCard";
import { SyncConfigsDTO } from "../../domain/types";
import { getConfigsAction, syncConfigsAction, getNotificationTypesAction } from "../../infra/actions";
import { ScheduledRulesSection } from "./ScheduledRulesSection";
import type { IntegrationSimple } from "@/services/integrations/core/domain/types";

interface IntegrationRulesFormProps {
  businessId: number;
  onSuccess: () => void;
  onCancel: () => void;
}

interface ChannelSection {
  id: number;
  name: string;
  subtitle: string;
  imageUrl?: string;
  fallbackLetter: string;
}

function generateTempId(): string {
  return Math.random().toString(36).substring(2, 11);
}

const CATEGORY_ORDER: Record<string, number> = {
  platform: 0,
  ecommerce: 1,
};

function sectionFromIntegration(integration: IntegrationSimple): ChannelSection {
  return {
    id: integration.id,
    name: integration.name,
    subtitle: integration.category_name || integration.type || "",
    imageUrl: integration.image_url,
    fallbackLetter: (integration.name || integration.type || "?").charAt(0).toUpperCase(),
  };
}

export function IntegrationRulesForm({
  businessId,
  onSuccess,
  onCancel,
}: IntegrationRulesFormProps) {
  const [rules, setRules] = useState<LocalRule[]>([]);
  const [originalIntegrationIds, setOriginalIntegrationIds] = useState<number[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingExisting, setLoadingExisting] = useState(true);
  const { showToast } = useToast();
  const { orderStatuses } = useOrderStatuses(true);
  const { integrations, loading: loadingIntegrations } = useIntegrationsSimple(
    businessId ? { businessId } : undefined,
  );
  const [whatsappTypeId, setWhatsappTypeId] = useState(0);
  const [activeTab, setActiveTab] = useState<"status" | "scheduled">("status");

  useEffect(() => {
    const loadTypes = async () => {
      const result = await getNotificationTypesAction();
      if (result.success && result.data) {
        const whatsapp = result.data.find(
          (type: { code: string; id: number }) => type.code === "whatsapp",
        );
        if (whatsapp) setWhatsappTypeId(whatsapp.id);
      }
    };
    loadTypes();
  }, []);

  useEffect(() => {
    const loadExisting = async () => {
      setLoadingExisting(true);
      try {
        const result = await getConfigsAction(
          businessId ? { business_id: businessId } : {},
        );
        if (result.success && result.data) {
          const existingRules: LocalRule[] = result.data.map((config: any) => ({
            _tempId: generateTempId(),
            id: config.id,
            integration_id: config.integration_id,
            notification_type_id: config.notification_type_id,
            notification_event_type_id: config.notification_event_type_id,
            enabled: config.enabled,
            description: config.description || "",
            order_status_ids: config.order_status_ids || [],
            _deleted: false,
          }));
          setRules(existingRules);
          setOriginalIntegrationIds([
            ...new Set(existingRules.map((r) => r.integration_id)),
          ]);
        }
      } catch (error) {
        showToast("Error al cargar reglas existentes", "error");
      } finally {
        setLoadingExisting(false);
      }
    };

    loadExisting();
  }, [businessId]);

  const sections = useMemo<ChannelSection[]>(() => {
    const active = integrations
      .filter((i) => i.is_active && (i.category === "platform" || i.category === "ecommerce"))
      .sort((a, b) => {
        const orderA = CATEGORY_ORDER[a.category ?? ""] ?? 99;
        const orderB = CATEGORY_ORDER[b.category ?? ""] ?? 99;
        if (orderA !== orderB) return orderA - orderB;
        return a.name.localeCompare(b.name);
      })
      .map(sectionFromIntegration);

    const known = new Set(active.map((s) => s.id));
    const orphans = [...new Set(rules.map((r) => r.integration_id))]
      .filter((id) => id > 0 && !known.has(id))
      .map((id) => ({
        id,
        name: `Integraci\u00f3n #${id}`,
        subtitle: "Ya no est\u00e1 activa",
        imageUrl: undefined,
        fallbackLetter: "?",
      }));

    return [...active, ...orphans];
  }, [integrations, rules]);

  const handleAddRule = (integrationId: number) => {
    setRules((prev) => [
      ...prev,
      {
        _tempId: generateTempId(),
        integration_id: integrationId,
        notification_type_id: 0,
        notification_event_type_id: 0,
        enabled: true,
        description: "",
        order_status_ids: [],
        _deleted: false,
      },
    ]);
  };

  const handleRuleChange = (tempId: string, updated: LocalRule) => {
    setRules((prev) => prev.map((r) => (r._tempId === tempId ? updated : r)));
  };

  const handleRuleDelete = (tempId: string) => {
    setRules((prev) =>
      prev
        .map((r) => (r._tempId === tempId && r.id ? { ...r, _deleted: true } : r))
        .filter((r) => !(r._tempId === tempId && !r.id)),
    );
  };

  const handleSave = async () => {
    const activeRules = rules.filter((r) => !r._deleted);

    for (const section of sections) {
      const sectionRules = activeRules.filter((r) => r.integration_id === section.id);

      for (let i = 0; i < sectionRules.length; i++) {
        const rule = sectionRules[i];
        if (!rule.notification_type_id) {
          showToast(`${section.name}, regla ${i + 1}: selecciona un canal`, "error");
          return;
        }
        if (!rule.notification_event_type_id) {
          showToast(`${section.name}, regla ${i + 1}: selecciona un evento`, "error");
          return;
        }
      }

      const seen = new Set<string>();
      for (const rule of sectionRules) {
        const key = `${rule.notification_type_id}-${rule.notification_event_type_id}`;
        if (seen.has(key)) {
          showToast(`${section.name}: hay reglas duplicadas (mismo canal + evento)`, "error");
          return;
        }
        seen.add(key);
      }
    }

    const touched = new Set<number>(originalIntegrationIds);
    for (const rule of activeRules) touched.add(rule.integration_id);

    const targets = sections.filter((s) => touched.has(s.id));

    setLoading(true);
    try {
      let created = 0;
      let updated = 0;
      let deleted = 0;

      for (const section of targets) {
        const dto: SyncConfigsDTO = {
          integration_id: section.id,
          rules: activeRules
            .filter((r) => r.integration_id === section.id)
            .map((r) => ({
              id: r.id,
              notification_type_id: r.notification_type_id,
              notification_event_type_id: r.notification_event_type_id,
              enabled: r.enabled,
              description: r.description,
              order_status_ids: r.order_status_ids,
            })),
        };

        const result = await syncConfigsAction(dto, businessId);

        if (!result.success) {
          showToast(`${section.name}: ${result.error || "error al sincronizar"}`, "error");
          return;
        }

        created += result.data?.created || 0;
        updated += result.data?.updated || 0;
        deleted += result.data?.deleted || 0;
      }

      showToast(
        `Sincronizado: ${created} creadas, ${updated} actualizadas, ${deleted} eliminadas`,
        "success",
      );
      onSuccess();
    } catch (error: any) {
      showToast(error.message || "Error inesperado", "error");
    } finally {
      setLoading(false);
    }
  };

  const activeRulesCount = rules.filter((r) => !r._deleted).length;
  const loadingSections = loadingExisting || loadingIntegrations;

  return (
    <div className="space-y-4">

      <div className="flex gap-1 border-b border-gray-200 dark:border-gray-600">
        <button
          type="button"
          onClick={() => setActiveTab("status")}
          className={`px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === "status"
              ? "border-b-2 border-[var(--color-primary)] text-[var(--color-primary)]"
              : "text-gray-500 hover:text-gray-700 dark:text-gray-400"
          }`}
        >
          {"Por estado de la orden"}
        </button>
        <button
          type="button"
          onClick={() => setActiveTab("scheduled")}
          className={`px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === "scheduled"
              ? "border-b-2 border-[var(--color-primary)] text-[var(--color-primary)]"
              : "text-gray-500 hover:text-gray-700 dark:text-gray-400"
          }`}
        >
          {"Programadas por segmento"}
        </button>
      </div>

      {activeTab === "scheduled" ? (
        <div className="pt-2">
          <p className="mb-4 text-xs text-gray-500 dark:text-gray-400">
            {"No dependen de un evento ni de un canal de venta: salen solas a un grupo de clientes, por ejemplo a los que llevan m\u00e1s de 30 d\u00edas sin comprar."}
          </p>
          {whatsappTypeId > 0 ? (
            <ScheduledRulesSection
              businessId={businessId}
              whatsappTypeId={whatsappTypeId}
            />
          ) : (
            <p className="py-6 text-center text-sm text-gray-500">{"Cargando..."}</p>
          )}
        </div>
      ) : loadingSections ? (
        <div className="text-center py-8 text-gray-500 dark:text-gray-400">Cargando reglas...</div>
      ) : sections.length === 0 ? (
        <div className="py-10 text-center">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {"No hay canales de venta activos"}
          </p>
          <p className="mt-1 text-xs text-gray-400">
            {"Conecta una integraci\u00f3n primero desde el m\u00f3dulo de Integraciones"}
          </p>
        </div>
      ) : (
        <div className="space-y-5">
          <p className="text-xs text-gray-500 dark:text-gray-400">
            {"Cada canal de venta tiene sus propias reglas. Agrega la regla en la secci\u00f3n del canal al que quieres que aplique."}
          </p>

          {sections.map((section) => {
            const sectionRules = rules.filter(
              (r) => r.integration_id === section.id && !r._deleted,
            );

            return (
              <div
                key={section.id}
                className="rounded-lg border border-gray-200 dark:border-gray-600 overflow-hidden"
              >
                <div className="flex items-center gap-3 px-3 py-2 bg-gray-50 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600">
                  {section.imageUrl ? (
                    <img
                      src={section.imageUrl}
                      alt={section.name}
                      className="w-7 h-7 object-contain rounded"
                    />
                  ) : (
                    <div className="w-7 h-7 rounded bg-gray-200 dark:bg-gray-600 flex items-center justify-center">
                      <span className="text-xs font-bold text-gray-500 dark:text-gray-400">
                        {section.fallbackLetter}
                      </span>
                    </div>
                  )}
                  <div className="min-w-0 flex-1">
                    <h3 className="text-sm font-medium text-gray-900 dark:text-white truncate">
                      {section.name}
                    </h3>
                    <p className="text-[11px] text-gray-500 dark:text-gray-400">
                      {section.subtitle}
                    </p>
                  </div>
                  <span className="shrink-0 rounded-full bg-white dark:bg-gray-800 px-2 py-0.5 text-[10px] font-semibold text-gray-500 dark:text-gray-400 border border-gray-200 dark:border-gray-600">
                    {sectionRules.length} {sectionRules.length === 1 ? "regla" : "reglas"}
                  </span>
                </div>

                {sectionRules.length === 0 ? (
                  <div className="py-5 text-center text-xs text-gray-400">
                    {"Sin reglas para este canal"}
                  </div>
                ) : (
                  <table className="w-full">
                    <thead style={{ backgroundColor: "var(--color-primary)" }}>
                      <tr>
                        <th style={{ color: "var(--color-on-primary, white)" }} className="py-2 px-3 text-[10px] font-semibold uppercase text-center w-10">#</th>
                        <th style={{ color: "var(--color-on-primary, white)" }} className="py-2 px-3 text-[10px] font-semibold uppercase text-left w-[120px]">Canal</th>
                        <th style={{ color: "var(--color-on-primary, white)" }} className="py-2 px-3 text-[10px] font-semibold uppercase text-left w-[160px]">Evento</th>
                        <th style={{ color: "var(--color-on-primary, white)" }} className="py-2 px-3 text-[10px] font-semibold uppercase text-left">Estados</th>
                        <th style={{ color: "var(--color-on-primary, white)" }} className="py-2 px-3 text-[10px] font-semibold uppercase text-center w-[70px]">Activo</th>
                        <th style={{ color: "var(--color-on-primary, white)" }} className="py-2 px-3 text-[10px] font-semibold uppercase text-center w-[50px]"></th>
                        <th style={{ color: "var(--color-on-primary, white)" }} className="py-2 px-3 text-[10px] font-semibold uppercase text-center w-[50px]"></th>
                      </tr>
                    </thead>
                    <tbody>
                      {sectionRules.map((rule, index) => (
                        <RuleCard
                          key={rule._tempId}
                          rule={rule}
                          index={index}
                          orderStatuses={orderStatuses}
                          businessId={businessId}
                          onChange={(updated) => handleRuleChange(rule._tempId, updated)}
                          onDelete={() => handleRuleDelete(rule._tempId)}
                        />
                      ))}
                    </tbody>
                  </table>
                )}

                <div className="flex justify-center border-t border-gray-100 dark:border-gray-700 py-2">
                  <button
                    type="button"
                    onClick={() => handleAddRule(section.id)}
                    className="flex items-center gap-1.5 px-3 py-1 rounded-lg bg-[var(--color-primary)]/10 text-[var(--color-primary)] hover:bg-[var(--color-primary)]/20 transition-colors text-xs font-medium"
                    title={`Agregar regla para ${section.name}`}
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                    </svg>
                    {"Agregar regla"}
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      <div className="flex items-center justify-between pt-4 border-t">
        <span className="text-sm text-gray-500 dark:text-gray-400">
          {activeTab === "status" ? `${activeRulesCount} regla(s) activa(s)` : ""}
        </span>
        <div className="flex gap-2">
          <button
            type="button"
            onClick={onCancel}
            disabled={loading}
            className="p-2 rounded-lg bg-gray-100 text-gray-500 dark:text-gray-400 hover:bg-gray-200 transition-colors disabled:opacity-40"
            title="Cancelar"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          {activeTab === "status" && (
            <button
              type="button"
              onClick={handleSave}
              disabled={loading || loadingSections}
              className="p-2 rounded-lg bg-green-50 text-green-600 hover:bg-green-100 transition-colors disabled:opacity-40"
              title={loading ? "Guardando..." : "Guardar reglas"}
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
