'use client';

import React from 'react';

interface VirtualCardProps {
  balance?: number;
  cardLastFour?: string;
  cardholder?: string;
  isActive?: boolean;
  brand?: string;
  brandTag?: string;
}

export function VirtualCard({
  balance = 107992,
  cardLastFour = '7842',
  cardholder = 'ProbabilityIA',
  isActive = true,
  brand = 'ProbabilityIA',
  brandTag = 'FINTECH'
}: VirtualCardProps) {
  const formatBalance = (amount: number) => {
    return `$${amount.toLocaleString('en-US', { maximumFractionDigits: 0 })}`;
  };

  return (
    <div
      className="relative rounded-[28px] overflow-hidden transition-all duration-300 cursor-pointer hover:shadow-2xl hover:scale-105"
      style={{
        width: '720px',
        height: '420px',
        // Mesh Gradient Background
        backgroundImage: `
          linear-gradient(to bottom right,
            #0A0118 0%, #1A0845 25%, #2D0B6A 33%,
            #120530 20%, #4C1D95 50%, #6D28D9 80%,
            #0D0220 40%, #2E1065 70%, #7C3AED 100%)
        `,
        backgroundSize: '100% 100%',
        backgroundPosition: '0 0',
        boxShadow: `
          0 20px 60px rgba(124, 58, 237, 0.21),
          0 4px 20px rgba(0, 0, 0, 0.31),
          0 0 80px rgba(167, 139, 250, 0.08)
        `,
        border: '1.5px solid',
        borderImageSource: `linear-gradient(135deg,
          rgba(167, 139, 250, 0.31) 0%,
          rgba(255, 255, 255, 0.13) 40%,
          rgba(192, 132, 252, 0.25) 70%,
          rgba(255, 255, 255, 0.06) 100%)`,
        borderImageSlice: 1
      }}
    >
      {/* DECORATIVE BACKGROUND LAYER */}

      {/* Orb Core */}
      <div
        className="absolute pointer-events-none"
        style={{
          width: '130px',
          height: '130px',
          borderRadius: '50%',
          background: '#DDD6FE1C',
          filter: 'blur(22px)',
          left: '550px',
          top: '35px'
        }}
      />

      {/* Orb Mid */}
      <div
        className="absolute pointer-events-none"
        style={{
          width: '220px',
          height: '220px',
          borderRadius: '50%',
          background: '#C4B5FD14',
          filter: 'blur(40px)',
          opacity: 0.9,
          left: '480px',
          top: '30px'
        }}
      />

      {/* Orb Glow */}
      <div
        className="absolute pointer-events-none"
        style={{
          width: '165px',
          height: '165px',
          borderRadius: '50%',
          background: '#A78BFA10',
          filter: 'blur(55px)',
          left: '460px',
          top: '95px'
        }}
      />

      {/* Orb Hot */}
      <div
        className="absolute pointer-events-none"
        style={{
          width: '70px',
          height: '70px',
          borderRadius: '50%',
          background: '#EDE9FE28',
          filter: 'blur(14px)',
          left: '590px',
          top: '55px'
        }}
      />

      {/* Glow Bottom */}
      <div
        className="absolute pointer-events-none"
        style={{
          width: '350px',
          height: '200px',
          borderRadius: '50%',
          background: '#7C3AED0A',
          filter: 'blur(60px)',
          opacity: 0.6,
          left: '-80px',
          top: '300px'
        }}
      />

      {/* Glow Left */}
      <div
        className="absolute pointer-events-none"
        style={{
          width: '287px',
          height: '237px',
          borderRadius: '50%',
          background: '#6D28D90A',
          filter: 'blur(55px)',
          opacity: 0.5,
          left: '-100px',
          top: '113px'
        }}
      />

      {/* Decorative Rings */}
      <div
        className="absolute pointer-events-none"
        style={{
          width: '260px',
          height: '260px',
          borderRadius: '50%',
          border: '0.5px solid rgba(255, 255, 255, 0.06)',
          opacity: 0.6,
          left: '490px',
          top: '-30px'
        }}
      />

      <div
        className="absolute pointer-events-none"
        style={{
          width: '180px',
          height: '180px',
          borderRadius: '50%',
          border: '0.5px solid rgba(255, 255, 255, 0.08)',
          opacity: 0.5,
          left: '530px',
          top: '10px'
        }}
      />

      <div
        className="absolute pointer-events-none"
        style={{
          width: '110px',
          height: '110px',
          borderRadius: '50%',
          border: '0.5px solid rgba(255, 255, 255, 0.12)',
          opacity: 0.4,
          left: '565px',
          top: '45px'
        }}
      />

      {/* Accent Dots */}
      <div
        className="absolute pointer-events-none rounded-full"
        style={{
          width: '4px',
          height: '4px',
          background: '#C4B5FD40',
          left: '650px',
          top: '150px'
        }}
      />

      <div
        className="absolute pointer-events-none rounded-full"
        style={{
          width: '3px',
          height: '3px',
          background: '#C4B5FD30',
          left: '620px',
          top: '180px'
        }}
      />

      <div
        className="absolute pointer-events-none rounded-full"
        style={{
          width: '5px',
          height: '5px',
          background: '#C4B5FD25',
          left: '670px',
          top: '200px'
        }}
      />

      <div
        className="absolute pointer-events-none rounded-full"
        style={{
          width: '3px',
          height: '3px',
          background: '#C4B5FD20',
          left: '640px',
          top: '230px'
        }}
      />

      {/* CONTENT LAYER */}

      {/* Top Bar - Logo */}
      <div className="absolute top-0 left-0 right-0 h-14 flex items-center justify-between px-9">
        <div className="flex items-center gap-2.5">
          {/* Logo Icon */}
          <div className="relative w-7 h-7">
            <div
              className="absolute"
              style={{
                width: '8px',
                height: '22px',
                background: '#C4B5FD',
                borderRadius: '4px',
                left: '3px',
                top: '3px',
                transform: 'rotate(-12deg)',
                transformOrigin: 'center'
              }}
            />
            <div
              className="absolute"
              style={{
                width: '8px',
                height: '16px',
                background: '#8B5CF6',
                borderRadius: '4px',
                right: '0px',
                top: '8px',
                transform: 'rotate(12deg)',
                transformOrigin: 'center'
              }}
            />
          </div>

          {/* Logo Text */}
          <span
            style={{
              fontFamily: "'DM Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
              fontSize: '15px',
              fontWeight: 700,
              color: '#FFFFFFCC',
              letterSpacing: '0.5px'
            }}
          >
            Probability
          </span>
        </div>
      </div>

      {/* Balance Area */}
      <div
        className="absolute flex flex-col gap-1"
        style={{
          left: '36px',
          top: '58px',
          width: '500px'
        }}
      >
        <p
          style={{
            fontFamily: "'DM Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
            fontSize: '10px',
            fontWeight: 600,
            color: '#C4B5FDA0',
            letterSpacing: '4px',
            textTransform: 'uppercase',
            margin: 0
          }}
        >
          SALDO DISPONIBLE
        </p>

        <h2
          style={{
            fontFamily: "'DM Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
            fontSize: '52px',
            fontWeight: 800,
            color: '#FFFFFF',
            letterSpacing: '-2px',
            margin: 0,
            lineHeight: '1'
          }}
        >
          {formatBalance(balance)}
        </h2>
      </div>

      {/* Chip Area */}
      <div
        className="absolute flex items-center gap-3"
        style={{
          left: '563px',
          top: '95px'
        }}
      >
        {/* Contactless Icon */}
        <svg
          width="18"
          height="18"
          viewBox="0 0 0 24"
          fill="none"
          style={{
            color: '#FFFFFF40',
            transform: 'rotate(90deg)'
          }}
        >
          <path
            d="M1 9l2 2c4.97-4.97 13.03-4.97 18 0l2-2M9 17l3 3 3-3M5 13l2 2c2.76-2.76 7.24-2.76 10 0l2-2"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>

        {/* Chip */}
        <div
          style={{
            width: '48px',
            height: '34px',
            borderRadius: '7px',
            backgroundImage: 'linear-gradient(145deg, #F5D060 0%, #D4A843 50%, #C49A38 100%)',
            border: '0.5px solid rgba(245, 208, 96, 0.38)',
            boxShadow: 'inset 0 -1px 2px rgba(0, 0, 0, 0.1)'
          }}
        />
      </div>

      {/* Status Indicator */}
      <div className="absolute flex items-center gap-3" style={{ left: '36px', top: '195px' }}>
        <div
          style={{
            width: '7px',
            height: '7px',
            borderRadius: '50%',
            background: isActive ? '#34D399' : '#9CA3AF',
            boxShadow: isActive ? '0 0 8px rgba(52, 211, 153, 0.6)' : 'none'
          }}
        />
        <span
          style={{
            fontFamily: "'DM Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
            fontSize: '10px',
            fontWeight: 500,
            color: isActive ? '#34D399A0' : '#9CA3AF',
            letterSpacing: '1px'
          }}
        >
          {isActive ? 'Activa' : 'Inactiva'}
        </span>
      </div>

      {/* Bottom Area - Card Data */}
      <div
        className="absolute left-0 right-0 flex items-end justify-between px-9"
        style={{
          bottom: '30px',
          height: '60px'
        }}
      >
        {/* Left: Card Number & Wallet Label */}
        <div className="flex flex-col gap-2">
          <div className="flex items-center gap-4">
            {[0, 1, 2].map((i) => (
              <span
                key={`dots-${i}`}
                style={{
                  fontFamily: "'DM Sans', monospace",
                  fontSize: '16px',
                  fontWeight: 600,
                  color: '#FFFFFF30',
                  letterSpacing: '2px'
                }}
              >
                ••••
              </span>
            ))}
            <span
              style={{
                fontFamily: "'DM Sans', monospace",
                fontSize: '16px',
                fontWeight: 600,
                color: '#FFFFFF70',
                letterSpacing: '2px'
              }}
            >
              {cardLastFour}
            </span>
          </div>

          <p
            style={{
              fontFamily: "'DM Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
              fontSize: '9px',
              fontWeight: 600,
              color: '#FFFFFF30',
              letterSpacing: '3px',
              textTransform: 'uppercase',
              margin: 0
            }}
          >
            BILLETERA EMPRESARIAL
          </p>
        </div>

        {/* Right: Brand Info */}
        <div className="flex flex-col items-end gap-0.5">
          <p
            style={{
              fontFamily: "'DM Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
              fontSize: '20px',
              fontWeight: 800,
              color: '#FFFFFF',
              letterSpacing: '-0.5px',
              margin: 0
            }}
          >
            {brand}
          </p>

          <p
            style={{
              fontFamily: "'DM Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
              fontSize: '8px',
              fontWeight: 600,
              color: '#C4B5FD60',
              letterSpacing: '4px',
              textTransform: 'uppercase',
              margin: 0
            }}
          >
            {brandTag}
          </p>
        </div>
      </div>
    </div>
  );
}
