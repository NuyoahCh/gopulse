import client from './client';

export interface GraphRequest {
  query: string;  // 后端期望 query 字段
  type?: string;
}

export interface GraphResponse {
  result: string;
  steps?: Array<{
    node: string;
    output: string;
  }>;
  execution_time?: number;
}

// 运行 Graph
export const runGraph = async (request: GraphRequest): Promise<GraphResponse> => {
  const response = await client.post<GraphResponse>('/v1/graph/run', request);
  return response.data;
};
