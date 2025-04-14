import { NAVIGATION } from "@/shared/navigation";
import { Redirect, Route, RouteProps } from "wouter";
import { useCheckAuth } from "./use-check-auth";

export const ProtectedRoute = ({ ...props }: RouteProps) => {
  const isAuthenticated = useCheckAuth();

  if (isAuthenticated === null) {
    return (
      <Redirect
        to={NAVIGATION.auth({
          searchParams: { state: "login" },
        })}
      />
    );
  }

  if (isAuthenticated) {
    return <Route {...props} />;
  }
};
