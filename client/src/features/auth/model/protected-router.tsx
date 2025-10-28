import { NAVIGATION } from "@/shared/navigation";
import { ReactNode, useState } from "react";
import { Redirect, Router, RouterProps } from "wouter";
import { useCheckAuth } from "./use-check-auth";

export const ProtectedRouter = ({ ...props }: RouterProps) => {
  const [routeComponent, setRouteComponent] = useState<ReactNode>(null);

  useCheckAuth({
    onSuccess: () => setRouteComponent(<Router {...props} />),
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
