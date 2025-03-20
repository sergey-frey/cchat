import { userSelector, useUserStore } from "@/entities/user";
import { authService, useLogout } from "@/features/auth";
import { BottomNavigation } from "@/features/navigation";
import {
  ArrowRightStartOnRectangleIcon,
  PencilSquareIcon,
} from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Card, CardBody, CardHeader } from "@heroui/card";
import { Divider } from "@heroui/divider";

export const ProfilePage = () => {
  const user = useUserStore(userSelector);
  const logout = useLogout();

  if (!user) return null;

  const { username, email } = user;

  const handleLogoutClick = () => {
    authService.logout().then(logout);
  };

  return (
    <>
      <section className="h-full flex items-center justify-center p-4">
        <Card className="w-full max-w-[340px] mx-auto">
          <CardHeader className="text-primary-400">
            {username}

            <ul className="ml-auto flex gap-2">
              <li>
                <Button isIconOnly size="sm" variant="flat">
                  <PencilSquareIcon className="w-4 h-4" />
                </Button>
              </li>

              <li>
                <Button
                  isIconOnly
                  size="sm"
                  color="danger"
                  title="Logout"
                  onPress={handleLogoutClick}
                >
                  <ArrowRightStartOnRectangleIcon className="w-4 h-4" />
                </Button>
              </li>
            </ul>
          </CardHeader>

          <Divider />

          <CardBody>
            <ul>
              <li>{email}</li>
            </ul>
          </CardBody>
        </Card>
      </section>

      <BottomNavigation className="fixed bottom-4 left-1/2 -translate-x-1/2" />
    </>
  );
};
