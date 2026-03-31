import type {
  ApiResponse,
  ArgConfig,
  AuthStatus,
  CommonToolItemInfo,
  CreateMcpServerToolResponse,
  GetToolsInfoByUUIDResponse,
  TestMcpServerToolResponse,
  ToolProtocolInfo,
} from '../types';
import { http } from './client';

const prefix = '/v1';

function sanitizeArgs(args?: ArgConfig[]): ArgConfig[] {
  if (!args?.length) return [];
  return args.map((a) => ({
    name: a.name,
    position: a.position,
    required: a.required,
    type: a.type || 'string',
    description: a.description ?? '',
    default: a.default ?? '',
    items: a.items ?? { type: '' },
    enum: a.enum,
    explode: a.explode ?? false,
  }));
}

/** 后端 GetMcpServerTools 的 data 直接为 tools 数组 */
export function getMcpServerTools(uuid: string) {
  return http.post<ApiResponse<ToolProtocolInfo[]>>(`${prefix}/mcpServer/getMcpServerTools`, { uuid });
}

export function getMcpServerToolsByUUID(uuid: string) {
  return http.post<ApiResponse<{ tools: CommonToolItemInfo[] }>>(`${prefix}/mcpServer/getMcpServerToolsByUUID`, {
    uuid,
  });
}

/**
 * 创建工具：后端 CreateMcpServerToolRequest 的 json 标签与字段语义交叉，
 * 请求体中 isAuth 表示平台鉴权，isPlatformAuth 表示接口鉴权。
 */
export function createMcpServerTool(payload: {
  mcpServerUUID: string;
  name: string;
  description: string;
  method: string;
  path: string;
  platformAuth: AuthStatus;
  toolAuth: AuthStatus;
  args?: ArgConfig[];
}) {
  const body = {
    mcpServerUUID: payload.mcpServerUUID,
    name: payload.name,
    description: payload.description,
    method: payload.method,
    path: payload.path,
    isAuth: payload.platformAuth,
    isPlatformAuth: payload.toolAuth,
    args: sanitizeArgs(payload.args),
  };
  return http.post<ApiResponse<CreateMcpServerToolResponse>>(`${prefix}/mcpServer/createMcpServerTool`, body);
}

export function updateMcpServerTool(payload: {
  uuid: string;
  name: string;
  description?: string;
  method: string;
  path: string;
  isPlatformAuth: AuthStatus;
  isAuth: AuthStatus;
  args?: ArgConfig[];
}) {
  const body = {
    uuid: payload.uuid,
    name: payload.name,
    description: payload.description ?? '',
    method: payload.method,
    path: payload.path,
    isPlatformAuth: payload.isPlatformAuth,
    isAuth: payload.isAuth,
    args: sanitizeArgs(payload.args),
  };
  return http.post<ApiResponse<unknown>>(`${prefix}/mcpServer/updateMcpServerTool`, body);
}

export function getToolsInfoByUUID(toolUUID: string) {
  return http.post<ApiResponse<GetToolsInfoByUUIDResponse>>(`${prefix}/mcpServer/getToolsInfoByUUID`, {
    uuid: toolUUID,
  });
}

/** 后端 TestMcpServerTool 当前未返回业务数据，仅占位 */
export function testMcpServerTool(payload: { uuid: string; arguments?: Record<string, unknown> }) {
  return http.post<ApiResponse<TestMcpServerToolResponse>>(`${prefix}/mcpServer/testMcpServerTool`, {
    uuid: payload.uuid,
    arguments: payload.arguments ?? {},
  });
}
