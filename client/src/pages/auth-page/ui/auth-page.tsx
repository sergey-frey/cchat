/* dependencies */
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import { Controller } from "react-hook-form";
import { Link } from "wouter";

import { useAuthForm } from "../model/use-auth-form";
import { useAuthFormState } from "../model/use-auth-form-state";

export const AuthPage = () => {
  const {
    formState,
    formTitle,
    changeFormLinkText,
    helperText,
    changeFormLinkTarget,
  } = useAuthFormState();

  const { onSubmit, control } = useAuthForm(formState);

  return (
    <section className="p-4">
      <form onSubmit={onSubmit}>
        <h1 className="text-xl">{formTitle}</h1>

        <div className="mt-4 grid gap-3">
          <Controller
            name="email"
            control={control}
            render={({ field: { ...field } }) => {
              return <Input type="text" label="Email" {...field} />;
            }}
          />

          <Controller
            name="password"
            control={control}
            render={({ field: { ...field } }) => {
              return <Input type="password" label="Password" {...field} />;
            }}
          />
        </div>

        <Button type="submit" className="mt-4 justify-self-end">
          Submit
        </Button>

        <p className="mt-3 text-sm text-center">
          {helperText}{" "}
          <Link
            to={`/auth?state=${changeFormLinkTarget}`}
            className={"underline underline-offset-2"}
          >
            {changeFormLinkText}
          </Link>
        </p>
      </form>
    </section>
  );
};
