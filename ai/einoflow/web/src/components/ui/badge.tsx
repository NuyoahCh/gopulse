import * as React from 'react';
import { cn } from '../../lib/utils';

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement> {}

const Badge: React.FC<BadgeProps> = ({ className, ...props }) => {
  return (
    <div
      className={cn(
        'inline-flex items-center rounded-full border border-slate-200 bg-slate-50 px-3 py-1 text-xs font-medium text-slate-600',
        className
      )}
      {...props}
    />
  );
};

export { Badge };
