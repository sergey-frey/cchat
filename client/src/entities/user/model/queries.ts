import { userApi } from "@/shared/api/instance/instance";
import { IUser } from "@/shared/api/types";
import { useQuery } from "@tanstack/react-query";

export const useUsersQuery = (name?: IUser["name"]) => {
  return useQuery<IUser[]>({
    queryKey: ["users"],
    queryFn: () => {
      return userApi.get(`/${name ?? ""}`).json();
      // return [
      //   {
      //     username: "username_1",
      //     name: "User1",
      //   },
      //   {
      //     username: "username_2",
      //     name: "User2",
      //   },
      // ].map((user, i) => ({ id: i.toString(), ...user }));
    },
  });
};
