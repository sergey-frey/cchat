import { NAVIGATION } from "@/shared/navigation";
import { ReactNode, useState } from "react";
import { Redirect, Route, RouteProps } from "wouter";
import { useCheckAuth } from "./use-check-auth";

export const ProtectedRoute = ({ ...props }: RouteProps) => {
  const [routeComponent, setRouteComponent] = useState<ReactNode>(null);

  useCheckAuth({
    onSuccess: () => setRouteComponent(<Route {...props} />),
    onError: () =>
      setRouteComponent(
        <Redirect
          to={NAVIGATION.auth({
            searchParams: { state: "login" },
          })}
        />,
      ),
  });

  return routeComponent;
};
