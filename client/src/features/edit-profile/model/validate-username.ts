import { userService } from "@/entities/user";
import { HTTPError } from "ky";

export const checkUniqueUsername = async (
  username: string,
): Promise<boolean> => {
  return userService
    .getByUsername(username)
    .then((response) => {
      return response.status === 404;
    })
    .catch((response: HTTPError) => {
      return response.response.status === 404;
    });
};
