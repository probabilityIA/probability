/**
 * Barrel de componentes UI compartidos
 */

export * from './alert';
export * from './avatar-upload';
export * from './badge';
export * from './button';
export * from './filters';
export * from './dynamic-filters';
export * from './confirm-modal';
export * from './date-picker';
export * from './date-range-picker';
export * from './file-input';
export * from './form-modal';
export * from './input';
export * from './modal';
export * from './select';
export * from './sidebar';
export * from './spinner';
export * from './table';
export * from './iam-sidebar';
export * from './orders-sidebar';
export * from './user-profile-modal';

// Re-exportar tipos útiles
export type { 
  TableColumn, 
  PaginationProps, 
  TableFiltersProps,
  FilterOption,
  ActiveFilter
} from './table';

