import { v1 as uuid } from 'uuid';
import { isArray } from 'lodash';
import { Operations, VariantType } from './copy';
import { cast } from '../../helpers';

export const OperationTypes = Object.keys(Operations)
  .reduce((ops, op) => ({ ...ops, [op]: op }), {});

export const VariantTypes = Object.keys(VariantType)
  .reduce((vts, vt) => ({ ...vts, [vt]: vt.toLowerCase() }), {});

export const newFlag = (flag = {}, flagType) => {
  const isNew = !flag.id;
  const variants = isArray(flag.variants) && flag.variants.length > 0 ?
    flag.variants.map(v => newVariant(v)) :
    createNewVariants(flagType);
  return {
    __new: isNew,
    id: flag.id || uuid(),
    name: flag.name || '',
    key: flag.key || '',
    description: flag.description || '',
    variants,
    rules: flag.rules ? flag.rules.map(r => newRule(r)) : [],
    defaultVariantWhenOn: flag.defaultVariantWhenOn || ((isNew && variants[0]) || { id: '' }),
    defaultVariantWhenOff: flag.defaultVariantWhenOff || ((isNew && variants[1]) || { id: '' }),
  }
};

export const newVariant = (variant = {}) => ({
  __new: !variant.id,
  id: variant.id || uuid(),
  description: variant.description || '',
  value: variant.value !== undefined ? variant.value : '',
  type: variant.type !== undefined ? variant.type :
    variant.value !== undefined ? typeof variant.value : VariantTypes.BOOLEAN,
});

export const newRule = (rule = {}) => ({
  __new: !rule.id,
  id: rule.id || uuid(),
  constraints: rule.constraints ?
    rule.constraints.map(c => newConstraint(c)) :
    [newConstraint()],
  distributions: rule.distributions ?
    rule.distributions.map(d => newDistribution(d)) :
    [newDistribution()],
});

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
  percentage: distribution.percentage || 100,
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
  distributions: rule.distributions.map(d => formatDistribution(d, variantsRef)),
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