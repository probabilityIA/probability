'use client';

import { useCallback, useEffect, useState } from 'react';
import { Button } from '@/shared/ui/button';
import { Spinner } from '@/shared/ui/spinner';
import { acceptLegalDocumentsAction, getPendingLegalDocumentsAction } from '../../infra/actions';
import type { LegalDocument } from '../../domain/types';
import { LegalDocumentText } from './LegalDocumentText';

export function LegalAcceptanceGate() {
    const [documentos, setDocumentos] = useState<LegalDocument[]>([]);
    const [cargando, setCargando] = useState(true);
    const [enviando, setEnviando] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [aceptado, setAceptado] = useState(false);
    const [activo, setActivo] = useState<string | null>(null);

    const consultarPendientes = useCallback(async () => {
        try {
            const resultado = await getPendingLegalDocumentsAction();
            if (resultado.requires_acceptance && resultado.documents.length > 0) {
                setDocumentos(resultado.documents);
                setActivo(resultado.documents[0].code);
            } else {
                setDocumentos([]);
            }
        } catch {
            setDocumentos([]);
        } finally {
            setCargando(false);
        }
    }, []);

    useEffect(() => {
        consultarPendientes();
    }, [consultarPendientes]);

    if (cargando || documentos.length === 0) return null;

    const documentoActivo = documentos.find((d) => d.code === activo) ?? documentos[0];

    const confirmar = async () => {
        setEnviando(true);
        setError(null);
        try {
            await acceptLegalDocumentsAction(documentos.map((d) => d.id));
            setDocumentos([]);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'No se pudo registrar la aceptación');
            setEnviando(false);
        }
    };

    return (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-gray-900/70 backdrop-blur-sm p-4">
            <div className="w-full max-w-4xl max-h-[92vh] flex flex-col bg-white dark:bg-gray-800 rounded-lg shadow-xl">
                <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                    <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                        {'Términos legales de Probability'}
                    </h2>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mt-1">
                        {'Para continuar usando la plataforma debes leer y aceptar los siguientes documentos.'}
                    </p>
                </div>

                <div className="px-6 pt-4 flex flex-wrap gap-2">
                    {documentos.map((doc) => (
                        <button
                            key={doc.id}
                            type="button"
                            onClick={() => setActivo(doc.code)}
                            className={`px-3 py-1.5 rounded text-xs font-medium border transition-colors ${
                                doc.code === documentoActivo.code
                                    ? 'bg-[var(--color-primary)] text-white border-transparent'
                                    : 'bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-200 border-gray-300 dark:border-gray-600'
                            }`}
                        >
                            {doc.title} {'v'}{doc.version}
                        </button>
                    ))}
                </div>

                <div className="flex-1 min-h-0 px-6 py-4">
                    <div className="h-[52vh] overflow-y-auto rounded border border-gray-200 bg-gray-50 px-5 py-4 dark:border-gray-700 dark:bg-gray-900/40">
                        <LegalDocumentText contentHtml={documentoActivo.content_html} />
                    </div>
                </div>

                <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 space-y-3">
                    <label className="flex items-start gap-3 text-sm text-gray-700 dark:text-gray-200 cursor-pointer">
                        <input
                            type="checkbox"
                            checked={aceptado}
                            onChange={(e) => setAceptado(e.target.checked)}
                            disabled={enviando}
                            className="mt-0.5 h-4 w-4"
                        />
                        <span>
                            {'Declaro que leí y acepto '}
                            {documentos.map((doc, idx) => (
                                <span key={doc.id}>
                                    {idx > 0 ? ' y ' : ''}
                                    <span className="font-medium">{doc.title}</span>
                                    {' (v'}{doc.version}{')'}
                                </span>
                            ))}
                            {'. Esta aceptación queda registrada con mi usuario, fecha y hora.'}
                        </span>
                    </label>

                    {error && <p className="text-sm text-red-600">{error}</p>}

                    <div className="flex justify-end">
                        <Button onClick={confirmar} disabled={!aceptado || enviando}>
                            {enviando ? <Spinner /> : 'Acepto y continuo'}
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default LegalAcceptanceGate;
