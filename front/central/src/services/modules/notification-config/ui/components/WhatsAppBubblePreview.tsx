"use client";

interface WhatsAppBubblePreviewProps {
  header: string;
  body: string;
  footer: string;
  buttons: string[];
  sampleValues: Record<number, string>;
}

const PLACEHOLDER = /\{\{(\d+)\}\}/g;

function fill(text: string, samples: Record<number, string>): string {
  return text.replace(PLACEHOLDER, (match, index) => {
    const value = samples[Number(index)];
    return value ? value : match;
  });
}

export function WhatsAppBubblePreview({
  header,
  body,
  footer,
  buttons,
  sampleValues,
}: WhatsAppBubblePreviewProps) {
  const filledBody = fill(body, sampleValues);
  const now = new Date().toLocaleTimeString("es-CO", {
    hour: "2-digit",
    minute: "2-digit",
  });

  return (
    <div className="rounded-xl bg-[#ece5dd] p-4 dark:bg-[#0b141a]">
      <p className="mb-2 text-center text-[11px] text-gray-500 dark:text-gray-400">
        {"Así lo va a ver el cliente"}
      </p>

      <div className="ml-auto max-w-[85%] rounded-lg rounded-tr-none bg-[#d9fdd3] p-3 shadow-sm dark:bg-[#005c4b]">
        {header && (
          <p className="mb-1 text-sm font-bold text-gray-900 dark:text-white">
            {fill(header, sampleValues)}
          </p>
        )}

        {filledBody ? (
          <p className="whitespace-pre-wrap break-words text-sm text-gray-800 dark:text-gray-100">
            {filledBody}
          </p>
        ) : (
          <p className="text-sm italic text-gray-400">
            {"Escribí el mensaje para ver cómo queda"}
          </p>
        )}

        {footer && (
          <p className="mt-2 text-[11px] text-gray-500 dark:text-gray-300">
            {fill(footer, sampleValues)}
          </p>
        )}

        <p className="mt-1 text-right text-[10px] text-gray-500 dark:text-gray-300">
          {now}
        </p>
      </div>

      {buttons.length > 0 && (
        <div className="ml-auto mt-1 max-w-[85%] space-y-1">
          {buttons.map((label) => (
            <div
              key={label}
              className="rounded-lg bg-white py-2 text-center text-sm font-medium text-[#00a5f4] shadow-sm dark:bg-[#1f2c34]"
            >
              {label}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
