import { gql } from 'apollo-boost';

export const FLAGS_QUERY = gql`
  query listFlags($search: String, $offset: Int, $limit: Int) {
    flags(search: $search, offset: $offset, limit: $limit) {
      flags {
        id
        key
        name
        description
        enabled
        createdAt
      }
      total
    }
  }
`;

export const TOGGLE_FLAG_QUERY = gql`
  mutation toggleFlag($id: ID!, $input: UpdateFlag!) {
    updateFlag(id: $id, input: $input) {
      id
      enabled
    }
  }
`;
