export const scrollPageToBottom = (buffer: number) => {
  if (
    window.scrollY + window.innerHeight >=
    document.body.scrollHeight - buffer
  ) {
    setTimeout(() => {
      window.scrollTo(0, document.body.scrollHeight);
    }, 200);
  }
};
