import { AboutContent } from '../../domain/types';

interface AboutSectionProps {
    content: AboutContent;
}

export function AboutSection({ content }: AboutSectionProps) {
    return (
        <section className="py-16 px-4 max-w-7xl mx-auto">
            <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">Sobre Nosotros</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
                <div>
                    {content.text && <p className="text-gray-600 text-lg leading-relaxed mb-6">{content.text}</p>}
                    {content.mission && (
                        <div className="mb-4">
                            <h3 className="font-semibold text-gray-900 mb-2">Misión</h3>
                            <p className="text-gray-600">{content.mission}</p>
                        </div>
                    )}
                    {content.vision && (
                        <div>
                            <h3 className="font-semibold text-gray-900 mb-2">Visión</h3>
                            <p className="text-gray-600">{content.vision}</p>
                        </div>
                    )}
                </div>
                {content.image && (
                    <div className="rounded-xl overflow-hidden">
                        <img src={content.image} alt="Sobre nosotros" className="w-full h-64 md:h-80 object-cover" />
                    </div>
                )}
            </div>
        </section>
    );
}
