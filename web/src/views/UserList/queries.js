import { gql } from 'apollo-boost';

export const USERS_QUERY = gql`
  query listUsers($search: String, $offset: Int, $limit: Int) {
    users(search: $search, offset: $offset, limit: $limit) {
      users {
        id
        context
        updatedAt
      }
      total
    }
  }
`;
