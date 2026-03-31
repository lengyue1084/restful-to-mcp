// Auth types
export type AuthStatus = 1 | 2; // 1=否, 2=是
export type AuthTypeStatus = 1 | 2 | 3 | 4; // 1=不开启, 2=平台授权, 3=service授权, 4=全部授权
export type ServerStatus = 1 | 2 | 3; // 1=未设置token, 2=已设置, 3=正常
export type Status = 1 | 2; // 1=隐藏, 2=显示
export type CommonStatus = 1 | 2; // 1=否, 2=是
export type SourceType = 1 | 2; // 1=文件, 2=表单
export type McpServerTypeStatus = 1 | 2; // 1=OpenAPI, 2=gRPC

export type AuthMode = 'apiKey' | 'http' | '';
export type AuthPosition = 'header' | 'query' | 'cookie';

export interface Security {
  securityKey: string;
  mode: AuthMode;
  name: string;
  scheme: string;
  in: AuthPosition;
  description: string;
  bearerFormat: string;
}

export interface ArgConfig {
  name: string;
  position: string; // header, query, path, body
  required: boolean;
  type: string;
  description: string;
  default: string;
  items?: ItemsConfig;
  enum?: string[];
  explode: boolean;
}

export interface ItemsConfig {
  type: string;
  enum?: string[];
  properties?: Record<string, any>;
  items?: ItemsConfig;
  required?: string[];
}

export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

// OpenAPI Upload
export interface OpenapiUploadRequest {
  uuid: string;
  name: string;
  fileContent: string;
  description: string;
  suffix: string;
  isAuth: AuthTypeStatus;
  serviceToken?: string;
  platformToken?: string;
}

export interface ToolInfo {
  ID: number;
  McpServerId: number;
  uuid: string;
  name: string;
  description: string;
  McpServerType: McpServerTypeStatus;
  method: string;
  endpoint: string;
  headers: string;
  args: string;
  requestBody: string;
  responseBody: string;
  toolSchema: any;
  annotations: string;
  security: string;
  isAuth: AuthStatus;
  authMode: string;
  isPlatformAuth: AuthStatus;
  isShow: Status;
  createdAt: string;
  updatedAt: string;
}

export interface OpenapiUploadResponse {
  id: number;
  name: string;
  uuid: string;
  description: string;
  urls: string[];
  createdAt: string;
  updatedAt: string;
  tools: ToolInfo[];
  version: string;
  allTools: string[];
  status: ServerStatus;
  headers: Record<string, string>;
}

// Update Auth
export interface OpenapiUpdateForAuthRequest {
  uuid: string;
  tools: { id: number; isAuth: AuthStatus }[];
}

// MCP Server
export interface GetMcpServerInfoByUUIDResponse {
  id: number;
  createdAt: string;
  updatedAt: string;
  status: ServerStatus;
  source: SourceType;
  uuid: string;
  name: string;
  description: string;
  urls: string[];
  version: string;
  isAuth: AuthTypeStatus;
  platformToken: string;
  serviceToken: string;
  headers: Record<string, string>;
  security: Security;
}

export interface UpdateMcpServerByUUIDRequest {
  uuid: string;
  name: string;
  description: string;
}

export interface CreateMcpServerByFormRequest {
  uuid: string;
  name: string;
  description: string;
  urls: string[];
  version: string;
  isAuth: AuthTypeStatus;
  platformToken?: string;
  serviceToken?: string;
  headers?: Record<string, string>;
  security?: Security;
}

export interface CreateMcpServerByFormResponse {
  id: number;
  uuid: string;
  createdAt: string;
}

export type UpdateMcpServerByFormRequest = CreateMcpServerByFormRequest;

export interface McpServerListItem {
  id: number;
  uuid: string;
  name: string;
  description: string;
  version: string;
  status: ServerStatus;
  source: SourceType;
  isAuth: AuthTypeStatus;
  toolCount: number;
  createdAt: string;
  updatedAt: string;
}

// MCP Tools
export interface ToolProtocolInfo {
  id: number;
  uuid: string;
  isAuth: AuthStatus;
  authMode: string;
  isPlatformAuth: AuthStatus;
  isShow: Status;
  isRepeat: CommonStatus;
  fullName: string;
  name: string;
  description: string;
  inputSchema: any;
  outputSchema: any;
  annotations: any;
}

export interface CommonToolItemInfo {
  id: number;
  uuid: string;
  createdAt: string;
  updatedAt: string;
  mcpServerId: number;
  mcpServerUUID: string;
  name: string;
  description: string;
  mcpServerType: McpServerTypeStatus;
  method: string;
  endpoint: string;
  headers: string;
  security: string;
  isAuth: AuthStatus;
  authMode: string;
  isPlatformAuth: AuthStatus;
  isShow: Status;
  serialNumber: string;
  isRepeat: CommonStatus;
}

export interface CreateMcpServerToolRequest {
  mcpServerUUID: string;
  name: string;
  description: string;
  method: string;
  path: string;
  isAuth: AuthStatus;
  isPlatformAuth: AuthStatus;
  args?: ArgConfig[];
}

export interface CreateMcpServerToolResponse {
  id: number;
  uuid: string;
  createdAt: string;
}

export interface UpdateMcpServerToolRequest {
  uuid: string;
  name: string;
  description?: string;
  method: string;
  path: string;
  isPlatformAuth: AuthStatus;
  isAuth: AuthStatus;
  args?: ArgConfig[];
}

export interface GetToolsInfoByUUIDResponse {
  id: number;
  uuid: string;
  createdAt: string;
  updatedAt: string;
  mcpServerId: number;
  mcpServerUUID: string;
  name: string;
  description: string;
  mcpServerType: McpServerTypeStatus;
  method: string;
  baseUrl: string;
  path: string;
  headers: Record<string, string>;
  args: ArgConfig[];
  security: Security;
  isAuth: AuthStatus;
  authMode: string;
  isPlatformAuth: AuthStatus;
  isShow: Status;
  serialNumber: string;
  isRepeat: CommonStatus;
}

export interface TestMcpServerToolResponse {
  name: string;
  requestUrl: string;
  requestMethod: string;
  requestBody?: unknown;
  responseText: string;
  responseJson?: Record<string, unknown>;
}

// Labels
export const AUTH_TYPE_LABELS: Record<number, string> = {
  1: '不开启', 2: '平台授权', 3: 'Service授权', 4: '全部授权',
};
export const AUTH_STATUS_LABELS: Record<number, string> = { 1: '否', 2: '是' };
export const SERVER_STATUS_LABELS: Record<number, string> = { 1: '未设置Token', 2: '已设置Token', 3: '正常工作' };
export const STATUS_LABELS: Record<number, string> = { 1: '隐藏', 2: '显示' };
export const COMMON_STATUS_LABELS: Record<number, string> = { 1: '否', 2: '是' };
export const SOURCE_TYPE_LABELS: Record<number, string> = { 1: 'OpenAPI文档', 2: '表单创建' };
export const SERVER_STATUS_COLORS: Record<number, string> = { 1: 'warning', 2: 'processing', 3: 'success' };

export const FIELD_TYPES = ['String', 'Number', 'Boolean', 'Object', 'Array<String>', 'Array<Number>', 'Array<Boolean>', 'Array<Object>'];
export const ARG_POSITIONS = ['header', 'query', 'path', 'body'];
export const HTTP_METHODS = ['GET', 'POST', 'PUT', 'DELETE'];
