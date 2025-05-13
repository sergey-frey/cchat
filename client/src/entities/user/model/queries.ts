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
    initialPageParam: 0,
    getNextPageParam: (lastPage) => {
      if (!lastPage?.cursors.next_cursor) {
        return undefined;
      }

      return lastPage?.cursors.next_cursor;
    },
    select: (data) => {
      return data.pages.map((page) => page.profiles).flat();
    },
    placeholderData,
  });
};
