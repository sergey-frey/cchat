import { IUser, IChatPreview } from "@/shared/api/types";

const getChatTitle = (
  currentUser: IUser,
  chatPreview: IChatPreview,
  isPersonalChat: boolean,
) => {
  if (!isPersonalChat) {
    return chatPreview.id;
  }

  const companion = chatPreview.users.find(
    (user) => user.id !== currentUser.id,
  );

  if (!companion) {
    return chatPreview.id;
  }

  return companion.username;
};

type UseChatInfoOptions = {
  currentUser: IUser;
  chatPreview: IChatPreview;
};

export const useChatInfo = ({
  chatPreview,
  currentUser,
}: UseChatInfoOptions) => {
  const isPersonalChat = chatPreview.users.length === 2;
  const chatTitle = getChatTitle(currentUser, chatPreview, isPersonalChat);
  const lastMessage = chatPreview.messages[chatPreview.messages.length - 1];

  return { isPersonalChat, chatTitle, lastMessage };
};
