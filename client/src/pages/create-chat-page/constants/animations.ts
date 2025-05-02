import { CONFIG } from "@/shared/config";
import { AnimationsCollection } from "@/shared/types/animation";

const HEADER_Y_DISTANCE = 100;
const SCROLL_TOP_BUTTON_X_DISTANCE = 100;

export const CREATE_CHAT_PAGE_ANIMATIONS: AnimationsCollection<
  "HEADER" | "SCROLL_TOP_BUTTON"
> = {
  HEADER: {
    initial: {
      y: -HEADER_Y_DISTANCE,
    },
    animate: {
      y: 0,
    },
    transition: {
      duration: CONFIG.ANIMATION_TRANSITION.BASE,
    },
  },

  SCROLL_TOP_BUTTON: {
    initial: {
      x: SCROLL_TOP_BUTTON_X_DISTANCE,
    },
    animate: {
      x: 0,
    },
    transition: {
      duration: CONFIG.ANIMATION_TRANSITION.BASE,
    },
  },
};
