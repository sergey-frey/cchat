import { HTMLAttributes } from "react";
import { twJoin } from "tailwind-merge";

type UserItemProps = HTMLAttributes<HTMLElement>;

export const UserItem = ({ className, ...props }: UserItemProps) => {
  return (
    <article
      {...props}
      className={twJoin("flex items-center gap-2", className)}
    />
  );
};

const Avatar = ({ className, ...props }: HTMLAttributes<HTMLElement>) => {
  return (
    <div
      {...props}
      className={twJoin(
        "w-10 h-10 rounded-full bg-primary-400",
        "md:w-14 md:h-14",
        className,
      )}
    />
  );
};

const Content = ({ ...props }: HTMLAttributes<HTMLElement>) => {
  return <div {...props} />;
};

const Name = ({ className, ...props }: HTMLAttributes<HTMLElement>) => {
  return (
    <h2
      {...props}
      className={twJoin("font-medium text-small", "md:text-medium", className)}
    />
  );
};

const Username = ({ className, ...props }: HTMLAttributes<HTMLElement>) => {
  return (
    <p
      {...props}
      className={twJoin("text-black/40 text-xs", "md:text-medium", className)}
    />
  );
};

UserItem.Avatar = Avatar;
UserItem.Content = Content;
UserItem.Name = Name;
UserItem.Username = Username;
