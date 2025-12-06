export type WorkflowStatus = 'healthy' | 'degraded' | 'failed';

export interface Workflow {
  id: string;
  name: string;
  owner: string;
  lastRun: string;
  status: WorkflowStatus;
  throughput: number;
  description: string;
  tags: string[];
}
