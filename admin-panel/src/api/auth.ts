import request from '../utils/request';
import type { LoginReq, LoginResp, User } from '../types';

export const login = (data: LoginReq) =>
  request.post<any, { code: number; message: string; data: LoginResp }>('/admin/login', data);

export const getProfile = () =>
  request.get<any, { code: number; message: string; data: User }>('/admin/profile');
