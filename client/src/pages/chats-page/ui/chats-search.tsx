import { Input, InputProps } from "@heroui/input";
import { FormEvent } from "react";

type ChatsSearchProps = InputProps & {
  onSubmitSearch: (search: string) => void;
};

export const ChatsSearch = ({
  onSubmitSearch,
  className,
  ...props
}: ChatsSearchProps) => {
  const handleSearchSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);

    onSubmitSearch(formData.get("search")?.toString() ?? "");
  };

  return (
    <form onSubmit={handleSearchSubmit} className={className}>
      <Input
        {...props}
        name="search"
        placeholder="Search..."
        className="bg-transparent"
      />
    </form>
  );
};
