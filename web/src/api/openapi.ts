import type { ApiResponse, OpenapiUpdateForAuthRequest, OpenapiUploadRequest, OpenapiUploadResponse } from '../types';
import { http } from './client';

const prefix = '/v1';

export function uploadOpenapi(payload: OpenapiUploadRequest) {
  return http.post<ApiResponse<OpenapiUploadResponse>>(`${prefix}/openapi/upload`, payload);
}

export function updateForAuth(payload: OpenapiUpdateForAuthRequest) {
  return http.post<ApiResponse<unknown>>(`${prefix}/openapi/updateForAuth`, payload);
}
