import { useUserStore, setUserSelector } from "@/entities/user";
import { authService } from "./auth-service";
import { useLocation } from "wouter";
import { NAVIGATION } from "@/shared/navigation";

export const useLogout = () => {
  const setUser = useUserStore(setUserSelector);
  const setLocation = useLocation()[1];

  return () =>
    authService.logout().finally(() => {
      setUser(null);
      setLocation(NAVIGATION.auth("login"), { replace: true });
    });
};
