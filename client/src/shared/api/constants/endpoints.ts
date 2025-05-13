import urlJoin from "url-join";

export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: "login",
    REGISTER: "register",
    LOGOUT: "logout",
    SESSION: "session",
  },
  USER: {
    MY_PROFILE: "myprofile",
    profile: (username: string) => urlJoin("profile", username),
    UPDATE: "update",
    SEARCH: "list-profiles",
  },
};
