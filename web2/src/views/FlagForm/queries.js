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
            name
        }
        operations: __type(name: "Operation") {
            enumValues {
                name
            }
        }
    }
`;

export const CREATE_FLAG_QUERY = gql`
    mutation createFlag($input: NewFlag!){
        createFlag(input: $input) {
            id
        }
    }
`;

export const UPDATE_FLAG_QUERY = gql`
    mutation updateFlag($id: ID!, $input: UpdateFlag!){
        updateFlag(id: $id, input: $input) {
            id
        }
    }
`;

export const DELETE_FLAG_QUERY = gql`
    mutation deleteFlag($id: ID!){
        deleteFlag(id: $id)
    }
`;
