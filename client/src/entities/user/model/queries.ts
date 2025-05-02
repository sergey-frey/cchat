import { userApi } from "@/shared/api/instance/instance";
import { queryClient } from "@/shared/query-client";
import {
  InfiniteData,
  useInfiniteQuery,
  useMutation,
  useQuery,
  UseQueryOptions,
} from "@tanstack/react-query";
import { SearchUsersDto, UpdateUserDto } from "../types/dto";
import { ISearchUsersResponse, IUserProfileResponse } from "../types/responses";
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

export const useSearchUsersQuery = (
  { username, limit }: SearchUsersDto,
  {
    placeholderData,
  }: {
    placeholderData?: InfiniteData<ISearchUsersResponse["data"], number>;
  } = {},
) => {
  return useInfiniteQuery({
    queryKey: ["users-search"],
    queryFn: async ({ signal, pageParam }) => {
      const res = await userService.searchUsers(
        {
          username,
          limit,
          pagination: pageParam as number,
        },
        { signal },
      );

      return res.data;
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage, _, lastPageParam) => {
      if (lastPage.length === 0) {
        return undefined;
      }

      return (lastPageParam as number) + 1;
    },
    select: (data) => {console.log(data);return data.pages.flat()},
    placeholderData,
  });
};
