'use client';

import { useState } from 'react';
import { submitContactAction } from '../../infra/actions';
import { ContactContent } from '../../domain/types';

interface ContactSectionProps {
    slug: string;
    content: ContactContent | null;
}

export function ContactSection({ slug, content }: ContactSectionProps) {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [phone, setPhone] = useState('');
    const [message, setMessage] = useState('');
    const [loading, setLoading] = useState(false);
    const [success, setSuccess] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        setSuccess(false);

        const result = await submitContactAction(slug, { name, email, phone, message });

        setLoading(false);
        if (result && 'message' in result && !('success' in result)) {
            setSuccess(true);
            setName('');
            setEmail('');
            setPhone('');
            setMessage('');
        } else {
            setError((result as any)?.message || 'Error al enviar el mensaje');
        }
    };

    return (
        <section className="py-16 px-4 bg-gray-50">
            <div className="max-w-3xl mx-auto">
                <h2 className="text-3xl font-bold text-gray-900 mb-8 text-center">Contacto</h2>

                {content?.contacts && content.contacts.length > 0 && (
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-8">
                        {content.contacts.map((contact, i) => (
                            <div key={i} className="bg-white rounded-lg p-4 shadow-sm">
                                <p className="font-medium text-gray-900">{contact.name}</p>
                                <p className="text-sm text-gray-500">{contact.role}</p>
                                <p className="text-sm text-gray-600 mt-1">{contact.phone}</p>
                            </div>
                        ))}
                    </div>
                )}

                {content?.email && (
                    <p className="text-center text-gray-600 mb-4">
                        Email: <a href={`mailto:${content.email}`} className="text-blue-600 hover:underline">{content.email}</a>
                    </p>
                )}
                {content?.phone && (
                    <p className="text-center text-gray-600 mb-8">
                        Teléfono: <a href={`tel:${content.phone}`} className="text-blue-600 hover:underline">{content.phone}</a>
                    </p>
                )}

                <form onSubmit={handleSubmit} className="bg-white rounded-xl p-6 shadow-sm space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Nombre *</label>
                        <input
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            required
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        />
                    </div>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
                            <input
                                type="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Teléfono</label>
                            <input
                                type="tel"
                                value={phone}
                                onChange={(e) => setPhone(e.target.value)}
                                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            />
                        </div>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Mensaje *</label>
                        <textarea
                            value={message}
                            onChange={(e) => setMessage(e.target.value)}
                            required
                            rows={4}
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                        />
                    </div>

                    {error && <p className="text-red-600 text-sm">{error}</p>}
                    {success && <p className="text-green-600 text-sm">Mensaje enviado correctamente</p>}

                    <button
                        type="submit"
                        disabled={loading}
                        className="w-full py-3 rounded-lg text-white font-medium transition-colors disabled:opacity-50"
                        style={{ backgroundColor: 'var(--brand-secondary)' }}
                    >
                        {loading ? 'Enviando...' : 'Enviar Mensaje'}
                    </button>
                </form>
            </div>
        </section>
    );
}
