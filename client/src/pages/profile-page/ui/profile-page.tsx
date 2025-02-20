import { userSelector, useUserStore } from "@/entities/user";
import { authService, useLogout } from "@/features/auth";
import { BottomNavigation } from "@/features/navigation";
import { ArrowRightStartOnRectangleIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";

export const ProfilePage = () => {
  const user = useUserStore(userSelector);
  const logout = useLogout();

  if (!user) return null;

  const { username } = user;

  const handleLogoutClick = () => {
    authService.logout().then(logout);
  };

  return (
    <>
      <section className="p-4">
        <h1 className="text-2xl text-indigo-500 text-center mt-10">
          {username}
        </h1>

        <div className="mt-4">
          <Button
            isIconOnly
            color="danger"
            title="Logout"
            onPress={handleLogoutClick}
          >
            <ArrowRightStartOnRectangleIcon className="w-5 h-5" />
          </Button>
        </div>
      </section>

      <BottomNavigation className="fixed bottom-4 left-1/2 -translate-x-1/2" />
    </>
  );
};
