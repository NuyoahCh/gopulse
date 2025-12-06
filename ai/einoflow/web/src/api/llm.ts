import client from './client';

export interface ChatMessage {
  role: 'user' | 'assistant' | 'system';
  content: string;
}

export interface ChatRequest {
  provider: string;
  model: string;
  messages: ChatMessage[];
  stream?: boolean;
}

export interface ChatResponse {
  content: string;
  model: string;
  usage?: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
}

export interface Model {
  id: string;
  provider: string;
  name: string;
}

export interface ModelsResponse {
  models: Model[];
}

// 普通对话
export const chat = async (request: ChatRequest): Promise<ChatResponse> => {
  const response = await client.post<ChatResponse>('/v1/llm/chat', request);
  return response.data;
};

// 流式对话
export const chatStream = async (
  request: ChatRequest,
  onMessage: (content: string) => void,
  onComplete: () => void,
  onError: (error: Error) => void
): Promise<void> => {
  try {
    const response = await fetch('/api/v1/llm/chat/stream', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const reader = response.body?.getReader();
    const decoder = new TextDecoder();

    if (!reader) {
      throw new Error('No response body');
    }

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      const chunk = decoder.decode(value);
      const lines = chunk.split('\n');

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = line.slice(6);
          if (data === '[DONE]') {
            onComplete();
            return;
          }
          try {
            const parsed = JSON.parse(data);
            if (parsed.content) {
              onMessage(parsed.content);
            }
          } catch (e) {
            console.warn('Failed to parse SSE data:', data);
          }
        }
      }
    }

    onComplete();
  } catch (error) {
    onError(error as Error);
  }
};

// 获取可用模型列表
export const listModels = async (): Promise<ModelsResponse> => {
  const response = await client.get<ModelsResponse>('/v1/llm/models');
  return response.data;
};
