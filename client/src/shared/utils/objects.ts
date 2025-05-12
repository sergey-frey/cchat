export const keys = <T extends Record<string | number | symbol, unknown>>(
  obj: T,
) => Object.keys(obj) as Array<keyof T>;
