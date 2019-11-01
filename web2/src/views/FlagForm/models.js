import uuid from 'uuid/v1';
import { Operations, VariantType } from './copy';

export const OperationTypes = Object.keys(Operations)
  .reduce((ops, op) => ({ ...ops, [op]: op }), {});

export const VariantTypes = Object.keys(VariantType)
  .reduce((vts, vt) => ({ ...vts, [vt]: vt.toLowerCase() }), {});

export const newFlag = (flag = {}) => ({
  __new: !flag.id,
  id: flag.id || uuid(),
  name: flag.name || '',
  key: flag.key || '',
  description: flag.description || '',
  variants: flag.variants ? flag.variants.map(v => newVariant(v)) : [],
  rules: flag.rules ? flag.rules.map(r => newRule(r)) : [],
  defaultVariantWhenOn: flag.defaultVariantWhenOn,
  defaultVariantWhenOff: flag.defaultVariantWhenOff,
});

export const newVariant = (variant = {}) => ({
  __new: !variant.id,
  id: variant.id || uuid(),
  description: variant.description || '',
  value: variant.value !== undefined ? variant.value : '',
  type: variant.type !== undefined ? variant.type :
    variant.value !== undefined ? typeof variant.value : VariantTypes.STRING,
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
  values: constraint.values || [''],
});

export const newDistribution = (distribution = {}) => ({
  __new: !distribution.id,
  id: distribution.id || uuid(),
  variant: distribution.variant || { id: '' },
  percentage: distribution.percentage || 100,
});

export const formatFlag = flag => ({
  key: flag.key,
  name: flag.name,
  description: flag.description,
});

export const formatVariant = variant => ({
  description: variant.description,
  value: cast(variant.value, variant.type),
});

export const formatRule = rule => ({
  constraints: rule.constraints.map(c => formatConstraint(c)),
  distributions: rule.distributions.map(d => formatDistribution(d)),
});

export const formatConstraint = constraint => ({
  property: constraint.property,
  operation: constraint.operation,
  values: constraint.values,
});

export const formatDistribution = distribution => ({
  variantId: distribution.variant.id,
  percentage: distribution.percentage,
});

const cast = (value, type) => {
  switch (type) {
    case VariantTypes.STRING:
      return String(value);
    case VariantTypes.BOOLEAN:
      return Boolean(value);
    case VariantTypes.NUMBER:
      return Number(value);
  }
  return value;
};
