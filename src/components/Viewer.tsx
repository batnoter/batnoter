import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import FolderIcon from '@mui/icons-material/Folder';
import NotesIcon from '@mui/icons-material/Notes';
import { Box, Breadcrumbs, Button, Container, Divider, Grid, Link } from "@mui/material";
import { useModal } from 'mui-modal-provider';
import React, { ReactElement, useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../app/hooks";
import { deleteNoteAsync, getNoteAsync, selectNotesTree, TreeNode, TreeUtil } from "../reducer/noteSlice";
import { confirmDeleteNote, getDecodedPath, getTitleFromFilename, splitPath } from "../util/util";
import CustomReactMarkdown from './lib/CustomReactMarkdown';
import './Viewer.scss';

const Viewer: React.FC = (): ReactElement => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { showModal } = useModal();
  const [note, setNote] = useState<TreeNode>()
  const [searchParams] = useSearchParams();
  const path = getDecodedPath(searchParams.get('path'));
  const tree = useAppSelector(selectNotesTree);
  const dirPathArray = splitPath(path);
  const title = getTitleFromFilename(dirPathArray.pop() || '');

  const handleDelete = () => {
    confirmDeleteNote(showModal, () => {
      dispatch(deleteNoteAsync(note as TreeNode));
      navigate(`/?path=${encodeURIComponent(dirPathArray.join('/'))}`);
    });
  }

  useEffect(() => {
    const treeNode = TreeUtil.searchNode(tree, path);
    if (treeNode == null || treeNode.is_dir) {
      return;
    }
    if (!treeNode.cached) {
      dispatch(getNoteAsync(treeNode.path));
      return;
    }
    setNote(treeNode);
  }, [tree, path]);

  return (
    <Container maxWidth="lg">
      <Grid container direction="row" justifyContent="space-between" alignItems="center">
        <Box>
          <Breadcrumbs itemsAfterCollapse={2} sx={{ fontSize: '1.2rem' }}>
            <Link key="root" underline="hover" color="inherit"><FolderIcon fontSize="medium" sx={{ mr: 0.5, verticalAlign: 'middle', }} />root</Link>
            {dirPathArray.map((option) => (<Link key={option} underline="hover" color="inherit"> {option} </Link>))}
          </Breadcrumbs>
          <NotesIcon color="inherit" fontSize="medium" sx={{ mr: 0.5, verticalAlign: 'middle', }} />{title}
        </Box>
        <Box>
          <Button onClick={() => navigate('/')} variant="outlined" startIcon={<ArrowBackIcon />}>BACK</Button>
          <Button onClick={() => navigate(`/edit?path=${encodeURIComponent(note?.path || '')}`)} variant="contained" sx={{ mx: 2 }} startIcon={<EditIcon />}>EDIT</Button>
          <Button onClick={() => handleDelete()} variant="contained" startIcon={<DeleteIcon />} color="error">DELETE</Button>
        </Box>
      </Grid>
      <Divider sx={{ my: 3 }} />
      <Box className='viewer-markdown' sx={{ background: 'white', p: 2 }}>
        <CustomReactMarkdown className='custom-html-style'>{note?.content || ''}</CustomReactMarkdown>
      </Box>
    </Container>
  );
}

export default Viewer;
