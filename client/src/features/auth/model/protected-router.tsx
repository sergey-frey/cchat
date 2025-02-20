import { NAVIGATION } from "@/shared/navigation";
import { Redirect, Router, RouterProps } from "wouter";
import { useCheckAuth } from "./use-check-auth";

export const ProtectedRouter = ({ ...props }: RouterProps) => {
  const authCheckResponse = useCheckAuth();

  if (authCheckResponse === null) {
    return <Redirect to={NAVIGATION.auth("login")} />;
  }

  if (authCheckResponse) {
    return <Router {...props} />;
  }
};
