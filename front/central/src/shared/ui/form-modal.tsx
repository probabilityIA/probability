'use client';

import { ReactNode } from 'react';
import { Modal } from './modal';

interface FormModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
  size?: 'sm' | 'md' | 'lg' | 'xl';
}

export function FormModal({ 
  isOpen, 
  onClose, 
  title, 
  children, 
  size = 'lg' 
}: FormModalProps) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} title={title} size={size}>
      <div className="max-h-[80vh] overflow-y-auto">
        {children}
      </div>
    </Modal>
  );
}
