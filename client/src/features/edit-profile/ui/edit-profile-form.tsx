import { cn } from "@/shared/utils/cn";
import { CheckCircleIcon, XCircleIcon } from "@heroicons/react/24/outline";
import { Input } from "@heroui/input";
import { Skeleton } from "@heroui/skeleton";
import { FormHTMLAttributes } from "react";
import { useEditProfileForm } from "../model/use-edit-profile-form";
import { Button } from "@heroui/button";

type EditProfileFormProps = FormHTMLAttributes<HTMLFormElement> & {};

export const EditProfileForm = ({
  className,
  ...props
}: EditProfileFormProps) => {
  const {
    formData,
    handleUsernameChange,
    handleEmailChange,
    handleSubmit,
    isPending,
    errors,
    reset,
  } = useEditProfileForm();

  const isValidUsername = !errors.username.length;
  const isValidEmail = !errors.email.length;
  const isSubmitDisabled = !isValidUsername || !isValidEmail || isPending;

  const usernameInputValidIcon = isValidUsername ? (
    <CheckCircleIcon className="w-6 h-6 text-green-500" />
  ) : (
    <XCircleIcon className="w-6 h-6 text-red-500" />
  );

  return (
    <form {...props} onSubmit={handleSubmit} className={cn("", className)}>
      <div className="flex flex-col gap-3">
        <Skeleton className="rounded-medium" isLoaded={!isPending}>
          <Input
            label="Username"
            value={formData.username}
            onValueChange={handleUsernameChange}
            isInvalid={!isValidUsername}
            errorMessage={
              <ul>
                {errors.username.map((error, i) => (
                  <li key={i}>{error}</li>
                ))}
              </ul>
            }
            startContent={
              <span className="text-neutral-600 font-medium">{"@"}</span>
            }
            endContent={usernameInputValidIcon}
          />
        </Skeleton>

        <Skeleton className="rounded-medium" isLoaded={!isPending}>
          <Input
            label="Email"
            value={formData.email}
            onValueChange={handleEmailChange}
            isInvalid={!isValidEmail}
            errorMessage={
              <ul>
                {errors.email.map((error, i) => (
                  <li key={i}>{error}</li>
                ))}
              </ul>
            }
          />
        </Skeleton>
      </div>

      <div className="mt-4 grid grid-cols-3 gap-2">
        <Button
          type="submit"
          color="primary"
          className="col-span-2"
          isDisabled={isSubmitDisabled}
        >
          Save
        </Button>
        <Button type="button" color="default" onPress={reset}>
          Clear
        </Button>
      </div>
    </form>
  );
};
