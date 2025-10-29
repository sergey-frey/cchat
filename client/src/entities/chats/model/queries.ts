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
      //   return {
      //     data: [
      //       {
      //         id: "1",
      //         users: [
      //           {
      //             username: "username_1",
      //             name: "User1",
      //           },
      //           {
      //             username: "username_2",
      //             name: "User2",
      //           },
      //         ].map((user, i) => ({ id: i.toString(), ...user })),
      //         messages: [],
      //       },
      //     ],
      //   };
      return chatsApi
        .get("preview", { searchParams: { ...(search ? { search } : {}) } })
        .json();
    },
  });
};

export const useCreateChatMutation = ({ onSettled }: QueryHooks) => {
  return useMutation({
    mutationFn: ({ users }: CreateChatDto) => {
      return chatsApi.post("/", { json: { users } }).json();
    },
    mutationKey: ["createChat"],
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["chats"] });
      onSettled?.();
    },
  });
};
