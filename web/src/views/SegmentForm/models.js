import uuid from 'uuid/v1';
import { isArray } from 'lodash';
import { Operations } from './copy';
import { VariantTypes } from '../FlagForm/models';
import { cast } from '../../helpers';

export const OperationTypes = Object.keys(Operations)
  .reduce((ops, op) => ({ ...ops, [op]: op }), {});

export const newSegment = (segment = {}) => ({
  __new: !segment.id,
  id: segment.id || uuid(),
  name: segment.name || '',
  description: segment.description || '',
  rules: segment.rules ? segment.rules.map(r => newRule(r)) : [],
});

export const newRule = (rule = {}) => ({
  __new: !rule.id,
  id: rule.id || uuid(),
  constraints: rule.constraints ?
    rule.constraints.map(c => newConstraint(c)) :
    [newConstraint()],
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

export const formatSegment = segment => ({
  key: segment.key,
  name: segment.name,
  description: segment.description,
});

export const formatRule = rule => ({
  constraints: rule.constraints.map(c => formatConstraint(c)),
});

export const formatConstraint = constraint => ({
  property: constraint.property,
  operation: constraint.operation,
  values: isArray(constraint.values) && constraint.values.length > 0 ?
    constraint.values.map(v => cast(v, constraint.type)) : [],
});
