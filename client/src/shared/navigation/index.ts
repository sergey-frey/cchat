const withSearchParams = (navigationFn: (args: unknown) => string) => {
  return ({
    searchParams,
    ...args
  }: { searchParams?: Record<string, string> } & Parameters<
    typeof navigationFn
  >[0]) => {
    if (!searchParams) {
      return navigationFn(args);
    }

    return `${navigationFn(args)}?${new URLSearchParams(searchParams).toString()}`;
  };
};

export const NAVIGATION = {
  auth: withSearchParams(() => `~/auth`),
  profile: "~/app/profile",
  editProfile: "~/app/profile/edit",
  chats: (id?: string) => (id ? `~/app/chats/${id}` : "~/app/chats"),
  createChat: "~/app/chats/create",
  channels: () => "~/app/channels",
  settings: () => "~/app/settings",
};
