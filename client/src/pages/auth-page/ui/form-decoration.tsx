import { cn } from "@/shared/utils/cn";
import { HTMLAttributes } from "react";

type FormDecorationProps = HTMLAttributes<HTMLDivElement>;

export const FormDecoration = ({
  className,
  ...props
}: FormDecorationProps) => {
  return (
    <div
      {...props}
      className={cn("relative bg-primary-700 overflow-hidden", className)}
    >
      <p
        aria-hidden
        className={cn(
          "uppercase text-6xl font-black text-white tracking-[1.1rem] leading-snug",
          "absolute -left-4 top-10 rotate-6",
        )}
      >
        Welcome to CChat
      </p>
    </div>
  );
};
