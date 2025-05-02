import { IUser } from "@/entities/user";
import { faker } from "@faker-js/faker";

export const generateRandomUser = (): IUser => {
  return {
    id: faker.string.uuid(),
    email: faker.internet.email(),
    name: faker.person.fullName(),
    username: faker.internet.username(),
  };
};
