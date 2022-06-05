import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import FolderIcon from '@mui/icons-material/Folder';
import NotesIcon from '@mui/icons-material/Notes';
import { Alert, Box, Breadcrumbs, Button, CircularProgress, Container, Divider, Grid, Link } from "@mui/material";
import { unwrapResult } from '@reduxjs/toolkit';
import { useModal } from 'mui-modal-provider';
import React, { ReactElement, useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../app/hooks";
import { APIStatus, APIStatusType } from '../reducer/common';
import { deleteNoteAsync, getNoteAsync, resetStatus, selectNoteAPIStatus, selectNotesTree, TreeNode } from "../reducer/noteSlice";
import TreeUtil from '../util/TreeUtil';
import { confirmDeleteNote, getDecodedPath, getSanitizedErrorMessage, getTitleFromFilename, splitPath, URL_ISSUES } from "../util/util";
import CustomReactMarkdown from './lib/CustomReactMarkdown';

const isLoading = (apiStatus: APIStatus): boolean => {
  const { getNoteAsync, deleteNoteAsync } = apiStatus;
  return getNoteAsync === APIStatusType.LOADING || deleteNoteAsync === APIStatusType.LOADING;
}

const isGetNoteLoading = (apiStatus: APIStatus): boolean => {
  const { getNoteAsync } = apiStatus;
  return getNoteAsync === APIStatusType.LOADING;
}

const isFailed = (apiStatus: APIStatus): boolean => {
  const { getNoteAsync, deleteNoteAsync } = apiStatus;
  return getNoteAsync === APIStatusType.FAIL || deleteNoteAsync === APIStatusType.FAIL;
}

const Viewer: React.FC = (): ReactElement => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { showModal } = useModal();

  const [note, setNote] = useState<TreeNode>()
  const [searchParams] = useSearchParams();
  const path = getDecodedPath(searchParams.get('path'));
  const tree = useAppSelector(selectNotesTree);
  const apiStatus = useAppSelector(selectNoteAPIStatus);
  const [errorMessage, setErrorMessage] = React.useState("");
  const dirPathArray = splitPath(path);
  const title = getTitleFromFilename(dirPathArray.pop() || '');

  const handleDelete = () => {
    confirmDeleteNote(showModal, () => {
      dispatch(deleteNoteAsync(note as TreeNode)).then(unwrapResult)
        .then(() => navigate(`/?path=${encodeURIComponent(dirPathArray.join('/'))}`))
        .catch(err => setErrorMessage(getSanitizedErrorMessage(err)));
    });
  }

  useEffect(() => {
    // This should be the first useEffect hook. Declare other useEffect hooks below this one.
    dispatch(resetStatus());
  }, [path])

  useEffect(() => {
    const treeNode = TreeUtil.searchNode(tree, path);
    if (treeNode == null || treeNode.is_dir) {
      return;
    }
    dispatch(getNoteAsync(treeNode.path)).then(unwrapResult)
      .catch(err => setErrorMessage(getSanitizedErrorMessage(err)));
    setNote(treeNode);
  }, [tree, path])

  return (
    <Container maxWidth="lg">{isGetNoteLoading(apiStatus) ? <CircularProgress sx={{ position: "relative", top: "50%", left: "50%" }} /> :
      <Box>
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
            <Button onClick={() => navigate(`/edit?path=${encodeURIComponent(note?.path || '')}`)} disabled={isLoading(apiStatus)} variant="contained" sx={{ mx: 2 }} startIcon={<EditIcon />}>EDIT</Button>
            <Button onClick={() => handleDelete()} disabled={isLoading(apiStatus)} variant="contained" startIcon={<DeleteIcon />} color="error">DELETE</Button>
          </Box>
        </Grid>
        <Divider sx={{ my: 3 }} />
        {isFailed(apiStatus) && errorMessage && <Alert severity="error" sx={{ width: "100%", mb: 2 }}>{errorMessage} <span>please try again or <Link href={URL_ISSUES} target="_blank" rel="noopener">create an issue</Link></span></Alert>}
        <Box className='viewer-markdown' sx={{ p: 2 }}>
          <CustomReactMarkdown className='custom-html-style'>{note?.content || ''}</CustomReactMarkdown>
        </Box>
      </Box>
    }
    </Container>
  );
}

export default Viewer;
