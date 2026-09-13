"use client";

import { useState } from "react";
import { Modal } from "@/shared/ui/modal";
import { useToast } from "@/shared/providers/toast-provider";
import {
  TEMPLATE_STATUS_LABEL,
  TemplateFlow,
  WhatsappTemplate,
} from "../../domain/scheduled-types";
import { replaceTemplateFlowsAction } from "../../infra/actions/whatsapp-templates";
import { TemplateBubble, fillPlaceholders } from "./TemplateBubble";
import { TemplateForm } from "./TemplateForm";

interface TemplateFlowViewProps {
  templates: WhatsappTemplate[];
  flows: TemplateFlow[];
  businessId?: number;
  variableCatalog: Record<string, string>;
  onChanged: () => void;
}

interface PendingResponse {
  sourceID: number;
  sourceName: string;
  buttonText: string;
}

const OPT_OUT_TEXT = "Dejar de recibir";
const MAX_DEPTH = 4;

const STATUS_STYLE: Record<string, string> = {
  approved: "bg-green-100 text-green-700",
  pending: "bg-amber-100 text-amber-700",
  rejected: "bg-red-100 text-red-700",
  failed: "bg-red-100 text-red-700",
  draft: "bg-gray-100 text-gray-600",
  paused: "bg-orange-100 text-orange-700",
  disabled: "bg-gray-200 text-gray-600",
};

function userButtons(template: WhatsappTemplate): string[] {
  if (!template.Buttons) return [];
  return template.Buttons.map((button) => (button.Text ?? "").trim()).filter(
    (text) => text && text.toLowerCase() !== OPT_OUT_TEXT.toLowerCase(),
  );
}

export function TemplateFlowView({
  templates,
  flows,
  businessId,
  variableCatalog,
  onChanged,
}: TemplateFlowViewProps) {
  const { showToast } = useToast();
  const [pending, setPending] = useState<PendingResponse | null>(null);
  const [saving, setSaving] = useState(false);

  const byID = new Map(templates.map((item) => [item.ID, item]));

  const childrenBySource = new Map<number, TemplateFlow[]>();
  const parentCount = new Map<number, number>();

  for (const flow of flows) {
    const current = childrenBySource.get(flow.SourceTemplateID) || [];
    current.push(flow);
    childrenBySource.set(flow.SourceTemplateID, current);
    parentCount.set(flow.TargetTemplateID, (parentCount.get(flow.TargetTemplateID) || 0) + 1);
  }

  const participates = templates.filter(
    (item) =>
      userButtons(item).length > 0 || childrenBySource.has(item.ID) || parentCount.has(item.ID),
  );

  const roots = participates.filter((item) => !parentCount.has(item.ID));

  const linkResponse = async (created: WhatsappTemplate) => {
    if (!pending) return;

    const existing = flows
      .filter((flow) => flow.SourceTemplateID === pending.sourceID)
      .map((flow) => ({
        button_text: flow.ButtonText,
        target_template_id: flow.TargetTemplateID,
        enabled: flow.Enabled,
      }));

    const next = [
      ...existing.filter(
        (item) => item.button_text.toLowerCase() !== pending.buttonText.toLowerCase(),
      ),
      { button_text: pending.buttonText, target_template_id: created.ID },
    ];

    setSaving(true);
    const result = await replaceTemplateFlowsAction(pending.sourceID, next, businessId);
    setSaving(false);
    setPending(null);

    if (!result.success) {
      showToast(result.error || "La plantilla se creó pero no quedó enlazada", "error");
      return;
    }

    showToast(`"${pending.buttonText}" ahora responde con ${created.Name}`, "success");
    onChanged();
  };

  if (participates.length === 0) {
    return (
      <p className="rounded-md border border-dashed border-gray-300 p-6 text-center text-sm text-gray-500 dark:border-gray-600">
        {
          "Ninguna plantilla tiene botones todavía. Agregá botones de respuesta para armar un flujo."
        }
      </p>
    );
  }

  const renderNode = (
    template: WhatsappTemplate,
    depth: number,
    path: number[],
  ): React.ReactNode => {
    const buttons = userButtons(template);
    const children = childrenBySource.get(template.ID) || [];
    const flowByButton = new Map(
      children.map((flow) => [flow.ButtonText.trim().toLowerCase(), flow]),
    );
    const reused = (parentCount.get(template.ID) || 0) > 1;

    return (
      <div className="flex items-center gap-3">
        <div className="flex w-[230px] shrink-0 flex-col gap-1.5">
          <div className="flex flex-wrap items-center gap-1.5">
            <span className="truncate text-[12px] font-semibold text-gray-900 dark:text-white">
              {template.Name}
            </span>
            <span
              className={`shrink-0 rounded-full px-1.5 py-0.5 text-[10px] ${
                STATUS_STYLE[template.Status] || "bg-gray-100 text-gray-600"
              }`}
            >
              {TEMPLATE_STATUS_LABEL[template.Status] || template.Status}
            </span>
            {reused && (
              <span className="shrink-0 rounded-full bg-blue-50 px-1.5 py-0.5 text-[10px] text-blue-600">
                {"Reutilizada"}
              </span>
            )}
          </div>

          <TemplateBubble
            headerType={template.HeaderType}
            headerMediaURL={template.HeaderMediaURL}
            headerText={template.HeaderText}
            bodyText={fillPlaceholders(template.BodyText, template.Variables)}
            footerText={template.FooterText}
            buttons={template.Category === "MARKETING" ? [...buttons, OPT_OUT_TEXT] : buttons}
            className="w-full"
            compact
          />
        </div>

        {buttons.length > 0 && (
          <div className="flex flex-col gap-3">
            {buttons.map((buttonText) => {
              const flow = flowByButton.get(buttonText.toLowerCase());
              const target = flow ? byID.get(flow.TargetTemplateID) : undefined;
              const looping = flow ? path.includes(flow.TargetTemplateID) : false;
              const tooDeep = depth >= MAX_DEPTH;

              return (
                <div key={buttonText} className="flex items-center gap-3">
                  <div className="flex shrink-0 items-center">
                    <span className="h-px w-4 bg-gray-300 dark:bg-gray-600" />
                    <span
                      className="whitespace-nowrap rounded-full border px-2 py-0.5 text-[11px] font-medium"
                      style={{
                        borderColor: "var(--color-primary)",
                        color: "var(--color-primary)",
                      }}
                    >
                      {buttonText}
                    </span>
                    <span className="h-px w-5 bg-gray-300 dark:bg-gray-600" />
                    <span className="h-0 w-0 border-y-[4px] border-l-[6px] border-y-transparent border-l-gray-300 dark:border-l-gray-600" />
                  </div>

                  {flow && target && !looping && !tooDeep && renderNode(target, depth + 1, [
                    ...path,
                    target.ID,
                  ])}

                  {flow && looping && (
                    <span className="whitespace-nowrap text-[11px] text-red-600">
                      {`vuelve a ${flow.TargetName}`}
                    </span>
                  )}

                  {flow && tooDeep && (
                    <span className="whitespace-nowrap text-[11px] text-amber-600 dark:text-amber-400">
                      {`${flow.TargetName} (${MAX_DEPTH} niveles, se corta el dibujo)`}
                    </span>
                  )}

                  {!flow && !tooDeep && (
                    <button
                      type="button"
                      disabled={saving}
                      onClick={() =>
                        setPending({
                          sourceID: template.ID,
                          sourceName: template.Name,
                          buttonText,
                        })
                      }
                      className="whitespace-nowrap rounded-lg border border-dashed border-gray-300 px-2.5 py-1.5 text-[11px] font-medium text-[var(--color-primary)] transition-colors hover:border-[var(--color-primary)] disabled:opacity-40 dark:border-gray-600"
                    >
                      {"+ Crear respuesta"}
                    </button>
                  )}

                  {!flow && tooDeep && (
                    <span className="whitespace-nowrap text-[11px] text-gray-400">
                      {`Máximo ${MAX_DEPTH} niveles`}
                    </span>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    );
  };

  const notApproved = participates.filter((item) => item.Status !== "approved").length;
  const looseEnds = participates.reduce((total, item) => {
    const children = childrenBySource.get(item.ID) || [];
    const linked = new Set(children.map((flow) => flow.ButtonText.trim().toLowerCase()));
    return total + userButtons(item).filter((text) => !linked.has(text.toLowerCase())).length;
  }, 0);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap gap-4 px-1 text-[12px]">
        <span className="text-gray-500 dark:text-gray-400">
          {`${participates.length} plantilla(s) en el flujo`}
        </span>
        {notApproved > 0 && (
          <span className="text-amber-600 dark:text-amber-400">
            {`${notApproved} sin aprobar: esa rama no responde hasta que Meta las apruebe`}
          </span>
        )}
        {looseEnds > 0 && (
          <span className="text-amber-600 dark:text-amber-400">
            {`${looseEnds} botón(es) sin respuesta`}
          </span>
        )}
      </div>

      <div className="overflow-x-auto pb-2">
        <div className="flex min-w-max flex-col gap-5">
          {roots.map((root) => (
            <div
              key={root.ID}
              className="rounded-lg border border-gray-200 p-3 dark:border-gray-700"
            >
              {renderNode(root, 1, [root.ID])}
            </div>
          ))}
        </div>
      </div>

      <Modal
        isOpen={pending !== null}
        onClose={() => setPending(null)}
        title={(
          <span className="flex w-full flex-col items-start pr-8">
            <span className="text-lg font-semibold">{"Plantilla de respuesta"}</span>
            <span className="text-[13px] font-normal text-gray-400">
              {pending
                ? `Se envía cuando el cliente toca "${pending.buttonText}" en ${pending.sourceName}`
                : ""}
            </span>
          </span>
        )}
        size="4xl"
        zIndex={80}
        noPadding
        noBodyScroll
      >
        {pending !== null && (
          <TemplateForm
            businessId={businessId}
            variableCatalog={variableCatalog}
            asFlowResponse
            onSuccess={(created) => {
              if (created) {
                linkResponse(created);
                return;
              }
              setPending(null);
            }}
            onCancel={() => setPending(null)}
          />
        )}
      </Modal>
    </div>
  );
}
