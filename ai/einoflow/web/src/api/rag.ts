import client from './client';

export interface RAGIndexRequest {
  documents: string[];
}

export interface RAGIndexResponse {
  count: number;
  message: string;
  total?: number;
}

export interface RAGQueryRequest {
  query: string;
  top_k?: number;
}

export interface RAGQueryResponse {
  answer: string;
  documents: string[];  // 后端返回 documents 字段
  relevance_scores?: number[];
}

export interface RAGStatsResponse {
  total_documents: number;
  total_chunks: number;
  vector_dimension: number;
}

// 索引文档
export const indexDocuments = async (request: RAGIndexRequest): Promise<RAGIndexResponse> => {
  const response = await client.post<RAGIndexResponse>('/v1/rag/index', request);
  return response.data;
};

// 查询 RAG
export const queryRAG = async (request: RAGQueryRequest): Promise<RAGQueryResponse> => {
  const response = await client.post<RAGQueryResponse>('/v1/rag/query', request);
  return response.data;
};

// 获取统计信息
export const getRAGStats = async (): Promise<RAGStatsResponse> => {
  const response = await client.get<RAGStatsResponse>('/v1/rag/stats');
  return response.data;
};

// 清空文档
export const clearRAG = async (): Promise<{ message: string }> => {
  const response = await client.delete<{ message: string }>('/v1/rag/clear');
  return response.data;
};

// 上传文件
export interface RAGUploadResponse {
  message: string;
  filename: string;
  document_count: number;
  total_count: number;
}

export const uploadFile = async (file: File): Promise<RAGUploadResponse> => {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await client.post<RAGUploadResponse>('/v1/rag/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
};
