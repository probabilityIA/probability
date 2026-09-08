"use client";

const RULES: Array<{ title: string; detail: string }> = [
  {
    title: "Sale desde tu propio número",
    detail:
      "Una campaña de marketing solo se lanza si el negocio tiene su propio número de WhatsApp conectado. El número compartido de Probability no se usa para envíos masivos: si alguien reporta los mensajes, cae la calidad de esa línea y afecta a todos los negocios.",
  },
  {
    title: "Solo a quien acepta marketing",
    detail:
      "Los clientes que se dieron de baja quedan fuera de la audiencia siempre, sin excepción y sin forma de saltarlo desde el panel. Si alguien responde pidiendo no recibir más, queda excluido de todas tus campañas.",
  },
  {
    title: "Sale por tandas, no de golpe",
    detail:
      "Meta arranca los números nuevos en 1.000 conversaciones por día y sube el límite solo si la calidad se mantiene. Mandar todo junto el primer día no acelera nada: penaliza el número. Por eso la campaña envía de a tandas dentro de la franja horaria que elijas.",
  },
  {
    title: "La plantilla la aprueba Meta, no nosotros",
    detail:
      "Una campaña no se lanza con una plantilla en revisión o rechazada. Meta suele responder en minutos, pero puede tardar hasta 24 horas. Prepará la plantilla antes del día que querés salir.",
  },
  {
    title: "Marketing se cobra por mensaje",
    detail:
      "A diferencia de las notificaciones de una orden, cada mensaje de marketing tiene costo. Antes de lanzar ves a cuánta gente le vas a escribir, así sabés qué estás gastando.",
  },
  {
    title: "El chat libre dura 24 horas",
    detail:
      "La plantilla abre la conversación. Recién cuando el cliente responde podés escribirle texto libre, y esa ventana dura 24 horas desde su último mensaje. Si nunca responde, el único camino es otra plantilla.",
  },
];

export function CampaignRules() {
  return (
    <div className="rounded-lg border border-amber-200 bg-amber-50 p-4 dark:border-amber-900/40 dark:bg-amber-900/10">
      <h4 className="mb-3 text-sm font-semibold text-amber-900 dark:text-amber-200">
        {"Cómo funcionan las campañas, y qué límites tienen"}
      </h4>
      <ul className="space-y-2.5">
        {RULES.map((rule) => (
          <li key={rule.title}>
            <p className="text-xs font-medium text-amber-900 dark:text-amber-200">
              {rule.title}
            </p>
            <p className="text-xs leading-relaxed text-amber-800/80 dark:text-amber-200/70">
              {rule.detail}
            </p>
          </li>
        ))}
      </ul>
    </div>
  );
}
