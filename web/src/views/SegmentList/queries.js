import { gql } from 'apollo-boost';

export const SEGMENTS_QUERY = gql`
    query listSegments {
        segments {
            id
            name
            createdAt
        }
    }
`;
