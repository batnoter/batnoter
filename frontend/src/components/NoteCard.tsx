import { DeleteOutlined } from "@mui/icons-material";
import { Card, CardContent, CardHeader, IconButton, Typography } from "@mui/material";
import React, { ReactElement } from 'react';
import { TreeNode } from "../reducer/noteSlice";

interface Props {
  note: TreeNode
  handleDelete: (note: TreeNode) => void
}

const MaxCardTextLength = 20

const NoteCard: React.FC<Props> = ({ note, handleDelete }): ReactElement => {
  return (
    <Card elevation={1}>
      <CardHeader action={<IconButton onClick={() => handleDelete(note)}> <DeleteOutlined /> </IconButton>} title={note.name.replace(/(\.md)$/i, '')} />
      <CardContent>
        <Typography color="textSecondary"> {note.content?.substring(0, MaxCardTextLength)} {note.content && note.content.length > MaxCardTextLength && '...'}</Typography>
      </CardContent>
    </Card>
  )
}

export default NoteCard;