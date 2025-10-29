import { Button } from "@heroui/button";

type ProfileLoadingErrorMessageProps = {
  isRefetching: boolean;
  onTryAgain: () => void;
};

export const ProfileLoadingErrorMessage = ({
  isRefetching,
  onTryAgain,
}: ProfileLoadingErrorMessageProps) => {
  return (
    <section className="h-full flex items-center justify-center p-4">
      <article className="flex flex-col gap-2">
        <div>
          <h1 className="text-xl">Profile Error</h1>
          <p className="text-large">
            We can't loading your profile. Please, try again
          </p>
        </div>
        <div>
          <Button
            variant="flat"
            color="primary"
            isLoading={isRefetching}
            onPress={onTryAgain}
          >
            Try again
          </Button>
        </div>
      </article>
    </section>
  );
};
