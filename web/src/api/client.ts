import axios from 'axios';
import type { AxiosError, AxiosResponse } from 'axios';
import type { ApiResponse } from '../types';

const baseURL = import.meta.env.VITE_API_BASE_URL ?? '';

export const http = axios.create({
  baseURL,
  timeout: 120000,
  headers: { 'Content-Type': 'application/json' },
});

http.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const body = response.data;
    if (body && typeof body.code === 'number' && body.code !== 0) {
      return Promise.reject(new Error(body.message || '请求失败'));
    }
    return response;
  },
  (error: AxiosError) => {
    const message =
      typeof error.response?.data === 'object' &&
      error.response?.data &&
      'message' in error.response.data
        ? String(error.response.data.message)
        : error.message || '网络请求失败';
    return Promise.reject(new Error(message));
  },
);
