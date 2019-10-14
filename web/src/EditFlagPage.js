import React from 'react';
import {
  useParams
} from "react-router-dom";

function EditFlagPage() {
  let { id } = useParams();
  return (
    <div>
      Edit Flag
    </div>
  )
}

export default EditFlagPage;