import { DeleteOutlined } from "@mui/icons-material";
import { Card, CardContent, CardHeader, IconButton, Typography } from "@mui/material";
import * as React from 'react';
import { Note } from "../reducer/noteSlice";

interface Props {
  note: Note
  handleDelete: (note: Note) => void
}

const NoteCard: React.FC<Props> = ({ note, handleDelete }) => {
  const getFirstLine = (str: string) => {
    return str.split('\n', 1)[0]
  }
  return (<div>
    <Card elevation={1}>
      <CardHeader action={<IconButton onClick={() => handleDelete(note)}> <DeleteOutlined /> </IconButton>}
        title={getFirstLine(note.content)} />
      <CardContent>
        <Typography color="textSecondary"> {note.content} </Typography>
      </CardContent>
    </Card>
  </div>)
}

export default NoteCard;