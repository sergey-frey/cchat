/* dependencies */
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import { Controller } from "react-hook-form";
import { Link } from "wouter";

import { useAuthForm } from "../model/use-auth-form";
import { useAuthFormState } from "../model/use-auth-form-state";
import { cn } from "@/shared/utils/cn";
import { FormDecoration } from "./form-decoration";
import { NAVIGATION } from "@/shared/navigation";

export const AuthPage = () => {
  const {
    formState,
    formTitle,
    changeFormLinkText,
    helperText,
    changeFormLinkTarget,
  } = useAuthFormState();

  const { onSubmit, control, isValid } = useAuthForm(formState);

  return (
    <section className="p-4">
      <div
        className={cn(
          "grid max-w-md mx-auto",
          "mt-10",
          "md:grid-cols-2 md:max-w-2xl md:shadow",
        )}
      >
        <FormDecoration className="hidden md:block" />

        <form onSubmit={onSubmit} className={cn("grid items-center", "md:p-7")}>
          <h1 className="text-xl">CChat | {formTitle}</h1>

          <div className="mt-4 grid gap-3">
            <Controller
              name="email"
              control={control}
              render={({ field: { ...field }, fieldState: { error } }) => {
                return (
                  <Input
                    type="text"
                    label="Email"
                    isInvalid={Boolean(error)}
                    errorMessage={error?.message}
                    {...field}
                  />
                );
              }}
            />

            <Controller
              name="password"
              control={control}
              render={({ field: { ...field }, fieldState: { error } }) => {
                return (
                  <Input
                    type="password"
                    label="Password"
                    isInvalid={Boolean(error)}
                    errorMessage={error?.message}
                    {...field}
                  />
                );
              }}
            />
          </div>

          <Button
            type="submit"
            color="primary"
            isDisabled={!isValid}
            className={cn(
              "mt-4 justify-self-end mx-auto",
              "w-[80%]",
              "md:mt-8",
            )}
          >
            Submit
          </Button>

          <p className="mt-3 text-sm text-center">
            {helperText}{" "}
            <Link
              to={NAVIGATION.auth({
                searchParams: {
                  state: changeFormLinkTarget,
                },
              })}
              className={"underline underline-offset-2"}
            >
              {changeFormLinkText}
            </Link>
          </p>
        </form>
      </div>
    </section>
  );
};
