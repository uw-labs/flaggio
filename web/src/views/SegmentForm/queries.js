import { gql } from 'apollo-boost';

export const SEGMENT_QUERY = gql`
    query getSegment($id: ID!) {
        segment(id: $id) {
            id
            name
            description
            rules {
                id
                constraints {
                    id
                    property
                    operation
                    values
                }
            }
        }
    }
`;

export const OPERATIONS_QUERY = gql`
    query operations {
        operations: __type(name: "Operation") {
            enumValues {
                name
            }
        }
    }
`;

export const CREATE_SEGMENT_QUERY = gql`
    mutation createSegment($input: NewSegment!) {
        createSegment(input: $input) {
            id
            name
            description
            rules {
                id
                constraints {
                    id
                    operation
                    property
                    values
                }
            }
            createdAt
        }
    }
`;

export const UPDATE_SEGMENT_QUERY = gql`
    mutation updateSegment($id: ID!, $input: UpdateSegment!){
        updateSegment(id: $id, input: $input) {
            id
            name
            description
            rules {
                id
                constraints {
                    id
                    operation
                    property
                    values
                }
            }
            createdAt
        }
    }
`;

export const DELETE_SEGMENT_QUERY = gql`
    mutation deleteSegment($id: ID!) {
        deleteSegment(id: $id)
    }
`;

export const CREATE_SEGMENT_RULE_QUERY = gql`
    mutation createSegmentRule($segmentId: ID!, $input: NewSegmentRule!) {
        createSegmentRule(segmentId: $segmentId, input: $input) {
            id
            constraints {
                id
                operation
                property
                values
            }
        }
    }
`;

export const UPDATE_SEGMENT_RULE_QUERY = gql`
    mutation updateSegmentRule($id: ID!, $segmentId: ID!, $input: UpdateSegmentRule!) {
        updateSegmentRule(id: $id, segmentId: $segmentId, input: $input) {
            id
            constraints {
                id
                operation
                property
                values
            }
        }
    }
`;

export const DELETE_SEGMENT_RULE_QUERY = gql`
    mutation deleteSegmentRule($id: ID!, $segmentId: ID!) {
        deleteSegmentRule(id: $id, segmentId: $segmentId)
    }
`;
