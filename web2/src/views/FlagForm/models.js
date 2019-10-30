import uuid from 'uuid/v1';
import { Operations, VariantType } from './copy';

export const OperationTypes = Object.keys(Operations)
  .reduce((ops, op) => ({ ...ops, [op]: op }), {});

export const VariantTypes = Object.keys(VariantType)
  .reduce((vts, vt) => ({ ...vts, [vt]: vt.toLowerCase() }), {});

export const newFlag = (flag = {}) => ({
  id: flag.id || uuid(),
  name: flag.name || '',
  key: flag.key || '',
  description: flag.description || '',
  variants: flag.variants ? flag.variants.map(v => newVariant(v)) : [],
  rules: flag.rules ? flag.rules.map(r => newRule(r)) : [],
});

export const newVariant = (variant = {}) => ({
  id: variant.id || uuid(),
  description: variant.description || '',
  value: variant.value !== undefined ? variant.value : '',
  type: variant.type !== undefined ? variant.type :
    variant.value !== undefined ? typeof variant.value : VariantTypes.STRING,
  defaultWhenOn: false,
  defaultWhenOff: false,
});

export const newRule = (rule = {}) => ({
  id: rule.id || uuid(),
  constraints: rule.constraints ?
    rule.constraints.map(c => newConstraint(c)) :
    [newConstraint()],
  distributions: rule.distributions ?
    rule.distributions.map(d => newDistribution(d)) :
    [newDistribution()],
});

export const newConstraint = (constraint = {}) => ({
  id: constraint.id || uuid(),
  property: constraint.property || '',
  operation: constraint.operation || OperationTypes.ONE_OF,
  values: constraint.values || [''],
});

export const newDistribution = (distribution = {}) => ({
  id: distribution.id || uuid(),
  variant: distribution.variant || { id: '' },
  percentage: distribution.percentage || 100,
});