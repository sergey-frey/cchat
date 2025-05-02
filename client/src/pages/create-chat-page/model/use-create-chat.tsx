import { IUser, useSearchUsersQuery } from "@/entities/user";
import { useIntersection } from "@/shared/utils/use-intersection";
import { Spinner } from "@heroui/spinner";
import debounce from "debounce";
import { useCallback } from "react";
import { CREATE_CHAT_SEARCH_LIMIT } from "../constants";
import { PLACEHOLDER_USERS } from "../constants/placeholder";

type UseCreateChatOptions = {
  search: string;
  chatMembers: IUser[];
};

export const useCreateChat = ({
  search,
  chatMembers,
}: UseCreateChatOptions) => {
  const {
    data: users,
    isPending: isSearchPending,
    isRefetching: isSearchRefetching,
    isPlaceholderData: isSearchPlaceholderData,
    refetch: refetchUsers,
    fetchNextPage: fetchNextUsersPage,
    hasNextPage: hasNextUsersPage,
    isFetchingNextPage: isFetchingNextUsersPage,
    error: fetchUsersError,
  } = useSearchUsersQuery(
    {
      username: search,
      limit: CREATE_CHAT_SEARCH_LIMIT,
    },
    {
      placeholderData: PLACEHOLDER_USERS,
    },
  );

  const paginationTriggerRef = useIntersection<HTMLDivElement>({
    onIntersect: fetchNextUsersPage,
  });

  const debouncedRefetchUsers = useCallback(
    () => debounce(refetchUsers, 500)(),
    [refetchUsers],
  );

  const isSearchInputLoading =
    isSearchPending || isSearchRefetching || isFetchingNextUsersPage;

  const searchInputEndContent = isSearchInputLoading ? (
    <Spinner size="sm" />
  ) : null;

  const isShowCreateChatButton = chatMembers.length > 0;

  const isShowPlaceholders = isSearchPlaceholderData || isSearchRefetching;

  return {
    users,
    hasNextUsersPage,
    debouncedRefetchUsers,
    paginationTriggerRef,
    searchInputEndContent,
    isShowCreateChatButton,
    fetchUsersError,
    isShowPlaceholders,
  };
};
