"use client";

import {
  TEMPLATE_STATUS_LABEL,
  TemplateFlow,
  WhatsappTemplate,
} from "../../domain/scheduled-types";

interface TemplateFlowViewProps {
  templates: WhatsappTemplate[];
  flows: TemplateFlow[];
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

export function TemplateFlowView({ templates, flows }: TemplateFlowViewProps) {
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
      userButtons(item).length > 0 ||
      childrenBySource.has(item.ID) ||
      parentCount.has(item.ID),
  );

  const roots = participates.filter((item) => !parentCount.has(item.ID));

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
      <div key={`${template.ID}-${depth}`} className="flex flex-col gap-1.5">
        <div className="flex items-center gap-2">
          <span className="truncate text-[13px] font-semibold text-gray-900 dark:text-white">
            {template.Name}
          </span>
          <span
            className={`shrink-0 rounded-full px-2 py-0.5 text-[11px] ${
              STATUS_STYLE[template.Status] || "bg-gray-100 text-gray-600"
            }`}
          >
            {TEMPLATE_STATUS_LABEL[template.Status] || template.Status}
          </span>
          {reused && (
            <span className="shrink-0 rounded-full bg-blue-50 px-2 py-0.5 text-[11px] text-blue-600">
              {"Reutilizada"}
            </span>
          )}
        </div>

        {buttons.length === 0 ? (
          <span className="pl-3 text-[12px] text-gray-400">{"Sin botones"}</span>
        ) : (
          <div className="flex flex-col gap-1.5 border-l border-gray-200 pl-3 dark:border-gray-700">
            {buttons.map((buttonText) => {
              const flow = flowByButton.get(buttonText.toLowerCase());
              const target = flow ? byID.get(flow.TargetTemplateID) : undefined;
              const looping = flow ? path.includes(flow.TargetTemplateID) : false;
              const tooDeep = depth >= MAX_DEPTH;

              return (
                <div key={buttonText} className="flex flex-col gap-1.5">
                  <div className="flex items-center gap-2">
                    <span className="rounded-full border border-dashed border-gray-300 px-2 py-0.5 text-[12px] text-gray-700 dark:border-gray-600 dark:text-gray-200">
                      {buttonText}
                    </span>
                    {!flow && (
                      <span className="text-[12px] text-amber-600 dark:text-amber-400">
                        {"→ sin respuesta"}
                      </span>
                    )}
                    {flow && !flow.Enabled && (
                      <span className="text-[12px] text-gray-400">{"→ desactivado"}</span>
                    )}
                  </div>

                  {flow && target && !looping && !tooDeep && (
                    <div className="pl-3">
                      {renderNode(target, depth + 1, [...path, target.ID])}
                    </div>
                  )}

                  {flow && looping && (
                    <span className="pl-3 text-[12px] text-red-600">
                      {`→ vuelve a ${flow.TargetName}`}
                    </span>
                  )}

                  {flow && tooDeep && (
                    <span className="pl-3 text-[12px] text-amber-600 dark:text-amber-400">
                      {`→ ${flow.TargetName} (se corta el dibujo, son ${MAX_DEPTH} niveles)`}
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

  const pending = participates.filter((item) => item.Status !== "approved").length;
  const looseEnds = participates.reduce((total, item) => {
    const children = childrenBySource.get(item.ID) || [];
    const linked = new Set(children.map((flow) => flow.ButtonText.trim().toLowerCase()));
    return total + userButtons(item).filter((text) => !linked.has(text.toLowerCase())).length;
  }, 0);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap gap-4 text-[12px]">
        <span className="text-gray-500 dark:text-gray-400">
          {`${participates.length} plantilla(s) en el flujo`}
        </span>
        {pending > 0 && (
          <span className="text-amber-600 dark:text-amber-400">
            {`${pending} sin aprobar: esa rama no responde hasta que Meta las apruebe`}
          </span>
        )}
        {looseEnds > 0 && (
          <span className="text-amber-600 dark:text-amber-400">
            {`${looseEnds} botón(es) sin respuesta`}
          </span>
        )}
      </div>

      <div className="flex flex-col gap-4">
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
  );
}
