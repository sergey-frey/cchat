import { useProfileQuery } from "@/entities/user";
import debounce from "debounce";
import { FormEvent, useCallback, useState } from "react";
import { safeParseAsync } from "valibot";
import { EditProfileFormSchema, emailSchema, usernameSchema } from "./schemas";

const initialFormState: EditProfileFormSchema = {
  username: "",
  email: "",
};

export const useEditProfileForm = () => {
  const [userFormState, setUserFormState] = useState<
    Partial<EditProfileFormSchema>
  >({});
  const [errors, setErrors] = useState<
    Record<keyof EditProfileFormSchema, string[]>
  >({
    username: [],
    email: [],
  });

  const profileQuery = useProfileQuery({
    staleTime: 0,
  });

  const formData = {
    ...initialFormState,
    ...profileQuery.data,
    ...userFormState,
  };

  const checkUsername = useCallback(
    debounce(async (username: string, onError?: (errors: string[]) => void) => {
      const parsedUsername = await safeParseAsync(usernameSchema, username);
      const usernameIssues =
        parsedUsername.issues?.map((issue) => issue.message) ?? [];

      onError?.(usernameIssues);
    }, 300),
    [],
  );

  const handleUsernameChange = (username: string) => {
    setUserFormState((prev) => ({ ...prev, username }));
    checkUsername(username, (errors) => {
      setErrors((prev) => ({
        ...prev,
        username: errors,
      }));
    });
  };

  const handleEmailChange = async (email: string) => {
    setUserFormState((prev) => ({ ...prev, email }));

    const parsedEmail = await safeParseAsync(emailSchema, email);
    const emailIssues = parsedEmail.issues?.map((issue) => issue.message) ?? [];

    setErrors((prev) => ({
      ...prev,
      email: emailIssues,
    }));
  };

  const reset = () => {
    setUserFormState({});
    setErrors({ username: [], email: [] });
  };

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!profileQuery.data) return;
    if (errors.username.length > 0) return;
    if (errors.email.length > 0) return;

    console.log(formData);
  };

  return {
    formData,
    isPending: profileQuery.isPending,
    handleUsernameChange,
    handleEmailChange,
    handleSubmit,
    errors,
    reset,
  };
};
