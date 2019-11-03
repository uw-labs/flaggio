export const cast = (value, type) => {
  switch (type) {
    case 'string':
      return String(value);
    case 'boolean':
      return Boolean(value);
    case 'number':
      return Number(value);
    default:
      return value;
  }
};

export const inferCast = value => {
  switch (true) {
    case !isNaN(Number(value)):
      return Number(value);
    case value === 'true' || value === 'false':
      return Boolean(value);
    default:
      return value;
  }
};
