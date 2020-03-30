import { gql } from 'apollo-boost';

export const FLAGS_QUERY = gql`
  query listFlags($search: String) {
    flags(search: $search) {
      id
      key
      name
      description
      enabled
      createdAt
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
