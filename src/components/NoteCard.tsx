import { DeleteOutlined } from "@mui/icons-material";
import { Card, CardContent, CardHeader, IconButton, Typography } from "@mui/material";
import * as React from 'react';
import { Note } from "../reducer/note/noteSlice";

interface Props {
  note: Note
  handleDelete: (noteId: number) => void
}

const NoteCard: React.FC<Props> = ({ note, handleDelete }) => {
  return (<div>
    <Card elevation={1}>
      <CardHeader action={<IconButton onClick={() => handleDelete(note.id as number)}> <DeleteOutlined /> </IconButton>}
        title={note.title} />
      <CardContent>
        <Typography color="textSecondary"> {note.content} </Typography>
      </CardContent>
    </Card>
  </div>)
}

export default NoteCard;