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
  variants: flag.variants || [],
  rules: flag.rules || [],
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
  constraints: rule.constraints || [newConstraint()],
  distributions: rule.distributions || [],
});

export const newConstraint = (constraint = {}) => ({
  id: constraint.id || uuid(),
  property: constraint.property || '',
  operation: constraint.operation || OperationTypes.ONE_OF,
  values: constraint.values || [''],
});
