import { ReactNode } from "react";
import { useConfirmStore } from "./confirm-store";
import {
  setConfirmContentSelector,
  setConfirmSelector,
} from "./confirm-store-selectors";

type ConfirmCb = (closeConfirm: () => void) => void;
type UseConfirmOptions = {
  content?: ReactNode;
};

export const useConfirm = ({
  content = "Are you sure?",
}: UseConfirmOptions = {}) => {
  const setConfirm = useConfirmStore(setConfirmSelector);
  const setConfirmContent = useConfirmStore(setConfirmContentSelector);

  const closeConfirm = () => setConfirm(null);

  return (confirmCb: ConfirmCb) => {
    return () => {
      setConfirmContent(content);
      setConfirm(() => confirmCb(closeConfirm));
    };
  };
};
