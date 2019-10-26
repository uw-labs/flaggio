export const Operations = {
  ONE_OF: "Equals",
  NOT_ONE_OF: "Not equals",
  GREATER: "Greater",
  GREATER_OR_EQUAL: "Greater or equal",
  LOWER: "Lower",
  LOWER_OR_EQUAL: "Lower or equal",
  EXISTS: "Exists",
  DOESNT_EXIST: "Doesn't exist",
  CONTAINS: "Contains",
  DOESNT_CONTAIN: "Doesn't Contain",
  STARTS_WITH: "Starts with",
  DOESNT_START_WITH: "Doesn't start with",
  ENDS_WITH: "Ends with",
  DOESNT_END_WITH: "Doesn't end with",
  MATCHES_REGEX: "Matches regex",
  DOESNT_MATCH_REGEX: "Doesn't match regex",
  BEFORE_DATE: "Before date",
  BEFORE_OR_SAME_DATE: "Before or same date",
  AFTER_DATE: "After date",
  AFTER_OR_SAME_DATE: "After or same date",
  IS_IN_SEGMENT: "Is in segment",
  ISNT_IN_SEGMENT: "Isn't in segment",
};

export const OperationTypes = Object.keys(Operations)
  .reduce((ops, op) => ({...ops, [op]: op}), {});

export const VariantType = {
  BOOLEAN: "Boolean",
  NUMBER: "Number",
  STRING: "String",
};

export const VariantTypes = Object.keys(VariantType)
  .reduce((vts, vt) => ({...vts, [vt]: vt.toLowerCase()}), {});

export const BooleanType = {
  [true]: "True",
  [false]: "False",
};
