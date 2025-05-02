import { API_ENDPOINTS } from "@/shared/api/constants/endpoints";
import { userApi } from "@/shared/api/instance/instance";
import { KyInstance } from "ky";
import { SearchUsersDto, UpdateUserDto } from "../types/dto";
import {
  ISearchUsersResponse,
  IUpdateUserResponse,
  IUserProfileResponse,
} from "../types/responses";

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

  async searchUsers(
    { username, limit, pagination }: SearchUsersDto,
    { signal }: { signal?: AbortSignal } = {},
  ) {
    const response = await this._instance.get<ISearchUsersResponse>(
      API_ENDPOINTS.USER.SEARCH,
      {
        searchParams: {
          username,
          ...(limit ? { limit: limit.toString() } : {}),
          ...(pagination ? { pagination: pagination.toString() } : {}),
        },
        signal,
      },
    );

    if (response.ok) {
      return await response.json();
    }

    throw new Error(response.statusText);
  }
}

export const userService = new UserService(userApi);
