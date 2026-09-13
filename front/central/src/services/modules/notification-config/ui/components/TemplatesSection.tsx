"use client";

import { useCallback, useEffect, useState } from "react";
import { Modal } from "@/shared/ui/modal";
import { useToast } from "@/shared/providers/toast-provider";
import { TEMPLATE_STATUS_LABEL, WhatsappTemplate } from "../../domain/scheduled-types";
import {
  deleteTemplateAction,
  getTemplateVariablesAction,
  listTemplatesAction,
  submitTemplateForReviewAction,
} from "../../infra/actions/whatsapp-templates";
import { TemplateBubble, fillPlaceholders } from "./TemplateBubble";
import { TemplateForm } from "./TemplateForm";

interface TemplatesSectionProps {
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

const SCOPE_LABEL: Record<string, string> = {
  order_event: "Evento de pedido",
  scheduled: "Programada",
  campaign: "Campaña",
  internal: "Interna",
};

function templateButtons(template: WhatsappTemplate): string[] {
  const own = (template.Buttons ?? [])
    .map((button) => (button.Text ?? "").trim())
    .filter((text) => text && text.toLowerCase() !== OPT_OUT_TEXT.toLowerCase());

  return template.Category === "MARKETING" ? [...own, OPT_OUT_TEXT] : own;
}

export function TemplatesSection({ businessId }: TemplatesSectionProps) {
  const { showToast } = useToast();

  const [templates, setTemplates] = useState<WhatsappTemplate[]>([]);
  const [catalog, setCatalog] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState<WhatsappTemplate | null>(null);
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null);

  const load = useCallback(async () => {
    setLoading(true);

    const [templatesResult, catalogResult] = await Promise.all([
      listTemplatesAction(businessId, "all", undefined, 1, 100),
      getTemplateVariablesAction(businessId),
    ]);

    if (templatesResult.success) setTemplates(templatesResult.data);
    if (catalogResult.success) setCatalog(catalogResult.data);

    setConfirmDelete(null);
    setLoading(false);
  }, [businessId]);

  useEffect(() => {
    load();
  }, [load]);

  const handleSubmitForReview = async (template: WhatsappTemplate) => {
    const result = await submitTemplateForReviewAction(template.ID, businessId);
    if (!result.success) {
      showToast(result.error || "No se pudo enviar la plantilla", "error");
      return;
    }
    showToast("Plantilla enviada a revisión de Meta", "success");
    load();
  };

  const handleDelete = async (template: WhatsappTemplate) => {
    const result = await deleteTemplateAction(template.ID, businessId);
    if (!result.success) {
      showToast(result.error || "No se pudo eliminar la plantilla", "error");
      return;
    }
    showToast("Plantilla eliminada", "success");
    load();
  };

  if (loading) {
    return <p className="py-10 text-center text-sm text-gray-500">{"Cargando..."}</p>;
  }

  const own = templates.filter((item) => item.Origin === "business");
  const system = templates.filter((item) => item.Origin !== "business");

  const renderCard = (template: WhatsappTemplate, editable: boolean) => (
    <div
      key={template.ID}
      className="flex flex-col gap-2.5 rounded-xl border border-gray-200 bg-white p-3 dark:border-gray-700 dark:bg-gray-800"
    >
      <div className="flex items-start justify-between gap-2">
        <p className="min-w-0 truncate text-[13px] font-semibold text-gray-900 dark:text-white">
          {template.Name}
        </p>
        <span
          className={`shrink-0 rounded-full px-2 py-0.5 text-[10px] ${
            STATUS_STYLE[template.Status] || "bg-gray-100 text-gray-600"
          }`}
        >
          {TEMPLATE_STATUS_LABEL[template.Status] || template.Status}
        </span>
      </div>

      <div className="flex flex-wrap gap-1.5">
        <span className="rounded-full bg-gray-100 px-2 py-0.5 text-[10px] text-gray-600 dark:bg-gray-700 dark:text-gray-300">
          {SCOPE_LABEL[template.Scope] || template.Scope}
        </span>
        <span className="rounded-full bg-gray-100 px-2 py-0.5 text-[10px] text-gray-600 dark:bg-gray-700 dark:text-gray-300">
          {template.Category === "MARKETING" ? "Marketing" : "Utilidad"}
        </span>
      </div>

      <div className="rounded-lg bg-[#e9e2d9] p-2.5 dark:bg-[#2a2724]">
        <TemplateBubble
          headerType={template.HeaderType}
          headerMediaURL={template.HeaderMediaURL}
          headerText={template.HeaderText}
          bodyText={fillPlaceholders(template.BodyText, template.Variables)}
          footerText={template.FooterText}
          buttons={templateButtons(template)}
          className="w-full"
          compact
        />
      </div>

      {template.Status === "rejected" && template.RejectedReason && (
        <p className="text-[11px] text-red-600">{template.RejectedReason}</p>
      )}

      <div className="mt-auto flex flex-wrap items-center justify-end gap-3 border-t border-gray-100 pt-2 dark:border-gray-700">
        {!editable && (
          <span className="mr-auto text-[11px] text-gray-400">{"Del sistema"}</span>
        )}

        {editable && confirmDelete === template.ID && (
          <>
            <span className="mr-auto text-[11px] text-red-600">
              {"¿Eliminar? También se borra en Meta"}
            </span>
            <button
              type="button"
              onClick={() => handleDelete(template)}
              className="text-[11px] font-medium text-red-600 hover:underline"
            >
              {"Sí, eliminar"}
            </button>
            <button
              type="button"
              onClick={() => setConfirmDelete(null)}
              className="text-[11px] font-medium text-gray-500 hover:underline"
            >
              {"Cancelar"}
            </button>
          </>
        )}

        {editable && confirmDelete !== template.ID && (
          <>
            <button
              type="button"
              onClick={() => setEditing(template)}
              className="text-[11px] font-medium text-[var(--color-primary)] hover:underline"
            >
              {"Editar"}
            </button>
            {["draft", "rejected", "failed"].includes(template.Status) && (
              <button
                type="button"
                onClick={() => handleSubmitForReview(template)}
                className="text-[11px] font-medium text-[var(--color-primary)] hover:underline"
              >
                {"Enviar a revisión"}
              </button>
            )}
            <button
              type="button"
              onClick={() => setConfirmDelete(template.ID)}
              className="text-[11px] font-medium text-red-500 hover:underline"
            >
              {"Eliminar"}
            </button>
          </>
        )}
      </div>
    </div>
  );

  return (
    <div className="space-y-8">
      <section>
        <div className="mb-3 flex items-baseline gap-3">
          <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
            {"Propias del negocio"}
          </h3>
          <span className="text-xs text-gray-400">{`${own.length}`}</span>
        </div>

        {own.length === 0 ? (
          <p className="rounded-md border border-dashed border-gray-300 p-4 text-sm text-gray-500 dark:border-gray-600">
            {"Todavía no creaste plantillas propias."}
          </p>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {own.map((template) => renderCard(template, true))}
          </div>
        )}
      </section>

      <section>
        <div className="mb-3 flex items-baseline gap-3">
          <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
            {"Predeterminadas del sistema"}
          </h3>
          <span className="text-xs text-gray-400">{`${system.length}`}</span>
        </div>

        {system.length === 0 ? (
          <p className="rounded-md border border-dashed border-gray-300 p-4 text-sm text-gray-500 dark:border-gray-600">
            {"No hay plantillas del sistema disponibles."}
          </p>
        ) : (
          <>
            <p className="mb-3 text-xs text-gray-400">
              {"Vienen con Probability y las comparten todos los negocios: no se editan ni se eliminan."}
            </p>
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
              {system.map((template) => renderCard(template, false))}
            </div>
          </>
        )}
      </section>

      <Modal
        isOpen={editing !== null}
        onClose={() => setEditing(null)}
        title={(
          <span className="flex w-full flex-col items-start pr-8">
            <span className="text-lg font-semibold">{"Editar plantilla"}</span>
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
        {editing !== null && (
          <TemplateForm
            businessId={businessId}
            variableCatalog={catalog}
            template={editing}
            onSuccess={() => {
              setEditing(null);
              load();
            }}
            onCancel={() => setEditing(null)}
          />
        )}
      </Modal>
    </div>
  );
}
