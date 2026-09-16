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

type ViewMode = "diagram" | "simulation";

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

function stepKey(trail: Step[]): string {
  return trail
    .slice(1)
    .map((step) => step.label)
    .join(" / ");
}

function samePath(a: Step[], b: Step[]): boolean {
  if (a.length > b.length) return false;
  return a.every((step, index) => step.node === b[index].node);
}

const pillPrimary = "border-[var(--color-primary)] text-[var(--color-primary)]";
const pillDecision = "border-amber-400 text-amber-700 dark:text-amber-300";

export function SystemFlowView({ flow, templates, collapsed, onToggleCollapsed }: SystemFlowViewProps) {
  const rootPath = useMemo<Step[]>(() => [{ label: "", node: flow.root }], [flow.root]);
  const [mode, setMode] = useState<ViewMode>("diagram");
  const [path, setPath] = useState<Step[]>(rootPath);
  const [focused, setFocused] = useState<string | null>(null);
  const chatRef = useRef<HTMLDivElement>(null);
  const diagramRef = useRef<HTMLDivElement>(null);
  const nodeRefs = useRef(new Map<string, HTMLDivElement>());

  useEffect(() => {
    setPath(rootPath);
    setFocused(null);
  }, [rootPath]);

  useEffect(() => {
    const chat = chatRef.current;
    if (mode === "simulation" && chat) chat.scrollTo({ top: chat.scrollHeight, behavior: "smooth" });
  }, [path.length, mode]);

  useEffect(() => {
    if (mode !== "diagram" || collapsed) return;
    const container = diagramRef.current;
    const target = focused === null ? undefined : nodeRefs.current.get(focused);
    if (!container || !target) return;
    const top = target.getBoundingClientRect().top - container.getBoundingClientRect().top + container.scrollTop - 12;
    container.scrollTo({ top, behavior: "smooth" });
  }, [focused, mode, collapsed]);

  const byName = useMemo(() => {
    const map = new Map<string, WhatsappTemplate>();
    for (const template of templates) {
      if (!map.has(template.Name) || template.Origin === "system") map.set(template.Name, template);
    }
    return map;
  }, [templates]);

  const current = path[path.length - 1].node;

  const onFocusedPath = (key: string) =>
    mode === "diagram" && focused !== null && (key === "" || focused === key || focused.startsWith(`${key} / `));

  const goTo = (trail: Step[]) => {
    setPath(trail);
    setFocused(stepKey(trail));
  };

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
        const active = branch.node
          ? mode === "diagram"
            ? focused === stepKey(nextTrail)
            : samePath(nextTrail, path)
          : false;
        return (
          <li key={`${depth}-${branch.label}`}>
            <button
              type="button"
              onClick={() => (branch.node ? goTo(nextTrail) : undefined)}
              disabled={!branch.node}
              className={`flex w-full items-start gap-1.5 rounded-md px-1.5 py-1 text-left text-[12px] transition-colors ${
                active
                  ? "bg-violet-50 text-[#3E0FA8] dark:bg-violet-500/15 dark:text-violet-200"
                  : "text-gray-600 hover:bg-gray-50 disabled:cursor-default disabled:hover:bg-transparent dark:text-gray-300 dark:hover:bg-gray-700/50"
              }`}
            >
              <span className={`mt-0.5 shrink-0 rounded-full border px-1.5 text-[10px] font-medium ${node.kind === "decision" ? pillDecision : pillPrimary}`}>
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

  const renderCard = (node: SystemFlowNode, highlighted: boolean, onPath = false) =>
    node.kind === "decision" ? (
      <div className={`w-full max-w-sm rounded-lg border border-amber-300 bg-amber-50 p-3 transition-shadow ${highlighted ? "ring-2 ring-[var(--color-primary)] ring-offset-2" : onPath ? "ring-2 ring-violet-300" : ""}`}>
        <p className="text-[10px] font-semibold uppercase tracking-wide text-amber-700">{"Decisi\u00f3n del sistema"}</p>
        <p className="text-[13px] font-semibold text-gray-900">{node.title}</p>
        {node.detail && <p className="text-[11px] text-gray-600">{node.detail}</p>}
      </div>
    ) : (
      <div className={`flex w-full max-w-[340px] flex-col gap-1 rounded-xl p-1.5 transition-all ${highlighted ? "bg-violet-100 ring-2 ring-[var(--color-primary)] dark:bg-violet-500/20" : onPath ? "bg-violet-50 ring-1 ring-violet-300 dark:bg-violet-500/10" : ""}`}>
        <div className="flex flex-wrap items-center gap-1.5">
          <span className="text-[12px] font-semibold text-gray-900 dark:text-white">{nodeTitle(node)}</span>
          <span className="rounded-full bg-violet-100 px-1.5 py-0.5 text-[10px] text-[#3E0FA8] dark:bg-violet-500/20 dark:text-violet-200">
            {node.template}
          </span>
        </div>
        <TemplateBubble
          bodyText={bodyFor(node)}
          buttons={(node.branches || []).length > 0 ? (node.branches || []).map((b) => b.label) : node.buttons || []}
          className="w-full"
          compact
        />
        {node.detail && <span className="text-[11px] text-gray-500 dark:text-gray-400">{node.detail}</span>}
      </div>
    );

  const renderDiagramNode = (node: SystemFlowNode, trail: Step[]): React.ReactNode => {
    const key = stepKey(trail);
    const branches = node.branches || [];
    return (
      <div
        ref={(element) => {
          if (element) nodeRefs.current.set(key, element);
          else nodeRefs.current.delete(key);
        }}
        className="flex flex-col gap-2"
      >
        {renderCard(node, mode === "diagram" && focused === key, onFocusedPath(key))}
        {(node.effects || []).length > 0 && (
          <div className="flex flex-wrap gap-1.5">
            {(node.effects || []).map((effect) => (
              <span key={effect} className="rounded-md bg-sky-50 px-1.5 py-0.5 text-[11px] text-sky-800 dark:bg-sky-500/10 dark:text-sky-200">
                {effect}
              </span>
            ))}
          </div>
        )}
        {branches.length > 0 && (
          <div className="ml-3 flex flex-col gap-4 border-l-2 border-dashed border-gray-300 pl-5 pt-1 dark:border-gray-600">
            {branches.map((branch) => {
              const childKey = stepKey([...trail, { label: branch.label, node: branch.node || node }]);
              const branchOnPath = Boolean(branch.node) && onFocusedPath(childKey);
              return (
              <div key={branch.label} className="flex flex-col gap-2">
                <div className="-ml-5 flex items-center gap-1.5">
                  <span className={`h-0.5 w-4 ${branchOnPath ? "bg-[var(--color-primary)]" : "bg-gray-300 dark:bg-gray-600"}`} />
                  <button
                    type="button"
                    disabled={!branch.node}
                    onClick={() => branch.node && goTo([...trail, { label: branch.label, node: branch.node }])}
                    className={`whitespace-nowrap rounded-full border px-2 py-0.5 text-[11px] font-medium ${
                      branchOnPath
                        ? "border-[var(--color-primary)] bg-[var(--color-primary)] text-white"
                        : `bg-white dark:bg-gray-800 ${node.kind === "decision" ? pillDecision : pillPrimary}`
                    }`}
                  >
                    {branch.label}
                  </button>
                  {branch.back_to && (
                    <button
                      type="button"
                      onClick={() => {
                        const index = trail.findIndex((step) => step.node.template === branch.back_to);
                        if (index >= 0) goTo(trail.slice(0, index + 1));
                      }}
                      className="text-[11px] text-gray-500 hover:underline dark:text-gray-400"
                    >
                      {`vuelve a ${findTitle(flow.root, branch.back_to) || branch.back_to}`}
                    </button>
                  )}
                </div>
                {branch.node && renderDiagramNode(branch.node, [...trail, { label: branch.label, node: branch.node }])}
              </div>
              );
            })}
          </div>
        )}
      </div>
    );
  };

  const total = countTemplates(flow.root, new Set()).size;
  const tabClass = (active: boolean) =>
    `rounded-md px-2.5 py-1 text-[12px] font-medium transition-colors ${
      active
        ? "bg-white text-gray-900 shadow-sm dark:bg-gray-700 dark:text-white"
        : "text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-200"
    }`;

  return (
    <div className="rounded-xl border border-violet-200 bg-white p-4 dark:border-violet-500/30 dark:bg-gray-800">
      <div className="flex flex-wrap items-center gap-3">
        <button
          type="button"
          onClick={onToggleCollapsed}
          aria-label={collapsed ? "Expandir" : "Contraer"}
          className="flex h-6 w-6 shrink-0 items-center justify-center rounded text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
        >
          <svg className={`h-4 w-4 transition-transform ${collapsed ? "" : "rotate-90"}`} fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
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
        <div className="mt-4 grid gap-4 lg:grid-cols-[minmax(0,300px)_minmax(0,1fr)]">
          <div className="max-h-[680px] overflow-y-auto rounded-lg border border-gray-200 p-3 dark:border-gray-700">
            <p className="mb-2 text-[11px] font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">{"Caminos del flujo"}</p>
            <button
              type="button"
              onClick={() => goTo(rootPath)}
              className={`mb-1 w-full rounded-md px-1.5 py-1 text-left text-[12px] font-semibold ${
                (mode === "diagram" ? focused === "" : path.length === 1)
                  ? "bg-violet-50 text-[#3E0FA8] dark:bg-violet-500/15 dark:text-violet-200"
                  : "text-gray-800 hover:bg-gray-50 dark:text-gray-100 dark:hover:bg-gray-700/50"
              }`}
            >
              {nodeTitle(flow.root)}
            </button>
            {renderOutline(flow.root, rootPath, 0)}
          </div>

          <div className="flex min-h-0 flex-col overflow-hidden rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="flex flex-wrap items-center gap-2 border-b border-gray-200 bg-gray-50 px-3 py-2 dark:border-gray-700 dark:bg-gray-900/40">
              <div className="flex gap-1 rounded-lg bg-gray-100 p-0.5 dark:bg-gray-900">
                <button type="button" onClick={() => setMode("diagram")} className={tabClass(mode === "diagram")}>
                  {"Flujo completo"}
                </button>
                <button type="button" onClick={() => setMode("simulation")} className={tabClass(mode === "simulation")}>
                  {"Simular conversaci\u00f3n"}
                </button>
              </div>
              {mode === "simulation" && (
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
              )}
              {mode === "diagram" && (
                <span className="ml-auto text-[11px] text-gray-400">{"Toca un camino a la izquierda para ir a su plantilla"}</span>
              )}
            </div>

            {mode === "diagram" ? (
              <div ref={diagramRef} className="relative max-h-[640px] overflow-y-auto p-4">
                {renderDiagramNode(flow.root, rootPath)}
              </div>
            ) : (
              <>
                <div ref={chatRef} className="max-h-[600px] overflow-y-auto p-4" style={{ backgroundColor: "#efeae2" }}>
                  <div className="mx-auto flex w-full max-w-xl flex-col gap-3">
                    {path.map((step, index) => {
                      const parent = index > 0 ? path[index - 1].node : undefined;
                      return (
                        <div key={`${index}-${nodeTitle(step.node)}`} className="flex flex-col gap-2">
                          {parent && parent.kind === "template" && (
                            <div className="self-end rounded-[10px] rounded-br-[3px] bg-white px-3 py-1.5 text-[13px] text-gray-800 shadow-sm">{step.label}</div>
                          )}
                          {parent && parent.kind === "decision" && (
                            <div className="self-center rounded-full bg-amber-100 px-3 py-1 text-[11px] font-medium text-amber-800">{`Caso: ${step.label}`}</div>
                          )}
                          <div className={step.node.kind === "decision" ? "self-center" : ""}>{renderCard(step.node, false)}</div>
                          {(step.node.effects || []).map((effect) => (
                            <div key={effect} className="self-center rounded-full bg-sky-100 px-3 py-1 text-[11px] font-medium text-sky-800">
                              {effect}
                            </div>
                          ))}
                        </div>
                      );
                    })}
                    {(current.branches || []).length === 0 && (
                      <div className="self-center rounded-full bg-gray-200 px-3 py-1 text-[11px] text-gray-600">{"Fin de la conversaci\u00f3n"}</div>
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
                          current.kind === "decision" ? `${pillDecision} hover:bg-amber-50` : `${pillPrimary} hover:bg-violet-50 dark:hover:bg-violet-500/10`
                        }`}
                      >
                        {branch.label}
                      </button>
                    ))}
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
