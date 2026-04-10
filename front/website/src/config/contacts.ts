export const WHATSAPP = {
  phoneNumber: '573138241302',
  baseUrl: 'https://wa.me',
  messages: {
    home: 'Hola, me interesa conocer más sobre ProbabilityIA',
    producto: 'Quiero saber más detalles sobre las características del producto',
    tracking: 'Me gustaría información sobre el sistema de rastreo',
    precios: 'Necesito información sobre los planes y precios',
    integraciones: 'Quiero conocer las integraciones disponibles',
    contacto: 'Me gustaría hablar con alguien del equipo',
  }
};

export const DEMO = {
  url: 'https://calendly.com/probabilityia/demo',
};

export function generateWhatsAppUrl(message: string = ''): string {
  const encodedMessage = encodeURIComponent(message);
  return `${WHATSAPP.baseUrl}/${WHATSAPP.phoneNumber}?text=${encodedMessage}`;
}
