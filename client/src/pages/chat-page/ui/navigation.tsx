import { NAVIGATION } from "@/shared/navigation";
import { cn } from "@/shared/utils/cn";
import { ChevronLeftIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Link } from "wouter";

type NavigationProps = {
  chatTitle: string;
};

export const Navigation = ({ chatTitle }: NavigationProps) => {
  return (
    <nav
      className={cn(
        "sticky top-0 shadow px-3 py-2 bg-white rounded-b-medium",
        "grid grid-cols-3 items-center",
      )}
    >
      <Button
        as={Link}
        isIconOnly
        variant="light"
        size="sm"
        href={NAVIGATION.chats()}
      >
        <ChevronLeftIcon className="w-5 h-5" />
      </Button>

      <h1 className="text-center text-indigo-400 font-semibold">{chatTitle}</h1>
    </nav>
  );
};
