// Dependencies
import { ArrowUpIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { motion } from "framer-motion";

import { cn } from "@/shared/utils/cn";
import { CREATE_CHAT_PAGE_ANIMATIONS } from "../constants/animations";

type ScrollTopButtonProps = { onPress: () => void; className?: string };

export const ScrollTopButton = ({
  onPress,
  className,
}: ScrollTopButtonProps) => {
  return (
    <motion.div
      className={cn("fixed bottom-10 right-5 z-10", className)}
      {...CREATE_CHAT_PAGE_ANIMATIONS.SCROLL_TOP_BUTTON}
    >
      <Button color="primary" variant="flat" isIconOnly onPress={onPress}>
        <ArrowUpIcon className="w-5 h-5" />
      </Button>
    </motion.div>
  );
};
