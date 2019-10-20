import {gql} from "apollo-boost";

export const FLAGS_QUERY = gql`
    query listFlags {
        flags {
            id
            key
            name
            enabled
            createdAt
        }
    }
`;

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
                key
                description
                value
                defaultWhenOn
                defaultWhenOff
            }
        }
        operations: __type(name: "Operation") {
            enumValues {
                name
            }
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

export const CREATE_FLAG_QUERY = gql`
    mutation createFlag($name: String!, $key: String!, $description:String){
        createFlag(input: {
            name: $name
            key: $key
            description: $description
        }) {
            id
            key
            name
            enabled
            createdAt
        }
    }
`;

export const DELETE_FLAG_QUERY = gql`
    mutation deleteFlag($id: ID!){
        deleteFlag(id: $id)
    }
`;