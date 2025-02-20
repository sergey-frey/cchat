import { useUsersQuery } from "@/entities/user";
import { BottomNavigation } from "@/features/navigation";
import { ChevronLeftIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import { FormEvent } from "react";
import { Link, useLocation } from "wouter";

export const CreateChatPage = () => {
  const { data } = useUsersQuery();

  const { state } = useLocation();

  const handleUsersSearchSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const form = e.currentTarget;

    if (!form) return;

    const name = new FormData(form).get("name")?.toString();
  };

  return (
    <>
      <section className="p-4">
        <div className="flex gap-2">
          <Button as={Link} isIconOnly variant="flat" href={state?.origin}>
            <ChevronLeftIcon className="w-5 h-5" />
          </Button>

          <form className="grow">
            <Input placeholder="Find users..." name="name" />
          </form>
        </div>

        <ul>{data?.map((user) => <li key={user.id}>{user.name}</li>)}</ul>
      </section>

      <BottomNavigation className="fixed bottom-4 left-1/2 -translate-x-1/2" />
    </>
  );
};
