import { HttpResponse, http } from "msw";

const delay = (ms: number) => new Promise((res) => setTimeout(res, ms));

const profileResolver = async () => {
  await delay(3000);

  return HttpResponse.json({
    id: "123",
    username: "Parker2903",
    email: "admin@gmail.com",
  });
};

const profileHandler = http.get(
  "http://localhost:8040/cchat/user/profile",
  profileResolver,
);

const checkUsernameHandler = http.get(
  "http://localhost:8040/cchat/user/check-username",
  () => {
    return HttpResponse.json({
      isUnique: Math.random() > 0.5,
    });
  },
);

export const handlers = [profileHandler, checkUsernameHandler];
