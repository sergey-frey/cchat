import { HttpResponse, http } from "msw";

let isAuthorized = false;

const loginResolver = () => {
  isAuthorized = true;

  return HttpResponse.json({
    id: "1",
    username: "@admin",
    email: "admin@test.com",
  });
};

const sessionResolver = () => {
  if (!isAuthorized) {
    return HttpResponse.json(null, {
      status: 401,
    });
  }

  return HttpResponse.json({
    id: "1",
    username: "@admin",
    email: "admin@test.com",
  });
};

const logoutResolver = () => {
  isAuthorized = false;

  return HttpResponse.json({ message: "Logout successfully" });
};

const loginHandler = http.post(
  "http://localhost:3000/auth/login",
  loginResolver,
);

const sessionHandler = http.post(
  "http://localhost:3000/auth/session",
  sessionResolver,
);

const logoutHandler = http.get(
  "http://localhost:3000/auth/logout",
  logoutResolver,
);

export const handlers = [loginHandler, sessionHandler, logoutHandler];
