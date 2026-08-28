import type { LegalSection } from './types';
import { terminosV1 } from './terminos-v1';
import { politicaDatosV1 } from './politica-datos-v1';

const CONTENIDO: Record<string, LegalSection[]> = {
    'terms_of_service:1.0': terminosV1,
    'privacy_policy:1.0': politicaDatosV1,
};

export function getLegalSections(code: string, version: string): LegalSection[] | undefined {
    return CONTENIDO[`${code}:${version}`];
}

export type { LegalSection };
