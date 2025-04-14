import { authService } from "@/entities/auth/model/auth-service";
import { useEffect } from "react";

type UseCheckAuthOptions = {
  onSuccess: () => void;
  onError: () => void;
};

export const useCheckAuth = ({ onSuccess, onError }: UseCheckAuthOptions) => {
  useEffect(() => {
    authService.checkSession().then(onSuccess).catch(onError);
  }, []);
};
