export interface ApiEnvelope<T> {
  code: number;
  message: string;
  data: T;
}

// ApiError 保留 HTTP 状态码与后端 message，便于调用方区分 404 等场景。
export class ApiError extends Error {
  status: number;
  code?: number;

  constructor(status: number, message: string, code?: number) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
  }
}

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) },
    ...init,
  });

  if (!response.ok) {
    let message = `Request failed: ${response.status}`;
    let bizCode: number | undefined;
    try {
      const payload = (await response.json()) as ApiEnvelope<unknown>;
      if (payload && typeof payload === 'object' && 'message' in payload && typeof payload.message === 'string') {
        message = payload.message;
        bizCode = payload.code;
      }
    } catch {
      // 响应体不是 JSON，保留默认 message
    }
    throw new ApiError(response.status, message, bizCode);
  }

  const payload = (await response.json()) as ApiEnvelope<T> | T;
  if (payload && typeof payload === 'object' && 'data' in payload) {
    return (payload as ApiEnvelope<T>).data;
  }
  return payload as T;
}
