import { authApi } from "@/shared/api/instance/instance";
import { IResponse } from "@/shared/api/types";
import { KyInstance } from "ky";
import { LoginDto, RegisterDto } from "../types/dto";
import {
  ILoginResponse,
  IRegisterResponse,
  ISessionResponse,
} from "../types/responses";
import { API_ENDPOINTS } from "@/shared/api/constants/endpoints";

class AuthService {
  private _instance: KyInstance;

  constructor(instance: KyInstance) {
    this._instance = instance;
  }

  async register({ email, password }: RegisterDto) {
    const response = await this._instance.post<IRegisterResponse>(
      API_ENDPOINTS.AUTH.REGISTER,
      {
        json: { email, password },
      },
    );

    if (response.ok) {
      return response.json();
    }

    throw new Error(response.statusText);
  }

  async login({ email, password }: LoginDto) {
    const response = await this._instance.post<ILoginResponse>(
      API_ENDPOINTS.AUTH.LOGIN,
      {
        json: { email, password },
      },
    );

    if (response.ok) {
      return response.json();
    }

    throw new Error(response.statusText);
  }

  async checkSession() {
    const response = await this._instance.post<ISessionResponse>(
      API_ENDPOINTS.AUTH.SESSION,
    );

    if (response.ok) {
      return response.json();
    }

    throw new Error(response.statusText);
  }

  async logout() {
    const response = await this._instance.post<IResponse<undefined>>(
      API_ENDPOINTS.AUTH.LOGOUT,
    );

    if (response.ok) {
      return response.json();
    }

    throw new Error(response.statusText);
  }
}

export const authService = new AuthService(authApi);
