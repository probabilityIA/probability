'use client';

interface Props {
    contentHtml: string;
}

export function LegalDocumentText({ contentHtml }: Props) {
    if (!contentHtml) {
        return (
            <p className="text-sm text-gray-600 dark:text-gray-300">
                {'El texto de este documento no esta disponible.'}
            </p>
        );
    }

    return (
        <article
            className="legal-document text-[13px] leading-relaxed text-gray-700 dark:text-gray-300"
            dangerouslySetInnerHTML={{ __html: contentHtml }}
        />
    );
}
