import { userApi } from "@/shared/api/instance/instance";

export const checkUniqueUsername = async (
  username: string,
): Promise<boolean> => {
  return userApi
    .get<{ isUnique: boolean }>("check-username", {
      searchParams: { username },
    })
    .json()
    .then((res) => res.isUnique)
    .catch(() => {
      return false;
    });
};
