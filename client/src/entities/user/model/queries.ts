import { userApi } from "@/shared/api/instance/instance";
import { IUser } from "@/shared/api/types";
import { useQuery, UseQueryOptions } from "@tanstack/react-query";

export const useProfileQuery = (
  options: Omit<
    UseQueryOptions<Omit<IUser, "id">>,
    "queryFn" | "queryKey"
  > = {},
) => {
  return useQuery<Omit<IUser, "id">>({
    queryKey: ["profile"],
    queryFn: async () => {
      const res = await userApi.get<IUser>("profile").json();

      return {
        username: res.username,
        email: res.email,
      };
    },
    ...options,
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
