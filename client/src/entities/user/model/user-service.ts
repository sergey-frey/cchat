import { API_ENDPOINTS } from "@/shared/api/constants/endpoints";
import { userApi } from "@/shared/api/instance/instance";
import { KyInstance } from "ky";
import { UpdateUserDto } from "../types/dto";
import { IUpdateUserResponse, IUserProfileResponse } from "../types/responses";

class UserService {
  private _instance: KyInstance;

  constructor(instance: KyInstance) {
    this._instance = instance;
  }

  async getProfile(): Promise<IUserProfileResponse> {
    const response = await this._instance.get<IUserProfileResponse>(
      API_ENDPOINTS.USER.PROFILE,
    );

    if (response.ok) {
      return await response.json();
    }

    throw new Error(response.statusText);
  }

  async updateProfile(updateData: UpdateUserDto) {
    const response = await this._instance.patch<IUpdateUserResponse>(
      API_ENDPOINTS.USER.UPDATE,
      {
        json: updateData,
      },
    );

    if (response.ok) {
      return await response.json();
    }

    throw new Error(response.statusText);
  }
}

export const userService = new UserService(userApi);
