"use client";

import { SystemFlow, SystemFlowNode, WhatsappTemplate } from "../../domain/scheduled-types";
import { TemplateBubble, fillPlaceholders } from "./TemplateBubble";

interface SystemFlowViewProps {
  flow: SystemFlow;
  templates: WhatsappTemplate[];
  collapsed: boolean;
  onToggleCollapsed: () => void;
}

function countTemplates(node: SystemFlowNode | undefined, seen: Set<string>): Set<string> {
  if (!node) return seen;
  if (node.kind === "template" && node.template) seen.add(node.template);
  for (const branch of node.branches || []) countTemplates(branch.node, seen);
  return seen;
}

function Connector({ label, tone = "primary" }: { label: string; tone?: "primary" | "decision" }) {
  const style =
    tone === "decision"
      ? { borderColor: "#d97706", color: "#b45309" }
      : { borderColor: "var(--color-primary)", color: "var(--color-primary)" };
  return (
    <div className="flex shrink-0 items-center">
      <span className="h-px w-4 bg-gray-300 dark:bg-gray-600" />
      <span className="whitespace-nowrap rounded-full border px-2 py-0.5 text-[11px] font-medium" style={style}>
        {label}
      </span>
      <span className="h-px w-5 bg-gray-300 dark:bg-gray-600" />
      <span className="h-0 w-0 border-y-[4px] border-l-[6px] border-y-transparent border-l-gray-300 dark:border-l-gray-600" />
    </div>
  );
}

export function SystemFlowView({ flow, templates, collapsed, onToggleCollapsed }: SystemFlowViewProps) {
  const byName = new Map<string, WhatsappTemplate>();
  for (const template of templates) {
    if (!byName.has(template.Name) || template.Origin === "system") byName.set(template.Name, template);
  }

  const renderNode = (node: SystemFlowNode): React.ReactNode => {
    const branches = node.branches || [];

    if (node.kind === "decision") {
      return (
        <div className="flex items-center gap-3">
          <div className="flex w-[230px] shrink-0 flex-col gap-1 rounded-lg border border-amber-300 bg-amber-50 p-2.5 dark:border-amber-500/40 dark:bg-amber-500/10">
            <span className="text-[10px] font-semibold uppercase tracking-wide text-amber-700 dark:text-amber-300">
              {"Decisión del sistema"}
            </span>
            <span className="text-[12px] font-semibold text-gray-900 dark:text-white">{node.title}</span>
            {node.detail && <span className="text-[11px] text-gray-600 dark:text-gray-300">{node.detail}</span>}
          </div>
          <div className="flex flex-col gap-3">
            {branches.map((branch) => (
              <div key={branch.label} className="flex items-center gap-3">
                <Connector label={branch.label} tone="decision" />
                {branch.node && renderNode(branch.node)}
              </div>
            ))}
          </div>
        </div>
      );
    }

    const template = node.template ? byName.get(node.template) : undefined;
    const body = template ? fillPlaceholders(template.BodyText, template.Variables) : node.body || node.description || "";
    const buttons = branches.length > 0 ? branches.map((branch) => branch.label) : node.buttons || [];

    return (
      <div className="flex items-center gap-3">
        <div className="flex w-[230px] shrink-0 flex-col gap-1.5">
          <div className="flex flex-wrap items-center gap-1.5">
            <span className="truncate text-[12px] font-semibold text-gray-900 dark:text-white">{node.template}</span>
            <span className="shrink-0 rounded-full bg-violet-100 px-1.5 py-0.5 text-[10px] text-[#3E0FA8] dark:bg-violet-500/20 dark:text-violet-200">
              {"Sistema"}
            </span>
          </div>
          <TemplateBubble bodyText={body} buttons={buttons} className="w-full" compact />
          {node.detail && <p className="text-[11px] text-gray-500 dark:text-gray-400">{node.detail}</p>}
          {(node.effects || []).map((effect) => (
            <span
              key={effect}
              className="w-fit rounded-md bg-sky-50 px-1.5 py-0.5 text-[11px] text-sky-800 dark:bg-sky-500/10 dark:text-sky-200"
            >
              {effect}
            </span>
          ))}
        </div>
        {branches.length > 0 && (
          <div className="flex flex-col gap-3">
            {branches.map((branch) => (
              <div key={branch.label} className="flex items-center gap-3">
                <Connector label={branch.label} />
                {branch.node && renderNode(branch.node)}
                {branch.back_to && (
                  <span className="whitespace-nowrap text-[11px] text-gray-500 dark:text-gray-400">
                    {`vuelve a ${branch.back_to}`}
                  </span>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    );
  };

  const total = countTemplates(flow.root, new Set()).size;

  return (
    <div className="rounded-xl border border-violet-200 bg-white p-4 dark:border-violet-500/30 dark:bg-gray-800">
      <div className="flex flex-wrap items-center gap-3">
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
        <div className="min-w-0">
          <div className="flex items-center gap-2">
            <p className="truncate text-sm font-semibold text-gray-900 dark:text-white">{flow.name}</p>
            <span className="rounded-full bg-violet-100 px-2 py-0.5 text-[10px] font-semibold text-[#3E0FA8] dark:bg-violet-500/20 dark:text-violet-200">
              {"Del sistema"}
            </span>
          </div>
          <p className="truncate text-xs text-gray-500 dark:text-gray-400">{flow.description}</p>
        </div>
        <span className="ml-auto text-[11px] text-gray-400">{"Solo lectura"}</span>
      </div>

      <div className={`mt-2 flex flex-wrap gap-4 text-[12px] ${collapsed ? "" : "mb-4"}`}>
        <span className="text-gray-500 dark:text-gray-400">{`${total} plantilla(s) en el flujo`}</span>
        <span className="text-gray-500 dark:text-gray-400">{flow.trigger}</span>
      </div>

      {!collapsed && <div className="overflow-x-auto pb-2">{renderNode(flow.root)}</div>}
    </div>
  );
}
