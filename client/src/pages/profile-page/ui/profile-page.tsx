import { useProfileQuery } from "@/entities/user";
import { authService, useLogout } from "@/features/auth";
import { NAVIGATION } from "@/shared/navigation";
import { useConfirm } from "@/shared/utils/confirm";
import {
  ArrowRightStartOnRectangleIcon,
  PencilSquareIcon,
} from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Card, CardBody, CardHeader } from "@heroui/card";
import { Divider } from "@heroui/divider";
import { Link } from "wouter";

export const ProfilePage = () => {
  const profileQuery = useProfileQuery();
  const logout = useLogout();
  const confirm = useConfirm({
    content: "Are you sure you want to logout?",
  });

  if (!profileQuery.isSuccess) return null;

  const { username, email } = profileQuery.data;

  const handleLogoutClick = confirm((closeConfirm) => {
    authService.logout().then(logout).finally(closeConfirm);
  });

  return (
    <section className="h-full flex items-center justify-center p-4">
      <Card className="w-full max-w-[340px] mx-auto">
        <CardHeader className="text-primary-400">
          {username}

          <ul className="ml-auto flex gap-2">
            <li>
              <Link
                asChild
                href={NAVIGATION.editProfile}
                state={{ origin: NAVIGATION.profile }}
              >
                <Button as={Link} isIconOnly size="sm" variant="flat">
                  <PencilSquareIcon className="w-4 h-4" />
                </Button>
              </Link>

              <Link href="/" state={{ origin: NAVIGATION.editProfile }} />
            </li>

            <li>
              <Button
                isIconOnly
                size="sm"
                color="danger"
                title="Logout"
                variant="flat"
                onPress={handleLogoutClick}
              >
                <ArrowRightStartOnRectangleIcon className="w-4 h-4" />
              </Button>
            </li>
          </ul>
        </CardHeader>

        <Divider />

        <CardBody>
          <ul className="text-base">
            <li>{email}</li>
          </ul>
        </CardBody>
      </Card>
    </section>
  );
};
