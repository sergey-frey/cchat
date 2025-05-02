// Dependencies
import {
  ChatBubbleOvalLeftEllipsisIcon,
  Cog6ToothIcon,
  RssIcon,
  UserIcon,
} from "@heroicons/react/24/outline";
import { Button, ButtonGroup, ButtonProps } from "@heroui/button";
import { motion } from "framer-motion";
import { HTMLAttributes } from "react";
import { Link } from "wouter";
import { usePathname } from "wouter/use-browser-location";

import { CONFIG } from "@/shared/config";
import { NAVIGATION } from "@/shared/navigation";
import { cn } from "@/shared/utils/cn";
import { absolutePath } from "@/shared/utils/path";

const BOTTOM_NAVIGATION_PATHS = [
  { path: NAVIGATION.profile, icon: UserIcon, enabled: true },
  {
    path: NAVIGATION.chats(),
    icon: ChatBubbleOvalLeftEllipsisIcon,
    enabled: true,
  },
  { path: NAVIGATION.channels(), icon: RssIcon, enabled: false },
  { path: NAVIGATION.settings(), icon: Cog6ToothIcon, enabled: false },
];

const Y_DISTANCE = 100;

export const BottomNavigation = ({
  className,
}: HTMLAttributes<HTMLElement>) => {
  const pathname = usePathname();

  return (
    <motion.nav
      className="fixed left-0 bottom-2 w-full flex items-center justify-center"
      initial={{ y: Y_DISTANCE }}
      animate={{ y: 0 }}
      transition={{ duration: CONFIG.ANIMATION_TRANSITION.BASE }}
    >
      <ButtonGroup className={cn("bg-white", className)}>
        {BOTTOM_NAVIGATION_PATHS.map(({ path, icon, enabled }, i) => {
          const Icon = icon;
          const isActive = pathname.startsWith(absolutePath(path));
          const color: ButtonProps["color"] = isActive ? "primary" : "default";
          const isDisabled = enabled === false;

          return (
            <Button
              key={i}
              as={Link}
              isIconOnly
              size="lg"
              href={path}
              color={color}
              variant="flat"
              isDisabled={isDisabled}
            >
              <Icon className="w-6 h-6" />
            </Button>
          );
        })}
      </ButtonGroup>
    </motion.nav>
  );
};
