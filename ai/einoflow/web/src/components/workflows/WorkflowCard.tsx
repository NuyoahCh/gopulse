import { Activity, AlertTriangle, ArrowUpRight, Signal, User } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Badge } from '../ui/badge';
import type { Workflow } from '../../types/workflow';
import { cn } from '../../lib/utils';

const statusConfig = {
  healthy: { label: 'Healthy', color: 'text-emerald-600', icon: Activity },
  degraded: { label: 'Degraded', color: 'text-amber-500', icon: Signal },
  failed: { label: 'Failed', color: 'text-rose-600', icon: AlertTriangle }
} as const;

interface WorkflowCardProps {
  workflow: Workflow;
}

const WorkflowCard = ({ workflow }: WorkflowCardProps) => {
  const config = statusConfig[workflow.status];
  const StatusIcon = config.icon;

  return (
    <Card className="flex flex-col justify-between">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>{workflow.name}</CardTitle>
            <p className="mt-1 text-sm text-slate-500">{workflow.description}</p>
          </div>
          <Badge className={cn('gap-1 text-slate-700')}>
            <StatusIcon size={14} className={config.color} />
            {config.label}
          </Badge>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex flex-wrap items-center gap-4 border-b border-dashed border-slate-200 pb-4 text-sm">
          <div className="flex items-center gap-2 text-slate-500">
            <User size={16} />
            Owner <span className="font-medium text-slate-900">{workflow.owner}</span>
          </div>
          <div className="flex items-center gap-2 text-slate-500">
            <Activity size={16} />
            Throughput <span className="font-medium text-slate-900">{workflow.throughput} / min</span>
          </div>
          <div className="flex items-center gap-2 text-slate-500">
            <Signal size={16} />
            Last run <span className="font-medium text-slate-900">{workflow.lastRun}</span>
          </div>
        </div>
        <div className="mt-4 flex flex-wrap items-center justify-between gap-4">
          <div className="flex flex-wrap gap-2">
            {workflow.tags.map((tag) => (
              <Badge key={tag} className="bg-slate-100 text-slate-600">
                {tag}
              </Badge>
            ))}
          </div>
          <button className="inline-flex items-center gap-1 text-sm font-semibold text-brand">
            Inspect
            <ArrowUpRight size={16} />
          </button>
        </div>
      </CardContent>
    </Card>
  );
};

export default WorkflowCard;
