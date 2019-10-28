import { gql } from 'apollo-boost';

export const FLAG_QUERY = gql`
    query getFlag($id: ID!){
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
                defaultWhenOn
                defaultWhenOff
            }
        }
        segments {
            id
            key
            name
        }
        operations: __type(name: "Operation") {
            enumValues {
                name
            }
        }
    }
`;

export const DELETE_FLAG_QUERY = gql`
    mutation deleteFlag($id: ID!){
        deleteFlag(id: $id)
    }
`;
