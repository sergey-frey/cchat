import { CONFIG } from "@/shared/config";
import { AnimationsCollection } from "@/shared/types/animation";

const HEADER_Y_DISTANCE = 100;
const SCROLL_TOP_BUTTON_X_DISTANCE = 100;
const FETCH_USERS_ERROR_MESSAGE_Y_DISTANCE = 30;
const CREATE_CHAT_BUTTON_X_DISTANCE = 30;

export const CREATE_CHAT_PAGE_ANIMATIONS: AnimationsCollection<
  | "HEADER"
  | "SCROLL_TOP_BUTTON"
  | "FETCH_USERS_ERROR_MESSAGE"
  | "CHAT_MEMBER_BADGE"
  | "CREATE_CHAT_BUTTON"
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
  FETCH_USERS_ERROR_MESSAGE: {
    initial: {
      y: FETCH_USERS_ERROR_MESSAGE_Y_DISTANCE,
      opacity: 0,
    },
    animate: {
      y: 0,
      opacity: 1,
    },
    transition: {
      duration: CONFIG.ANIMATION_TRANSITION.BASE,
    },
  },
  CHAT_MEMBER_BADGE: {
    initial: {
      scale: 0.6,
      opacity: 0,
    },
    animate: {
      scale: 1,
      opacity: 1,
    },
    transition: {
      duration: CONFIG.ANIMATION_TRANSITION.FAST,
    },
  },
  CREATE_CHAT_BUTTON: {
    initial: {
      x: -CREATE_CHAT_BUTTON_X_DISTANCE,
      opacity: 0,
    },
    animate: {
      x: 0,
      opacity: 1,
    },
    transition: {
      duration: CONFIG.ANIMATION_TRANSITION.FAST,
    },
  },
};
