import { authApi } from "@/shared/api/instance/instance";
import { IResponse, IUser, IUserServerData } from "@/shared/api/types";
import { KyInstance, KyResponse } from "ky";
import { LoginDto, RegisterDto } from "../types/dto";

class AuthService {
  private _instance: KyInstance;

  constructor(instance: KyInstance) {
    this._instance = instance;
  }

  private async _userResponseTransform(
    response: KyResponse<IResponse<IUserServerData>>,
  ): Promise<IResponse<IUser>> {
    const parsedResponse = await response.json();

    return {
      status: parsedResponse.status,
      data: {
        id: parsedResponse.data.ID,
        username: parsedResponse.data.Username,
        email: parsedResponse.data.Email,
      },
    };
  }

  async register({ email, password }: RegisterDto) {
    const response = await this._instance.post<IResponse<IUserServerData>>(
      "register",
      {
        json: { email, password },
      },
    );

    if (response.ok) {
      return this._userResponseTransform(response);
    }

    throw new Error(response.statusText);
  }

  async login({ email, password }: LoginDto) {
    const response = await this._instance.post<IResponse<IUserServerData>>(
      "login",
      {
        json: { email, password },
      },
    );

    if (response.ok) {
      return this._userResponseTransform(response);
    }

    throw new Error(response.statusText);
  }

  async checkSession() {
    const response =
      await this._instance.post<IResponse<IUserServerData>>("session");

    if (response.ok) {
      return this._userResponseTransform(response);
    }

    throw new Error(response.statusText);
  }

  async logout() {
    const response = await this._instance.post<IResponse<undefined>>("logout");

    if (response.ok) {
      return response.json();
    }

    throw new Error(response.statusText);
  }
}

export const authService = new AuthService(authApi);
