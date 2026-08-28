'use client';

import { useEffect, useState } from 'react';
import { Spinner } from '@/shared/ui/spinner';

interface Props {
    fileUrl: string;
}

function textUrl(fileUrl: string): string {
    return fileUrl.replace(/\.pdf$/i, '.md');
}

export function LegalDocumentText({ fileUrl }: Props) {
    const [texto, setTexto] = useState<string | null>(null);
    const [error, setError] = useState(false);

    useEffect(() => {
        let cancelado = false;
        setTexto(null);
        setError(false);

        fetch(textUrl(fileUrl), { cache: 'force-cache' })
            .then((res) => {
                if (!res.ok) throw new Error('no disponible');
                return res.text();
            })
            .then((contenido) => {
                if (!cancelado) setTexto(contenido);
            })
            .catch(() => {
                if (!cancelado) setError(true);
            });

        return () => {
            cancelado = true;
        };
    }, [fileUrl]);

    if (error) {
        return (
            <div className="p-4 text-sm text-gray-700 dark:text-gray-200">
                {'No se pudo cargar el documento. '}
                <a
                    href={fileUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-blue-600 underline"
                >
                    {'Abrirlo en otra pestana'}
                </a>
            </div>
        );
    }

    if (texto === null) {
        return (
            <div className="flex h-full items-center justify-center py-12">
                <Spinner size="md" />
            </div>
        );
    }

    const bloques = texto.split('\n').filter((linea) => linea.trim().length > 0);

    return (
        <article className="space-y-3">
            {bloques.map((linea, idx) => {
                if (linea.startsWith('## ')) {
                    return (
                        <h3
                            key={idx}
                            className="pt-3 text-sm font-bold text-gray-900 dark:text-white first:pt-0"
                        >
                            {linea.slice(3)}
                        </h3>
                    );
                }
                return (
                    <p key={idx} className="text-[13px] leading-relaxed text-gray-700 dark:text-gray-300">
                        {linea}
                    </p>
                );
            })}
        </article>
    );
}
