import { cn } from "@/shared/utils/cn";
import { Skeleton } from "@heroui/skeleton";
import { HTMLAttributes } from "react";
import { twJoin } from "tailwind-merge";

type UserItemProps = HTMLAttributes<HTMLElement>;
type PropsWithLoadedState<T> = { isLoaded?: boolean } & T;

export const UserItem = ({ className, ...props }: UserItemProps) => {
  return (
    <article
      {...props}
      className={twJoin("flex items-center gap-2", className)}
    />
  );
};

const Avatar = ({
  className,
  isLoaded,
  ...props
}: PropsWithLoadedState<HTMLAttributes<HTMLElement>>) => {
  return (
    <Skeleton isLoaded={isLoaded} className={cn("rounded-full", className)}>
      <div
        {...props}
        className={twJoin(
          "w-10 h-10 rounded-full bg-primary-400",
          "md:w-14 md:h-14",
          className,
        )}
      />
    </Skeleton>
  );
};

const Content = ({ ...props }: HTMLAttributes<HTMLElement>) => {
  return <div {...props} />;
};

const Name = ({
  className,
  isLoaded,
  ...props
}: PropsWithLoadedState<HTMLAttributes<HTMLElement>>) => {
  return (
    <Skeleton isLoaded={isLoaded} className={cn("rounded-small", className)}>
      <h2
        {...props}
        className={twJoin(
          "font-medium text-small",
          "md:text-medium",
          className,
        )}
      />
    </Skeleton>
  );
};

const Username = ({
  className,
  isLoaded,
  ...props
}: PropsWithLoadedState<HTMLAttributes<HTMLElement>>) => {
  return (
    <Skeleton isLoaded={isLoaded} className={cn("rounded-small", className)}>
      <p
        {...props}
        className={twJoin("text-black/40 text-xs", "md:text-medium", className)}
      />
    </Skeleton>
  );
};

UserItem.Avatar = Avatar;
UserItem.Content = Content;
UserItem.Name = Name;
UserItem.Username = Username;
