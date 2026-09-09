"use client";

import { useEffect, useRef, useState } from "react";

export interface StatusOption {
  id: number;
  name: string;
  code: string;
  color?: string;
}

interface StatusMultiSelectProps {
  options: StatusOption[];
  selected: number[];
  onChange: (ids: number[]) => void;
  disabled?: boolean;
}

export function StatusMultiSelect({
  options,
  selected,
  onChange,
  disabled,
}: StatusMultiSelectProps) {
  const [open, setOpen] = useState(false);
  const contenedor = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const alClicAfuera = (e: MouseEvent) => {
      if (contenedor.current && !contenedor.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    const alEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", alClicAfuera);
    document.addEventListener("keydown", alEscape);
    return () => {
      document.removeEventListener("mousedown", alClicAfuera);
      document.removeEventListener("keydown", alEscape);
    };
  }, [open]);

  const alternar = (id: number) => {
    onChange(
      selected.includes(id)
        ? selected.filter((s) => s !== id)
        : [...selected, id],
    );
  };

  const elegidos = options.filter((o) => selected.includes(o.id));

  return (
    <div className="relative" ref={contenedor}>
      <button
        type="button"
        disabled={disabled}
        onClick={() => setOpen((v) => !v)}
        className="flex w-full min-w-[170px] items-center justify-between gap-1 rounded border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-2 py-1 text-left transition-colors hover:border-gray-300 disabled:cursor-not-allowed disabled:opacity-60"
      >
        <span className="flex flex-wrap gap-1">
          {elegidos.length === 0 ? (
            <span className="text-[10px] italic text-gray-400">
              Todos los estados
            </span>
          ) : (
            elegidos.map((o) => (
              <span
                key={o.id}
                className="inline-flex items-center gap-1 rounded-full border border-blue-300 bg-blue-50 px-1.5 py-0.5 text-[10px] font-medium text-blue-700"
              >
                <span
                  className="h-1.5 w-1.5 rounded-full"
                  style={{ backgroundColor: o.color || "#9CA3AF" }}
                />
                {o.name}
              </span>
            ))
          )}
        </span>
        <svg
          className={`h-3 w-3 shrink-0 text-gray-400 transition-transform ${open ? "rotate-180" : ""}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </button>

      {open && (
        <div className="absolute z-30 mt-1 max-h-56 w-max min-w-full overflow-y-auto rounded-md border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 py-1 shadow-lg">
          {options.length === 0 ? (
            <p className="px-3 py-2 text-[11px] text-gray-500 dark:text-gray-400">
              Este evento no tiene estados configurables.
            </p>
          ) : (
            <>
              {options.map((o) => {
                const marcado = selected.includes(o.id);
                return (
                  <button
                    key={o.id}
                    type="button"
                    onClick={() => alternar(o.id)}
                    className="flex w-full items-center gap-2 px-3 py-1.5 text-left text-[11px] text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700"
                  >
                    <span
                      className={`flex h-3.5 w-3.5 shrink-0 items-center justify-center rounded border ${
                        marcado
                          ? "border-blue-500 bg-blue-500"
                          : "border-gray-300 dark:border-gray-600"
                      }`}
                    >
                      {marcado && (
                        <svg
                          className="h-2.5 w-2.5 text-white"
                          fill="currentColor"
                          viewBox="0 0 20 20"
                        >
                          <path
                            fillRule="evenodd"
                            d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                            clipRule="evenodd"
                          />
                        </svg>
                      )}
                    </span>
                    <span
                      className="h-1.5 w-1.5 shrink-0 rounded-full"
                      style={{ backgroundColor: o.color || "#9CA3AF" }}
                    />
                    {o.name}
                  </button>
                );
              })}
              {selected.length > 0 && (
                <button
                  type="button"
                  onClick={() => onChange([])}
                  className="mt-1 w-full border-t border-gray-100 dark:border-gray-700 px-3 py-1.5 text-left text-[10px] text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700"
                >
                  Quitar el filtro (dispara en todos los estados)
                </button>
              )}
            </>
          )}
        </div>
      )}
    </div>
  );
}
