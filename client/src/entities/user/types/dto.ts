export type UpdateUserDto = {
  email?: string;
  name?: string;
  new_password?: string;
  previous_password?: string;
  username?: string;
};

export type SearchUsersDto = {
  username: string;
  limit?: number;
  pagination?: number;
};
