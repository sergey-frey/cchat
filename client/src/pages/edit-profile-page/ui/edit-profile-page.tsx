import { EditProfileForm } from "@/features/edit-profile";
import { NAVIGATION } from "@/shared/navigation";
import { ChevronLeftIcon } from "@heroicons/react/24/outline";
import { Button } from "@heroui/button";
import { Card, CardBody, CardHeader } from "@heroui/card";
import { Divider } from "@heroui/divider";
import { Link } from "wouter";
import { useHistoryState } from "wouter/use-browser-location";

export const EditProfilePage = () => {
  const history = useHistoryState<{ origin?: string }>();

  const backHref = history?.origin ?? NAVIGATION.profile;

  return (
    <section className="h-full p-4 flex items-center justify-center">
      <Card className="w-full max-w-[340px]">
        <CardHeader className="flex items-center gap-2">
          <Button variant="light" as={Link} href={backHref} size="sm" isIconOnly>
            <ChevronLeftIcon className="w-4 h-4" />
          </Button>

          <h1>Edit profile</h1>
        </CardHeader>
        <Divider />
        <CardBody>
          <EditProfileForm />
        </CardBody>
      </Card>
    </section>
  );
};
