import { chatsApi } from "@/shared/api/instance/instance";
import { IChatPreview, IResponse } from "@/shared/api/types";
import { queryClient } from "@/shared/query-client";
import { QueryHooks } from "@/shared/types/query";
import { useMutation, useQuery } from "@tanstack/react-query";
import { CreateChatDto } from "./dto";

export const useChatsQuery = (search?: string) => {
  return useQuery<IResponse<IChatPreview[]>>({
    queryKey: ["chats", search],
    queryFn: () => {
      return chatsApi
        .get("preview", { searchParams: { ...(search ? { search } : {}) } })
        .json();
    },
  });
};

export const useCreateChatMutation = ({ onSettled }: QueryHooks = {}) => {
  return useMutation({
    mutationFn: ({ users }: CreateChatDto) => {
      return chatsApi.post("new", { json: { users } }).json();
    },
    mutationKey: ["createChat"],
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["chats"] });
      onSettled?.();
    },
  });
};
