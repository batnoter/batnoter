import DeleteIcon from "@mui/icons-material/Delete";
import { Button, Card, CardActions, CardContent, CardHeader, IconButton, Typography } from "@mui/material";
import React, { ReactElement } from 'react';
import { TreeNode } from "../reducer/noteSlice";
import { getTitleFromFilename } from "../util/util";

interface Props {
  note: TreeNode
  handleEdit: (note: TreeNode) => void
  handleDelete: (note: TreeNode) => void
}

const MAX_CARD_TEXT_LENGTH = 20


const NoteCard: React.FC<Props> = ({ note, handleEdit, handleDelete }): ReactElement => {
  return (
    <Card elevation={1}>
      <CardHeader action={
        <>
          <IconButton onClick={() => handleDelete(note)}> <DeleteIcon /> </IconButton>
        </>
      } title={getTitleFromFilename(note.name)} />
      <CardContent>
        <Typography color="textSecondary"> {note.content?.substring(0, MAX_CARD_TEXT_LENGTH)} {note.content && note.content.length > MAX_CARD_TEXT_LENGTH && '...'}</Typography>
      </CardContent>
      <CardActions>
        <Button onClick={() => handleEdit(note)} size="small"> EDIT</Button>
      </CardActions>
    </Card>
  )
}

export default NoteCard;