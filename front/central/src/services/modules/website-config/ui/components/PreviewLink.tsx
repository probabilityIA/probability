'use client';

interface PreviewLinkProps {
    businessCode: string;
}

export function PreviewLink({ businessCode }: PreviewLinkProps) {
    if (!businessCode) return null;

    const url = `/tienda/${businessCode}`;

    return (
        <div className="bg-green-50 border border-green-200 rounded-lg p-4 flex items-center justify-between">
            <div>
                <p className="font-medium text-green-800">Tu sitio web</p>
                <p className="text-sm text-green-600">{url}</p>
            </div>
            <a
                href={url}
                target="_blank"
                rel="noopener noreferrer"
                className="px-4 py-2 bg-green-600 text-white rounded-lg text-sm font-medium hover:bg-green-700 transition-colors"
            >
                Previsualizar sitio
            </a>
        </div>
    );
}
