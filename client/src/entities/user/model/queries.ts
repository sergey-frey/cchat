import { userApi } from "@/shared/api/instance/instance";
import { IUser } from "@/shared/api/types";
import { useQuery } from "@tanstack/react-query";

export const useProfileQuery = () => {
  return useQuery<Omit<IUser, "id">>({
    queryKey: ["profile"],
    queryFn: () => {
      return userApi
        .get<IUser>("profile")
        .json()
        .then((res) => ({
          username: res.username,
          email: res.email,
        }));
    },
  });
};

export const useCheckUsernameQuery = () => {
  return useQuery<{ isUnique: boolean }>({
    queryKey: [],
    queryFn: () => {
      return userApi.get("profile").json();
    },
  });
};
