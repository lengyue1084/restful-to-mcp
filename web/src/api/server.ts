import type {
  ApiResponse,
  CreateMcpServerByFormRequest,
  CreateMcpServerByFormResponse,
  GetMcpServerInfoByUUIDResponse,
  McpServerListItem,
  UpdateMcpServerByFormRequest,
  UpdateMcpServerByUUIDRequest,
} from '../types';
import { http } from './client';

const prefix = '/v1';

export function listMcpServers() {
  return http.post<ApiResponse<{ items: McpServerListItem[] }>>(`${prefix}/mcpServer/list`, {});
}

export function getMcpServerInfoByUUID(uuid: string) {
  return http.post<ApiResponse<GetMcpServerInfoByUUIDResponse>>(`${prefix}/mcpServer/getMcpServerInfoByUUID`, {
    uuid,
  });
}

export function updateMcpServerByUUID(payload: UpdateMcpServerByUUIDRequest) {
  return http.post<ApiResponse<unknown>>(`${prefix}/mcpServer/updateByUUID`, payload);
}

export function getMcpConnectTokenByUUID(uuid: string) {
  return http.post<ApiResponse<{ connectToken: string }>>(`${prefix}/mcpServer/getMcpConnectTokenByUUID`, {
    uuid,
  });
}

export function deleteMcpServerByUUID(uuid: string) {
  return http.post<ApiResponse<unknown>>(`${prefix}/mcpServer/deleteMcpServerByUUID`, { uuid });
}

export function createMcpServerByForm(payload: CreateMcpServerByFormRequest) {
  return http.post<ApiResponse<CreateMcpServerByFormResponse>>(`${prefix}/mcpServer/createByForm`, payload);
}

export function updateMcpServerByForm(payload: UpdateMcpServerByFormRequest) {
  return http.post<ApiResponse<CreateMcpServerByFormResponse>>(`${prefix}/mcpServer/updateMcpServerByForm`, payload);
}
