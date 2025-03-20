import { NAVIGATION } from "@/shared/navigation";
import { Redirect, Route, RouteProps } from "wouter";
import { useCheckAuth } from "./use-check-auth";

export const ProtectedRoute = ({ ...props }: RouteProps) => {
  const authCheckResponse = useCheckAuth();

  if (authCheckResponse === null) {
    return (
      <Redirect to={NAVIGATION.auth({ searchParams: { state: "login" } })} />
    );
  }

  if (authCheckResponse) {
    return <Route {...props} />;
  }
};
