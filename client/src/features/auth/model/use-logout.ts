import { authService } from "@/entities/auth/model/auth-service";
import { NAVIGATION } from "@/shared/navigation";
import { useLocation } from "wouter";

export const useLogout = () => {
  const setLocation = useLocation()[1];

  return () =>
    authService.logout().finally(() => {
      setLocation(NAVIGATION.auth({ searchParams: { state: "login" } }), {
        replace: true,
      });
    });
};
