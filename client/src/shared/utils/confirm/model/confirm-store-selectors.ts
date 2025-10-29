import { UseConfirmStore } from "./confirm-store";

export const confirmSelector = (state: UseConfirmStore) => state.confirm;
export const confirmContentSelector = (state: UseConfirmStore) =>
  state.confirmContent;
export const setConfirmSelector = (state: UseConfirmStore) => state.setConfirm;
export const setConfirmContentSelector = (state: UseConfirmStore) =>
  state.setConfirmContent;
