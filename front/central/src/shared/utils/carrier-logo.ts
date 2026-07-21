const carrierLogos: { [key: string]: string } = {
    'SERVIENTREGA': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_servientrega.png',
    'COORDINADORA': 'https://images-cam93.s3.us-east-1.amazonaws.com/imagen_coordinadora.png',
    'DHLEXPRESS': 'https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png',
    'DHL': 'https://logodownload.org/wp-content/uploads/2015/12/dhl-logo-2.png',
    'FEDEX': 'https://upload.wikimedia.org/wikipedia/commons/thumb/9/9d/FedEx_Express.svg/960px-FedEx_Express.svg.png',
    'INTERRAPIDISIMO': 'https://probability-media-assets.s3.us-east-1.amazonaws.com/carriers/interrapidisimo.jpg',
    '472LOGISTICA': 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTnDF0ozRHf3s5BPqLsr7Vg-X8JRzECvFvwBQ&s',
    'SPEED': 'https://speedcargopa.com/wp-content/uploads/2021/03/Logo-mejorado-transparencia.png',
    'SPEEDCARGO': 'https://speedcargopa.com/wp-content/uploads/2021/03/Logo-mejorado-transparencia.png',
    'ENVIA': 'https://images.seeklogo.com/logo-png/31/1/envia-mensajeria-logo-png_seeklogo-311137.png',
    'PIBOX': 'https://play-lh.googleusercontent.com/r_zPLkaHZK4Odu1yp6dqIdUnVAmIiLc3s18F9gUFqcz8IyHqCb_aGHP4iJSesXxnUyU',
    'TCC': 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Logo_TCC.svg/1280px-Logo_TCC.svg.png',
    'TRANSPORTADORADECARACOLOMBIA': 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Logo_TCC.svg/1280px-Logo_TCC.svg.png',
    '99MINUTOS': 'https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Logo-99minutos.svg/3840px-Logo-99minutos.svg.png',
    'DEPRISA': 'https://www.specialcolombia.com/wp-content/uploads/2023/05/Logo_azul_concepto_azul-deprisa.png',
    'MELONN': 'https://cdn.prod.website-files.com/63dd1adcb3e9e11ce6deef7d/6420fdec2bf3e7cef8db96b1_logo-melonn.svg',
};

export function normalizeCarrierKey(carrierName: string): string {
    if (!carrierName) return '';
    return carrierName
        .normalize('NFD')
        .replace(/[̀-ͯ]/g, '')
        .replace(/[^a-zA-Z0-9]/g, '')
        .toUpperCase();
}

export function getCarrierLogo(carrierName: string): string | null {
    const key = normalizeCarrierKey(carrierName);
    return carrierLogos[key] || null;
}
