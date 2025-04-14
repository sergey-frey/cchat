import { queryClient } from "@/shared/query-client";
import { HeroUIProvider } from "@heroui/system";
import { ToastProvider } from "@heroui/toast";
import { QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";

export const Providers = ({ children }: { children: ReactNode }) => {
  return (
    <QueryClientProvider client={queryClient}>
      <HeroUIProvider className="contents">
        <ToastProvider />
        {children}
      </HeroUIProvider>
    </QueryClientProvider>
  );
};
