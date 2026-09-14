const VALID_PLACEHOLDER = /\{\{\d+\}\}/g;
const LEADING_PLACEHOLDER = /^\{\{\d+\}\}/;
const TRAILING_PLACEHOLDER = /\{\{\d+\}\}$/;
const ADJACENT_PLACEHOLDERS = /\{\{\d+\}\}\s*\{\{\d+\}\}/;
const EMOJI = /[\u{1F000}-\u{1FAFF}\u{2600}-\u{27BF}\u{2B00}-\u{2BFF}\u{FE0F}]/u;
const WORD = /[\p{L}\p{N}]+/gu;

export interface TemplateRuleInput {
  body: string;
  header: string;
  footer: string;
  buttons: string[];
  variables: number;
}

export function metaRuleError({ body, header, footer, buttons, variables }: TemplateRuleInput): string | null {
  const stripped = body.replace(VALID_PLACEHOLDER, " ");

  if (stripped.includes("{{") || stripped.includes("}}")) {
    return "Las variables se escriben con n\u00famero entre llaves dobles: {{1}}, {{2}}";
  }
  if (LEADING_PLACEHOLDER.test(body)) {
    return "El cuerpo no puede empezar con una variable: Meta lo rechaza. Pon texto antes.";
  }
  if (TRAILING_PLACEHOLDER.test(body)) {
    return "El cuerpo no puede terminar con una variable: Meta lo rechaza. Pon texto o un punto despu\u00e9s.";
  }
  if (ADJACENT_PLACEHOLDERS.test(body)) {
    return "Hay dos variables seguidas: Meta exige texto entre una variable y otra.";
  }
  if (variables > 0) {
    const fixedWords = stripped.match(WORD)?.length ?? 0;
    const needed = 2 * variables + 1;
    if (fixedWords < needed) {
      return `El cuerpo tiene demasiadas variables para tan poco texto: con ${variables} variable(s) escribe al menos ${needed} palabras fijas.`;
    }
  }

  if (header) {
    if (header.includes("{{")) return "El encabezado no admite variables.";
    if (/[\n\r]/.test(header)) return "El encabezado no admite saltos de l\u00ednea.";
    if (EMOJI.test(header)) return "El encabezado no admite emojis: Meta lo rechaza.";
    if (/[*_~]/.test(header)) return "El encabezado no admite formato (*, _, ~).";
  }

  if (footer) {
    if (footer.includes("{{")) return "El pie no admite variables.";
    if (/[\n\r]/.test(footer)) return "El pie no admite saltos de l\u00ednea.";
  }

  for (const text of buttons) {
    if (text.includes("{{")) return `El bot\u00f3n "${text}" no admite variables.`;
  }

  return null;
}
