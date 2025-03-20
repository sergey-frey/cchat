import { email, InferInput, minLength, object, pipe, string } from "valibot";

export const AuthFormSchema = object({
  email: pipe(string(), email("Invalid email")),
  password: pipe(
    string(),
    minLength(8, "Password must be at least 8 characters long"),
  ),
});

export type AuthFormSchemaType = InferInput<typeof AuthFormSchema>;
