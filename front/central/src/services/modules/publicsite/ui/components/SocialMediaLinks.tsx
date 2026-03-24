import { SocialMediaContent } from '../../domain/types';

interface SocialMediaLinksProps {
    content: SocialMediaContent;
}

export function SocialMediaLinks({ content }: SocialMediaLinksProps) {
    const links = [
        { url: content.facebook, label: 'Facebook', icon: 'M18 2h-3a5 5 0 00-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 011-1h3z' },
        { url: content.instagram, label: 'Instagram', icon: 'M16 4H8a4 4 0 00-4 4v8a4 4 0 004 4h8a4 4 0 004-4V8a4 4 0 00-4-4zm-4 11a3 3 0 110-6 3 3 0 010 6zm3.5-6.5a1 1 0 110-2 1 1 0 010 2z' },
        { url: content.twitter, label: 'Twitter', icon: 'M23 3a10.9 10.9 0 01-3.14 1.53 4.48 4.48 0 00-7.86 3v1A10.66 10.66 0 013 4s-4 9 5 13a11.64 11.64 0 01-7 2c9 5 20 0 20-11.5a4.5 4.5 0 00-.08-.83A7.72 7.72 0 0023 3z' },
        { url: content.tiktok, label: 'TikTok', icon: 'M9 12a4 4 0 104 4V4a5 5 0 005 5' },
    ].filter(l => l.url);

    if (links.length === 0) return null;

    return (
        <section className="py-8 px-4">
            <div className="max-w-7xl mx-auto text-center">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Síguenos</h3>
                <div className="flex justify-center gap-4">
                    {links.map((link) => (
                        <a
                            key={link.label}
                            href={link.url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="w-10 h-10 rounded-full bg-gray-100 flex items-center justify-center text-gray-600 dark:text-gray-300 hover:bg-gray-200 transition-colors"
                            title={link.label}
                        >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" d={link.icon} />
                            </svg>
                        </a>
                    ))}
                </div>
            </div>
        </section>
    );
}
