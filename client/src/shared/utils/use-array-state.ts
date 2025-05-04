import { Dispatch, SetStateAction, useState } from "react";

export type UseArrayStateMethods<T> = {
  setState: Dispatch<SetStateAction<T[]>>;
  push: (item: T) => void;
  pushUnique: (item: T) => void;
  remove: (deleteItem: T) => void;
  pop: () => void;
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
