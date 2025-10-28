import { faker } from "@faker-js/faker";
import { HttpResponse } from "msw";
import { generateRandomUser } from "./faker";
import { MockAPI } from "./mock-api";

export const mockApi = new MockAPI("http://localhost:8040/cchat");

const delay = (ms: number) => new Promise((res) => setTimeout(res, ms));

const users = faker.helpers.multiple(generateRandomUser, { count: 10000 });

mockApi.get("user/search", async ({ request }) => {
  await delay(1000);

  const searchParams = new URL(request.url).searchParams;

  const username = searchParams.get("username");
  const limit = searchParams.get("limit");
  const pagination = searchParams.get("pagination");

  let filteredUsers = users.filter((user) => {
    return user.username.startsWith(username as string);
  });

  if (!username) {
    filteredUsers = users;
  }

  if (limit && pagination) {
    filteredUsers = filteredUsers.slice(
      (Number(pagination) - 1) * Number(limit),
      Number(pagination) * Number(limit),
    );
  }

  return HttpResponse.json({
    status: 200,
    data: filteredUsers,
  });
});
