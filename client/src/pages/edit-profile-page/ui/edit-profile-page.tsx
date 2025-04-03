import { EditProfileForm } from "@/features/edit-profile";
import { Card, CardBody, CardHeader } from "@heroui/card";
import { Divider } from "@heroui/divider";

export const EditProfilePage = () => {
  return (
    <section className="h-full p-4 flex items-center justify-center">
      <Card className="w-full max-w-[340px]">
        <CardHeader>
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
