import client from './client';

export interface AgentRequest {
  task: string;
  provider?: string;
  model?: string;
}

export interface AgentResponse {
  answer: string;
  steps?: string[];
  execution_time?: number;
}

// 运行 Agent
export const runAgent = async (request: AgentRequest): Promise<AgentResponse> => {
  const response = await client.post<AgentResponse>('/v1/agent/run', request);
  return response.data;
};
