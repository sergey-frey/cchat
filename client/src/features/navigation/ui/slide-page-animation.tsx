// Dependencies
import { motion } from "framer-motion";
import { ReactNode } from "react";
import { useLocation } from "wouter";

import { CONFIG } from "@/shared/config";

type SlidePageAnimationProps = { children: ReactNode };

const X_DISTANCE = 50;
export const SlidePageAnimation = ({ children }: SlidePageAnimationProps) => {
  const [location] = useLocation();

  return (
    <motion.div
      key={location}
      className="bg-background"
      style={{ position: "absolute", width: "100%", height: "100%" }}
      initial={{ opacity: 0, x: X_DISTANCE }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 1, x: -X_DISTANCE }}
      transition={{ duration: CONFIG.ANIMATION_TRANSITION.BASE }}
    >
      {children}
    </motion.div>
  );
};
