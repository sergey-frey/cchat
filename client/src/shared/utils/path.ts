export const absolutePath = (path: string) => {
  if (path.startsWith("/")) {
    return path;
  }

  if (path.startsWith("~")) {
    return path.slice(1);
  }

  return `/${path}`;
};
