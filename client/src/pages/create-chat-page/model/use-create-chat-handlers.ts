import { IUser } from "@/entities/user";
import {
  containerRefSelector,
  useAppContainer,
} from "@/shared/utils/app-container";
import { UseArrayStateMethods } from "@/shared/utils/use-array-state";
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

  const handleSearchChange = (value: string) => {
    setSearch(value);
    debouncedRefetchUsers();
  };

  const handleScrollTopClick = () => {
    if (!containerRef?.current) return;
    containerRef.current.scrollTo({ top: 0, behavior: "smooth" });
  };

  const handleUserClick = (user: IUser) => {
    if (chatMembers.includes(user)) {
      chatMembersMethods.remove(user);
    } else {
      chatMembersMethods.pushUnique(user);
    }
  };

  return { handleSearchChange, handleScrollTopClick, handleUserClick };
};
