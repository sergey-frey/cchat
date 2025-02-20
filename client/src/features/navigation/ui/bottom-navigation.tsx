import { NAVIGATION } from "@/shared/navigation";
import { cn } from "@/shared/utils/cn";
import { absolutePath } from "@/shared/utils/path";
import {
  ChatBubbleOvalLeftEllipsisIcon,
  UserIcon,
} from "@heroicons/react/24/outline";
import { Tab, Tabs } from "@heroui/tabs";
import { HTMLAttributes } from "react";
import { Link } from "wouter";
import { usePathname } from "wouter/use-browser-location";

const BOTTOM_NAVIGATION_PATHS = [
  { path: NAVIGATION.profile, icon: UserIcon },
  { path: NAVIGATION.chats(), icon: ChatBubbleOvalLeftEllipsisIcon },
];

export const BottomNavigation = ({
  className,
}: HTMLAttributes<HTMLElement>) => {
  const pathname = usePathname();

  return (
    <Tabs className={cn(className)} selectedKey={pathname}>
      {BOTTOM_NAVIGATION_PATHS.map(({ path, icon }) => {
        const Icon = icon;

        return (
          <Tab
            as={Link}
            key={absolutePath(path)}
            href={path}
            title={<Icon className="w-6 h-6" />}
          />
        );
      })}
    </Tabs>
  );
};
