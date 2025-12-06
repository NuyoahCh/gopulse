import { SlidersHorizontal } from 'lucide-react';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';

interface WorkflowFilterProps {
  activeStatus: string;
  onStatusChange: (status: string) => void;
}

const filters = ['all', 'healthy', 'degraded', 'failed'];

const WorkflowFilter = ({ activeStatus, onStatusChange }: WorkflowFilterProps) => {
  return (
    <div className="flex flex-wrap items-center gap-3">
      <Button variant="outline" className="gap-2 text-slate-600">
        <SlidersHorizontal size={16} />
        Filters
      </Button>
      <div className="flex flex-wrap gap-2">
        {filters.map((filter) => (
          <Badge
            key={filter}
            className={`cursor-pointer border-transparent px-4 py-2 text-sm capitalize ${
              activeStatus === filter
                ? 'bg-slate-900 text-white shadow'
                : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
            }`}
            onClick={() => onStatusChange(filter)}
          >
            {filter}
          </Badge>
        ))}
      </div>
    </div>
  );
};

export default WorkflowFilter;
