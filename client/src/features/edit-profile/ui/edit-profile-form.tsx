import { cn } from "@/shared/utils/cn";
import { CheckCircleIcon, XCircleIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import { Skeleton } from "@heroui/skeleton";
import { FormHTMLAttributes } from "react";
import { useEditProfileForm } from "../model/use-edit-profile-form";

type EditProfileFormProps = FormHTMLAttributes<HTMLFormElement> & {
  onSuccess?: () => void;
  onError?: () => void;
};

export const EditProfileForm = ({
  className,
  onSuccess,
  onError,
  ...props
}: EditProfileFormProps) => {
  const {
    formData,
    handleUsernameChange,
    handleEmailChange,
    handleSubmit,
    fetchingQueryState,
    updatingMutationState,
    errors,
    reset,
  } = useEditProfileForm({
    onSuccess,
    onError,
  });

  const isValidUsername = !errors.username.length;
  const isValidEmail = !errors.email.length;
  const isInputsDisabled = updatingMutationState.isPending;
  const isSubmitDisabled =
    !isValidUsername ||
    !isValidEmail ||
    fetchingQueryState.isPending ||
    updatingMutationState.isPending;

  const usernameInputValidIcon = isValidUsername ? (
    <CheckCircleIcon className="w-6 h-6 text-green-500" />
  ) : (
    <XCircleIcon className="w-6 h-6 text-red-500" />
  );

  return (
    <form {...props} onSubmit={handleSubmit} className={cn("", className)}>
      <div className="flex flex-col gap-3">
        <Skeleton
          className="rounded-medium"
          isLoaded={!fetchingQueryState.isPending}
        >
          <Input
            label="Username"
            value={formData.username}
            onValueChange={handleUsernameChange}
            isInvalid={!isValidUsername}
            isDisabled={isInputsDisabled}
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

        <Skeleton
          className="rounded-medium"
          isLoaded={!fetchingQueryState.isPending}
        >
          <Input
            label="Email"
            value={formData.email}
            onValueChange={handleEmailChange}
            isInvalid={!isValidEmail}
            isDisabled={isInputsDisabled}
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
          Reset
        </Button>
      </div>
    </form>
  );
};
