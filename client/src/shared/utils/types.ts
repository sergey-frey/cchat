import { RefObject } from "react";
import { RefCallBack } from "react-hook-form";

export type RefProp<T> = { elemRef?: RefObject<T> | RefCallBack };
