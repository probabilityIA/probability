import { Testimonial } from '../../domain/types';

interface TestimonialsSectionProps {
    content: Testimonial[];
}

export function TestimonialsSection({ content }: TestimonialsSectionProps) {
    if (!content || content.length === 0) return null;

    return (
        <section className="py-16 px-4 bg-gray-50">
            <div className="max-w-7xl mx-auto">
                <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">Lo que dicen nuestros clientes</h2>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    {content.map((testimonial, index) => (
                        <div key={index} className="bg-white rounded-xl p-6 shadow-sm">
                            <div className="flex items-center gap-3 mb-4">
                                {testimonial.avatar ? (
                                    <img src={testimonial.avatar} alt={testimonial.name} className="w-12 h-12 rounded-full object-cover" />
                                ) : (
                                    <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 font-bold">
                                        {testimonial.name.charAt(0).toUpperCase()}
                                    </div>
                                )}
                                <div>
                                    <p className="font-semibold text-gray-900">{testimonial.name}</p>
                                    {testimonial.rating && (
                                        <div className="flex gap-0.5">
                                            {Array.from({ length: 5 }).map((_, i) => (
                                                <svg key={i} className={`w-4 h-4 ${i < testimonial.rating! ? 'text-yellow-400' : 'text-gray-200'}`} fill="currentColor" viewBox="0 0 20 20">
                                                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                                                </svg>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            </div>
                            <p className="text-gray-600 italic">&ldquo;{testimonial.text}&rdquo;</p>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
}
