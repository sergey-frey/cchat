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
import { USER_QUERY_KEYS } from "./query-keys";
import { userService } from "./user-service";

export const useProfileQuery = (
  options: Omit<
    UseQueryOptions<IUserProfileResponse["data"]>,
    "queryFn" | "queryKey"
  > = {},
) => {
  return useQuery<IUserProfileResponse["data"]>({
    queryKey: [USER_QUERY_KEYS.PROFILE],
    queryFn: async () => {
      const res = await userService.getProfile();
      return res.data;
    },
    ...options,
  });
};

export const useUpdateProfileQuery = () => {
  return useMutation<IUserProfileResponse["data"], Error, UpdateUserDto>({
    mutationFn: async (updateData) => {
      const res = await userService.updateProfile(updateData);
      return res.data;
    },
    onSuccess: (data) => {
      queryClient.setQueryData([USER_QUERY_KEYS.PROFILE], data);
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
    queryKey: [USER_QUERY_KEYS.USERS_SEARCH],
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
    select: (data) => {
      console.log(data);
      return data.pages.flat();
    },
    placeholderData,
  });
};
