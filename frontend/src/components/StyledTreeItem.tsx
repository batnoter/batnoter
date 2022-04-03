import styled from '@emotion/styled';
import AddBoxIcon from '@mui/icons-material/AddBox';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import { TreeItem, TreeItemProps } from '@mui/lab';
import { IconButton } from '@mui/material';
import React, { SyntheticEvent } from 'react';

type StyledTreeItemProps = TreeItemProps & {
  isDir: boolean,
  handleCreate: (e: SyntheticEvent, dirPath: string) => void,
  handleEdit: (e: SyntheticEvent, filepath: string) => void,
  handleDelete: (e: SyntheticEvent, filepath: string) => void
}

const StyledTreeItem = styled((props: StyledTreeItemProps) => {
  const { isDir, handleCreate, handleEdit, handleDelete, ...otherProps } = props;
  return (
    <TreeItem  {...otherProps} label={
      <>{otherProps.label + ' '}
        {isDir ? <IconButton size="small" onClick={(e) => handleCreate(e, otherProps.nodeId)}><AddBoxIcon sx={{ display: 'none', verticalAlign: 'text-bottom' }} fontSize='inherit' /> </IconButton> :
          <>
            <IconButton size="small" onClick={(e) => handleEdit(e, otherProps.nodeId)}><EditIcon sx={{ display: 'none', verticalAlign: 'text-bottom' }} fontSize='inherit' /></IconButton>
            <IconButton size="small" onClick={(e) => handleDelete(e, otherProps.nodeId)}><DeleteIcon className="delete" sx={{ display: 'none', verticalAlign: 'text-bottom' }} fontSize='inherit' /></IconButton>
          </>
        }
      </>
    } />
  )
})(() => ({
  [`& .MuiTreeItem-content .MuiTreeItem-label`]: {
    '&:hover svg': {
      display: 'inline-block',
    },
    '& svg:hover': {
      color: 'blue'
    },
    '& svg.delete:hover': {
      color: 'red'
    }
  },
}));

export default StyledTreeItem;