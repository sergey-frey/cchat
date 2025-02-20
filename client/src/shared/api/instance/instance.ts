import ky from "ky";

const baseUrl = "http://localhost:3000";

const api = ky.create({
  prefixUrl: baseUrl,
  headers: {
    "Content-Type": "application/json",
  },
  credentials: "include",
});

export const authApi = api.extend({
  prefixUrl: new URL("/auth", baseUrl).href,
});

export const chatsApi = api.extend({
  prefixUrl: new URL("/chats", baseUrl).href,
});

export const userApi = api.extend({
  prefixUrl: new URL("/user", baseUrl).href,
});
