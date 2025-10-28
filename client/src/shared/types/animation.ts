import { TargetAndTransition, Transition } from "framer-motion";

export type AnimationsCollection<Keys extends string> = Record<
  Keys,
  {
    initial?: TargetAndTransition;
    animate?: TargetAndTransition;
    exit?: TargetAndTransition;
    transition?: Transition;
  }
>;
