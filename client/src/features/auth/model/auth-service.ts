import { authApi } from "@/shared/api/instance/instance";
import { IUser } from "@/shared/api/types";
import { KyInstance } from "ky";
import { LoginDto, RegisterDto } from "../types/dto";

class AuthService {
  private _instance: KyInstance;

  constructor(instance: KyInstance) {
    this._instance = instance;
  }

  async register({ email, password }: RegisterDto) {
    const response = await this._instance.post<IUser>("register", {
      json: { email, password },
    });

    if (response.ok) {
      return response.json();
    }

    throw new Error(response.statusText);
  }

  async login({ email, password }: LoginDto) {
    const response = await this._instance.post<IUser>("login", {
      json: { email, password },
    });

    if (response.ok) {
      return response.json();
    }

    throw new Error(response.statusText);
  }

  async checkSession() {
    const response = await this._instance.post<IUser>("session");

    if (response.ok) {
      return response.json();
    }

    throw new Error(response.statusText);
  }

  async logout() {
    const response = await this._instance.get<IUser>("logout");

    if (response.ok) {
      return response.json();
    }

    throw new Error(response.statusText);
  }
}

export const authService = new AuthService(authApi);
