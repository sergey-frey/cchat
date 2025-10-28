import { ISearchUsersResponse } from "@/entities/user";
import { InfiniteData } from "@tanstack/react-query";

export const PLACEHOLDER_USERS: InfiniteData<
  ISearchUsersResponse["data"],
  number
> = {
  pageParams: [],
  pages: [
    {
      id: "d8687f38-0c7b-4afc-8a3c-ea1b2ddee0b0",
      email: "Rosalia_Stroman@gmail.com",
      username: "Guadalupe31",
      name: "Gerardo",
    },
    {
      id: "6c9ff07a-d5f8-4bf2-a99f-1a0fe5e12065",
      email: "Alvah.Rosenbaum@gmail.com",
      username: "Ramon_Johns0",
      name: "Dino",
    },
    {
      id: "b221a43f-21de-4c52-8426-c41e48f90e6c",
      email: "Arnoldo_Grimes@hotmail.com",
      username: "Matt_Feil",
      name: "Leonora",
    },
    {
      id: "5088cd9d-d696-4949-8660-c290af68c949",
      email: "Winifred_Crist30@gmail.com",
      username: "Sedrick32",
      name: "Heidi",
    },
    {
      id: "ad6ccafe-b00f-40c8-9111-393e69587072",
      email: "Domenica.Windler69@yahoo.com",
      username: "Deja49",
      name: "Kaelyn",
    },
    {
      id: "3c95bca8-4296-431d-94fd-f47f86abce33",
      email: "Elissa_McClure91@gmail.com",
      username: "Isac87",
      name: "Rod",
    },
    {
      id: "695bd08c-7321-4905-b4dc-8413990f33af",
      email: "Birdie.Wolff41@gmail.com",
      username: "Beaulah_Mohr20",
      name: "Cornelius",
    },
    {
      id: "50259a5d-36d3-4932-a081-0d952ddd74ec",
      email: "Dejuan69@yahoo.com",
      username: "Karen47",
      name: "Ramiro",
    },
  ] as unknown as InfiniteData<ISearchUsersResponse["data"], number>["pages"],
};
