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
  const { handleSubmit, ...rest } = useForm<AuthFormSchemaType>({
    defaultValues: {
      email: "",
      password: "",
    },

    resolver: valibotResolver(AuthFormSchema),
  });

  const setLocation = useLocation()[1];

  const authMethod = authMethodMap[formState];

  const onSubmit = handleSubmit((data) => {
    authService[authMethod](data).then(() => {
      setLocation(NAVIGATION.profile);
    });
  });

  return { onSubmit, ...rest };
};
