"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { SystemFlow, SystemFlowNode, WhatsappTemplate } from "../../domain/scheduled-types";
import { TemplateBubble, fillPlaceholders } from "./TemplateBubble";

interface SystemFlowViewProps {
  flow: SystemFlow;
  templates: WhatsappTemplate[];
  collapsed: boolean;
  onToggleCollapsed: () => void;
}

interface Step {
  label: string;
  node: SystemFlowNode;
}

function nodeTitle(node: SystemFlowNode): string {
  return node.title || node.template || "";
}

function countTemplates(node: SystemFlowNode | undefined, seen: Set<string>): Set<string> {
  if (!node) return seen;
  if (node.kind === "template" && node.template) seen.add(node.template);
  for (const branch of node.branches || []) countTemplates(branch.node, seen);
  return seen;
}

function findTitle(node: SystemFlowNode | undefined, template: string): string {
  if (!node) return "";
  if (node.template === template) return nodeTitle(node);
  for (const branch of node.branches || []) {
    const found = findTitle(branch.node, template);
    if (found) return found;
  }
  return "";
}

function samePath(a: Step[], b: Step[]): boolean {
  if (a.length > b.length) return false;
  return a.every((step, index) => step.node === b[index].node);
}

export function SystemFlowView({ flow, templates, collapsed, onToggleCollapsed }: SystemFlowViewProps) {
  const rootPath = useMemo<Step[]>(() => [{ label: "", node: flow.root }], [flow.root]);
  const [path, setPath] = useState<Step[]>(rootPath);
  const chatRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setPath(rootPath);
  }, [rootPath]);

  useEffect(() => {
    const chat = chatRef.current;
    if (chat) chat.scrollTo({ top: chat.scrollHeight, behavior: "smooth" });
  }, [path.length]);

  const byName = useMemo(() => {
    const map = new Map<string, WhatsappTemplate>();
    for (const template of templates) {
      if (!map.has(template.Name) || template.Origin === "system") map.set(template.Name, template);
    }
    return map;
  }, [templates]);

  const current = path[path.length - 1].node;

  const choose = (label: string) => {
    const branch = (current.branches || []).find((item) => item.label === label);
    if (!branch) return;
    if (branch.node) {
      setPath((steps) => [...steps, { label, node: branch.node as SystemFlowNode }]);
      return;
    }
    if (branch.back_to) {
      const index = path.findIndex((step) => step.node.template === branch.back_to);
      if (index >= 0) setPath((steps) => [...steps.slice(0, index + 1)]);
    }
  };

  const bodyFor = (node: SystemFlowNode) => {
    const template = node.template ? byName.get(node.template) : undefined;
    if (template) return fillPlaceholders(template.BodyText, template.Variables);
    return node.body || node.description || "";
  };

  const renderOutline = (node: SystemFlowNode, trail: Step[], depth: number): React.ReactNode => (
    <ul className={depth === 0 ? "space-y-1" : "ml-3 space-y-1 border-l border-gray-200 pl-3 dark:border-gray-700"}>
      {(node.branches || []).map((branch) => {
        const nextTrail = branch.node ? [...trail, { label: branch.label, node: branch.node }] : trail;
        const active = branch.node ? samePath(nextTrail, path) : false;
        return (
          <li key={`${depth}-${branch.label}`}>
            <button
              type="button"
              onClick={() => (branch.node ? setPath(nextTrail) : undefined)}
              disabled={!branch.node}
              className={`flex w-full items-start gap-1.5 rounded-md px-1.5 py-1 text-left text-[12px] transition-colors ${
                active
                  ? "bg-violet-50 text-[#3E0FA8] dark:bg-violet-500/15 dark:text-violet-200"
                  : "text-gray-600 hover:bg-gray-50 disabled:cursor-default disabled:hover:bg-transparent dark:text-gray-300 dark:hover:bg-gray-700/50"
              }`}
            >
              <span
                className={`mt-0.5 shrink-0 rounded-full border px-1.5 text-[10px] font-medium ${
                  node.kind === "decision"
                    ? "border-amber-400 text-amber-700 dark:text-amber-300"
                    : "border-[var(--color-primary)] text-[var(--color-primary)]"
                }`}
              >
                {branch.label}
              </span>
              <span className="min-w-0">
                {branch.node ? nodeTitle(branch.node) : `vuelve a ${findTitle(flow.root, branch.back_to || "") || branch.back_to}`}
              </span>
            </button>
            {branch.node && (branch.node.branches || []).length > 0 && renderOutline(branch.node, nextTrail, depth + 1)}
          </li>
        );
      })}
    </ul>
  );

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

      <div className="mt-2 flex flex-wrap gap-4 text-[12px]">
        <span className="text-gray-500 dark:text-gray-400">{`${total} plantilla(s) en el flujo`}</span>
        <span className="text-gray-500 dark:text-gray-400">{flow.trigger}</span>
      </div>

      {!collapsed && (
        <div className="mt-4 grid gap-4 lg:grid-cols-[minmax(0,320px)_minmax(0,1fr)]">
          <div className="rounded-lg border border-gray-200 p-3 dark:border-gray-700">
            <p className="mb-2 text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
              {"Caminos del flujo"}
            </p>
            <button
              type="button"
              onClick={() => setPath(rootPath)}
              className={`mb-1 w-full rounded-md px-1.5 py-1 text-left text-[12px] font-semibold ${
                path.length === 1
                  ? "bg-violet-50 text-[#3E0FA8] dark:bg-violet-500/15 dark:text-violet-200"
                  : "text-gray-800 hover:bg-gray-50 dark:text-gray-100 dark:hover:bg-gray-700/50"
              }`}
            >
              {nodeTitle(flow.root)}
            </button>
            {renderOutline(flow.root, rootPath, 0)}
          </div>

          <div className="flex min-h-0 flex-col overflow-hidden rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="flex items-center gap-2 border-b border-gray-200 bg-gray-50 px-3 py-2 dark:border-gray-700 dark:bg-gray-900/40">
              <span className="text-[12px] font-semibold text-gray-700 dark:text-gray-200">{"Simulaci\u00f3n de la conversaci\u00f3n"}</span>
              <span className="text-[11px] text-gray-400">{"Toca los botones como lo har\u00eda el cliente"}</span>
              <div className="ml-auto flex gap-2">
                <button
                  type="button"
                  disabled={path.length <= 1}
                  onClick={() => setPath((steps) => steps.slice(0, -1))}
                  className="rounded-md border border-gray-200 px-2 py-0.5 text-[11px] text-gray-600 hover:bg-white disabled:opacity-40 dark:border-gray-600 dark:text-gray-300"
                >
                  {"Paso atr\u00e1s"}
                </button>
                <button
                  type="button"
                  disabled={path.length <= 1}
                  onClick={() => setPath(rootPath)}
                  className="rounded-md border border-gray-200 px-2 py-0.5 text-[11px] text-gray-600 hover:bg-white disabled:opacity-40 dark:border-gray-600 dark:text-gray-300"
                >
                  {"Reiniciar"}
                </button>
              </div>
            </div>

            <div ref={chatRef} className="max-h-[560px] overflow-y-auto p-4" style={{ backgroundColor: "#efeae2" }}>
              <div className="mx-auto flex w-full max-w-xl flex-col gap-3">
              {path.map((step, index) => {
                const parent = index > 0 ? path[index - 1].node : undefined;
                return (
                  <div key={`${index}-${nodeTitle(step.node)}`} className="flex flex-col gap-2">
                    {parent && parent.kind === "template" && (
                      <div className="self-end rounded-[10px] rounded-br-[3px] bg-white px-3 py-1.5 text-[13px] text-gray-800 shadow-sm">
                        {step.label}
                      </div>
                    )}
                    {parent && parent.kind === "decision" && (
                      <div className="self-center rounded-full bg-amber-100 px-3 py-1 text-[11px] font-medium text-amber-800">
                        {`Caso: ${step.label}`}
                      </div>
                    )}

                    {step.node.kind === "decision" ? (
                      <div className="self-center w-full max-w-sm rounded-lg border border-amber-300 bg-amber-50 p-3 text-center">
                        <p className="text-[10px] font-semibold uppercase tracking-wide text-amber-700">{"Decisi\u00f3n del sistema"}</p>
                        <p className="text-[13px] font-semibold text-gray-900">{step.node.title}</p>
                        {step.node.detail && <p className="text-[11px] text-gray-600">{step.node.detail}</p>}
                      </div>
                    ) : (
                      <div className="flex max-w-[340px] flex-col gap-1">
                        <span className="text-[10px] font-medium text-gray-500">{nodeTitle(step.node)}</span>
                        <TemplateBubble
                          bodyText={bodyFor(step.node)}
                          buttons={(step.node.branches || []).length > 0 ? (step.node.branches || []).map((b) => b.label) : step.node.buttons || []}
                          className="w-full"
                          compact
                        />
                        {step.node.detail && <span className="text-[11px] text-gray-500">{step.node.detail}</span>}
                      </div>
                    )}

                    {(step.node.effects || []).map((effect) => (
                      <div key={effect} className="self-center rounded-full bg-sky-100 px-3 py-1 text-[11px] font-medium text-sky-800">
                        {effect}
                      </div>
                    ))}
                  </div>
                );
              })}

              {(current.branches || []).length === 0 && (
                <div className="self-center rounded-full bg-gray-200 px-3 py-1 text-[11px] text-gray-600">
                  {"Fin de la conversaci\u00f3n"}
                </div>
              )}
              </div>
            </div>

            {(current.branches || []).length > 0 && (
              <div className="flex flex-wrap items-center gap-2 border-t border-gray-200 bg-white px-3 py-2 dark:border-gray-700 dark:bg-gray-800">
                <span className="text-[11px] text-gray-500 dark:text-gray-400">
                  {current.kind === "decision" ? "Elige el caso:" : "Responde como el cliente:"}
                </span>
                {(current.branches || []).map((branch) => (
                  <button
                    key={branch.label}
                    type="button"
                    onClick={() => choose(branch.label)}
                    className={`rounded-full border px-3 py-1 text-[12px] font-medium transition-colors ${
                      current.kind === "decision"
                        ? "border-amber-400 text-amber-700 hover:bg-amber-50 dark:text-amber-300"
                        : "border-[var(--color-primary)] text-[var(--color-primary)] hover:bg-violet-50 dark:hover:bg-violet-500/10"
                    }`}
                  >
                    {branch.label}
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
