import { gql } from 'apollo-boost';

export const FLAG_QUERY = gql`
    query getFlag($id: ID!) {
        flag(id: $id) {
            id
            key
            name
            description
            enabled
            rules {
                id
                constraints {
                    id
                    property
                    operation
                    values
                }
                distributions {
                    variant {
                        id
                    }
                    percentage
                }
            }
            variants {
                id
                description
                value
            }
            defaultVariantWhenOn {
                id
            }
            defaultVariantWhenOff {
                id
            }
        }
    }
`;

export const OPERATIONS_SEGMENTS_QUERY = gql`
    query operationsAndSegments {
        operations: __type(name: "Operation") {
            enumValues {
                name
            }
        }
        segments {
            id
            name
        }
    }
`;

export const CREATE_FLAG_QUERY = gql`
    mutation createFlag($input: NewFlag!) {
        createFlag(input: $input) {
            id
            key
            name
            enabled
            description
            createdAt
            __typename
        }
    }
`;

export const UPDATE_FLAG_QUERY = gql`
    mutation updateFlag($id: ID!, $input: UpdateFlag!){
        updateFlag(id: $id, input: $input) {
            id
            key
            name
            description
            enabled
        }
    }
`;

export const DELETE_FLAG_QUERY = gql`
    mutation deleteFlag($id: ID!) {
        deleteFlag(id: $id)
    }
`;

export const CREATE_VARIANT_QUERY = gql`
    mutation createVariant($flagId: ID!, $input: NewVariant!) {
        createVariant(flagId: $flagId, input: $input) {
            id
            description
            value
        }
    }
`;

export const UPDATE_VARIANT_QUERY = gql`
    mutation updateVariant($id: ID!, $flagId: ID!, $input: UpdateVariant!) {
        updateVariant(id: $id, flagId: $flagId, input: $input) {
            id
            description
            value
        }
    }
`;

export const DELETE_VARIANT_QUERY = gql`
    mutation deleteVariant($id: ID!, $flagId: ID!) {
        deleteVariant(id: $id, flagId: $flagId)
    }
`;

export const CREATE_FLAG_RULE_QUERY = gql`
    mutation createFlagRule($flagId: ID!, $input: NewFlagRule!) {
        createFlagRule(flagId: $flagId, input: $input)
    }
`;

export const UPDATE_FLAG_RULE_QUERY = gql`
    mutation updateFlagRule($id: ID!, $flagId: ID!, $input: UpdateFlagRule!) {
        updateFlagRule(id: $id, flagId: $flagId, input: $input)
    }
`;

export const DELETE_FLAG_RULE_QUERY = gql`
    mutation deleteFlagRule($id: ID!, $flagId: ID!) {
        deleteFlagRule(id: $id, flagId: $flagId)
    }
`;
