import client from './client';
import type { Workflow } from '../types/workflow';

export const fetchWorkflows = async (): Promise<Workflow[]> => {
  try {
    const response = await client.get<Workflow[]>('/workflows');
    return response.data;
  } catch (error) {
    // graceful fallback for demo mode
    console.warn('Falling back to mock workflows', error);
    return mockWorkflows;
  }
};

const now = new Date();
const format = (date: Date) => date.toLocaleString('en-US', { hour12: false });

const mockWorkflows: Workflow[] = [
  {
    id: 'wf-analytics',
    name: 'Daily Product Analytics',
    owner: 'Insights Team',
    lastRun: format(new Date(now.getTime() - 1000 * 60 * 5)),
    status: 'healthy',
    throughput: 560,
    description: 'Aggregates metrics from the data warehouse and syncs dashboards.',
    tags: ['analytics', 'batch']
  },
  {
    id: 'wf-alerting',
    name: 'Realtime Alert Orchestrator',
    owner: 'SRE',
    lastRun: format(new Date(now.getTime() - 1000 * 60 * 42)),
    status: 'degraded',
    throughput: 220,
    description: 'Listens to platform events and triggers incident communication workflows.',
    tags: ['streaming', 'alerts']
  },
  {
    id: 'wf-growth',
    name: 'Lifecycle Journey Sync',
    owner: 'Growth Ops',
    lastRun: format(new Date(now.getTime() - 1000 * 60 * 120)),
    status: 'failed',
    throughput: 0,
    description: 'Pushes curated segments to marketing destinations every 2 hours.',
    tags: ['crm', 'sync']
  }
];
