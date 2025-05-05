import { Button } from "@heroui/button";
import { useConfirmStore } from "../model/confirm-store";
import {
  confirmContentSelector,
  confirmSelector,
  setConfirmSelector,
} from "../model/confirm-store-selectors";
import {
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
} from "@heroui/modal";

export const Confirm = () => {
  const confirmFn = useConfirmStore(confirmSelector);
  const confirmContent = useConfirmStore(confirmContentSelector);
  const setConfirm = useConfirmStore(setConfirmSelector);

  const handleCloseConfirm = () => setConfirm(null);

  if (!confirmFn) return null;

  return (
    <Modal
      isOpen={Boolean(confirmFn)}
      onClose={handleCloseConfirm}
      size="xs"
      backdrop="blur"
      hideCloseButton
    >
      <ModalContent>
        <ModalHeader>Confirm</ModalHeader>

        <ModalBody>{confirmContent}</ModalBody>

        <ModalFooter className="grid grid-cols-2">
          <Button color="primary" onPress={confirmFn}>
            Approve
          </Button>
          <Button onPress={handleCloseConfirm}>Close</Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};
