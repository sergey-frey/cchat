import { IUser } from "@/entities/user";
import { HTMLAttributes, ReactNode } from "react";
import { twJoin } from "tailwind-merge";

type ChatMembersProps = Omit<HTMLAttributes<HTMLDivElement>, "children"> & {
  users: IUser[];
  endContent: ReactNode;
  children: (user: IUser) => ReactNode;
};

export const ChatMembers = ({
  users,
  endContent,
  className,
  children,
  ...props
}: ChatMembersProps) => {
  return (
    <div {...props} className={twJoin("flex justify-between gap-2", className)}>
      <ul className="flex gap-1 flex-wrap">
        {users.map((user) => {
          return <li key={user.id}>{children?.(user)}</li>;
        })}
      </ul>

      <div>{endContent}</div>
    </div>
  );
};
