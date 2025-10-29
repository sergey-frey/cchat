export const getChatSocket = (chatId: string) => {
  return new WebSocket(`ws://localhost:3000/chats/ws/${chatId}`);
};
