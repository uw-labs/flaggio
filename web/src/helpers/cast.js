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

export const inferCast = value => {
  switch (true) {
    case typeof value !== 'string':
      return value;
    case value === 'true':
      return true;
    case value === 'false':
      return false;
    case !isNaN(Number(value)):
      return Number(value);
    default:
      return value;
  }
};
