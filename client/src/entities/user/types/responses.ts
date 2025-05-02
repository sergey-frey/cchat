import { IResponse } from "@/shared/api/types";

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

export type ISearchUsersResponse = IResponse<
  Array<{
    id: string;
    username: string;
    email: string;
    name: string;
  }>
>;
