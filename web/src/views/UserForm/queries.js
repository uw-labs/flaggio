import { gql } from 'apollo-boost';

export const USER_QUERY = gql`
    query getUser($id: ID!, $search: String, $offset: Int, $limit: Int) {
        user(id: $id) {
            id
            context
            updatedAt
            evaluations(search: $search, offset: $offset, limit: $limit) {
                evaluations {
                    id
                    flagId
                    flagKey
                    flagVersion
                    value
                    createdAt
                }
                total
            }
        }
    }
`;

export const DELETE_USER_QUERY = gql`
    mutation deleteUser($id: ID!) {
        deleteUser(id: $id)
    }
`;

export const DELETE_EVALUATION_QUERY = gql`
    mutation deleteEvaluation($id: ID!) {
        deleteEvaluation(id: $id)
    }
`;
