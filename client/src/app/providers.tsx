import { queryClient } from "@/shared/query-client";
import { HeroUIProvider } from "@heroui/system";
import { QueryClientProvider } from "@tanstack/react-query";
import { ReactNode } from "react";

export const Providers = ({ children }: { children: ReactNode }) => {
  return (
    <QueryClientProvider client={queryClient}>
      <HeroUIProvider className="contents">{children}</HeroUIProvider>
    </QueryClientProvider>
  );
};
