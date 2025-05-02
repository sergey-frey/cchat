import { IGetUserByUsernameResponse, userService } from "@/entities/user";
import { KyResponse } from "ky";

export const checkUniqueUsername = async (
  username: string,
): Promise<boolean> => {
  return userService
    .getByUsername(username)
    .then((response) => {
      return response.status === 404;
    })
    .catch((response: KyResponse<IGetUserByUsernameResponse>) => {
      return response.status === 404;
    });
};
