import { userApi } from "@/shared/api/instance/instance";
import { queryClient } from "@/shared/query-client";
import { useMutation, useQuery, UseQueryOptions } from "@tanstack/react-query";
import { UpdateUserDto } from "../types/dto";
import { IUserProfileResponse } from "../types/responses";
import { userService } from "./user-service";

export const useProfileQuery = (
  options: Omit<
    UseQueryOptions<IUserProfileResponse["data"]>,
    "queryFn" | "queryKey"
  > = {},
) => {
  return useQuery<IUserProfileResponse["data"]>({
    queryKey: ["profile"],
    queryFn: async () => {
      const res = await userService.getProfile();
      return res.data;
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

export const useUpdateProfileQuery = () => {
  return useMutation<IUserProfileResponse["data"], Error, UpdateUserDto>({
    mutationFn: async (updateData) => {
      const res = await userService.updateProfile(updateData);
      return res.data;
    },
    onSuccess: (data) => {
      queryClient.setQueryData(["profile"], data);
    },
  });
};
