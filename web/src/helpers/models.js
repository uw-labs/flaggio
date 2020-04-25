import { v1 as uuid } from 'uuid';
import { get, isArray, isNil } from 'lodash';
import { cast } from '.';

const Operations = [
  'ONE_OF', 'NOT_ONE_OF',
  'GREATER', 'GREATER_OR_EQUAL',
  'LOWER', 'LOWER_OR_EQUAL',
  'EXISTS', 'DOESNT_EXIST',
  'CONTAINS', 'DOESNT_CONTAIN',
  'STARTS_WITH', 'DOESNT_START_WITH',
  'ENDS_WITH', 'DOESNT_END_WITH',
  'MATCHES_REGEX', 'DOESNT_MATCH_REGEX',
  'IS_IN_SEGMENT', 'ISNT_IN_SEGMENT',
  'IS_IN_NETWORK',
];
export const OperationTypes = Operations.reduce((ops, op) => (
  { ...ops, [op]: op }
), {});

export const VariantTypes = { BOOLEAN: 'boolean', NUMBER: 'number', STRING: 'string' };

export const PercentageRollout = 'ROLLOUT';

export const newFlag = (flag = {}, flagType) => {
  const isNew = flag.__new || !flag.id;
  const variants = isArray(flag.variants) && flag.variants.length > 0 ?
    flag.variants.map(v => newVariant(v)) :
    createNewVariants(flagType);
  const [var1, var2] = variants;

  return {
    __new: isNew,
    id: flag.id || uuid(),
    name: flag.name || '',
    key: flag.key || '',
    description: flag.description || '',
    variants,
    rules: flag.rules ? flag.rules.map(r => newRule(r, variants)) : [],
    defaultVariantWhenOn: flag.defaultVariantWhenOn || { id: (isNew && var1 && var1.id) || '' },
    defaultVariantWhenOff: flag.defaultVariantWhenOff || { id: (isNew && var2 && var2.id) || '' },
  }
};

export const newVariant = (variant = {}) => ({
  __new: variant.__new || !variant.id,
  id: variant.id || uuid(),
  description: variant.description || '',
  value: variant.value !== undefined ? variant.value : '',
  type: variant.type !== undefined ? variant.type :
    variant.value !== undefined ? typeof variant.value : VariantTypes.BOOLEAN,
});

export const newRule = (rule = {}, variants) => {
  const isNew = rule.__new || !rule.id;

  // ensure there is one distribution for each variant
  const distributions = variants.map(v => {
    const d = rule.distributions && rule.distributions.find(d => d.variant.id === v.id);
    return newDistribution(d ? d : { variant: v, percentage: 0 });
  });

  // check if there is a distribution with 100% chance (means this was the value selected to be returned)
  let returnVariant = get(distributions.find(d => d.percentage === 100), 'variant.id', '');
  // no return variant found, if the flag is not new then the return value must be a percentage rollout
  if (!returnVariant && !isNew) returnVariant = PercentageRollout;

  return {
    __new: isNew,
    id: rule.id || uuid(),
    returnVariant,
    constraints: rule.constraints ?
      rule.constraints.map(c => newConstraint(c)) :
      [newConstraint()],
    distributions,
  }
};

export const newConstraint = (constraint = {}) => {
  const values = constraint.values || [];
  return {
    __new: !constraint.id,
    id: constraint.id || uuid(),
    property: constraint.property || '',
    operation: constraint.operation || OperationTypes.ONE_OF,
    values,
    type: constraint.type !== undefined ? constraint.type :
      values[0] !== undefined ? typeof values[0] : VariantTypes.STRING,
  }
};

export const newDistribution = (distribution = {}) => ({
  __new: !distribution.id,
  id: distribution.id || uuid(),
  variant: distribution.variant || { id: '' },
  percentage: isNil(distribution.percentage) ? 100 : distribution.percentage,
});

export const formatFlag = (flag, variantsRef) => ({
  key: flag.key,
  name: flag.name,
  description: flag.description,
  defaultVariantWhenOn: variantsRef[flag.defaultVariantWhenOn.id] || flag.defaultVariantWhenOn.id,
  defaultVariantWhenOff: variantsRef[flag.defaultVariantWhenOff.id] || flag.defaultVariantWhenOff.id,
});

export const formatVariant = variant => ({
  description: variant.description,
  value: cast(variant.value, variant.type),
});

export const formatRule = (rule, variantsRef) => ({
  constraints: rule.constraints.map(c => formatConstraint(c)),
  distributions: rule.distributions.map(d => {
    if (rule.returnVariant !== PercentageRollout) {
      d.percentage = d.variant.id === rule.returnVariant ? 100 : 0;
    }
    return formatDistribution(d, variantsRef);
  }),
});

export const formatConstraint = constraint => ({
  property: constraint.property,
  operation: constraint.operation,
  values: isArray(constraint.values) && constraint.values.length > 0 ?
    constraint.values.map(v => cast(v, constraint.type)) : [],
});

export const formatDistribution = (distribution, variantsRef) => ({
  variantId: variantsRef[distribution.variant.id] || distribution.variant.id,
  percentage: distribution.percentage,
});

export const newSegment = (segment = {}) => ({
  __new: !segment.id,
  id: segment.id || uuid(),
  name: segment.name || '',
  description: segment.description || '',
  rules: segment.rules ? segment.rules.map(r => newSegmentRule(r)) : [],
});

export const newSegmentRule = (rule = {}) => ({
  __new: !rule.id,
  id: rule.id || uuid(),
  constraints: rule.constraints ?
    rule.constraints.map(c => newConstraint(c)) :
    [newConstraint()],
});

export const formatSegment = segment => ({
  key: segment.key,
  name: segment.name,
  description: segment.description,
});

export const formatSegmentRule = rule => ({
  constraints: rule.constraints.map(c => formatConstraint(c)),
});

const createNewVariants = (type) => {
  switch (type) {
    case VariantTypes.NUMBER:
      return [
        newVariant({ value: 1, type: VariantTypes.NUMBER }),
        newVariant({ value: 2, type: VariantTypes.NUMBER }),
      ];
    case VariantTypes.STRING:
      return [
        newVariant({ value: 'a', type: VariantTypes.STRING }),
        newVariant({ value: 'b', type: VariantTypes.STRING }),
      ];
    case VariantTypes.BOOLEAN:
    default:
      return [
        newVariant({ value: true, type: VariantTypes.BOOLEAN }),
        newVariant({ value: false, type: VariantTypes.BOOLEAN }),
      ];
  }
};