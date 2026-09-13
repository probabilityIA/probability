"use client";

import { useCallback, useEffect, useState } from "react";
import { Button } from "@/shared/ui/button";
import { Modal } from "@/shared/ui/modal";
import { useToast } from "@/shared/providers/toast-provider";
import {
  ScheduledRule,
  TEMPLATE_STATUS_LABEL,
  TemplateFlow,
  WhatsappTemplate,
} from "../../domain/scheduled-types";
import {
  getTemplateVariablesAction,
  listAllTemplateFlowsAction,
  listTemplatesAction,
  submitTemplateForReviewAction,
} from "../../infra/actions/whatsapp-templates";
import {
  deleteScheduledRuleAction,
  listScheduledRulesAction,
  runScheduledRuleNowAction,
} from "../../infra/actions/scheduled-rules";
import { TemplateForm } from "./TemplateForm";
import { TemplateFlowView } from "./TemplateFlowView";
import { ScheduledRuleForm } from "./ScheduledRuleForm";

interface ScheduledRulesSectionProps {
  businessId?: number;
  whatsappTypeId: number;
}

type Panel = "none" | "template" | "rule";

const STATUS_STYLE: Record<string, string> = {
  approved: "bg-green-100 text-green-700",
  pending: "bg-amber-100 text-amber-700",
  rejected: "bg-red-100 text-red-700",
  failed: "bg-red-100 text-red-700",
  draft: "bg-gray-100 text-gray-600",
  paused: "bg-orange-100 text-orange-700",
  disabled: "bg-gray-200 text-gray-600",
};

export function ScheduledRulesSection({
  businessId,
  whatsappTypeId,
}: ScheduledRulesSectionProps) {
  const { showToast } = useToast();
  const [loading, setLoading] = useState(true);
  const [panel, setPanel] = useState<Panel>("none");
  const [templates, setTemplates] = useState<WhatsappTemplate[]>([]);
  const [rules, setRules] = useState<ScheduledRule[]>([]);
  const [catalog, setCatalog] = useState<Record<string, string>>({});
  const [flows, setFlows] = useState<TemplateFlow[]>([]);
  const [editingTemplate, setEditingTemplate] = useState<WhatsappTemplate | null>(null);
  const [isFlowOpen, setIsFlowOpen] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);

    const [templatesResult, rulesResult, catalogResult, flowsResult] = await Promise.all([
      listTemplatesAction(businessId, "scheduled"),
      listScheduledRulesAction(businessId),
      getTemplateVariablesAction(businessId),
      listAllTemplateFlowsAction(businessId),
    ]);

    if (templatesResult.success) setTemplates(templatesResult.data);
    if (rulesResult.success) setRules(rulesResult.data);
    if (catalogResult.success) setCatalog(catalogResult.data);
    if (flowsResult.success) setFlows(flowsResult.data);

    setLoading(false);
  }, [businessId]);

  useEffect(() => {
    load();
  }, [load]);

  const handleSubmitTemplate = async (template: WhatsappTemplate) => {
    const result = await submitTemplateForReviewAction(template.ID, businessId);
    if (!result.success) {
      showToast(result.error || "No se pudo enviar la plantilla", "error");
      return;
    }
    showToast("Plantilla enviada a revisión de Meta", "success");
    load();
  };

  const handleRunNow = async (rule: ScheduledRule) => {
    const result = await runScheduledRuleNowAction(rule.ID, businessId);
    if (!result.success) {
      showToast(result.error || "No se pudo ejecutar la regla", "error");
      return;
    }
    const run = result.data;
    showToast(
      `Encolados ${run?.QueuedCount ?? 0} de ${run?.MatchedCount ?? 0} clientes`,
      "success",
    );
    load();
  };

  const handleDelete = async (rule: ScheduledRule) => {
    const result = await deleteScheduledRuleAction(rule.ID, businessId);
    if (!result.success) {
      showToast(result.error || "No se pudo eliminar la regla", "error");
      return;
    }
    showToast("Regla eliminada", "success");
    load();
  };

  if (loading) {
    return <p className="py-6 text-center text-sm text-gray-500">{"Cargando..."}</p>;
  }

  if (panel === "rule") {
    return (
      <ScheduledRuleForm
        businessId={businessId}
        whatsappTypeId={whatsappTypeId}
        templates={templates}
        onSuccess={() => {
          setPanel("none");
          load();
        }}
        onCancel={() => setPanel("none")}
      />
    );
  }

  return (
    <div className="space-y-6">
      <section>
        <div className="mb-2 flex items-center justify-between">
          <h3 className="text-sm font-semibold">{"Plantillas propias"}</h3>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => setIsFlowOpen(true)}
              className="rounded-md border border-gray-300 px-2.5 py-1 text-xs font-medium text-gray-700 transition-colors hover:bg-gray-100 dark:border-gray-600 dark:text-gray-200"
            >
              {"Ver flujo"}
            </button>
            <Button
              type="button"
              size="sm"
              onClick={() => {
                setEditingTemplate(null);
                setPanel("template");
              }}
            >
              {"Nueva plantilla"}
            </Button>
          </div>
        </div>

        {templates.length === 0 ? (
          <p className="rounded-md border border-dashed border-gray-300 p-4 text-sm text-gray-500">
            {"Todavía no creaste ninguna plantilla."}
          </p>
        ) : (
          <ul className="divide-y divide-gray-200 rounded-md border border-gray-200">
            {templates.map((template) => (
              <li key={template.ID} className="flex items-center justify-between p-3">
                <div className="min-w-0">
                  <p className="truncate text-sm font-medium">{template.Name}</p>
                  <p className="truncate text-xs text-gray-500">{template.BodyText}</p>
                  {template.Status === "rejected" && template.RejectedReason && (
                    <p className="mt-1 text-xs text-red-600">{template.RejectedReason}</p>
                  )}
                </div>
                <div className="ml-3 flex shrink-0 items-center gap-2">
                  <span
                    className={`rounded-full px-2 py-1 text-xs ${
                      STATUS_STYLE[template.Status] || "bg-gray-100 text-gray-600"
                    }`}
                  >
                    {TEMPLATE_STATUS_LABEL[template.Status] || template.Status}
                  </span>
                  <button
                    type="button"
                    onClick={() => {
                      setEditingTemplate(template);
                      setPanel("template");
                    }}
                    className="text-xs font-medium text-[var(--color-primary)] hover:underline"
                  >
                    {"Editar"}
                  </button>
                  {["draft", "rejected", "failed"].includes(template.Status) && (
                    <button
                      type="button"
                      onClick={() => handleSubmitTemplate(template)}
                      className="text-xs font-medium text-[var(--color-primary)] hover:underline"
                    >
                      {"Enviar a revisión"}
                    </button>
                  )}
                </div>
              </li>
            ))}
          </ul>
        )}
      </section>

      <section>
        <div className="mb-2 flex items-center justify-between">
          <h3 className="text-sm font-semibold">{"Reglas programadas"}</h3>
          <Button type="button" size="sm" onClick={() => setPanel("rule")}>
            {"Nueva regla"}
          </Button>
        </div>

        {rules.length === 0 ? (
          <p className="rounded-md border border-dashed border-gray-300 p-4 text-sm text-gray-500">
            {
              "Sin reglas programadas. Una regla envía a un segmento de clientes cada cierto tiempo, sin depender de un evento."
            }
          </p>
        ) : (
          <ul className="divide-y divide-gray-200 rounded-md border border-gray-200">
            {rules.map((rule) => (
              <li key={rule.ID} className="flex items-center justify-between p-3">
                <div className="min-w-0">
                  <p className="truncate text-sm font-medium">{rule.Name}</p>
                  <p className="text-xs text-gray-500">
                    {`Más de ${rule.SegmentParams?.DaysWithoutPurchase ?? 0} días sin comprar · ${rule.SendWindowStart}-${rule.SendWindowEnd} · máx. ${rule.DailySendCap}/día`}
                  </p>
                </div>
                <div className="ml-3 flex shrink-0 items-center gap-2">
                  <span
                    className={`rounded-full px-2 py-1 text-xs ${
                      rule.Enabled
                        ? "bg-green-100 text-green-700"
                        : "bg-gray-100 text-gray-600"
                    }`}
                  >
                    {rule.Enabled ? "Activa" : "Pausada"}
                  </span>
                  <button
                    type="button"
                    onClick={() => handleRunNow(rule)}
                    className="text-xs font-medium text-[var(--color-primary)] hover:underline"
                  >
                    {"Ejecutar ahora"}
                  </button>
                  <button
                    type="button"
                    onClick={() => handleDelete(rule)}
                    className="text-xs text-red-500 hover:underline"
                  >
                    {"Eliminar"}
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </section>

      <Modal
        isOpen={isFlowOpen}
        onClose={() => setIsFlowOpen(false)}
        title={(
          <span className="flex w-full flex-col items-start pr-8">
            <span className="text-lg font-semibold">{"Flujo de plantillas"}</span>
            <span className="text-[13px] font-normal text-gray-400">
              {"Qué responde cada botón"}
            </span>
          </span>
        )}
        size="4xl"
        zIndex={60}
      >
        {isFlowOpen && <TemplateFlowView templates={templates} flows={flows} />}
      </Modal>

      <Modal
        isOpen={panel === "template"}
        onClose={() => {
          setPanel("none");
          setEditingTemplate(null);
        }}
        title={(
          <span className="flex w-full flex-col items-start pr-8">
            <span className="text-lg font-semibold">
              {editingTemplate ? "Editar plantilla" : "Nueva plantilla"}
            </span>
            <span className="text-[13px] font-normal text-gray-400">
              {"Plantilla de mensaje para WhatsApp · Meta"}
            </span>
          </span>
        )}
        size="4xl"
        zIndex={60}
        noPadding
        noBodyScroll
      >
        {panel === "template" && (
          <TemplateForm
            businessId={businessId}
            variableCatalog={catalog}
            template={editingTemplate}
            templates={templates}
            flows={flows.filter((flow) => flow.SourceTemplateID === editingTemplate?.ID)}
            onSuccess={() => {
              setPanel("none");
              setEditingTemplate(null);
              load();
            }}
            onCancel={() => {
              setPanel("none");
              setEditingTemplate(null);
            }}
          />
        )}
      </Modal>
    </div>
  );
}
