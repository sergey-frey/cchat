export const getMessageDate = (date: Date) => {
  const d = new Date(date);
  return `${d.getHours()}:${d.getMinutes()}`;
};
