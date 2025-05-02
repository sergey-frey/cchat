import { useLocation } from "wouter";

export const useLocationState = <T extends object>() => {
  return useLocation() as unknown as { state: T | undefined };
};
