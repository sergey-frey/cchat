import { USER_QUERY_KEYS } from "@/entities/user";
import { NAVIGATION } from "@/shared/navigation";
import { queryClient } from "@/shared/query-client";
import { addToast } from "@heroui/toast";
import { useLocation } from "wouter";

export const useEditProfileHandlers = () => {
  const setLocation = useLocation()[1];

  const handleEditProfileSuccess = () => {
    queryClient.invalidateQueries({
      queryKey: [USER_QUERY_KEYS.PROFILE],
    });
    setLocation(NAVIGATION.profile);
    addToast({
      title: "Profile updated",
      description: "Your profile has been updated successfully",
      color: "success",
    });
  };

  const handleEditProfileError = () => {
    addToast({
      title: "Profile update",
      description: "Your profile has not been updated",
      color: "danger",
    });
  };

  return { handleEditProfileSuccess, handleEditProfileError };
};
