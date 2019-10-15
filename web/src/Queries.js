import {gql} from "apollo-boost";

export const FLAGS_QUERY = gql`
    {
        flags {
            id
            key
            name
            enabled
            createdAt
        }
    }
`;

export const TOGGLE_FLAG_QUERY = gql`
    mutation($id: ID!, $input: UpdateFlag!) {
        updateFlag(id: $id, input: $input) {
            id
            enabled
        }
    }
`;

export const CREATE_FLAG_QUERY = gql`
    mutation ($name: String!, $key: String!, $description:String){
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