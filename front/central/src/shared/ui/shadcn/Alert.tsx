import React from 'react';

interface Props { type?: 'info' | 'success' | 'error' | 'warning'; children: React.ReactNode }

const Alert: React.FC<Props> = ({ type = 'info', children }) => {
  const base = 'rounded-md p-3 text-sm';
  const styles = {
    info: 'bg-blue-50 text-blue-800 border border-blue-100',
    success: 'bg-green-50 text-green-800 border border-green-100',
    error: 'bg-red-50 text-red-800 border border-red-100',
    warning: 'bg-yellow-50 text-yellow-800 border border-yellow-100',
  } as const;

  return <div className={`${base} ${styles[type]}`}>{children}</div>;
};

export default Alert;
