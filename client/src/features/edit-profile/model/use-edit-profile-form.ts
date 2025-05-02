import { useProfileQuery, useUpdateProfileQuery } from "@/entities/user";
import debounce from "debounce";
import { FormEvent, useCallback, useEffect, useMemo, useState } from "react";
import { safeParseAsync } from "valibot";
import { EditProfileFormSchema, emailSchema, usernameSchema } from "./schemas";

const initialFormState: EditProfileFormSchema = {
  username: "",
  email: "",
};

export type ValidationErrors = Record<keyof EditProfileFormSchema, string[]>;

type UseEditProfileFormOptions = {
  onSuccess?: () => void;
  onError?: (error?: unknown) => void;
};

export const useEditProfileForm = ({
  onSuccess,
  onError,
}: UseEditProfileFormOptions = {}) => {
  const [userFormState, setUserFormState] = useState<
    Partial<EditProfileFormSchema>
  >({});
  const [errors, setErrors] = useState<ValidationErrors>({
    username: [],
    email: [],
  });
  const [isFormDirty, setIsFormDirty] = useState(false);

  const profileQuery = useProfileQuery({
    staleTime: 0,
    gcTime: 0,
    refetchOnMount: "always",
    refetchOnWindowFocus: true,
  });

  const updateProfileMutation = useUpdateProfileQuery();

  // Memoize form data to prevent unnecessary re-renders
  const formData = useMemo(
    () => ({
      ...initialFormState,
      ...profileQuery.data,
      ...userFormState,
    }),
    [profileQuery.data, userFormState],
  );

  // Sync profile data when loaded/updated
  useEffect(() => {
    if (profileQuery.data && !isFormDirty) {
      setUserFormState({});
    }
  }, [profileQuery.data, isFormDirty]);

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
    setIsFormDirty(true);

    checkUsername(username, (errors) => {
      setErrors((prev) => ({
        ...prev,
        username: errors,
      }));
    });
  };

  const handleEmailChange = async (email: string) => {
    setUserFormState((prev) => ({ ...prev, email }));
    setIsFormDirty(true);

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
    setIsFormDirty(false);
  };

  const hasErrors = Object.values(errors).some(
    (fieldErrors) => fieldErrors.length > 0,
  );
  const hasChanges =
    isFormDirty &&
    (formData.username !== profileQuery.data?.username ||
      formData.email !== profileQuery.data?.email);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!profileQuery.data || hasErrors || !hasChanges) return;

    try {
      await updateProfileMutation.mutateAsync(formData);
      setIsFormDirty(false);
      onSuccess?.();
    } catch (error) {
      onError?.(error);
    }
  };

  return {
    formData,
    fetchingQueryState: {
      isPending: profileQuery.isPending,
    },
    updatingMutationState: {
      isPending: updateProfileMutation.isPending,
    },
    handleUsernameChange,
    handleEmailChange,
    handleSubmit,
    errors,
    reset,
    hasErrors,
    hasChanges,
    isFormDirty,
  };
};
