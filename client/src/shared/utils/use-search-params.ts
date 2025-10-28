import { useSearch } from "wouter";

export const useSearchParams = () => {
  const search = useSearch();
  return new URLSearchParams(search);
};
