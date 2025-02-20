import { useSearchParams } from "@/shared/utils/use-search-params";
import { AuthFormState } from "../types";

const formTitle: Record<AuthFormState, string> = {
  reg: "Sign Up",
  login: "Sign In",
};

const changeFormLinkText: Record<AuthFormState, string> = {
  reg: "Sign In",
  login: "Sign Up",
};

const helperText: Record<AuthFormState, string> = {
  reg: "Already have an account?",
  login: "Don't have an account?",
};

const changeFormLinkTarget: Record<AuthFormState, string> = {
  reg: "login",
  login: "reg",
};

export const useAuthFormState = () => {
  const search = useSearchParams();
  const formState = (search.get("state") ?? "reg") as AuthFormState;

  return {
    formState,
    formTitle: formTitle[formState],
    changeFormLinkText: changeFormLinkText[formState],
    helperText: helperText[formState],
    changeFormLinkTarget: changeFormLinkTarget[formState],
  };
};
