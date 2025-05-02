import { NAVIGATION } from "@/shared/navigation";
import { PlusCircleIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { useState } from "react";
import { Link } from "wouter";
import { ChatsSearch } from "./chats-search";

export const ChatsPage = () => {
  const [search, setSearch] = useState("");

  const handleSearchSubmit = (searchValue: string) => {
    setSearch(searchValue);
  };

  return (
    <section className="p-4 pt-14">
      <div className="flex gap-2 items-center">
        <ChatsSearch onSubmitSearch={handleSearchSubmit} className="grow" />

        <Link
          asChild
          href={NAVIGATION.createChat}
          state={{ origin: NAVIGATION.chats() }}
        >
          <Button as={Link} variant="flat" isIconOnly color="primary">
            <PlusCircleIcon className="w-6 h-6" />
          </Button>
        </Link>
      </div>
    </section>
  );
};
