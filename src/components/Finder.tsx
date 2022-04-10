import { Masonry } from '@mui/lab';
import { Alert, CircularProgress, Container } from '@mui/material';
import { unwrapResult } from '@reduxjs/toolkit';
import { useModal } from 'mui-modal-provider';
import React, { ReactElement, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { APIStatus, APIStatusType } from '../reducer/common';
import { deleteNoteAsync, getNotesAsync, resetStatus, selectNoteStatus, selectNotesTree, TreeNode, TreeUtil } from '../reducer/noteSlice';
import { confirmDeleteNote, getDecodedPath, getSanitizedErrorMessage } from '../util/util';
import NoteCard from './NoteCard';

const isGetNotesLoading = (apiStatus: APIStatus): boolean => {
  const { getNotesAsync } = apiStatus;
  return getNotesAsync === APIStatusType.LOADING;
}

const isGetNotesFailed = (apiStatus: APIStatus): boolean => {
  const { getNotesAsync } = apiStatus;
  return getNotesAsync === APIStatusType.FAIL;
}

const Finder = (): ReactElement => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { showModal } = useModal();
  const tree = useAppSelector(selectNotesTree);
  const apiStatus = useAppSelector(selectNoteStatus);
  const [searchParams] = useSearchParams();
  const path = getDecodedPath(searchParams.get('path'));
  const [errorMessage, setErrorMessage] = React.useState("");

  useEffect(() => {
    // This should be the first useEffect hook. Declare other useEffect hooks below this one.
    dispatch(resetStatus());
  }, [path])

  useEffect(() => {
    dispatch(getNotesAsync(path)).then(unwrapResult)
      .catch(err => setErrorMessage(getSanitizedErrorMessage(err)));
  }, [tree, path])

  const handleDelete = (note: TreeNode) => {
    confirmDeleteNote(showModal, () => dispatch(deleteNoteAsync(note as TreeNode)));
  }

  const handleView = (note: TreeNode) => {
    navigate(`/view?path=${encodeURIComponent(note.path)}`);
  }

  const handleEdit = (note: TreeNode) => {
    navigate(`/edit?path=${encodeURIComponent(note.path)}`);
  }

  const getChildren = (path: string): TreeNode[] | undefined => {
    const node = TreeUtil.searchNode(tree, path);
    if (node?.cached) {
      return node.children;
    }
  }

  const notes = getChildren(path) || [] as TreeNode[];

  return (
    <Container>
      {isGetNotesFailed(apiStatus) && errorMessage && <Alert severity="error" sx={{ width: "100%", mb: 2 }}>{errorMessage}</Alert>}

      <Masonry columns={{ xs: 1, md: 3, xl: 4 }} spacing={2}>
        {isGetNotesLoading(apiStatus) ? <CircularProgress sx={{ position: "relative", top: "50%", left: "50%" }} /> :

          notes.filter(n => !n.is_dir).map(note => (
            <div key={note.path}> <NoteCard note={note} handleView={handleView} handleEdit={handleEdit} handleDelete={handleDelete} /> </div>
          ))}
      </Masonry>
    </Container>
  );
}

export default Finder;
