export const CARRIER_LOGOS: Record<string, string> = {
    'SERVIENTREGA': 'https://i.revistapym.com.co/old/2021/09/WhatsApp-Image-2021-09-25-at-1.08.55-PM.jpeg?w=400&r=1_1',
    'COORDINADORA': 'https://olartemoure.com/wp-content/uploads/2023/05/coordinadora-logo.png',
    'DHLEXPRESS': 'https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png',
    'DHL': 'https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png',
    'FEDEX': 'https://upload.wikimedia.org/wikipedia/commons/thumb/9/9d/FedEx_Express.svg/960px-FedEx_Express.svg.png',
    'INTERRAPIDISIMO': 'https://interrapidisimo.com/wp-content/uploads/Logo-Inter-Rapidisimo-Vv-400x431-1.png',
    '472LOGISTICA': 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTnDF0ozRHf3s5BPqLsr7Vg-X8JRzECvFvwBQ&s',
    'SPEED': 'https://speedcargopa.com/wp-content/uploads/2021/03/Logo-mejorado-transparencia.png',
    'SPEEDCARGO': 'https://speedcargopa.com/wp-content/uploads/2021/03/Logo-mejorado-transparencia.png',
    'ENVIA': 'https://images.seeklogo.com/logo-png/31/1/envia-mensajeria-logo-png_seeklogo-311137.png',
    'PIBOX': 'https://play-lh.googleusercontent.com/r_zPLkaHZK4Odu1yp6dqIdUnVAmIiLc3s18F9gUFqcz8IyHqCb_aGHP4iJSesXxnUyU',
    'TCC': 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Logo_TCC.svg/1280px-Logo_TCC.svg.png',
    'TRANSPORTADORADECARACOLOMBIA': 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Logo_TCC.svg/1280px-Logo_TCC.svg.png',
    '99MINUTOS': 'https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Logo-99minutos.svg/3840px-Logo-99minutos.svg.png',
    'DEPRISA': 'https://www.specialcolombia.com/wp-content/uploads/2023/05/Logo_azul_concepto_azul-deprisa.png',
    'ENVIOCLICK': 'https://www.envioclickpro.com.co/assets/images/envioclick-logo.png',
    'ENVIAME': 'https://enviame.io/wp-content/uploads/2021/01/logo-enviame.svg',
    'MIPAQUETE': 'https://mipaquete.com/wp-content/uploads/2021/03/mipaquete-logo.png',
};

export function getCarrierLogo(carrier?: string | null): string | null {
    if (!carrier) return null;
    const key = carrier.toUpperCase().replace(/[\s\-_]/g, '');
    return CARRIER_LOGOS[key] || null;
}
