import { ReactNode, useState } from "react";
import { Redirect, Route, RouteProps } from "wouter";
import { useCheckAuth } from "./use-check-auth";

type UnauthorizedRouteProps = RouteProps & {
  redirectPath: string;
};

export const UnauthorizedRoute = ({
  redirectPath,
  ...props
}: UnauthorizedRouteProps) => {
  const [routeComponent, setRouteComponent] = useState<ReactNode>(null);

  useCheckAuth({
    onSuccess: () => setRouteComponent(<Redirect to={redirectPath} />),
    onError: () => setRouteComponent(<Route {...props} />),
  });

  return routeComponent;
};
