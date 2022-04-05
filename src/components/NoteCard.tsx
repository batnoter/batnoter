import DeleteIcon from "@mui/icons-material/Delete";
import { Button, Card, CardActions, CardContent, CardHeader, IconButton } from "@mui/material";
import React, { ReactElement } from 'react';
import { TreeNode } from "../reducer/noteSlice";
import { getTitleFromFilename } from "../util/util";
import CustomReactMarkdown from "./lib/CustomReactMarkdown";

interface Props {
  note: TreeNode
  handleView: (note: TreeNode) => void
  handleEdit: (note: TreeNode) => void
  handleDelete: (note: TreeNode) => void
}

const MAX_CARD_TEXT_LENGTH = 300;

const NoteCard: React.FC<Props> = ({ note, handleView, handleEdit, handleDelete }): ReactElement => {
  const getCardText = (text?: string): string => {
    if (text == null) return '';
    return text.substring(0, MAX_CARD_TEXT_LENGTH) + (text.length > MAX_CARD_TEXT_LENGTH ? '...' : '');
  }
  return (
    <Card elevation={1}>
      <CardHeader action={
        <>
          <IconButton sx={{ "&:hover": { color: "red" } }} onClick={() => handleDelete(note)}> <DeleteIcon /> </IconButton>
        </>
      } title={getTitleFromFilename(note.name)} />
      <CardContent>
        <CustomReactMarkdown className='custom-html-style'>{getCardText(note.content)}</CustomReactMarkdown>
      </CardContent>
      <CardActions>
        <Button onClick={() => handleView(note)} size="small">VIEW</Button>
        <Button onClick={() => handleEdit(note)} size="small">EDIT</Button>
      </CardActions>
    </Card>
  )
}

export default NoteCard;