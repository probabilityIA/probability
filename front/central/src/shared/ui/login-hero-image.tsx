'use client';

const MEDIA_BASE =
    process.env.NEXT_PUBLIC_S3_BASE_URL ||
    'https://probability-media-assets.s3.us-east-1.amazonaws.com';

const ASSET = `${MEDIA_BASE}/public/login/login-hero-v1`;

export function LoginHeroImage() {
    return (
        <div className="absolute inset-0">
            <video
                className="absolute inset-0 w-full h-full object-cover"
                poster={`${ASSET}-poster.jpg`}
                autoPlay
                muted
                loop
                playsInline
                preload="auto"
                aria-hidden="true"
            >
                <source src={`${ASSET}.webm`} type="video/webm" />
                <source src={`${ASSET}.mp4`} type="video/mp4" />
            </video>
        </div>
    );
}
