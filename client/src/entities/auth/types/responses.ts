import { IResponse } from "@/shared/api/types";

export type IRegisterResponse = IResponse<{
  id: string;
  username: string;
  email: string;
}>;

export type ILoginResponse = IResponse<{
  id: string;
  username: string;
  email: string;
}>;

export type ISessionResponse = IResponse<{
  id: string;
  username: string;
  email: string;
}>;
