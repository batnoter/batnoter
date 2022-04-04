import { Masonry } from '@mui/lab';
import { Container } from '@mui/material';
import { useModal } from 'mui-modal-provider';
import React, { ReactElement, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { deleteNoteAsync, getNotesAsync, selectNotesTree, TreeNode, TreeUtil } from '../reducer/noteSlice';
import { getDecodedPath } from '../util/util';
import ConfirmDialog from './ConfirmDialog';
import NoteCard from './NoteCard';

const Finder = (): ReactElement => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { showModal } = useModal();
  const tree = useAppSelector(selectNotesTree);
  const [searchParams] = useSearchParams();
  const path = getDecodedPath(searchParams.get('path'));

  useEffect(() => {
    dispatch(getNotesAsync(path));
  }, [tree, path]);

  const handleDelete = (note: TreeNode) => {
    showModal(ConfirmDialog, {
      desc: 'Are you sure you want to delete this note?',
      onConfirm: () => dispatch(deleteNoteAsync(note))
    });
  }

  const handleEdit = (note: TreeNode) => {
    navigate("/edit?path=" + encodeURIComponent(note.path));
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
      <Masonry columns={{ xs: 1, md: 3, xl: 4 }} spacing={2}>
        {notes.filter(n => !n.is_dir).map(note => (
          <div key={note.path}> <NoteCard note={note} handleEdit={handleEdit} handleDelete={handleDelete} /> </div>
        ))}
      </Masonry>
    </Container>
  );
}

export default Finder;
