export const cast = (value, type) => {
  switch (type) {
    case 'boolean':
      return Boolean(value);
    case 'number':
      return Number(value);
    case 'string':
      return String(value);
    default:
      return value;
  }
};
