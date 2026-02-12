import { LucideIcon } from 'lucide-react';
import { Card, CardContent } from '../ui/card';

interface MetricCardProps {
  label: string;
  value: string;
  trend: string;
  icon: LucideIcon;
  accent?: string;
}

const MetricCard = ({ label, value, trend, icon: Icon, accent = 'bg-blue-100 text-blue-700' }: MetricCardProps) => {
  return (
    <Card>
      <CardContent className="flex items-center justify-between">
        <div>
          <p className="text-sm text-slate-500">{label}</p>
          <p className="mt-2 text-3xl font-semibold text-slate-900">{value}</p>
          <p className="mt-1 text-sm text-emerald-600">{trend}</p>
        </div>
        <div className={`flex h-12 w-12 items-center justify-center rounded-2xl ${accent}`}>
          <Icon size={24} />
        </div>
      </CardContent>
    </Card>
  );
};

export default MetricCard;
