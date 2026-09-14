"use client";

import { useCallback, useEffect, useState } from "react";
import { createPortal } from "react-dom";
import { Modal } from "@/shared/ui/modal";
import { useToast } from "@/shared/providers/toast-provider";
import { Flow, TemplateFlow, WhatsappTemplate } from "../../domain/scheduled-types";
import {
  createFlowAction,
  deleteFlowAction,
  getTemplateVariablesAction,
  listAllTemplateFlowsAction,
  listFlowsAction,
  listTemplatesAction,
  updateFlowAction,
} from "../../infra/actions/whatsapp-templates";
import { TemplateFlowView } from "./TemplateFlowView";

export const NOTIFICATIONS_ACTIONS_SLOT_ID = "notifications-actions-slot";

interface FlowsSectionProps {
  businessId?: number;
}

const inputCls =
  "w-full rounded-lg border border-gray-300 px-3 py-2 text-sm text-gray-900 outline-none focus:border-[var(--color-primary)] dark:border-gray-600 dark:bg-gray-800 dark:text-white";

export function FlowsSection({ businessId }: FlowsSectionProps) {
  const { showToast } = useToast();

  const [flows, setFlows] = useState<Flow[]>([]);
  const [transitions, setTransitions] = useState<TemplateFlow[]>([]);
  const [templates, setTemplates] = useState<WhatsappTemplate[]>([]);
  const [catalog, setCatalog] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);

  const [editing, setEditing] = useState<Flow | null>(null);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [rootTemplateID, setRootTemplateID] = useState(0);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [confirmDelete, setConfirmDelete] = useState(0);
  const [studioFlow, setStudioFlow] = useState<Flow | null>(null);
  const [actionsSlot, setActionsSlot] = useState<HTMLElement | null>(null);
  const [expanded, setExpanded] = useState<Record<number, boolean>>({});

  const toggleExpanded = (id: number) =>
    setExpanded((current) => ({ ...current, [id]: !current[id] }));

  useEffect(() => {
    setActionsSlot(document.getElementById(NOTIFICATIONS_ACTIONS_SLOT_ID));
  }, []);

  const load = useCallback(async () => {
    setLoading(true);

    const [flowsResult, transitionsResult, templatesResult, catalogResult] = await Promise.all([
      listFlowsAction(businessId),
      listAllTemplateFlowsAction(businessId),
      listTemplatesAction(businessId, "all", undefined, 1, 100),
      getTemplateVariablesAction(businessId),
    ]);

    if (flowsResult.success) setFlows(flowsResult.data);
    if (transitionsResult.success) setTransitions(transitionsResult.data);
    if (templatesResult.success) setTemplates(templatesResult.data);
    if (catalogResult.success) setCatalog(catalogResult.data);

    setConfirmDelete(0);
    setLoading(false);
  }, [businessId]);

  useEffect(() => {
    load();
  }, [load]);

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

  const flowActions = (flow: Flow) =>
    confirmDelete === flow.ID ? (
      <>
        <span className="text-[12px] text-red-600">{"¿Eliminar el flujo?"}</span>
        <button
          type="button"
          onClick={() => handleDelete(flow)}
          className="text-xs font-medium text-red-600 hover:underline"
        >
          {"Sí"}
        </button>
        <button
          type="button"
          onClick={() => setConfirmDelete(0)}
          className="text-xs font-medium text-gray-500 hover:underline"
        >
          {"No"}
        </button>
      </>
    ) : (
      <>
        <button
          type="button"
          onClick={() => setStudioFlow(flow)}
          className="text-xs font-medium text-[var(--color-primary)] hover:underline"
        >
          {"Editar flujo"}
        </button>
        <button
          type="button"
          onClick={() => startEdit(flow)}
          className="text-xs font-medium text-gray-600 hover:underline dark:text-gray-300"
        >
          {"Ajustes"}
        </button>
        <button
          type="button"
          onClick={() => setConfirmDelete(flow.ID)}
          className="text-xs font-medium text-red-500 hover:underline"
        >
          {"Eliminar"}
        </button>
      </>
    );

  return (
    <div className="space-y-5">
      {actionsSlot &&
        createPortal(
          <button
            type="button"
            onClick={startCreate}
            style={{
              backgroundColor: "var(--color-primary)",
              color: "var(--color-on-primary, white)",
            }}
            className="rounded-lg px-3 py-1.5 text-sm font-medium transition-opacity hover:opacity-90"
            data-tour="notif-new-flow"
          >
            {"Nuevo flujo"}
          </button>,
          actionsSlot,
        )}

      {flows.length === 0 ? (
        <p className="rounded-md border border-dashed border-gray-300 p-6 text-center text-sm text-gray-500 dark:border-gray-600">
          {
            "Todavía no hay flujos. Un flujo agrupa una plantilla inicial y las respuestas que se encadenan a sus botones."
          }
        </p>
      ) : (
        flows.map((flow) =>
          flow.RootTemplateID ? (
            <TemplateFlowView
              key={flow.ID}
              templates={templates}
              flows={transitions.filter((item) => item.FlowID === flow.ID)}
              businessId={businessId}
              variableCatalog={catalog}
              flowId={flow.ID}
              rootTemplateId={flow.RootTemplateID}
              title={flow.Name}
              subtitle={flow.Description}
              actions={flowActions(flow)}
              collapsed={!expanded[flow.ID]}
              onToggleCollapsed={() => toggleExpanded(flow.ID)}
              readOnly
              onChanged={load}
            />
          ) : (
            <div
              key={flow.ID}
              className="rounded-xl border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800"
            >
              <div className="flex flex-wrap items-center gap-3">
                <button
                  type="button"
                  onClick={() => toggleExpanded(flow.ID)}
                  aria-label={expanded[flow.ID] ? "Contraer" : "Expandir"}
                  className="flex h-6 w-6 shrink-0 items-center justify-center rounded text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
                >
                  <svg
                    className={`h-4 w-4 transition-transform ${
                      expanded[flow.ID] ? "rotate-90" : ""
                    }`}
                    fill="none"
                    stroke="currentColor"
                    strokeWidth={2}
                    viewBox="0 0 24 24"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
                  </svg>
                </button>
                <div className="min-w-0">
                  <p className="truncate text-sm font-semibold text-gray-900 dark:text-white">
                    {flow.Name}
                  </p>
                  {flow.Description && (
                    <p className="truncate text-xs text-gray-500 dark:text-gray-400">
                      {flow.Description}
                    </p>
                  )}
                </div>
                <div className="ml-auto flex items-center gap-3">{flowActions(flow)}</div>
              </div>

              {expanded[flow.ID] && (
                <>
                  <p className="mt-3 text-sm text-gray-500">
                    {"Este flujo no tiene plantilla inicial: es el mensaje con el que arranca."}
                  </p>
                  <button
                    type="button"
                    onClick={() => startEdit(flow)}
                    className="mt-1 text-sm font-medium text-[var(--color-primary)] hover:underline"
                  >
                    {"Elegir plantilla inicial"}
                  </button>
                </>
              )}
            </div>
          ),
        )
      )}

      <Modal
        isOpen={studioFlow !== null}
        onClose={() => setStudioFlow(null)}
        title={(
          <span className="flex w-full flex-col items-start pr-8">
            <span className="text-lg font-semibold">{studioFlow?.Name ?? ""}</span>
            <span className="text-[13px] font-normal text-gray-400">
              {"Arrastrá para moverte, rueda para acercar"}
            </span>
          </span>
        )}
        size="full"
        zIndex={60}
      >
        {studioFlow !== null && studioFlow.RootTemplateID && (
          <div className="flex h-full min-h-0 flex-col p-6">
            <TemplateFlowView
              templates={templates}
              flows={transitions.filter((item) => item.FlowID === studioFlow.ID)}
              businessId={businessId}
              variableCatalog={catalog}
              flowId={studioFlow.ID}
              rootTemplateId={studioFlow.RootTemplateID}
              studio
              onChanged={load}
            />
          </div>
        )}
      </Modal>

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
    </div>
  );
}
