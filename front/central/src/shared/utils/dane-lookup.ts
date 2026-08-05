import danes from "@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json";

const ISO_STATE_TO_NAME: Record<string, string> = {
    "CO-AMA": "AMAZONAS",
    "CO-ANT": "ANTIOQUIA",
    "CO-ARA": "ARAUCA",
    "CO-ATL": "ATLANTICO",
    "CO-BOL": "BOLIVAR",
    "CO-BOY": "BOYACA",
    "CO-CAL": "CALDAS",
    "CO-CAQ": "CAQUETA",
    "CO-CAS": "CASANARE",
    "CO-CAU": "CAUCA",
    "CO-CES": "CESAR",
    "CO-CHO": "CHOCO",
    "CO-COR": "CORDOBA",
    "CO-CUN": "CUNDINAMARCA",
    "CO-DC": "BOGOTA",
    "CO-GUA": "GUAINIA",
    "CO-GUV": "GUAVIARE",
    "CO-HUI": "HUILA",
    "CO-LAG": "LA GUAJIRA",
    "CO-MAG": "MAGDALENA",
    "CO-MET": "META",
    "CO-NAR": "NARINO",
    "CO-NSA": "NORTE DE SANTANDER",
    "CO-PUT": "PUTUMAYO",
    "CO-QUI": "QUINDIO",
    "CO-RIS": "RISARALDA",
    "CO-SAN": "SANTANDER",
    "CO-SAP": "SAN ANDRES Y PROVIDENCIA",
    "CO-SUC": "SUCRE",
    "CO-TOL": "TOLIMA",
    "CO-VAC": "VALLE DEL CAUCA",
    "CO-VAU": "VAUPES",
    "CO-VID": "VICHADA",
};

export const normalizeLocationName = (str: string) => {
    if (!str) return "";
    let s = str.normalize("NFD").replace(/[̀-ͯ]/g, "").toUpperCase().trim();
    s = s.replace(/,\s*D\.C\./g, "").replace(/\sD\.C\./g, "").replace(/\sDC\b/g, "").trim();
    return s;
};

const normalizeState = (state: string) => {
    const upper = (state || "").trim().toUpperCase();
    if (ISO_STATE_TO_NAME[upper]) return normalizeLocationName(ISO_STATE_TO_NAME[upper]);
    return normalizeLocationName(state);
};

export const findDaneCode = (city: string, state: string) => {
    const targetCity = normalizeLocationName(city);
    const targetState = normalizeState(state);
    if (!targetCity) return null;

    const entries = Object.entries(danes);

    const exactMatch = entries.find(([_, data]: [string, any]) => {
        const dCity = normalizeLocationName(data.ciudad);
        const dState = normalizeLocationName(data.departamento);
        return dCity === targetCity && dState === targetState;
    });
    if (exactMatch) return exactMatch[0];

    const cityMatch = entries.find(([_, data]: [string, any]) => {
        const dCity = normalizeLocationName(data.ciudad);
        return dCity === targetCity;
    });
    if (cityMatch) return cityMatch[0];

    const partialMatch = entries.find(([_, data]: [string, any]) => {
        const dCity = normalizeLocationName(data.ciudad);
        return dCity.includes(targetCity) || targetCity.includes(dCity);
    });
    if (partialMatch) return partialMatch[0];

    return null;
};

type DaneEntry = { ciudad: string; departamento: string };
const daneMap = danes as Record<string, DaneEntry>;

export const resolveCityState = (
    city: string,
    state: string,
    verifiedDaneCode?: string | null
): DaneEntry | null => {
    const verifiedKey = verifiedDaneCode ? `${verifiedDaneCode}000` : null;
    if (verifiedKey && daneMap[verifiedKey]) return daneMap[verifiedKey];

    const mappedCode = findDaneCode(city, state);
    if (mappedCode && daneMap[mappedCode]) return daneMap[mappedCode];

    return null;
};
