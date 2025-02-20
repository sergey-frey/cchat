import { useUserStore } from "@/entities/user";
import { authService } from "@/features/auth/model/auth-service";
import { useEffect } from "react";

export const useCheckAuth = () => {
  const { user, setUser } = useUserStore();

  useEffect(() => {
    authService
      .checkSession()
      .then(setUser)
      .catch(() => setUser(null));
  }, [setUser]);

  return user;
};
