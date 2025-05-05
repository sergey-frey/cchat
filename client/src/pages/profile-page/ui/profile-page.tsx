import { useProfileQuery } from "@/entities/user";
import { authService, useLogout } from "@/features/auth";
import { NAVIGATION } from "@/shared/navigation";
import { cn } from "@/shared/utils/cn";
import { useConfirm } from "@/shared/utils/confirm";
import {
  ArrowRightStartOnRectangleIcon,
  PencilSquareIcon,
} from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Card, CardBody, CardHeader } from "@heroui/card";
import { Divider } from "@heroui/divider";
import { Link } from "wouter";
import { ProfileLoadingErrorMessage } from "./profile-loading-error-message";

export const ProfilePage = () => {
  const profileQuery = useProfileQuery();
  const logout = useLogout();
  const confirm = useConfirm({
    content: "Are you sure you want to logout?",
  });

  if (!profileQuery.isSuccess) {
    return (
      <ProfileLoadingErrorMessage
        onTryAgain={profileQuery.refetch}
        isRefetching={profileQuery.isRefetching}
      />
    );
  }

  const { username, email, name } = profileQuery.data;

  const handleLogoutClick = confirm((closeConfirm) => {
    authService.logout().then(logout).finally(closeConfirm);
  });

  return (
    <section className="h-full flex items-center justify-center p-4">
      <div
        className={cn("w-full max-w-[340px] mx-auto", "flex flex-col gap-2")}
      >
        <Card>
          <CardHeader as={"h1"} className="text-primary-400">
            @{username}
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
            <ul className={cn("text-base", "flex flex-col gap-2")}>
              <li>
                <h2 className="text-small font-semibold">Name:</h2>
                <div>{name}</div>
              </li>

              <li>
                <h2 className="text-small font-semibold">Email:</h2>
                <div>{email}</div>
              </li>
            </ul>
          </CardBody>
        </Card>

        <Card>
          <CardBody>
            <div className={cn("flex gap-1", "text-xs text-foreground-500")}>
              Registered at: <time>{new Date().toLocaleDateString()}</time>
            </div>
          </CardBody>
        </Card>
      </div>
    </section>
  );
};
