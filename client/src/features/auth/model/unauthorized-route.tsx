import { Redirect, Route, RouteProps } from "wouter";
import { useCheckAuth } from "./use-check-auth";

type UnauthorizedRouteProps = RouteProps & {
  redirectPath: string;
};

export const UnauthorizedRoute = ({
  redirectPath,
  ...props
}: UnauthorizedRouteProps) => {
  const authCheckResponse = useCheckAuth();

  if (authCheckResponse === undefined) {
    return null;
  }

  if (authCheckResponse) {
    return <Redirect to={redirectPath} />;
  }

  return <Route {...props} />;
};
