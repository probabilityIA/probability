'use client';

import { useEffect, useState } from 'react';

const MEDIA_BASE =
    process.env.NEXT_PUBLIC_S3_BASE_URL ||
    'https://probability-media-assets.s3.us-east-1.amazonaws.com';

const CLIPS = ['login-hero-hq1', 'login-hero-hq2', 'login-hero-hq3b', 'login-hero-hq4'];

export function LoginHeroImage() {
    const [clip, setClip] = useState<string | null>(null);

    useEffect(() => {
        setClip(CLIPS[Math.floor(Math.random() * CLIPS.length)]);
    }, []);

    if (!clip) return null;

    const asset = `${MEDIA_BASE}/public/login/${clip}`;

    return (
        <video
            key={clip}
            className="absolute inset-0 w-full h-full object-cover"
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
}
