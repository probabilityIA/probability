import React from 'react';

interface Props { size?: 'sm' | 'md' | 'lg' }

const Spinner: React.FC<Props> = ({ size = 'md' }) => {
  const s = size === 'sm' ? 16 : size === 'lg' ? 32 : 24;
  return (
    <svg className="animate-spin" width={s} height={s} viewBox="0 0 24 24" fill="none">
      <circle cx="12" cy="12" r="10" stroke="#E6EEF8" strokeWidth="4" />
      <path d="M22 12a10 10 0 00-10-10" stroke="#0F172A" strokeWidth="4" strokeLinecap="round" />
    </svg>
  );
};

export default Spinner;
