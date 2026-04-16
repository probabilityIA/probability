import { LocationContent } from '../../domain/types';

interface LocationSectionProps {
    content: LocationContent;
}

export function LocationSection({ content }: LocationSectionProps) {
    return (
        <section className="py-16 px-4 max-w-7xl mx-auto">
            <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-8 text-center">Ubicación</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div className="space-y-4">
                    {content.address && (
                        <div className="flex items-start gap-3">
                            <svg className="w-6 h-6 text-gray-400 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                            </svg>
                            <p className="text-gray-600 dark:text-gray-300">{content.address}</p>
                        </div>
                    )}
                    {content.hours && (
                        <div className="flex items-start gap-3">
                            <svg className="w-6 h-6 text-gray-400 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                            <p className="text-gray-600 dark:text-gray-300 whitespace-pre-line">{content.hours}</p>
                        </div>
                    )}
                </div>
                {content.lat && content.lng && (
                    <div className="h-64 bg-gray-100 rounded-xl overflow-hidden">
                        <iframe
                            src={`https://maps.google.com/maps?q=${content.lat},${content.lng}&z=15&output=embed`}
                            className="w-full h-full border-0"
                            loading="lazy"
                            title="Ubicación"
                        />
                    </div>
                )}
            </div>
        </section>
    );
}
