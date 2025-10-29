import { HTMLAttributes } from "react";
import { twJoin } from "tailwind-merge";

type UserBadgeProps = HTMLAttributes<HTMLElement>;

export const UserBadge = ({ className, ...props }: UserBadgeProps) => {
  return (
    <span
      {...props}
      className={twJoin(
        "flex gap-1 items-center",
        "border rounded-full p-1 pr-2",
        className,
      )}
    />
  );
};

const Avatar = ({ className, ...props }: HTMLAttributes<HTMLElement>) => {
  return (
    <span
      {...props}
      className={twJoin("w-6 h-6 rounded-full", "bg-primary-400", className)}
    />
  );
};

const AvatarPlaceholder = ({
  className,
  ...props
}: HTMLAttributes<HTMLElement>) => {
  return (
    <span {...props} className={twJoin("w-6 h-6 rounded-full", className)} />
  );
};

const Username = ({ className, ...props }: HTMLAttributes<HTMLElement>) => {
  return (
    <span
      {...props}
      className={twJoin("font-medium text-black/40 text-xs", className)}
    />
  );
};

const EndContent = ({ className, ...props }: HTMLAttributes<HTMLElement>) => {
  return <span {...props} className={twJoin("text-neutral-600", className)} />;
};

UserBadge.Avatar = Avatar;
UserBadge.AvatarPlaceholder = AvatarPlaceholder;
UserBadge.Username = Username;
UserBadge.EndContent = EndContent;
