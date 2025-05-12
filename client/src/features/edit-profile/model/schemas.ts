import {
  customAsync,
  email,
  minLength,
  pipe,
  pipeAsync,
  string,
} from "valibot";
import { checkUniqueUsername } from "./validate-username";

export const usernameSchema = pipeAsync(
  string(),
  minLength(3, "Username must be at least 3 characters long"),
  customAsync(async (input) => {
    const username = input as string;
    return await checkUniqueUsername(username);
  }, "Username is already taken"),
);

export const emailSchema = pipe(string(), email("Invalid email"));

export const nameSchema = pipe(
  string(),
  minLength(4, "Name must be at latest 4 characters long"),
);

export type EditProfileFormSchema = {
  username: string;
  email: string;
  name: string;
};
