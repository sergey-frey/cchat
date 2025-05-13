import { IResponse, IResponseWithCursor } from "@/shared/api/types";

export type IUserProfileResponse = IResponse<{
  id: string;
  username: string;
  email: string;
  name: string;
}>;

export type IUpdateUserResponse = IResponse<{
  id: string;
  username: string;
  email: string;
  name: string;
}>;

export type IGetUserByUsernameResponse = IResponse<{
  id: string;
  username: string;
  email: string;
  name: string;
}>;

export type ISearchUsersResponse = IResponseWithCursor<{
  profiles: Array<{
    id: string;
    username: string;
    email: string;
    name: string;
  }>;
}>;
