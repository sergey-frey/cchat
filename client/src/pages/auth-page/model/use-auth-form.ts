import { authService } from "@/features/auth/model/auth-service";
import { NAVIGATION } from "@/shared/navigation";
import { valibotResolver } from "@hookform/resolvers/valibot";
import { useForm } from "react-hook-form";
import { useLocation } from "wouter";
import { AuthFormState } from "../types";
import { AuthFormSchema, AuthFormSchemaType } from "./auth-form-schema";

const authMethodMap: Record<AuthFormState, keyof typeof authService> = {
  reg: "register",
  login: "login",
};

export const useAuthForm = (formState: AuthFormState) => {
  const {
    handleSubmit,
    control,
    formState: authFormState,
    getValues,
  } = useForm<AuthFormSchemaType>({
    defaultValues: {
      email: "",
      password: "",
    },

    mode: "onBlur",
    resolver: valibotResolver(AuthFormSchema),
  });

  console.log(authFormState.errors, authFormState.isValid, getValues());

  const setLocation = useLocation()[1];

  const authMethod = authMethodMap[formState];

  const onSubmit = handleSubmit((data) => {
    if (!authFormState.isValid) return;

    authService[authMethod](data).then(() => {
      setLocation(NAVIGATION.profile);
    });
  });

  return {
    onSubmit,
    control,
    errors: authFormState.errors,
    isValid: authFormState.isValid,
  };
};
