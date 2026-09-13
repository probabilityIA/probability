"use client";

import { useRef, useState } from "react";
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
  flowId?: number;
  rootTemplateId?: number;
  title?: string;
  subtitle?: string;
  actions?: React.ReactNode;
  collapsed?: boolean;
  onToggleCollapsed?: () => void;
  readOnly?: boolean;
  studio?: boolean;
  onChanged: () => void;
}

interface PendingResponse {
  sourceID: number;
  sourceName: string;
  buttonText: string;
}

const OPT_OUT_TEXT = "Dejar de recibir";
const MAX_DEPTH = 10;
const FLOW_BLOCKED_SOURCES = ["sender.name", "campaign.name"];
const MIN_SCALE = 0.3;
const MAX_SCALE = 2;

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
  flowId,
  rootTemplateId,
  title,
  subtitle,
  actions,
  collapsed = false,
  onToggleCollapsed,
  readOnly = false,
  studio = false,
  onChanged,
}: TemplateFlowViewProps) {
  const { showToast } = useToast();
  const [pending, setPending] = useState<PendingResponse | null>(null);
  const [picking, setPicking] = useState<PendingResponse | null>(null);
  const [saving, setSaving] = useState(false);

  const [scale, setScale] = useState(1);
  const [offset, setOffset] = useState({ x: 0, y: 0 });
  const dragFrom = useRef<{ x: number; y: number; offsetX: number; offsetY: number } | null>(null);
  const [dragging, setDragging] = useState(false);

  const clampScale = (value: number) => Math.min(MAX_SCALE, Math.max(MIN_SCALE, value));

  const handleWheel = (event: React.WheelEvent) => {
    if (event.deltaY === 0) return;
    setScale((current) => clampScale(current - event.deltaY * 0.0015));
  };

  const handleMouseDown = (event: React.MouseEvent) => {
    if ((event.target as HTMLElement).closest("button")) return;
    dragFrom.current = {
      x: event.clientX,
      y: event.clientY,
      offsetX: offset.x,
      offsetY: offset.y,
    };
    setDragging(true);
  };

  const handleMouseMove = (event: React.MouseEvent) => {
    const from = dragFrom.current;
    if (!from) return;
    setOffset({
      x: from.offsetX + (event.clientX - from.x),
      y: from.offsetY + (event.clientY - from.y),
    });
  };

  const stopDragging = () => {
    dragFrom.current = null;
    setDragging(false);
  };

  const resetView = () => {
    setScale(1);
    setOffset({ x: 0, y: 0 });
  };

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

  const explicitRoot = rootTemplateId ? byID.get(rootTemplateId) : undefined;
  const roots = explicitRoot
    ? [explicitRoot]
    : participates.filter((item) => !parentCount.has(item.ID));

  const linkResponse = async (created: WhatsappTemplate, target?: PendingResponse) => {
    const pendingLink = target ?? pending;
    if (!pendingLink) return;

    const existing = flows
      .filter((flow) => flow.SourceTemplateID === pendingLink.sourceID)
      .map((flow) => ({
        button_text: flow.ButtonText,
        target_template_id: flow.TargetTemplateID,
        enabled: flow.Enabled,
      }));

    const next = [
      ...existing.filter(
        (item) => item.button_text.toLowerCase() !== pendingLink.buttonText.toLowerCase(),
      ),
      { button_text: pendingLink.buttonText, target_template_id: created.ID },
    ];

    setSaving(true);
    const result = await replaceTemplateFlowsAction(
      pendingLink.sourceID,
      next,
      businessId,
      flowId,
    );
    setSaving(false);
    setPending(null);
    setPicking(null);

    if (!result.success) {
      showToast(result.error || "No quedó enlazada", "error");
      return;
    }

    showToast(`"${pendingLink.buttonText}" ahora responde con ${created.Name}`, "success");
    onChanged();
  };

  if (!explicitRoot && participates.length === 0) {
    return (
      <p className="rounded-md border border-dashed border-gray-300 p-6 text-center text-sm text-gray-500 dark:border-gray-600">
        {
          "Ninguna plantilla tiene botones todavía. Agregá botones de respuesta para armar un flujo."
        }
      </p>
    );
  }

  if (rootTemplateId && !explicitRoot) {
    return (
      <p className="rounded-md border border-dashed border-gray-300 p-6 text-center text-sm text-gray-500 dark:border-gray-600">
        {"La plantilla inicial de este flujo ya no existe. Editá el flujo y elegí otra."}
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

                  {!flow && readOnly && (
                    <span className="whitespace-nowrap text-[11px] text-amber-600 dark:text-amber-400">
                      {"sin respuesta"}
                    </span>
                  )}

                  {!flow && !readOnly && !tooDeep && (
                    <div className="flex shrink-0 flex-col gap-1 rounded-lg border border-dashed border-gray-300 p-1.5 dark:border-gray-600">
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
                        className="whitespace-nowrap px-1 text-[11px] font-medium text-[var(--color-primary)] hover:underline disabled:opacity-40"
                      >
                        {"+ Crear respuesta"}
                      </button>
                      <button
                        type="button"
                        disabled={saving}
                        onClick={() =>
                          setPicking({
                            sourceID: template.ID,
                            sourceName: template.Name,
                            buttonText,
                          })
                        }
                        className="whitespace-nowrap px-1 text-[11px] font-medium text-gray-600 hover:underline disabled:opacity-40 dark:text-gray-300"
                      >
                        {"Elegir existente"}
                      </button>
                    </div>
                  )}

                  {!flow && !readOnly && tooDeep && (
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

  const reachable = new Set<number>();
  const collect = (id: number, depth: number) => {
    if (reachable.has(id) || depth > MAX_DEPTH) return;
    reachable.add(id);
    for (const step of childrenBySource.get(id) || []) {
      collect(step.TargetTemplateID, depth + 1);
    }
  };
  for (const root of roots) {
    collect(root.ID, 1);
  }

  const scoped = templates.filter((item) => reachable.has(item.ID));

  const notApproved = scoped.filter((item) => item.Status !== "approved").length;
  const looseEnds = scoped.reduce((total, item) => {
    const children = childrenBySource.get(item.ID) || [];
    const linked = new Set(children.map((flow) => flow.ButtonText.trim().toLowerCase()));
    return total + userButtons(item).filter((text) => !linked.has(text.toLowerCase())).length;
  }, 0);

  return (
    <div className="flex flex-col gap-4">
      <div
        className={
          studio
            ? "flex min-h-0 flex-1 flex-col"
            : "rounded-xl border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800"
        }
      >
        {(title || actions || onToggleCollapsed) && (
          <div className="mb-3 flex flex-wrap items-center gap-3">
            {onToggleCollapsed && (
              <button
                type="button"
                onClick={onToggleCollapsed}
                aria-label={collapsed ? "Expandir" : "Contraer"}
                className="flex h-6 w-6 shrink-0 items-center justify-center rounded text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
              >
                <svg
                  className={`h-4 w-4 transition-transform ${collapsed ? "" : "rotate-90"}`}
                  fill="none"
                  stroke="currentColor"
                  strokeWidth={2}
                  viewBox="0 0 24 24"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
                </svg>
              </button>
            )}
            {title && (
              <div className="min-w-0">
                <p className="truncate text-sm font-semibold text-gray-900 dark:text-white">
                  {title}
                </p>
                {subtitle && (
                  <p className="truncate text-xs text-gray-500 dark:text-gray-400">{subtitle}</p>
                )}
              </div>
            )}
            {actions && <div className="ml-auto flex items-center gap-3">{actions}</div>}
          </div>
        )}

        <div className={`flex flex-wrap gap-4 text-[12px] ${collapsed ? "" : "mb-4"}`}>
          <span className="text-gray-500 dark:text-gray-400">
            {`${scoped.length} plantilla(s) en el flujo`}
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

        {!studio && (
          <div
            className={`max-h-[420px] overflow-auto rounded-lg border border-gray-100 bg-gray-50 p-4 dark:border-gray-700 dark:bg-gray-900/30 ${
              collapsed ? "hidden" : ""
            }`}
          >
            <div className="flex min-w-max flex-col gap-5">
              {roots.map((root) => (
                <div key={root.ID}>{renderNode(root, 1, [root.ID])}</div>
              ))}
            </div>
          </div>
        )}

        <div className={`relative ${studio ? "" : "hidden"}`}>
          <div className="absolute right-2 top-2 z-10 flex items-center gap-1 rounded-lg border border-gray-200 bg-white p-1 shadow-sm dark:border-gray-700 dark:bg-gray-800">
            <button
              type="button"
              onClick={() => setScale((current) => clampScale(current - 0.15))}
              className="h-6 w-6 rounded text-sm font-semibold text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
            >
              {"-"}
            </button>
            <span className="w-10 text-center text-[11px] text-gray-500">
              {`${Math.round(scale * 100)}%`}
            </span>
            <button
              type="button"
              onClick={() => setScale((current) => clampScale(current + 0.15))}
              className="h-6 w-6 rounded text-sm font-semibold text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
            >
              {"+"}
            </button>
            <button
              type="button"
              onClick={resetView}
              className="rounded px-1.5 text-[11px] font-medium text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
            >
              {"Centrar"}
            </button>
          </div>

          <div
            onWheel={handleWheel}
            onMouseDown={handleMouseDown}
            onMouseMove={handleMouseMove}
            onMouseUp={stopDragging}
            onMouseLeave={stopDragging}
            style={{
              cursor: dragging ? "grabbing" : "grab",
              backgroundImage: studio
                ? "radial-gradient(circle, rgba(17,27,33,0.16) 1px, transparent 1px)"
                : undefined,
              backgroundSize: studio ? `${24 * scale}px ${24 * scale}px` : undefined,
              backgroundPosition: studio ? `${offset.x}px ${offset.y}px` : undefined,
            }}
            className="h-[calc(90vh-230px)] overflow-hidden rounded-lg border border-gray-100 bg-[#f4f1ec] dark:border-gray-700 dark:bg-gray-900"
          >
            <div
              style={{
                transform: `translate(${offset.x}px, ${offset.y}px) scale(${scale})`,
                transformOrigin: "top left",
              }}
              className="flex min-w-max flex-col gap-5 p-4"
            >
              {roots.map((root) => (
                <div key={root.ID}>{renderNode(root, 1, [root.ID])}</div>
              ))}
            </div>
          </div>
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

      <Modal
        isOpen={picking !== null}
        onClose={() => setPicking(null)}
        title={(
          <span className="flex w-full flex-col items-start pr-8">
            <span className="text-lg font-semibold">{"Elegir plantilla existente"}</span>
            <span className="text-[13px] font-normal text-gray-400">
              {picking
                ? `Se envía cuando el cliente toca "${picking.buttonText}" en ${picking.sourceName}`
                : ""}
            </span>
          </span>
        )}
        size="4xl"
        zIndex={80}
      >
        {picking !== null && (
          <PickTemplate
            templates={templates.filter(
              (item) =>
                item.ID !== picking.sourceID &&
                item.Origin === "business" &&
                !(item.Variables ?? []).some((variable) =>
                  FLOW_BLOCKED_SOURCES.includes(variable.Source),
                ),
            )}
            saving={saving}
            onPick={(template) => linkResponse(template, picking)}
          />
        )}
      </Modal>
    </div>
  );
}

interface PickTemplateProps {
  templates: WhatsappTemplate[];
  saving: boolean;
  onPick: (template: WhatsappTemplate) => void;
}

function PickTemplate({ templates, saving, onPick }: PickTemplateProps) {
  if (templates.length === 0) {
    return (
      <p className="rounded-md border border-dashed border-gray-300 p-6 text-center text-sm text-gray-500 dark:border-gray-600">
        {"No hay plantillas que sirvan como respuesta. Creá una nueva desde el botón anterior."}
      </p>
    );
  }

  return (
    <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      {templates.map((template) => (
        <button
          key={template.ID}
          type="button"
          disabled={saving}
          onClick={() => onPick(template)}
          className="flex flex-col gap-2 rounded-xl border border-gray-200 p-2.5 text-left transition-colors hover:border-[var(--color-primary)] disabled:opacity-40 dark:border-gray-700"
        >
          <div className="flex items-center justify-between gap-2">
            <span className="min-w-0 truncate text-[12px] font-semibold text-gray-900 dark:text-white">
              {template.Name}
            </span>
            <span
              className={`shrink-0 rounded-full px-1.5 py-0.5 text-[10px] ${
                STATUS_STYLE[template.Status] || "bg-gray-100 text-gray-600"
              }`}
            >
              {TEMPLATE_STATUS_LABEL[template.Status] || template.Status}
            </span>
          </div>

          <div className="rounded-lg bg-[#e9e2d9] p-2 dark:bg-[#2a2724]">
            <TemplateBubble
              headerType={template.HeaderType}
              headerMediaURL={template.HeaderMediaURL}
              headerText={template.HeaderText}
              bodyText={fillPlaceholders(template.BodyText, template.Variables)}
              footerText={template.FooterText}
              buttons={userButtons(template)}
              className="w-full"
              compact
            />
          </div>
        </button>
      ))}
    </div>
  );
}
