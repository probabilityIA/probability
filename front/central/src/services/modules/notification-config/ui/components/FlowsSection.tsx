"use client";

import { useCallback, useEffect, useState } from "react";
import { Modal } from "@/shared/ui/modal";
import { useToast } from "@/shared/providers/toast-provider";
import {
  Flow,
  TEMPLATE_STATUS_LABEL,
  TemplateFlow,
  WhatsappTemplate,
} from "../../domain/scheduled-types";
import {
  createFlowAction,
  deleteFlowAction,
  getTemplateVariablesAction,
  listFlowTransitionsAction,
  listFlowsAction,
  listTemplatesAction,
  updateFlowAction,
} from "../../infra/actions/whatsapp-templates";
import { TemplateBubble, fillPlaceholders } from "./TemplateBubble";
import { TemplateFlowView } from "./TemplateFlowView";

interface FlowsSectionProps {
  businessId?: number;
}

const OPT_OUT_TEXT = "Dejar de recibir";

const STATUS_STYLE: Record<string, string> = {
  approved: "bg-green-100 text-green-700",
  pending: "bg-amber-100 text-amber-700",
  rejected: "bg-red-100 text-red-700",
  failed: "bg-red-100 text-red-700",
  draft: "bg-gray-100 text-gray-600",
  paused: "bg-orange-100 text-orange-700",
  disabled: "bg-gray-200 text-gray-600",
};

const inputCls =
  "w-full rounded-lg border border-gray-300 px-3 py-2 text-sm text-gray-900 outline-none focus:border-[var(--color-primary)] dark:border-gray-600 dark:bg-gray-800 dark:text-white";

function templateButtons(template: WhatsappTemplate): string[] {
  const own = (template.Buttons ?? [])
    .map((button) => (button.Text ?? "").trim())
    .filter((text) => text && text.toLowerCase() !== OPT_OUT_TEXT.toLowerCase());

  return template.Category === "MARKETING" ? [...own, OPT_OUT_TEXT] : own;
}

export function FlowsSection({ businessId }: FlowsSectionProps) {
  const { showToast } = useToast();

  const [flows, setFlows] = useState<Flow[]>([]);
  const [templates, setTemplates] = useState<WhatsappTemplate[]>([]);
  const [catalog, setCatalog] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);

  const [openFlow, setOpenFlow] = useState<Flow | null>(null);
  const [transitions, setTransitions] = useState<TemplateFlow[]>([]);

  const [editing, setEditing] = useState<Flow | null>(null);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [rootTemplateID, setRootTemplateID] = useState(0);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  const load = useCallback(async () => {
    setLoading(true);

    const [flowsResult, templatesResult, catalogResult] = await Promise.all([
      listFlowsAction(businessId),
      listTemplatesAction(businessId, "all", undefined, 1, 100),
      getTemplateVariablesAction(businessId),
    ]);

    if (flowsResult.success) setFlows(flowsResult.data);
    if (templatesResult.success) setTemplates(templatesResult.data);
    if (catalogResult.success) setCatalog(catalogResult.data);

    setConfirmDelete(null);
    setLoading(false);
  }, [businessId]);

  useEffect(() => {
    load();
  }, [load]);

  const loadTransitions = useCallback(
    async (flowId: number) => {
      const result = await listFlowTransitionsAction(flowId, businessId);
      if (result.success) setTransitions(result.data);
    },
    [businessId],
  );

  const openDiagram = async (flow: Flow) => {
    setOpenFlow(flow);
    await loadTransitions(flow.ID);
  };

  const startCreate = () => {
    setEditing(null);
    setName("");
    setDescription("");
    setRootTemplateID(0);
    setError("");
    setIsFormOpen(true);
  };

  const startEdit = (flow: Flow) => {
    setEditing(flow);
    setName(flow.Name);
    setDescription(flow.Description);
    setRootTemplateID(flow.RootTemplateID ?? 0);
    setError("");
    setIsFormOpen(true);
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!name.trim()) {
      setError("El nombre del flujo es obligatorio");
      return;
    }

    const input = {
      name: name.trim(),
      description: description.trim(),
      root_template_id: rootTemplateID > 0 ? rootTemplateID : null,
    };

    setSaving(true);
    const result = editing
      ? await updateFlowAction(editing.ID, input, businessId)
      : await createFlowAction(input, businessId);
    setSaving(false);

    if (!result.success) {
      setError(result.error || "No se pudo guardar el flujo");
      return;
    }

    setIsFormOpen(false);
    showToast(editing ? "Flujo actualizado" : "Flujo creado", "success");
    load();
  };

  const handleDelete = async (flow: Flow) => {
    const result = await deleteFlowAction(flow.ID, businessId);
    if (!result.success) {
      showToast(result.error || "No se pudo eliminar el flujo", "error");
      return;
    }
    showToast("Flujo eliminado", "success");
    load();
  };

  if (loading) {
    return <p className="py-10 text-center text-sm text-gray-500">{"Cargando..."}</p>;
  }

  const rootOf = (flow: Flow) =>
    flow.RootTemplateID ? templates.find((item) => item.ID === flow.RootTemplateID) : undefined;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-baseline gap-3">
          <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
            {"Flujos de conversación"}
          </h3>
          <span className="text-xs text-gray-400">{`${flows.length}`}</span>
        </div>
        <button
          type="button"
          onClick={startCreate}
          style={{
            backgroundColor: "var(--color-primary)",
            color: "var(--color-on-primary, white)",
          }}
          className="rounded-lg px-3 py-1.5 text-sm font-medium transition-opacity hover:opacity-90"
        >
          {"Nuevo flujo"}
        </button>
      </div>

      {flows.length === 0 ? (
        <p className="rounded-md border border-dashed border-gray-300 p-6 text-center text-sm text-gray-500 dark:border-gray-600">
          {
            "Todavía no hay flujos. Un flujo agrupa una plantilla inicial y las respuestas que se encadenan a sus botones."
          }
        </p>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {flows.map((flow) => {
            const root = rootOf(flow);

            return (
              <div
                key={flow.ID}
                className="flex flex-col gap-2.5 rounded-xl border border-gray-200 bg-white p-3 dark:border-gray-700 dark:bg-gray-800"
              >
                <div className="flex items-start justify-between gap-2">
                  <p className="min-w-0 truncate text-[13px] font-semibold text-gray-900 dark:text-white">
                    {flow.Name}
                  </p>
                  {flow.RootTemplateStatus && (
                    <span
                      className={`shrink-0 rounded-full px-2 py-0.5 text-[10px] ${
                        STATUS_STYLE[flow.RootTemplateStatus] || "bg-gray-100 text-gray-600"
                      }`}
                    >
                      {TEMPLATE_STATUS_LABEL[flow.RootTemplateStatus] || flow.RootTemplateStatus}
                    </span>
                  )}
                </div>

                {flow.Description && (
                  <p className="line-clamp-2 text-[11px] text-gray-500 dark:text-gray-400">
                    {flow.Description}
                  </p>
                )}

                {root ? (
                  <div className="rounded-lg bg-[#e9e2d9] p-2.5 dark:bg-[#2a2724]">
                    <TemplateBubble
                      headerType={root.HeaderType}
                      headerMediaURL={root.HeaderMediaURL}
                      headerText={root.HeaderText}
                      bodyText={fillPlaceholders(root.BodyText, root.Variables)}
                      footerText={root.FooterText}
                      buttons={templateButtons(root)}
                      className="w-full"
                      compact
                    />
                  </div>
                ) : (
                  <p className="rounded-lg border border-dashed border-gray-300 p-3 text-center text-[11px] text-gray-400 dark:border-gray-600">
                    {"Sin plantilla inicial"}
                  </p>
                )}

                <div className="flex flex-wrap gap-1.5">
                  <span className="rounded-full bg-gray-100 px-2 py-0.5 text-[10px] text-gray-600 dark:bg-gray-700 dark:text-gray-300">
                    {`${flow.StepCount} paso(s)`}
                  </span>
                  {flow.PendingCount > 0 && (
                    <span className="rounded-full bg-amber-100 px-2 py-0.5 text-[10px] text-amber-700">
                      {`${flow.PendingCount} sin aprobar`}
                    </span>
                  )}
                </div>

                <div className="mt-auto flex flex-wrap items-center justify-end gap-3 border-t border-gray-100 pt-2 dark:border-gray-700">
                  {confirmDelete === flow.ID ? (
                    <>
                      <span className="mr-auto text-[11px] text-red-600">{"¿Eliminar?"}</span>
                      <button
                        type="button"
                        onClick={() => handleDelete(flow)}
                        className="text-[11px] font-medium text-red-600 hover:underline"
                      >
                        {"Sí"}
                      </button>
                      <button
                        type="button"
                        onClick={() => setConfirmDelete(null)}
                        className="text-[11px] font-medium text-gray-500 hover:underline"
                      >
                        {"No"}
                      </button>
                    </>
                  ) : (
                    <>
                      <button
                        type="button"
                        onClick={() => openDiagram(flow)}
                        className="text-[11px] font-medium text-[var(--color-primary)] hover:underline"
                      >
                        {"Abrir diagrama"}
                      </button>
                      <button
                        type="button"
                        onClick={() => startEdit(flow)}
                        className="text-[11px] font-medium text-[var(--color-primary)] hover:underline"
                      >
                        {"Editar"}
                      </button>
                      <button
                        type="button"
                        onClick={() => setConfirmDelete(flow.ID)}
                        className="text-[11px] font-medium text-red-500 hover:underline"
                      >
                        {"Eliminar"}
                      </button>
                    </>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}

      <Modal
        isOpen={isFormOpen}
        onClose={() => setIsFormOpen(false)}
        title={(
          <span className="flex w-full flex-col items-start pr-8">
            <span className="text-lg font-semibold">
              {editing ? "Editar flujo" : "Nuevo flujo"}
            </span>
            <span className="text-[13px] font-normal text-gray-400">
              {"Una plantilla inicial y las respuestas encadenadas a sus botones"}
            </span>
          </span>
        )}
        size="xl"
        zIndex={60}
      >
        <form onSubmit={handleSave} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <label className="text-[13px] font-semibold text-gray-900 dark:text-white">
              {"Nombre"}
            </label>
            <input
              value={name}
              maxLength={120}
              onChange={(e) => setName(e.target.value)}
              placeholder="Reactivación de clientes"
              className={inputCls}
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-[13px] font-semibold text-gray-900 dark:text-white">
              {"Descripción "}
              <span className="font-normal text-gray-400">{"(opcional)"}</span>
            </label>
            <input
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Para qué sirve este flujo"
              className={inputCls}
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-[13px] font-semibold text-gray-900 dark:text-white">
              {"Plantilla inicial"}
            </label>
            <select
              value={rootTemplateID}
              onChange={(e) => setRootTemplateID(Number(e.target.value))}
              style={{ height: 38, fontSize: 14, padding: "0 10px", borderRadius: 8 }}
              className="w-full border border-gray-300 bg-white text-gray-900 outline-none dark:border-gray-600 dark:bg-gray-800 dark:text-white"
            >
              <option value={0}>{"Sin definir"}</option>
              {templates
                .filter((item) => item.Origin === "business")
                .map((item) => (
                  <option key={item.ID} value={item.ID}>
                    {item.Name}
                  </option>
                ))}
            </select>
            <span className="text-[12px] text-gray-400">
              {"Es el mensaje con el que arranca la conversación."}
            </span>
          </div>

          {error && (
            <span className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-[13px] font-medium text-red-700 dark:border-red-900 dark:bg-red-950/40 dark:text-red-300">
              {error}
            </span>
          )}

          <div className="flex justify-end gap-2.5">
            <button
              type="button"
              onClick={() => setIsFormOpen(false)}
              disabled={saving}
              className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-900 transition-colors hover:bg-gray-100 disabled:opacity-40 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            >
              {"Cancelar"}
            </button>
            <button
              type="submit"
              disabled={saving}
              style={{
                backgroundColor: "var(--color-primary)",
                color: "var(--color-on-primary, white)",
              }}
              className="rounded-lg px-4 py-2 text-sm font-semibold transition-opacity hover:opacity-90 disabled:opacity-40"
            >
              {saving ? "Guardando..." : "Guardar"}
            </button>
          </div>
        </form>
      </Modal>

      <Modal
        isOpen={openFlow !== null}
        onClose={() => setOpenFlow(null)}
        title={(
          <span className="flex w-full flex-col items-start pr-8">
            <span className="text-lg font-semibold">{openFlow?.Name ?? ""}</span>
            <span className="text-[13px] font-normal text-gray-400">
              {"Qué responde cada botón"}
            </span>
          </span>
        )}
        size="6xl"
        zIndex={60}
      >
        {openFlow !== null && (
          <TemplateFlowView
            templates={templates}
            flows={transitions}
            businessId={businessId}
            variableCatalog={catalog}
            flowId={openFlow.ID}
            rootTemplateId={openFlow.RootTemplateID ?? undefined}
            onChanged={() => {
              loadTransitions(openFlow.ID);
              load();
            }}
          />
        )}
      </Modal>
    </div>
  );
}
