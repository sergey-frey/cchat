import { EditProfileForm } from "@/features/edit-profile";
import { NAVIGATION } from "@/shared/navigation";
import { NavigationOriginState } from "@/shared/types/navigation";
import { ChevronLeftIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Card, CardBody, CardHeader } from "@heroui/card";
import { Divider } from "@heroui/divider";
import { Link } from "wouter";
import { useHistoryState } from "wouter/use-browser-location";
import { useEditProfileHandlers } from "../model/use-edit-profile-handlers";

export const EditProfilePage = () => {
  const history = useHistoryState<NavigationOriginState>();

  const { handleEditProfileSuccess, handleEditProfileError } =
    useEditProfileHandlers();

  const backHref = history?.origin ?? NAVIGATION.profile;

  return (
    <section className="h-full p-4 flex items-center justify-center">
      <Card className="w-full max-w-[340px]">
        <CardHeader className="flex items-center gap-2">
          <Button
            variant="light"
            as={Link}
            href={backHref}
            size="sm"
            isIconOnly
          >
            <ChevronLeftIcon className="w-4 h-4" />
          </Button>

          <h1>Edit profile</h1>
        </CardHeader>
        <Divider />
        <CardBody>
          <EditProfileForm
            onSuccess={handleEditProfileSuccess}
            onError={handleEditProfileError}
          />
        </CardBody>
      </Card>
    </section>
  );
};
