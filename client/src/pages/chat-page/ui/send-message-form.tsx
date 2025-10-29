import { MessageType } from "@/entities/chats";
import { SendMessageDto } from "@/features/chat";
import { Textarea } from "@heroui/input";
import { ReactNode } from "react";
import { Controller, useForm } from "react-hook-form";

type SendMessageFormProps = {
  onSubmit: (message: SendMessageDto) => Promise<void>;
  submitButton: ReactNode;
};

export const SendMessageForm = ({
  onSubmit,
  submitButton,
}: SendMessageFormProps) => {
  const { handleSubmit, control, reset } = useForm<SendMessageDto>({
    defaultValues: { type: MessageType.TEXT, content: "" },
  });

  return (
    <form
      className="flex gap-2"
      onSubmit={handleSubmit((data) => onSubmit(data).then(() => reset()))}
    >
      <Controller
        name="content"
        control={control}
        render={({ field: { value, onChange } }) => {
          return (
            <Textarea
              minRows={1}
              maxRows={4}
              placeholder="Type a message..."
              value={value}
              onChange={onChange}
            />
          );
        }}
      />

      {submitButton}
    </form>
  );
};
