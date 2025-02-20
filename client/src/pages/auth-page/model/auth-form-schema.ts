import { email, InferInput, object, pipe, string } from "valibot";

export const AuthFormSchema = object({
  email: pipe(string(), email()),
  password: string(),
});

export type AuthFormSchemaType = InferInput<typeof AuthFormSchema>;
