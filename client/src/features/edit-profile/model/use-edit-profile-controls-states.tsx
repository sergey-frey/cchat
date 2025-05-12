import { CheckCircleIcon, XCircleIcon } from "@heroicons/react/24/outline";
import { ValidationErrors } from "./use-edit-profile-form";

type UseEditProfileControlsStatesOptions = {
  errors: ValidationErrors;
  updatingMutationState: {
    isPending: boolean;
  };
  fetchingQueryState: {
    isPending: boolean;
  };
};

export const useEditProfileControlsStates = ({
  errors,
  updatingMutationState,
  fetchingQueryState,
}: UseEditProfileControlsStatesOptions) => {
  const isValidUsername = !errors.username.length;
  const isValidEmail = !errors.email.length;
  const isValidName = !errors.name.length;
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

  return {
    isInputsDisabled,
    isSubmitDisabled,
    isValidEmail,
    isValidUsername,
    isValidName,
    usernameInputValidIcon,
  };
};
