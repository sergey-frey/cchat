import { Dispatch, SetStateAction, useState } from "react";

export type UseArrayStateMethods<T> = {
  /**
   * Set the state of the array
   */
  setState: Dispatch<SetStateAction<T[]>>;

  /**
   * Push the item to the array
   */
  push: (item: T) => void;

  /**
   * Push the item to the array if it is not already in the array
   */
  pushUnique: (item: T) => void;

  /**
   * Remove the item from the array
   */
  remove: (deleteItem: T) => void;

  /**
   * Remove the last item from the array
   */
  pop: () => void;

  /**
   * Clear the array
   */
  clear: () => void;
};

export const useArrayState = <T>(
  initialState: T[] | (() => T[]),
): [T[], UseArrayStateMethods<T>] => {
  const [state, setState] = useState(initialState);

  const push = (item: T) => {
    setState((prev) => [...prev, item]);
  };

  const remove = (deleteItem: T) => {
    setState((prev) => prev.filter((item) => item !== deleteItem));
  };

  const pushUnique = (item: T) => {
    setState((prev) => (prev.includes(item) ? prev : [...prev, item]));
  };

  const pop = () => {
    setState((prev) => {
      prev.pop();

      return [...prev];
    });
  };

  const clear = () => {
    setState([]);
  };

  return [state, { setState, push, pushUnique, pop, remove, clear }];
};
