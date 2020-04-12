import { v1 as uuid } from 'uuid';
import { get, isArray, isNil } from 'lodash';
import { Operations, VariantType } from './copy';
import { cast } from '../../helpers';

export const OperationTypes = Object.keys(Operations)
  .reduce((ops, op) => ({ ...ops, [op]: op }), {});

export const VariantTypes = Object.keys(VariantType)
  .reduce((vts, vt) => ({ ...vts, [vt]: vt.toLowerCase() }), {});

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

export const newConstraint = (constraint = {}) => ({
  __new: !constraint.id,
  id: constraint.id || uuid(),
  property: constraint.property || '',
  operation: constraint.operation || OperationTypes.ONE_OF,
  values: constraint.values || [],
  type: constraint.type !== undefined ? constraint.type :
    constraint.value !== undefined ? typeof constraint.value : VariantTypes.STRING,
});

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
  distributions: rule.distributions.map((d, idx) => {
    if (rule.returnVariant !== PercentageRollout) {
      d.percentage = idx === 0 ? 100 : 0;
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