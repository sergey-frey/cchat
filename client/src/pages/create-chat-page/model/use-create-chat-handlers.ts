import { useCreateChatMutation } from "@/entities/chats";
import { IUser } from "@/entities/user";
import {
  containerRefSelector,
  useAppContainer,
} from "@/shared/utils/app-container";
import { UseArrayStateMethods } from "@/shared/utils/hooks";
import { Dispatch, SetStateAction } from "react";

type UseCreateChatHandlersOptions = {
  setSearch: Dispatch<SetStateAction<string>>;
  debouncedRefetchUsers: () => void;
  chatMembers: IUser[];
  chatMembersMethods: UseArrayStateMethods<IUser>;
};

export const useCreateChatHandlers = ({
  setSearch,
  debouncedRefetchUsers,
  chatMembers,
  chatMembersMethods,
}: UseCreateChatHandlersOptions) => {
  const containerRef = useAppContainer(containerRefSelector);

  const createChatMutation = useCreateChatMutation();

  const handleSearchChange = (value: string) => {
    setSearch(value);
    debouncedRefetchUsers();
  };

  const handleScrollTopClick = () => {
    if (!containerRef?.current) return;
    containerRef.current.scrollTo({ top: 0, behavior: "smooth" });
  };

  const handleUserClick = (user: IUser) => {
    if (chatMembers.indexOf(user) !== -1) {
      chatMembersMethods.remove(user);
    } else {
      chatMembersMethods.pushUnique(user);
    }
  };

  const handleClickCreateChatButton = async () => {
    // TODO: navigate to chat page
    await createChatMutation.mutateAsync({
      users: chatMembers.map((user) => user.id),
    });
  };

  return {
    handleSearchChange,
    handleScrollTopClick,
    handleUserClick,
    handleClickCreateChatButton,
  };
};
