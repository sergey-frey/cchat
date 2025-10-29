import ky from "ky";
import urlJoin from "url-join";

const baseUrl = import.meta.env.VITE_BASE_API_URL;

if (!baseUrl) {
  throw new Error("VITE_BASE_API_URL is not provide");
}

console.log(baseUrl);

const api = ky.create({
  prefixUrl: baseUrl,
  headers: {
    "Content-Type": "application/json",
  },
  credentials: "include",
});

export const authApi = api.extend({
  prefixUrl: urlJoin(baseUrl, "auth"),
});

export const chatsApi = api.extend({
  prefixUrl: urlJoin(baseUrl, "chats"),
});

export const userApi = api.extend({
  prefixUrl: urlJoin(baseUrl, "user"),
});
