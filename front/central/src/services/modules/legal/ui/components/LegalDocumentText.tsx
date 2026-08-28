'use client';

import { getLegalSections } from '../../content';

interface Props {
    code: string;
    version: string;
}

export function LegalDocumentText({ code, version }: Props) {
    const secciones = getLegalSections(code, version);

    if (!secciones || secciones.length === 0) {
        return (
            <p className="text-sm text-gray-600 dark:text-gray-300">
                {'El texto de este documento no esta disponible en esta version de la plataforma.'}
            </p>
        );
    }

    return (
        <article>
            {secciones.map((seccion, idx) => (
                <section key={idx} className="mb-5 last:mb-0">
                    {seccion.title && (
                        <h3 className="mb-2 text-sm font-bold text-gray-900 dark:text-white">
                            {seccion.title}
                        </h3>
                    )}
                    {seccion.body.split('\n\n').map((parrafo, i) => (
                        <p
                            key={i}
                            className="mb-2 text-[13px] leading-relaxed text-gray-700 last:mb-0 dark:text-gray-300"
                        >
                            {parrafo}
                        </p>
                    ))}
                </section>
            ))}
        </article>
    );
}
