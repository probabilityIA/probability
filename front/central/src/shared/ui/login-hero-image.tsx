'use client';

import { useEffect, useRef } from 'react';

const MEDIA_BASE =
    process.env.NEXT_PUBLIC_S3_BASE_URL ||
    'https://probability-media-assets.s3.us-east-1.amazonaws.com';

export const LOGIN_CLIPS = [
    'login-hero-hq1',
    'login-hero-hq2',
    'login-hero-hq3',
    'login-hero-hq4',
];

interface Props {
    index: number;
}

export function LoginHeroImage({ index }: Props) {
    const refs = useRef<Record<string, HTMLVideoElement | null>>({});

    useEffect(() => {
        Object.entries(refs.current).forEach(([clip, el]) => {
            if (!el) return;
            const i = LOGIN_CLIPS.indexOf(clip);
            if (i === index || i === (index + 1) % LOGIN_CLIPS.length) {
                const p = el.play();
                if (p && typeof p.catch === 'function') p.catch(() => undefined);
            } else {
                el.pause();
            }
        });
    }, [index]);

    return (
        <div className="absolute inset-0 overflow-hidden">
            {LOGIN_CLIPS.map((clip, i) => {
                const active = i === index;
                const next = i === (index + 1) % LOGIN_CLIPS.length;
                if (!active && !next) return null;
                const asset = `${MEDIA_BASE}/public/login/${clip}`;
                return (
                    <video
                        key={clip}
                        ref={(el) => { refs.current[clip] = el; }}
                        className="absolute w-full h-full object-cover transition-opacity duration-1000"
                        style={{
                            opacity: active ? 1 : 0,
                            filter: 'blur(6px) saturate(112%)',
                            transform: 'scale(1.06)',
                        }}
                        poster={`${asset}-poster.jpg`}
                        autoPlay
                        muted
                        loop
                        playsInline
                        preload="auto"
                        aria-hidden="true"
                    >
                        <source src={`${asset}.webm`} type="video/webm" />
                        <source src={`${asset}.mp4`} type="video/mp4" />
                    </video>
                );
            })}

            <div
                className="absolute inset-0"
                style={{ background: 'linear-gradient(100deg, rgba(20,12,45,0.74) 0%, rgba(20,12,45,0.5) 45%, rgba(20,12,45,0.3) 100%)' }}
            />
        </div>
    );
}
